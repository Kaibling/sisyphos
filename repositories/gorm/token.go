package gormrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sisyphos/lib/reqctx"
	"sisyphos/lib/utils"
	"sisyphos/models"

	"gorm.io/gorm"
)

const TokenExpiration = 3 // days

type Token struct {
	DBModel
	Token   string `gorm:"index:idx_name,unique"`
	Expires time.Time
	UserID  string `gorm:"size:255"`
}

type TokenRepo struct {
	db       *gorm.DB
	username string
}

func NewTokenRepo(db *gorm.DB, username string) *TokenRepo {
	ctx := context.WithValue(context.TODO(), reqctx.String("username"), username)
	db = db.WithContext(ctx)
	return &TokenRepo{db, username}
}

func (r *TokenRepo) getDB() *gorm.DB {
	d := r.db.Session(&gorm.Session{NewDB: true})
	return d
}

func (r *TokenRepo) Create(userID string) (*models.Token, error) {
	newToken := Token{
		Token:   utils.NewULID().String(),
		Expires: time.Now().Add(TokenExpiration * time.Hour * 24),
		UserID:  userID,
	}
	err := r.getDB().Create(&newToken).Error
	if err != nil {
		return nil, err
	}
	// update assosiation table
	readToken, err := r.ReadByToken(newToken.Token)
	if err != nil {
		return nil, err
	}
	return readToken, nil
}

func (r *TokenRepo) ReadByToken(token interface{}) (*models.Token, error) {
	var a Token
	err := r.getDB().Model(&Token{}).Where(&Token{Token: token.(string)}).First(&a).Error
	if err != nil {
		return nil, err
	}
	m := UnmarshalToken(a)
	return &m, nil
}

func (r *TokenRepo) GetID(token string) (string, error) {
	var t Token
	if err := r.getDB().Model(&Token{}).Where(&Token{Token: token}).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("GetID: id of '%s' not found", token)
		}
		return "", err
	}
	return t.ID, nil
}

func (r *TokenRepo) ReadAll() ([]models.Token, error) {
	var a []Token
	err := r.getDB().Model(&Token{}).Find(&a).Error
	if err != nil {
		return nil, err
	}
	m := UnmarshalArrayToken(a)
	return m, nil
}

func MarshalToken(a models.Token) Token {
	return Token{
		DBModel: DBModel(a.DBInfo),
		Expires: a.Expires,
		Token:   a.Token,
	}
}

func UnmarshalToken(a Token) models.Token {
	return models.Token{
		DBInfo:  models.DBInfo(a.DBModel),
		Token:   a.Token,
		Expires: a.Expires,
	}
}

func MarshalArrayToken(m []models.Token) []Token {
	actions := []Token{}
	for _, a := range m {
		actions = append(actions, MarshalToken(a))
	}
	return actions
}

func UnmarshalArrayToken(a []Token) []models.Token {
	actions := []models.Token{}
	for _, m := range a {
		actions = append(actions, UnmarshalToken(m))
	}
	return actions
}

type TokenDBMigrator struct {
	db *gorm.DB
}

func (s TokenDBMigrator) Migrate() error {
	err := s.db.AutoMigrate(&Token{})
	if err != nil {
		return err
	}
	return nil
}

func NewTokenDBMigrator(db *gorm.DB) *TokenDBMigrator {
	return &TokenDBMigrator{db: db}
}
