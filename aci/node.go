package aci

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const (
	nodesPath           = "api/node/mo/uni/controller/nodeidentpol.json"
	nodeAddPath         = "api/node/mo/uni/controller/nodeidentpol/nodep-%s.json" // requires node serial
	nodeDeletePath      = "api/node/mo/uni/controller/nodeidentpol.json"
	nodeDecomissionPath = "api/node/mo/uni/fabric/outofsvc.json"
	nodeListPath        = "api/node/class/fabricNode.json"
	nodeDN              = "uni/controller/nodeidentpol/nodep-%s" // requires node serial
	nodeRN              = "nodep-%s"                             // requires node serial
	nodeTDN             = "topology/pod-%s/node-%s"              // requires pod id, node id
)

// Node is an ACI fabric membership node
type Node struct {
	ID     string
	Name   string
	PodID  string
	Serial string
}

// String returns the string representation of an ACI node
func (node *Node) String() string {
	return fmt.Sprintf("%s %s %s", node.Name, node.ID, node.Serial)
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
		ID:     nodeID,
		PodID:  podID,
		Serial: serial,
	}, nil
}

// NodeIdentProfContainer is a container for a NodeIdentityProfile
type NodeIdentProfContainer struct {
	NodeIdentityProfile `json:"fabricNodeIdentPol"`
}

// NodeIdentityProfile is a container for the node identity profile
type NodeIdentityProfile struct {
	NodeIdentProfAttributes `json:"attributes"`
	Children                []*FabricNodeContainer `json:"children"`
}

// NodeIdentProfAttributes contains all the attributes of an ACI fabric node API response
type NodeIdentProfAttributes struct {
	AdSt             string    `json:"adSt,omitempty"`
	ChildAction      string    `json:"childAction,omitempty"`
	DelayedHeartbeat string    `json:"delayedHeartbeat,omitempty"`
	DN               string    `json:"dn,omitempty"`
	FabricSt         string    `json:"fabricSt,omitempty"`
	ID               string    `json:"id,omitempty"`
	LcOwn            string    `json:"lcOwn,omitempty"`
	ModTs            time.Time `json:"modTs,omitempty"`
	Model            string    `json:"model,omitempty"`
	MonPolDn         string    `json:"monPolDn,omitempty"`
	Name             string    `json:"name,omitempty"`
	NameAlias        string    `json:"nameAlias,omitempty"`
	Role             string    `json:"role,omitempty"`
	Serial           string    `json:"serial,omitempty"`
	Status           string    `json:"status,omitempty"`
	UID              string    `json:"uid,omitempty"`
	Vendor           string    `json:"vendor,omitempty"`
	Version          string    `json:"version,omitempty"`
}

// FabricNodeContainer is a container for a Fabric Node
type FabricNodeContainer struct {
	FabricNode `json:"fabricNodeIdentP"`
}

// FabricNode is the node identity profile
type FabricNode struct {
	*FabricNodeIdentPAttributes `json:"attributes"`
}

// FabricNodeIdentPAttributes contains all the attributes of an ACI fabric node API request
type FabricNodeIdentPAttributes struct {
	DN     string `json:"dn,omitempty"`
	Serial string `json:"serial,omitempty"`
	NodeID string `json:"nodeId,omitempty"`
	Name   string `json:"name,omitempty"`
	Role   string `json:"role,omitempty"`
	RN     string `json:"rn,omitempty"`
	Status string `json:"status,omitempty"`
}

// NodesResponse contains the response for ACI fabric nodes requests
type NodesResponse struct {
	NodesImdata []NodesImdata `json:"imdata"`
}

// NodesImdata is describes the node in the nodes response structure
type NodesImdata struct {
	FabricNode `json:"fabricNode"`
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

// NewFabricNodeContainer instantiates a FabricNodeContainer
func NewFabricNodeContainer(node *Node, action string) *FabricNodeContainer {
	return &FabricNodeContainer{
		FabricNode: FabricNode{
			FabricNodeIdentPAttributes: &FabricNodeIdentPAttributes{
				Status: action,
				DN:     fmt.Sprintf(nodeDN, node.Serial),
				RN:     fmt.Sprintf(nodeRN, node.Serial),
				Name:   node.Name,
				NodeID: node.ID,
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
			NodeIdentProfAttributes: NodeIdentProfAttributes{
				Status: createModify,
			},
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
	nr, err := nodeDo(c, http.MethodGet, nodeListPath, nil)
	if err != nil {
		return nil, fmt.Errorf("list nodes: %v", err)
	}

	var ns []*Node
	for _, node := range nr.NodesImdata {
		n := &Node{
			Name:   node.Name,
			ID:     node.NodeID,
			Serial: node.Serial,
		}
		ns = append(ns, n)
	}
	return ns, nil
}

// DecommissionNode decommisions a fabric membership node
func (c *Client) DecommissionNode(node *Node) error {
	payload := DecommissionNodeContainer{
		DecommissionNode: DecommissionNode{
			DecommissionAttributes: DecommissionAttributes{
				TDN:                  fmt.Sprintf(nodeTDN, node.PodID, node.ID),
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
func (node *Node) Key() string {
	return node.Serial + node.ID + node.Name
}

// Value implements the Value method of the Mapper interface
func (node *Node) Value() *Node {
	return node
}
