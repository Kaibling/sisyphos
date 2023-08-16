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
	Name           *string               `gorm:"unique"`
	Actions        []string              `gorm:"-"`
	ActionsRef     []Action              `gorm:"many2many:actions_actions;"`
	OrderedActions []models.OrderdAction `gorm:"-"`
	Groups         []string              `gorm:"-"`
	GroupsRef      []Group               `gorm:"many2many:groups_actions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	OrderedHosts   []models.OrderedHost  `gorm:"-"`
	HostsRef       []Host                `gorm:"many2many:actions_hosts;"`
	TagsRef        []Tag                 `gorm:"many2many:actions_tags;"`
	Tags           []string              `gorm:"-"`
	Script         *string
	FailOnErrors   *bool
	Variables      datatypes.JSON
}

func (a *Action) BeforeSave(tx *gorm.DB) (err error) {
	triggers := []Action{}
	for _, t := range a.Actions {
		actionRepo := NewActionRepo(tx)
		tid, err := actionRepo.GetID(t)
		if err != nil {
			return err
		}
		triggers = append(triggers, Action{DBModel: DBModel{ID: tid}})
	}
	a.ActionsRef = triggers

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
	for _, s := range a.ActionsRef {
		triggers = append(triggers, *s.Name)
	}
	a.Actions = triggers

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
}

type ActionsActions struct {
	ActionsRefID string `gorm:"type:varchar(191);primaryKey"`
	ActionID     string `gorm:"type:varchar(191);primaryKey"`
	Order        int
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
		if action.Script == nil {
			action.Script = utils.ToPointer("")
		}
		if action.Variables == nil {
			action.Variables.UnmarshalJSON([]byte("{}"))
		}
		err := tx.Omit("HostsRef").Omit("ActionsRef").Omit("TagsRef.*").Omit("GroupsRef.*").Create(&action).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		for _, c := range action.OrderedHosts {
			hostRepo := NewHostRepo(tx)
			hostID, err := hostRepo.GetID(utils.PtrRead(&c.Name))
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

		for _, a := range action.OrderedActions {
			actionRepo := NewActionRepo(tx)
			actionID, err := actionRepo.GetID(a.Name)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			newActionAction := ActionsActions{
				ActionsRefID: actionID,
				ActionID:     action.ID,
				Order:        a.Order,
			}
			if err := tx.Model(&ActionsActions{}).Clauses(clause.OnConflict{UpdateAll: true}).Create(&newActionAction).Error; err != nil {
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
	var a Action
	err := r.db.Model(&Action{}).Where("name = ?", name.(string)).Preload("HostsRef").Preload("ActionsRef").Preload("TagsRef").Preload("GroupsRef").First(&a).Error
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
				conns = append(conns, models.OrderedHost{Name: *host.Name, Order: rel.Order})
			}
		}
	}
	a.OrderedHosts = conns

	aa := []ActionsActions{}
	err = r.db.Model(&ActionsActions{}).Where("action_id = ?", a.ID).Find(&aa).Error
	if err != nil {
		return nil, err
	}

	orderActions := []models.OrderdAction{}
	for _, action := range a.ActionsRef {
		for _, rel := range aa {
			if rel.ActionsRefID == action.ID {
				orderActions = append(orderActions, models.OrderdAction{Name: *action.Name, Order: rel.Order})
			}
		}
	}
	a.OrderedActions = orderActions
	m := UnmarshalAction(a)
	return &m, nil
}
func (r *ActionRepo) ReadIDs(ids []interface{}) ([]models.Action, error) {
	var actions []Action
	if err := r.db.Model(&Action{}).Where("id IN ?", ids).Preload("HostsRef").Preload("ActionsRef").Preload("TagsRef").Preload("GroupsRef").Find(&actions).Error; err != nil {
		return nil, err
	}
	actionExt := []models.Action{}
	for _, action := range actions {
		ah := []ActionsHosts{}
		if err := r.db.Model(&ActionsHosts{}).Where("action_id = ?", action.ID).Find(&ah).Error; err != nil {
			return nil, err
		}

		conns := []models.OrderedHost{}
		for _, host := range action.HostsRef {
			for _, rel := range ah {
				if rel.HostID == host.ID {
					conns = append(conns, models.OrderedHost{Name: *host.Name, Order: rel.Order})
				}
			}
		}
		action.OrderedHosts = conns

		aa := []ActionsActions{}
		if err := r.db.Model(&ActionsActions{}).Where("action_id = ?", action.ID).Find(&aa).Error; err != nil {
			return nil, err
		}

		orderActions := []models.OrderdAction{}
		for _, action := range action.ActionsRef {
			for _, rel := range aa {
				if rel.ActionsRefID == action.ID {
					orderActions = append(orderActions, models.OrderdAction{Name: *action.Name, Order: rel.Order})
				}
			}
		}
		action.OrderedActions = orderActions
		actionExt = append(actionExt, UnmarshalAction(action))
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

func (r *ActionRepo) ReadRuns(actionname interface{}) ([]models.Run, error) {
	var a []Run
	err := r.getDB().Model(&Run{}).Preload("Action").Preload("Host").Joins("JOIN actions on runs.action_id = actions.id").Where("actions.name = ?", actionname.(string)).Find(&a).Error
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

	if uAction.OrderedActions != nil {
		for _, conn := range uAction.OrderedActions {
			actionRepo := NewActionRepo(tx)
			actionID, err := actionRepo.GetID(conn.Name)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			aa := ActionsActions{ActionsRefID: actionID, ActionID: uAction.ID, Order: conn.Order}
			err = tx.Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]interface{}{"order": conn.Order}),
			}).Model(&ActionsActions{}).Create(&aa).Error
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}
	if uAction.Groups != nil {
		if err := tx.Model(&uAction).Association("GroupsRef").Replace(uAction.GroupsRef); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if uAction.OrderedHosts != nil {
		for _, conn := range uAction.OrderedHosts {
			hostRepo := NewHostRepo(tx)
			hostID, err := hostRepo.GetID(conn.Name)
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
	orderHost := []models.OrderedHost{}
	for _, h := range a.Hosts {
		orderHost = append(orderHost, models.OrderedHost{Name: h.Name, Order: h.Order})
	}
	orderAction := []models.OrderdAction{}
	for _, a := range a.Actions {
		orderAction = append(orderAction, models.OrderdAction{Name: a.Name, Order: a.Order})
	}

	triggers := []string{}
	for _, h := range a.Actions {
		triggers = append(triggers, h.Name)
	}

	return Action{
		Name:           a.Name,
		Script:         a.Script,
		Groups:         a.Groups,
		Actions:        triggers,
		OrderedHosts:   orderHost,
		OrderedActions: orderAction,
		Tags:           a.Tags,
		Variables:      b,
		FailOnErrors:   a.FailOnErrors,
	}
}

func UnmarshalAction(a Action) models.Action {
	v := map[string]interface{}{}
	// if a.Variables.String() != "null" {
	byteVar, err := a.Variables.MarshalJSON()
	if err != nil {
		fmt.Println(err.Error())
	}
	err = json.Unmarshal(byteVar, &v)
	if err != nil {
		fmt.Println(err.Error())
	}
	// }
	if len(a.Tags) == 0 {
		a.Tags = []string{}
	}

	return models.Action{
		Name:         a.Name,
		Script:       a.Script,
		Actions:      a.OrderedActions,
		Tags:         a.Tags,
		Variables:    v,
		Groups:       a.Groups,
		FailOnErrors: a.FailOnErrors,
		Hosts:        a.OrderedHosts,
	}
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

	err = s.db.AutoMigrate(&ActionsActions{})
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
