package gormrepo

import (
	"errors"
	"fmt"
	"sisyphos/lib/utils"
	"sisyphos/models"

	"gorm.io/gorm"
)

type Host struct {
	DBModel
	Name     *string `gorm:"index:idx_name,unique"`
	User     *string
	Password *string
	SSHKey   *string
	Address  *string `gorm:"index:idx_name,unique"`
	Port     *int
	TagsRef  []Tag    `gorm:"many2many:hosts_hosts;"`
	Tags     []string `gorm:"-"`
}

func (h *Host) BeforeSave(tx *gorm.DB) (err error) {
	hosts := []Tag{}
	for _, t := range h.Tags {
		hostRepo := NewTagRepo(tx)
		tid, err := hostRepo.GetID(t)
		if err != nil {
			return err
		}
		hosts = append(hosts, Tag{DBModel: DBModel{ID: tid}})
	}
	h.TagsRef = hosts
	return
}

func (h *Host) AfterFind(tx *gorm.DB) (err error) {
	hosts := []string{}
	for _, t := range h.TagsRef {
		hosts = append(hosts, t.Name)
	}
	h.Tags = hosts
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

func (r *HostRepo) Create(hosts []models.Host) ([]models.Host, error) {
	resp := []models.Host{}
	for _, a := range hosts {
		host := MarshalHost(a)

		err := r.getDB().Omit("TagsRef.*").Create(&host).Error
		if err != nil {
			return nil, err
		}
		newHost, err := r.ReadByName(utils.PtrRead(host.Name))
		if err != nil {
			return nil, err
		}
		resp = append(resp, *newHost)
	}
	return resp, nil
}

func (r *HostRepo) ReadByName(name interface{}) (*models.Host, error) {
	var a Host
	err := r.getDB().Model(&Host{}).Where(&Host{Name: utils.ToPointer(name.(string))}).Preload("TagsRef").First(&a).Error
	if err != nil {
		return nil, err
	}
	m := UnmarshalHost(a)
	return &m, nil
}

func (r *HostRepo) Update(name any, d *models.Host) (*models.Host, error) {
	uHost := MarshalHost(*d)
	if hid, err := r.GetID(name); err != nil {
		return nil, err
	} else {
		uHost.ID = hid
	}
	if uHost.Tags != nil {
		if err := r.getDB().Model(&uHost).Association("TagsRef").Replace(uHost.TagsRef); err != nil {
			return nil, err
		}
	}

	if err := r.getDB().Omit("TagsRef.*").Updates(&uHost).Error; err != nil {
		return nil, err
	}
	return r.ReadByName(name)
}

func (r *HostRepo) GetID(name any) (string, error) {
	var a Host
	if err := r.getDB().Model(&Host{}).Where(&Host{Name: utils.ToPointer(name.(string))}).First(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("GetID: id of '%s' not found", name)
		}
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
		Name:     a.Name,
		User:     a.Username,
		Password: a.Password,
		SSHKey:   a.SSHKey,
		Address:  a.Address,
		Port:     a.Port,
		Tags:     a.Tags,
	}
}

func UnmarshalHost(a Host) models.Host {
	if len(a.Tags) == 0 {
		a.Tags = []string{}
	}
	return models.Host{
		Name:     a.Name,
		Username: a.User,
		Password: a.Password,
		SSHKey:   a.SSHKey,
		Address:  a.Address,
		Port:     a.Port,
		Tags:     a.Tags,
	}
}

func MarshalArrayHost(m []models.Host) []Host {
	hosts := []Host{}
	for _, a := range m {
		hosts = append(hosts, MarshalHost(a))
	}
	return hosts
}

func UnmarshalArrayHost(a []Host) []models.Host {
	hosts := []models.Host{}
	for _, m := range a {
		hosts = append(hosts, UnmarshalHost(m))
	}
	return hosts
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
