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

	// login
	err = client.Login()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// // create node
	// name := "wedname06"
	// serial := "wedser06"
	// nodeID := "3006"
	// podID := "1"
	// node, err := aci.NewNode(name, nodeID, podID, serial)
	// if err != nil {
	//         log.Fatal(err)
	//         os.Exit(1)
	// }
	//
	// // add node
	// err = client.AddNode(node)
	// if err != nil {
	//         log.Fatal(err)
	//         os.Exit(1)
	// }
	//
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
	for _, v := range nodes {
		fmt.Printf("v = %+v\n", v)
	}
}
