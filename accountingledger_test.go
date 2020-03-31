package netvisor_test

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	_ "github.com/joho/godotenv/autoload"
	netvisor "github.com/omniboost/go-visma-netvisor"
)

func TestAccountingledger(t *testing.T) {
	req := client.NewAccountingledgerRequest()
	req.QueryParams().StartDate = netvisor.DateTime{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}
	req.QueryParams().EndDate = netvisor.DateTime{Time: time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC)}

	resp, err := req.Do()
	if err != nil {
		t.Error(err)
	}

	b, _ := json.MarshalIndent(resp, "", "  ")
	log.Println(string(b))
}
