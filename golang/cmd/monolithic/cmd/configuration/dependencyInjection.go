package configuration

import (
	"fmt"
	media "mntreamer/media/cmd/configuration"
	monitor "mntreamer/monitor/cmd/configuration"
	mediaCtrl "mntreamer/monolithic/cmd/api/media/presentation/controller"
	monitorCtrl "mntreamer/monolithic/cmd/api/monitor/presentation/controller"
	platform "mntreamer/platform/cmd/configuration"
	"mntreamer/shared/database"
	streamer "mntreamer/streamer/cmd/configuration"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type MonolithicContainer struct {
	Variable     *Variable
	Router       *gin.Engine
	MysqlWrapper *database.MysqlWrapper
	MonitorCtnr  *monitor.MonolithicContainer
	PlatformCtnr *platform.MonolithicContainer
	StreamerCtnr *streamer.MonolithicContainer
	MediaCtnr    *media.MonolithicContainer
}

func (ctnr *MonolithicContainer) InitVariable() error {
	ctnr.Variable = NewVariable()
	return nil
}

func (ctnr *MonolithicContainer) SetRouter(router any) {
	ctnr.Router = router.(*gin.Engine)
	ctnr.Router.Static("/media", "/zzz/mntreamer")
	ctnr.Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Allow frontend domain
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
}

func (ctnr *MonolithicContainer) RunRouter() error {
	ctnr.Router.Run(fmt.Sprintf(":%d", ctnr.Variable.Api.Port))
	return nil
}

func (ctnr *MonolithicContainer) DefineRoute() error {
	monitorGroup := ctnr.Router.Group("/api/v1/monitor")
	{
		monitorGroup.POST("", ctnr.MonitorCtnr.Handler.Add)
	}
	mediaGroup := ctnr.Router.Group("/api/v1/media")
	{
		mediaGroup.POST("", ctnr.MediaCtnr.Handler.GetFiles)
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
	ctnr.MediaCtnr = media.NewMonolithicContainer(ctnr.MysqlWrapper)
	ctnr.MonitorCtnr.Controller = monitorCtrl.NewControllerMono(ctnr.MonitorCtnr.Service, ctnr.PlatformCtnr.Service, ctnr.StreamerCtnr.Service, ctnr.MediaCtnr.Service)
	ctnr.MediaCtnr.Controller = mediaCtrl.NewControllerMono(ctnr.MediaCtnr.Service, ctnr.StreamerCtnr.Service)
	ctnr.MonitorCtnr.Handler = ctnr.MonitorCtnr.NewHandler(ctnr.MonitorCtnr.Controller)
	ctnr.MediaCtnr.Handler = ctnr.MediaCtnr.NewHandler(ctnr.MediaCtnr.Controller)
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
