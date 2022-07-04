package internal

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getLatchStats(conn *sql.DB) []prometheus.Metric {
	var metrics []prometheus.Metric

	rows := performQuery(`select 
	  latch_class,
	  waiting_requests_count,
	  wait_time_ms,
	  max_wait_time_ms 
	FROM sys.dm_os_latch_stats;`,
		conn,
	)

	for rows.Next() {

		var latch_class string
		var waiting_requests_count int
		var wait_time_ms int
		var max_wait_time_ms int

		err := rows.Scan(
			&latch_class,
			&waiting_requests_count,
			&wait_time_ms,
			&max_wait_time_ms,
		)
		if err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
		}

		metrics = append(metrics, returnMetric(
			"sql_latch_waiting_requests_count",
			"Current value of waiting_requests_count in dm_os_latch_stats",
			"latch_class",
			latch_class,
			float64(waiting_requests_count),
		))

		metrics = append(metrics, returnMetric(
			"sql_latch_wait_time_ms",
			"Current value of wait_time_ms in dm_os_latch_stats",
			"latch_class",
			latch_class,
			float64(wait_time_ms),
		))

		metrics = append(metrics, returnMetric(
			"sql_latch_max_wait_time_ms",
			"Current value of max_wait_time_ms in dm_os_latch_stats",
			"latch_class",
			latch_class,
			float64(max_wait_time_ms),
		))

	}
	err := rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}
