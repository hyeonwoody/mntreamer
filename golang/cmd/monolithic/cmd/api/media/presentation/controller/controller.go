package controller

import (
	media "mntreamer/media/cmd/api/domain/service"
	"mntreamer/media/cmd/model"
	platform "mntreamer/platform/cmd/api/domain/service"
	streamer "mntreamer/streamer/cmd/api/domain/service"
	"os"
	"strings"
	"sync"
)

type ControllerMono struct {
	svc         media.IService
	platformSvc platform.IService
	streamerSvc streamer.IService
}

func NewControllerMono(svc media.IService, platformSvc platform.IService, streamerSvc streamer.IService) *ControllerMono {
	return &ControllerMono{svc: svc, platformSvc: platformSvc, streamerSvc: streamerSvc}
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
			platformName := c.platformSvc.GetPlatformNameById(media.PlatformId)
			filePath := c.svc.GetFilePath(&media, platformName, streamer.ChannelName)
			m3u8Infos, err := c.svc.GetM3u8(filePath, media.Sequence)
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

func (c *ControllerMono) Stream(filePath string) (*os.File, string, error) {
	if strings.HasSuffix(filePath, ".m3u8") {
		fullPath, err := c.svc.StreamMediaPlaylist(filePath)
		return nil, fullPath, err
	}
	if strings.HasSuffix(filePath, ".ts") {
		fullPath, err := c.svc.StreamSegment(filePath)
		return nil, fullPath, err
	}
	if strings.HasSuffix(filePath, ".mp4") {
		file, err := c.svc.StreamMp4(filePath)
		return file, file.Name(), err
	}
	return nil, "", nil
}

func (c *ControllerMono) Excise(path string, begin float64, end float64) error {
	return c.svc.Excise(path, begin, end)
}

func (c *ControllerMono) Confirm(filePath string) error {
	platformName, err := c.svc.GetPlatformNameByFilePath(filePath)
	if err != nil {
		return err
	}
	channelName, err := c.svc.GetChannelNameByFilePath(filePath)
	if err != nil {
		return err
	}
	platformId, _ := c.platformSvc.GetPlatformIdByName(platformName)
	if err != nil {
		return err
	}
	streamer, err := c.streamerSvc.FindByPlatformIdAndChannelName(platformId, channelName)
	if err != nil {
		return err
	}
	_, err = c.svc.Confirm(streamer.PlatformId, streamer.Id, filePath)
	return err
}

func (c *ControllerMono) Delete(filePath string) error {
	platformName, err := c.svc.GetPlatformNameByFilePath(filePath)
	channelName, err := c.svc.GetChannelNameByFilePath(filePath)
	platformId, _ := c.platformSvc.GetPlatformIdByName(platformName)
	streamer, err := c.streamerSvc.FindByPlatformIdAndChannelName(platformId, channelName)
	_, err = c.svc.Delete(streamer.PlatformId, streamer.Id, filePath)
	if err != nil {
		return err
	}
	return nil
}
