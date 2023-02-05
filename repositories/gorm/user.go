package gormrepo

import (
	"errors"
	"fmt"
	"time"

	"sisyphos/lib/apperrors"
	"sisyphos/lib/config"
	"sisyphos/lib/utils"
	"sisyphos/models"

	"gorm.io/gorm"
)

type User struct {
	DBModel
	Name         string `gorm:"index:idx_name,unique"`
	PasswordHash string
	Token        []string `gorm:"-"`
	TokenRef     []Token
	Groups       []string `gorm:"-"`
	GroupsRef    []Group  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;many2many:groups_users;"`
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	groups := []Group{}
	for _, g := range u.Groups {
		groupRepo := NewGroupRepo(tx)
		gid, err := groupRepo.GetID(g)
		if err != nil {
			return err
		}
		groups = append(groups, Group{DBModel: DBModel{ID: gid}})
	}
	u.GroupsRef = groups
	return
}

func (u *User) AfterFind(tx *gorm.DB) (err error) {
	for _, tr := range u.TokenRef {
		u.Token = append(u.Token, tr.Token)
	}
	for _, gr := range u.GroupsRef {
		u.Groups = append(u.Groups, *gr.Name)
	}
	return
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) getDB() *gorm.DB {
	d := r.db.Session(&gorm.Session{NewDB: true})
	return d
}

func (r *UserRepo) Create(users []models.User) ([]models.User, error) {
	resp := []models.User{}
	for _, a := range users {
		user := MarshalUser(a)
		var err error
		user.PasswordHash, err = utils.HashPassword(user.PasswordHash)
		if err != nil {
			return nil, err
		}
		err = r.getDB().Create(&user).Error
		if err != nil {
			return nil, err
		}
		newUser, err := r.ReadByName(user.Name)
		if err != nil {
			return nil, err
		}
		resp = append(resp, *newUser)
	}
	return resp, nil
}

func (r *UserRepo) ReadByName(name interface{}) (*models.User, error) {
	if name.(string) == "" {
		return nil, fmt.Errorf("empty name not allowed")
	}
	var a User
	err := r.db.Model(&User{}).Where(&User{Name: name.(string)}).Preload("TokenRef").Preload("GroupsRef").First(&a).Error
	if err != nil {
		return nil, err
	}
	return UnmarshalUser(a), nil
}

func (r *UserRepo) GetID(name string) (string, error) {
	var a User
	err := r.db.Model(&User{}).Where(&User{Name: name}).Find(&a).Error
	if err != nil {
		return "", err
	}
	return a.ID, nil
}

func (r *UserRepo) ReadIDs(ids []any) ([]*models.User, error) {
	var a []User
	err := r.db.Model(&User{}).Where("id IN ?", ids).Preload("TokenRef").Preload("GroupsRef").Find(&a).Error
	if err != nil {
		return nil, err
	}

	return UnmarshalArrayUser(a), nil
}

func (r *UserRepo) Authenticate(auth models.Authentication) (*models.User, error) {
	var u User
	err := r.db.Model(&User{}).Where(&User{Name: auth.Username}).Preload("TokenRef").Preload("GroupsRef").First(&u).Error
	if err != nil {
		return nil, err
	}

	if !utils.CheckPasswordHash(auth.Password, u.PasswordHash) {
		return nil, fmt.Errorf("username/password incorrect")
	}
	tr := NewTokenRepo(r.db)
	newToken, err := tr.Create(u.ID)
	if err != nil {
		return nil, err
	}
	u.Token = []string{newToken.Token}
	return UnmarshalUser(u), nil
}

func (r *UserRepo) ValidateToken(token string) (*models.User, error) {
	var t Token
	if err := r.getDB().Model(&Token{}).Where(&Token{Token: token}).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("ValidateToken:token not in database")
			return nil, errors.New(apperrors.TokenInvalid)
		}
		return nil, err
	}
	if t.Expires.Before(time.Now()) {
		fmt.Println("ValidateToken:token expired")
		return nil, errors.New(apperrors.TokenInvalid)
	}
	var u User
	if err := r.getDB().Model(&User{}).Where(&User{DBModel: DBModel{ID: t.UserID}}).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("ValidateToken:user not found")
			return nil, errors.New(apperrors.TokenInvalid)
		}
		return nil, err
	}

	return UnmarshalUser(u), nil
}

func MarshalUser(a models.User) User {
	return User{
		Name:         a.Name,
		PasswordHash: a.Password,
		Groups:       a.Groups,
	}
}

func UnmarshalUser(a User) *models.User {
	if a.Token == nil {
		a.Token = []string{}
	}
	return &models.User{
		Name:   a.Name,
		Token:  a.Token,
		Groups: a.Groups,
	}
}

func MarshalArrayUser(m []models.User) []User {
	users := []User{}
	for _, a := range m {
		users = append(users, MarshalUser(a))
	}
	return users
}

func UnmarshalArrayUser(a []User) []*models.User {
	users := []*models.User{}
	for _, m := range a {
		users = append(users, UnmarshalUser(m))
	}
	return users
}

type UserDBMigrator struct {
	db *gorm.DB
}

func (s UserDBMigrator) Migrate() error {
	err := s.db.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	var adminUser User
	if err := s.db.Model(&User{Name: "admin"}).First(&adminUser).Error; err != nil {
		var pwd string
		if config.Config.AdminPassword == "" {
			pwd = utils.NewULID().String()
		} else {
			pwd = config.Config.AdminPassword
		}

		userRepo := NewUserRepo(s.db)
		if _, err := userRepo.Create([]models.User{{Name: "admin", Password: pwd}}); err != nil {
			return fmt.Errorf("creating of user Admin failed: %w", err)
		}
		fmt.Printf("admin password: %s\n", pwd)
	}
	return nil
}

func NewUserDBMigrator(db *gorm.DB) *UserDBMigrator {
	return &UserDBMigrator{db: db}
}
