package instance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	var tenantId string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Returns a list of instances",
		Long: `This subcommand returns a list containing a summary of each of your Aura instances. To find out more about a specific instance, retrieve the details using the get subcommand.

You can filter instances in a particular tenant using --tenant-id. If the tenant flag is not specified, this subcommand lists all instances a user has access to across all tenants.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			aura, err := api.GetApiFromConfig(cmd)
			if err != nil {
				return fmt.Errorf("error in command %s", os.Args[1:])
			}

			instances, statusCode, err := aura.Instances.List(tenantId)
			if err != nil {
				return fmt.Errorf("error in command %s: %v", os.Args[1:], err)
			}

			if statusCode == http.StatusOK {
				jsonResponse, err := json.Marshal(instances)
				if err != nil {
					return fmt.Errorf("error in command %s: %v", os.Args[1:], err)
				}

				err = output.PrintBody(cmd, jsonResponse, []string{"id", "name", "tenant_id", "cloud_provider"})
				if err != nil {
					return fmt.Errorf("error in command %s: %v", os.Args[1:], err)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&tenantId, "tenant-id", "", "An optional Tenant ID to filter instances in a tenant")

	return cmd
}
