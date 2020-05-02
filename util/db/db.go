package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB

	ErrMissingTx = fmt.Errorf("missing transaction in context")
)

func newDB(connectionInfo string) error {
	newDB, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return err
	}
	newDB.SetNowFuncOverride(func() time.Time { return time.Now().UTC() })

	db = newDB
	return nil
}

func Get() *gorm.DB {
	return db
}

func Start() {
	err := newDB(os.Getenv("PG_DATASOURCE"))
	if err != nil {
		log.Fatal("Error connect db :", err)
	}
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func Healthz() error {
	return db.DB().Ping()
}
