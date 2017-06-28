# go-aci
Go API wrapper for Cisco ACI

[![GoDoc](https://godoc.org/github.com/robphoenix/go-aci/aci?status.svg)](http://godoc.org/github.com/robphoenix/go-aci/aci)
[![Go Report Card](https://goreportcard.com/badge/github.com/robphoenix/go-aci)](https://goreportcard.com/report/github.com/robphoenix/go-aci)

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/robphoenix/go-aci/aci"
)

func main() {

	// set config options
	cfg := aci.Config{
		Host:     "sandboxapicdc.cisco.com",
		Username: "admin",
		Password: "ciscopsdt",
	}

	// create client
	client := aci.NewClient(cfg)

	// login
	err := client.Login()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// create node
	name := "leaf-101"
	serial := "FOC0849N1BD"
	nodeID := "101"
	podID := "1"
	node, err := client.FabricMembership.NewNode(name, nodeID, podID, serial)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// add node
	err = client.FabricMembership.AddNode(node)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// // delete node
	// err = client.FabricMembership.DeleteNode(node)
	// if err != nil {
	//         log.Fatal(err)
	//         os.Exit(1)
	// }

	// list nodes
	nodes, err := client.FabricMembership.ListNodes()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	for i, node := range nodes {
		fmt.Printf("%d: %s\n", i+1, node)
	}
}
```
