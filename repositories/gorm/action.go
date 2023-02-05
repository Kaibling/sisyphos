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
	Name        string `gorm:"unique"`
	Script      string
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
	return
}

func (a *Action) AfterFind(tx *gorm.DB) (err error) {
	triggers := []string{}
	for _, s := range a.TriggersRef {
		triggers = append(triggers, s.Name)
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
	resp := []models.Action{}
	for _, a := range actions {
		action := MarshalAction(a)
		utils.PrettyJSON(action)
		err := r.getDB().Omit("HostsRef.*").Omit("TriggersRef.*").Omit("TagsRef.*").Omit("GroupsRef.*").Create(&action).Error
		if err != nil {
			return nil, err
		}
		for _, c := range action.Connections {
			hostRepo := NewHostRepo(r.db)
			hostID, err := hostRepo.GetID(c.Name)
			if err != nil {
				return nil, err
			}
			newActionHost := ActionsHosts{
				HostID:   hostID,
				ActionID: action.ID,
				Port:     c.Port,
				Order:    c.Order,
			}

			if err := r.db.Model(&ActionsHosts{}).Create(&newActionHost).Error; err != nil {
				if err != nil {
					return nil, err
				}
			}
		}

		newAction, err := r.ReadByName(action.Name)
		if err != nil {
			return nil, err
		}
		resp = append(resp, *newAction)
	}
	return resp, nil
}

func (r *ActionRepo) ReadByName(name interface{}) (*models.Action, error) {
	a, err := r.ReadExtendedv3(name)
	if err != nil {
		return nil, err
	}
	return ReduceExtended(a), nil
}

func (r *ActionRepo) ReadExtendedv3(name interface{}) (*models.ActionExtendedv3, error) {
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
	m := UnmarshalActionExtendedv3(a)
	return &m, nil
}

func (r *ActionRepo) GetID(name string) (string, error) {
	var a Action
	if err := r.db.Model(&Action{}).Where(&Action{Name: name}).First(&a).Error; err != nil {
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
	var a []Action
	err := r.db.Model(&Action{}).Where("id IN ?", ids).Preload("HostsRef").Preload("TriggersRef").Preload("TagsRef").Preload("GroupsRef").Find(&a).Error
	if err != nil {
		return nil, err
	}
	m := UnmarshalArrayAction(a)
	return m, nil
}

func (r *ActionRepo) ReadRuns(actionname interface{}) ([]models.Run, error) {
	var a []Run
	err := r.getDB().Model(&Run{}).Joins("JOIN actions on runs.action_id = actions.id").Where("actions.name = ?", actionname.(string)).Find(&a).Error
	if err != nil {
		return nil, err
	}
	return UnmarshalArrayRun(a), nil
}

func (r *ActionRepo) ReadAll() ([]models.Action, error) {
	var a []Action
	err := r.db.Model(&Action{}).Preload("HostsRef").Preload("TriggersRef").Preload("TagsRef").Preload("GroupsRef").Find(&a).Error
	if err != nil {
		return nil, err
	}
	m := UnmarshalArrayAction(a)
	return m, nil
}

func (r *ActionRepo) Update(name string, data map[string]interface{}) (*models.Action, error) {
	actionID, err := r.GetID(name)
	if err != nil {
		return nil, err
	}
	if val, ok := data["hosts"].([]interface{}); ok {
		for _, service := range val {
			ms := models.Service{}
			ms.FromJson(service.(map[string]interface{}))
			id, err := r.GetHostID(ms.HostName)
			if err != nil {
				return nil, err
			}
			ah := ActionsHosts{HostID: id, ActionID: actionID, Port: ms.Port}
			err = r.getDB().Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]interface{}{"port": ms.Port}),
			}).Model(&ActionsHosts{}).Create(&ah).Error
			if err != nil {
				return nil, err
			}
		}
	}
	if val, ok := data["triggers"]; ok {
		switch val.(type) {
		case []interface{}:
			break
		default:
			return nil, errors.New("triggers not an string array")
		}
		actionNames := []string{}
		for _, in := range val.([]interface{}) {
			actionNames = append(actionNames, in.(string))
		}

		actions := []Action{}
		for _, actionName := range actionNames {
			triggerID, err := r.GetID(actionName)
			if err != nil {
				return nil, err
			}
			actions = append(actions, Action{DBModel: DBModel{ID: triggerID}})
		}
		err = r.getDB().Model(&Action{DBModel: DBModel{ID: actionID}}).Omit("TriggersRef.*").Association("TriggersRef").Replace(&actions)
		if err != nil {
			return nil, err
		}
	}
	if val, ok := data["tags"]; ok {
		switch val.(type) {
		case []interface{}:
			break
		default:
			return nil, errors.New("tags not an string array")
		}
		tagNames := []string{}
		for _, in := range val.([]interface{}) {
			tagNames = append(tagNames, in.(string))
		}

		tags := []Tag{}
		for _, tagName := range tagNames {
			tagID, err := r.GetTagID(tagName)
			if err != nil {
				return nil, err
			}
			tags = append(tags, Tag{DBModel: DBModel{ID: tagID}})
		}
		err = r.getDB().Model(&Action{DBModel: DBModel{ID: actionID}}).Omit("TagsRef.*").Association("TagsRef").Replace(&tags)
		if err != nil {
			return nil, err
		}
	}
	if val, ok := data["groups"]; ok {
		switch val.(type) {
		case []interface{}:
			break
		default:
			return nil, errors.New("groups not an string array")
		}
		groupNames := []string{}
		for _, in := range val.([]interface{}) {
			groupNames = append(groupNames, in.(string))
		}

		groups := []Group{}
		for _, groupName := range groupNames {
			groupID, err := r.GetTagID(groupName)
			if err != nil {
				return nil, err
			}
			groups = append(groups, Group{DBModel: DBModel{ID: groupID}})
		}
		err = r.getDB().Model(&Action{DBModel: DBModel{ID: actionID}}).Omit("GroupsRef.*").Association("GroupsRef").Replace(&groups)
		if err != nil {
			return nil, err
		}
	}

	if val, ok := data["variables"]; ok {
		b, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}
		data["variables"] = b
	}
	delete(data, "hosts")
	delete(data, "triggers")
	delete(data, "tags")
	delete(data, "groups")

	err = r.db.Model(&Action{DBModel: DBModel{ID: actionID}}).Updates(data).Error
	if err != nil {
		return nil, err
	}
	new, err := r.ReadByName(name)
	if err != nil {
		return nil, err
	}
	return new, nil
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
	err := json.Unmarshal(a.Variables.MarshalJSON())
	if err != nil {
		fmt.Println(err.Error())
	}
	return models.Action{
		Name:     a.Name,
		Script:   a.Script,
		Triggers: a.Triggers,
		// Hosts:     services,
		Tags:      a.Tags,
		Variables: v,
	}
}

func ReduceExtended(m *models.ActionExtendedv3) *models.Action {
	triggers := []string{}
	for _, t := range m.Triggers {
		triggers = append(triggers, t.Name)
	}
	services := []models.Service{}
	for _, c := range m.Hosts {
		services = append(services, models.Service{HostName: c.Name, Port: c.Port})
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

func UnmarshalActionExtendedv3(a Action) models.ActionExtendedv3 {
	v := map[string]interface{}{}
	if a.Variables.String() != "null" {
		err := json.Unmarshal(a.Variables.MarshalJSON())
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return models.ActionExtendedv3{
		Name:      a.Name,
		Script:    a.Script,
		Groups:    a.Groups,
		Triggers:  UnmarshalArrayActionExtendedv3(a.TriggersRef),
		Hosts:     a.Connections,
		Variables: v,
	}
}

func UnmarshalArrayActionExtendedv3(a []Action) []models.ActionExtendedv3 {
	actions := []models.ActionExtendedv3{}
	for _, m := range a {
		actions = append(actions, UnmarshalActionExtendedv3(m))
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
