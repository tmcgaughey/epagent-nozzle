package processors

import (
	"github.com/tmcgaughey/epagent-nozzle/metrics"
	"github.com/cloudfoundry/noaa/events"
)

type HeartbeatProcessor struct{}

func NewHeartbeatProcessor() *HeartbeatProcessor {
	return &HeartbeatProcessor{}
}

func (p *HeartbeatProcessor) Process(e *events.Envelope) []metrics.WMetric {
	processedMetrics := make([]metrics.WMetric, 4)
	heartbeat := e.GetHeartbeat()
	origin := e.GetOrigin()

	processedMetrics[0] = *p.ProcessHeartbeatCount(heartbeat, origin)
	processedMetrics[1] = *p.ProcessHeartbeatEventsSentCount(heartbeat, origin)
	processedMetrics[2] = *p.ProcessHeartbeatEventsReceivedCount(heartbeat, origin)
	processedMetrics[3] = *p.ProcessHeartbeatEventsErrorCount(heartbeat, origin)

	return processedMetrics
}

func (p *HeartbeatProcessor) ProcessHeartbeatCount(e *events.Heartbeat, origin string) *metrics.WMetric {
	stat := "ops." + origin + ".heartbeats.count"
	metric := metrics.NewPerintervalCounterMetric(stat, int32(1))

	return metric
}

func (p *HeartbeatProcessor) ProcessHeartbeatEventsSentCount(e *events.Heartbeat, origin string) *metrics.WMetric {
	stat := "ops." + origin + ".heartbeats.eventsSentCount"
	metric := metrics.NewPerintervalCounterMetric(stat, int32(e.GetSentCount()))

	return metric
}

func (p *HeartbeatProcessor) ProcessHeartbeatEventsReceivedCount(e *events.Heartbeat, origin string) *metrics.WMetric {
	stat := "ops." + origin + ".heartbeats.eventsReceivedCount"
	metric := metrics.NewPerintervalCounterMetric(stat, int32(e.GetReceivedCount()))

	return metric
}

func (p *HeartbeatProcessor) ProcessHeartbeatEventsErrorCount(e *events.Heartbeat, origin string) *metrics.WMetric {
	stat := "ops." + origin + ".heartbeats.eventsErrorCount"
	metric := metrics.NewLongCounterMetric(stat, int64(e.GetErrorCount()))

	return metric
}
