package configuration

import (
	//"mntreamer/monitor/cmd/api/presentation/handler"
	//mntreamerConfiguration "mntreamer/shared/configuration"

	"mntreamer/platform/cmd/api/domain/business"
	"mntreamer/platform/cmd/api/domain/service"
	"mntreamer/platform/cmd/api/infrastructure/externalApi"
	"mntreamer/platform/cmd/api/infrastructure/repository"
	"mntreamer/platform/cmd/model"
	"mntreamer/shared/database"
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
	err := ctnrMysqlWrapper.Driver.AutoMigrate(&model.Platform{})
	if err != nil {
		return err
	}
	platforms := []model.Platform{
		{Id: 1, Name: "치지직"},
		{Id: 1, Name: "chzzk"},
		{Id: 1, Name: "1"},
	}
	result := ctnrMysqlWrapper.Driver.Create(&platforms)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (ctnr *MonolithicContainer) DefineGrpc() error {

	return nil
}

func (ctnr *MonolithicContainer) InitDependency(mysql any) error {
	ctnrMysqlWrapper := mysql.(*database.MysqlWrapper)
	ctnr.Repository = repository.NewRepository(ctnrMysqlWrapper)
	businessMap := map[uint16]business.IBusiness{
		1: business.NewChzzkBusiness(externalApi.NewChzzkClient()),
	}

	ctnr.Service = service.NewService(business.NewBusinessStrategy(businessMap), ctnr.Repository)
	return nil
}

func NewMonolithicContainer(mysqlWrapper *database.MysqlWrapper) *MonolithicContainer {
	ctnr := &MonolithicContainer{}
	ctnr.InitVariable()
	ctnr.DefineDatabase(mysqlWrapper)
	ctnr.InitDependency(mysqlWrapper)
	return ctnr
}
