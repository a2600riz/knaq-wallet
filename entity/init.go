package entity

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type CommonTime struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

var db *gorm.DB

func SetDB(d *gorm.DB) {
	db = d
	if err := db.AutoMigrate(
		&Wallet{},
		&FtTransaction{},
	); err != nil {
		fmt.Printf("Error occurred while migrating entity: %s\n", err)
	}
}

type Common struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`
	CommonTime

	tx *gorm.DB
}

type DatabaseTransaction interface {
	GetTx() *gorm.DB
	InitTx() error
}
