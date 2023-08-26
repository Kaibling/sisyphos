package services

import (
	"fmt"

	"sisyphos/lib/ssh"
	"sisyphos/lib/utils"
	"sisyphos/models"
)

type hostRepo interface {
	Create(hosts []*models.Host) ([]models.Host, error)
	ReadByName(name any) (*models.Host, error)
	Update(name any, d *models.Host) (*models.Host, error)
	ReadAll() ([]models.Host, error)
}

type HostService struct {
	repo hostRepo
}

func NewHostService(repo hostRepo) *HostService {
	return &HostService{repo: repo}
}

func (s *HostService) Create(models []*models.Host) ([]models.Host, error) {
	for _, h := range models {
		sshcfg := h.ToSSHConfig()
		if err := sshcfg.Validate(); err != nil {
			return nil, fmt.Errorf("hostname:'%s':%w", utils.PtrRead(h.Name), err)
		}
		sshc := ssh.NewSSHConnector()
		sshService := NewSSHService(sshc)
		if hk, err := sshService.ReadHostKey(sshcfg.Address, sshcfg.Port); err == nil {
			h.KnownKey = &hk
		} else {
			fmt.Println(err.Error())
		}
	}
	return s.repo.Create(models)
}

func (s *HostService) ReadByName(name any) (*models.Host, error) {
	return s.repo.ReadByName(name)
}

func (s *HostService) ReadAll() ([]models.Host, error) {
	return s.repo.ReadAll()
}

func (s *HostService) Update(name any, d *models.Host) (*models.Host, error) {
	return s.repo.Update(name, d)
}

func (s *HostService) TestConnection(name any) error {
	h, err := s.repo.ReadByName(name)
	if err != nil {
		return err
	}
	sshc := ssh.NewSSHConnector()
	sshService := NewSSHService(sshc)

	cfg := h.ToSSHConfig()
	if err := cfg.Validate(); err != nil {
		return err
	}
	cmd := "test true"
	output, err := sshService.RunCommand(cfg, cmd)
	if err != nil {
		return err
	}
	fmt.Printf("'%s'", output)

	return nil
}

func (s *HostService) GetSSHConfig(name any) (*models.SSHConfig, error) {
	h, err := s.ReadByName(name)
	if err != nil {
		return nil, err
	}
	return &models.SSHConfig{
		Address:    utils.PtrRead(h.Address),
		Port:       utils.PtrRead(h.Port),
		Username:   utils.PtrRead(h.Username),
		Password:   utils.PtrRead(h.Password),
		PrivateKey: utils.PtrRead(h.SSHKey),
		KnownKey:   utils.PtrRead(h.KnownKey),
	}, nil
}
