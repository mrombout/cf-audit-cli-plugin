package client

import "encoding/json"

type AuditEventResponse struct {
	Pagination Pagination
	Resources  []AuditEventModel
}

type Pagination struct {
	TotalResults int `json:"total_results"`
	TotalPages   int `json:"total_pages"`
}

type AuditEventModel struct {
	Guid         string
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Type         string
	Actor        Actor
	Target       Target
	Data         json.RawMessage
	Space        Space
	Organization Organization
	Links        json.RawMessage
}

type Actor struct {
	Guid string
	Type string
	Name string
}

type Target struct {
	Guid string
	Type string
	Name string
}

type Space struct {
	Guid string
}

type Organization struct {
	Organization string
}

type AuditServiceBindingDeleteData struct {
	Request AuditServiceBindingDeleteDataRequest
}

type AuditServiceBindingDeleteDataRequest struct {
	AppGuid             string `json:"app_guid"`
	ServiceInstanceGuid string `json:"service_instance_guid"`
}

type AuditServiceBindingCreateData struct {
	Request AuditServiceBindingCreateDataRequest
}

type AuditServiceBindingCreateDataRequest struct {
	Data          string
	Name          string
	Relationships map[string]AuditServiceBindingCreateDataRelationship
}

type AuditServiceBindingCreateDataRelationship struct {
	Data RelationShipData
}

type RelationShipData struct {
	Guid string
}
