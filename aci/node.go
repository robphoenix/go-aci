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

// Node contains all the attributes of an ACI fabric node API request/response
type Node struct {
	Status       string `json:"status,omitempty"`
	DN           string `json:"dn,omitempty"`
	FabricStatus string `json:"fabricSt,omitempty"`
	ID           string `json:"id,omitempty"`
	Model        string `json:"model,omitempty"`
	Name         string `json:"name,omitempty"`
	NodeID       string `json:"nodeId,omitempty"`
	Role         string `json:"role,omitempty"`
	RN           string `json:"rn,omitempty"`
	Serial       string `json:"serial,omitempty"`
	Version      string `json:"version,omitempty"`
}

// String returns the string representation of an ACI node
func (n *Node) String() string {
	return fmt.Sprintf("%s %s %s", n.Name, n.ID, n.Serial)
}

// NewNode instanstatiates a valid ACI fabric membership node
func NewNode(name, id, serial string) (*Node, error) {
	// A valid serial number has a maximum length of 16
	// and contains only letters and numbers
	validSerial := regexp.MustCompile(`^[a-zA-Z0-9]{0,16}$`)
	if !validSerial.MatchString(serial) {
		return nil, fmt.Errorf("invalid serial number: %s can only contain letters and numbers and have a max length of 16", serial)
	}
	// Node ID must be between 101 and 4000
	i, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid node id: %s %v", id, err)
	}
	if i < 101 || i > 4000 {
		return nil, fmt.Errorf("out of range: %d node id must be between 101 & 4000", i)
	}
	return &Node{
		Name:   name,
		NodeID: id,
		Serial: serial,
		DN:     nodeDN + serial,
		RN:     nodeRN + serial,
	}, nil
}

// NewFabricNodeContainer instantiates a FabricNodeContainer
func NewFabricNodeContainer(n *Node) *FabricNodeContainer {
	return &FabricNodeContainer{FabricNode: FabricNode{Node: n}}
}

// NewNodeIdentProfContainer instantiates a FabricNodeIdentProfContainer
func NewNodeIdentProfContainer(nodes []Node, action string) *NodeIdentProfContainer {
	var children []*FabricNodeContainer
	for _, node := range nodes {
		node.Status = action
		children = append(children, NewFabricNodeContainer(&node))
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
	node.Status = createModify
	_, err := nodeDo(c, http.MethodPost, fmt.Sprintf(nodeAddPath, node.Serial), NewFabricNodeContainer(node))
	return err
}

// DeleteNode deletes a fabric membership node
func (c *Client) DeleteNode(node *Node) error {
	node.Status = delete
	_, err := nodeDo(c, http.MethodPost, nodeDeletePath, NewFabricNodeContainer(node))
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

func nodeDo(c *Client, method, URL string, payload interface{}) (NodesResponse, error) {
	var nr NodesResponse
	req, err := c.NewRequest(method, URL, payload)
	if err != nil {
		return nr, err
	}
	_, err = c.Do(req, &nr)
	return nr, err
}

// func (c *Client) DecomissionNode() error {
// https://supportforums.cisco.com/discussion/13296271/decommissioning-fabric-nodes-api
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
//         return nil
// }

// Key implements the Key method of the Mapper interface
func (n *Node) Key() string {
	return n.Serial + n.NodeID + n.Name
}

// Value implements the Value method of the Mapper interface
func (n *Node) Value() Node {
	return *n
}
