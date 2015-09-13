package epagentclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
//	"time"

	"strings"
	"strconv"
	//"errors"
	"github.com/cloudfoundry/noaa/events" //"github.com/cloudfoundry/sonde-go/events"
	"log"
	"io/ioutil"
    "github.com/tmcgaughey/epagent-nozzle/metrics"
    "github.com/tmcgaughey/epagent-nozzle/processors"
)

const slowConsumer = "slowConsumerAlert" 

type Client struct {
	apiURL                string
	metricPoints          map[string]metrics.WMetric
//	prefix                string
//	deployment            string
//	ip                    string
	totalMessagesReceived int64
	totalMetricsSent      int64
	httpStartStopProcessor *processors.HttpStartStopProcessor
	valueMetricProcessor  *processors.ValueMetricProcessor
	containerMetricProcessor *processors.ContainerMetricProcessor
	heartbeatProcessor *processors.HeartbeatProcessor
	counterProcessor *processors.CounterProcessor
}


func New(apiURL string) *Client { //, apiKey string, prefix string, deployment string, ip string) *Client {
	return &Client{
		apiURL:       apiURL,
//		apiKey:       apiKey,
		metricPoints: make(map[string]metrics.WMetric),
//		prefix:       prefix,
//		deployment:   deployment,
//		ip:           ip,
		httpStartStopProcessor : processors.NewHttpStartStopProcessor(),
		valueMetricProcessor : processors.NewValueMetricProcessor(),
		containerMetricProcessor : processors.NewContainerMetricProcessor(),
		heartbeatProcessor : processors.NewHeartbeatProcessor(),
		counterProcessor : processors.NewCounterProcessor(),
	}
}

func (c *Client) AlertSlowConsumerError() {
	c.addInternalMetric(slowConsumer, int64(1))
}

func (c *Client) AddMetric(envelope *events.Envelope) {
	c.totalMessagesReceived++
	eventType := envelope.GetEventType()
	
	processedMetrics := []metrics.WMetric{}

	// epagent-nozzle can handle CounterEvent, ContainerMetric, Heartbeat,
		// HttpStartStop and ValueMetric events
		switch eventType {
		case events.Envelope_ContainerMetric:
			processedMetrics = c.containerMetricProcessor.Process(envelope)
		case events.Envelope_CounterEvent:
			processedMetrics = c.counterProcessor.Process(envelope)
		case events.Envelope_Heartbeat:
			processedMetrics = c.heartbeatProcessor.Process(envelope)
		case events.Envelope_HttpStartStop:
			processedMetrics = c.httpStartStopProcessor.Process(envelope)
		case events.Envelope_ValueMetric:
			processedMetrics = c.valueMetricProcessor.Process(envelope)
		default:
			// do nothing
		}

			if len(processedMetrics) > 0 {
				for _, metric := range processedMetrics {
					if metric.Type == "PerintervalCounter" {
						_,ok := c.metricPoints[metric.Name]
						if ok {
							//log.Println("incrementing counter")
							oldmetric := c.metricPoints[metric.Name]
							oldint,_ := strconv.ParseInt(oldmetric.Value,10,64)
							newint,_ := strconv.ParseInt(metric.Value,10,64)
							newint = oldint + newint
							
							metric.Value = strconv.FormatInt(newint,10)
						}
					}
					c.metricPoints[metric.Name] = metric
					//log.Printf("metric: %s:::%s\n", metric.Name, metric.Value)
				}
			}

}

func (c *Client) PostMetrics() error {

	c.populateInternalMetrics()
	numMetrics := len(c.metricPoints)
	log.Printf("Posting %d metrics", numMetrics)

	seriesBytes, metricsCount := c.formatMetrics()
    
    //fmt.Println(string(seriesBytes[:]))

    //log.Print(c.apiURL)
	req, err := http.NewRequest("POST", c.apiURL, bytes.NewBuffer(seriesBytes))
	req.Header.Set("Content-Type", "application/json")
	
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

    c.totalMessagesReceived = 0
	c.totalMetricsSent = metricsCount
	c.metricPoints = make(map[string]metrics.WMetric)

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		body,_ :=ioutil.ReadAll(resp.Body)
		log.Printf("Response for %s code:\n%s\n", resp.Status, string(body))
		return fmt.Errorf("EPAgent returned HTTP response: %s", resp.Status)
	}
	return nil
}

func (c *Client) populateInternalMetrics() {
	c.addInternalMetric("MessagesReceived", c.totalMessagesReceived)
	c.addInternalMetric("MetricsSent", c.totalMetricsSent)

	if !c.containsSlowConsumerAlert() {
		c.addInternalMetric(slowConsumer, int64(0))
	}
}

func (c *Client) containsSlowConsumerAlert() bool {
	_, ok := c.metricPoints[slowConsumer]
	return ok
}

func (c *Client) formatMetrics() ([]byte, int64) {
	//log.Println("formatting metrics");
	wmetrics := []metrics.WMetric{}
	for _, mVal := range c.metricPoints {
		mVal.Name = "Firehose." + mVal.Name
		str := strings.Replace(mVal.Name, ".","|",-1)
		idx := strings.LastIndexAny(str,"|")
		substring := str[idx:len(str)]
		revsubstring := ":"+substring[1:len(substring)]
		mVal.Name = strings.Replace(str, substring, revsubstring, -1)
		//log.Printf("metric: %s\tvalue: %s", mVal.Name, mVal.Value)
		wmetrics = append(wmetrics, mVal)
	}

    metriclist := new (metrics.MetricFeed)
    metriclist.Metrics = wmetrics
	encodedMetric, _ := json.Marshal(metriclist)
	return encodedMetric, int64(len(wmetrics))
	
}

func (c *Client) addInternalMetric(name string, value int64) {
	
	c.metricPoints[name] = *metrics.NewPerintervalCounterMetric(name, int32(value))
}

