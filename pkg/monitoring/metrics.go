package monitoring

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"chaosmonkey.monke/chaos/pkg/logger"
	"github.com/VictoriaMetrics/metrics"
)

var log = logger.Logger.Sugar()

type MetricServer struct {
	set *metrics.Set
}

var MetricsServerInstance = &MetricServer{metrics.NewSet()}

func (ms *MetricServer) IncWithLabels(name string, labels map[string]string) {
	log.Debugf("Incrementing metric [%s] with labels [%s]", name, labels)
	metricName := getMetricNameWithLabels(name, labels)
	log.Debugf("Metric name: [%s]", metricName)
	ms.Inc(metricName)
}

func (ms *MetricServer) Inc(name string) {
	log.Debugf("Incrementing metric %s by 1", name)
	ms.set.GetOrCreateCounter(name).Inc()
}

func getMetricNameWithLabels(name string, labels map[string]string) string {
	if len(labels) == 0 {
		return strings.ReplaceAll(name, "-", "_")
	}

	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for indx, val := range keys {
		keys[indx] = fmt.Sprintf(`%s="%s"`, strings.ReplaceAll(val, "-", "_"), labels[val])
	}

	return fmt.Sprintf(`%s{%s}`, strings.ReplaceAll(name, "-", "_"), strings.Join(keys, ", "))
}

func Start() {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		MetricsServerInstance.set.WritePrometheus(w)
		metrics.WriteProcessMetrics(w)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`OK`))
	})

	log.Info("Starting metric server listening on port 8000")
	http.ListenAndServe(":8000", nil)
}
