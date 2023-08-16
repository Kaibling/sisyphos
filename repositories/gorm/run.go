package gormrepo

import (
	"time"

	"sisyphos/lib/utils"
	"sisyphos/models"

	"gorm.io/gorm"
)

type Run struct {
	DBModel
	RunID     string
	RequestID string
	ParentID  string
	User      string
	HostID    *string `gorm:"size:255"`
	Host      Host
	ActionID  string `gorm:"size:255"`
	Action    Action
	StartTime time.Time
	EndTime   time.Time
	Duration  string
	Output    string
	Error     string
	Status    string
}

type RunRepo struct {
	db        *gorm.DB
	username  string
	requestid string
}

func NewRunRepo(db *gorm.DB, requestID, username string) *RunRepo {
	return &RunRepo{db: db, username: username, requestid: requestID}
}

func (r *RunRepo) getDB() *gorm.DB {
	d := r.db.Session(&gorm.Session{NewDB: true})
	return d
}

func (r *RunRepo) Create(runs []models.Run) ([]models.Run, error) {
	resp := []models.Run{}
	for _, a := range runs {
		run, err := MarshalRun(a, r.getDB())
		if err != nil {
			return nil, err
		}
		err = r.getDB().Omit("ActionRef.*").Omit("HostRef").Create(&run).Error
		if err != nil {
			return nil, err
		}
		newRun, err := r.ReadByRunID(run.RunID)
		if err != nil {
			return nil, err
		}
		resp = append(resp, *newRun)
	}
	return resp, nil
}

func (r *RunRepo) ReadByRunID(runid interface{}) (*models.Run, error) {
	var a Run
	err := r.db.Model(&Run{}).Where(&Run{RunID: runid.(string)}).Preload("Action").Preload("Host").First(&a).Error
	if err != nil {
		return nil, err
	}
	m := UnmarshalRun(a)
	return &m, nil
}

func (r *RunRepo) ReadByReqID(reqID interface{}) ([]models.Run, error) {
	var a []Run
	err := r.getDB().Model(&Run{}).Where(&Run{RequestID: reqID.(string)}).Preload("Action").Preload("Host").Find(&a).Error
	if err != nil {
		return nil, err
	}
	return UnmarshalArrayRun(a), nil
}

func (r *RunRepo) GetRequestID() string {
	return r.requestid
}

func (r *RunRepo) GetUsername() string {
	return r.username
}

func (r *RunRepo) ReadAll() ([]models.Run, error) {
	var a []Run
	err := r.db.Model(&Run{}).Preload("HostRef").Preload("ActionRef").Find(&a).Error
	if err != nil {
		return nil, err
	}
	m := UnmarshalArrayRun(a)
	return m, nil
}

func MarshalRun(a models.Run, db *gorm.DB) (*Run, error) {
	actionRepo := NewActionRepo(db)
	actionID, err := actionRepo.GetID(a.Action)
	if err != nil {
		return nil, err
	}
	hostRepo := NewHostRepo(db)
	hID, err := hostRepo.GetID(utils.PtrRead(a.Host))
	var hostID *string
	if err == nil {
		hostID = &hID
	}
	return &Run{
		RunID:     a.RunID,
		RequestID: a.RequestID,
		ParentID:  a.ParentID,
		HostID:    hostID,
		User:      a.User,
		ActionID:  actionID,
		StartTime: a.StartTime,
		EndTime:   a.EndTime,
		Duration:  a.Duration,
		Output:    a.Output,
		Error:     a.Error,
		Status:    a.Status,
	}, nil
}

func UnmarshalRun(a Run) models.Run {
	return models.Run{
		RunID:     a.RunID,
		RequestID: a.RequestID,
		ParentID:  a.ParentID,
		User:      a.User,
		Host:      a.Host.Name,
		Action:    *a.Action.Name,
		StartTime: a.StartTime,
		EndTime:   a.EndTime,
		Duration:  a.Duration,
		Output:    a.Output,
		Error:     a.Error,
		Status:    a.Status,
	}
}

func MarshalArrayRun(m []models.Run, db *gorm.DB) ([]Run, error) {
	runs := []Run{}
	for _, a := range m {
		mr, err := MarshalRun(a, db)
		if err != nil {
			return nil, err
		}
		runs = append(runs, *mr)
	}
	return runs, nil
}

func UnmarshalArrayRun(a []Run) []models.Run {
	runs := []models.Run{}
	for _, m := range a {
		runs = append(runs, UnmarshalRun(m))
	}
	return runs
}

type RunDBMigrator struct {
	db *gorm.DB
}

func (s RunDBMigrator) Migrate() error {
	err := s.db.AutoMigrate(&Run{})
	if err != nil {
		return err
	}
	return nil
}

func NewRunDBMigrator(db *gorm.DB) *RunDBMigrator {
	return &RunDBMigrator{db: db}
}
