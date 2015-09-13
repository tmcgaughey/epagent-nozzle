package processors

import (
	"github.com/tmcgaughey/epagent-nozzle/metrics"
	"github.com/cloudfoundry/noaa/events"
)

type CounterProcessor struct{}

func NewCounterProcessor() *CounterProcessor {
	return &CounterProcessor{}
}

func (p *CounterProcessor) Process(e *events.Envelope) []metrics.WMetric {
	processedMetrics := make([]metrics.WMetric, 1)
	counterEvent := e.GetCounterEvent()

	processedMetrics[0] = *p.ProcessCounter(counterEvent)

	return processedMetrics
}

func (p *CounterProcessor) ProcessCounter(event *events.CounterEvent) *metrics.WMetric {
	stat := "ops." + event.GetName()
	metric := metrics.NewLongCounterMetric(stat, int64(event.GetDelta()))

	return metric
}
