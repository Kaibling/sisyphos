package services

import (
	"errors"

	"sisyphos/models"
)

type userRepo interface {
	Create(users []models.User) ([]models.User, error)
	ReadByName(name interface{}) (*models.User, error)
	ReadIDs([]any) ([]*models.User, error)
	Authenticate(models.Authentication) (*models.User, error)
	ValidateToken(token string) (*models.User, error)
}

type UserService struct {
	repo        userRepo
	permService *PermissionService
}

func NewUserService(repo userRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) AddPermissionService(p *PermissionService) {
	s.permService = p
}

func (s *UserService) Create(models []models.User) ([]models.User, error) {
	for i := 0; i < len(models); i++ {
		// add to default group, if no group is provided
		if len(models[i].Groups) == 0 {
			models[i].Groups = []string{defaultGroupName}
		}
	}
	return s.repo.Create(models)
}

func (s *UserService) ReadByName(name interface{}) (*models.User, error) {
	return s.repo.ReadByName(name)
}

func (s *UserService) ReadIDs(ids []any) ([]*models.User, error) {
	return s.repo.ReadIDs(ids)
}

func (s *UserService) Update(name interface{}, newElements JSON) ([]*models.User, error) {
	return []*models.User{}, nil
}

func (s *UserService) Authenticate(auth models.Authentication) (*models.User, error) {
	return s.repo.Authenticate(auth)
}

func (s *UserService) ValidateToken(token string) (*models.User, error) {
	return s.repo.ValidateToken(token)
}

func (s *UserService) ReadAllPermission(username string) ([]*models.User, error) {
	if s.permService == nil {
		return nil, errors.New("no permission service instantiated")
	}
	ids, err := s.permService.GetUserIDs(username)
	if err != nil {
		return nil, err
	}
	return s.repo.ReadIDs(ids)
}
