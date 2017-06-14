package aci

import (
	"fmt"
	"net/http"
)

const (
	tenantsPath     = "/api/mo/uni/tn-%s.json"
	listTenantsPath = "/api/node/class/fvTenant.json"
)

// Tenant ...
type Tenant struct {
	Name string
}

// Tenants ...
type Tenants struct {
	TenantsData []TenantsData `json:"imdata"`
}

// TenantsData ...
type TenantsData struct {
	FvTenant `json:"fvTenant"`
}

// FvTenant ...
type FvTenant struct {
	TenantAttributes `json:"attributes"`
}

// TenantAttributes ...
type TenantAttributes struct {
	ChildAction string `json:"childAction,omitempty"`
	Descr       string `json:"descr,omitempty"`
	Dn          string `json:"dn,omitempty"`
	LcOwn       string `json:"lcOwn,omitempty"`
	ModTs       string `json:"modTs,omitempty"`
	MonPolDn    string `json:"monPolDn,omitempty"`
	Name        string `json:"name,omitempty"`
	NameAlias   string `json:"nameAlias,omitempty"`
	OwnerKey    string `json:"ownerKey,omitempty"`
	OwnerTag    string `json:"ownerTag,omitempty"`
	Rn          string `json:"rn,omitempty"`
	Status      string `json:"status,omitempty"`
	UID         string `json:"uid,omitempty"`
}

// editTenant takes a createModify or delete action and performs the
// necessary API request
func editTenant(c *Client, tenant Tenant, action string) error {
	td := TenantsData{}
	td.Dn = "uni/tn-" + tenant.Name
	td.Name = tenant.Name
	td.Rn = "tn-" + tenant.Name
	td.Status = action

	p := fmt.Sprintf(tenantsPath, tenant.Name)
	req, err := c.newRequest(http.MethodPost, p, td)
	if err != nil {
		return err
	}

	var f interface{}

	_, err = c.do(req, &f)
	return err
}

// CreateTenant adds a slice of nodes to the ACI fabric memebership
func CreateTenant(c *Client, tenant Tenant) error {
	err := editTenant(c, tenant, createModify)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTenant adds a slice of nodes to the ACI fabric memebership
func DeleteTenant(c *Client, tenant Tenant) error {
	err := editTenant(c, tenant, delete)
	if err != nil {
		return err
	}
	return nil
}

// ListTenants lists all node members of the ACI fabric
func ListTenants(c *Client) ([]Tenant, error) {
	req, err := c.newRequest(http.MethodGet, listTenantsPath, nil)
	if err != nil {
		return nil, fmt.Errorf("list tenants: %v", err)
	}

	var ts Tenants
	_, err = c.do(req, &ts)
	if err != nil {
		return nil, fmt.Errorf("list tenants: %v", err)
	}

	var tenants []Tenant
	for _, td := range ts.TenantsData {
		tenants = append(tenants, Tenant{td.Name})
	}
	return tenants, nil
}

// Key ...
func (t Tenant) Key() string {
	return t.Name
}

// Value ...
func (t Tenant) Value() Tenant {
	return t
}
