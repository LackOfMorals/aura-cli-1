package authprovider

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "auth-provider",
		Short: "Allows you to programmatically manage authentication providers for a specific GraphQL Data API",
	}

	cmd.AddCommand(NewListCmd(cfg))
	cmd.AddCommand(NewGetCmd(cfg))
	cmd.AddCommand(NewUpdateCmd(cfg))

	return cmd
}
