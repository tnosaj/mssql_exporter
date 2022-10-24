package internal

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getSettings(conn *sql.DB) []prometheus.Metric {
	metrics := []prometheus.Metric{}

	rows, err := performQuery(`Select name, is_auto_update_stats_async_on from sys.databases where name !='master';`,
		conn,
	)
	if err != nil {
		logrus.Errorf("Error in query execution, skipping metrics")
		return metrics
	}
	var dbname string
	var enabled bool
	for rows.Next() {

		if err := rows.Scan(
			&dbname,
			&enabled,
		); err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
			continue
		}

		intenabled := 0
		if enabled {
			intenabled = 1
		}
		metrics = append(metrics, returnMetric(
			"sql_settings_async_stats_update",
			"Current value of is_auto_update_stats_async_on for the database",
			"database",
			dbname,
			float64(intenabled),
		))

	}
	err = rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}
