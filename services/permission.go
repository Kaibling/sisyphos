package services

type PermissionRepo interface {
	GetActionIDs(username any) ([]any, error)
	GetUserIDs(username any) ([]any, error)
}

type PermissionService struct {
	repo PermissionRepo
}

func NewPermissionService(repo PermissionRepo) *PermissionService {
	return &PermissionService{repo}
}

func (s *PermissionService) GetActionIDs(username any) ([]any, error) {
	return s.repo.GetActionIDs(username)
}

func (s *PermissionService) GetUserIDs(username any) ([]any, error) {
	return s.repo.GetUserIDs(username)
}
