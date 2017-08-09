package aci

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Node is an ACI fabric membership node
type Node struct {
	id     string
	name   string
	podID  string
	serial string
	status string
}

// ID returns the node ID
func (n *Node) ID() string {
	return n.id
}

// SetID validates and sets the node ID.
//
// A node ID must be a number between 101 and 4000 inclusive.
func (n *Node) SetID(id string) error {
	// Node id must be between 101 and 4000
	idN, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid node id: %s", id)
	}
	if idN < 101 || idN > 4000 {
		return fmt.Errorf("invalid node id: %s", id)
	}
	n.id = id
	return nil
}

// Name returns the node name.
func (n *Node) Name() string {
	return n.name
}

// SetName validates and sets the node name.
//
// A node name can be no longer than 64 characters
// and can contain only letters, numbers, hyphen
// and underscore. It cannot end with a hyphen or
// underscore character
func (n *Node) SetName(name string) error {
	valid := regexp.MustCompile(`^[a-zA-Z0-9_-]{0,64}$`)
	if !valid.MatchString(name) {
		return fmt.Errorf("invalid name: %s", name)
	}
	if strings.HasSuffix(name, "-") || strings.HasSuffix(name, "_") {
		return fmt.Errorf("invalid name: %s", name)
	}
	n.name = name
	return nil
}

// PodID returns the id of the pod the node is attached to.
func (n *Node) PodID() string {
	return n.podID
}

// SetPodID validates and sets the id of the pod the node
// is attached to.
//
// A pod id must be a number between 0 and 255.
func (n *Node) SetPodID(id string) error {
	// Pod id must be between 0 and 255
	idN, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid pod id: %s", id)
	}
	if idN < 0 || idN > 255 {
		return fmt.Errorf("invalid pod id: %s", id)
	}
	n.podID = id
	return nil
}

// Serial returns the node's serial number.
func (n *Node) Serial() string {
	return n.serial
}

// SetSerial validates and sets the node serial number.
//
// A valid serial number has a maximum length of 16
// and contains only letters and numbers.
func (n *Node) SetSerial(serial string) error {
	valid := regexp.MustCompile(`^[a-zA-Z0-9]{0,16}$`)
	if !valid.MatchString(serial) {
		return fmt.Errorf("invalid serial number: %s", serial)
	}
	n.serial = serial
	return nil
}

// SetCreated sets the status of the node to "created,modified".
func (n *Node) SetCreated() string {
	n.status = createdModified
	return n.status
}

// SetDeleted sets the status of the node to "deleted".
func (n *Node) SetDeleted() string {
	n.status = deleted
	return n.status
}

// Status returns the status of the node.
func (n *Node) Status() string {
	return n.status
}

// String returns the string representation of an ACI node
func (n *Node) String() string {
	return fmt.Sprintf("%s %s %s", n.name, n.id, n.serial)
}

// NodesResponse contains the response for ACI fabric nodes requests
type NodesResponse struct {
	NodesImdata []NodesImdata `json:"imdata"`
}

// NodesImdata describes the node in the nodes response structure
type NodesImdata struct {
	FabricNode `json:"fabricNode"`
}

// FabricNode is the node identity profile
type FabricNode struct {
	*NodeIdentProfAttributes `json:"attributes"`
}

// NodeIdentProfContainer is a container for a NodeIdentityProfile
type NodeIdentProfContainer struct {
	NodeIdentityProfile `json:"fabricNodeIdentPol"`
}

