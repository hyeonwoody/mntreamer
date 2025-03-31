package configuration

import (
	//"mntreamer/monitor/cmd/api/presentation/handler"
	//mntreamerConfiguration "mntreamer/shared/configuration"
	"mntreamer/shared/database"
	mntreamerModel "mntreamer/shared/model"
	"mntreamer/streamer/cmd/api/domain/service"
	"mntreamer/streamer/cmd/api/infrastructure/repository"
	"net/http"

	"gorm.io/gorm"
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

func (ctnr *MonolithicContainer) DefineRoute(router any) error {
	return nil
}
func (ctnr *MonolithicContainer) GetHttpHandler() http.Handler {
	return nil
}

func (ctnr *MonolithicContainer) DefineDatabase(mysqlWrapper any) error {
	ctnrMysqlWrapper := mysqlWrapper.(*database.MysqlWrapper)

	err := ctnrMysqlWrapper.Driver.AutoMigrate(&mntreamerModel.Streamer{})
	if err != nil {
		return err
	}
	result := ctnrMysqlWrapper.Driver.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&mntreamerModel.Streamer{}).Update("status", 1)
	if result.Error != nil {
		return result.Error
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
	ctnr.DefineDatabase(mysqlWrapper)
	ctnr.InitDependency(mysqlWrapper)
	return ctnr
}
