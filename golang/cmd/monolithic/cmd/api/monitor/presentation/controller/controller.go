package controller

import (
	"fmt"
	media "mntreamer/media/cmd/api/domain/service"
	monitor "mntreamer/monitor/cmd/api/domain/service"
	platform "mntreamer/platform/cmd/api/domain/service"
	mntreamerModel "mntreamer/shared/model"
	streamer "mntreamer/streamer/cmd/api/domain/service"
	"time"
)

type ControllerMono struct {
	platformSvc platform.IService
	streamerSvc streamer.IService
	monitorSvc  monitor.IService
	mediaSvc    media.IService
}

func NewControllerMono(monitorSvc monitor.IService, platformSvc platform.IService, streamerSvc streamer.IService, mediaSvc media.IService) *ControllerMono {
	ctrl := &ControllerMono{monitorSvc: monitorSvc, platformSvc: platformSvc, streamerSvc: streamerSvc, mediaSvc: mediaSvc}
	go ctrl.beginMonitor()
	return ctrl
}

func (c *ControllerMono) Add(platformName, nickname string) error {

	streamer, err := c.platformSvc.BuildStreamer(platformName, nickname)
	if err != nil {
		return err
	}
	streamer, err = c.streamerSvc.Save(streamer)
	if err != nil {
		return err
	}
	_, err = c.monitorSvc.AddToMonitoring(streamer.PlatformId, streamer.Id)
	if err != nil {
		return err
	}
	return nil
}

func (c *ControllerMono) beginMonitor() {

	for {
		monitor, err := c.monitorSvc.Checkout()
		if err != nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		streamer, err := c.streamerSvc.FindByPlatformIdAndStreamerId(monitor.PlatformId, monitor.StreamerId)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if !c.streamerSvc.CheckMonitoringEligibility(streamer) {
			c.monitorSvc.AddMissCount(monitor)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		c.streamerSvc.UpdateStatus(streamer, mntreamerModel.PROCESS)

		media, err := c.platformSvc.GetLiveDetail(streamer)
		if err != nil {
			c.streamerSvc.UpdateStatus(streamer, mntreamerModel.IDLE)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		c.streamerSvc.UpdateStatus(streamer, mntreamerModel.RECORDING)
		go func() {
			c.mediaSvc.Download(media, streamer)
			c.streamerSvc.UpdateStatus(streamer, mntreamerModel.IDLE)
			c.monitorSvc.ResetMissCount(monitor)
		}()
		fmt.Sprintf(media.Title)
		time.Sleep(5 * time.Second)
	}
}
