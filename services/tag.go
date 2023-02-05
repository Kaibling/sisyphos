package services

import (
	"sisyphos/models"
)

type tagRepo interface {
	Create(tags []models.Tag) ([]models.Tag, error)
	ReadByName(name interface{}) (*models.Tag, error)
	ReadAll() ([]models.Tag, error)
}

type TagService struct {
	repo tagRepo
}

func NewTagService(repo tagRepo) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) Create(models []models.Tag) ([]models.Tag, error) {
	return s.repo.Create(models)
}

func (s *TagService) ReadByName(name interface{}) (*models.Tag, error) {
	return s.repo.ReadByName(name)
}

func (s *TagService) ReadAll() ([]models.Tag, error) {
	return s.repo.ReadAll()
}

func (s *TagService) Update(name interface{}, newElements JSON) ([]models.Tag, error) {
	return []models.Tag{}, nil
}
