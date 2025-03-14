package controller

import (
	monitor "mntreamer/monitor/cmd/api/domain/service"
	platform "mntreamer/platform/cmd/api/domain/service"
	streamer "mntreamer/streamer/cmd/api/domain/service"
)

type ControllerMono struct {
	platformSvc platform.IService
	streamerSvc streamer.IService
	monitorSvc  monitor.IService
}

func NewControllerMono(monitorSvc monitor.IService, platformSvc platform.IService, streamerSvc streamer.IService) *ControllerMono {
	return &ControllerMono{monitorSvc: monitorSvc, platformSvc: platformSvc, streamerSvc: streamerSvc}
}

func (c *ControllerMono) Add(platformName, nickname string) error {
	platformId, err := c.platformSvc.GetPlatformIdByName(platformName)
	if err != nil {
		return err
	}

	streamerModel, err := c.streamerSvc.SaveStreamer(nickname, platformId)
	if err != nil {
		return err
	}
	_, err = c.monitorSvc.AddToMonitoring(platformId, streamerModel.Id)
	if err != nil {
		return err
	}
	return nil
}
