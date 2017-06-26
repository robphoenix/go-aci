package aci

import (
	"fmt"
	"net/http"
)

const (
	vrfPath     = "/api/mo/uni/tn-%s/ctx-%s.json"
	vrfListPath = "/api/node/mo/uni/tn-%s.json?query-target=children&target-subtree-class=fvCtx"
)

// A VRF represents a single ACI VRF.
//
// A VRF has a name and is associated
// with an ACI Tenant.
type VRF struct {
	Name   string
	Tenant string
}

// VRFs ...
type VRFs struct {
	VRFsData []VRFsData `json:"imdata"`
}

// VRFsData ...
type VRFsData struct {
	FvCtx `json:"fvCtx"`
}

// FvCtx ...
type FvCtx struct {
	VRFAttributes `json:"attributes"`
}

// VRFAttributes ...
type VRFAttributes struct {
	ChildAction     string `json:"childAction,omitempty"`
	Descr           string `json:"descr,omitempty"`
	Dn              string `json:"dn,omitempty"`
	KnwMcastAct     string `json:"knwMcastAct,omitempty"`
	LcOwn           string `json:"lcOwn,omitempty"`
	ModTs           string `json:"modTs,omitempty"`
	MonPolDn        string `json:"monPolDn,omitempty"`
	Name            string `json:"name,omitempty"`
	NameAlias       string `json:"nameAlias,omitempty"`
	OwnerKey        string `json:"ownerKey,omitempty"`
	OwnerTag        string `json:"ownerTag,omitempty"`
	PcEnfDir        string `json:"pcEnfDir,omitempty"`
	PcEnfDirUpdated string `json:"pcEnfDirUpdated,omitempty"`
	PcEnfPref       string `json:"pcEnfPref,omitempty"`
	PcTag           string `json:"pcTag,omitempty"`
	Rn              string `json:"rn,omitempty"`
	Scope           string `json:"scope,omitempty"`
	Seg             string `json:"seg,omitempty"`
	Status          string `json:"status,omitempty"`
	UID             string `json:"uid,omitempty"`
}

// editVRF ...
func editVRF(c *Client, v VRF, action string) error {
	vd := VRFsData{}
	vd.Dn = "uni/tn-" + v.Tenant + "/ctx-" + v.Name
	vd.Name = v.Name
	vd.Rn = "ctx-" + v.Name
	vd.Status = action

	p := fmt.Sprintf(vrfPath, v.Tenant, v.Name)
	req, err := c.NewRequest(http.MethodPost, p, vd)
	if err != nil {
		return err
	}

	var f interface{}

	_, err = c.Do(req, &f)
	return err
}

// CreateVRF ...
func CreateVRF(c *Client, v VRF) error {
	err := editVRF(c, v, createModify)
	if err != nil {
		return err
	}
	return nil
}

// DeleteVRF ...
func DeleteVRF(c *Client, v VRF) error {
	err := editVRF(c, v, delete)
	if err != nil {
		return err
	}
	return nil
}

// ListVRFs lists all node members of the ACI fabric
func ListVRFs(c *Client, tenant string) ([]VRF, error) {
	p := fmt.Sprintf(vrfListPath, tenant)
	req, err := c.NewRequest(http.MethodGet, p, nil)
	if err != nil {
		return nil, fmt.Errorf("list VRFs: %v", err)
	}

	var vs VRFs
	_, err = c.Do(req, &vs)
	if err != nil {
		return nil, fmt.Errorf("list VRFs: %v", err)
	}

	var vrfs []VRF
	for _, v := range vs.VRFsData {
		vrfs = append(vrfs, VRF{Name: v.Name, Tenant: tenant})
	}
	return vrfs, nil
}
