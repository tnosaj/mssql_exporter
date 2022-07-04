package internal

import (
	"database/sql"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getPerformanceCountersStats(conn *sql.DB) []prometheus.Metric {
	var metrics []prometheus.Metric

	rows := performQuery(`SELECT 
	object_name, 
	counter_name, 
	CASE WHEN instance_name != "" THEN instance_name ELSE "empty" END as instance_name,
	cntr_value
	FROM sys.dm_os_performance_counters
	WHERE 
	cntr_type = 65792 OR
	cntr_type = 537003264 OR
	cntr_type = 1073939712 OR
	cntr_type = 1073874176 OR
	cntr_type=272696576;`,
		conn,
	)

	/*

	    https://techcommunity.microsoft.com/t5/sql-server-support-blog/interpreting-the-counter-values-from-sys-dm-os-performance/ba-p/317824
	   	https://troubleshootingsql.com/2011/03/03/what-does-cntr_type-mean/

	    65792 - PERF_COUNTER_LARGE_RAWCOUNT (This counter value shows the last observed value directly. Primarily used to track counts of objects)

	    537003264 - PERF_LARGE_RAW_FRACTION (This counter value represents a fractional value as a ratio to its corresponding PERF_LARGE_RAW_BASE counter value)
	   	1073939712 - PERF_LARGE_RAW_BASE (This counter value is raw data that is used as the denominator of a counter that presents a instantaneous arithmetic fraction)
	   	1073874176 - PERF_AVERAGE_BULK (This counter value represents an average metric)

	   	272696320 - PERF_COUNTER_COUNTER (Average number of operations completed during each second of the sample interval)
	   	272696576 - PERF_COUNTER_BULK_COUNT (The value is obtained by taking two samples of the PERF_COUNTER_BULK_COUNT value.)
	*/

	for rows.Next() {

		var object_name, counter_name, instance_name string
		var cntr_value int64

		err := rows.Scan(
			&object_name,
			&counter_name,
			&instance_name,
			&cntr_value,
		)
		if err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
		}

		labelName := fmt.Sprintf(
			"%s_%s_%s",
			sanitizeStringForLabel(object_name),
			sanitizeStringForLabel(counter_name),
			sanitizeStringForLabel(instance_name),
		)

		metrics = append(metrics, returnMetric(
			"sql_performance_counters",
			"Current value of the counter type in dm_os_performance_counters for the type",
			"counter_type",
			labelName,
			float64(cntr_value),
		))

	}
	err := rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}
