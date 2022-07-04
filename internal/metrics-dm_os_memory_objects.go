package internal

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getMemoryObjectsStats(conn *sql.DB) []prometheus.Metric {
	var metrics []prometheus.Metric

	rows := performQuery(`SELECT type, SUM(pages_in_bytes)   
	FROM sys.dm_os_memory_objects  
	GROUP BY type;`,
		conn,
	)

	for rows.Next() {

		var ttype string
		var sum_pages_in_bytes int64

		err := rows.Scan(
			&ttype,
			&sum_pages_in_bytes,
		)
		if err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
		}

		metrics = append(metrics, returnMetric(
			"sql_memory_clerks_sum_pages_in_bytes",
			"Current value of sum_pages_in_bytes in dm_os_memory_objects for the type",
			"type",
			ttype,
			float64(sum_pages_in_bytes),
		))

	}
	err := rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}
