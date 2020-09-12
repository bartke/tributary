package gormfilter

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"time"

	"github.com/bartke/tributary"
	"github.com/bartke/tributary/pipeline/interceptor"
	"github.com/bartke/tributary/sink/handler"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Record struct {
	Hash       string `gorm:"type:varchar(255);unique"`
	CreateTime int64
}

type Filter struct {
	db      *gorm.DB
	builder tributary.EventConstructor
}

func New(dbhandle gorm.Dialector, cfg *gorm.Config, builder tributary.EventConstructor) (*Filter, error) {
	cfg.Logger = logger.Default.LogMode(logger.Silent)
	db, err := gorm.Open(dbhandle, cfg)
	if err != nil {
		return nil, err
	}
	return &Filter{
		db:      db,
		builder: builder,
	}, nil
}

func (f *Filter) Create(name string) (interceptor.Fn, error) {
	err := f.db.Table(name).AutoMigrate(&Record{})
	if err != nil {
		return nil, err
	}
	return func(msg tributary.Event) tributary.Event {
		hash := md5.Sum(msg.Payload())
		result := f.db.Table(name).Create(&Record{Hash: hex.EncodeToString(hash[:]), CreateTime: time.Now().UnixNano()})
		if result.Error != nil {
			return f.builder(msg.Context(), msg.Payload(), result.Error)
		}
		return f.builder(msg.Context(), msg.Payload(), nil)
	}, nil
}

func (f *Filter) Clean(name string, s int) handler.Fn {
	return func(e tributary.Event) {
		result := f.db.Table(name).Delete(Record{}, "create_time < ?", time.Now().Add(-1*time.Duration(s)*time.Second).UnixNano())
		if result.Error != nil {
			log.Println(result.Error)
		}
		//var tick time.Time
		//tick.UnmarshalBinary(e.Payload())
		//fmt.Println(result.RowsAffected, "purged", tick, "after", time.Since(tick))
	}
}
