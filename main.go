package main

import (
	"context"
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
	node, err := client.FabricMembership.NewNode(
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
