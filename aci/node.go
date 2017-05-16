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
	Node     `json:"attributes"`
	Children []FabricNodeIdentPContainer `json:"children"`
}

// FabricNodeIdentP ...
type FabricNodeIdentP struct {
	Node `json:"attributes"`
}

// Node ...
type Node struct {
	Name   string `json:"name,omitempty"`
	ID     string `json:"nodeId,omitempty"`
	Serial string `json:"serial,omitempty"`
	Status string `json:"status,omitempty"`
	// Role is provisioned by ACI, we only need it when fetching node info
	// Role   string `json:"role,omitempty"`
}

// GetNodes ...
type GetNodes struct {
	Imdata []struct {
		FabricNode `json:"fabricNode"`
	} `json:"imdata"`
	TotalCount string `json:"totalCount"`
}

// FabricNode ...
type FabricNode struct {
	FabricNodeAttributes `json:"attributes"`
}

// FabricNodeAttributes ...
type FabricNodeAttributes struct {
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
	Role             string `json:"role,omitempty"`
	Serial           string `json:"serial,omitempty"`
	UID              string `json:"uid,omitempty"`
	Vendor           string `json:"vendor,omitempty"`
	Version          string `json:"version,omitempty"`
}

// AddNodes ...
func (c *Client) AddNodes(ns []Node) error {
	var children []FabricNodeIdentPContainer
	// add individual nodes to the nodes struct
	for _, n := range ns {
		n.Status = create
		fp := FabricNodeIdentPContainer{FabricNodeIdentP: FabricNodeIdentP{Node: n}}
		children = append(children, fp)
	}
	fpol := FabricNodeIdentPolContainer{
		FabricNodeIdentPol{
			Node: Node{
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
	fmt.Printf("b = %+v\n", b)

	// nodes endpoint
	nodesURL := url.URL{Scheme: c.Host.Scheme, Host: c.Host.Host, Path: nodesPath}
	req, err := http.NewRequest("POST", nodesURL.String(), b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
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
	// add individual nodes to the nodes struct
	for _, n := range ns {
		n.Status = delete
		fp := FabricNodeIdentPContainer{FabricNodeIdentP: FabricNodeIdentP{Node: n}}
		children = append(children, fp)
	}
	fpol := FabricNodeIdentPolContainer{
		FabricNodeIdentPol{
			Node: Node{
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
	fmt.Printf("b = %+v\n", b)

	// nodes endpoint
	nodesURL := url.URL{Scheme: c.Host.Scheme, Host: c.Host.Host, Path: nodesPath}
	req, err := http.NewRequest("POST", nodesURL.String(), b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
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

// ModifyNodes ...
func (c *Client) ModifyNodes(ns []Node) error {
	var children []FabricNodeIdentPContainer
	// add individual nodes to the nodes struct
	for _, n := range ns {
		n.Status = modify
		fp := FabricNodeIdentPContainer{FabricNodeIdentP: FabricNodeIdentP{Node: n}}
		children = append(children, fp)
	}
	fpol := FabricNodeIdentPolContainer{
		FabricNodeIdentPol{
			Node: Node{
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
	fmt.Printf("b = %+v\n", b)

	nodesURL := url.URL{Scheme: c.Host.Scheme, Host: c.Host.Host, Path: nodesPath}
	req, err := http.NewRequest("POST", nodesURL.String(), b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
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

// ListNodes ...
func (c *Client) ListNodes() (GetNodes, error) {
	listNodesURL := url.URL{Scheme: c.Host.Scheme, Host: c.Host.Host, Path: listNodesPath}
	req, err := http.NewRequest("GET", listNodesURL.String(), nil)
	if err != nil {
		return GetNodes{}, err
	}
	req.Header.Set("Cookie", c.Cookie)
	resp, err := c.Client.Do(req)
	if err != nil {
		return GetNodes{}, err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	nodesList, _ := ioutil.ReadAll(resp.Body)
	var n GetNodes
	err = json.NewDecoder(bytes.NewReader(nodesList)).Decode(&n)
	if err != nil {
		return GetNodes{}, err
	}
	return n, nil
}

func (c *Client) decomissionNode() error {
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
