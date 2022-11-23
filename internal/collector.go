package internal

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
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

	conn, err := connect(fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;encrypt=true",
		dbConnectionInfo.HostName,
		dbConnectionInfo.User,
		dbConnectionInfo.Password,
		dbConnectionInfo.Port,
		dbConnectionInfo.DBName),
	)
	if err != nil {
		logrus.Fatalf("Failed to connect: %s", err)
	}

	return &collector{
		databaseConnection: conn,
		ctx:                context.Background(),
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
	// 	ch <- c.up.Desc()
}

func (c collector) Collect(ch chan<- prometheus.Metric) {
	c.up.Set(0)

	if c.isConnected() {
		c.up.Set(1)

		t := time.Now()
		metrics := collectMetrics(c.databaseConnection, c.dbname, c.dbhost, c.enabledMetrics)
		logrus.Debugf("Collected %d metrics after %s", len(metrics), time.Since(t))

		numberOfMetrics := len(metrics)
		for i, metric := range metrics {
			tm := time.Now()
			ch <- metric
			logrus.Debugf("%d/%d Added metric %s to the registry after %s", i, numberOfMetrics, metric.Desc().String(), time.Since(tm))
		}

	}
	ch <- c.up
}
