package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AuraInstancesApi struct {
	config *AuraApiConfig
}

type InstanceGetResponse struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Status                string `json:"status"`
	TenantID              string `json:"tenant_id"`
	CloudProvider         string `json:"cloud_provider"`
	ConnectionURL         string `json:"connection_url"`
	MetricsIntegrationURL string `json:"metrics_integration_url"`
	Region                string `json:"region"`
	Type                  string `json:"type"`
	Memory                string `json:"memory"`
	Storage               string `json:"storage"`
	CustomerManagedKeyId  string `json:"customer_managed_key_id,omitempty"`
	SecondariesCount      *int   `json:"secondaries_count,omitempty"` // Using pointer to handle null cases
	CDCEnrichmentMode     string `json:"cdc_enrichment_mode,omitempty"`
}

type InstancesListResponse []struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	TenantId      string `json:"tenant_id"`
	CloudProvider string `json:"cloud_provider"`
}

func (api *AuraInstancesApi) List(tenantId string) (response *InstancesListResponse, status int, err error) {
	var path string

	if tenantId != "" {
		path = fmt.Sprintf("/instances?tenantId=%s", tenantId)
	} else {
		path = "/instances"
	}

	responseBody, statusCode, err := api.config.makeRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, statusCode, err
	}

	var parsedBody struct {
		Data InstancesListResponse
	}
	err = json.Unmarshal(responseBody, &parsedBody)
	if err != nil {
		return nil, statusCode, fmt.Errorf("error parsing response of GET %s: %v", path, err)
	}
	return &parsedBody.Data, statusCode, nil
}

func (api *AuraInstancesApi) Get(instanceId string) (response *InstanceGetResponse, status int, err error) {
	path := fmt.Sprintf("/instances/%s", instanceId)

	responseBody, statusCode, err := api.config.makeRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, statusCode, err
	}

	var parsedBody struct {
		Data InstanceGetResponse
	}

	err = json.Unmarshal(responseBody, &parsedBody)
	if err != nil {
		return nil, statusCode, fmt.Errorf("error parsing response of GET %s: %v", path, err)
	}

	return &parsedBody.Data, statusCode, nil
}
