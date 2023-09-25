package database

import (
	"knaq-wallet/config"
	"knaq-wallet/entity"
)

const (
	typeMySQL = iota + 1
)

var DBs Databases

func New() {
	DBs = startDatabase()
	entity.SetDB(DBs.GetDB(typeMySQL).DB())
}

func startDatabase() Databases {
	var databases Databases
	for _, database := range config.Config.GetDatabases() {
		switch database.GetKind() {
		case typeMySQL:
			my := MySQL{}
			setEndpoint(&my, database.GetEndpoint())
			setLogLevel(&my, database.GetLogLevel())
			databases.list = append(databases.list, getDatabase(&my, database.GetMaxIdleConns(), database.GetMaxOpenConns()))
		}
	}

	return databases
}

func setEndpoint(db Database, endpoint string) {
	db.SetEndpoint(endpoint)
}

func setLogLevel(db Database, level int) {
	db.SetLogLevel(level)
}

func setMaxIdleConnections(db Database, n int) {
	db.SetMaxIdleConnections(n)
}

func setMaxOpenConnections(db Database, n int) {
	db.SetMaxOpenConnections(n)
}

func getDatabase(db Database, maxIdleConns, maxOpenConns int) Database {
	if err := db.New(); err != nil {
		panic(err)
	}
	setMaxIdleConnections(db, maxIdleConns)
	setMaxOpenConnections(db, maxOpenConns)
	return db
}
