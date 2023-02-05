package gormrepo

import (
	"encoding/json"
	"errors"
	"fmt"

	"sisyphos/lib/utils"
	"sisyphos/models"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Action struct {
	DBModel
	Name        *string `gorm:"unique"`
	Script      *string
	Triggers    []string            `gorm:"-"`
	TriggersRef []Action            `gorm:"many2many:action_triggers;"`
	Groups      []string            `gorm:"-"`
	GroupsRef   []Group             `gorm:"many2many:groups_actions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Connections []models.Connection `gorm:"-"`
	HostsRef    []Host              `gorm:"many2many:actions_hosts;"`
	Variables   datatypes.JSON
	TagsRef     []Tag    `gorm:"many2many:actions_tags;"`
	Tags        []string `gorm:"-"`
}

func (a *Action) BeforeSave(tx *gorm.DB) (err error) {
	triggers := []Action{}
	for _, t := range a.Triggers {
		actionRepo := NewActionRepo(tx)
		tid, err := actionRepo.GetID(t)
		if err != nil {
			return err
		}
		triggers = append(triggers, Action{DBModel: DBModel{ID: tid}})
	}
	a.TriggersRef = triggers

	groups := []Group{}
	for _, g := range a.Groups {
		groupRepo := NewGroupRepo(tx)
		gid, err := groupRepo.GetID(g)
		if err != nil {
			return err
		}
		groups = append(groups, Group{DBModel: DBModel{ID: gid}})
	}
	a.GroupsRef = groups

	tags := []Tag{}
	for _, t := range a.Tags {
		tagRepo := NewTagRepo(tx)
		tid, err := tagRepo.GetID(t)
		if err != nil {
			return err
		}
		tags = append(tags, Tag{DBModel: DBModel{ID: tid}})
	}
	a.TagsRef = tags
	return
}

func (a *Action) AfterFind(tx *gorm.DB) (err error) {
	triggers := []string{}
	for _, s := range a.TriggersRef {
		triggers = append(triggers, *s.Name)
	}
	a.Triggers = triggers

	groups := []string{}
	for _, g := range a.GroupsRef {
		groups = append(groups, *g.Name)
	}
	a.Groups = groups

	tags := []string{}
	for _, t := range a.TagsRef {
		tags = append(tags, t.Name)
	}
	a.Tags = tags
	return
}

type ActionsHosts struct {
	HostID   string `gorm:"type:varchar(191);primaryKey"`
	ActionID string `gorm:"type:varchar(191);primaryKey"`
	Port     string
	Order    int
	Name     string
}

type ActionRepo struct {
	db *gorm.DB
}

func NewActionRepo(db *gorm.DB) *ActionRepo {
	return &ActionRepo{db}
}

func (r *ActionRepo) getDB() *gorm.DB {
	d := r.db.Session(&gorm.Session{NewDB: true})
	return d
}

func (r *ActionRepo) Create(actions []models.Action) ([]models.Action, error) {
	newIDs := []any{}
	tx := r.getDB().Begin()
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e.(error).Error())
			tx.Rollback()
		}
	}()
	for _, a := range actions {
		action := MarshalAction(a)
		err := tx.Omit("HostsRef.*").Omit("TriggersRef.*").Omit("TagsRef.*").Omit("GroupsRef.*").Create(&action).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		for _, c := range action.Connections {
			hostRepo := NewHostRepo(tx)
			hostID, err := hostRepo.GetID(c.Name)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			newActionHost := ActionsHosts{
				HostID:   hostID,
				ActionID: action.ID,
				Port:     c.Port,
				Order:    c.Order,
			}

			if err := tx.Model(&ActionsHosts{}).Create(&newActionHost).Error; err != nil {
				if err != nil {
					tx.Rollback()
					return nil, err
				}
			}

		}
		newIDs = append(newIDs, action.ID)
	}
	tx.Commit()
	return r.ReadIDs(newIDs)
}

func (r *ActionRepo) ReadByName(name interface{}) (*models.Action, error) {
	a, err := r.ReadExt(name)
	if err != nil {
		return nil, err
	}
	return ReduceExtended(a), nil
}

func (r *ActionRepo) ReadExt(name interface{}) (*models.ActionExt, error) {
	var a Action
	err := r.db.Model(&Action{}).Where("name = ?", name.(string)).Preload("HostsRef").Preload("TriggersRef").Preload("TagsRef").Preload("GroupsRef").First(&a).Error
	if err != nil {
		return nil, err
	}
	ah := []ActionsHosts{}
	err = r.db.Model(&ActionsHosts{}).Where("action_id = ?", a.ID).Find(&ah).Error
	if err != nil {
		return nil, err
	}

	conns := []models.Connection{}
	for _, host := range a.HostsRef {
		for _, rel := range ah {
			if rel.HostID == host.ID {
				mh := UnmarshalHost(host)
				conns = append(conns, models.Connection{Host: mh, Port: rel.Port, Order: rel.Order})
			}
		}
	}
	a.Connections = conns
	m := UnmarshalActionExt(a)
	return &m, nil
}
func (r *ActionRepo) ReadExtIDs(ids []interface{}) ([]models.ActionExt, error) {
	var actions []Action
	if err := r.db.Model(&Action{}).Where("id IN ?", ids).Preload("HostsRef").Preload("TriggersRef").Preload("TagsRef").Preload("GroupsRef").Find(&actions).Error; err != nil {
		return nil, err
	}
	actionExt := []models.ActionExt{}
	for _, action := range actions {
		ah := []ActionsHosts{}
		if err := r.db.Model(&ActionsHosts{}).Where("action_id = ?", action.ID).Find(&ah).Error; err != nil {
			return nil, err
		}

		conns := []models.Connection{}
		for _, host := range action.HostsRef {
			for _, rel := range ah {
				if rel.HostID == host.ID {
					mh := UnmarshalHost(host)
					conns = append(conns, models.Connection{Host: mh, Port: rel.Port, Order: rel.Order})
				}
			}
		}
		action.Connections = conns
		actionExt = append(actionExt, UnmarshalActionExt(action))
	}
	return actionExt, nil
}
func (r *ActionRepo) GetID(name string) (string, error) {
	var a Action
	if err := r.db.Model(&Action{}).Where(&Action{Name: &name}).First(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("GetID: id of '%s' not found", name)
		}
		return "", err
	}
	return a.ID, nil
}

func (r *ActionRepo) GetHostID(name string) (string, error) {
	var a Host
	err := r.db.Model(&Host{}).Where(&Host{Name: name}).Find(&a).Error
	if err != nil {
		return "", err
	}
	return a.ID, nil
}

func (r *ActionRepo) GetTagID(name string) (string, error) {
	var a Tag
	err := r.db.Model(&Tag{}).Where(&Tag{Name: name}).First(&a).Error
	if err != nil {
		return "", err
	}
	return a.ID, nil
}

func (r *ActionRepo) ReadIDs(ids []interface{}) ([]models.Action, error) {
	ext, err := r.ReadExtIDs(ids)
	if err != nil {
		return nil, err
	}
	actions := []models.Action{}
	for _, a := range ext {
		actions = append(actions, *ReduceExtended(&a))
	}
	return actions, nil
}

func (r *ActionRepo) ReadRuns(actionname interface{}) ([]models.Run, error) {
	var a []Run
	err := r.getDB().Model(&Run{}).Joins("JOIN actions on runs.action_id = actions.id").Where("actions.name = ?", actionname.(string)).Find(&a).Error
	if err != nil {
		return nil, err
	}
	return UnmarshalArrayRun(a), nil
}

// func (r *ActionRepo) ReadAll() ([]models.Action, error) {
// 	var a []Action
// 	err := r.db.Model(&Action{}).Preload("HostsRef").Preload("TriggersRef").Preload("TagsRef").Preload("GroupsRef").Find(&a).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	m := UnmarshalArrayAction(a)
// 	return m, nil
// }

func (r *ActionRepo) Update(name string, d *models.Action) (*models.Action, error) {
	uAction := MarshalAction(*d)
	if uid, err := r.GetID(name); err != nil {
		return nil, err
	} else {
		uAction.ID = uid
	}
	if uAction.Triggers != nil {
		if err := r.getDB().Model(&uAction).Association("TriggersRef").Replace(uAction.TriggersRef); err != nil {
			return nil, err
		}
	}
	if uAction.Groups != nil {
		if err := r.getDB().Model(&uAction).Association("GroupsRef").Replace(uAction.GroupsRef); err != nil {
			return nil, err
		}
	}
	if uAction.Connections != nil {
		for _, conn := range uAction.Connections {
			hostID, err := r.GetHostID(conn.Name)
			if err != nil {
				return nil, err
			}
			ah := ActionsHosts{HostID: hostID, ActionID: uAction.ID, Port: conn.Port}
			err = r.getDB().Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]interface{}{"port": conn.Port}),
			}).Model(&ActionsHosts{}).Create(&ah).Error
			if err != nil {
				return nil, err
			}
		}
	}

	if err := r.getDB().Omit("AllowsRef.*").Omit("UsersRef.*").Updates(&uAction).Error; err != nil {
		return nil, err
	}
	return r.ReadByName(name)
}

func MarshalAction(a models.Action) Action {
	b, err := json.Marshal(a.Variables)
	if err != nil {
		fmt.Println(err.Error())
	}
	conns := []models.Connection{}
	for _, h := range a.Hosts {
		conns = append(conns, models.Connection{Host: models.Host{Name: h.HostName}, Order: h.Order, Port: h.Port})
	}

	return Action{
		Name:        a.Name,
		Script:      a.Script,
		Groups:      a.Groups,
		Triggers:    a.Triggers,
		Connections: conns,
		Tags:        a.Tags,
		Variables:   b,
	}
}

func UnmarshalAction(a Action) models.Action {
	v := map[string]interface{}{}
	if a.Variables.String() != "null" {
		err := json.Unmarshal(a.Variables.MarshalJSON())
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	if len(a.Tags) == 0 {
		a.Tags = []string{}
	}

	return models.Action{
		Name:      utils.PtrDefault(a.Name),
		Script:    utils.PtrDefault(a.Script),
		Triggers:  a.Triggers,
		Tags:      a.Tags,
		Variables: v,
		Groups:    a.Groups,
		//Hosts: a.Connections,
	}
}

func ReduceExtended(m *models.ActionExt) *models.Action {
	triggers := []string{}
	for _, t := range m.Triggers {
		triggers = append(triggers, *t.Name)
	}
	services := []models.Service{}
	for _, c := range m.Hosts {
		services = append(services, models.Service{HostName: c.Name, Port: c.Port})
	}
	if len(m.Tags) == 0 {
		m.Tags = []string{}
	}
	return &models.Action{
		Name:      m.Name,
		Groups:    m.Groups,
		Script:    m.Script,
		Tags:      m.Tags,
		Triggers:  triggers,
		Hosts:     services,
		Variables: m.Variables,
	}
}

func UnmarshalActionExt(a Action) models.ActionExt {
	v := map[string]interface{}{}
	if a.Variables.String() != "null" {
		err := json.Unmarshal(a.Variables.MarshalJSON())
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return models.ActionExt{
		Name:      a.Name,
		Script:    a.Script,
		Groups:    a.Groups,
		Triggers:  UnmarshalArrayActionExt(a.TriggersRef),
		Hosts:     a.Connections,
		Tags:      a.Tags,
		Variables: v,
	}
}

func UnmarshalArrayActionExt(a []Action) []models.ActionExt {
	actions := []models.ActionExt{}
	for _, m := range a {
		actions = append(actions, UnmarshalActionExt(m))
	}
	return actions
}

func UnmarshalArrayAction(a []Action) []models.Action {
	actions := []models.Action{}
	for _, m := range a {
		actions = append(actions, UnmarshalAction(m))
	}
	return actions
}

type ActionDBMigrator struct {
	db *gorm.DB
}

func (s ActionDBMigrator) Migrate() error {
	err := s.db.AutoMigrate(&ActionsHosts{})
	if err != nil {
		return err
	}

	err = s.db.AutoMigrate(&Action{})
	if err != nil {
		return err
	}

	return nil
}

func NewActionMigrator(db *gorm.DB) *ActionDBMigrator {
	return &ActionDBMigrator{db: db}
}
