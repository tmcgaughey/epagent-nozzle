package uaatokenfetcher

// Stolen from https://github.com/cloudfoundry-incubator/datadog-firehose-nozzle
// Just copied to simplify godeps

import (
	"github.com/cloudfoundry-incubator/uaago"
	"log"
)

type UAATokenFetcher struct {
	UaaUrl                string
	Username              string
	Password              string
	InsecureSSLSkipVerify bool
}

func (uaa *UAATokenFetcher) FetchAuthToken() string {
	uaaClient, err := uaago.NewClient(uaa.UaaUrl)
	if err != nil {
		log.Fatalf("Error creating uaa client: %s", err.Error())
	}

	var authToken string
	authToken, err = uaaClient.GetAuthToken(uaa.Username, uaa.Password, uaa.InsecureSSLSkipVerify)
	if err != nil {
		log.Fatalf("Error getting oauth token: %s. Please check your username and password.", err.Error())
	}
	return authToken
}
