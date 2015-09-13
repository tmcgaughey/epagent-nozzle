package metrics

//CA APM (Wily) Metrics

import (
	"strconv"
	)


type WMetric struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Value string `json:"value"`
}

type MetricFeed struct {
	Metrics []WMetric `json:"metrics"`
}

//Useful for metric such as miles per hour or errors per interval.
//Resets to zero at each new interval.
func NewPerintervalCounterMetric(stat string, value int32) *WMetric {
	return &WMetric{ Type: "PerintervalCounter", Name: stat, Value: strconv.FormatInt(int64(value),10) }
}

//This is an int value that can rise or fall 
//Useful for tally metrics, such as msgs in queue.
//The value does not change until a new value is reported
func NewIntCounterMetric(stat string, value int32) *WMetric {
	return &WMetric{ Type: "IntCounter", Name: stat, Value: strconv.FormatInt(int64(value),10) }
}

//This is an int value that is averaged over time.
//Useful for response times, such as average time in seconds.
//The value reports all the applicable metrics, such as in a loop and automatically performs the calculation at the end of the interval.
func NewIntAverageMetric(stat string, value int32) *WMetric {
	return &WMetric{ Type: "IntAverage", Name: stat, Value: strconv.FormatInt(int64(value),10) }
}

//The value is a per second rate
//These metrics are aggregated over time from an average of the value
func NewIntRateMetric(stat string, value int32) *WMetric {
	return &WMetric{ Type: "IntRate", Name: stat, Value: strconv.FormatInt(int64(value),10) }
}

//This is a long value that can rise or fall 
//Useful for tally metrics, such as msgs in queue.
//The value does not change until a new value is reported
func NewLongCounterMetric(stat string, value int64) *WMetric {
	return &WMetric{ Type: "LongCounter", Name: stat, Value: strconv.FormatInt(value,10) }
}

//This is a long value that is averaged over time.
//Useful for response times, such as average time in seconds.
//The value reports all the applicable metrics, such as in a loop and automatically performs the calculation at the end of the interval.
func NewLongAverageMetric(stat string, value int64) *WMetric {
	return &WMetric{ Type: "LongAverage", Name: stat, Value: strconv.FormatInt(value,10) }
}

//current latest string value (not stored historically)
func NewStringEventMetric(stat string, value string) *WMetric {
	return &WMetric{ Type: "StringEvent", Name: stat, Value: value }
}

//A type which generates successively increasing timestamps
func NewTimestampMetric(stat string, value int64) *WMetric {
	return &WMetric{ Type: "Timestamp", Name: stat, Value: strconv.FormatInt(value,10) }
}

