package configuration

import (
	//"mntreamer/monitor/cmd/api/presentation/handler"
	//mntreamerConfiguration "mntreamer/shared/configuration"
	"mntreamer/shared/database"
	"mntreamer/streamer/cmd/api/domain/service"
	"mntreamer/streamer/cmd/api/infrastructure/repository"
	"mntreamer/streamer/cmd/model"
	"net/http"
)

type MonolithicContainer struct {
	Variable   *Variable
	Service    service.IService
	Repository repository.IRepository
	//MonitorHandler    *handler.IHandler

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

	err := ctnrMysqlWrapper.Driver.AutoMigrate(&model.Streamer{})
	if err != nil {
		return err
	}

	return nil
}

func (ctnr *MonolithicContainer) DefineGrpc() error {

	return nil
}

func (ctnr *MonolithicContainer) InitDependency(mysql any) error {
	ctnr.Repository = repository.NewRepository(mysql.(*database.MysqlWrapper))
	ctnr.Service = service.NewService(ctnr.Repository)
	return nil
}

func NewMonolithicContainer(mysqlWrapper *database.MysqlWrapper) *MonolithicContainer {
	ctnr := &MonolithicContainer{}
	ctnr.InitVariable()
	ctnr.InitDependency(mysqlWrapper)
	return ctnr
}
