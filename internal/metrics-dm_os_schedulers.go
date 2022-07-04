package internal

import (
	"database/sql"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getSchedulersStats(conn *sql.DB) []prometheus.Metric {
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
		var is_idle bool
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
		var failed_to_create_worker bool
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
			&total_cpu_usage_ms,
			&total_cpu_idle_capped_ms,
			&total_scheduler_delay_ms,
			&ideal_workers_limit,
		)
		if err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
		}
		schedulerIdName := strconv.FormatInt(scheduler_id, 10)

		metrics = append(metrics, returnMetric(
			"sql_scheduler_is_idle",
			"Current value of is_idle in dm_os_schedulers",
			"scheduler_id",
			schedulerIdName,
			float64(bool2int(is_idle)),
		))

		metrics = append(metrics, returnMetric(
			"sql_scheduler_preemptive_switches_count",
			"Current value of preemptive_switches_count in dm_os_schedulers",
			"scheduler_id",
			schedulerIdName,
			float64(preemptive_switches_count),
		))

		metrics = append(metrics, returnMetric(
			"sql_scheduler_context_switches_count",
			"Current value of context_switches_count in dm_os_schedulers",
			"scheduler_id",
			schedulerIdName,
			float64(context_switches_count),
		))

		metrics = append(metrics, returnMetric(
			"sql_scheduler_idle_switches_count",
			"Current value of idle_switches_count in dm_os_schedulers",
			"scheduler_id",
			schedulerIdName,
			float64(idle_switches_count),
		))

		metrics = append(metrics, returnMetric(
			"sql_scheduler_current_tasks_count",
			"Current value of current_tasks_count in dm_os_schedulers",
			"scheduler_id",
			schedulerIdName,
			float64(current_tasks_count),
		))

		metrics = append(metrics, returnMetric(
			"sql_scheduler_runnable_tasks_count",
			"Current value of runnable_tasks_count in dm_os_schedulers",
			"scheduler_id",
			schedulerIdName,
			float64(runnable_tasks_count),
		))

		metrics = append(metrics, returnMetric(
			"sql_scheduler_current_workers_count",
			"Current value of current_workers_count in dm_os_schedulers",
			"scheduler_id",
			schedulerIdName,
			float64(current_workers_count),
		))

		metrics = append(metrics, returnMetric(
			"sql_scheduler_active_workers_count",
			"Current value of active_workers_count in dm_os_schedulers",
			"scheduler_id",
			schedulerIdName,
			float64(active_workers_count),
		))

		metrics = append(metrics, returnMetric(
			"sql_scheduler_work_queue_count",
			"Current value of work_queue_count in dm_os_schedulers",
			"scheduler_id",
			schedulerIdName,
			float64(work_queue_count),
		))

		metrics = append(metrics, returnMetric(
			"sql_scheduler_pending_disk_io_count",
			"Current value of pending_disk_io_count in dm_os_schedulers",
			"scheduler_id",
			schedulerIdName,
			float64(pending_disk_io_count),
		))

		metrics = append(metrics, returnMetric(
			"sql_scheduler_queued_disk_io_count",
			"Current value of queued_disk_io_count in dm_os_schedulers",
			"scheduler_id",
			schedulerIdName,
			float64(queued_disk_io_count),
		))

		metrics = append(metrics, returnMetric(
			"sql_scheduler_load_factor",
			"Current value of load_factor in dm_os_schedulers",
			"scheduler_id",
			schedulerIdName,
			float64(load_factor),
		))

		metrics = append(metrics, returnMetric(
			"sql_scheduler_yield_count",
			"Current value of yield_count in dm_os_schedulers",
			"scheduler_id",
			schedulerIdName,
			float64(yield_count),
		))

		metrics = append(metrics, returnMetric(
			"sql_scheduler_last_timer_activity",
			"Current value of last_timer_activity in dm_os_schedulers",
			"scheduler_id",
			schedulerIdName,
			float64(last_timer_activity),
		))

		metrics = append(metrics, returnMetric(
			"sql_scheduler_failed_to_create_worker",
			"Current value of failed_to_create_worker in dm_os_schedulers",
			"scheduler_id",
			schedulerIdName,
			float64(bool2int(failed_to_create_worker)),
		))

		metrics = append(metrics, returnMetric(
			"sql_scheduler_total_cpu_usage_ms",
			"Current value of total_cpu_usage_ms in dm_os_schedulers",
			"scheduler_id",
			schedulerIdName,
			float64(total_cpu_usage_ms),
		))

		metrics = append(metrics, returnMetric(
			"sql_scheduler_total_cpu_idle_capped_ms",
			"Current value of total_cpu_idle_capped_ms in dm_os_schedulers",
			"scheduler_id",
			schedulerIdName,
			float64(total_cpu_idle_capped_ms),
		))

		metrics = append(metrics, returnMetric(
			"sql_scheduler_total_scheduler_delay_ms",
			"Current value of total_scheduler_delay_ms in dm_os_schedulers",
			"scheduler_id",
			schedulerIdName,
			float64(total_scheduler_delay_ms),
		))

		metrics = append(metrics, returnMetric(
			"sql_scheduler_ideal_workers_limit",
			"Current value of ideal_workers_limit in dm_os_schedulers",
			"scheduler_id",
			schedulerIdName,
			float64(ideal_workers_limit),
		))

	}
	err := rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}
