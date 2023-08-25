package mysql

import (
	"fmt"
	"time"

	"sisyphos/lib/cluster/repos"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const tableName = "cluster_lock"
const dbName = "cluster_lock"

type MysqlBackend struct {
	cfg MysqlConfig
	db  *sqlx.DB
	log repos.Logger
}

func New(cfg MysqlConfig) *MysqlBackend {
	cfg.Defaults()
	return &MysqlBackend{cfg: cfg}
}

type MysqlConfig struct {
	User      string
	Password  string
	Host      string
	Port      string
	Database  string
	Tablename string
}

func (mcfg *MysqlConfig) Defaults() {
	if mcfg.Tablename == "" {
		mcfg.Tablename = tableName
	}
	if mcfg.Database == "" {
		mcfg.Database = dbName
	}
}

func (be *MysqlBackend) Lock(key string, lockDuration time.Duration) (bool, error) {
	query := fmt.Sprintf(`
	UPDATE %s SET 
		master_id = :master_id,
		last_heartbeat = :last_heartbeat
	WHERE 
		last_heartbeat < :min_time
		or
		master_id = :master_id`,
		be.cfg.Tablename)
	ud := map[string]any{
		"last_heartbeat": time.Now(),
		"min_time":       time.Now().Add(-lockDuration),
		"master_id":      key,
	}

	if res, err := be.db.NamedExec(query, ud); err != nil {
		return false, fmt.Errorf("failed to obtain a master lock: %v", err)
	} else {

		rAff, err := res.RowsAffected()
		if err != nil {
			return false, err
		}
		if rAff == 0 {
			return false, nil
		}
	}

	return true, nil
}
func (be *MysqlBackend) UnLock(key string) error {
	return nil
}

func (be *MysqlBackend) Connect() error {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", be.cfg.User, be.cfg.Password, be.cfg.Host, be.cfg.Port, be.cfg.Database))
	if err != nil {
		db.Close()
		return err
	}
	be.db = db
	be.log.Debugf("connected to DB")

	if err := be.createTable(); err != nil {
		be.log.Errorf("table creation failed: %s", err.Error())
	}
	return nil
}

func (be *MysqlBackend) createTable() error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (`+
		`id VARCHAR(255) NOT NULL,`+
		`master_id VARCHAR(255) UNIQUE,`+
		`last_heartbeat TIMESTAMP(3) NULL,`+
		`PRIMARY KEY(id));`, be.cfg.Tablename)
	be.log.Debugf("Attempting to create lock table")

	_, err := be.db.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("unable to create lock table: %v", err)
	}
	be.log.Debugf("ensured table '%v' exists", be.cfg.Tablename)

	return nil
}

func (be *MysqlBackend) AddEmptyLock(key string) error {
	query := fmt.Sprintf(`
	INSERT IGNORE  INTO %s (id, master_id, last_heartbeat) 
	VALUES (:id, :master_id, :last_heartbeat)`,
		be.cfg.Tablename)
	ud := map[string]any{
		"id":             "masterlock",
		"last_heartbeat": time.Now(),
		"master_id":      key,
	}
	if _, err := be.db.NamedExec(query, ud); err != nil {
		return fmt.Errorf("failed to obtain a master lock: %v", err)
	}
	return nil
}
