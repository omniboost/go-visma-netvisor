package netvisor

type ResponseStatus struct {
	Status    string
	TimeStamp DateTime
}

type Vouchers []Voucher

type Voucher struct {
	Status                  string   `xml:"Status,attr"`
	NetvisorKey             string   `xml:"NetvisorKey"`
	VoucherDate             DateTime `xml:"VoucherDate"`
	VoucherNumber           string   `xml:"VoucherNumber"`
	VoucherDescription      string   `xml:"VoucherDescription"`
	VoucherClass            string   `xml:"VoucherClass"`
	LinkedSourceNetvisorKey struct {
		Text string `xml:",chardata"`
		Type string `xml:"type,attr"`
	} `xml:"LinkedSourceNetvisorKey"`
	VoucherNetvisorURI string       `xml:"VoucherNetvisorURI"`
	VoucherLine        VoucherLines `xml:"VoucherLine"`
}

type VoucherLines []VoucherLine

type VoucherLine struct {
	NetvisorKey   int     `xml:"NetvisorKey"`
	LineSum       float64 `xml:"LineSum"`
	Description   string  `xml:"Description"`
	AccountNumber string  `xml:"AccountNumber"`
	VatPercent    float64 `xml:"VatPercent"`
	VatCode       string  `xml:"VatCode"`
}
