package api

import (
	"net/http"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/signatures"
	"github.com/emicklei/go-restful"
	"github.com/syradium/libvirt_rest_go/manager"
)

// VPSService ...
type VPSService struct {
	mq   *machinery.Server
	conn *manager.LibvirtFacade
}

// Register ...
func (v VPSService) Register() {
	ws := new(restful.WebService)

	ws.
		Path("/vps").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("list").To(v.listAll).
		// docs
		Doc("retrieve all available VPS").
		Operation("listAll").
		Returns(200, "OK", []manager.VPS{}))

	ws.Route(ws.POST("").To(v.createVPS).
		// docs
		Doc("create VPS").
		Operation("createVPS").
		Reads(manager.VPS{}))

	ws.Route(ws.GET("stop/{vps-name}").To(v.stopVPS).
		// docs
		Doc("stop VPS").
		Operation("stopVPS").
		Param(ws.PathParameter("vps-name", "identifier of the vps").DataType("string")).
		Writes(manager.VPS{}))

	ws.Route(ws.GET("start/{vps-name}").To(v.startVPS).
		// docs
		Doc("start VPS").
		Operation("startVPS").
		Param(ws.PathParameter("vps-name", "identifier of the vps").DataType("string")).
		Writes(manager.VPS{}))

	ws.Route(ws.GET("task/{task-uuid}").To(v.getTaskInfo).
		// docs
		Doc("get task info").
		Operation("getTaskInfo").
		Param(ws.PathParameter("task-uuid", "UUID of a task").DataType("string")))

	restful.Add(ws)
}

func (v VPSService) listAll(request *restful.Request, response *restful.Response) {
	vpsList := []manager.VPS{manager.VPS{Name: "foo"}}
	response.WriteEntity(vpsList)
}

func (v VPSService) createVPS(request *restful.Request, response *restful.Response) {
	vps := new(manager.VPS)
	err := request.ReadEntity(vps)
	if err != nil {
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	task := signatures.TaskSignature{
		Name: "createVPS",
		Args: []signatures.TaskArg{
			{
				Type:  "string",
				Value: vps.Name,
			},
			{
				Type:  "uint64",
				Value: vps.RAM,
			},
			{
				Type:  "uint64",
				Value: vps.DiskSize,
			},
		},
	}

	asyncResult, err := v.mq.SendTask(&task)

	if err != nil {
		response.WriteHeaderAndJson(http.StatusNotFound, err, "application/json")
		return
	}

	response.WriteHeaderAndEntity(http.StatusCreated, asyncResult.GetState().TaskUUID)
}

func (v VPSService) stopVPS(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("vps-name")
	vps, err := v.conn.GetVPSByName(name)

	if err != nil {
		response.WriteHeaderAndJson(http.StatusNotFound, err, "application/json")
		return
	}
	v.conn.DestroyVPS(vps)
}

func (v VPSService) startVPS(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("vps-name")
	vps, err := v.conn.GetVPSByName(name)

	if err != nil {
		response.WriteHeaderAndJson(http.StatusNotFound, err, "application/json")
		return
	}
	v.conn.StartVPS(vps)
}

func (v VPSService) getTaskInfo(request *restful.Request, response *restful.Response) {
	taskUUID := request.PathParameter("task-uuid")
	backend := v.mq.GetBackend()
	taskState, err := backend.GetState(taskUUID)

	if err != nil {
		response.WriteHeaderAndJson(http.StatusNotFound, err, "application/json")
		return
	}

	response.WriteAsJson(taskState)
}
