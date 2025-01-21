package implementation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"terraform-provider-infomaniak/internal/apis/endpoints"
	"terraform-provider-infomaniak/internal/apis/kaas"
)

// Ensure that our client implements Api
var (
	_ kaas.Api = (*Client)(nil)
)

type Client struct {
	baseUri    string
	httpClient *http.Client
}

func New(baseUri string) *Client {
	return &Client{
		baseUri:    baseUri,
		httpClient: http.DefaultClient,
	}
}

type ApiResponse[K any] struct {
	Result string `json:"result"`
	Data   K
}

// UnmarshalResponse unmarshal a http response into a struct,
// The body must not be closed for this to work properly
func UnmarshalResponse[K any](response *http.Response, result *K) error {
	if result == nil {
		return fmt.Errorf("result musn't be nil")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var parsedResponse ApiResponse[K]

	err = json.Unmarshal(body, &parsedResponse)
	*result = parsedResponse.Data
	return err
}

func (client *Client) Do(route *endpoints.CompiledEndpoint, data any) (*http.Response, error) {
	uri, err := url.JoinPath(client.baseUri, route.URL)
	if err != nil {
		return nil, err
	}

	var body io.Reader = nil
	if data != nil {
		rawData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		body = bytes.NewReader(rawData)
	}

	req, err := http.NewRequest(route.Endpoint.Method, uri, body)
	if err != nil {
		return nil, err
	}

	return client.httpClient.Do(req)
}

func (client *Client) GetPacks() ([]*kaas.KaasPack, error) {
	compiledRoute, err := GetPacks.Compile(nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(compiledRoute, nil)
	if err != nil {
		return nil, err
	}

	var result []*kaas.KaasPack
	err = UnmarshalResponse(response, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (client *Client) GetKaas(pcpId, kaasId string) (*kaas.Kaas, error) {
	compiledRoute, err := GetKaas.Compile(nil, pcpId, kaasId)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(compiledRoute, nil)
	if err != nil {
		return nil, err
	}

	var result kaas.Kaas
	err = UnmarshalResponse(response, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *Client) CreateKaas(input *kaas.Kaas) (*kaas.Kaas, error) {
	if input.PcpId == "" {
		return nil, fmt.Errorf("kaas is missing public cloud project id")
	}

	compiledRoute, err := CreateKaas.Compile(nil, input.PcpId)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(compiledRoute, nil)
	if err != nil {
		return nil, err
	}

	var result kaas.InstancePool
	err = UnmarshalResponse(response, &result)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (client *Client) UpdateKaas(input *kaas.Kaas) (*kaas.Kaas, error) {
	return nil, nil
}

func (client *Client) DeleteKaas(pcpId, kaasId string) error {
	return nil
}

func (client *Client) GetInstancePool(pcpId, kaasId, instancePoolId string) (*kaas.InstancePool, error) {
	compiledRoute, err := GetInstancePool.Compile(nil, pcpId, kaasId, instancePoolId)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(compiledRoute, nil)
	if err != nil {
		return nil, err
	}

	var result kaas.InstancePool
	err = UnmarshalResponse(response, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *Client) CreateInstancePool(input *kaas.InstancePool) (*kaas.InstancePool, error) {
	return nil, nil
}

func (client *Client) UpdateInstancePool(input *kaas.InstancePool) (*kaas.InstancePool, error) {
	return nil, nil
}

func (client *Client) DeleteInstancePool(pcpId, kaasId, instancePoolId string) error {
	return nil
}
