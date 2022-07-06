package internal

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getExecRequestStatusStats(conn *sql.DB) []prometheus.Metric {
	metrics := []prometheus.Metric{}

	rows, err := performQuery(
		"SELECT [status], COUNT(*) AS cnt FROM sys.dm_exec_requests WHERE session_id > 50 GROUP BY [status];",
		conn,
	)
	if err != nil {
		logrus.Errorf("Error in query ExecRequestStatusStats execution, skipping metrics: %s", err)
		return metrics
	}

	// Don't generate them over and over, that's more allocations
	var status string
	var count int64

	for rows.Next() {
		if err := rows.Scan(&status, &count); err != nil {
			logrus.Errorf("Failed scanning request status stats with error: %s", err)
			continue // Skip, otherwise you are repeating crap
		}

		metrics = append(metrics, returnMetric(
			"sql_user_sessions",
			"Current user sessions",
			"status",
			status,
			float64(count),
		))
	}

	return metrics
}

func getExecRequestSuspendedStats(conn *sql.DB) []prometheus.Metric {
	metrics := []prometheus.Metric{}

	rows, err := performQuery(
		"SELECT wait_type, COUNT(*) AS cnt FROM sys.dm_exec_requests WHERE session_id > 50 AND status = 'suspended' GROUP BY wait_type;",
		conn,
	)
	if err != nil {
		logrus.Errorf("Error in query execution, skipping metrics")
		return metrics
	}
	var waitTypes string
	var count int64
	for rows.Next() {

		if err := rows.Scan(&waitTypes, &count); err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
			continue
		}
		metrics = append(metrics, returnMetric(
			"sql_suspended_sessions",
			"Current suspended user sessions",
			"wait_type",
			waitTypes,
			float64(count),
		))
	}
	err = rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}
