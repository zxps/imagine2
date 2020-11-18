package storage

import (
	"imagine2/config"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// DB - database connection
var DB *sqlx.DB

// Initialize ...
func Initialize() {
	connection, err := sqlx.Connect("mysql", config.Context.Service.DatabaseDSN)
	if err != nil {
		logrus.Fatal("unable connect to storage ", err.Error())
	}

	DB = connection
}

func getFilesTable() string {
	return config.Context.Service.DatabaseFilesTable
}

func getDialect() string {
	return config.Context.Service.DatabaseDialect
}
