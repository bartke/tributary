package gormdedupe

import (
	"time"

	"github.com/bartke/tributary"
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

type Deduper struct {
	db   *gorm.DB
	name string
}

func New(dbhandle gorm.Dialector, cfg *gorm.Config) (*Filter, error) {
	cfg.Logger = logger.Default.LogMode(logger.Silent)
	db, err := gorm.Open(dbhandle, cfg)
	if err != nil {
		return nil, err
	}
	return &Filter{db: db}, nil
}

func (f *Filter) Create(name string) (*Deduper, error) {
	d := &Deduper{
		db:   f.db,
		name: name,
	}
	err := f.db.Table(name).AutoMigrate(&Record{})
	if err != nil {
		return d, err
	}
	return d, nil
}

func (d *Deduper) Dedupe(msg tributary.Event) (tributary.Event, error) {
	result := d.db.Table(d.name).Create(&Record{Payload: msg.Payload()})
	if result.Error != nil {
		return nil, result.Error
	}
	return msg, nil
}
