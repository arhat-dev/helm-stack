package cmd

import (
	"context"
	"fmt"

	"arhat.dev/helm-stack/pkg/conf"

	"github.com/spf13/cobra"

	"arhat.dev/helm-stack/pkg/constant"
)

func NewGenCommand(appCtx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "gen <environment name 1> ... <environment name N>",
		Short:         "generate manifests according to your custom values",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.MinimumNArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			config := (*appCtx).Value(constant.ContextKeyConfig).(*conf.ResolvedConfig)
			return runGen(*appCtx, config, args)
		},
	}

	return cmd
}

func runGen(ctx context.Context, config *conf.ResolvedConfig, names []string) error {
	toGen, err := GetEnvironmentsToRun(names, config)
	if err != nil {
		return err
	}

	for _, e := range toGen {
		fmt.Println("--- Generating Manifests:", e.Name)

		if err := e.Gen(ctx, config.App.ChartsDir, config.App.EnvironmentsDir, config.Charts); err != nil {
			return fmt.Errorf("failed to generate manifests %q: %w", e.Name, err)
		}
	}

	return nil
}
