package processors

import (
	"strconv"
	"strings"

	"github.com/tmcgaughey/epagent-nozzle/metrics"
	"github.com/cloudfoundry/noaa/events"
)

type HttpStartStopProcessor struct{}

func NewHttpStartStopProcessor() *HttpStartStopProcessor {
	return &HttpStartStopProcessor{}
}

func (p *HttpStartStopProcessor) Process(e *events.Envelope) []metrics.WMetric {
	processedMetrics := make([]metrics.WMetric, 4)
	httpStartStopEvent := e.GetHttpStartStop()

	processedMetrics[0] = *p.ProcessHttpStartStopResponseTime(httpStartStopEvent)
	processedMetrics[1] = *p.ProcessHttpStartStopStatusCodeCount(httpStartStopEvent)
	processedMetrics[2] = *p.ProcessHttpStartStopHttpErrorCount(httpStartStopEvent)
	processedMetrics[3] = *p.ProcessHttpStartStopHttpRequestCount(httpStartStopEvent)

	return processedMetrics
}

func (p *HttpStartStopProcessor) ProcessHttpStartStopResponseTime(event *events.HttpStartStop) *metrics.WMetric {
	statPrefix := "http.responsetimes."
	hostname := strings.Replace(strings.Split(event.GetUri(), "/")[0], ".", "_", -1)
	hostport := strings.Replace(hostname, ":", "_",-1)
	stat := statPrefix + hostport

	startTimestamp := event.GetStartTimestamp()
	stopTimestamp := event.GetStopTimestamp()
	durationNanos := stopTimestamp - startTimestamp
	durationMillis := durationNanos / 1000000 // NB: loss of precision here
	metric := metrics.NewIntCounterMetric(stat, int32(durationMillis))

	return metric
}

func (p *HttpStartStopProcessor) ProcessHttpStartStopStatusCodeCount(event *events.HttpStartStop) *metrics.WMetric {
	statPrefix := "http.statuscodes."
	hostname := strings.Replace(strings.Split(event.GetUri(), "/")[0], ".", "_", -1)
	hostport := strings.Replace(hostname, ":", "_",-1)
	stat := statPrefix + hostport + "." + strconv.Itoa(int(event.GetStatusCode()))

	metric := metrics.NewPerintervalCounterMetric(stat, int32(isPeer(event)))

	return metric
}

func (p *HttpStartStopProcessor) ProcessHttpStartStopHttpErrorCount(event *events.HttpStartStop) *metrics.WMetric {
	var incrementValue int64

	statPrefix := "http.errors."
	hostname := strings.Replace(strings.Split(event.GetUri(), "/")[0], ".", "_", -1)
	hostport := strings.Replace(hostname, ":", "_",-1)
	stat := statPrefix + hostport

	if 299 < event.GetStatusCode() && 1 == isPeer(event) {
		incrementValue = 1
	} else {
		incrementValue = 0
	}

	metric := metrics.NewPerintervalCounterMetric(stat, int32(incrementValue))

	return metric
}

func (p *HttpStartStopProcessor) ProcessHttpStartStopHttpRequestCount(event *events.HttpStartStop) *metrics.WMetric {
	statPrefix := "http.requests."
	hostname := strings.Replace(strings.Split(event.GetUri(), "/")[0], ".", "_", -1)
	hostport := strings.Replace(hostname, ":", "_",-1)
	stat := statPrefix + hostport
	metric := metrics.NewPerintervalCounterMetric(stat, int32(isPeer(event)))

	return metric
}

func isPeer(event *events.HttpStartStop) int64 {
	if event.GetPeerType() == events.PeerType_Client {
		return 1
	} else {
		return 0
	}
}
