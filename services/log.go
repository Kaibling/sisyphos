package services

import (
	"sisyphos/models"
)

type logRepo interface {
	Create(logs []models.Log) ([]models.Log, error)
	ReadByRequestID(name interface{}) (*models.Log, error)
	ReadAll() ([]models.Log, error)
}

type LogService struct {
	repo logRepo
}

func NewLogService(repo logRepo) *LogService {
	return &LogService{repo: repo}
}

func (s *LogService) Log(url, body, method, user, requestID string) error {
	n := []models.Log{
		{
			Url:       url,
			Body:      body,
			Method:    method,
			User:      user,
			RequestID: requestID,
		},
	}
	_, e := s.repo.Create(n)
	return e
}

func (s *LogService) ReadByRequestID(name interface{}) (*models.Log, error) {
	return s.repo.ReadByRequestID(name)
}

func (s *LogService) ReadAll() ([]models.Log, error) {
	return s.repo.ReadAll()
}

func (s *LogService) Update(name interface{}, newElements JSON) ([]models.Log, error) {
	return []models.Log{}, nil
}
