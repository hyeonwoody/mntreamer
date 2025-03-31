package configuration

import (
	parser "mntreamer/media/cmd/api/domain/business/parser"
	platform "mntreamer/media/cmd/api/domain/business/platform"
	"mntreamer/media/cmd/api/domain/service"
	"mntreamer/media/cmd/api/infrastructure/repository"
	"mntreamer/media/cmd/api/presentation/controller"
	"mntreamer/media/cmd/api/presentation/handler"
	"mntreamer/media/cmd/model"
	"mntreamer/shared/database"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type MonolithicContainer struct {
	Router       *gin.Engine
	Variable     *Variable
	Service      service.IService
	Controller   controller.IController
	Handler      handler.IHandler
	MysqlWrapper *database.MysqlWrapper
}

func (ctnr *MonolithicContainer) NewHandler(controller controller.IController) handler.IHandler {
	return handler.NewHandler(ctnr.Variable.RootPath, ctnr.Controller)
}

func (ctnr *MonolithicContainer) InitVariable() error {
	ctnr.Variable = NewVariable()
	return nil
}

func (ctnr *MonolithicContainer) SetRouter(router any) {

}

func (ctnr *MonolithicContainer) DefineRoute(router any) error {

	ginRouter, ok := router.(*gin.Engine)
	if !ok {
		ginRouter = gin.Default()
	}
	ctnr.Router = ginRouter
	mediaGroup := ctnr.Router.Group("/api/v1/media")
	{
		mediaGroup.POST("", ctnr.Handler.GetFiles)
		mediaGroup.GET("/stream/*filePath", ctnr.Handler.Stream)
		mediaGroup.GET("/target-duration/*filePath", ctnr.Handler.GetTargetDuration)
		mediaGroup.PATCH("/excise", ctnr.Handler.Excise)
		mediaGroup.PATCH("/confirm/*filePath", ctnr.Handler.Confirm)
		mediaGroup.DELETE("/*filePath", ctnr.Handler.Delete)
	}

	return nil
}
func (ctnr *MonolithicContainer) GetHttpHandler() http.Handler {
	return nil
}

func (ctnr *MonolithicContainer) DefineDatabase(mysqlWrapper any) error {
	ctnrMysqlWrapper := mysqlWrapper.(*database.MysqlWrapper)

	err := ctnrMysqlWrapper.Driver.AutoMigrate(&model.MediaRecord{})
	if err != nil {
		return err
	}
	ctnr.MysqlWrapper = ctnrMysqlWrapper
	return nil
}
func (ctnr *MonolithicContainer) DefineGrpc() error {

	return nil
}

func (ctnr *MonolithicContainer) InitDependency(mysql any) error {
	businessMap := map[uint16]platform.IBusiness{
		1: platform.NewChzzkBusiness(),
	}
	m3u8Biz := parser.NewM3u8Business()
	ctnr.Service = service.NewShellScriptService(platform.NewBusinessStrategy(businessMap), repository.NewRepository(ctnr.MysqlWrapper), m3u8Biz, ctnr.Variable.RootPath)
	ctnr.Handler = handler.NewHandler(ctnr.Variable.RootPath, ctnr.Controller)
	return nil
}

func (ctnr *MonolithicContainer) SignalProcess() error {
	cmd := exec.Command("bash", "-c", "ps -aux | grep ffmpeg | grep -v grep | awk '{print $2}' | xargs sudo kill -2")
	cmd.Stderr = os.Stderr
	return nil
}

func NewMonolithicContainer(mysqlWrapper *database.MysqlWrapper) *MonolithicContainer {
	ctnr := &MonolithicContainer{}
	ctnr.InitVariable()
	ctnr.SignalProcess()
	ctnr.DefineDatabase(mysqlWrapper)
	ctnr.InitDependency(mysqlWrapper)
	return ctnr
}
