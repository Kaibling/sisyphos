package gormrepo

import (
	"sisyphos/models"

	"gorm.io/gorm"
)

type Host struct {
	DBModel
	Name    string
	SSHKey  string
	Address string
	TagsRef []Tag    `gorm:"many2many:hosts_tags;"`
	Tags    []string `gorm:"-"`
}

func (h *Host) AfterFind(tx *gorm.DB) (err error) {
	tags := []string{}
	for _, t := range h.TagsRef {
		tags = append(tags, t.Name)
	}
	h.Tags = tags
	return
}

type HostRepo struct {
	db *gorm.DB
}

func NewHostRepo(db *gorm.DB) *HostRepo {
	return &HostRepo{db}
}

func (r *HostRepo) getDB() *gorm.DB {
	d := r.db.Session(&gorm.Session{NewDB: true})
	return d
}

func (r *HostRepo) Create(actions []models.Host) ([]models.Host, error) {
	resp := []models.Host{}
	for _, a := range actions {
		action := MarshalHost(a)

		err := r.getDB().Create(&action).Error
		if err != nil {
			return nil, err
		}
		// update assosiation table

		newHost, err := r.ReadByName(action.Name)
		if err != nil {
			return nil, err
		}
		resp = append(resp, *newHost)
	}
	return resp, nil
}

func (r *HostRepo) ReadByName(name interface{}) (*models.Host, error) {
	var a Host
	err := r.db.Model(&Host{Name: name.(string)}).Preload("TagsRef").First(&a).Error
	if err != nil {
		return nil, err
	}
	m := UnmarshalHost(a)
	return &m, nil
}

func (r *HostRepo) GetID(name string) (string, error) {
	var a Host
	err := r.db.Model(&Host{}).Where(&Host{Name: name}).Find(&a).Error
	if err != nil {
		return "", err
	}
	return a.ID, nil
}

func (r *HostRepo) ReadAll() ([]models.Host, error) {
	var a []Host
	err := r.db.Model(&Host{}).Preload("TagsRef").Find(&a).Error
	if err != nil {
		return nil, err
	}
	m := UnmarshalArrayHost(a)
	return m, nil
}

func MarshalHost(a models.Host) Host {
	return Host{
		Name:    a.Name,
		SSHKey:  a.SSHKey,
		Address: a.Address,
	}
}

func UnmarshalHost(a Host) models.Host {
	return models.Host{
		Name:    a.Name,
		SSHKey:  a.SSHKey,
		Address: a.Address,
	}
}

func MarshalArrayHost(m []models.Host) []Host {
	actions := []Host{}
	for _, a := range m {
		actions = append(actions, MarshalHost(a))
	}
	return actions
}

func UnmarshalArrayHost(a []Host) []models.Host {
	actions := []models.Host{}
	for _, m := range a {
		actions = append(actions, UnmarshalHost(m))
	}
	return actions
}

type HostDBMigrator struct {
	db *gorm.DB
}

func (s HostDBMigrator) Migrate() error {
	err := s.db.AutoMigrate(&Host{})
	if err != nil {
		return err
	}
	return nil
}

func NewHostDBMigrator(db *gorm.DB) *HostDBMigrator {
	return &HostDBMigrator{db: db}
}
