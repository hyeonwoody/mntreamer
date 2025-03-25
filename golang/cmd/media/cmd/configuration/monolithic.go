package configuration

import (
	"mntreamer/media/cmd/api/domain/business"
	"mntreamer/media/cmd/api/domain/infrastructure/repository"
	"mntreamer/media/cmd/api/domain/service"
	"mntreamer/media/cmd/model"
	"mntreamer/shared/database"
	"net/http"
	"os"
	"os/exec"
)

type MonolithicContainer struct {
	Variable     *Variable
	Service      service.IService
	MysqlWrapper *database.MysqlWrapper
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
	return nil
}

func (ctnr *MonolithicContainer) SignalProcess() error {
	cmd := exec.Command("bash", "-c", "ps -aux | grep ffmpeg | grep -v grep | awk '{print $2}' | xargs sudo kill -2")
	cmd.Stderr = os.Stderr
	return nil
}

func NewMonolithicContainer(mysqlWrapper *database.MysqlWrapper) *MonolithicContainer {
	ctnr := &MonolithicContainer{}
	ctnr.SignalProcess()
	ctnr.DefineDatabase(mysqlWrapper)
	ctnr.InitDependency(mysqlWrapper)
	return ctnr
}
