package implementation

import (
	"fmt"
	"strings"
	"terraform-provider-infomaniak/internal/apis/domain"
	"terraform-provider-infomaniak/internal/apis/helpers"

	"resty.dev/v3"
)

// Ensure that our client implements Api
var (
	_ domain.Api = (*Client)(nil)
)

type Client struct {
	resty *resty.Client
}

func New(baseUri, token, version string) *Client {
	return &Client{
		resty: resty.New().
			SetBaseURL(baseUri).
			SetAuthToken(token).
			SetHeader("User-Agent", helpers.GetUserAgent(version)),
	}
}

func (client *Client) GetZone(fqdn string) (*domain.Zone, error) {
	var result helpers.NormalizedApiResponse[*domain.Zone]

	resp, err := client.resty.R().
		SetPathParam("fqdn", fmt.Sprint(fqdn)).
		SetQueryParam("with", "records,idn").
		SetResult(&result).
		SetError(&result).
		Get(EndpointZone)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	return result.Data, nil
}

func (client *Client) CreateZone(fqdn string) (*domain.Zone, error) {
	var result helpers.NormalizedApiResponse[*domain.Zone]

	resp, err := client.resty.R().
		SetPathParam("fqdn", fmt.Sprint(fqdn)).
		SetQueryParam("with", "records,idn").
		SetResult(&result).
		SetError(&result).
		Post(EndpointZone)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	return result.Data, nil
}

func (client *Client) DeleteZone(fqdn string) (bool, error) {
	var result helpers.NormalizedApiResponse[bool]

	resp, err := client.resty.R().
		SetPathParam("fqdn", fmt.Sprint(fqdn)).
		SetResult(&result).
		SetError(&result).
		Delete(EndpointZone)
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, result.Error
	}

	return result.Data, nil
}

func (client *Client) GetRecord(zoneFqdn string, id int64) (*domain.Record, error) {
	var result helpers.NormalizedApiResponse[*domain.Record]

	resp, err := client.resty.R().
		SetPathParam("zone_fqdn", strings.TrimSuffix(zoneFqdn, ".")).
		SetPathParam("id", fmt.Sprint(id)).
		SetQueryParam("with", "idn,records_description").
		SetResult(&result).
		SetError(&result).
		Get(EndpointRecord)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	return result.Data, nil
}

type CreateRecordRequest struct {
	Type   string `json:"type"`
	Source string `json:"source"`
	Target string `json:"target"`
	TTL    int64  `json:"ttl"`
}

func (client *Client) CreateRecord(zoneFqdn, recordType, source, target string, ttl int64) (*domain.Record, error) {
	var result helpers.NormalizedApiResponse[*domain.Record]

	var input = CreateRecordRequest{
		Type:   recordType,
		Source: source,
		Target: target,
		TTL:    ttl,
	}

	resp, err := client.resty.R().
		SetPathParam("zone_fqdn", strings.TrimSuffix(zoneFqdn, ".")).
		SetQueryParam("with", "idn,records_description").
		SetResult(&result).
		SetBody(input).
		SetError(&result).
		Post(EndpointRecords)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	return result.Data, nil
}

func (client *Client) UpdateRecord(zoneFqdn string, id int64, recordType, source, target string, ttl int64) (*domain.Record, error) {
	var result helpers.NormalizedApiResponse[*domain.Record]

	var input = CreateRecordRequest{
		Type:   recordType,
		Source: source,
		Target: target,
		TTL:    ttl,
	}

	resp, err := client.resty.R().
		SetPathParam("zone_fqdn", strings.TrimSuffix(zoneFqdn, ".")).
		SetPathParam("id", fmt.Sprint(id)).
		SetQueryParam("with", "idn,records_description").
		SetResult(&result).
		SetBody(input).
		SetError(&result).
		Put(EndpointRecord)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, result.Error
	}

	return result.Data, nil
}

func (client *Client) DeleteRecord(zoneFqdn string, id int64) (bool, error) {
	var result helpers.NormalizedApiResponse[bool]

	resp, err := client.resty.R().
		SetPathParam("zone_fqdn", strings.TrimSuffix(zoneFqdn, ".")).
		SetPathParam("id", fmt.Sprint(id)).
		SetResult(&result).
		SetError(&result).
		Delete(EndpointRecord)
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, result.Error
	}

	return result.Data, nil
}
