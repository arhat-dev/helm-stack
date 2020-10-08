package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"arhat.dev/pkg/exechelper"
	"github.com/rogpeppe/go-internal/semver"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/version"

	"arhat.dev/helm-stack/pkg/conf"
	"arhat.dev/helm-stack/pkg/constant"
)

func NewApplyCommand(appCtx *context.Context) *cobra.Command {
	var (
		dryRun bool
	)

	cmd := &cobra.Command{
		Use:           "apply <environment name>",
		Short:         "run kubectl apply with generated manifests",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.MinimumNArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			config := (*appCtx).Value(constant.ContextKeyConfig).(*conf.ResolvedConfig)
			return runApply(*appCtx, config, dryRun, args)
		},
	}

	fs := cmd.Flags()

	fs.BoolVar(&dryRun, "dry-run", false, "run kubectl apply with --dry-run=client")

	return cmd
}

func runApply(ctx context.Context, config *conf.ResolvedConfig, dryRun bool, names []string) error {
	toApply, err := GetEnvironmentsToRun(names, config)
	if err != nil {
		return err
	}

	dryRunArg := ""
	if dryRun {
		dryRunArg = "--dry-run"

		func() {
			buf := new(bytes.Buffer)
			_, err = exechelper.Do(exechelper.Spec{
				Command: []string{"kubectl", "version", "--client", "--output", "json"},
				Stdout:  buf,
				Stderr:  os.Stderr,
			})

			if err != nil {
				return
			}

			dec := json.NewDecoder(buf)
			ver := new(version.Version)
			err := dec.Decode(ver)
			if err != nil && ver.ClientVersion == nil {
				return
			}

			if semver.Compare(ver.ClientVersion.String(), "v1.18") > 0 {
				dryRunArg = "--dry-run=client"
			}
		}()
	}

	for _, e := range toApply {
		fmt.Println("--- Applying:", e.Name)

		err := e.Apply(ctx, dryRunArg, config.App.EnvironmentsDir, config.Charts)
		if err != nil {
			return fmt.Errorf("failed to apply manifests %q: %w", e.Name, err)
		}
	}

	return nil
}
