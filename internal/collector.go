package internal

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type collector struct {
	databaseConnection *sql.DB
	ctx                context.Context
	up                 prometheus.Gauge
	dbname             string
	dbhost             string
	enabledMetrics     []string
}

func NewCollector(dbConnectionInfo DBConnectionInfo, enabledMetrics []string) *collector {

	return &collector{
		databaseConnection: createConnection(
			fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s",
				dbConnectionInfo.HostName,
				dbConnectionInfo.User,
				dbConnectionInfo.Password,
				dbConnectionInfo.Port,
				dbConnectionInfo.DBName),
		),
		ctx: context.Background(),
		up: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "up",
				Help: "Was the last query of mssql stats successful.",
			}),
		dbname:         dbConnectionInfo.DBName,
		dbhost:         dbConnectionInfo.HostName,
		enabledMetrics: enabledMetrics,
	}
}

func (c collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up.Desc()
}

func (c collector) Collect(ch chan<- prometheus.Metric) {
	if c.checkConnection() == false {
		ch <- prometheus.MustNewConstMetric(c.up.Desc(), prometheus.GaugeValue, 0)
		return
	}
	ch <- prometheus.MustNewConstMetric(c.up.Desc(), prometheus.GaugeValue, 1)
	for _, metricsReturned := range returnMetrics(
		c.databaseConnection,
		c.dbname,
		c.dbhost,
		c.enabledMetrics,
	) {
		ch <- metricsReturned
	}
}
