# go-aci
Go API wrapper for Cisco ACI

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/robphoenix/go-aci/aci"
)

func main() {

	// set client options
	cfg := aci.Config{
		Host:     "sandboxapicdc.cisco.com",
		Username: "admin",
		Password: "ciscopsdt",
	}

	// create client
	client, err := aci.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// authenticate
	err = client.Login()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// create node
	name := "leaf-101"
	serial := "FOC0849N1BD"
	nodeID := "101"
	podID := "1"

	node, err := aci.NewNode(name, nodeID, podID, serial)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// add node
	err = client.AddNode(node)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// // delete node
	// err = client.DeleteNode(node)
	// if err != nil {
	//         log.Fatal(err)
	//         os.Exit(1)
	// }

	// list nodes
	nodes, err := client.ListNodes()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	for i, node := range nodes {
		fmt.Printf("%d: %s\n", i+1, node)
	}
}
```
