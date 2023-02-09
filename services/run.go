package services

import (
	"sisyphos/models"
)

type runRepo interface {
	Create(runs []models.Run) ([]models.Run, error)
	ReadByRunID(runID interface{}) (*models.Run, error)
	ReadByReqID(reqID interface{}) ([]models.Run, error)
	ReadAll() ([]models.Run, error)
	GetRequestID() string
	GetUsername() string
}

type RunService struct {
	repo runRepo
}

func NewRunService(repo runRepo) *RunService {
	return &RunService{repo: repo}
}

func (s *RunService) Create(model *models.Run) ([]models.Run, error) {
	return s.repo.Create([]models.Run{*model})
}

func (s *RunService) ReadByRunID(runID interface{}) (*models.Run, error) {
	return s.repo.ReadByRunID(runID)
}

func (s *RunService) ReadByReqID() ([]models.Run, error) {
	return s.repo.ReadByReqID(s.repo.GetRequestID())
}

func (s *RunService) ReadAll() ([]models.Run, error) {
	return s.repo.ReadAll()
}
