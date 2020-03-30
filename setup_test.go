package netvisor_test

import (
	"log"
	"net/url"
	"os"
	"testing"

	netvisor "github.com/omniboost/go-visma-netvisor"
	// ntlm "github.com/vadimi/go-http-ntlm"
)

var (
	client *netvisor.Client
)

func TestMain(m *testing.M) {
	var (
		baseURL *url.URL
		err     error
	)

	baseURLString := os.Getenv("NETVISOR_BASE_URL")
	if baseURLString != "" {
		baseURL, err = url.Parse(baseURLString)
		if err != nil {
			log.Fatal(err)
		}
	}

	customerID := os.Getenv("NETVISOR_CUSTOMER_ID")
	partnerID := os.Getenv("NETVISOR_PARTNER_ID")
	organisationID := os.Getenv("NETVISOR_ORGANISATION_ID")
	privateKey := os.Getenv("NETVISOR_PRIVATE_KEY")
	partnerKey := os.Getenv("NETVISOR_PARTNER_KEY")
	debug := os.Getenv("NETVISOR_DEBUG")

	client = netvisor.NewClient(nil, customerID, partnerID, organisationID, partnerKey, privateKey)
	if debug != "" {
		client.SetDebug(true)
	}
	if baseURL != nil {
		client.SetBaseURL(*baseURL)
	}
	m.Run()
}
