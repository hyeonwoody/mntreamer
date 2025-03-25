package controller

import (
	media "mntreamer/media/cmd/api/domain/service"
	monitor "mntreamer/monitor/cmd/api/domain/service"
	monitorModel "mntreamer/monitor/cmd/model"
	platform "mntreamer/platform/cmd/api/domain/service"
	mntreamerModel "mntreamer/shared/model"
	streamer "mntreamer/streamer/cmd/api/domain/service"
	"sync"
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
		time.Sleep(200 * time.Millisecond)
		go c.monitorProcess()
	}
}

func (c *ControllerMono) monitorProcess() {
	monitor, err := c.monitorSvc.Checkout()
	if err != nil {
		return
	}

	streamer, err := c.streamerSvc.FindByPlatformIdAndStreamerId(monitor.PlatformId, monitor.StreamerId)
	if err != nil {
		return
	}

	if !c.streamerSvc.CheckMonitoringEligibility(streamer) {
		//c.monitorSvc.IncreaseMissCount(monitor)
		return
	}
	c.streamerSvc.UpdateStatus(streamer, mntreamerModel.PROCESS)

	media, err := c.platformSvc.GetLiveDetail(streamer)
	if err != nil {
		c.streamerSvc.UpdateStatus(streamer, mntreamerModel.IDLE)
		c.monitorSvc.IncreaseMissCount(monitor)
		return
	}
	c.streamerSvc.UpdateStatus(streamer, mntreamerModel.RECORDING)
	go func(channelName string, monitor *monitorModel.StreamerMonitor, media *mntreamerModel.Media) {
		c.mediaSvc.Download(media, streamer.ChannelName, monitor.PlatformId)
		go c.postStream(monitor)
	}(streamer.ChannelName, monitor, media)
}

func (c *ControllerMono) postStream(monitor *monitorModel.StreamerMonitor) {
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		c.handleStreamerAfterStream(monitor)
	}()
	go func() {
		defer wg.Done()
		c.handleMonitorAfterStream(monitor)
	}()
	go func() {
		defer wg.Done()
		c.handleMediaAfterStream(monitor)
	}()
	wg.Wait()
}

func (c *ControllerMono) handleMediaAfterStream(monitor *monitorModel.StreamerMonitor) {
	c.mediaSvc.Save(monitor.PlatformId, monitor.StreamerId)
}

func (c *ControllerMono) handleMonitorAfterStream(monitor *monitorModel.StreamerMonitor) {
	c.monitorSvc.ResetMissCount(monitor)
}

func (c *ControllerMono) handleStreamerAfterStream(monitor *monitorModel.StreamerMonitor) {
	streamer, err := c.streamerSvc.FindByPlatformIdAndStreamerId(monitor.PlatformId, monitor.StreamerId)
	if err != nil {
		return
	}
	c.streamerSvc.UpdateStatus(streamer, mntreamerModel.IDLE)
	c.streamerSvc.UpdateLastStreamAt(streamer)
	c.streamerSvc.UpdateLastRecordedAt(streamer)
}
