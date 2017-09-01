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
	node401, err := client.FabricMembership.NewNode(
		"leaf-401",    // name
		"401",         // id
		"1",           // pod id
		"FOC0849N1BD", // serial number
		"leaf",        // role
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	node402, err := client.FabricMembership.NewNode(
		"leaf-402",    // name
		"402",         // id
		"1",           // pod id
		"FOC0456N2BC", // serial number
		"leaf",        // role
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	// mark node for creation
	node401.SetCreated()

	// mark node for deletion
	node402.SetDeleted()

	nodes := []*aci.Node{
		node401,
		node402,
	}

	// update ACI Fabric Membership
	resp, err := client.FabricMembership.Update(ctx, nodes...)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("resp = %+v\n", resp)

	// see a nodes current status (defaults to create)
	status401 := node401.Status()
	fmt.Println(status401) // Output: "created"
	status402 := node402.Status()
	fmt.Println(status402) // Output: "deleted"

	node402.SetCreated()

	resp, err = client.FabricMembership.Update(ctx, node402)
	if err != nil {
		log.Fatal(err)
		return
	}

	// list nodes
	nodes, err = client.FabricMembership.List(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	for i, node := range nodes {
		fmt.Printf("%d: %s\n", i+1, node)
	}
}
