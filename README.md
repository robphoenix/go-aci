# go-aci
Go API wrapper for Cisco ACI

[![GoDoc](https://godoc.org/github.com/robphoenix/go-aci/aci?status.svg)](http://godoc.org/github.com/robphoenix/go-aci/aci)
[![Go Report Card](https://goreportcard.com/badge/github.com/robphoenix/go-aci)](https://goreportcard.com/report/github.com/robphoenix/go-aci)
[![stability-experimental](https://img.shields.io/badge/stability-experimental-orange.svg)](https://github.com/emersion/stability-badges#experimental)

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/robphoenix/go-aci/aci"
)

func main() {

	// create client
	client, err := aci.NewClient(aci.Config{
		Host:     "sandboxapicdc.cisco.com",
		Username: "admin",
		Password: "ciscopsdt",
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx := context.Background()

	// login
	err = client.Login(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	// define nodes
	node101, err := client.FabricMembership.NewNode(
		"leaf-101",    // name
		"101",         // id
		"1",           // pod id
		"FOC0849N1BD", // serial number
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	node102, err := client.FabricMembership.NewNode(
		"leaf-102",    // name
		"102",         // id
		"1",           // pod id
		"FOC0456N2BC", // serial number
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	// mark node for creation
    node101.SetCreate()

	// mark node for deletion
    node102.SetDelete()

    nodes := []aci.FabricMembership.Node{
        node101,
        node102,
    }

    // update ACI Fabric Membership
	resp, err := client.FabricMembership.Update(ctx, nodes...)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("resp = %+v\n", resp)

    // see a nodes current status (defaults to create)
    status101 := node101.Status()
    fmt.Println(status101) // Output: "created"
    status102 := node102.Status()
    fmt.Println(status102) // Output: "deleted"

    node102.SetCreate()

	resp, err := client.FabricMembership.Update(ctx, node102)
	if err != nil {
		log.Fatal(err)
		return
	}

	// list nodes
	nodes, err := client.FabricMembership.List(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	for i, node := range nodes {
		fmt.Printf("%d: %s\n", i+1, node)
	}
}
```
