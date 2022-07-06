package internal

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getBufferDescriptorsStats(conn *sql.DB) []prometheus.Metric {
	metrics := []prometheus.Metric{}

	rows, err := performQuery(`SELECT COUNT(*) AS cached_pages_count  
	FROM sys.dm_os_buffer_descriptors  
	WHERE is_in_bpool_extension <> 0;`,
		conn,
	)
	if err != nil {
		logrus.Errorf("Error in query execution, skipping metrics")
		return metrics
	}
	var cached_pages_count int64
	for rows.Next() {

		if err := rows.Scan(
			&cached_pages_count,
		); err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
			continue
		}

		metrics = append(metrics, returnMetric(
			"sql_cached_pages_count",
			"Current value of cached_pages_count in dm_os_buffer_descriptors",
			"none",
			"none",
			float64(cached_pages_count),
		))

	}
	err = rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}
