package internal

import (
	"database/sql"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getReplicationStats(conn *sql.DB) []prometheus.Metric {
	metrics := []prometheus.Metric{}

	// Only for primaries
	rows, err := performQuery(`SELECT 
	partner_server, 
	role_desc,
	replication_state,
	replication_lag_sec 
	FROM sys.dm_geo_replication_link_status
	WHERE role=0;`,
		conn,
	)
	if err != nil {
		logrus.Errorf("Error in query execution, skipping metrics")
		return metrics
	}
	var partner_server string
	var role_desc string
	var replication_state int
	var replication_lag_sec int
	for rows.Next() {

		if err := rows.Scan(
			&partner_server,
			&role_desc,
			&replication_state,
			&replication_lag_sec,
		); err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
			continue
		}

		metrics = append(metrics, returnMetric(
			"sql_replication_state",
			"State of the replicaion",
			"type",
			fmt.Sprintf("%s_for_%s", role_desc, partner_server),
			float64(replication_state),
		))

		metrics = append(metrics, returnMetric(
			"sql_replication_lag_seconds",
			"Lag in seconds of the replica",
			"type",
			fmt.Sprintf("%s_for_%s", role_desc, partner_server),
			float64(replication_lag_sec),
		))

	}
	err = rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}
