package internal

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getSpinLockStats(conn *sql.DB) []prometheus.Metric {
	var metrics []prometheus.Metric

	rows := performQuery(`SELECT
		name,
		collisions,
		spins,
		spins_per_collision,
		sleep_time,
		backoffs
	  FROM sys.dm_os_spinlock_stats;`,
		conn,
	)

	for rows.Next() {

		var name string
		var collisions int64
		var spins int64
		var spins_per_collision float32
		var sleep_time int64
		var backoffs int64

		err := rows.Scan(
			&name,
			&collisions,
			&spins,
			&spins_per_collision,
			&sleep_time,
			&backoffs,
		)
		if err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
		}

		metrics = append(metrics, returnMetric(
			"sql_spinlock_stats_collisions",
			"Current value of collisions in dm_os_spinlock_stats",
			"spinlock_type",
			name,
			float64(collisions),
		))

		metrics = append(metrics, returnMetric(
			"sql_spinlock_stats_spins",
			"Current value of spins in dm_os_spinlock_stats",
			"spinlock_type",
			name,
			float64(spins),
		))

		metrics = append(metrics, returnMetric(
			"sql_spinlock_stats_spins_per_collision",
			"Current value of spins_per_collision in dm_os_spinlock_stats",
			"spinlock_type",
			name,
			float64(spins_per_collision),
		))

		metrics = append(metrics, returnMetric(
			"sql_spinlock_stats_sleep_time",
			"Current value of sleep_time in dm_os_spinlock_stats",
			"spinlock_type",
			name,
			float64(sleep_time),
		))

		metrics = append(metrics, returnMetric(
			"sql_spinlock_stats_backoffs",
			"Current value of backoffs in dm_os_spinlock_stats",
			"spinlock_type",
			name,
			float64(backoffs),
		))
	}
	err := rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}
