package gormrepo

import (
	"errors"
	"fmt"
	"sisyphos/models"

	"gorm.io/gorm"
)

type Tag struct {
	DBModel
	Name        string `gorm:"index:idx_name,unique"`
	Description string
}

type TagRepo struct {
	db       *gorm.DB
	username string
}

func NewTagRepo(db *gorm.DB, username string) *TagRepo {
	return &TagRepo{db, username}
}

func (r *TagRepo) getDB() *gorm.DB {
	d := r.db.Session(&gorm.Session{NewDB: true})
	return d
}

func (r *TagRepo) Create(tags []models.Tag) ([]models.Tag, error) {
	resp := []models.Tag{}
	for _, a := range tags {
		tag := MarshalTag(a)
		err := r.getDB().Create(&tag).Error
		if err != nil {
			return nil, err
		}
		// update assosiation table

		newTag, err := r.ReadByName(tag.Name)
		if err != nil {
			return nil, err
		}
		resp = append(resp, *newTag)
	}
	return resp, nil
}

func (r *TagRepo) ReadByName(name interface{}) (*models.Tag, error) {
	var a Tag
	err := r.db.Model(&Tag{}).Where(&Tag{Name: name.(string)}).First(&a).Error
	if err != nil {
		return nil, err
	}
	m := UnmarshalTag(a)
	return &m, nil
}

func (r *TagRepo) GetID(name string) (string, error) {
	var a Tag
	if err := r.db.Model(&Tag{}).Where(&Tag{Name: name}).First(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("GetID: id of '%s' not found", name)
		}
		return "", err
	}
	return a.ID, nil
}

func (r *TagRepo) ReadAll() ([]models.Tag, error) {
	var a []Tag
	err := r.db.Model(&Tag{}).Find(&a).Error
	if err != nil {
		return nil, err
	}
	m := UnmarshalArrayTag(a)
	return m, nil
}

func MarshalTag(a models.Tag) Tag {
	return Tag{
		DBModel:     DBModel(a.DBInfo),
		Name:        a.Name,
		Description: a.Description,
	}
}

func UnmarshalTag(a Tag) models.Tag {
	return models.Tag{
		DBInfo:      models.DBInfo(a.DBModel),
		Name:        a.Name,
		Description: a.Description,
	}
}

func MarshalArrayTag(m []models.Tag) []Tag {
	tags := []Tag{}
	for _, a := range m {
		tags = append(tags, MarshalTag(a))
	}
	return tags
}

func UnmarshalArrayTag(a []Tag) []models.Tag {
	tags := []models.Tag{}
	for _, m := range a {
		tags = append(tags, UnmarshalTag(m))
	}
	return tags
}

type TagDBMigrator struct {
	db *gorm.DB
}

func (s TagDBMigrator) Migrate() error {
	err := s.db.AutoMigrate(&Tag{})
	if err != nil {
		return err
	}
	return nil
}

func NewTagDBMigrator(db *gorm.DB) *TagDBMigrator {
	return &TagDBMigrator{db: db}
}
