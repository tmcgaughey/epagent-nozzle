package processors

import (
	"github.com/tmcgaughey/epagent-nozzle/metrics"
	"github.com/cloudfoundry/noaa/events"
	)

type ValueMetricProcessor struct{}

func NewValueMetricProcessor() *ValueMetricProcessor {
	return &ValueMetricProcessor{}
}

func (p *ValueMetricProcessor) Process(e *events.Envelope) []metrics.WMetric {
	processedMetrics := make([]metrics.WMetric, 1)
	valueMetricEvent := e.GetValueMetric()

	processedMetrics[0] = *p.ProcessValueMetric(valueMetricEvent, e.GetOrigin())

	return processedMetrics
}

func (p *ValueMetricProcessor) ProcessValueMetric(event *events.ValueMetric, origin string) *metrics.WMetric {
	statPrefix := "ops." + origin + "."
	valueMetricName := event.GetName()
	stat := statPrefix + valueMetricName
	metric := metrics.NewLongCounterMetric(stat, int64(event.GetValue()))

	return metric
}
