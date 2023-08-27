package services

import (
	"context"
	"sisyphos/lib/tasker"
	"sisyphos/models"
)

type SchedulerService struct {
	t *tasker.Tasker
}

func NewSchedulerService(as *ActionService) (*SchedulerService, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	scheduler := tasker.New(ctx, as)
	return &SchedulerService{t: scheduler}, cancel
}

func (s *SchedulerService) Add(t models.Action) {
	// TODO check expression
	s.t.Add(t)
}

func (s *SchedulerService) Remove(id string) {
	s.t.Remove(id)
}

func (s *SchedulerService) Start() {
	s.t.Start()
}
