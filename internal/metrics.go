package internal

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func returnMetrics(db *sql.DB) []prometheus.Metric {
	var metrics []prometheus.Metric
	metrics = append(metrics, getCurrentUserSessions(db)...)
	metrics = append(metrics, getSuspendedSessions(db)...)
	return metrics
}

func returnMetric(name, desc, labelDesc, label string, value float64) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			name,
			desc,
			[]string{labelDesc}, nil),
		prometheus.GaugeValue,
		value,
		[]string{label}...,
	)
}

func getCurrentUserSessions(conn *sql.DB) []prometheus.Metric {
	var metrics []prometheus.Metric

	rows := performQuery(
		"SELECT [status], COUNT(*) AS cnt FROM sys.dm_exec_requests WHERE session_id > 50 GROUP BY [status];",
		conn,
	)

	for rows.Next() {
		var status string
		var count int64
		err := rows.Scan(&status, &count)
		if err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
		}
		metrics = append(metrics, returnMetric(
			"sql_user_sessions",
			"Current user sessions",
			"status",
			status,
			float64(count),
		))
	}
	err := rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}

func getSuspendedSessions(conn *sql.DB) []prometheus.Metric {
	var metrics []prometheus.Metric

	rows := performQuery(
		"SELECT wait_type, COUNT(*) AS cnt FROM sys.dm_exec_requests WHERE session_id > 50 AND status = 'suspended' GROUP BY wait_type;",
		conn,
	)

	for rows.Next() {
		var waitTypes string
		var count int64
		err := rows.Scan(&waitTypes, &count)
		if err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
		}
		metrics = append(metrics, returnMetric(
			"sql_suspended_sessions",
			"Current suspended user sessions",
			"wait_type",
			waitTypes,
			float64(count),
		))
	}
	err := rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}

/*
SELECT
    --scheduler_address
    parent_node_id
	, scheduler_id
	, cpu_id
	--,[status]
	, is_online
	, is_idle
	, preemptive_switches_count
	, context_switches_count
	, idle_switches_count
	, current_tasks_count
	, runnable_tasks_count
	, current_workers_count
	, active_workers_count
	, work_queue_count
	, pending_disk_io_count
	, load_factor yield_count
	, last_timer_activity
	, failed_to_create_worker
	--,active_worker_address
	--,memory_object_address
	--,task_memory_object_address
	, quantum_length_us
FROM sys.dm_os_schedulers
WHERE status = 'VISIBLE ONLINE';

SELECT
SUM(total_page_count *1.0/128) AS Total_space_MB ,
SUM(unallocated_extent_page_count*1.0/128) AS Unallocated_Space_MB,
SUM(user_object_reserved_page_count*1.0/128) AS User_Obj_Allocated_Space_MB,
SUM(internal_object_reserved_page_count*1.0/128) AS Internal_Obj_Allocated_Space_MB,
(SUM(total_page_count)-SUM(unallocated_extent_page_count)-SUM(user_object_reserved_page_count)- SUM(internal_object_reserved_page_count) )*1.0/128 AS Other_Obj_Space_MB
FROM tempdb.sys.dm_db_file_space_usage

SELECT
	[type]
	,SUM(pages_kb)					   AS sum_pages_kb
	,SUM(virtual_memory_reserved_kb)   AS sum_virtual_memory_reserved_kb
	,SUM(virtual_memory_committed_kb)  AS sum_virtual_memory_committed_kb
	,SUM(shared_memory_reserved_kb)    AS sum_shared_memory_reserved_kb
	,SUM(shared_memory_committed_kb)   AS sum_shared_memory_committed_kb
FROM sys.dm_os_memory_clerks
GROUP BY [type]

SELECT @@SERVERNAME, SERVERPROPERTY('ProductVersion');

SELECT *
FROM sys.dm_os_performance_counters

SELECT
  wait_type, waiting_tasks_count, wait_time_ms
FROM sys.dm_os_wait_stats
WHERE [wait_type] IN (
  %%WAITS%%
);

*/
