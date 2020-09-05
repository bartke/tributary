package main

import (
	"github.com/bartke/tributary/example/advanced/event"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func NewDB() error {
	inmemory := sqlite.Open("file::memory:?cache=shared")
	var err error
	db, err = gorm.Open(inmemory, &gorm.Config{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&event.BetORM{})
	if err != nil {
		return err
	}
	return nil
}
