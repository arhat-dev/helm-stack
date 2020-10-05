package cmd

import (
	"context"
	"fmt"

	"arhat.dev/helm-stack/pkg/conf"

	"github.com/spf13/cobra"

	"arhat.dev/helm-stack/pkg/constant"
)

func NewEnsureCommand(appCtx *context.Context) *cobra.Command {
	var (
		forcePull bool
	)

	cmd := &cobra.Command{
		Use:   "ensure",
		Short: "ensure directories and files for charts and environments",
		Long: "create charts and environments directories, " +
			"pull charts according to your configuration extract files to your environments directories",
		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: func(cmd *cobra.Command, args []string) error {
			config := (*appCtx).Value(constant.ContextKeyConfig).(*conf.ResolvedConfig)
			return runEnsure(*appCtx, config, forcePull)
		},
	}

	fs := cmd.Flags()
	fs.BoolVar(&forcePull, "force-pull", false, "pull chart even though already exists")

	return cmd
}

func runEnsure(ctx context.Context, config *conf.ResolvedConfig, forcePull bool) error {
	for _, c := range config.Charts {
		fmt.Println("--- Ensuring Chart:", c.Name)

		err := c.Ensure(ctx, forcePull, config.App.ChartsDir, config.Repos)
		if err != nil {
			return fmt.Errorf("failed to ensure chart %q: %w", c.Name, err)
		}
	}

	for _, e := range config.Environments {
		fmt.Println("--- Ensuring Environment:", e.Name)

		if err := e.Ensure(ctx, config.App.ChartsDir, config.App.EnvironmentsDir, config.Charts); err != nil {
			return fmt.Errorf("failed to ensure deployment environment %q: %w", e.Name, err)
		}
	}

	return nil
}
