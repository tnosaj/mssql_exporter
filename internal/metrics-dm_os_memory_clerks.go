package internal

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getMemoryClerksStats(conn *sql.DB) []prometheus.Metric {
	var metrics []prometheus.Metric

	rows := performQuery(`SELECT
	[type]
	,SUM(pages_kb)					   AS sum_pages_kb
	,SUM(virtual_memory_reserved_kb)   AS sum_virtual_memory_reserved_kb
	,SUM(virtual_memory_committed_kb)  AS sum_virtual_memory_committed_kb
	,SUM(shared_memory_reserved_kb)    AS sum_shared_memory_reserved_kb
	,SUM(shared_memory_committed_kb)   AS sum_shared_memory_committed_kb
    FROM sys.dm_os_memory_clerks
    GROUP BY [type];`,
		conn,
	)

	for rows.Next() {

		var ttype string
		var sum_pages_kb int64
		var sum_virtual_memory_reserved_kb int64
		var sum_virtual_memory_committed_kb int64
		var sum_shared_memory_reserved_kb int64
		var sum_shared_memory_committed_kb int64

		err := rows.Scan(
			&ttype,
			&sum_pages_kb,
			&sum_virtual_memory_reserved_kb,
			&sum_virtual_memory_committed_kb,
			&sum_shared_memory_reserved_kb,
			&sum_shared_memory_committed_kb,
		)
		if err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
		}

		metrics = append(metrics, returnMetric(
			"sql_memory_clerks_pages_kb",
			"Current value of sum_pages_kb in dm_os_memory_clerks for the type",
			"type",
			ttype,
			float64(sum_pages_kb),
		))

		metrics = append(metrics, returnMetric(
			"sql_memory_clerks_virtual_memory_reserved_kb",
			"Current value of sum_virtual_memory_reserved_kb in dm_os_memory_clerks for the type",
			"type",
			ttype,
			float64(sum_virtual_memory_reserved_kb),
		))

		metrics = append(metrics, returnMetric(
			"sql_memory_clerks_virtual_memory_committed_kb",
			"Current value of sum_virtual_memory_committed_kb in dm_os_memory_clerks for the type",
			"type",
			ttype,
			float64(sum_virtual_memory_committed_kb),
		))

		metrics = append(metrics, returnMetric(
			"sql_memory_clerks_shared_memory_reserved_kb",
			"Current value of sum_shared_memory_reserved_kb in dm_os_memory_clerks for the type",
			"type",
			ttype,
			float64(sum_shared_memory_reserved_kb),
		))

		metrics = append(metrics, returnMetric(
			"sql_memory_clerks_shared_memory_committed_kb",
			"Current value of sum_shared_memory_committed_kb in dm_os_memory_clerks for the type",
			"type",
			ttype,
			float64(sum_shared_memory_committed_kb),
		))

	}
	err := rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}
