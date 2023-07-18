package internal

import (
	"database/sql"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var dbname, dbhost string

func collectMetrics(db *sql.DB, database string, host string, enabledCollectors []string) []prometheus.Metric {
	dbname = database
	dbhost = host

	type metricsCollector = struct {
		name string
		f    func(*sql.DB) []prometheus.Metric
	}

	collectors := make([]metricsCollector, 0)

	var metrics []prometheus.Metric
	for _, collector := range enabledCollectors {
		switch collector {
		case "exec":
			collectors = append(collectors, metricsCollector{collector, getExecRequestStatusStats})
		case "filespace":
			collectors = append(collectors, metricsCollector{collector, getFileSpaceUsageStats})
		case "index":
			collectors = append(collectors, metricsCollector{collector, getIndexUsageStatsStats},
				metricsCollector{collector, getMissingIndexDetailsStats})
		case "locks":
			collectors = append(collectors, metricsCollector{collector, getSpinLockStats},
				metricsCollector{collector, getLatchStats})
		case "memory":
			collectors = append(collectors, metricsCollector{collector, getMemoryCacheHashtablesStats},
				metricsCollector{collector, getMemoryClerksStats},
				metricsCollector{collector, getMemoryObjectsStats})
		case "performance":
			collectors = append(collectors, metricsCollector{collector, getPerformanceCountersStats})
		case "replication":
			collectors = append(collectors, metricsCollector{collector, getReplicationStats})
		case "schedulers":
			collectors = append(collectors, metricsCollector{collector, getSchedulersStats})
		case "tasks":
			collectors = append(collectors, metricsCollector{collector, getTasksStats})
		case "waits":
			collectors = append(collectors, metricsCollector{collector, getWaitStatsStats})
		default:
			logrus.Errorf("invalid collector %s, skipping it", collector)
		}
	}

	// do this regardless
	collectors = append(collectors, metricsCollector{"settings", getSettings})

	for _, collector := range collectors {
		collectedMetrics := collector.f(db)
		if len(collectedMetrics) == 0 {
			logrus.Warnf("empty metrics list for collector %s", collector.name)
			continue
		}

		metrics = append(metrics, collectedMetrics...)
	}

	logrus.Debugf("Retrieved %d metrics", len(metrics))
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
	logrus.Debugf("labels: %v, desc: %v, value: %f", labelSanatized, labelDescSanatized, value)
	m, err := prometheus.NewConstMetric(
		prometheus.NewDesc(
			name,
			desc,
			labelDescSanatized, nil),
		prometheus.GaugeValue,
		value,
		labelSanatized...,
	)
	if err != nil {
		logrus.Errorf("creating metric failed: %s", err)
		return nil
	}
	return m
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}

func sanitizeStringForLabel(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Replace(s, " ", "_", -1)
	s = strings.Replace(s, "%", "pct", -1)
	s = strings.Replace(s, "/", "_per_", -1)
	s = strings.Replace(s, ".", "", -1)
	s = strings.Replace(s, ">=", "gte_", -1)
	s = strings.Replace(s, "<=", "lte_", -1)
	s = strings.Replace(s, ">", "gt_", -1)
	s = strings.Replace(s, "<", "lt_", -1)
	s = strings.Replace(s, "&", "_and_", -1)
	s = strings.Replace(s, "(", "", -1)
	s = strings.Replace(s, ")", "", -1)
	s = strings.Replace(s, ":", "_", -1)
	s = strings.Replace(s, "$", "_", -1)
	s = strings.Replace(s, ".", "", -1)
	s = strings.Replace(s, ".blob.core.windows.net", "", -1)

	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == '_' {
			result.WriteByte(b)
		}
	}
	//    return result.String()

	return strings.ToLower(result.String())
}
