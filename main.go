package main

// Borrows heavily from https://github.com/cloudfoundry-incubator/datadog-firehose-nozzle and https://github.com/CloudCredo/graphite-nozzle

import (
	"flag"
	"log"

	"github.com/tmcgaughey/epagent-nozzle/apmfirehosenozzle"
	"github.com/tmcgaughey/epagent-nozzle/nozzleconfig"
	"github.com/tmcgaughey/epagent-nozzle/uaatokenfetcher"
)

func main() {
	configFilePath := flag.String("config", "config/caapm-firehose-nozzle.json", "Location of the nozzle config json file")
	flag.Parse()

	config, err := nozzleconfig.Parse(*configFilePath)
	if err != nil {
		log.Fatalf("Error parsing config: %s", err.Error())
	}

	tokenFetcher := &uaatokenfetcher.UAATokenFetcher{
		UaaUrl:                config.UAAURL,
		Username:              config.Username,
		Password:              config.Password,
		InsecureSSLSkipVerify: config.InsecureSSLSkipVerify,
	}
	
	apm_nozzle := apmfirehosenozzle.NewAPMFirehoseNozzle(config, tokenFetcher)
	apm_nozzle.Start()

}