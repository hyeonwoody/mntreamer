package configuration

import (
	"mntreamer/media/cmd/api/domain/business"
	"mntreamer/media/cmd/api/domain/service"
	"mntreamer/media/cmd/api/infrastructure/repository"
	"mntreamer/media/cmd/api/presentation/controller"
	"mntreamer/media/cmd/api/presentation/handler"
	"mntreamer/media/cmd/model"
	"mntreamer/shared/database"
	"net/http"
	"os"
	"os/exec"
)

type MonolithicContainer struct {
	Variable     *Variable
	Service      service.IService
	Controller   controller.IController
	Handler      handler.IHandler
	MysqlWrapper *database.MysqlWrapper
}

func (ctnr *MonolithicContainer) NewHandler(controller controller.IController) handler.IHandler {
	return handler.NewHandler(ctnr.Variable.BasePath, ctnr.Controller)
}

func (ctnr *MonolithicContainer) InitVariable() error {
	ctnr.Variable = NewVariable()
	return nil
}

func (ctnr *MonolithicContainer) SetRouter(router any) {

}

func (ctnr *MonolithicContainer) DefineRoute() error {
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
	businessMap := map[uint16]business.IBusiness{
		1: business.NewChzzkBusiness(),
	}
	ctnr.Service = service.NewShellScriptService(business.NewBusinessStrategy(businessMap), repository.NewRepository(ctnr.MysqlWrapper))
	ctnr.Handler = handler.NewHandler(ctnr.Variable.BasePath, ctnr.Controller)
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
