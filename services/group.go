package services

import (
	"sisyphos/models"
)

const defaultGroupName = "default"

type groupRepo interface {
	Create(groups []models.Group) ([]models.Group, error)
	ReadByName(name interface{}) (*models.Group, error)
	ReadAll() ([]models.Group, error)
	Update(name any, d *models.Group) (*models.Group, error)
}

type GroupService struct {
	repo groupRepo
}

func NewGroupService(repo groupRepo) *GroupService {
	return &GroupService{repo: repo}
}

func (s *GroupService) Create(models []models.Group) ([]models.Group, error) {
	return s.repo.Create(models)
}

func (s *GroupService) ReadByName(name interface{}) (*models.Group, error) {
	return s.repo.ReadByName(name)
}

func (s *GroupService) ReadAll() ([]models.Group, error) {
	return s.repo.ReadAll()
}

func (s *GroupService) Update(name any, d *models.Group) (*models.Group, error) {
	return s.repo.Update(name, d)
}
