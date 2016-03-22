package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/syradium/libvirt_rest_go/api"
	"github.com/syradium/libvirt_rest_go/manager"
	"github.com/syradium/libvirt_rest_go/worker"
)

var (
	pool     = flag.String("pool", "default", "Pool name to use for VPS storage")
	poolPath = flag.String("pool-path", "/root/images", "Path to VPS images")
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please specify if it is a worker or api instance")
		return
	}

	flag.Parse()

	conn, err := manager.NewLibvirtFacade("qemu:///system", *pool, *poolPath)

	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	server, err := worker.NewServer(conn)

	if err != nil {
		log.Fatal(err)
		return
	}

	switch os.Args[1] {
	case "worker":
		{
			worker := server.NewWorker("newfoo")
			worker.Launch()
		}
	case "api":
		{
			api.StartAPI(conn, server)
		}
	}
}
