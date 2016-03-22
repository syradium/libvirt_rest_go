package worker

import (
	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/syradium/libvirt_rest_go/manager"
)

var (
	libvirtConn *manager.LibvirtFacade
)

// NewServer ...
func NewServer(conn *manager.LibvirtFacade) (*machinery.Server, error) {
	cnf := config.Config{
		Broker:        "redis://127.0.0.1:6379",
		ResultBackend: "redis://127.0.0.1:6379",
		Exchange:      "machinery_exchange",
		ExchangeType:  "direct",
		DefaultQueue:  "machinery_tasks",
		BindingKey:    "machinery_tasks",
	}

	server, err := machinery.NewServer(&cnf)
	if err != nil {
		return nil, err
	}

	// Not the best way to pass libvrt connection to tasks
	libvirtConn = conn

	tasks := map[string]interface{}{
		"createVPS": createVPS,
	}

	server.RegisterTasks(tasks)
	return server, nil
}
