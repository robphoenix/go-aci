package aci

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	nodesPath     = "api/node/mo/uni/controller/nodeidentpol.json"
	nodeAddPath   = "api/node/mo/uni/controller/nodeidentpol/nodep-%s.json" // requires node serial
	listNodesPath = "api/node/class/fabricNode.json"
)

// add
// url: https://sandboxapicdc.cisco.com/api/node/mo/uni/controller/nodeidentpol/nodep-serial1.json
// payload{"fabricNodeIdentP":{"attributes":{"dn":"uni/controller/nodeidentpol/nodep-serial1","serial":"serial1","nodeId":"110","name":"leaf-110","role":"leaf","rn":"nodep-serial1","status":"created"},"children":[]}}
// delete
// url: https://sandboxapicdc.cisco.com/api/node/mo/uni/controller/nodeidentpol.json
// payload{"fabricNodeIdentP":{"attributes":{"dn":"uni/controller/nodeidentpol/nodep-serial1","status":"deleted"},"children":[]}}

// Node is a member of an ACI fabric
type Node struct {
	FabricStatus string
	Model        string
	Name         string
	id           string
	serial       string
	Status       string
	Role         string
}

// ID returns the node ID
func (n *Node) ID() string {
	return n.id
}

// SetID validates and sets the node id
func (n *Node) SetID(s string) error {
	i, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("invalid node id: %d %v", i, err)
	}
	// Node ID must be between 101 and 4000
	if i < 101 || i > 4000 {
		return fmt.Errorf("node id: %d is out of range, must be between 101 & 4000", i)
	}
	n.id = s
	return nil
}

// Serial returns the node serial number
func (n *Node) Serial() string {
	return n.serial
}

// SetSerial validates and sets the node serial number
func (n *Node) SetSerial(s string) error {
	// A valid serial number has a maximum length of 16
	// and contains only letters and numbers
	validSerial := regexp.MustCompile(`^[a-zA-Z0-9]{0,16}$`)
	if !validSerial.MatchString(s) {
		return fmt.Errorf("invalid serial number: %s can only contain letters and numbers and have a max length of 16", s)
	}
	n.serial = strings.ToUpper(s)
	return nil
}

// NewNode instantiates a valid ACI fabric membership node
func NewNode(name, id, serial string) (*Node, error) {
	node := &Node{Name: name}
	err := node.SetID(id)
	if err != nil {
		return node, fmt.Errorf("cannot create node: %v", err)
	}
	err = node.SetSerial(serial)
	if err != nil {
		return node, fmt.Errorf("cannot create node: %v", err)
	}
	return node, nil
}

// NIPContainer is a container for a NodeIdentityProfile
type NIPContainer struct {
	NodeIdentityProfile `json:"fabricNodeIdentPol"`
}

// NodeIdentityProfile is a container for the node identity profile
type NodeIdentityProfile struct {
	NodeAttributes `json:"attributes"`
	Children       []FNContainer `json:"children"`
}

// FNContainer is a container for a Fabric Node
type FNContainer struct {
	FabricNode `json:"fabricNodeIdentP"`
}

// FabricNode is the node identity profile
type FabricNode struct {
	NodeAttributes `json:"attributes"`
}

// nodesResponse contains the response for ACI fabric nodes requests
type nodesResponse struct {
	NodesImdata []NodesImdata `json:"imdata"`
}

// NodesImdata is describes the node in the nodes response structure
type NodesImdata struct {
	FabricNode `json:"fabricNode"`
}

// NodeAttributes contains all the attributes of an ACI fabric node API request/response
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
	Rn               string `json:"rn,omitempty"`
	Serial           string `json:"serial,omitempty"`
	UID              string `json:"uid,omitempty"`
	Vendor           string `json:"vendor,omitempty"`
	Version          string `json:"version,omitempty"`
}

func newNIPContainer(ns []Node, action string) NIPContainer {
	var fs []FNContainer
	for _, n := range ns {
		var f FNContainer
		f.Name = n.Name
		f.NodeID = n.ID()
		f.Serial = n.Serial()
		f.Status = action
		fs = append(fs, f)
	}
	var nc NIPContainer
	nc.Status = createModify
	nc.Children = fs
	return nc
}

// editNodes takes a createModify or delete action and performs the
// necessary API request
func editNodes(c *Client, ns []Node, action string) error {
	nip := newNIPContainer(ns, action)
	req, err := c.NewRequest(http.MethodPost, nodesPath, nip)
	if err != nil {
		return err
	}

	var f nodesResponse

	_, err = c.Do(req, &f)
	return err
}

// AddNode adds a single node to the ACI fabric membership
func (c *Client) AddNode(n Node) error {
	var f FNContainer
	f.Name = n.Name
	f.NodeID = n.ID()
	f.Serial = n.Serial()
	f.Status = createModify
	f.Dn = "uni/controller/nodeidentpol/nodep-" + n.Serial()
	f.Rn = "nodep-" + n.Serial()

	path := fmt.Sprintf(nodeAddPath, n.Serial())

	req, err := c.NewRequest(http.MethodPost, path, f)
	if err != nil {
		return err
	}

	var nr nodesResponse

	_, err = c.Do(req, &nr)
	return err
}

// CreateNodes adds a slice of nodes to the ACI fabric membership
func (c *Client) CreateNodes(ns []Node) error {
	err := editNodes(c, ns, createModify)
	if err != nil {
		return err
	}
	return nil
}

// DeleteNodes deletes a slice of nodes from the ACI fabric membership
func (c *Client) DeleteNodes(ns []Node) error {
	err := editNodes(c, ns, delete)
	if err != nil {
		return err
	}
	return nil
}

// ListNodes lists all node members of the ACI fabric
func (c *Client) ListNodes() ([]Node, error) {
	req, err := c.NewRequest(http.MethodGet, listNodesPath, nil)
	if err != nil {
		return nil, fmt.Errorf("list nodes: %v", err)
	}

	var n nodesResponse
	_, err = c.Do(req, &n)
	if err != nil {
		return nil, fmt.Errorf("list nodes: %v", err)
	}

	var ns []Node
	for _, v := range n.NodesImdata {
		node := Node{
			FabricStatus: v.FabricSt,
			Model:        v.Model,
			Name:         v.Name,
			Role:         v.Role,
			Status:       v.Status,
		}
		// we don't need to check validity
		// as it's coming from the APIC
		_ = node.SetID(v.ID)
		_ = node.SetSerial(v.Serial)
		ns = append(ns, node)
	}
	return ns, nil
}

// func (c *Client) decomissionNode() error {
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
//         return nil
// }

// Key implements the Key method of the Mapper interface
func (n *Node) Key() string {
	return n.Serial() + n.ID() + n.Name
}

// Value implements the Value method of the Mapper interface
func (n *Node) Value() Node {
	return *n
}
