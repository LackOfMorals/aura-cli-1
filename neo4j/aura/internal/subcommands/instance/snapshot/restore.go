package snapshot

import (
	"fmt"
	"net/http"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewRestoreCmd(cfg *clicfg.Config) *cobra.Command {
	var instanceId string

	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Restores a snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := fmt.Sprintf("/instances/%s/snapshots/%s/restore", instanceId, args[0])

			resBody, statusCode, err := api.MakeRequest(cfg, path, &api.RequestConfig{
				Method: http.MethodPost,
			})
			if err != nil {
				return err
			}

			if statusCode == http.StatusAccepted {
				err = output.PrintBody(cmd, cfg, resBody, []string{"id", "name", "status", "tenant_id", "connection_url", "cloud_provider", "type", "region", "memory"})
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&instanceId, "instance-id", "", "The id of the instance to list its snapshots")
	cmd.MarkFlagRequired("instance-id")

	return cmd
}
