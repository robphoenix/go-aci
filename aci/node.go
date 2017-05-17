package aci

import (
	"fmt"
)

const (
	nodesPath     = "api/node/mo/uni/controller/nodeidentpol.json"
	listNodesPath = "api/node/class/fabricNode.json"
)

// Node ...
type Node struct {
	FabricStatus string
	Model        string
	Name         string
	ID           string
	Serial       string
	Status       string
	Role         string
}

// FabricNodeIdentPolContainer ...
type FabricNodeIdentPolContainer struct {
	FabricNodeIdentPol `json:"fabricNodeIdentPol"`
}

// FabricNodeIdentPContainer ...
type FabricNodeIdentPContainer struct {
	FabricNodeIdentP `json:"fabricNodeIdentP"`
}

// FabricNodeIdentPol ...
type FabricNodeIdentPol struct {
	NodeAttributes `json:"attributes"`
	Children       []FabricNodeIdentPContainer `json:"children"`
}

// FabricNodeIdentP ...
type FabricNodeIdentP struct {
	NodeAttributes `json:"attributes"`
}

// FabricNodes ...
type FabricNodes struct {
	Imdata []struct {
		FabricNode `json:"fabricNode"`
	} `json:"imdata"`
	TotalCount string `json:"totalCount"`
}

// FabricNode ...
type FabricNode struct {
	NodeAttributes `json:"attributes"`
}

// nodeResponse
type nodeResponse struct {
	Imdata     []interface{} `json:"imdata"`
	TotalCount string        `json:"totalCount"`
}

// NodeAttributes ...
type NodeAttributes struct {
	Status           string `json:"status,omitempty"`
	AdSt             string `json:"adSt,omitempty"`
	ChildAction      string `json:"childAction,omitempty"`
	DelayedHeartbeat string `json:"delayedHeartbeat,omitempty"`
	Dn               string `json:"dn,omitempty"`
	FabricSt         string `json:"fabricSt,omitempty"`
	ID               string `json:"id,omitempty"`
	LcOwn            string `json:"lcOwn,omitempty"`
	ModTs            string `json:"modTs,omitempty"`
	Model            string `json:"model,omitempty"`
	MonPolDn         string `json:"monPolDn,omitempty"`
	Name             string `json:"name,omitempty"`
	NameAlias        string `json:"nameAlias,omitempty"`
	NodeID           string `json:"nodeId,omitempty"`
	Role             string `json:"role,omitempty"`
	Serial           string `json:"serial,omitempty"`
	UID              string `json:"uid,omitempty"`
	Vendor           string `json:"vendor,omitempty"`
	Version          string `json:"version,omitempty"`
}

func createFNIPolC(ns []Node, action string) FabricNodeIdentPolContainer {
	var children []FabricNodeIdentPContainer
	for _, n := range ns {
		a := NodeAttributes{
			Name:   n.Name,
			NodeID: n.ID,
			Serial: n.Serial,
			Status: action,
		}
		fp := FabricNodeIdentPContainer{
			FabricNodeIdentP: FabricNodeIdentP{
				NodeAttributes: a,
			},
		}
		children = append(children, fp)
	}
	fpol := FabricNodeIdentPolContainer{
		FabricNodeIdentPol{
			NodeAttributes: NodeAttributes{
				Status: createModify,
			},
			Children: children,
		},
	}
	return fpol
}

func (c *Client) editNodes(ns []Node, action string) error {
	fpol := createFNIPolC(ns, action)
	req, err := c.newRequest("POST", nodesPath, fpol)
	if err != nil {
		return err
	}

	var f nodeResponse

	resp, err := c.do(req, &f)
	if err != nil {
		return err
	}
	fmt.Println("response Status:", resp.Status)
	return nil
}

// AddNodes ...
func (c *Client) AddNodes(ns []Node) error {
	err := c.editNodes(ns, create)
	if err != nil {
		return err
	}
	return nil
}

// DeleteNodes ...
func (c *Client) DeleteNodes(ns []Node) error {
	err := c.editNodes(ns, delete)
	if err != nil {
		return err
	}
	return nil
}

// ListNodes ...
func (c *Client) ListNodes() ([]Node, error) {
	req, err := c.newRequest("GET", listNodesPath, nil)
	if err != nil {
		return nil, err
	}

	var n FabricNodes
	resp, err := c.do(req, &n)
	if err != nil {
		return nil, err
	}

	fmt.Println("response Status:", resp.Status)
	var ns []Node
	for _, v := range n.Imdata {
		node := Node{
			FabricStatus: v.FabricSt,
			ID:           v.ID,
			Model:        v.Model,
			Name:         v.Name,
			Role:         v.Role,
			Serial:       v.Serial,
			Status:       v.Status,
		}
		ns = append(ns, node)
	}
	return ns, nil
}

func (c *Client) decomissionNode() error {
	// https://supportforums.cisco.com/discussion/13296271/decommissioning-fabric-nodes-api
	// url := "https://sandboxapicdc.cisco.com/api/node/mo/uni/fabric/outofsvc.json"
	// {
	//   "fabricRsDecommissionNode": {
	//     "attributes": {
	//       "tDn": "topology/pod-1/node-102",
	//       "status": "created,modified",
	//       "removeFromController": "true"
	//     },
	//     "children": []
	//   }
	// }
	return nil
}
