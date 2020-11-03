package conf

import (
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
	"sigs.k8s.io/yaml"

	"arhat.dev/helm-stack/pkg/constant"
)

type Environment struct {
	Name        string           `json:"name" yaml:"name"`
	KubeContext string           `json:"kubeContext" yaml:"kubeContext"`
	Deployments []DeploymentSpec `json:"deployments" yaml:"deployments"`
}

func (e Environment) ValuesDir(envDir string) string {
	return filepath.Join(envDir, e.Name)
}

func (e Environment) ManifestsDir(envDir string) string {
	return filepath.Join(e.ValuesDir(envDir), "manifests")
}

func (e Environment) CustomManifestsDir(envDir string, dep *DeploymentSpec) string {
	valuesDir := e.ValuesDir(envDir)

	if dep == nil {
		return filepath.Join(valuesDir, "manifests-custom")
	}

	namespace, name := dep.NamespaceAndName()
	return filepath.Join(valuesDir, "manifests-custom", namespace+"."+name)
}

func (e Environment) Validate(charts map[string]*ChartSpec) error {
	var err error
	if e.Name == "" {
		err = multierr.Append(err, fmt.Errorf("invalid empty deployment environment name"))
	}

	names := make(map[string]struct{})
	for _, d := range e.Deployments {
		if _, defined := names[d.Name]; defined {
			err = multierr.Append(err, fmt.Errorf("duplicate deployment item %q", d.Name))
		}

		err = multierr.Append(err, d.Validate(charts))
	}

	return nil
}

func (e Environment) Ensure(
	ctx context.Context,
	chartsDir, localChartsDir, envDir string,
	charts map[string]*ChartSpec,
) error {
	_ = ctx
	valuesDir := e.ValuesDir(envDir)

	err := os.MkdirAll(valuesDir, 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("failed to ensure deployment environment values dir %q: %w", valuesDir, err)
	}

	for i, d := range e.Deployments {
		chart := charts[d.Chart]
		if chart == nil {
			return fmt.Errorf("chart %s not found", d.Chart)
		}

		subChartNames, err := chart.SubChartNames(chartsDir, localChartsDir)
		if err != nil {
			return fmt.Errorf("failed to check sub charts")
		}

		cDir := e.CustomManifestsDir(envDir, &e.Deployments[i])
		err = os.MkdirAll(cDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create custom manifests dir %q: %w", cDir, err)
		}

		for _, subChartName := range append([]string{""}, subChartNames...) {
			destValuesFile := filepath.Join(valuesDir, d.Filename(subChartName))

			_, err := os.Stat(destValuesFile)
			if err == nil {
				// values file exists
				continue
			}

			if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("failed to probe values file %q: %w", destValuesFile, err)
			}

			baseValuesFile := d.BaseValues
			if baseValuesFile == "" {
				baseValuesFile = constant.DefaultValuesFile
			}

			srcValuesFile := filepath.Join(chart.Dir(chartsDir, localChartsDir, subChartName), baseValuesFile)
			if err = iohelper.CopyFile(srcValuesFile, destValuesFile); err == nil {
				continue
			}

			if subChartName != "" {
				// this is a sub chart, may not contain target values file
				subChartValuesFiles := []string{
					filepath.Join(chart.Dir(chartsDir, localChartsDir, subChartName), baseValuesFile),
				}
				if baseValuesFile != constant.DefaultValuesFile {
					// fallback to values.yaml
					subChartValuesFiles = append(subChartValuesFiles,
						filepath.Join(chart.Dir(chartsDir, localChartsDir, subChartName), constant.DefaultValuesFile),
					)
				}
				for _, srcValuesFile := range subChartValuesFiles {
					err = iohelper.CopyFile(srcValuesFile, destValuesFile)
					if os.IsNotExist(err) {
						// values file not found in sub chart, just ignore it
						err = nil
					}
				}
			}

			if err != nil {
				return fmt.Errorf("failed to copy values file %q: %w", srcValuesFile, err)
			}
		}
	}

	return nil
}

