package internal

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getFileSpaceUsageStats(conn *sql.DB) []prometheus.Metric {
	metrics := []prometheus.Metric{}

	rows, err := performQuery(`SELECT
	SUM(total_page_count *1.0/128) AS Total_space_MB ,
	SUM(unallocated_extent_page_count*1.0/128) AS Unallocated_Space_MB,
	SUM(user_object_reserved_page_count*1.0/128) AS User_Obj_Allocated_Space_MB,
	SUM(internal_object_reserved_page_count*1.0/128) AS Internal_Obj_Allocated_Space_MB,
	(SUM(total_page_count)-SUM(unallocated_extent_page_count)-SUM(user_object_reserved_page_count)- SUM(internal_object_reserved_page_count) )*1.0/128 AS Other_Obj_Space_MB
	FROM tempdb.sys.dm_db_file_space_usage;`,
		conn,
	)
	if err != nil {
		logrus.Errorf("Error in query execution, skipping metrics")
		return metrics
	}
	var total_space_mb float32
	var unallocated_space_mb float32
	var user_obj_allocated_space_mb float32
	var internal_obj_allocated_space_mb float32
	var other_obj_allocated_space_mb float32

	for rows.Next() {

		if err := rows.Scan(
			&total_space_mb,
			&unallocated_space_mb,
			&user_obj_allocated_space_mb,
			&internal_obj_allocated_space_mb,
			&other_obj_allocated_space_mb,
		); err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
			continue
		}

		metrics = append(metrics, returnMetric(
			"sql_file_space_usage_total_space_mb",
			"Current value of system total_space_mb in dm_os_file_space_usage",
			"none",
			"",
			float64(total_space_mb),
		))

		metrics = append(metrics, returnMetric(
			"sql_file_space_usage_unallocated_space_mb",
			"Current value of system unallocated_space_mb in dm_os_file_space_usage",
			"none",
			"",
			float64(unallocated_space_mb),
		))

		metrics = append(metrics, returnMetric(
			"sql_file_space_usage_user_obj_allocated_space_mb",
			"Current value of system user_obj_allocated_space_mb in dm_os_file_space_usage",
			"none",
			"",
			float64(user_obj_allocated_space_mb),
		))

		metrics = append(metrics, returnMetric(
			"sql_file_space_usage_internal_obj_allocated_space_mb",
			"Current value of system internal_obj_allocated_space_mb in dm_os_file_space_usage",
			"none",
			"",
			float64(internal_obj_allocated_space_mb),
		))

		metrics = append(metrics, returnMetric(
			"sql_file_space_usage_other_obj_allocated_space_mb",
			"Current value of system other_obj_allocated_space_mb in dm_os_file_space_usage",
			"none",
			"",
			float64(other_obj_allocated_space_mb),
		))
	}
	err = rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}
