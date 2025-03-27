package controller

import (
	media "mntreamer/media/cmd/api/domain/service"
	"mntreamer/media/cmd/model"
	platform "mntreamer/platform/cmd/api/domain/service"
	streamer "mntreamer/streamer/cmd/api/domain/service"
	"sync"
)

type ControllerMono struct {
	svc         media.IService
	platformSvc platform.IService
	streamerSvc streamer.IService
}

func NewControllerMono(svc media.IService, streamerSvc streamer.IService) *ControllerMono {
	return &ControllerMono{svc: svc, streamerSvc: streamerSvc}
}

func (c *ControllerMono) GetFiles(filePath string) ([]model.FileInfo, error) {
	return c.svc.GetFiles(filePath)
}

func (c *ControllerMono) GetTargetDuration(path string) (float64, error) {
	playlist, _ := c.svc.Decode(path)
	mpl, _ := playlist.(*model.MediaPlaylist)
	return mpl.TargetDuration, nil
}

func (c *ControllerMono) GetFilesToRefine() ([]model.FileInfo, error) {
	mediaList, err := c.svc.GetMediaToRefine()
	if err != nil {
		return nil, err
	}
	var m3u8List []model.FileInfo
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, media := range mediaList {
		wg.Add(1)
		go func(m model.MediaRecord) {
			defer wg.Done()
			streamer, err := c.streamerSvc.FindByPlatformIdAndStreamerId(media.PlatformId, media.StreamerId)
			if err != nil {
				return
			}
			filePath := c.svc.GetFilePath(&media, streamer.ChannelName)
			m3u8Infos, err := c.svc.GetM3u8(filePath)
			if err != nil {
				return
			}
			mu.Lock()
			m3u8List = append(m3u8List, m3u8Infos...)
			mu.Unlock()
		}(media)
	}
	wg.Wait()
	return m3u8List, nil
}

func (c *ControllerMono) Stream(filePath string) (string, error) {
	return c.svc.Stream(filePath)
}

func (c *ControllerMono) Excise(path string, begin float64, end float64) error {
	return c.svc.Excise(path, begin, end)
}
