package main

import (
	"log"

	"github.com/syradium/libvirt_rest_go/manager"
)

func main() {
	conn, err := manager.NewLibvirtFacade("qemu:///system", "mypool")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	err = conn.CreateVPSDisk("vps-123", 1024)

	if err != nil {
		log.Fatal(err)
	}
}
