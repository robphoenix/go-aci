package aci

import (
	"fmt"
	"regexp"
	"strconv"
)

// Node is an ACI fabric membership node
type Node struct {
	id     string
	name   string
	pod    string
	serial string
	role   string
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
// A node name can be up to 64 characters and can only contain
// alphanumeric, hyphen and underscore characters. It must begin
// and end with an alphanumeric character.
func (n *Node) SetName(name string) error {
	valid := regexp.MustCompile(`(^[a-zA-Z0-9]{1,2}$|^[a-zA-Z0-9][a-zA-Z0-9_-]{0,62}[a-zA-Z0-9]$)`)
	if !valid.MatchString(name) {
		return fmt.Errorf("invalid name: %s", name)
	}
	n.name = name
	return nil
}

// Pod returns the id of the pod the node is attached to.
func (n *Node) Pod() string {
	return n.pod
}

// SetPod validates and sets the id of the pod the node
// is attached to.
//
// A pod id must be a number between 0 and 255.
func (n *Node) SetPod(id string) error {
	// Pod id must be between 0 and 255
	idN, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid pod id: %s", id)
	}
	if idN < 1 || idN > 255 {
		return fmt.Errorf("invalid pod id: %s", id)
	}
	n.pod = id
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
	valid := regexp.MustCompile(`^[a-zA-Z0-9]{1,16}$`)
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

// SetRole sets the role of the node.
// Can only be "leaf" or "spine"
func (n *Node) SetRole(role string) error {
	if role != "leaf" && role != "spine" {
		return fmt.Errorf("invalid role: %s", role)
	}
	n.role = role
	return nil
}

// Role returns the role of the node.
func (n *Node) Role() string {
	return n.role
}
