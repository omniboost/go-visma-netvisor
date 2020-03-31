package netvisor_test

import (
	"encoding/json"
	"log"
	"testing"

	_ "github.com/joho/godotenv/autoload"
)

func TestAccounting(t *testing.T) {
	req := client.NewAccountingRequest()

	rb := `{
		"Voucher": {
			"CalculationMode": "net",
			"VoucherDate": {
				"Format": "ansi",
				"Text": "2017-09-11"
			},
			"Description": "test"
		}
	}`
	err := json.Unmarshal([]byte(rb), req.RequestBody())
	if err != nil {
		t.Error(err)
	}

	resp, err := req.Do()
	if err != nil {
		t.Error(err)
	}

	b, _ := json.MarshalIndent(resp, "", "  ")
	log.Println(string(b))
}
