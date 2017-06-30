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

	// create client
	client := aci.NewClient(aci.Config{
		Host:     "sandboxapicdc.cisco.com",
		Username: "admin",
		Password: "ciscopsdt",
	})

	ctx := context.Background()

	// login
	err := client.Login(ctx)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// create node
	node, err := client.FabricMembership.New(
		"leaf-101",    // name
		"101",         // id
		"1",           // pod id
		"FOC0849N1BD", // serial number
	)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// add node
	err = client.FabricMembership.Add(ctx, node)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// // delete node
	// err = client.FabricMembership.Delete(ctx, node)
	// if err != nil {
	//         log.Fatal(err)
	//         os.Exit(1)
	// }

	// list nodes
	nodes, err := client.FabricMembership.List(ctx)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	for i, node := range nodes {
		fmt.Printf("%d: %s\n", i+1, node)
	}
}
```
