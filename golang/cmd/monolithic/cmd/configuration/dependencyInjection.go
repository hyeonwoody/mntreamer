package configuration

import (
	"fmt"
	monitor "mntreamer/monitor/cmd/configuration"
	platform "mntreamer/platform/cmd/configuration"
	streamer "mntreamer/streamer/cmd/configuration"

	monitorCtrl "mntreamer/monolithic/cmd/api/monitor/presentation/controller"
	"mntreamer/shared/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MonolithicContainer struct {
	Variable     *Variable
	Router       *gin.Engine
	MysqlWrapper *database.MysqlWrapper
	MonitorCtnr  *monitor.MonolithicContainer
	PlatformCtnr *platform.MonolithicContainer
	StreamerCtnr *streamer.MonolithicContainer
}

func (ctnr *MonolithicContainer) InitVariable() error {
	ctnr.Variable = NewVariable()
	return nil
}

func (ctnr *MonolithicContainer) SetRouter(router any) {
	ctnr.Router = router.(*gin.Engine)
}

func (ctnr *MonolithicContainer) RunRouter() error {
	ctnr.Router.Run(fmt.Sprintf(":%d", ctnr.Variable.Api.Port))
	return nil
}

func (ctnr *MonolithicContainer) DefineRoute() error {
	monitorGroup := ctnr.Router.Group("/monitor")
	{
		monitorGroup.POST("", ctnr.MonitorCtnr.Handler.Add)
	}
	return nil
}
func (ctnr *MonolithicContainer) GetHttpHandler() http.Handler {
	return nil
}

func (ctnr *MonolithicContainer) DefineDatabase(dependency any) error {
	ctnr.MysqlWrapper = database.ConnectMysqlDatabase(ctnr.Variable.Database)

	return nil
}

func (ctnr *MonolithicContainer) DefineGrpc() error {

	return nil
}

func (ctnr *MonolithicContainer) InitDependency(dependency any) error {
	ctnr.MonitorCtnr = monitor.NewMonolithicContainer(ctnr.MysqlWrapper)
	ctnr.PlatformCtnr = platform.NewMonolithicContainer(ctnr.MysqlWrapper)
	ctnr.StreamerCtnr = streamer.NewMonolithicContainer(ctnr.MysqlWrapper)
	ctnr.MonitorCtnr.Controller = monitorCtrl.NewControllerMono(ctnr.MonitorCtnr.Service, ctnr.PlatformCtnr.Service, ctnr.StreamerCtnr.Service)

	ctnr.MonitorCtnr.Handler = ctnr.MonitorCtnr.NewHandler(ctnr.MonitorCtnr.Controller)
	return nil
}

func NewMonolithicContainer() *MonolithicContainer {
	ctnr := &MonolithicContainer{}
	ctnr.InitVariable()
	ctnr.DefineDatabase(nil)
	ctnr.InitDependency(nil)

	router := gin.Default()
	ctnr.SetRouter(router)

	ctnr.DefineRoute()
	return ctnr
}
