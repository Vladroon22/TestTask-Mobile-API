package database

import (
	"database/sql"
	"time"

	golog "github.com/Vladroon22/GoLog"
	"github.com/Vladroon22/TestTask-Mobile-API/config"
	_ "github.com/lib/pq"
)

type DataBase struct {
	logger *golog.Logger
	config *config.Config
	sqlDB  *sql.DB
}

func NewDB(conf *config.Config, logg *golog.Logger) *DataBase {
	return &DataBase{
		config: conf,
		logger: logg,
	}
}

func (d *DataBase) Connect() error {
	if err := d.openDB(*d.config); err != nil {
		return err
	}
	return nil
}

func (d *DataBase) openDB(conf config.Config) error {
	str := "postgresql://" + conf.DB
	db, err := sql.Open("postgres", str)
	d.logger.Infoln(str)
	if err != nil {
		return err
	}
	if err := RetryPing(db); err != nil {
		return err
	}
	d.sqlDB = db

	return nil
}

func RetryPing(db *sql.DB) error {
	var err error
	for i := 0; i < 5; i++ {
		if err = db.Ping(); err == nil {
			return nil
		}
		time.Sleep(time.Second)
	}
	return err
}

func (db *DataBase) CloseDB() error {
	return db.sqlDB.Close()
}
