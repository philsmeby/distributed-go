package main

import (
	"github.com/distributed-go/coordinator"
	"github.com/distributed-go/monitoring"
)

func main() {
	ea := coordinator.NewEventAggregator()
	_ = coordinator.NewDatabaseConsumer(ea)
	_ = coordinator.NewMetricConsumer(ea)

	ql := coordinator.NewQueueListener(ea)
	go ql.ListenForNewSource()

	monitoring.MetricExporter()
}
