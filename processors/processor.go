package processors

import (
	"github.com/tmcgaughey/epagent-nozzle/metrics"
	"github.com/cloudfoundry/noaa/events"
)

type Processor interface {
	Process(e *events.Envelope) []metrics.WMetric
}
