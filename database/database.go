package database

import (
	"gorm.io/gorm"
)

type Databases struct {
	list []Database
}

func (d *Databases) GetDB(kind int) Database {
	for _, db := range d.list {
		if db.GetKind() == kind {
			return db
		}
	}
	return nil
}

type Database interface {
	New() error
	DB() *gorm.DB
	SetEndpoint(endpoint string)
	SetLogLevel(level int)
	SetMaxIdleConnections(n int)
	SetMaxOpenConnections(n int)
	GetKind() int
}
