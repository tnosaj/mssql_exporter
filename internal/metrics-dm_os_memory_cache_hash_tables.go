package internal

import (
	"database/sql"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getMemoryCacheHashtablesStats(conn *sql.DB) []prometheus.Metric {
	var metrics []prometheus.Metric

	rows := performQuery(`SELECT    
	name, 
	type,     
	table_level,    
	sum(buckets_count) as buckets_countr,     
	sum(buckets_in_use_count),    
	sum(buckets_min_length),    
	sum(buckets_max_length),    
	sum(buckets_avg_length),    
	sum(buckets_max_length_ever),     
	sum(hits_count),    
	sum(misses_count),    
	sum(buckets_avg_scan_hit_length),     
	sum(buckets_avg_scan_miss_length)     
	FROM sys.dm_os_memory_cache_hash_tables 
	GROUP BY name, type, table_level;`,
		conn,
	)

	for rows.Next() {
		var name string
		var ttype string
		var table_level int
		var buckets_count int
		var buckets_in_use_count int
		var buckets_min_length int
		var buckets_max_length int
		var buckets_avg_length int
		var buckets_max_length_ever int
		var hits_count int64
		var misses_count int64
		var buckets_avg_scan_hit_length int
		var buckets_avg_scan_miss_length int

		err := rows.Scan(
			&name,
			&ttype,
			&table_level,
			&buckets_count,
			&buckets_in_use_count,
			&buckets_min_length,
			&buckets_max_length,
			&buckets_avg_length,
			&buckets_max_length_ever,
			&hits_count,
			&misses_count,
			&buckets_avg_scan_hit_length,
			&buckets_avg_scan_miss_length,
		)
		if err != nil {
			logrus.Errorf("Failed to scan with error: %s", err)
		}

		metrics = append(metrics, returnMetric(
			"sql_memory_cache_hash_tables_buckets_count",
			"Current value of buckets_count in dm_os_memory_cache_hash_tables for the type",
			"type",
			fmt.Sprintf("%s_%s_%d", name, ttype, table_level),
			float64(buckets_count),
		))

		metrics = append(metrics, returnMetric(
			"sql_memory_cache_hash_tables_buckets_in_use_count",
			"Current value of buckets_in_use_count in dm_os_memory_cache_hash_tables for the type",
			"type",
			fmt.Sprintf("%s_%s_%d", name, ttype, table_level),
			float64(buckets_in_use_count),
		))

		metrics = append(metrics, returnMetric(
			"sql_memory_cache_hash_tables_buckets_min_length",
			"Current value of buckets_min_length in dm_os_memory_cache_hash_tables for the type",
			"type",
			fmt.Sprintf("%s_%s_%d", name, ttype, table_level),
			float64(buckets_min_length),
		))

		metrics = append(metrics, returnMetric(
			"sql_memory_cache_hash_tables_buckets_max_length",
			"Current value of buckets_max_length in dm_os_memory_cache_hash_tables for the type",
			"type",
			fmt.Sprintf("%s_%s_%d", name, ttype, table_level),
			float64(buckets_max_length),
		))

		metrics = append(metrics, returnMetric(
			"sql_memory_cache_hash_tables_buckets_avg_length",
			"Current value of buckets_avg_length in dm_os_memory_cache_hash_tables for the type",
			"type",
			fmt.Sprintf("%s_%s_%d", name, ttype, table_level),
			float64(buckets_avg_length),
		))

		metrics = append(metrics, returnMetric(
			"sql_memory_cache_hash_tables_buckets_max_length_ever",
			"Current value of buckets_max_length_ever in dm_os_memory_cache_hash_tables for the type",
			"type",
			fmt.Sprintf("%s_%s_%d", name, ttype, table_level),
			float64(buckets_max_length_ever),
		))

		metrics = append(metrics, returnMetric(
			"sql_memory_cache_hash_tables_hits_count",
			"Current value of hits_count in dm_os_memory_cache_hash_tables for the type",
			"type",
			fmt.Sprintf("%s_%s_%d", name, ttype, table_level),
			float64(hits_count),
		))

		metrics = append(metrics, returnMetric(
			"sql_memory_cache_hash_tables_misses_count",
			"Current value of misses_count in dm_os_memory_cache_hash_tables for the type",
			"type",
			fmt.Sprintf("%s_%s_%d", name, ttype, table_level),
			float64(misses_count),
		))

		metrics = append(metrics, returnMetric(
			"sql_memory_cache_hash_tables_buckets_avg_scan_hit_length",
			"Current value of buckets_avg_scan_hit_length in dm_os_memory_cache_hash_tables for the type",
			"type",
			fmt.Sprintf("%s_%s_%d", name, ttype, table_level),
			float64(buckets_avg_scan_hit_length),
		))

		metrics = append(metrics, returnMetric(
			"sql_memory_cache_hash_tables_buckets_avg_scan_miss_length",
			"Current value of buckets_avg_scan_miss_length in dm_os_memory_cache_hash_tables for the type",
			"type",
			fmt.Sprintf("%s_%s_%d", name, ttype, table_level),
			float64(buckets_avg_scan_miss_length),
		))

	}
	err := rows.Err()
	if err != nil {
		logrus.Errorf("Scan failed %s:", err)
	}
	return metrics
}
