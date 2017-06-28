package aci

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

const (
	nodesPath           = "api/node/mo/uni/controller/nodeidentpol.json"
	nodeAddPath         = "api/node/mo/uni/controller/nodeidentpol/nodep-%s.json" // requires node serial
	nodeDeletePath      = "api/node/mo/uni/controller/nodeidentpol.json"
	nodeDecomissionPath = "api/node/mo/uni/fabric/outofsvc.json"
	listNodesPath       = "api/node/class/fabricNode.json"
	nodeDN              = "uni/controller/nodeidentpol/nodep-"
	nodeRN              = "nodep-"
)

// NodeIdentProfContainer is a container for a NodeIdentityProfile
type NodeIdentProfContainer struct {
	NodeIdentityProfile `json:"fabricNodeIdentPol"`
}

// NodeIdentityProfile is a container for the node identity profile
type NodeIdentityProfile struct {
	Node     `json:"attributes"`
	Children []*FabricNodeContainer `json:"children"`
}

// FabricNodeContainer is a container for a Fabric Node
type FabricNodeContainer struct {
	FabricNode `json:"fabricNodeIdentP"`
}

// NodesResponse contains the response for ACI fabric nodes requests
type NodesResponse struct {
	NodesImdata []NodesImdata `json:"imdata"`
}

// NodesImdata is describes the node in the nodes response structure
type NodesImdata struct {
	FabricNode `json:"fabricNode"`
}

// FabricNode is the node identity profile
type FabricNode struct {
	*Node `json:"attributes"`
}

// DecommissionNodeContainer is a container for
// the request to decommission a fabric membership node
type DecommissionNodeContainer struct {
	DecommissionNode `json:"fabricRsDecommissionNode"`
}

// DecommissionNode describes the node to decommission
type DecommissionNode struct {
	DecommissionAttributes `json:"attributes"`
}

// DecommissionAttributes are the attributes of the node to be decomissioned
type DecommissionAttributes struct {
	TDN                  string `json:"tDn"`                  // "topology/pod-<podID>/node-<nodeID>"
	Status               string `json:"status"`               // createModify
	RemoveFromController string `json:"removeFromController"` // "true"
}

// Node contains all the attributes of an ACI fabric node API request/response
type Node struct {
	Status       string `json:"status,omitempty"`
	DN           string `json:"dn,omitempty"`
	FabricStatus string `json:"fabricSt,omitempty"`
	ID           string `json:"id,omitempty"`
	Model        string `json:"model,omitempty"`
	Name         string `json:"name,omitempty"`
	NodeID       string `json:"nodeId,omitempty"`
	PodID        string `json:"podId,omitempty"`
	Role         string `json:"role,omitempty"`
	RN           string `json:"rn,omitempty"`
	Serial       string `json:"serial,omitempty"`
	Version      string `json:"version,omitempty"`
}

// // FabricNodeAttributes ...
// type FabricNodeAttributes struct {
//         AdSt             string    `json:"adSt"`
//         ChildAction      string    `json:"childAction"`
//         DelayedHeartbeat string    `json:"delayedHeartbeat"`
//         Dn               string    `json:"dn"`
//         FabricSt         string    `json:"fabricSt"`
//         ID               string    `json:"id"`
//         LcOwn            string    `json:"lcOwn"`
//         ModTs            time.Time `json:"modTs"`
//         Model            string    `json:"model"`
//         MonPolDn         string    `json:"monPolDn"`
//         Name             string    `json:"name"`
//         NameAlias        string    `json:"nameAlias"`
//         Role             string    `json:"role"`
//         Serial           string    `json:"serial"`
//         Status           string    `json:"status"`
//         UID              string    `json:"uid"`
//         Vendor           string    `json:"vendor"`
//         Version          string    `json:"version"`
// } // `json:"attributes"`
//
// // FabricNodeIdentPAttributes ...
// type FabricNodeIdentPAttributes struct {
//         Dn     string `json:"dn"`
//         Serial string `json:"serial"`
//         NodeID string `json:"nodeId"`
//         Name   string `json:"name"`
//         Role   string `json:"role"`
//         Rn     string `json:"rn"`
//         Status string `json:"status"`
// } // `json:"attributes"`

// String returns the string representation of an ACI node
func (n *Node) String() string {
	return fmt.Sprintf("%s %s %s", n.Name, n.ID, n.Serial)
}

