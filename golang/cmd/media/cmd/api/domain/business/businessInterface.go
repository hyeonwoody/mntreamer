package business

type IBusiness interface {
	GetDownloadPath() string
}

type BusinessStrategy struct {
	businessMap map[uint16]IBusiness
}

func NewBusinessStrategy(businessmap map[uint16]IBusiness) *BusinessStrategy {
	return &BusinessStrategy{businessMap: businessmap}
}

func (bs *BusinessStrategy) GetDownloadPath(platformId uint16) string {
	if business, ok := bs.businessMap[platformId]; ok {
		return business.GetDownloadPath()
	}
	return "/zzz/mntreamer/DEFAULT/"
}
