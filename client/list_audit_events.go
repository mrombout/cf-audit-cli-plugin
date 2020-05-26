package client

import (
	"code.cloudfoundry.org/cli/plugin"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type CloudFoundryClient struct {
	HttpClient              http.Client
	CliConnection           plugin.CliConnection
	ApiEndpoint             *url.URL
	AccessToken             string
	ListAuditEventsEndpoint *url.URL
}

type ListAuditEventsRequest struct {
	Types             []string
	TargetGuids       []string
	SpaceGuids        []string
	OrganizationGuids []string
	Page              int
	PerPage           int
	OrderBy           string
}

func (r *ListAuditEventsRequest) typesString() string {
	return strings.Join(r.Types, ",")
}

func (r *ListAuditEventsRequest) targetGuidsString() string {
	return strings.Join(r.TargetGuids, ",")
}

func (r *ListAuditEventsRequest) spaceGuidsString() string {
	return strings.Join(r.SpaceGuids, ",")
}

func (r *ListAuditEventsRequest) organizationGuidsString() string {
	return strings.Join(r.OrganizationGuids, ",")
}

func (r *ListAuditEventsRequest) query(baseValues url.Values) url.Values {
	orderBy := r.OrderBy
	if orderBy == "" {
		orderBy = "-created_at"
	}
	baseValues.Set("order_by", orderBy)

	if len(r.OrganizationGuids) > 0 {
		baseValues.Set("organization_guids", r.organizationGuidsString())
	}

	if len(r.SpaceGuids) > 0 {
		baseValues.Set("space_guids", r.spaceGuidsString())
	}

	if len(r.TargetGuids) > 0 {
		baseValues.Set("target_guids", r.targetGuidsString())
	}

	if len(r.Types) > 0 {
		baseValues.Set("types", r.typesString())
	}

	return baseValues
}

func (r *ListAuditEventsRequest) createHttpRequest(endpoint url.URL, accessToken string) (*http.Request, error) {
	endpoint.RawQuery = r.query(endpoint.Query()).Encode()

	httpRequest, err := http.NewRequest("GET", endpoint.String(), strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("Authorization", accessToken)

	return httpRequest, nil
}

func (c *CloudFoundryClient) ListAuditEvents(request ListAuditEventsRequest) ([]AuditEventModel, error) {
	httpRequest, err := request.createHttpRequest(*c.ListAuditEventsEndpoint, c.AccessToken)
	if err != nil {
		return []AuditEventModel{}, err
	}

	rawResponse, err := c.HttpClient.Do(httpRequest)
	if err != nil {
		return []AuditEventModel{}, err
	}

	responseBody, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return []AuditEventModel{}, err
	}

	var response AuditEventResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return []AuditEventModel{}, err
	}

	return response.Resources, nil
}
