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

	// create node
	node, err := client.FabricMembership.NewNode(
		"leaf-101",    // name
		"101",         // id
		"1",           // pod id
		"FOC0849N1BD", // serial number
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	// add node
	resp, err := client.FabricMembership.AddNode(ctx, node)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("resp = %+v\n", resp)

	// // delete node
	// resp, err = client.FabricMembership.DeleteNode(ctx, node)
	// if err != nil {
	//         log.Fatal(err)
	//         return
	// }
	//
	// fmt.Printf("resp = %+v\n", resp)

	// list nodes
	nodes, err := client.FabricMembership.ListNodes(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	for i, node := range nodes {
		fmt.Printf("%d: %s\n", i+1, node)
	}
}
