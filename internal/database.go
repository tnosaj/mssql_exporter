package internal

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/sirupsen/logrus"
)

func connect(connectionString string) (*sql.DB, error) {
	logrus.Info("Create new connection")
	dbConn, err := sql.Open("mssql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %s", err)
	}
	if err := dbConn.Ping(); err != nil {
		return nil, fmt.Errorf("initial database Ping failed: %s", err)
	}

	logrus.Info("Connected to the MSSQL database, and initial Ping was successful")
	return dbConn, nil
}

func (c collector) isConnected() bool {
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