// NodeIdentityProfile describes the node identity profile
type NodeIdentityProfile struct {
	NodeIdentProfAttributes `json:"attributes"`
	Children                []FabricNodeContainer `json:"children"`
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

// FabricNodeContainer is a container for a Fabric Node Identity Profile
type FabricNodeContainer struct {
	FabricNodeIdentP `json:"fabricNodeIdentP"`
}

// newFabricNodeContainer instantiates a FabricNodeContainer
func newFabricNodeContainer(n Node) FabricNodeContainer {
	rn := fmt.Sprintf("nodep-%s", n.Serial())
	dn := fmt.Sprintf("uni/controller/nodeidentpol/%s", rn)
	return FabricNodeContainer{
		FabricNodeIdentP: FabricNodeIdentP{
			FabricNodeIdentPAttributes: &FabricNodeIdentPAttributes{
				Status: n.Status(),
				DN:     dn,
				RN:     rn,
				Name:   n.Name(),
				NodeID: n.ID(),
				Serial: n.Serial(),
			},
		},
	}
}

// FabricNodeIdentP is the node identity profile
type FabricNodeIdentP struct {
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

// FabricMembershipService handles communication with the fabric membership related
// methods of the APIC API.
type FabricMembershipService service

// NewNode instanstatiates a valid ACI fabric membership node
func (s *FabricMembershipService) NewNode(name, nodeID, podID, serial string) (*Node, error) {
	node := &Node{}
	err := node.SetID(nodeID)
	if err != nil {
		return node, err
	}
	err = node.SetSerial(serial)
	if err != nil {
		return node, err
	}
	err = node.SetPodID(podID)
	if err != nil {
		return node, err
	}
	return &Node{
		name:   name,
		id:     nodeID,
		podID:  podID,
		serial: serial,
	}, nil
}

func (s *FabricMembershipService) Update(ctx context.Context, nodes ...Node) (NodesResponse, error) {
	payload := newNodeIdentProfContainer(nodes)
	var nr NodesResponse
	req, err := s.client.NewRequest(http.MethodPost, "api/node/mo/uni/controller/nodeidentpol.json", payload)
	if err != nil {
		return nr, err
	}
	_, err = s.client.Do(ctx, req, &nr)
	return nr, err
}

// // AddNode adds a single node to the ACI fabric membership
// func (s *FabricMembershipService) AddNode(ctx context.Context, n *Node) (NodesResponse, error) {
//
//         path := fmt.Sprintf("api/node/mo/uni/controller/nodeidentpol/nodep-%s.json", n.serial)
//
//         var nr NodesResponse
//
//         req, err := s.client.NewRequest(http.MethodPost, path, newFabricNodeContainer(n, createModify))
//         if err != nil {
//                 return nr, err
//         }
//
//         _, err = s.client.Do(ctx, req, &nr)
//         if err != nil {
//                 return nr, err
//         }
//
//         return nr, nil
// }

func newNodeIdentProfContainer(nodes []Node) NodeIdentProfContainer {
	var children []FabricNodeContainer
	for _, node := range nodes {
		children = append(children, newFabricNodeContainer(node))
	}
	return NodeIdentProfContainer{
		NodeIdentityProfile: NodeIdentityProfile{
			NodeIdentProfAttributes: NodeIdentProfAttributes{
				Status: createdModified,
			},
			Children: children,
		},
	}
}

//
// // DeleteNode deletes a fabric membership node
// func (s *FabricMembershipService) DeleteNode(ctx context.Context, n *Node) (NodesResponse, error) {
//
//         path := "api/node/mo/uni/controller/nodeidentpol.json"
//
//         var nr NodesResponse
//
//         req, err := s.client.NewRequest(http.MethodPost, path, newFabricNodeContainer(n, delete))
//         if err != nil {
//                 return nr, err
//         }
//
//         _, err = s.client.Do(ctx, req, &nr)
//         if err != nil {
//                 return nr, err
//         }
//
//         return nr, nil
// }

// List lists all node members of the ACI fabric
func (s *FabricMembershipService) List(ctx context.Context) ([]*Node, error) {

	path := "api/node/class/fabricNode.json"

	var ns []*Node
	var nr NodesResponse

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list: %v", err)
	}

	_, err = s.client.Do(ctx, req, &nr)
	if err != nil {
		return nil, fmt.Errorf("list: %v", err)
	}

	for _, n := range nr.NodesImdata {
		ns = append(ns, &Node{
			name:   n.Name,
			id:     n.ID,
			serial: n.Serial,
			status: n.Status,
		})
	}
	return ns, nil
}

// NodeDecommissionContainer is a container for
// the request to decommission a fabric membership node
type NodeDecommissionContainer struct {
	NodeDecommission `json:"fabricRsDecommissionNode"`
}

// NodeDecommission describes the node to decommission
type NodeDecommission struct {
	DecommissionAttributes `json:"attributes"`
}

// DecommissionAttributes are the attributes of the node to be decomissioned
type DecommissionAttributes struct {
	TDN                  string `json:"tDn"`                  // "topology/pod-<podID>/node-<nodeID>"
	Status               string `json:"status"`               // createModify
	RemoveFromController string `json:"removeFromController"` // "true"
}

// DecommissionNode decommisions a fabric membership node
func (s *FabricMembershipService) DecommissionNode(ctx context.Context, node *Node) (NodesResponse, error) {

	path := "api/node/mo/uni/fabric/outofsvc.json"

	tdn := fmt.Sprintf("topology/pod-%s/node-%s", node.podID, node.id)
	payload := NodeDecommissionContainer{
		NodeDecommission: NodeDecommission{
			DecommissionAttributes: DecommissionAttributes{
				TDN:                  tdn,
				Status:               createdModified,
				RemoveFromController: "true",
			},
		},
	}

	var nr NodesResponse

	req, err := s.client.NewRequest(http.MethodPost, path, payload)
	if err != nil {
		return nr, err
	}

	_, err = s.client.Do(ctx, req, &nr)
	if err != nil {
		return nr, err
	}
	return nr, nil
}
