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
	Name         *string              `gorm:"unique"`
	Triggers     []string             `gorm:"-"`
	TriggersRef  []Action             `gorm:"many2many:action_triggers;"`
	Groups       []string             `gorm:"-"`
	GroupsRef    []Group              `gorm:"many2many:groups_actions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	OrderedHost  []models.OrderedHost `gorm:"-"`
	HostsRef     []Host               `gorm:"many2many:actions_hosts;"`
	TagsRef      []Tag                `gorm:"many2many:actions_tags;"`
	Tags         []string             `gorm:"-"`
	Script       *string
	FailOnErrors *bool
	Variables    datatypes.JSON
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
	Order    int
	Name     string
}

// type ActionsActionss struct {
// 	HostID   string `gorm:"type:varchar(191);primaryKey"`
// 	ActionID string `gorm:"type:varchar(191);primaryKey"`
// 	Order    int
// 	Name     string
// }

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
		for _, c := range action.OrderedHost {
			hostRepo := NewHostRepo(tx)
			hostID, err := hostRepo.GetID(utils.PtrRead(&c.HostName))
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			newActionHost := ActionsHosts{
				HostID:   hostID,
				ActionID: action.ID,
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

	conns := []models.OrderedHost{}
	for _, host := range a.HostsRef {
		for _, rel := range ah {
			if rel.HostID == host.ID {
				//mh := UnmarshalHost(host)
				conns = append(conns, models.OrderedHost{HostName: *host.Name, Order: rel.Order})
			}
		}
	}
	a.OrderedHost = conns
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

		conns := []models.OrderedHost{}
		for _, host := range action.HostsRef {
			for _, rel := range ah {
				if rel.HostID == host.ID {
					// mh := UnmarshalHost(host)
					conns = append(conns, models.OrderedHost{HostName: *host.Name, Order: rel.Order})
				}
			}
		}
		action.OrderedHost = conns
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

func (r *ActionRepo) Update(name string, d *models.Action) (*models.Action, error) {
	uAction := MarshalAction(*d)
	if uid, err := r.GetID(name); err != nil {
		return nil, err
	} else {
		uAction.ID = uid
	}
	tx := r.getDB().Begin()
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e.(error).Error())
			tx.Rollback()
		}
	}()
	if uAction.Triggers != nil {
		if err := tx.Model(&uAction).Association("TriggersRef").Replace(uAction.TriggersRef); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if uAction.Groups != nil {
		if err := tx.Model(&uAction).Association("GroupsRef").Replace(uAction.GroupsRef); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if uAction.OrderedHost != nil {
		for _, conn := range uAction.OrderedHost {
			hostRepo := NewHostRepo(tx)
			hostID, err := hostRepo.GetID(conn.HostName)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			ah := ActionsHosts{HostID: hostID, ActionID: uAction.ID, Order: conn.Order}
			err = tx.Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]interface{}{"order": conn.Order}),
			}).Model(&ActionsHosts{}).Create(&ah).Error
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if err := tx.Omit("AllowsRef.*").Omit("UsersRef.*").Omit("GroupsRef.*").Updates(&uAction).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return r.ReadByName(name)
}

func MarshalAction(a models.Action) Action {
	b, err := json.Marshal(a.Variables)
	if err != nil {
		fmt.Println(err.Error())
	}
	conns := []models.OrderedHost{}
	for _, h := range a.Hosts {
		conns = append(conns, models.OrderedHost{HostName: h.HostName, Order: h.Order})
	}

	triggers := []string{}
	for _, h := range a.Actions {
		triggers = append(triggers, h.Action)
	}

	return Action{
		Name:         a.Name,
		Script:       a.Script,
		Groups:       a.Groups,
		Triggers:     triggers,
		OrderedHost:  conns,
		Tags:         a.Tags,
		Variables:    b,
		FailOnErrors: a.FailOnErrors,
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

	triggers := []models.OrderdAction{}
	for _, h := range a.Triggers {
		triggers = append(triggers, models.OrderdAction{Action: h})
	}

	return models.Action{
		Name:         a.Name,
		Script:       a.Script,
		Actions:      triggers,
		Tags:         a.Tags,
		Variables:    v,
		Groups:       a.Groups,
		FailOnErrors: a.FailOnErrors,
		//Hosts: a.Connections,
	}
}

func ReduceExtended(m *models.ActionExt) *models.Action {
	triggers := []models.OrderdAction{}
	for _, h := range m.Triggers {
		triggers = append(triggers, models.OrderdAction{Action: h.ActionExt})
	}
	services := []models.OrderedHost{}
	for _, c := range m.Hosts {
		services = append(services, models.OrderedHost{HostName: c.HostName, Order: c.Order})
	}
	if len(m.Tags) == 0 {
		m.Tags = []string{}
	}
	return &models.Action{
		Name:         m.Name,
		Groups:       m.Groups,
		Script:       m.Script,
		Tags:         m.Tags,
		Actions:      triggers,
		Hosts:        services,
		Variables:    m.Variables,
		FailOnErrors: m.FailOnErrors,
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

	orderedActionExt := []models.OrderdActionExt{}
	actions := UnmarshalArrayActionExt(a.TriggersRef)
	for _, a := range actions {
		orderedActionExt = append(orderedActionExt, models.OrderdActionExt{ActionExt: *a.Name})
	}

	return models.ActionExt{
		Name:         a.Name,
		Script:       a.Script,
		Groups:       a.Groups,
		Triggers:     orderedActionExt,
		Hosts:        a.OrderedHost,
		Tags:         a.Tags,
		Variables:    v,
		FailOnErrors: a.FailOnErrors,
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
