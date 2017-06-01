package main

import (
	"fmt"
	"log"

	"github.com/robphoenix/go-aci/aci"
)

func main() {
	apicURL := "sandboxapicdc.cisco.com"
	apicUser := "admin"
	apicPwd := "ciscopsdt"

	apicClient, err := aci.NewClient(apicURL, apicUser, apicPwd)
	if err != nil {
		log.Fatal(err)
	}
	err = apicClient.Login()
	if err != nil {
		log.Fatal(err)
	}
	// vrf := aci.VRF{Name: "Jo_Cox", Tenant: "CORBYN"}
	// err = aci.CreateVRF(apicClient, vrf)
	// if err != nil {
	//         log.Fatal(err)
	// }
	vrfs, err := aci.ListVRFs(apicClient, "CORBYN")
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range vrfs {
		fmt.Printf("v = %+v\n", v)
	}
}