// nolint:gocyclo
func (e Environment) Gen(
	ctx context.Context,
	chartsDir, localChartsDir, envDir string,
	charts map[string]*ChartSpec,
) error {
	manifestsDir := e.ManifestsDir(envDir)

	_ = os.RemoveAll(manifestsDir)

	err := os.MkdirAll(manifestsDir, 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to ensure manifests dir")
	}

	for _, d := range e.Deployments {
		chart := charts[d.Chart]
		if chart == nil {
			return fmt.Errorf("chart %s not found", d.Chart)
		}

		subChartNames, err := chart.SubChartNames(chartsDir, localChartsDir)
		if err != nil {
			return fmt.Errorf("failed to check sub charts")
		}

		var (
			namespace, name = d.NamespaceAndName()
			manifestFile    = filepath.Join(manifestsDir, d.Filename(""))
			baseValuesFile  = d.BaseValues
		)

		if baseValuesFile == "" {
			baseValuesFile = constant.DefaultValuesFile
		}

		cmd := []string{
			"helm", "template", "--namespace", namespace, "--debug",
			"--values", filepath.Join(chart.Dir(chartsDir, localChartsDir, ""), baseValuesFile),
			"--set", "fullnameOverride=" + name,
		}

		if !isHelmV2() {
			cmd = append(cmd, "--no-hooks")
			if d.ExcludeChartCRDs {
				cmd = append(cmd, "--skip-crds")
			} else {
				cmd = append(cmd, "--include-crds")
			}
			cmd = append(cmd, name, chart.Dir(chartsDir, localChartsDir, ""))
		} else {
			cmd = append(cmd, chart.Dir(chartsDir, localChartsDir, ""), name)
		}

		allValues := make(map[string]interface{})
		for _, subChartName := range append([]string{""}, subChartNames...) {
			var (
				valuesFile = filepath.Join(e.ValuesDir(envDir), d.Filename(subChartName))
			)

			currentValues := map[string]interface{}{}

			subChartValuesFiles := []string{
				filepath.Join(chart.Dir(chartsDir, localChartsDir, subChartName), baseValuesFile),
			}
			if subChartName != "" && baseValuesFile != constant.DefaultValuesFile {
				// fallback to values.yaml
				subChartValuesFiles = append(subChartValuesFiles,
					filepath.Join(chart.Dir(chartsDir, localChartsDir, subChartName), constant.DefaultValuesFile),
				)
			}

			data, fErr := ioutil.ReadFile(valuesFile)
			if fErr != nil {
				if subChartName == "" || !os.IsNotExist(fErr) {
					return fmt.Errorf("failed to read values from file %q: %w", valuesFile, fErr)
				}

				// some sub chart may not contain any values file, check if has values in its dir
				for _, f := range subChartValuesFiles {
					_, fErr = os.Stat(f)
					if fErr == nil {
						return fmt.Errorf("inconsistent values file, please run `helm-stack ensure` to fix it")
					}
				}

				continue
			}

			if mErr := yaml.Unmarshal(data, &currentValues); err != nil {
				return fmt.Errorf("failed to parse values from file %q: %w", valuesFile, mErr)
			}

			if subChartName != "" {
				// get sub chart base values
				for _, subChartBaseValuesFile := range subChartValuesFiles {
					data, fErr = ioutil.ReadFile(subChartBaseValuesFile)
					// ignore this error
					if fErr != nil {
						if !os.IsNotExist(fErr) {
							return fmt.Errorf(
								"failed to check sub chart base values %q: %w",
								subChartBaseValuesFile, fErr,
							)
						}

						continue
					}

					subChartBaseValues := make(map[string]interface{})
					if fErr == nil {
						if mErr := yaml.Unmarshal(data, &subChartBaseValues); mErr != nil {
							return fmt.Errorf(
								"failed to parse sub chart base values from file %q: %w",
								valuesFile, mErr,
							)
						}

						switch t := allValues[subChartName].(type) {
						case map[string]interface{}:
							allValues[subChartName] = mergeMaps(t, subChartBaseValues)
						case nil:
							allValues[subChartName] = subChartBaseValues
						default:
							return fmt.Errorf("invalid sub chart values in main values file: %v", t)
						}
					}
				}
			}

			if subChartName == "" {
				allValues = mergeMaps(allValues, currentValues)
			} else {
				// merge sub chart values
				switch t := allValues[subChartName].(type) {
				case map[string]interface{}:
					allValues[subChartName] = mergeMaps(t, currentValues)
				case nil:
					allValues[subChartName] = currentValues
				default:
					return fmt.Errorf("invalid sub chart values in main values file: %v", t)
				}
			}
		}

		valuesBytes, mErr := yaml.Marshal(allValues)
		if mErr != nil {
			return fmt.Errorf("failed to marshal values: %w", mErr)
		}

		tempValuesFile, fErr := ioutil.TempFile(os.TempDir(), "helm-stack-values-*.yaml")
		if fErr != nil {
			return fmt.Errorf("failed to create temporary values file: %w", fErr)
		}

		err = func() error {
			defer func() {
				// best effort
				_ = os.Remove(tempValuesFile.Name())
			}()

			err = func() error {
				defer func() {
					_ = tempValuesFile.Close()
				}()

				_, err2 := tempValuesFile.Write(valuesBytes)
				if err2 != nil {
					return fmt.Errorf("failed to write temporary values file %q: %w", tempValuesFile.Name(), err2)
				}

				return nil
			}()
			if err != nil {
				return err
			}

			cmd = assembleCommandWithoutEmptyString(cmd, "--values", tempValuesFile.Name())

			fmt.Println("Executing:", strings.Join(cmd, " "))
			manifestFile, err2 := os.OpenFile(manifestFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err2 != nil {
				return fmt.Errorf("failed to open manifest file: %w", err2)
			}
			defer func() { _ = manifestFile.Close() }()

			_, err = exechelper.Do(exechelper.Spec{
				Context: ctx,
				Command: cmd,
				Stdout:  manifestFile,
				Stderr:  os.Stderr,
			})
			if err != nil {
				return fmt.Errorf("failed to generate manifest: %w", err)
			}

			return nil
		}()

		if err != nil {
			return err
		}
	}

	return nil
}

func (e Environment) Apply(ctx context.Context, dryRunArg, envDir string, charts map[string]*ChartSpec) error {
	kubectlCmd := []string{"kubectl"}

	if e.KubeContext != "" {
		kubectlCmd = append(kubectlCmd, "--context", e.KubeContext)
	}

	var failedCustomMainfests []string
	for i, d := range e.Deployments {
		var action []string

		s := d.GetState()

		switch {
		case s.Present && !s.CRDPresent:
			// TODO: filter out crd
			fallthrough
		case s.Present && s.CRDPresent:
			action = []string{"apply"}
		case !s.Present && !s.CRDPresent:
			// TODO: filter out crd
			fallthrough
		case !s.Present && s.CRDPresent:
			action = []string{"delete", "--ignore-not-found=true"}
		}

		chart := charts[d.Chart]
		if chart == nil {
			return fmt.Errorf("chart %s not found", d.Chart)
		}

		namespace, _ := d.NamespaceAndName()

		// ensure namespace (best effort)
		nsCreateCmd := assembleCommandWithoutEmptyString(kubectlCmd, "create", dryRunArg, "namespace", namespace)
		fmt.Println("Executing:", strings.Join(nsCreateCmd, " "))
		_, _ = exechelper.Do(exechelper.Spec{
			Context: ctx,
			Command: nsCreateCmd,
			Stdout:  ioutil.Discard,
			Stderr:  ioutil.Discard,
		})

		manifestFile := filepath.Join(e.ManifestsDir(envDir), d.Filename(""))

		applyCmd := assembleCommandWithoutEmptyString(kubectlCmd,
			append(action, dryRunArg, "--filename", manifestFile)...)
		if !chart.NamespaceInTemplate {
			applyCmd = append(applyCmd, "--namespace", namespace)
		}

		fmt.Println("Executing:", strings.Join(applyCmd, " "))
		_, err := exechelper.Do(exechelper.Spec{
			Context: ctx,
			Command: applyCmd,
			Stdout:  os.Stdout,
			Stderr:  os.Stderr,
		})
		if err != nil {
			return fmt.Errorf("failed to execute kubectl apply: %w", err)
		}

		cDir := e.CustomManifestsDir(envDir, &e.Deployments[i])
		_, err = os.Stat(cDir)
		if err != nil {
			continue
		}

		files, err := ioutil.ReadDir(cDir)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to check custom manifests: %w", err)
		}

		if len(files) == 0 {
			continue
		}

		customApply := assembleCommandWithoutEmptyString(kubectlCmd,
			append(action, dryRunArg, "--recursive", "--filename", cDir)...)
		fmt.Println("Executing:", strings.Join(customApply, " "))
		_, err = exechelper.Do(exechelper.Spec{
			Context: ctx,
			Command: customApply,
			Stdout:  os.Stdout,
			Stderr:  os.Stderr,
		})
		if err != nil {
			if s.Present {
				// failed to apply (error!)
				return fmt.Errorf("failed to apply custom manifests: %w", err)
			}

			// failed to delete (no big deal)
			failedCustomMainfests = append(failedCustomMainfests, fmt.Sprintf("%s: %v", cDir, err))
		}
	}

	if len(failedCustomMainfests) > 0 {
		_, _ = fmt.Fprintf(os.Stderr,
			"failed to delete following custom manifests, you may want to delete them manually:\n%s\n",
			strings.Join(failedCustomMainfests, "\n"))
	}

	return nil
}

