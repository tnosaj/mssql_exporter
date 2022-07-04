package internal

import (
	"database/sql"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getMissingIndexDetailsStats(conn *sql.DB) []prometheus.Metric {
	var metrics []prometheus.Metric

	rows := performQuery(`select 
	  database_id, 
	  object_id, 
	  count(statement) 
	from sys.dm_db_missing_index_details 
	group by database_id, object_id;`,
		conn,
	)

	for rows.Next() {
		var database_id int
		var object_id int
		var count int

		err := rows.Scan(
			&database_id,
			&object_id,
			&count,
		)
		if err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
		}

		metrics = append(metrics, returnMetric(
			"sql_missing_indexes_count",
			"Potential number of missing indexes on a table",
			"index",
			fmt.Sprintf("%d_%d", database_id, object_id),
			float64(count),
		))

	}
	err := rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}
