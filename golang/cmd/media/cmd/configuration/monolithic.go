package configuration

import (
	"mntreamer/media/cmd/api/domain/business"
	"mntreamer/media/cmd/api/domain/service"
	"mntreamer/shared/database"
	"net/http"
)

type MonolithicContainer struct {
	Service service.IService
}

func (ctnr *MonolithicContainer) InitVariable() error {
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
	return nil
}
func (ctnr *MonolithicContainer) DefineGrpc() error {

	return nil
}

func (ctnr *MonolithicContainer) InitDependency(mysql any) error {
	businessMap := map[uint16]business.IBusiness{
		1: business.NewChzzkBusiness(),
	}

	ctnr.Service = service.NewShellScriptService(business.NewBusinessStrategy(businessMap))
	return nil
}

func NewMonolithicContainer(mysqlWrapper *database.MysqlWrapper) *MonolithicContainer {
	ctnr := &MonolithicContainer{}
	ctnr.InitVariable()
	ctnr.DefineDatabase(mysqlWrapper)
	ctnr.InitDependency(mysqlWrapper)
	return ctnr
}
