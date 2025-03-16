package service

import (
	"mntreamer/monitor/cmd/api/infrastructure/repository"
	"mntreamer/monitor/cmd/model"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	repo       repository.IRepository
	backoff    time.Duration
	minBackoff time.Duration
	maxBackoff time.Duration
}

func NewService(repo repository.IRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddToMonitoring(platformId uint16, streamerId uint32) (*model.StreamerMonitor, error) {
	return s.repo.Create(model.NewStreamerMonitor(platformId, streamerId))
}

func (s *Service) Checkout() (*model.StreamerMonitor, error) {
	target, tx, err := s.repo.FindByCheckAtLock(time.Now())
	if err != nil {
		tx.Commit()
		return nil, err
	}
	s.UpdateCheckAt(tx, target)
	tx.Commit()
	return target, nil
}

func (s *Service) UpdateCheckAt(tx *gorm.DB, streamerMonitor *model.StreamerMonitor) {

	switch {
	case streamerMonitor.MissCount < 12:
		streamerMonitor.CheckAt = time.Now().Add(1 * time.Minute)
	case streamerMonitor.MissCount < 24:
		streamerMonitor.CheckAt = time.Now().Add(3 * time.Minute)
	case streamerMonitor.MissCount < 36:
		streamerMonitor.CheckAt = time.Now().Add(10 * time.Minute)
	case streamerMonitor.MissCount < 48:
		streamerMonitor.CheckAt = time.Now().Add(15 * time.Minute)
	case streamerMonitor.MissCount < 60:
		streamerMonitor.CheckAt = time.Now().Add(30 * time.Minute)
	case streamerMonitor.MissCount < 62:
		streamerMonitor.CheckAt = time.Now().Add(40 * time.Minute)
	default:
		streamerMonitor.CheckAt = time.Now().Add(50 * time.Minute)
		streamerMonitor.MissCount = 0
	}
	s.repo.UpdateTx(tx, streamerMonitor)
}

func (s *Service) AddMissCount(streamerMonitor *model.StreamerMonitor) {
	streamerMonitor.MissCount++
	s.repo.Save(streamerMonitor)
}

func (s *Service) ResetMissCount(streamerMonitor *model.StreamerMonitor) {
	streamerMonitor.MissCount = 0
	s.repo.Save(streamerMonitor)
}

func (s *Service) clampBackoff(duration time.Duration, minBackoff time.Duration, maxBackoff time.Duration) time.Duration {
	if duration < minBackoff {
		return minBackoff
	}
	if duration > maxBackoff {
		return maxBackoff
	}
	return duration
}
