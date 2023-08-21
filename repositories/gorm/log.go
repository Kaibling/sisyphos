package gormrepo

import (
	"sisyphos/models"
	"time"

	"gorm.io/gorm"
)

type Log struct {
	Url       string
	Body      string
	Method    string
	User      string
	RequestID string    `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"created_at"`
}

type LogRepo struct {
	db       *gorm.DB
	username string
}

func NewLogRepo(db *gorm.DB, username string) *LogRepo {
	return &LogRepo{db, username}
}

func (r *LogRepo) getDB() *gorm.DB {
	d := r.db.Session(&gorm.Session{NewDB: true})
	return d
}

func (r *LogRepo) Create(logs []models.Log) ([]models.Log, error) {
	resp := []models.Log{}
	for _, a := range logs {
		log := marshalLog(a)
		log.CreatedAt = time.Now()
		err := r.getDB().Create(&log).Error
		if err != nil {
			return nil, err
		}

		newLog, err := r.ReadByRequestID(log.RequestID)
		if err != nil {
			return nil, err
		}
		resp = append(resp, *newLog)
	}
	return resp, nil
}

func (r *LogRepo) ReadByRequestID(request_id interface{}) (*models.Log, error) {
	var a Log
	err := r.db.Model(&Log{}).Where(&Log{RequestID: request_id.(string)}).First(&a).Error
	if err != nil {
		return nil, err
	}
	m := unmarshalLog(a)
	return &m, nil
}

func (r *LogRepo) ReadAll() ([]models.Log, error) {
	var a []Log
	err := r.db.Model(&Log{}).Find(&a).Error
	if err != nil {
		return nil, err
	}
	m := unmarshalArrayLog(a)
	return m, nil
}

func marshalLog(l models.Log) Log {
	return Log{
		Url:       l.Url,
		Body:      l.Body,
		Method:    l.Method,
		User:      l.User,
		RequestID: l.RequestID,
	}
}

func unmarshalLog(l Log) models.Log {
	return models.Log{
		Url:       l.Url,
		Body:      l.Body,
		Method:    l.Method,
		User:      l.User,
		RequestID: l.RequestID,
	}
}

// func marshalArrayLog(m []models.Log) []Log {
// 	logs := []Log{}
// 	for _, a := range m {
// 		logs = append(logs, marshalLog(a))
// 	}
// 	return logs
// }

func unmarshalArrayLog(a []Log) []models.Log {
	logs := []models.Log{}
	for _, m := range a {
		logs = append(logs, unmarshalLog(m))
	}
	return logs
}

type LogDBMigrator struct {
	db *gorm.DB
}

func (s LogDBMigrator) Migrate() error {
	err := s.db.AutoMigrate(&Log{})
	if err != nil {
		return err
	}
	return nil
}

func NewLogDBMigrator(db *gorm.DB) *LogDBMigrator {
	return &LogDBMigrator{db: db}
}
