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
