package internal

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/sirupsen/logrus"
)

func createConnection(connectionString string) *sql.DB {
	logrus.Debug("Create new connection")
	dbConn, connectionError := sql.Open("mssql", connectionString)
	if connectionError != nil {
		logrus.Errorf("error opening database: %v", connectionError)
	}
	return dbConn
}

func (c collector) checkConnection() bool {
	err := c.databaseConnection.PingContext(c.ctx)
	if err != nil {
		logrus.Debug("Connection Ping failed")
		return false
	}
	logrus.Debug("Connection Ping succeeded")
	return true
}

func performQuery(query string, conn *sql.DB) *sql.Rows {
	rows, err := conn.Query(query)
	if err != nil {
		logrus.Errorf("query %s failed with error: %s", query, err)
		return &sql.Rows{}
	}
	return rows
}
