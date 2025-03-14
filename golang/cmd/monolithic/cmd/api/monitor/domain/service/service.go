package service

import "mntreamer/monitor/cmd/api/infrastructure/repository"



type Service struct {
	repo repository.IRepository
}

func NewService(repo repository.IRepository) *Service {
	return &Service{repo: repo}
}