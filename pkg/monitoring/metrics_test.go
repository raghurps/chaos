package monitoring

import (
	"testing"

	"github.com/VictoriaMetrics/metrics"
	"github.com/stretchr/testify/assert"
)

func TestCounters(t *testing.T) {
	metricName := "some_random_metric"
	metricWithLabelName := "some_random_metric_with_labels"
	labels := map[string]string{
		"someLabel":  "someValue",
		"otherLabel": "otherValue",
	}
	srv := &MetricServer{metrics.NewSet()}
	t.Run("should increment metric", func(t *testing.T) {
		srv.Inc(metricName)
		assert.Equal(t, uint64(1), srv.set.GetOrCreateCounter(metricName).Get())
	})

	t.Run("should increment metric with labels", func(t *testing.T) {
		generatedMetricName := getMetricNameWithLabels(metricWithLabelName, labels)
		srv.IncWithLabels(metricWithLabelName, labels)
		assert.Equal(t, uint64(1), srv.set.GetOrCreateCounter(generatedMetricName).Get())
	})
}

func TestMetricNameGeneration(t *testing.T) {
	metricName := "some_random_metric"

	t.Run("should return same metric name when no labels provided", func(t *testing.T) {
		assert.Equal(t, metricName, getMetricNameWithLabels(metricName, map[string]string{}))
	})

	labels := map[string]string{
		"someLabel":  "someValue",
		"otherLabel": "otherValue",
	}

	expectedMetricName := `some_random_metric{otherLabel="otherValue", someLabel="someValue"}`
	t.Run("should return metric name with labels", func(t *testing.T) {
		assert.Equal(t, expectedMetricName, getMetricNameWithLabels(metricName, labels))
	})

	metricName = "some-random-metric"

	t.Run("should replace - with _ in metric name", func(t *testing.T) {
		assert.Equal(t, "some_random_metric", getMetricNameWithLabels(metricName, map[string]string{}))
	})

	labels = map[string]string{
		"some-Label":  "someValue",
		"other-Label": "otherValue",
	}

	expectedMetricName = `some_random_metric{other_Label="otherValue", some_Label="someValue"}`
	t.Run("should replace - with _ in metric labels", func(t *testing.T) {
		assert.Equal(t, expectedMetricName, getMetricNameWithLabels(metricName, labels))
	})

}
