package internal

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getPerformanceCounterStats(conn *sql.DB) []prometheus.Metric {
	var metrics []prometheus.Metric

	rows := performQuery(`SELECT
		scheduler_id,
		is_idle,
		preemptive_switches_count,
		context_switches_count,
		idle_switches_count,
		current_tasks_count,
		runnable_tasks_count,
		current_workers_count,
		active_workers_count,
		work_queue_count,
		pending_disk_io_count,
		queued_disk_io_count,
		load_factor,
		yield_count,
		last_timer_activity,
		failed_to_create_worker,
		total_cpu_usage_ms,
		total_cpu_idle_capped_ms,
		total_scheduler_delay_ms,
		ideal_workers_limit
	  FROM sys.dm_os_schedulers
	  WHERE status = 'VISIBLE ONLINE';`,
		conn,
	)

	for rows.Next() {
		var scheduler_id int64
		var is_idle int64
		var preemptive_switches_count int64
		var context_switches_count int64
		var idle_switches_count int64
		var current_tasks_count int64
		var runnable_tasks_count int64
		var current_workers_count int64
		var active_workers_count int64
		var work_queue_count int64
		var pending_disk_io_count int64
		var queued_disk_io_count int64
		var load_factor int64
		var yield_count int64
		var last_timer_activity int64
		var failed_to_create_worker int64
		var active_worker_address int64
		var memory_object_address int64
		var task_memory_object_address int64
		var quantum_length_us int64
		var total_cpu_usage_ms int64
		var total_cpu_idle_capped_ms int64
		var total_scheduler_delay_ms int64
		var ideal_workers_limit int64

		err := rows.Scan(
			&scheduler_id,
			&is_idle,
			&preemptive_switches_count,
			&context_switches_count,
			&idle_switches_count,
			&current_tasks_count,
			&runnable_tasks_count,
			&current_workers_count,
			&active_workers_count,
			&work_queue_count,
			&pending_disk_io_count,
			&queued_disk_io_count,
			&load_factor,
			&yield_count,
			&last_timer_activity,
			&failed_to_create_worker,
			&active_worker_address,
			&memory_object_address,
			&task_memory_object_address,
			&quantum_length_us,
			&total_cpu_usage_ms,
			&total_cpu_idle_capped_ms,
			&total_scheduler_delay_ms,
			&ideal_workers_limit,
		)
		if err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
		}
		// metrics = append(metrics, returnMetric(
		// 	"sql_suspended_sessions",
		// 	"Current suspended user sessions",
		// 	"wait_type",
		// 	waitTypes,
		// 	float64(count),
		// ))
	}
	err := rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}

/*
scheduler_id is_idle preemptive_switches_count context_switches_count idle_switches_count current_tasks_count runnable_tasks_count current_workers_count active_workers_count work_queue_count pending_disk_io_count queued_disk_io_count load_factor yield_count last_timer_activity failed_to_create_worker active_worker_address memory_object_address task_memory_object_address quantum_length_us total_cpu_usage_ms total_cpu_idle_capped_ms total_scheduler_delay_ms ideal_workers_limit
SELECT
  scheduler_id,
  is_idle,
  preemptive_switches_count,
  context_switches_count,
  idle_switches_count,
  current_tasks_count,
  runnable_tasks_count,
  current_workers_count,
  active_workers_count,
  work_queue_count,
  pending_disk_io_count,
  queued_disk_io_count,
  load_factor,
  yield_count,
  last_timer_activity,
  failed_to_create_worker,
  total_cpu_usage_ms,
  total_cpu_idle_capped_ms,
  total_scheduler_delay_ms,
  ideal_workers_limit
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

SELECT * FROM sys.dm_os_performance_counters

SELECT
  wait_type, waiting_tasks_count, wait_time_ms
FROM sys.dm_os_wait_stats
WHERE [wait_type] IN (
  %%WAITS%%
);

*/