type DeploymentSpec struct {
	// Name value for <namespace>/<fullnameOverride>
	Name string `json:"name" yaml:"name"`
	// Chart name listed in the configuration
	Chart string `json:"chart" yaml:"chart"`

	State string `json:"state" yaml:"state"`

	// BaseValues the values file name
	BaseValues string `json:"baseValues" yaml:"baseValues"`

	// ExcludeChartCRDs to apply crds dir in chart
	ExcludeChartCRDs bool `json:"excludeChartCRDs" yaml:"excludeChartCRDs"`
}

func (c DeploymentSpec) Filename(subChart string) string {
	repoName, chartName, chartVersion := getChartRepoNameChartNameChartVersion(c.Chart)
	var (
		namespace, name = c.NamespaceAndName()
	)

	if subChart != "" {
		subChart = "_" + subChart
	}
	if repoName == "" {
		return fmt.Sprintf("%s.%s[%s@%s].yaml", namespace, name, chartName+subChart, chartVersion)
	}

	return fmt.Sprintf("%s.%s[%s.%s@%s].yaml", namespace, name, repoName, chartName+subChart, chartVersion)
}

func (c DeploymentSpec) NamespaceAndName() (namespace, name string) {
	parts := strings.SplitN(c.Name, "/", 2)
	if len(parts) != 2 {
		panic("invalid deployment name without namespace")
	}

	return parts[0], parts[1]
}

