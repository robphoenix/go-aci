package aci

import (
	"fmt"
	"net/http"
)

const (
	tenantsPath     = "/api/mo/uni.json"
	listTenantsPath = "/api/node/class/fvTenant.json"
)

// tenantContainer contains the data structure needed
// to create a tenant
type tenantContainer struct {
	FvTenant `json:"fvTenant"`
}

// Tenant describes a tenant
type Tenant struct {
	Name string
}

// Tenants ...
type Tenants struct {
	TenantsImdata []TenantsImdata `json:"imdata"`
}

// TenantsImdata ...
type TenantsImdata struct {
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
	Status      string `json:"status,omitempty"`
	UID         string `json:"uid,omitempty"`
}

// editTenant takes a createModify or delete action and performs the
// necessary API request
func editTenant(c *Client, t Tenant, action string) error {
	ta := TenantAttributes{
		Name:   t.Name,
		Status: action,
	}
	req, err := c.newRequest(http.MethodPost, tenantsPath, ta)
	if err != nil {
		return err
	}

	var f interface{}

	_, err = c.do(req, &f)
	return err
}

// AddTenant adds a slice of nodes to the ACI fabric memebership
func AddTenant(c *Client, t Tenant) error {
	err := editTenant(c, t, createModify)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTenant adds a slice of nodes to the ACI fabric memebership
func DeleteTenant(c *Client, t Tenant) error {
	err := editTenant(c, t, delete)
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
	for _, t := range ts.TenantsImdata {
		tenants = append(tenants, Tenant{Name: t.Name})
	}
	return tenants, nil
}
