package internal

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getBufferDescriptorsStats(conn *sql.DB) []prometheus.Metric {
	var metrics []prometheus.Metric

	rows := performQuery(`SELECT COUNT(*) AS cached_pages_count  
	FROM sys.dm_os_buffer_descriptors  
	WHERE is_in_bpool_extension <> 0;`,
		conn,
	)

	for rows.Next() {

		var cached_pages_count int64

		err := rows.Scan(
			&cached_pages_count,
		)
		if err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
		}

		metrics = append(metrics, returnMetric(
			"sql_cached_pages_count",
			"Current value of cached_pages_count in dm_os_buffer_descriptors",
			"none",
			"none",
			float64(cached_pages_count),
		))

	}
	err := rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}
