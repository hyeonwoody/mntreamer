package configuration

import (
	"mntreamer/monitor/cmd/api/domain/service"
	"mntreamer/monitor/cmd/api/infrastructure/repository"
	"mntreamer/monitor/cmd/api/presentation/controller"
	"mntreamer/monitor/cmd/api/presentation/handler"
	"mntreamer/monitor/cmd/model"
	"mntreamer/shared/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MonolithicContainer struct {
	Variable   *Variable
	Router     *gin.Engine
	Handler    handler.IHandler
	Controller controller.IController
	Service    service.IService
	Repository repository.IRepository
}

func (ctnr *MonolithicContainer) NewHandler(controller controller.IController) handler.IHandler {
	return handler.NewHandler(ctnr.Controller)
}

func (ctnr *MonolithicContainer) InitVariable() error {
	ctnr.Variable = NewVariable()
	return nil
}

func (ctnr *MonolithicContainer) SetRouter(router any) {
	ctnr.Router = router.(*gin.Engine)
}

func (ctnr *MonolithicContainer) DefineRoute() error {
	return nil
}
func (ctnr *MonolithicContainer) GetHttpHandler() http.Handler {
	return nil
}

func (ctnr *MonolithicContainer) DefineDatabase(mysqlWrapper any) error {
	ctnrMysqlWrapper := mysqlWrapper.(*database.MysqlWrapper)

	err := ctnrMysqlWrapper.Driver.AutoMigrate(&model.StreamerMonitor{})
	if err != nil {
		return err
	}

	return nil
}

func (ctnr *MonolithicContainer) DefineMysqlWrapper(mysqlWrapper *database.MysqlWrapper) error {
	return nil
}

func (ctnr *MonolithicContainer) DefineGrpc() error {

	return nil
}

func (ctnr *MonolithicContainer) InitDependency(mysqlWrapper any) error {
	ctnr.Repository = repository.NewRepository(mysqlWrapper.(*database.MysqlWrapper))
	ctnr.Service = service.NewService(ctnr.Repository)
	ctnr.Handler = handler.NewHandler(ctnr.Controller)
	return nil
}

func NewMonolithicContainer(mysqlWrapper *database.MysqlWrapper) *MonolithicContainer {
	ctnr := &MonolithicContainer{}
	ctnr.InitVariable()
	ctnr.DefineDatabase(mysqlWrapper)
	ctnr.InitDependency(mysqlWrapper)
	return ctnr
}