func (c DeploymentSpec) Validate(charts map[string]*ChartSpec) error {
	var err error
	if !strings.Contains(c.Name, "/") {
		err = multierr.Append(err, fmt.Errorf("invalid deployment name without namespace"))
	}

	if c.Chart == "" {
		err = multierr.Append(err, fmt.Errorf("invalid empty chart name for deployment"))
	} else if _, ok := charts[c.Chart]; !ok {
		err = multierr.Append(err, fmt.Errorf("deployment chart %q not listed", c.Chart))
	}

	return multierr.Append(err, c.GetState().Validate())
}

func (c DeploymentSpec) GetState() DeploymentState {
	ret := new(DeploymentState)
	ret.Present = true
	ret.CRDPresent = true

	for _, s := range strings.Split(c.State, ",") {
		switch strings.ToLower(s) {
		case "present", "":
			ret.Present = true
		case "absent":
			ret.Present = false
		case "crds":
			ret.CRDPresent = true
		case "nocrds":
			ret.CRDPresent = false
		default:
			if ret.UnknownStates == "" {
				ret.UnknownStates = s
			} else {
				ret.UnknownStates = ret.UnknownStates + "," + s
			}
		}
	}

	return *ret
}

type DeploymentState struct {
	Present    bool
	CRDPresent bool

	UnknownStates string
}

func (s DeploymentState) Validate() error {
	var err error
	if s.UnknownStates != "" {
		err = multierr.Append(err, fmt.Errorf("coantains unknown states %q", s.UnknownStates))
	}

	return err
}
