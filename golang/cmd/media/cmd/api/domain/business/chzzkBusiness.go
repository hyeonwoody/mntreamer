package business

type ChzzkBusiness struct {
}

func NewChzzkBusiness() *ChzzkBusiness {
	return &ChzzkBusiness{}
}

func (chzz *ChzzkBusiness) GetDownloadPath() string {
	return "/zzz/mntreamer/chzzk/"
}
