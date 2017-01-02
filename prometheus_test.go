package prometheusmetrics

import (
	"fmt"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

func TestUpdatePrometheusMetricsOnce(t *testing.T) {
	metricsRegistry := metrics.NewRegistry()
	pClient := NewPrometheusProvider(metricsRegistry, "test", "subsys", 1*time.Second)
	metricsRegistry.Register("counter", metrics.NewCounter())
	pClient.UpdatePrometheusMetricsOnce()
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "test",
		Subsystem: "subsys",
		Name:      "counter",
		Help:      "counter",
	})
	err := pClient.promRegistry.Register(gauge)
	if err == nil {
		t.Fatalf("Go-metrics registry didn't get registered to prometheus registry")
	}

}

func TestUpdatePrometheusMetrics(t *testing.T) {
	metricsRegistry := metrics.NewRegistry()
	pClient := NewPrometheusProvider(metricsRegistry, "test", "subsys", 1*time.Second)
	metricsRegistry.Register("counter", metrics.NewCounter())
	go pClient.UpdatePrometheusMetrics()
	time.Sleep(2 * time.Second)
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "test",
		Subsystem: "subsys",
		Name:      "counter",
		Help:      "counter",
	})
	err := pClient.promRegistry.Register(gauge)
	if err == nil {
		t.Fatalf("Go-metrics registry didn't get registered to prometheus registry")
	}

}

func TestPrometheusMetricsGetUpdated(t *testing.T) {
	metricsRegistry := metrics.NewRegistry()
	pClient := NewPrometheusProvider(metricsRegistry, "a", "a", 1*time.Second)
	cntr := metrics.NewCounter()
	metricsRegistry.Register("counter", cntr)
	cntr.Inc(2)
	go pClient.UpdatePrometheusMetrics()
	cntr.Inc(13)
	time.Sleep(5 * time.Second)
	metrics, _ := pClient.promRegistry.Gather()
	serialized := fmt.Sprint(metrics[0])
	expected := fmt.Sprintf("name:\"a_a_counter\" help:\"counter\" type:GAUGE metric:<gauge:<value:%d > > ", cntr.Count())
	if serialized != expected {
		t.Fatalf("Go-metrics value and prometheus metics value do not match, got %s , expected %s", serialized, expected)
	}
}
