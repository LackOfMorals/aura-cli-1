package session

import (
	"net/http"

	"github.com/neo4j/cli/aura-cli/internal/api"
	"github.com/neo4j/cli/aura-cli/internal/output"
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewListCmd(cfg *clicfg.Config) *cobra.Command {
	var tenantId string
	var instanceId string
	var organizationId string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Returns a list of Graph Analytics Serverless sessions",
		Long: `This subcommand returns a list containing a summary of each of your Graph Analytics Serverless session
				By default, this subcommand lists all sessions a user has access to across all projects.
				You can filter sessions in a particular project/tenant using:
				--organization-id <organization-id>
				--tenant-id <tenant-id>
				--instance-id <instance-id>
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "/graph-analytics/sessions"

			queryParams := make(map[string]string)
			if organizationId != "" {
				queryParams["organizationId"] = organizationId
			}
			if tenantId != "" {
				queryParams["tenantId"] = tenantId
			}
			if instanceId != "" {
				queryParams["instanceId"] = instanceId
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
				output.PrintBody(cmd, cfg, resBody, []string{"id", "name", "status", "tenant_id", "cloud_provider", "ttl"})
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&tenantId, "tenant-id", "", "An optional Project ID to filter sessions in a project/tenant")
	cmd.Flags().StringVar(&organizationId, "organization-id", "", "An optional Organization ID to filter sessions in an organization")
	cmd.Flags().StringVar(&instanceId, "instance-id", "", "An optional Instance ID to filter for sessions attached to an instance")

	return cmd
}
