package processors

import (
	"strconv"

	"github.com/tmcgaughey/epagent-nozzle/metrics"
	"github.com/cloudfoundry/noaa/events"
)

type ContainerMetricProcessor struct{}

func NewContainerMetricProcessor() *ContainerMetricProcessor {
	return &ContainerMetricProcessor{}
}

func (p *ContainerMetricProcessor) Process(e *events.Envelope) []metrics.WMetric {
	processedMetrics := make([]metrics.WMetric, 3)
	containerMetricEvent := e.GetContainerMetric()

	processedMetrics[0] = *p.ProcessContainerMetricCPU(containerMetricEvent)
	processedMetrics[1] = *p.ProcessContainerMetricMemory(containerMetricEvent)
	processedMetrics[2] = *p.ProcessContainerMetricDisk(containerMetricEvent)

	return processedMetrics
}

func (p *ContainerMetricProcessor) ProcessContainerMetricCPU(e *events.ContainerMetric) *metrics.WMetric {
	appID := e.GetApplicationId()
	instanceIndex := strconv.Itoa(int(e.GetInstanceIndex()))

	stat := "apps." + appID + ".cpu." + instanceIndex
	metric := metrics.NewLongCounterMetric(stat, int64(e.GetCpuPercentage()))

	return metric
}

func (p *ContainerMetricProcessor) ProcessContainerMetricMemory(e *events.ContainerMetric) *metrics.WMetric {
	appID := e.GetApplicationId()
	instanceIndex := strconv.Itoa(int(e.GetInstanceIndex()))

	stat := "apps." + appID + ".memoryBytes." + instanceIndex
	metric := metrics.NewLongCounterMetric(stat, int64(e.GetMemoryBytes()))
	
	return metric
}

func (p *ContainerMetricProcessor) ProcessContainerMetricDisk(e *events.ContainerMetric) *metrics.WMetric {
	appID := e.GetApplicationId()
	instanceIndex := strconv.Itoa(int(e.GetInstanceIndex()))

	stat := "apps." + appID + ".diskBytes." + instanceIndex
	metric := metrics.NewLongCounterMetric(stat, int64(e.GetDiskBytes()))

	return metric
}
