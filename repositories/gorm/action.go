package gormrepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"sisyphos/lib/reqctx"
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
	ScheduleExpr   *string
	Script         *string
	FailOnErrors   *bool
	Variables      datatypes.JSON
}

func (a *Action) BeforeSave(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	username, ok := ctx.Value(reqctx.String("username")).(string)
	if !ok {
		return fmt.Errorf("before hook: username is missing in transaction context")
	}
	triggers := []Action{}
	for _, t := range a.Actions {
		actionRepo := NewActionRepo(tx, username)
		tid, err := actionRepo.GetID(t)
		if err != nil {
			return err
		}
		triggers = append(triggers, Action{DBModel: DBModel{ID: tid}})
	}
	a.ActionsRef = triggers

	groups := []Group{}
	for _, g := range a.Groups {
		groupRepo := NewGroupRepo(tx, username)
		gid, err := groupRepo.GetID(g)
		if err != nil {
			return err
		}
		groups = append(groups, Group{DBModel: DBModel{ID: gid}})
	}
	a.GroupsRef = groups

	tags := []Tag{}
	for _, t := range a.Tags {
		tagRepo := NewTagRepo(tx, username)
		tid, err := tagRepo.GetID(t)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				nt, err := tagRepo.Create([]models.Tag{{Name: t}})
				if err != nil {
					return err
				}
				tid, err = tagRepo.GetID(nt[0].Name)
				if err != nil {
					return err
				}
			} else {
				return err
			}
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
	db       *gorm.DB
	username string
}

func NewActionRepo(db *gorm.DB, username string) *ActionRepo {
	ctx := context.WithValue(context.TODO(), reqctx.String("username"), username)
	db = db.WithContext(ctx)
	return &ActionRepo{db, username}
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
			fmt.Printf("creation failed. rollback.: %s\n", e.(error).Error())
			tx.Rollback()
		}
	}()
	for _, a := range actions {
		action := MarshalAction(a)
		action.UpdatedBy = r.username
		action.CreatedBy = r.username
		if action.Script == nil {
			action.Script = utils.ToPointer("")
		}
		if action.Variables == nil {
			_ = action.Variables.UnmarshalJSON([]byte("{}"))
		}
		err := tx.Omit("HostsRef").Omit("ActionsRef").Omit("TagsRef.*").Omit("GroupsRef.*").Create(&action).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		for _, c := range action.OrderedHosts {
			hostRepo := NewHostRepo(tx, r.username)
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
			actionRepo := NewActionRepo(tx, r.username)
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
	err := r.getDB().Model(&Action{}).Where("name = ?", name.(string)).Preload("HostsRef").Preload("ActionsRef").Preload("TagsRef").Preload("GroupsRef").First(&a).Error
	if err != nil {
		return nil, err
	}
	ah := []ActionsHosts{}
	err = r.getDB().Model(&ActionsHosts{}).Where("action_id = ?", a.ID).Find(&ah).Error
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
	err = r.getDB().Model(&ActionsActions{}).Where("action_id = ?", a.ID).Find(&aa).Error
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
	if err := r.getDB().Model(&Action{}).Where("id IN ?", ids).Preload("HostsRef").Preload("ActionsRef").Preload("TagsRef").Preload("GroupsRef").Find(&actions).Error; err != nil {
		return nil, err
	}
	actionExt := []models.Action{}
	for _, action := range actions {
		ah := []ActionsHosts{}
		if err := r.getDB().Model(&ActionsHosts{}).Where("action_id = ?", action.ID).Find(&ah).Error; err != nil {
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
		if err := r.getDB().Model(&ActionsActions{}).Where("action_id = ?", action.ID).Find(&aa).Error; err != nil {
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
	if err := r.getDB().Model(&Action{}).Where(&Action{Name: &name}).First(&a).Error; err != nil {
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
	uAction.UpdatedBy = r.username
	uAction.UpdatedAt = time.Now()
	tx := r.getDB().Begin()
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("update failed.rollback.: %s\n", e.(error).Error())
			tx.Rollback()
		}
	}()

	if uAction.OrderedActions != nil {
		// delete all
		err := tx.Exec("delete from actions_actions where action_id = ?", uAction.ID).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		// add remaining
		for _, conn := range uAction.OrderedActions {
			actionRepo := NewActionRepo(tx, r.username)
			actionID, err := actionRepo.GetID(conn.Name)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			aa := ActionsActions{ActionsRefID: actionID, ActionID: uAction.ID, Order: conn.Order}
			err = tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "action__ref_id"},
					{Name: "action_id"},
				},
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

	if uAction.Tags != nil {
		if err := tx.Model(&uAction).Association("TagsRef").Replace(uAction.TagsRef); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if uAction.OrderedHosts != nil {
		// delete all
		err := tx.Exec("delete from actions_hosts where action_id = ?", uAction.ID).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		// add remaining
		hostRepo := NewHostRepo(tx, r.username)
		for _, conn := range uAction.OrderedHosts {
			hostID, err := hostRepo.GetID(conn.Name)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			ah := ActionsHosts{HostID: hostID, ActionID: uAction.ID, Order: conn.Order}
			err = tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "host_id"},
					{Name: "action_id"},
				},
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

func (r *ActionRepo) ReadToBeScheduled() ([]models.Action, error) {
	var a []Action
	err := r.getDB().Model(&Action{}).Preload("HostsRef").Where("actions.schedule_expr is not null").Find(&a).Error
	if err != nil {
		return nil, err
	}
	return UnmarshalArrayAction(a), nil
}

func MarshalAction(a models.Action) Action {
	b, err := json.Marshal(a.Variables)
	if err != nil {
		fmt.Printf("MarshalAction: %s\n", err.Error())
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
		DBModel:        DBModel(a.DBInfo),
		Name:           a.Name,
		Script:         a.Script,
		Groups:         a.Groups,
		Actions:        triggers,
		OrderedHosts:   orderHost,
		OrderedActions: orderAction,
		Tags:           a.Tags,
		Variables:      b,
		FailOnErrors:   a.FailOnErrors,
		ScheduleExpr:   a.ScheduleExpr,
	}
}

func UnmarshalAction(a Action) models.Action {
	v := map[string]any{}
	byteVar, err := a.Variables.MarshalJSON()
	if err != nil {
		// TODO error handling
		fmt.Printf("UnmarshalAction: %s\n", err.Error())
	}
	aa := map[string]any{}
	err = json.Unmarshal(byteVar, &aa)
	if err != nil {
		// TODO error handling
		fmt.Printf("UnmarshalAction: %s\n", err.Error())
	}
	if len(a.Tags) == 0 {
		a.Tags = []string{}
	}

	return models.Action{
		DBInfo:       models.DBInfo(a.DBModel),
		Name:         a.Name,
		Script:       a.Script,
		Actions:      a.OrderedActions,
		Tags:         a.Tags,
		Variables:    v,
		Groups:       a.Groups,
		FailOnErrors: a.FailOnErrors,
		Hosts:        a.OrderedHosts,
		ScheduleExpr: a.ScheduleExpr,
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
