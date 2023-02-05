package gormrepo

import (
	"fmt"

	"sisyphos/lib/utils"
	"sisyphos/models"

	"gorm.io/gorm"
)

type Group struct {
	DBModel
	Name        *string `gorm:"index:idx_name,unique"`
	Description *string
	Allows      []string `gorm:"-"`
	AllowsRef   []Action `gorm:"many2many:groups_actions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Users       []string `gorm:"-"`
	UsersRef    []User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;many2many:groups_users;"`
	// Deny []string
}

func (g *Group) BeforeSave(tx *gorm.DB) (err error) {
	if len(g.Users) > 0 {
		userRepo := NewUserRepo(tx)
		users := []User{}
		for _, u := range g.Users {
			userID, err := userRepo.GetID(u)
			if err != nil {
				return err
			}
			if userID == "" {
				return fmt.Errorf("user '%s' not found", u)
			}
			users = append(users, User{
				DBModel: DBModel{ID: userID},
			})
		}
		g.UsersRef = users
	}

	if len(g.Allows) > 0 {
		actionRepo := NewActionRepo(tx)
		actions := []Action{}
		for _, u := range g.Allows {
			actionID, err := actionRepo.GetID(u)
			if err != nil {
				return err
			}
			if actionID == "" {
				return fmt.Errorf("action '%s' not found", u)
			}
			actions = append(actions, Action{
				DBModel: DBModel{ID: actionID},
			})
		}
		g.AllowsRef = actions
	}
	return
}

func (g *Group) AfterFind(tx *gorm.DB) (err error) {
	for _, u := range g.UsersRef {
		g.Users = append(g.Users, u.Name)
	}
	for _, a := range g.AllowsRef {
		g.Allows = append(g.Allows, a.Name)
	}
	return
}

type GroupRepo struct {
	db *gorm.DB
}

func NewGroupRepo(db *gorm.DB) *GroupRepo {
	return &GroupRepo{db}
}

func (r *GroupRepo) getDB() *gorm.DB {
	d := r.db.Session(&gorm.Session{NewDB: true})
	return d
}

func (r *GroupRepo) Create(groups []models.Group) ([]models.Group, error) {
	resp := []models.Group{}
	for _, a := range groups {
		group := MarshalGroup(a)

		err := r.getDB().Omit("UsersRef.*").Omit("AllowsRef.*").Create(&group).Error
		if err != nil {
			return nil, err
		}
		newGroup, err := r.ReadByName(*group.Name)
		if err != nil {
			return nil, err
		}
		resp = append(resp, *newGroup)
	}
	return resp, nil
}

func (r *GroupRepo) ReadByName(name interface{}) (*models.Group, error) {
	var a Group
	err := r.db.Model(&Group{}).Where(&Group{Name: utils.ToPointer(name.(string))}).Preload("UsersRef").Preload("AllowsRef").First(&a).Error
	if err != nil {
		return nil, err
	}
	m := UnmarshalGroup(a)
	return &m, nil
}

func (r *GroupRepo) Update(name any, d *models.Group) (*models.Group, error) {
	uGroup := MarshalGroup(*d)
	if gid, err := r.GetID(name); err != nil {
		return nil, err
	} else {
		uGroup.ID = gid
	}
	if uGroup.Users != nil {
		if err := r.getDB().Model(&uGroup).Association("UsersRef").Replace(uGroup.UsersRef); err != nil {
			return nil, err
		}
	}
	if uGroup.Allows != nil {
		if err := r.getDB().Model(&uGroup).Association("AllowsRef").Replace(uGroup.AllowsRef); err != nil {
			return nil, err
		}
	}

	if err := r.getDB().Omit("AllowsRef.*").Omit("UsersRef.*").Updates(&uGroup).Error; err != nil {
		return nil, err
	}
	return r.ReadByName(name)
}

func (r *GroupRepo) GetID(name any) (string, error) {
	var a Group
	err := r.db.Model(&Group{}).Where(&Group{Name: utils.ToPointer(name.(string))}).First(&a).Error
	if err != nil {
		return "", err
	}
	return a.ID, nil
}

func (r *GroupRepo) ReadAll() ([]models.Group, error) {
	var a []Group
	err := r.db.Model(&Group{}).Preload("UsersRef").Preload("AllowsRef").Find(&a).Error
	if err != nil {
		return nil, err
	}
	m := UnmarshalArrayGroup(a)
	return m, nil
}

func MarshalGroup(a models.Group) Group {
	return Group{
		Name:        a.Name,
		Allows:      a.Allows,
		Users:       a.Users,
		Description: a.Description,
	}
}

func UnmarshalGroup(a Group) models.Group {
	if a.Allows == nil {
		a.Allows = []string{}
	}
	if a.Users == nil {
		a.Users = []string{}
	}
	return models.Group{
		Name:        a.Name,
		Allows:      a.Allows,
		Users:       a.Users,
		Description: a.Description,
	}
}

func MarshalArrayGroup(m []models.Group) []Group {
	groups := []Group{}
	for _, a := range m {
		groups = append(groups, MarshalGroup(a))
	}
	return groups
}

func UnmarshalArrayGroup(a []Group) []models.Group {
	groups := []models.Group{}
	for _, m := range a {
		groups = append(groups, UnmarshalGroup(m))
	}
	return groups
}

type GroupDBMigrator struct {
	db *gorm.DB
}

func (s GroupDBMigrator) Migrate() error {
	fmt.Println("Migrating Group")
	err := s.db.AutoMigrate(&Group{})
	if err != nil {
		return err
	}
	var adminGroup Group
	if err := s.db.Model(&Group{Name: utils.ToPointer("Admin")}).First(&adminGroup).Error; err != nil {
		repo := NewGroupRepo(s.db)
		if _, err := repo.Create([]models.Group{
			{
				Name:  utils.ToPointer("admin"),
				Users: []string{"admin"},
			},
			{
				Name:  utils.ToPointer("default"),
				Users: []string{"admin"},
			},
		}); err != nil {
			return err
		}
	}
	return nil
}

func NewGroupDBMigrator(db *gorm.DB) *GroupDBMigrator {
	return &GroupDBMigrator{db: db}
}
