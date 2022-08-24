package internal

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getWaitStatsStats(conn *sql.DB) []prometheus.Metric {
	metrics := []prometheus.Metric{}

	rows, err := performQuery(`SELECT
	      wait_type,
		  waiting_tasks_count,
		  wait_time_ms,
		  max_wait_time_ms,
		  signal_wait_time_ms
		FROM sys.dm_os_wait_stats where wait_type in ('ASYNC_NETWORK_IO','CXPACKET','DTC','NETWORKIO','OLEDB','SOS_SCHEDULER_YIELD','WRITELOG','IO_COMPLETION','IO_RETRY','WAIT_ON_SYNC_STATISTICS_REFRESH') 
		OR wait_type like 'LCK_M_%' 
		OR wait_type like 'PAGEIOLATCH_%';`,
		conn,
	)
	if err != nil {
		logrus.Errorf("Error in query execution, skipping metrics")
		return metrics
	}
	var wait_type string
	var waiting_tasks_count int64
	var wait_time_ms int64
	var max_wait_time_ms int64
	var signal_wait_time_ms int64
	for rows.Next() {

		if err := rows.Scan(
			&wait_type,
			&waiting_tasks_count,
			&wait_time_ms,
			&max_wait_time_ms,
			&signal_wait_time_ms,
		); err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
			continue
		}

		metrics = append(metrics, returnMetric(
			"sql_wait_stats_waiting_tasks_count",
			"Current value of waiting_tasks_count in dm_os_wait_stats by wait_type",
			"wait_type",
			wait_type,
			float64(waiting_tasks_count),
		))

		metrics = append(metrics, returnMetric(
			"sql_wait_stats_wait_time_ms",
			"Current value of wait_time_ms in dm_os_wait_stats by wait_type",
			"wait_type",
			wait_type,
			float64(wait_time_ms),
		))

		metrics = append(metrics, returnMetric(
			"sql_wait_stats_max_wait_time_ms",
			"Current value of max_wait_time_ms in dm_os_wait_stats by wait_type",
			"wait_type",
			wait_type,
			float64(max_wait_time_ms),
		))

		metrics = append(metrics, returnMetric(
			"sql_wait_stats_signal_wait_time_ms",
			"Current value of signal_wait_time_ms in dm_os_wait_stats by wait_type",
			"wait_type",
			wait_type,
			float64(signal_wait_time_ms),
		))
	}
	err = rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}
