package gormrepo

import (
	"errors"
	"fmt"

	"sisyphos/lib/utils"

	"gorm.io/gorm"
)

type PermissionRepo struct {
	db *gorm.DB
}

func NewPermissionRepo(db *gorm.DB) *PermissionRepo {
	return &PermissionRepo{db}
}

func (r *PermissionRepo) getDB() *gorm.DB {
	d := r.db.Session(&gorm.Session{NewDB: true})
	return d
}

func (r *PermissionRepo) GetActionIDs(username interface{}) ([]any, error) {
	var ids []any
	q := r.getDB().Select("actions.id")
	if username.(string) != "admin" {
		q.Model(&User{}).
			Joins("LEFT JOIN groups_users ON groups_users.user_id = users.id").
			Joins("LEFT JOIN `groups` ON groups_users.group_id = groups.id").
			Joins("LEFT JOIN groups_actions ON groups.id = groups_actions.group_id").
			Joins("LEFT JOIN actions ON actions.id = groups_actions.action_id").
			Where(&User{Name: username.(string)})
	} else {
		q.Model(&Action{})
	}
	if err := q.Find(&ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *PermissionRepo) GetUserIDs(username interface{}) ([]any, error) {
	var ids []any
	if username.(string) != "admin" {
		var tempUser User
		if err := r.getDB().Model(&User{}).
			Joins("LEFT JOIN groups_users ON groups_users.user_id = users.id").
			Joins("LEFT JOIN `groups` ON groups_users.group_id = groups.id").
			Where(&Group{Name: utils.ToPointer("admin")}).
			Where(&User{Name: username.(string)}).First(&tempUser).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("forbidden")
			}
			return nil, err
		}
	}
	if err := r.getDB().Select("users.id").Model(&User{}).Find(&ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}
