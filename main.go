package main

import (
	"fmt"
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

	vps, err := conn.CreateVPS("test-vps", 128, 1)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(vps)
}
