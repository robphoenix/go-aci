package aci

import (
	"fmt"
)

const (
	nodesPath     = "api/node/mo/uni/controller/nodeidentpol.json"
	listNodesPath = "api/node/class/fabricNode.json"
)

// Node is a member of an ACI fabric
type Node struct {
	FabricStatus string
	Model        string
	Name         string
	ID           string
	Serial       string
	Status       string
	Role         string
}

// FabricNodeIdentPolContainer is a container for a FabricNodeIdentPol
type FabricNodeIdentPolContainer struct {
	FabricNodeIdentPol `json:"fabricNodeIdentPol"`
}

// FabricNodeIdentPol is a container for an ACI node identity profile
type FabricNodeIdentPol struct {
	NodeAttributes `json:"attributes"`
	Children       []FabricNodeIdentPContainer `json:"children"`
}

// FabricNodeIdentPContainer is a container for a FabricNodeIdentP
type FabricNodeIdentPContainer struct {
	FabricNodeIdentP `json:"fabricNodeIdentP"`
}

// FabricNodeIdentP is the node identity profile,
// that assigns IDs to the fabric nodes
type FabricNodeIdentP struct {
	NodeAttributes `json:"attributes"`
}

// FabricNodes contains is the response body for requests
// about current ACI fabric nodes
type FabricNodes struct {
	Imdata []Imdata `json:"imdata"`
}

// Imdata is a container for the response structure
type Imdata struct {
	FabricNodeIdentP `json:"fabricNode"`
}

// NodeAttributes contains all the attributes of an ACI fabric node
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

func buildFabricNodeContainer(ns []Node, action string) FabricNodeIdentPolContainer {
	var c []FabricNodeIdentPContainer
	for _, n := range ns {
		var fp FabricNodeIdentPContainer
		fp.Name = n.Name
		fp.NodeID = n.ID
		fp.Serial = n.Serial
		fp.Status = action
		c = append(c, fp)
	}
	var fpc FabricNodeIdentPolContainer
	fpc.Status = createModify
	fpc.Children = c
	return fpc
}

// editNodes takes a createModify or delete action and performs the
// necessary API request
func (c *Client) editNodes(ns []Node, action string) error {
	fpol := buildFabricNodeContainer(ns, action)
	req, err := c.newRequest("POST", nodesPath, fpol)
	if err != nil {
		return err
	}

	var f FabricNodes

	_, err = c.do(req, &f)
	return err
}

// AddNodes adds a slice of nodes to the ACI fabric memebership
func (c *Client) AddNodes(ns []Node) error {
	err := c.editNodes(ns, createModify)
	if err != nil {
		return err
	}
	return nil
}

// DeleteNodes deletes a slice of nodes from the ACI fabric membership
func (c *Client) DeleteNodes(ns []Node) error {
	err := c.editNodes(ns, delete)
	if err != nil {
		return err
	}
	return nil
}

// ListNodes lists all node members of the ACI fabric
func (c *Client) ListNodes() ([]Node, error) {
	req, err := c.newRequest("GET", listNodesPath, nil)
	if err != nil {
		return nil, fmt.Errorf("list nodes: %v", err)
	}

	var n FabricNodes
	_, err = c.do(req, &n)
	if err != nil {
		return nil, fmt.Errorf("list nodes: %v", err)
	}

	var ns []Node
	for _, v := range n.Imdata {
		if v.Role != "controller" {
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
