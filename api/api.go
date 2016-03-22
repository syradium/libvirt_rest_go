package api

import (
	"log"
	"net/http"

	"github.com/RichardKnop/machinery/v1"
	"github.com/syradium/libvirt_rest_go/manager"
)

// StartAPI spins up rest API service
func StartAPI(conn *manager.LibvirtFacade, server *machinery.Server) {
	vps := VPSService{conn: conn, mq: server}
	vps.Register()

	log.Printf("start listening on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
