package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type MySQL struct {
	db       *gorm.DB
	endpoint string
	logLevel int
}

func (m *MySQL) New() error {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,                 // Slow SQL threshold
			LogLevel:      logger.LogLevel(m.logLevel), // Log level
			Colorful:      false,                       // Disable color
		},
	)

	var err error
	if m.db, err = gorm.Open(mysql.Open(m.endpoint), &gorm.Config{
		Logger: newLogger,
	}); err != nil {
		return err
	}

	return nil
}
func (m *MySQL) DB() *gorm.DB {
	return m.db
}
func (m *MySQL) SetEndpoint(endpoint string) {
	m.endpoint = endpoint
}
func (m *MySQL) SetLogLevel(level int) {
	m.logLevel = level
}
func (m *MySQL) SetMaxIdleConnections(n int) {
	db, err := m.db.DB()
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(n)
}

func (m *MySQL) SetMaxOpenConnections(n int) {
	db, err := m.db.DB()
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(n)
}
func (m *MySQL) GetKind() int {
	return typeMySQL
}
