package main

import (
	"fmt"
	"log"
	"os"

	"github.com/robphoenix/go-aci/aci"
)

func main() {

	// set client options
	opts := aci.ClientOptions{
		Host:     "sandboxapicdc.cisco.com",
		Username: "admin",
		Password: "ciscopsdt",
	}

	// create client
	client, err := aci.NewClient(opts)
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

	// create node
	// n := aci.Node{Name: "IAMNODE03"}
	// err = n.SetID("1252")
	// if err != nil {
	//         log.Fatal(err)
	//         os.Exit(1)
	// }
	n := aci.Node{}
	err = n.SetSerial("serial1")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// add node
	err = client.DeleteNode(n)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

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
