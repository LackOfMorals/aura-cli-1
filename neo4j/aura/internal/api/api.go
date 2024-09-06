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

type AuraApiConfig struct {
	BaseUrl   string
	Token     string
	UserAgent string
}

type AuraApi struct {
	Instances *AuraInstancesApi
}

func New(config *AuraApiConfig) *AuraApi {
	return &AuraApi{
		Instances: &AuraInstancesApi{
			config: config,
		},
	}
}

func (config *AuraApiConfig) makeRequest(method string, path string, data map[string]any) (responseBody []byte, statusCode int, err error) {
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

	u, _ := url.ParseRequestURI(config.BaseUrl)
	u = u.JoinPath(path)
	urlString := u.String()

	req, err := http.NewRequest(method, urlString, body)
	if err != nil {
		return responseBody, 0, err
	}

	req.Header = config.getHeaders()

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

	return nil, res.StatusCode, config.handleError(res.StatusCode, responseBody)
}

func (config *AuraApiConfig) getHeaders() http.Header {
	return http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {fmt.Sprintf("Bearer %s", config.Token)},
		"User-Agent":    {config.UserAgent},
	}
}

func (config *AuraApiConfig) handleError(_statusCode int, responseBody []byte) error {
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
	if err != nil {
		return nil, err
	}

	token, err := getToken(cmd.Context())

	if err != nil {
		return nil, err
	}

	version, ok := clictx.Version(cmd.Context())
	if !ok {
		return nil, errors.New("error fetching version from context")
	}

	return New(&AuraApiConfig{
		BaseUrl:   baseUrl,
		Token:     token,
		UserAgent: fmt.Sprintf(userAgent, version),
	}), nil
}
