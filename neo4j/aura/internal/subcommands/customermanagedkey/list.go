package customermanagedkey

import (
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewListCmd(cfg *clicfg.Config) *cobra.Command {
	var projectId string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Returns a list of customer managed keys",
		Long: `This subcommand returns a list containing a summary of each of your customer managed keys. To find out more about a specific key, retrieve the details using the get subcommand.

You can filter keys in a particular project using --project-id. If the project flag is not specified, this endpoint lists all keys a user has access to across all projects.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "/customer-managed-keys"
			queryParams := make(map[string]string)
			if projectId != "" {
				queryParams["tenantId"] = projectId
			}
			cmd.SilenceUsage = true
			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method:      http.MethodGet,
				QueryParams: queryParams,
			})
			if err != nil {
				return err
			}

			if statusCode == http.StatusOK {
				output.PrintBody(cmd, cfg, resBody, []string{"id", "name", "tenant_id"})

			}

			return nil
		},
	}

	cmd.Flags().StringVar(&projectId, "project-id", "", "An optional Project ID to filter customer managed keys in a project")

	return cmd
}
