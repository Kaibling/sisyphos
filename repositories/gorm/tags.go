package gormrepo

import (
	"sisyphos/models"

	"gorm.io/gorm"
)

type Tag struct {
	DBModel
	Name string
	Text string
}

type TagRepo struct {
	db *gorm.DB
}

func NewTagRepo(db *gorm.DB) *TagRepo {
	return &TagRepo{db}
}

func (r *TagRepo) getDB() *gorm.DB {
	d := r.db.Session(&gorm.Session{NewDB: true})
	return d
}

func (r *TagRepo) Create(actions []models.Tag) ([]models.Tag, error) {
	resp := []models.Tag{}
	for _, a := range actions {
		action := MarshalTag(a)

		err := r.getDB().Create(&action).Error
		if err != nil {
			return nil, err
		}
		// update assosiation table

		newTag, err := r.ReadByName(action.Name)
		if err != nil {
			return nil, err
		}
		resp = append(resp, *newTag)
	}
	return resp, nil
}

func (r *TagRepo) ReadByName(name interface{}) (*models.Tag, error) {
	var a Tag
	err := r.db.Model(&Tag{Name: name.(string)}).First(&a).Error
	if err != nil {
		return nil, err
	}
	m := UnmarshalTag(a)
	return &m, nil
}

func (r *TagRepo) GetID(name string) (string, error) {
	var a Tag
	err := r.db.Model(&Tag{}).Where(&Tag{Name: name}).Find(&a).Error
	if err != nil {
		return "", err
	}
	return a.Name, nil
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
		Name: a.Name,
		Text: a.Text,
	}
}

func UnmarshalTag(a Tag) models.Tag {
	return models.Tag{
		Name: a.Name,
		Text: a.Text,
	}
}

func MarshalArrayTag(m []models.Tag) []Tag {
	actions := []Tag{}
	for _, a := range m {
		actions = append(actions, MarshalTag(a))
	}
	return actions
}

func UnmarshalArrayTag(a []Tag) []models.Tag {
	actions := []models.Tag{}
	for _, m := range a {
		actions = append(actions, UnmarshalTag(m))
	}
	return actions
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
