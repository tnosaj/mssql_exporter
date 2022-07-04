package internal

import (
	"database/sql"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var dbname, dbhost string

func returnMetrics(db *sql.DB, database string, host string, enabledMetrics []string) []prometheus.Metric {
	dbname = database
	dbhost = host
	var metrics []prometheus.Metric

	for _, enabledMetric := range enabledMetrics {
		switch enabledMetric {
		case "exec":
			metrics = append(metrics, checkToAppend(getExecRequestStatusStats(db))...)
		case "filespace":
			metrics = append(metrics, checkToAppend(getFileSpaceUsageStats(db))...)
		case "index":
			metrics = append(metrics, checkToAppend(getIndexUsageStatsStats(db))...)
			metrics = append(metrics, checkToAppend(getMissingIndexDetailsStats(db))...)
		case "memory":
			metrics = append(metrics, checkToAppend(getMemoryCacheHashtablesStats(db))...)
			metrics = append(metrics, checkToAppend(getMemoryClerksStats(db))...)
			metrics = append(metrics, checkToAppend(getMemoryObjectsStats(db))...)
		case "performance":
			metrics = append(metrics, checkToAppend(getPerformanceCountersStats(db))...)
			metrics = append(metrics, checkToAppend(getLatchStats(db))...)
			metrics = append(metrics, checkToAppend(getSpinLockStats(db))...)
			metrics = append(metrics, checkToAppend(getBufferDescriptorsStats(db))...)
		case "schedulers":
			metrics = append(metrics, checkToAppend(getSchedulersStats(db))...)
		case "tasks":
			metrics = append(metrics, checkToAppend(getTasksStats(db))...)
		case "waits":
			metrics = append(metrics, checkToAppend(getWaitStatsStats(db))...)
		}
	}
	logrus.Infof("Retrieved %d metrics", len(metrics))
	return metrics
}

func returnMetric(name, desc, labelDesc, label string, value float64) prometheus.Metric {
	// add dbhost and dbname
	labelDescSanatized := []string{"database_name", "database_host"}
	labelSanatized := []string{dbname, dbhost}
	if labelDesc != "none" {
		labelDescSanatized = append(labelDescSanatized, []string{sanitizeStringForLabel(labelDesc)}...)
		labelSanatized = append(labelSanatized, []string{sanitizeStringForLabel(label)}...)
	}
	return prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			name,
			desc,
			labelDescSanatized, nil),
		prometheus.GaugeValue,
		value,
		labelSanatized...,
	)
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}

func sanitizeStringForLabel(s string) string {
	s = strings.Replace(s, " ", "_", -1)
	s = strings.Replace(s, "%", "percent", -1)
	return strings.ToLower(s)
}

func checkToAppend(m []prometheus.Metric) []prometheus.Metric {
	if len(m) > 0 {
		return m
	}
	return []prometheus.Metric{}
}
