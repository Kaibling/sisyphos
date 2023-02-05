package services

import (
	"sisyphos/models"
)

type hostRepo interface {
	Create(hosts []models.Host) ([]models.Host, error)
	ReadByName(name interface{}) (*models.Host, error)
	ReadAll() ([]models.Host, error)
}

type HostService struct {
	repo hostRepo
}

func NewHostService(repo hostRepo) *HostService {
	return &HostService{repo: repo}
}

func (s *HostService) Create(models []models.Host) ([]models.Host, error) {
	return s.repo.Create(models)
}

func (s *HostService) ReadByName(name interface{}) (*models.Host, error) {
	return s.repo.ReadByName(name)
}

func (s *HostService) ReadAll() ([]models.Host, error) {
	return s.repo.ReadAll()
}

func (s *HostService) Update(name interface{}, newElements JSON) ([]models.Host, error) {
	return []models.Host{}, nil
}
