package internal

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/sirupsen/logrus"
)

func createConnection(connectionString string) *sql.DB {
	logrus.Info("Create new connection")
	dbConn, connectionError := sql.Open("mssql", connectionString)
	if connectionError != nil {
		logrus.Errorf("error opening database: %v", connectionError)
	}
	err := dbConn.Ping()
	if err != nil {
		logrus.Error("Initial connection Ping failed")
	}
	logrus.Info("Initial connection Ping succeeded")
	return dbConn
}

func (c collector) checkConnection() bool {
	err := c.databaseConnection.Ping()
	if err != nil {
		logrus.Debug("Connection Ping failed")
		return false
	}
	logrus.Debug("Connection Ping succeeded")
	return true
}

func performQuery(query string, conn *sql.DB) (*sql.Rows, error) {
	rows, err := conn.Query(query)
	if err != nil {
		logrus.Errorf("query %s failed with error: %s", query, err)
		return &sql.Rows{}, err
	}
	return rows, nil
}
