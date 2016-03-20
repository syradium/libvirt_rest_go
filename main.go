package main

import (
	"log"

	"github.com/syradium/libvirt_rest_go/api"
	"github.com/syradium/libvirt_rest_go/manager"
)

func main() {
	conn, err := manager.NewLibvirtFacade("qemu:///system", "mypool")
	if err != nil {
		log.Fatal(err)
		return
	}

	api.StartAPI(*conn)
}
