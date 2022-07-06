package internal

import (
	"database/sql"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getMissingIndexDetailsStats(conn *sql.DB) []prometheus.Metric {
	metrics := []prometheus.Metric{}

	rows, err := performQuery(`select 
	  database_id, 
	  object_id, 
	  count(statement) 
	from sys.dm_db_missing_index_details 
	group by database_id, object_id;`,
		conn,
	)
	if err != nil {
		logrus.Errorf("Error in query execution, skipping metrics")
		return metrics
	}
	var database_id int
	var object_id int
	var count int
	for rows.Next() {

		if err := rows.Scan(
			&database_id,
			&object_id,
			&count,
		); err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
			continue
		}

		metrics = append(metrics, returnMetric(
			"sql_missing_indexes_count",
			"Potential number of missing indexes on a table",
			"index",
			fmt.Sprintf("%d_%d", database_id, object_id),
			float64(count),
		))

	}
	err = rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}
