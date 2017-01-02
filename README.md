# go-metrics-prometheus

This is a reporter for the go-metrics library which will post the metrics to the prometheus client registry.
It just updates the registry, taking care of exporting the metrics is still your responsibility.
It is based on https://github.com/deathowl/go-metrics-prometheus

Usage:

```
import (
  "github.com/MeteoGroup/go-metrics-prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  "github.com/prometheus/client_golang/prometheus"
  metrics "github.com/rcrowley/go-metrics"
  )

// create metrics registry
metricsRegistry := metrics.NewRegistry()
appGauge := metrics.GetOrRegisterGauge("m1", metricsRegistry)
appGauge.Update(1)

// register in prometheus
pClient := NewPrometheusProvider(metricsRegistry, "test", "subsys", 1*time.Second)
go pClient.UpdatePrometheusMetrics()

// export http for prometheus

go func() {
  http.Handle("/metrics", promhttp.HandlerFor(promClient.promRegistry, promhttp.HandlerOpts{}))
  log.Fatal(http.ListenAndServe(":8080", nil))
}()

```
