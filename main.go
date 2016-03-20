package main

import (
	"flag"
	"log"

	"github.com/syradium/libvirt_rest_go/api"
	"github.com/syradium/libvirt_rest_go/manager"
)

var (
	pool     = flag.String("pool", "default", "Pool name to use for VPS storage")
	poolPath = flag.String("pool-path", "/root/images", "Path to VPS images")
)

func main() {
	flag.Parse()
	conn, err := manager.NewLibvirtFacade("qemu:///system", *pool, *poolPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	api.StartAPI(*conn)
}
