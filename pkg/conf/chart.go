package conf

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"arhat.dev/pkg/exechelper"
	"arhat.dev/pkg/iohelper"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/util/json"
)

type ChartSpec struct {
	Name string `json:"name" yaml:"name"`

	// Source of the chart, if nil, the Name MUST be in <repo-name>/<chart-name> format
	*ChartSource `json:",inline" yaml:",inline"`

	// NamespaceInTemplate means there is already proper namespace information defined in its templates
	// and apply with `kubectl --namespace` will fail (mostly for rbac resources)
	NamespaceInTemplate bool `json:"namespaceInTemplate" yaml:"namespaceInTemplate"`
}

func (c ChartSpec) SubChartNames(chartsDir, localChartsDir string) ([]string, error) {
	var baseDir string
	switch {
	case c.ChartSource != nil && c.ChartSource.Local != nil:
		baseDir = localChartsDir
	default:
		baseDir = chartsDir
	}

	files, err := ioutil.ReadDir(filepath.Join(c.Dir(baseDir, localChartsDir, ""), "charts"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("failed to check sub charts")
	}

	result := []string{""}
	for _, f := range files {
		if f.IsDir() {
			result = append(result, f.Name())
		}
	}

	return result, nil
}

func (c ChartSpec) Dir(chartsDir, localChartsDir string, subChartName string) string {
	var baseDir string
	switch {
	case c.ChartSource != nil && c.ChartSource.Local != nil:
		baseDir = localChartsDir
	default:
		baseDir = chartsDir
	}

	repoName, chartName, chartVersion := getChartRepoNameChartNameChartVersion(c.Name)

	var dir string
	if repoName != "" {
		dir = filepath.Join(baseDir, fmt.Sprintf("%s_%s", repoName, chartName))
	} else {
		dir = filepath.Join(baseDir, chartName)
	}

	result := filepath.Join(dir, chartVersion)
	if subChartName == "" {
		return result
	}

	return filepath.Join(result, "charts", subChartName)
}

func (c ChartSpec) Validate(repos map[string]*RepoSpec) error {
	var err error
	if !strings.Contains(c.Name, "@") {
		err = multierr.Append(err, fmt.Errorf("chart name must inclue version info"))
	}

	if c.ChartSource == nil || (c.ChartSource.Git == nil && c.ChartSource.Local == nil) {
		// no custom source (git/local), using repo
		if !strings.Contains(c.Name, "/") {
			err = multierr.Append(err, fmt.Errorf("invalid chart without repo or custom source"))
		} else {
			parts := strings.SplitN(c.Name, "/", 2)
			if _, ok := repos[parts[0]]; !ok {
				err = multierr.Append(err, fmt.Errorf("repo %q for chart %q not found", parts[0], c.Name))
			}
		}
	} else {
		switch {
		case c.ChartSource.Git != nil:
			err = multierr.Append(err, c.ChartSource.Git.Validate())
		case c.ChartSource.Local != nil:
			err = multierr.Append(err, c.ChartSource.Local.Validate())
		}
	}

	return err
}

// nolint:gocyclo
func (c ChartSpec) Ensure(
	ctx context.Context,
	forcePull bool,
	chartsDir, localChartsDir string,
	repos map[string]*RepoSpec,
) error {
	targetDir := c.Dir(chartsDir, localChartsDir, "")

	_, err := os.Stat(targetDir)
	if err == nil && !forcePull {
		return nil
	}

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to probe chart file %q: %w", targetDir, err)
	}

	err = os.MkdirAll(filepath.Dir(targetDir), 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("failed to ensure chart dir %q: %w", targetDir, err)
	}

	// not found, add it
	repoName, chartName, chartVersion := getChartRepoNameChartNameChartVersion(c.Name)

	tmpDir, err := ioutil.TempDir(os.TempDir(), "helm-stack-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary for chart fetch: %w", err)
	}

	defer func() { _ = os.RemoveAll(tmpDir) }()

	switch {
	case c.ChartSource == nil || (c.ChartSource.Git == nil && c.ChartSource.Local == nil):
		// use helm repo
		repo := repos[repoName]
		if repo == nil {
			return fmt.Errorf("repo %q for chart %q not found", repoName, c.Name)
		}

		fetchCmd := []string{"helm", "fetch", "--repo", repo.URL, "--untar", "--untardir", tmpDir}

		var repoExtraArgs []string
		if u := repo.Auth.HTTPBasic.Username; u != "" {
			repoExtraArgs = append(repoExtraArgs, "--username", u)
		}

		if p := repo.Auth.HTTPBasic.Password; p != "" {
			repoExtraArgs = append(repoExtraArgs, "--password", p)
		}

		if repo.TLS.InsecureSkipVerify {
			repoExtraArgs = append(repoExtraArgs, "--insecure-skip-tls-verify")
		}
		if repo.TLS.CaCert != "" {
			repoExtraArgs = append(repoExtraArgs, "--ca-file", repo.TLS.CaCert)
		}

		if repo.TLS.Cert != "" {
			repoExtraArgs = append(repoExtraArgs, "--cert-file", repo.TLS.Cert)
		}

		if repo.TLS.Key != "" {
			repoExtraArgs = append(repoExtraArgs, "--key-file", repo.TLS.Key)
		}

		switch chartVersion {
		case "devel", "latest":
			// add a repo with random name (do not mess with existing repos)
			fakeRepoName := generateRandomName(repo.Name)

			_, err := exechelper.Do(exechelper.Spec{
				Context: ctx,
				Command: append([]string{
					"helm", "repo", "add", "--no-update",
					fakeRepoName, repo.URL}, repoExtraArgs...),
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			})
			if err != nil {
				return fmt.Errorf("failed to add temporary repo %q for version probe: %w", repo.Name, err)
			}

			clearTempRepo := true
			// repo added, ensure we delete it at last
			defer func() {
				if !clearTempRepo {
					return
				}

				_, rmErr := exechelper.Do(exechelper.Spec{
					Context: ctx,
					Command: []string{"helm", "repo", "remove", fakeRepoName},
					Stdout:  os.Stdout,
					Stderr:  os.Stderr,
				})
				if rmErr != nil {
					_, _ = fmt.Fprintf(os.Stderr,
						"failed to delete temporary repo %q after version probe, clean it yourself: %v",
						fakeRepoName, rmErr,
					)
				}
			}()

			_, err = exechelper.Do(exechelper.Spec{
				Context: ctx,
				Command: []string{"helm", "repo", "update"},
				Stdout:  os.Stdout,
				Stderr:  os.Stderr,
			})
			if err != nil {
				return fmt.Errorf("failed to update helm repos: %w", err)
			}

			// check helm version, search command is not compatible between helm2 and helm3
			// currently helm3 will ignore --client flag so it's fine
			// default to helm3
			searchCmd := []string{"helm", "search", "repo"}
			if isHelmV2() {
				searchCmd = []string{"helm", "search"}
			}
			searchCmd = append(searchCmd, "--output", "json")

			if chartVersion == "devel" {
				searchCmd = append(searchCmd, "--devel")
			}

			buf := new(bytes.Buffer)
			_, err = exechelper.Do(exechelper.Spec{
				Context: ctx,
				Command: append(searchCmd, filepath.Join(fakeRepoName, chartName)),
				Stdout:  buf,
				Stderr:  buf,
			})
			if err != nil {
				return fmt.Errorf("failed to get latest chart version of %q: %w", c.Name, err)
			}

			clearTempRepo = false
			_, err = exechelper.Do(exechelper.Spec{
				Context: ctx,
				Command: []string{"helm", "repo", "remove", fakeRepoName},
				Stdout:  os.Stdout,
				Stderr:  os.Stderr,
			})
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr,
					"failed to delete temporary repo %q after version probe, clean it yourself: %v", fakeRepoName, err)
			}

			var data []map[string]interface{}
			err = json.Unmarshal(buf.Bytes(), &data)
			if err != nil {
				return fmt.Errorf("failed to unmarshal chart %q version info: %w", c.Name, err)
			}

			if len(data) == 0 {
				return fmt.Errorf("unable to determin chart %q version", c.Name)
			}

			if len(data[0]) == 0 {
				return fmt.Errorf("unable to parse chart %q version info", c.Name)
			}

			fetchVersion, ok := data[0]["version"].(string)
			if !ok {
				fetchVersion, ok = data[0]["Version"].(string)
			}

			if !ok {
				return fmt.Errorf("unable to get version info for chart %q from json search result", c.Name)
			}

			fetchCmd = append(fetchCmd, "--version", fetchVersion)
		default:
			fetchCmd = append(fetchCmd, "--version", chartVersion)
		}

		exitCode, err := exechelper.Do(exechelper.Spec{
			Context: ctx,
			Command: append(fetchCmd, append(repoExtraArgs, chartName)...),
			Stdout:  os.Stdout,
			Stderr:  os.Stderr,
		})
		if err != nil {
			return fmt.Errorf("helm fetch command exited with code %d: %w", exitCode, err)
		}

		if forcePull {
			err = os.RemoveAll(targetDir)
			if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove existing chart dir %q: %w", targetDir, err)
			}
		}

		// rename fetched chart to version only
		err = iohelper.CopyDir(filepath.Join(tmpDir, chartName), targetDir)
		if err != nil {
			return fmt.Errorf("failed to move fetched chart dir %q: %w", tmpDir, err)
		}

		return nil
	case c.ChartSource.Git != nil:
		config := c.ChartSource.Git
		_, err := exechelper.Do(exechelper.Spec{
			Context: ctx,
			Command: []string{"git", "clone", "--branch", chartVersion, "--depth", "1", config.URL, tmpDir},
			Stdout:  os.Stdout,
			Stderr:  os.Stderr,
		})
		if err != nil {
			return fmt.Errorf("failed to clone repo %q: %w", config.URL, err)
		}

		if forcePull {
			err = os.RemoveAll(targetDir)
			if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove existing chart dir %q: %w", targetDir, err)
			}
		}

		err = iohelper.CopyDir(filepath.Join(tmpDir, config.Path), targetDir)
		if err != nil {
			return fmt.Errorf("failed to move cloned repo: %w", err)
		}

		return nil
	case c.ChartSource.Local != nil:
		// check if chart exists
		chartDir := c.Dir(chartsDir, localChartsDir, "")
		_, err := os.Stat(chartDir)
		if err != nil {
			return fmt.Errorf("failed to check local chart: %w", err)
		}
		return nil
	}

	return nil
}

type ChartSource struct {
	Git   *ChartFromGitRepo   `json:"git" yaml:"git"`
	Local *ChartFromLocalPath `json:"local" yaml:"local"`
}

type ChartFromGitRepo struct {
	// URL for repo (git/http)
	URL string `json:"url" yaml:"url"`
	// Path in the repo
	Path string `json:"path" yaml:"path"`
}

func (g *ChartFromGitRepo) Validate() error {
	var err error
	// u, err := url.Parse(g.URL)
	// if u == nil && err != nil {
	// 	return err
	// }
	// switch u.Scheme {
	// case "git", "https", "http":
	// default:
	// 	err = multierr.Append(err, fmt.Errorf("unsupported git scheme %q, only git/http/https allowed", u.Scheme))
	// }

	return err
}

type ChartFromLocalPath struct {
}

func (l *ChartFromLocalPath) Validate() error {
	return nil
}
