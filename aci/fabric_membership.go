package aci

import (
	"context"
	"fmt"
	"net/http"
)

// NodesResponse contains the response for ACI fabric nodes requests
type NodesResponse struct {
	NodesImdata []NodesImdata `json:"imdata"`
}

// NodesImdata contains the node in the nodes response structure
type NodesImdata struct {
	FabricNode `json:"fabricNode"`
}

// FabricNode is the node identity profile
type FabricNode struct {
	NodeResponseAttrs `json:"attributes"`
}

// NodeIdentProfContainer is a container for a NodeIdentityProfile
type NodeIdentProfContainer struct {
	NodeIdentityProfile `json:"fabricNodeIdentPol"`
}

// NodeIdentityProfile describes the node identity profile
type NodeIdentityProfile struct {
	NodeResponseAttrs `json:"attributes"`
	Children          []FabricNodeContainer `json:"children"`
}

// NodeResponseAttrs contains all the attributes of an ACI fabric node API response
type NodeResponseAttrs struct {
	DN       string `json:"dn,omitempty"`
	FabricSt string `json:"fabricSt,omitempty"`
	ID       string `json:"id,omitempty"`
	Model    string `json:"model,omitempty"`
	Name     string `json:"name,omitempty"`
	Role     string `json:"role,omitempty"`
	Serial   string `json:"serial,omitempty"`
	Status   string `json:"status,omitempty"`
	UID      string `json:"uid,omitempty"`
	Vendor   string `json:"vendor,omitempty"`
	Version  string `json:"version,omitempty"`
}

// FabricNodeContainer is a container for a Fabric Node Identity Profile
type FabricNodeContainer struct {
	FabricNodeIdentP `json:"fabricNodeIdentP"`
}

// FabricNodeIdentP is the node identity profile
type FabricNodeIdentP struct {
	NodeRequestAttrs `json:"attributes"`
}

// NodeRequestAttrs contains all the attributes of an ACI fabric node API request
type NodeRequestAttrs struct {
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
func (s *FabricMembershipService) NewNode(name, ID, pod, serial, role string) (*Node, error) {
	node := &Node{}

	if err := node.SetName(name); err != nil {
		return node, err
	}

	if err := node.SetID(ID); err != nil {
		return node, err
	}

	if err := node.SetSerial(serial); err != nil {
		return node, err
	}

	if err := node.SetPod(pod); err != nil {
		return node, err
	}

	if err := node.SetRole(role); err != nil {
		return node, err
	}

	return node, nil
}

// Update ...
func (s *FabricMembershipService) Update(ctx context.Context, nodes ...*Node) (NodesResponse, error) {

	path := "api/node/mo/uni/controller/nodeidentpol.json"
	payload := newNodeIdentProfContainer(nodes)

	var nr NodesResponse

	req, err := s.client.NewRequest(http.MethodPost, path, payload)
	if err != nil {
		return nr, err
	}

	_, err = s.client.Do(ctx, req, &nr)
	return nr, err
}

func newNodeIdentProfContainer(nodes []*Node) NodeIdentProfContainer {

	var children []FabricNodeContainer

	for _, node := range nodes {
		rn := fmt.Sprintf("nodep-%s", node.Serial())
		dn := fmt.Sprintf("uni/controller/nodeidentpol/%s", rn)
		child := FabricNodeContainer{
			FabricNodeIdentP: FabricNodeIdentP{
				NodeRequestAttrs: NodeRequestAttrs{
					Status: node.Status(),
					DN:     dn,
					RN:     rn,
					Name:   node.Name(),
					NodeID: node.ID(),
					Serial: node.Serial(),
				},
			},
		}
		children = append(children, child)
	}

	return NodeIdentProfContainer{
		NodeIdentityProfile: NodeIdentityProfile{
			NodeResponseAttrs: NodeResponseAttrs{
				Status: createdModified,
			},
			Children: children,
		},
	}
}

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
		// we shouldn't need to validate nodes
		// returned by the APIC server, as APIC
		// should have already done this.
		ns = append(ns, &Node{
			name:   n.Name,
			id:     n.ID,
			role:   n.Role,
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
	TDN                  string `json:"tDn"`                  // "topology/pod-<pod>/node-<nodeID>"
	Status               string `json:"status"`               // createModify
	RemoveFromController string `json:"removeFromController"` // "true"
}

// DecommissionNode decommisions a fabric membership node
func (s *FabricMembershipService) DecommissionNode(ctx context.Context, node *Node) (NodesResponse, error) {

	path := "api/node/mo/uni/fabric/outofsvc.json"

	tdn := fmt.Sprintf("topology/pod-%s/node-%s", node.pod, node.id)
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
