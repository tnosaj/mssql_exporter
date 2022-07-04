package internal

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getTasksStats(conn *sql.DB) []prometheus.Metric {
	var metrics []prometheus.Metric

	rows := performQuery(`SELECT
		task_state, 
		sum(context_switches_count) as sum_context_switches,
		sum(pending_io_count) as sum_io_operation_count,
		sum(pending_io_byte_count) as sum_io_bytes_count, 
		sum(pending_io_byte_average) as sum_io_bytes_average 
		FROM sys.dm_os_tasks where task_state !='DONE' group by task_state;`,
		conn,
	)

	for rows.Next() {

		var task_state string
		var sum_context_switches int
		var sum_io_operation_count int
		var sum_io_bytes_count int64
		var sum_io_bytes_average int

		err := rows.Scan(
			&task_state,
			&sum_context_switches,
			&sum_io_operation_count,
			&sum_io_bytes_count,
			&sum_io_bytes_average,
		)
		if err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
		}

		metrics = append(metrics, returnMetric(
			"sql_task_sum_context_switches",
			"Current value of the sum of context_switches_count in dm_os_tasks by task_state",
			"task_state",
			task_state,
			float64(sum_context_switches),
		))

		metrics = append(metrics, returnMetric(
			"sql_task_sum_io_operation_count",
			"Current value of the sum of pending_io_count in dm_os_tasks by task_state",
			"task_state",
			task_state,
			float64(sum_io_operation_count),
		))

		metrics = append(metrics, returnMetric(
			"sql_task_sum_io_bytes_count",
			"Current value of the sum of pending_io_byte_count in dm_os_tasks by task_state",
			"task_state",
			task_state,
			float64(sum_io_bytes_count),
		))

		metrics = append(metrics, returnMetric(
			"sql_task_sum_io_bytes_average",
			"Current value of the sum of pending_io_byte_average in dm_os_tasks by task_state",
			"task_state",
			task_state,
			float64(sum_io_bytes_average),
		))

	}
	err := rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}
