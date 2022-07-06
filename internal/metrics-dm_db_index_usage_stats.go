package internal

import (
	"database/sql"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getIndexUsageStatsStats(conn *sql.DB) []prometheus.Metric {
	metrics := []prometheus.Metric{}

	rows, err := performQuery(`SELECT
	database_id,
	object_id,
	index_id,
	user_seeks,
	user_scans,
	user_lookups,
	user_updates,
	system_seeks,
	system_scans,
	system_lookups,
	system_updates
	FROM sys.dm_db_index_usage_stats;`,
		conn,
	)
	if err != nil {
		logrus.Errorf("Error in query execution, skipping metrics")
		return metrics
	}

	var database_id int
	var object_id int
	var index_id int
	var user_seeks int64
	var user_scans int64
	var user_lookups int64
	var user_updates int64
	var system_seeks int64
	var system_scans int64
	var system_lookups int64
	var system_updates int64
	for rows.Next() {

		if err := rows.Scan(
			&database_id,
			&object_id,
			&index_id,
			&user_seeks,
			&user_scans,
			&user_lookups,
			&user_updates,
			&system_seeks,
			&system_scans,
			&system_lookups,
			&system_updates,
		); err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
			continue
		}

		metrics = append(metrics, returnMetric(
			"sql_index_usage_stats_user_seeks",
			"Current value of user_seeks in dm_db_index_usage_stats for the type",
			"index",
			fmt.Sprintf("%d_%d_%d", database_id, object_id, index_id),
			float64(user_seeks),
		))

		metrics = append(metrics, returnMetric(
			"sql_index_usage_stats_user_scans",
			"Current value of user_scans in dm_db_index_usage_stats for the type",
			"index",
			fmt.Sprintf("%d_%d_%d", database_id, object_id, index_id),
			float64(user_scans),
		))

		metrics = append(metrics, returnMetric(
			"sql_index_usage_stats_user_lookups",
			"Current value of user_lookups in dm_db_index_usage_stats for the type",
			"index",
			fmt.Sprintf("%d_%d_%d", database_id, object_id, index_id),
			float64(user_lookups),
		))

		metrics = append(metrics, returnMetric(
			"sql_index_usage_stats_user_updates",
			"Current value of user_updates in dm_db_index_usage_stats for the type",
			"index",
			fmt.Sprintf("%d_%d_%d", database_id, object_id, index_id),
			float64(user_updates),
		))

		metrics = append(metrics, returnMetric(
			"sql_index_usage_stats_system_seeks",
			"Current value of system_seeks in dm_db_index_usage_stats for the type",
			"index",
			fmt.Sprintf("%d_%d_%d", database_id, object_id, index_id),
			float64(system_seeks),
		))

		metrics = append(metrics, returnMetric(
			"sql_index_usage_stats_system_scans",
			"Current value of system_scans in dm_db_index_usage_stats for the type",
			"index",
			fmt.Sprintf("%d_%d_%d", database_id, object_id, index_id),
			float64(system_scans),
		))

		metrics = append(metrics, returnMetric(
			"sql_index_usage_stats_system_lookups",
			"Current value of system_lookups in dm_db_index_usage_stats for the type",
			"index",
			fmt.Sprintf("%d_%d_%d", database_id, object_id, index_id),
			float64(system_lookups),
		))

		metrics = append(metrics, returnMetric(
			"sql_index_usage_stats_system_updates",
			"Current value of system_updates in dm_db_index_usage_stats for the type",
			"index",
			fmt.Sprintf("%d_%d_%d", database_id, object_id, index_id),
			float64(system_updates),
		))

	}
	err = rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}
