package client

import (
	"code.cloudfoundry.org/cli/plugin"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type CloudFoundryClient struct {
	CliConnection           plugin.CliConnection
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

func (c *CloudFoundryClient) ListAuditEvents(request ListAuditEventsRequest) ([]AuditEventModel, error) {
	endpoint := url.URL{Path: "/v3/audit_events"}
	endpoint.RawQuery = request.query(endpoint.Query()).Encode()

	output, err := c.CliConnection.CliCommandWithoutTerminalOutput("curl", endpoint.String())
	if err != nil {
		return nil, err
	}
	if len(output) == 0 {
		return nil, fmt.Errorf("no response received from endpoint: %s", endpoint.String())
	}
	outputBody := strings.Join(output, "")

	var response AuditEventResponse
	err = json.Unmarshal([]byte(outputBody), &response)
	if err != nil {
		return nil, err
	}

	return response.Resources, nil
}
