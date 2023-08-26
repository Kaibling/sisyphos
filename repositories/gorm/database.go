package gormrepo

import (
	"fmt"
	"time"

	slog "sisyphos/lib/log"
	"sisyphos/lib/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	User     string
	Port     string
	Password string
	Host     string
	Database string
	Dialect  string
}

func InitDatabase(cfg DBConfig, l *slog.Logger) (*gorm.DB, error) {
	newLogger := logger.New(
		l,
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	// var dialector  gorm.Dialector
	// dialector = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database))
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=europe/Amsterdam", cfg.Host, cfg.User, cfg.Password, cfg.Database, cfg.Port)
	dialector := postgres.Open(dsn)
	db, err := gorm.Open(dialector, &gorm.Config{Logger: newLogger})
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
		NewLogDBMigrator(db),
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
	ID        string    `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"created_at"`
	CreatedBy string    `gorm:"created_by"`
	UpdatedAt time.Time `gorm:"updated_at"`
	UpdatedBy string    `gorm:"updated_by"`
}

func (db *DBModel) BeforeCreate(tx *gorm.DB) error {
	if db.ID == "" {
		id := utils.NewULID()
		db.ID = id.String()
	}
	return nil
}
