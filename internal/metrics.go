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
