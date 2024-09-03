package customermanagedkey

import (
	"fmt"
	"net/http"

	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Returns a customer managed key details",
		Long:  `This subcommand returns details about a specific Customer Managed Key.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := fmt.Sprintf("/customer-managed-keys/%s", args[0])
			resBody, statusCode, err := api.MakeRequest(cmd, http.MethodGet, path, nil)
			if err != nil {
				return err
			}

			if statusCode == http.StatusOK {
				output.PrintBody2(cmd, resBody, []string{"id", "name", "tenant_id", "status", "created", "cloud_provider", "key_id", "region", "type"})

			}

			return nil
		},
	}
}
