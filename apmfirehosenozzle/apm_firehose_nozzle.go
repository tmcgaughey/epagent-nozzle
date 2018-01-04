package apmfirehosenozzle

import (
	"crypto/tls"
	"strings"
	"time"

	"log"

	"github.com/cloudfoundry/noaa/consumer"
	"github.com/cloudfoundry/sonde-go/events"

	"github.com/tmcgaughey/epagent-nozzle/epagentclient"
	"github.com/tmcgaughey/epagent-nozzle/nozzleconfig"
)

type APMFirehoseNozzle struct {
	config           *nozzleconfig.NozzleConfig
	errs             <-chan error
	messages         <-chan *events.Envelope
	authTokenFetcher AuthTokenFetcher
	consumer         *consumer.Consumer
	client           *epagentclient.Client
}

type AuthTokenFetcher interface {
	FetchAuthToken() string
}

func NewAPMFirehoseNozzle(config *nozzleconfig.NozzleConfig, tokenFetcher AuthTokenFetcher) *APMFirehoseNozzle {
	return &APMFirehoseNozzle{
		config:           config,
		errs:             make(chan error),
		messages:         make(chan *events.Envelope),
		authTokenFetcher: tokenFetcher,
	}
}

func (d *APMFirehoseNozzle) Start() {
	var authToken string

	if !d.config.DisableAccessControl {
		authToken = d.authTokenFetcher.FetchAuthToken()
		authToken = strings.TrimPrefix(authToken, "bearer ")
	}

	log.Print("Starting CA APM Firehose Nozzle...")
	d.createClient()
	d.consumeFirehose(authToken)
	d.postToAPM()
	log.Print("CA APM Firehose Nozzle shutting down...")
}

func (d *APMFirehoseNozzle) createClient() {

	d.client = epagentclient.New(d.config.EPAgentURL)
}

func (d *APMFirehoseNozzle) consumeFirehose(authToken string) {
	d.consumer = consumer.New(d.config.TrafficControllerURL, &tls.Config{InsecureSkipVerify: d.config.InsecureSSLSkipVerify}, nil)

	d.messages, d.errs = d.consumer.Firehose(d.config.FirehoseSubscriptionID, authToken)

}

func (d *APMFirehoseNozzle) postToAPM() {
	ticker := time.NewTicker(time.Duration(d.config.FlushDurationSeconds) * time.Second)
	for {
		select {
		case <-ticker.C:
			//log.Print("time trigger")
			d.postMetrics()
		case envelope := <-d.messages:
			//log.Print("adding metrics")
			d.handleMessage(envelope)
			d.client.AddMetric(envelope)
		case err := <-d.errs:
			//log.Print("handling error")
			d.handleError(err)
			return
		}
	}
}

func (d *APMFirehoseNozzle) postMetrics() {

	err := d.client.PostMetrics()

	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

}

func (d *APMFirehoseNozzle) handleError(err error) {

	log.Printf("Closing connection with traffic controller due to %v", err)
	d.consumer.Close()
	d.postMetrics()
}

func (d *APMFirehoseNozzle) handleMessage(envelope *events.Envelope) {

	if envelope.GetEventType() == events.Envelope_CounterEvent && envelope.CounterEvent.GetName() == "TruncatingBuffer.DroppedMessages" && envelope.GetOrigin() == "doppler" {
		log.Printf("We've intercepted an upstream message which indicates that the nozzle or the TrafficController is not keeping up. Please try scaling up the nozzle.")
		d.client.AlertSlowConsumerError()
	}
}
