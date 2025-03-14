package service

import (
	"mntreamer/platform/cmd/api/infrastructure/repository"
)

type Service struct {
	repo repository.IRepository
}

func NewService(repo repository.IRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetPlatformIdByName(name string) (uint16, error) {
	platform, err := s.repo.FindByName(name)
	if err != nil {
		return 0, err
	}
	return platform.Id, nil
}
