package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/neo4j/cli/common/clictx"
	"github.com/spf13/cobra"
)

type AuraApiContext struct {
	BaseUrl   string
	Token     string
	UserAgent string
}

type AuraApi struct {
	Instances *AuraInstancesApi
	Config    *AuraApiContext
}

type AuraInstancesApi struct {
	ctx *AuraApiContext
}

type InstancesGetResponse []struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	TenantId      string `json:"tenant_id"`
	CloudProvider string `json:"cloud_provider"`
}

func New(config *AuraApiContext) *AuraApi {
	return &AuraApi{
		Instances: &AuraInstancesApi{
			ctx: config,
		},
	}
}

func (api *AuraInstancesApi) List(tenantId string) (response *InstancesGetResponse, status int, err error) {
	var path string

	if tenantId != "" {
		path = fmt.Sprintf("/instances?tenantId=%s", tenantId)
	} else {
		path = "/instances"
	}

	responseBody, statusCode, err := api.ctx.makeRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, statusCode, err
	}

	var parsedBody struct {
		Data InstancesGetResponse
	}

	err = json.Unmarshal(responseBody, &parsedBody)
	if err != nil {
		return nil, statusCode, fmt.Errorf("error parsing response of GET %s: %v", path, err)
	}

	return &parsedBody.Data, statusCode, nil
}

func (ctx *AuraApiContext) makeRequest(method string, path string, data map[string]any) (responseBody []byte, statusCode int, err error) {
	client := http.Client{}
	var body io.Reader
	if data == nil {
		body = nil
	} else {
		jsonData, err := json.Marshal(data)

		if err != nil {
			return responseBody, 0, err
		}

		body = bytes.NewBuffer(jsonData)
	}

	u, _ := url.ParseRequestURI(ctx.BaseUrl)
	u = u.JoinPath(path)
	urlString := u.String()

	req, err := http.NewRequest(method, urlString, body)
	if err != nil {
		return responseBody, 0, err
	}

	req.Header = ctx.getHeaders()

	res, err := client.Do(req)
	if err != nil {
		return responseBody, 0, err
	}

	defer res.Body.Close()

	if isSuccessful(res.StatusCode) {
		responseBody, err = io.ReadAll(res.Body)

		if err != nil {
			return responseBody, 0, err
		}

		return responseBody, res.StatusCode, nil
	}

	return nil, res.StatusCode, ctx.handleError(res.StatusCode, responseBody)
}

func (config *AuraApiContext) getHeaders() http.Header {
	return http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {fmt.Sprintf("Bearer %s", config.Token)},
		// "User-Agent":    {fmt.Sprintf(userAgent, version)},
		"User-Agent": {config.UserAgent},
	}
}

func (config *AuraApiContext) handleError(_statusCode int, responseBody []byte) error {
	type ErrorResponse struct {
		Errors []struct {
			Message string
			Reason  string
			field   string
		}
	}

	var errorResponse ErrorResponse
	err := json.Unmarshal(responseBody, &errorResponse)
	if err != nil {
		return err
	}

	messages := []string{}
	for _, e := range errorResponse.Errors {
		messages = append(messages, e.Message)
	}
	return fmt.Errorf("%s", messages)
}

// Note: This is specific to the CLI, not part of the Aura API
func GetApiFromConfig(cmd *cobra.Command) (*AuraApi, error) {
	config, ok := clictx.Config(cmd.Context())
	if !ok {
		return nil, errors.New("error fetching cli configuration values")
	}

	baseUrl, err := config.GetString("aura.base-url")

	token, err := getToken(cmd.Context())

	if err != nil {
		return nil, err
	}

	version, ok := clictx.Version(cmd.Context())
	if !ok {
		return nil, errors.New("error fetching version from context")
	}

	return New(&AuraApiContext{
		BaseUrl:   baseUrl,
		Token:     token,
		UserAgent: fmt.Sprintf(userAgent, version),
	}), nil
}
