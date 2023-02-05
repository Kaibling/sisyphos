package gormrepo

import (
	"log"
	"os"
	"time"

	"sisyphos/lib/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase() (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	dsn := "db:example@tcp(db:3306)/db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, err
	}

	dbMirgators := []DBMigrator{
		NewTagDBMigrator(db),
		NewActionMigrator(db),
		NewTokenDBMigrator(db),
		NewHostDBMigrator(db),
		NewUserDBMigrator(db),
		NewGroupDBMigrator(db),
		NewRunDBMigrator(db),
	}
	err = Migrate(dbMirgators)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(dbMigs []DBMigrator) error {
	for i := range dbMigs {
		err := dbMigs[i].Migrate()
		if err != nil {
			return err
		}
	}
	return nil
}

type DBMigrator interface {
	Migrate() error
}
type DBModel struct {
	ID string `gorm:"primaryKey"`
}

func (db *DBModel) BeforeCreate(tx *gorm.DB) error {
	if db.ID == "" {
		id := utils.NewULID()
		db.ID = id.String()
	}
	return nil
}
