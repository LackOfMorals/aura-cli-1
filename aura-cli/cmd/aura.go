package aura

import (
	"github.com/neo4j/cli/aura-cli/internal/subcommands/graphanalytics"
	"github.com/spf13/cobra"

	"github.com/neo4j/cli/aura-cli/internal/subcommands/config"
	"github.com/neo4j/cli/aura-cli/internal/subcommands/credential"
	"github.com/neo4j/cli/aura-cli/internal/subcommands/customermanagedkey"
	"github.com/neo4j/cli/aura-cli/internal/subcommands/dataapi"
	"github.com/neo4j/cli/aura-cli/internal/subcommands/instance"
	"github.com/neo4j/cli/aura-cli/internal/subcommands/tenant"
	"github.com/neo4j/cli/common/clicfg"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "aura-cli",
		Short:   "Allows you to programmatically provision and manage your Aura resources",
		Version: cfg.Version,
	}

	cmd.AddCommand(config.NewCmd(cfg))
	cmd.AddCommand(credential.NewCmd(cfg))
	cmd.AddCommand(customermanagedkey.NewCmd(cfg))
	cmd.AddCommand(instance.NewCmd(cfg))
	cmd.AddCommand(tenant.NewCmd(cfg))
	cmd.AddCommand(graphanalytics.NewCmd(cfg))
	if cfg.Aura.AuraBetaEnabled() {
		cmd.AddCommand(dataapi.NewCmd(cfg))
	}

	return cmd
}
