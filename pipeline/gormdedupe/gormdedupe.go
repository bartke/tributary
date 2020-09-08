package gormdedupe

import (
	"log"
	"time"

	"github.com/bartke/tributary"
	"github.com/bartke/tributary/pipeline/interceptor"
	"github.com/bartke/tributary/sink/handler"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Record struct {
	Id         int       `gorm:"PrimaryKey;autoincrement:true"`
	Payload    []byte    `gorm:"unique"`
	CreateTime time.Time `gorm:"DEFAULT:current_timestamp"`
}

type Filter struct {
	db *gorm.DB
}

func New(dbhandle gorm.Dialector, cfg *gorm.Config) (*Filter, error) {
	cfg.Logger = logger.Default.LogMode(logger.Silent)
	db, err := gorm.Open(dbhandle, cfg)
	if err != nil {
		return nil, err
	}
	return &Filter{db: db}, nil
}

func (f *Filter) Create(name string) (interceptor.Fn, error) {
	err := f.db.Table(name).AutoMigrate(&Record{})
	if err != nil {
		return nil, err
	}
	return func(msg tributary.Event) (tributary.Event, error) {
		result := f.db.Table(name).Create(&Record{Payload: msg.Payload()})
		if result.Error != nil {
			return nil, result.Error
		}
		return msg, nil
	}, nil
}

func (f *Filter) Clean(name string, s int) handler.Fn {
	return func(e tributary.Event) {
		result := f.db.Table(name).Delete(Record{}, "create_time < ?", time.Now().Add(-1*time.Duration(s)*time.Second))
		if result.Error != nil {
			log.Println(result.Error)
		}
		//var tick time.Time
		//tick.UnmarshalBinary(e.Payload())
		//fmt.Println(result.RowsAffected, "purged", tick, "after", time.Since(tick))
	}
}
