package gormdeduper

import (
	"time"

	"github.com/bartke/tributary"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Record struct {
	Id         int       `gorm:"PrimaryKey;autoincrement:true"`
	Payload    []byte    `gorm:"index:91efe54f-03f9-4033-97a6-4bc98051d1eb_idx,unique"`
	CreateTime time.Time `gorm:"DEFAULT:current_timestamp"`
}

func (Record) TableName() string {
	return "91efe54f-03f9-4033-97a6-4bc98051d1eb"
}

type Deduper struct {
	db *gorm.DB
}

func New(dbhandle gorm.Dialector, cfg *gorm.Config) (*Deduper, error) {
	cfg.Logger = logger.Default.LogMode(logger.Silent)
	db, err := gorm.Open(dbhandle, cfg)
	if err != nil {
		return nil, err
	}
	deduper := &Deduper{
		db: db,
	}
	err = db.AutoMigrate(&Record{})
	if err != nil {
		return deduper, err
	}
	return deduper, nil
}

func (d *Deduper) Dedupe(msg tributary.Event) (tributary.Event, error) {
	result := d.db.Create(&Record{Payload: msg.Payload()})
	if result.Error != nil {
		return nil, result.Error
	}
	return msg, nil
}
