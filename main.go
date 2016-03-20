package main

import (
	"log"

	"github.com/syradium/libvirt_rest_go/manager"
)

func main() {
	conn, err := manager.NewLibvirtFacade("qemu:///system")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
}
