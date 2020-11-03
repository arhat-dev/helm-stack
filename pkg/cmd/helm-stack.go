/*
Copyright 2020 The arhat.dev Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/yaml"

	"arhat.dev/helm-stack/pkg/conf"
	"arhat.dev/helm-stack/pkg/constant"
)

func NewHelmStackCmd() *cobra.Command {
	var (
		appCtx      context.Context
		configFiles []string

		config = conf.NewEmptyResolvedConfig()
	)

	cmd := &cobra.Command{
		Use:           "helm-stack",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Use == "version" {
				return nil
			}

			ctx, exit := context.WithCancel(context.Background())
			sigCh := make(chan os.Signal)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
			go func() {
				i := 0
				for range sigCh {
					i++

					switch i {
					case 1:
						exit()
					default:
						os.Exit(1)
					}
				}
			}()

			userDefinedConfigs := cmd.Parent().PersistentFlags().Lookup("config").Changed

			for _, confFile := range configFiles {
				err := readConfigAndResolve(confFile, config)
				if err != nil {
					if !errors.Is(err, os.ErrNotExist) {
						return err
					}

					if userDefinedConfigs {
						return fmt.Errorf("failed to reolve config: %w", err)
					}
				}
			}

			if config.App.ChartsDir == "" {
				return fmt.Errorf("charts dir must not be empty")
			}

			// fallback to charts dir for local charts
			if config.App.LocalChartsDir == "" {
				config.App.LocalChartsDir = config.App.ChartsDir
			}

			if config.App.EnvironmentsDir == "" {
				return fmt.Errorf("environments dir must not be empty")
			}

			for _, r := range config.Repos {
				if err := r.Validate(); err != nil {
					return fmt.Errorf("repo %q not valid: %w", r.Name, err)
				}
			}

			for _, c := range config.Charts {
				if err := c.Validate(config.Repos); err != nil {
					return fmt.Errorf("chart %q not valid: %w", c.Name, err)
				}
			}

			for _, e := range config.Environments {
				if err := e.Validate(config.Charts); err != nil {
					return fmt.Errorf("environment %q not valid: %w", e.Name, err)
				}
			}

			// all configurations are valid

			appCtx = context.WithValue(ctx, constant.ContextKeyConfig, config)

			return nil
		},
	}

	fs := cmd.PersistentFlags()

	fs.StringSliceVarP(&configFiles, "config", "c", []string{
		constant.DefaultHelmStackConfigFile,
		constant.DefaultHelmStackConfigDir,
	}, "set config files")

	fs.BoolVar(&config.App.DebugHelm, "debugHelm", false, "debug helm commands")
	fs.StringVar(&config.App.ChartsDir, "chartsDir",
		constant.DefaultChartsDir, "set directory for all listed charts")
	fs.StringVar(&config.App.EnvironmentsDir, "environmentsDir",
		constant.DefaultEnvironmentsDir, "set directory for all deployment environments")

	cmd.AddCommand(
		NewEnsureCommand(&appCtx),
		NewGenCommand(&appCtx),
		NewApplyCommand(&appCtx),
		NewCleanCommand(&appCtx),
	)

	return cmd
}

func readConfigAndResolve(configFileOrDir string, rc *conf.ResolvedConfig) error {
	err := filepath.Walk(configFileOrDir, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return fmt.Errorf("failed to visit %q: %w", path, e)
		}

		if info.IsDir() {
			return nil
		}

		switch filepath.Ext(info.Name()) {
		case ".yaml", ".yml", ".json":
		default:
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %q: %w", path, err)
		}

		if len(data) == 0 || !bytes.Contains(data, []byte(":")) {
			return nil
		}

		dec := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(data), 100)

		config := new(conf.Config)

		err = dec.Decode(config)
		if err != nil {
			return fmt.Errorf("failed to decode config for file %q: %w", path, e)
		}

		for i, r := range config.Repos {
			if _, exists := rc.Repos[r.Name]; exists {
				return fmt.Errorf("duplicate repo name %q", r.Name)
			}
			rc.Repos[r.Name] = &config.Repos[i]
		}

		for i, c := range config.Charts {
			if _, exists := rc.Charts[c.Name]; exists {
				return fmt.Errorf("duplicate chart name %q", c.Name)
			}
			rc.Charts[c.Name] = &config.Charts[i]
		}

		for i, e := range config.Environments {
			if existingEnv, exists := rc.Environments[e.Name]; exists {
				if existingEnv.KubeContext != e.KubeContext {
					return fmt.Errorf("environment %q configured with multiple kubeContext", e.Name)
				}

				// merge environment deployments
				rc.Environments[e.Name].Deployments = append(
					rc.Environments[e.Name].Deployments,
					config.Environments[i].Deployments...,
				)
			} else {
				rc.Environments[e.Name] = &config.Environments[i]
			}
		}

		rc.App = rc.App.Override(&config.App)

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to resolve config: %w", err)
	}

	return nil
}