// NewNode instanstatiates a valid ACI fabric membership node
func NewNode(name, nodeID, podID, serial string) (*Node, error) {
	// A valid serial number has a maximum length of 16
	// and contains only letters and numbers
	validSerial := regexp.MustCompile(`^[a-zA-Z0-9]{0,16}$`)
	if !validSerial.MatchString(serial) {
		return nil, fmt.Errorf("invalid serial number: %s can only contain letters and numbers and have a max length of 16", serial)
	}
	// Node ID must be between 101 and 4000
	i, err := strconv.Atoi(nodeID)
	if err != nil {
		return nil, fmt.Errorf("invalid node id: %s %v", nodeID, err)
	}
	if i < 101 || i > 4000 {
		return nil, fmt.Errorf("out of range: %d node id must be between 101 & 4000", i)
	}
	// Pod ID must be between 0 and 255
	j, err := strconv.Atoi(podID)
	if err != nil {
		return nil, fmt.Errorf("invalid pod id: %s %v", podID, err)
	}
	if j < 0 || j > 255 {
		return nil, fmt.Errorf("out of range: %d pod id must be between 0 & 255", j)
	}
	return &Node{
		Name:   name,
		NodeID: nodeID,
		PodID:  podID,
		Serial: serial,
	}, nil
}

// NewFabricNodeContainer instantiates a FabricNodeContainer
func NewFabricNodeContainer(node *Node, action string) *FabricNodeContainer {
	return &FabricNodeContainer{
		FabricNode: FabricNode{
			Node: &Node{
				Status: action,
				DN:     nodeDN + node.Serial,
				RN:     nodeRN + node.Serial,
				Name:   node.Name,
				NodeID: node.NodeID,
				PodID:  node.PodID,
				Serial: node.Serial,
			},
		},
	}
}

// NewNodeIdentProfContainer instantiates a FabricNodeIdentProfContainer
func NewNodeIdentProfContainer(nodes []*Node, action string) *NodeIdentProfContainer {
	var children []*FabricNodeContainer
	for _, node := range nodes {
		children = append(children, NewFabricNodeContainer(node, action))
	}
	return &NodeIdentProfContainer{
		NodeIdentityProfile: NodeIdentityProfile{
			Node:     Node{Status: createModify},
			Children: children,
		},
	}
}

// AddNode adds a single node to the ACI fabric membership
func (c *Client) AddNode(node *Node) error {
	_, err := nodeDo(c, http.MethodPost, fmt.Sprintf(nodeAddPath, node.Serial), NewFabricNodeContainer(node, createModify))
	return err
}

// DeleteNode deletes a fabric membership node
func (c *Client) DeleteNode(node *Node) error {
	_, err := nodeDo(c, http.MethodPost, nodeDeletePath, NewFabricNodeContainer(node, delete))
	return err
}

// AddNodes adds a slice of nodes to the ACI fabric membership
func (c *Client) AddNodes(ns []*Node) error {
	_, err := nodeDo(c, http.MethodPost, nodesPath, NewNodeIdentProfContainer(ns, createModify))
	return err
}

// DeleteNodes deletes a slice of nodes from the ACI fabric membership
func (c *Client) DeleteNodes(ns []*Node) error {
	_, err := nodeDo(c, http.MethodPost, nodesPath, NewNodeIdentProfContainer(ns, delete))
	return err
}

// ListNodes lists all node members of the ACI fabric
func (c *Client) ListNodes() ([]*Node, error) {
	nr, err := nodeDo(c, http.MethodGet, listNodesPath, nil)
	if err != nil {
		return nil, fmt.Errorf("list nodes: %v", err)
	}

	var ns []*Node
	for _, node := range nr.NodesImdata {
		ns = append(ns, node.FabricNode.Node)
	}
	return ns, nil
}

// DecommissionNode decommisions a fabric membership node
func (c *Client) DecommissionNode(node *Node) error {
	payload := DecommissionNodeContainer{
		DecommissionNode: DecommissionNode{
			DecommissionAttributes: DecommissionAttributes{
				TDN:                  "topology/pod-" + node.PodID + "/node-" + node.NodeID,
				Status:               createModify,
				RemoveFromController: "true",
			},
		},
	}
	_, err := nodeDo(c, http.MethodPost, nodeDecomissionPath, payload)
	return err
}

func nodeDo(c *Client, method, URL string, payload interface{}) (NodesResponse, error) {
	var nr NodesResponse
	req, err := c.NewRequest(method, URL, payload)
	if err != nil {
		return nr, err
	}
	_, err = c.Do(req, &nr)
	return nr, err
}

// Key implements the Key method of the Mapper interface
func (n *Node) Key() string {
	return n.Serial + n.NodeID + n.Name
}

// Value implements the Value method of the Mapper interface
func (n *Node) Value() Node {
	return *n
}
