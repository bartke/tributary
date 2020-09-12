package gormwindow

import (
	"database/sql"
	"encoding/json"
	"log"
	"reflect"

	"github.com/bartke/tributary"
	"github.com/bartke/tributary/pipeline/injector"
	"github.com/bartke/tributary/pipeline/interceptor"
	"gorm.io/gorm"
)

type Window struct {
	db      *gorm.DB
	builder tributary.EventConstructor
}

func New(dbhandle gorm.Dialector, cfg *gorm.Config, builder tributary.EventConstructor, migrations ...interface{}) (*Window, error) {
	db, err := gorm.Open(dbhandle, cfg)
	if err != nil {
		return nil, err
	}
	window := &Window{
		db:      db,
		builder: builder,
	}
	for _, t := range migrations {
		err = db.AutoMigrate(t)
		if err != nil {
			return window, err
		}
	}
	return window, nil
}

func (w *Window) Create(v interface{}) interceptor.Fn {
	return func(msg tributary.Event) tributary.Event {
		// clone zero value from v
		p := reflect.New(reflect.TypeOf(v).Elem()).Interface()
		err := json.Unmarshal(msg.Payload(), &p)
		if err != nil {
			log.Println("inject error", err)
			return w.builder(msg.Context(), msg.Payload(), err)
		}
		result := w.db.Create(p)
		if result.Error != nil {
			log.Println("create error", result.Error)
			return w.builder(msg.Context(), msg.Payload(), result.Error)
		}
		return w.builder(msg.Context(), msg.Payload(), nil)
	}
}

func (w *Window) Query(query string) injector.Fn {
	return func(msg tributary.Event) []tributary.Event {
		records := []map[string]interface{}{}
		result := w.db.Raw(query).Scan(&records)
		if result.Error != nil {
			return []tributary.Event{w.builder(msg.Context(), msg.Payload(), result.Error)}
		}
		if len(records) == 0 && result.RowsAffected == 0 {
			return []tributary.Event{w.builder(msg.Context(), msg.Payload(), sql.ErrNoRows)}
		}
		var out []tributary.Event
		for i := range records {
			r := map[string]interface{}{}
			for k, v := range records[i] {
				b, ok := v.([]byte)
				if ok {
					r[k] = string(b)
				} else {
					r[k] = v
				}
			}
			oi, _ := json.Marshal(r)
			out = append(out, w.builder(msg.Context(), oi, nil))
		}
		return out
	}
}
