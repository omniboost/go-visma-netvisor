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

type NewVoucher struct {
	CalculationMode string `xml:"CalculationMode"` // net
	VoucherDate     struct {
		Text   string `xml:",chardata"`   // 2017-09-11
		Format string `xml:"format,attr"` // ansi
	} `xml:"VoucherDate"`
	Description        string                `xml:"Description"`  // Testiyritys Oy lasku 1
	VoucherClass       string                `xml:"VoucherClass"` // Myyntilasku
	VoucherLine        NewVoucherLines       `xml:"VoucherLine"`
	VoucherAttachments NewVoucherAttachments `xml:"VoucherAttachments"`
}

type NewVoucherAttachments []NewVoucherAttachment

type NewVoucherAttachment struct {
	MimeType              string `xml:"MimeType"`              // Application/Pdf
	AttachmentDescription string `xml:"AttachmentDescription"` // Testiliite
	FileName              string `xml:"FileName"`              // liite.pdf
	DocumentData          string `xml:"DocumentData"`          // JVBERi0xLjQNJeLjz9MNCjYgM...
}

type NewVoucherLines []NewVoucherLine

type NewVoucherLine struct {
	LineSum struct {
		Text float64 `xml:",chardata"` // -10000.00, -5000, 20000
		Type string  `xml:"type,attr"`
	} `xml:"LineSum"`
	Description   string `xml:"Description"`   // Testiyritys Oy, lasku 1, ...
	AccountNumber string `xml:"AccountNumber"` // 3000, 2939, 1701
	VatPercent    struct {
		Text    string `xml:",chardata"` // 24, 0, 0
		Vatcode string `xml:"vatcode,attr"`
	} `xml:"VatPercent"`
	Dimension struct {
		DimensionName string `xml:"DimensionName"` // Project
		DimensionItem string `xml:"DimensionItem"` // IIHF
	} `xml:"Dimension"`
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

type Customers []Customer

type Customer struct {
	Netvisorkey            string `xml:"Netvisorkey"`
	Name                   string `xml:"Name"`
	Code                   string `xml:"Code"`
	OrganisationIdentifier string `xml:"OrganisationIdentifier"`
	URI                    string `xml:"Uri"`
}
