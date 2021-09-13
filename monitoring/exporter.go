package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type ReadingCounter struct {
	prometheus.CounterVec
}

type ReadingCounterInterface interface {
	Increment(string, string)
}

type ReadingGauge struct {
	prometheus.GaugeVec
}

type ReadingGaugeInterface interface {
	Set(float64, string)
}

func NewReadingCounter() ReadingCounterInterface {
	var rc ReadingCounter

	rc.CounterVec = *promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "reading_count",
		Help: "The total number of processed readings by the coordinator"},
		[]string{"layer", "sensor"})

	return &rc
}

func (rc *ReadingCounter) Increment(layer, sensor string) {
	rc.With(prometheus.Labels{"layer": layer, "sensor": sensor}).Inc()
}

func NewReadingGauge() ReadingGaugeInterface {
	var rg ReadingGauge

	rg.GaugeVec = *promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sensor_reading",
		Help: "The current reading from the sensor"},
		[]string{"sensor"})

	return &rg

}

func (rg *ReadingGauge) Set(value float64, sensor string) {
	rg.With(prometheus.Labels{"sensor": sensor}).Set(value)
}

func MetricExporter() {
	http.Handle("/metrics", promhttp.Handler())
	_ = http.ListenAndServe(":2112", nil)
}
