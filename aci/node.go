package aci

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
	Attributes `json:"attributes"`
	Children   []FabricNodeIdentPContainer `json:"children"`
}

// FabricNodeIdentP ...
type FabricNodeIdentP struct {
	Attributes `json:"attributes"`
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
	Attributes `json:"attributes"`
}

// Attributes ...
type Attributes struct {
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

// AddNodes ...
func (c *Client) AddNodes(ns []Node) error {
	var children []FabricNodeIdentPContainer
	for _, n := range ns {
		a := Attributes{
			Name:   n.Name,
			NodeID: n.ID,
			Serial: n.Serial,
			Status: create,
		}
		fp := FabricNodeIdentPContainer{
			FabricNodeIdentP: FabricNodeIdentP{
				Attributes: a,
			},
		}
		children = append(children, fp)
	}
	fpol := FabricNodeIdentPolContainer{
		FabricNodeIdentPol{
			Attributes: Attributes{
				Status: createModify,
			},
			Children: children,
		},
	}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(fpol)
	if err != nil {
		return err
	}

	nodesURL := url.URL{Scheme: c.Host.Scheme, Host: c.Host.Host, Path: nodesPath}
	req, err := http.NewRequest("POST", nodesURL.String(), b)
	if err != nil {
		return err
	}
	req.Header.Set("Cookie", c.Cookie)
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	nodesBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("nodesBody = %+v\n", string(nodesBody))
	return nil
}

// DeleteNodes ...
func (c *Client) DeleteNodes(ns []Node) error {
	var children []FabricNodeIdentPContainer
	for _, n := range ns {
		a := Attributes{
			Name:   n.Name,
			NodeID: n.ID,
			Serial: n.Serial,
			Status: delete,
		}
		fp := FabricNodeIdentPContainer{
			FabricNodeIdentP: FabricNodeIdentP{
				Attributes: a,
			},
		}
		children = append(children, fp)
	}
	fpol := FabricNodeIdentPolContainer{
		FabricNodeIdentPol{
			Attributes: Attributes{
				Status: createModify,
			},
			Children: children,
		},
	}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(fpol)
	if err != nil {
		return err
	}

	nodesURL := url.URL{Scheme: c.Host.Scheme, Host: c.Host.Host, Path: nodesPath}
	req, err := http.NewRequest("POST", nodesURL.String(), b)
	if err != nil {
		return err
	}
	req.Header.Set("Cookie", c.Cookie)
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	nodesBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("nodesBody = ", string(nodesBody))
	return nil
}

// ListNodes ...
func (c *Client) ListNodes() ([]Node, error) {
	listNodesURL := url.URL{Scheme: c.Host.Scheme, Host: c.Host.Host, Path: listNodesPath}
	req, err := http.NewRequest("GET", listNodesURL.String(), nil)
	if err != nil {
		return []Node{}, err
	}
	req.Header.Set("Cookie", c.Cookie)
	resp, err := c.Client.Do(req)
	if err != nil {
		return []Node{}, err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	nodesList, _ := ioutil.ReadAll(resp.Body)
	var n FabricNodes
	err = json.NewDecoder(bytes.NewReader(nodesList)).Decode(&n)
	if err != nil {
		return []Node{}, err
	}
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
