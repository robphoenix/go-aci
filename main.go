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

	// // create node
	// name := "wedname07"
	// serial := "wedser07"
	// nodeID := "3007"
	// podID := "1"
	// node, err := client.FabricMembership.NewNode(name, nodeID, podID, serial)
	// if err != nil {
	//         log.Fatal(err)
	//         os.Exit(1)
	// }
	//
	// // add node
	// err = client.FabricMembership.AddNode(node)
	// if err != nil {
	//         log.Fatal(err)
	//         os.Exit(1)
	// }

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
