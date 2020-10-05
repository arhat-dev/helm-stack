package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"arhat.dev/helm-stack/pkg/conf"
	"arhat.dev/helm-stack/pkg/constant"
)

func NewCleanCommand(appCtx *context.Context) *cobra.Command {
	var (
		noAsk bool
	)

	cmd := &cobra.Command{
		Use:           "clean",
		Short:         "clean directories and files for charts and environments",
		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: func(cmd *cobra.Command, args []string) error {
			config := (*appCtx).Value(constant.ContextKeyConfig).(*conf.ResolvedConfig)
			return runClean(*appCtx, config, noAsk)
		},
	}

	fs := cmd.Flags()

	fs.BoolVarP(&noAsk, "yes", "y", false, "remove without asking")

	return cmd
}

func runClean(ctx context.Context, config *conf.ResolvedConfig, noAsk bool) error {
	_ = ctx
	s := bufio.NewScanner(os.Stdin)
	s.Split(bufio.ScanLines)

	filesToRemove, err := collectChartsToRemove(config)
	if err != nil {
		return fmt.Errorf("failed to collect charts to be removed: %w", err)
	}

	if len(filesToRemove) != 0 {
		fmt.Printf("will remove following file or directories for charts:\n  - %s\n",
			strings.Join(filesToRemove, "\n  - "))
		if noAsk || getUserSelection(s) {
			removeAll(filesToRemove)
		}
	} else {
		fmt.Println("charts dir clean")
	}

	for _, e := range config.Environments {
		filesToRemove, err = collectEnvironmentFilesToRemove(config, e)
		if err != nil {
			return fmt.Errorf("failed to collect environment files to be removed: %w", err)
		}

		if len(filesToRemove) != 0 {
			fmt.Printf("will remove following file or directories for environment %q:\n  - %s\n",
				e.Name, strings.Join(filesToRemove, "\n  - "))
			if noAsk || getUserSelection(s) {
				removeAll(filesToRemove)
			}
		} else {
			fmt.Printf("environment %q clean\n", e.Name)
		}
	}

	var (
		envsToRemove []string
	)
	envDirs, err := ioutil.ReadDir(config.App.EnvironmentsDir)
	if err != nil {
		return fmt.Errorf("failed to inspect environment directories: %w", err)
	}

	for _, d := range envDirs {
		if _, ok := config.Environments[d.Name()]; !ok {
			envsToRemove = append(envsToRemove, filepath.Join(config.App.EnvironmentsDir, d.Name()))
		}
	}

	if len(envsToRemove) != 0 {
		fmt.Printf("will remove following environment(s):\n  - %s\n", strings.Join(envsToRemove, "\n  - "))
		if noAsk || getUserSelection(s) {
			removeAll(envsToRemove)
		}
	} else {
		fmt.Println("environments clean")
	}

	return nil
}

func collectChartsToRemove(config *conf.ResolvedConfig) ([]string, error) {
	var (
		chartDirWanted = make(map[string]struct{})
		filesToRemove  []string
	)
	for _, c := range config.Charts {
		chartDirWanted[c.Dir(config.App.ChartsDir, "")] = struct{}{}
	}

	err := filepath.Walk(config.App.ChartsDir, func(path string, info os.FileInfo, e error) error {
		path = filepath.Clean(path)
		if path == filepath.Clean(config.App.ChartsDir) {
			return nil
		}

		if !info.IsDir() {
			filesToRemove = append(filesToRemove, path)
			return nil
		}

		chartDirs, err := ioutil.ReadDir(path)
		if err != nil {
			return fmt.Errorf("failed to read charts dir %q", path)
		}

		for _, d := range chartDirs {
			chartDir := filepath.Join(path, d.Name())
			if _, ok := chartDirWanted[chartDir]; !ok {
				filesToRemove = append(filesToRemove, chartDir)
			}
		}

		return filepath.SkipDir
	})

	if err != nil {
		return nil, fmt.Errorf("failed to inspect all charts: %w", err)
	}

	return filesToRemove, nil
}

func collectEnvironmentFilesToRemove(config *conf.ResolvedConfig, e *conf.Environment) ([]string, error) {
	var (
		filesToRemove    []string
		valuesFileWanted = make(map[string]struct{})
	)

	for _, d := range e.Deployments {
		chart := config.Charts[d.Chart]
		if chart == nil {
			return nil, fmt.Errorf("chart %q does not exists", d.Chart)
		}

		subChartNames, err := chart.SubChartNames(config.App.ChartsDir)
		if err != nil {
			return nil, fmt.Errorf("failed to check sub charts")
		}

		for _, subChartName := range append([]string{""}, subChartNames...) {
			f := filepath.Join(e.ValuesDir(config.App.EnvironmentsDir), d.Filename(subChartName))
			valuesFileWanted[f] = struct{}{}
		}
	}

	valuesDir := e.ValuesDir(config.App.EnvironmentsDir)

	filesInValuesDir, err := ioutil.ReadDir(valuesDir)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to inspect environment values dir %q: %w", valuesDir, err)
	}

	for _, f := range filesInValuesDir {
		path := filepath.Join(valuesDir, f.Name())

		if f.IsDir() {
			// do not remove directories
			continue
		}

		if _, ok := valuesFileWanted[path]; !ok {
			filesToRemove = append(filesToRemove, path)
		}
	}

	return filesToRemove, nil
}

func getUserSelection(s *bufio.Scanner) bool {
	fmt.Print("proceed? (y/N): ")
	for s.Scan() {
		switch strings.ToLower(s.Text()) {
		case "y", "yes":
			return true
		case "n", "no", "":
			return false
		default:
			fmt.Print("please answer y(es) or n(o): ")
		}
	}

	return false
}

func removeAll(files []string) {
	for _, d := range files {
		err := os.RemoveAll(d)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to remove %q: %v", d, err)
		}
	}
}
