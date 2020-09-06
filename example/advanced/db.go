package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/bartke/tributary"
	"github.com/bartke/tributary/example/advanced/event"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Window struct {
	db *gorm.DB
}

func NewDB() (*Window, error) {
	//inmemory := sqlite.Open("file::memory:?cache=shared")
	//db, err := gorm.Open(inmemory, &gorm.Config{})
	inmemory := sqlite.Open("file:test.db")
	db, err := gorm.Open(inmemory, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	window := &Window{db}
	err = db.AutoMigrate(&event.Bet{})
	if err != nil {
		return window, err
	}
	err = db.AutoMigrate(&event.Selection{})
	if err != nil {
		return window, err
	}
	return window, nil
}

func (w *Window) createInjector(msg tributary.Event) (tributary.Event, error) {
	bet := &event.Bet{}
	err := json.Unmarshal(msg.Payload(), &bet)
	if err != nil {
		log.Println("inject error", err)
		return nil, err
	}
	result := w.db.Create(bet)
	if result.Error != nil {
		log.Println("create error", result.Error)
		return nil, err
	}
	return msg, result.Error
}

func (w *Window) queryWindow(query string) func(msg tributary.Event) (tributary.Event, error) {
	return func(msg tributary.Event) (tributary.Event, error) {
		ress := []map[string]interface{}{}
		result := w.db.Raw(query).Scan(&ress)
		if result.Error != nil {
			log.Println("query error", result.Error)
			return nil, result.Error
		}
		if len(ress) == 0 {
			return nil, errors.New("no return")
		}
		out, _ := json.Marshal(ress)
		return Msg(out, msg.Context()), result.Error
	}
}
