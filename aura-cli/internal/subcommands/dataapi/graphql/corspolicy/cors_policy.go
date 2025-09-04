package corspolicy

import (
	"github.com/neo4j/cli/aura-cli/internal/subcommands/dataapi/graphql/corspolicy/allowedorigin"
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "cors-policy",
		Short: "Allows you to manage the Cross-Origin Resource Sharing (CORS) policy for a specific GraphQL Data API",
	}

	cmd.AddCommand(allowedorigin.NewCmd(cfg))

	return cmd
}
