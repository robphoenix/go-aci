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
	node, err := client.FabricMembership.New(
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
	err = client.FabricMembership.Add(ctx, node)
	if err != nil {
		log.Fatal(err)
		return
	}

	// // delete node
	// err = client.FabricMembership.Delete(ctx, node)
	// if err != nil {
	//         log.Fatal(err)
	//         return
	// }

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
