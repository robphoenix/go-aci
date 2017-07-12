package aci

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// Node is an ACI fabric membership node
type Node struct {
	ID     string
	Name   string
	PodID  string
	Serial string
}

// String returns the string representation of an ACI node
func (n *Node) String() string {
	return fmt.Sprintf("%s %s %s", n.Name, n.ID, n.Serial)
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
func newFabricNodeContainer(n *Node, action string) *FabricNodeContainer {
	rn := fmt.Sprintf("nodep-%s", n.Serial)
	dn := fmt.Sprintf("uni/controller/nodeidentpol/%s", rn)
	return &FabricNodeContainer{
		FabricNodeIdentP: FabricNodeIdentP{
			FabricNodeIdentPAttributes: &FabricNodeIdentPAttributes{
				Status: action,
				DN:     dn,
				RN:     rn,
				Name:   n.Name,
				NodeID: n.ID,
				Serial: n.Serial,
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

// New instanstatiates a valid ACI fabric membership node
func (s *FabricMembershipService) New(name, nodeID, podID, serial string) (*Node, error) {
	// A valid serial number has a maximum length of 16
	// and contains only letters and numbers
	validSerial := regexp.MustCompile(`^[a-zA-Z0-9]{0,16}$`)
	if !validSerial.MatchString(serial) {
		return nil, fmt.Errorf("invalid serial number: %s can only contain letters and numbers and have a max length of 16", serial)
	}
	// Node ID must be between 101 and 4000
	n, err := strconv.Atoi(nodeID)
	if err != nil {
		return nil, fmt.Errorf("invalid node id: %s %v", nodeID, err)
	}
	if n < 101 || n > 4000 {
		return nil, fmt.Errorf("out of range: %d node id must be between 101 & 4000", n)
	}
	// Pod ID must be between 0 and 255
	p, err := strconv.Atoi(podID)
	if err != nil {
		return nil, fmt.Errorf("invalid pod id: %s %v", podID, err)
	}
	if p < 0 || p > 255 {
		return nil, fmt.Errorf("out of range: %d pod id must be between 0 & 255", p)
	}
	return &Node{
		Name:   name,
		ID:     nodeID,
		PodID:  podID,
		Serial: serial,
	}, nil
}

// Add adds a single node to the ACI fabric membership
func (s *FabricMembershipService) Add(ctx context.Context, n *Node) (NodesResponse, error) {
	path := fmt.Sprintf("api/node/mo/uni/controller/nodeidentpol/nodep-%s.json", n.Serial)

	var nr NodesResponse

	req, err := s.client.NewRequest(http.MethodPost, path, newFabricNodeContainer(n, createModify))
	if err != nil {
		return nr, err
	}

	_, err = s.client.Do(ctx, req, &nr)
	if err != nil {
		return nr, err
	}

	return nr, nil
}

// Delete deletes a fabric membership node
func (s *FabricMembershipService) Delete(ctx context.Context, n *Node) (NodesResponse, error) {
	path := "api/node/mo/uni/controller/nodeidentpol.json"

	var nr NodesResponse

	req, err := s.client.NewRequest(http.MethodPost, path, newFabricNodeContainer(n, delete))
	if err != nil {
		return nr, err
	}

	_, err = s.client.Do(ctx, req, &nr)
	if err != nil {
		return nr, err
	}

	return nr, nil
}

// List lists all node members of the ACI fabric
func (s *FabricMembershipService) List(ctx context.Context) ([]*Node, error) {
	path := "api/node/class/fabricNode.json"

	var ns []*Node
	var nr NodesResponse

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("list nodes: %v", err)
	}

	_, err = s.client.Do(ctx, req, &nr)
	if err != nil {
		return nil, fmt.Errorf("list nodes: %v", err)
	}

	for _, node := range nr.NodesImdata {
		n := &Node{
			Name:   node.Name,
			ID:     node.ID,
			Serial: node.Serial,
		}
		ns = append(ns, n)
	}
	return ns, nil
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

// Decommission decommisions a fabric membership node
func (s *FabricMembershipService) Decommission(ctx context.Context, node *Node) (NodesResponse, error) {
	tdn := fmt.Sprintf("topology/pod-%s/node-%s", node.PodID, node.ID)
	payload := DecommissionNodeContainer{
		DecommissionNode: DecommissionNode{
			DecommissionAttributes: DecommissionAttributes{
				TDN:                  tdn,
				Status:               createModify,
				RemoveFromController: "true",
			},
		},
	}
	path := "api/node/mo/uni/fabric/outofsvc.json"

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
