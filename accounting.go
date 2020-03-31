package netvisor

import (
	"encoding/xml"
	"net/http"
	"net/url"
)

func (c *Client) NewAccountingRequest() AccountingRequest {
	return AccountingRequest{
		client:      c,
		queryParams: c.NewAccountingQueryParams(),
		pathParams:  c.NewAccountingPathParams(),
		method:      http.MethodPost,
		headers:     http.Header{},
		requestBody: c.NewAccountingRequestBody(),
	}
}

type AccountingRequest struct {
	client      *Client
	queryParams *AccountingQueryParams
	pathParams  *AccountingPathParams
	method      string
	headers     http.Header
	requestBody AccountingRequestBody
}

func (c *Client) NewAccountingQueryParams() *AccountingQueryParams {
	return &AccountingQueryParams{}
}

type AccountingQueryParams struct {
	// Finds vouchers that are added after given date (inclusive)
	StartDate DateTime `schema:"startdate"`
	// Finds vouchers that are added before given date(inclusive)
	EndDate DateTime `schema:"enddate"`
	// Finds vouchers which account number is equal or greater than given value
	AccountNumberStart int `schema:"accountnumberstart,omitemptyt"`
	// Finds vouchers which account number is equal or less than given value
	AccountNumberEnd int `schema:"accountnumberend,omitempty"`
	// Finds vouchers that are modified or added after given date
	ChangedSince DateTime `schema:"changedsince,omitempty"`
	// Finds vouchers that are modified after given date
	LastModifiedStart DateTime `schema:"lastmodifiedstart,omitempty"`
	// Finds vouchers that are modified before given date
	LastModifiedEnd DateTime `schema:"lastmodifiedend,omitempty"`
	// Returns the creator of the voucher (eg. user or system) 	1
	ShowGenerator bool `schema:"showgenerator,omitempty"`
	// Finds all, only valid or invalidated and deleted vouchers. Status -attribute shows which status voucher is: "valid", "invalidated" or "deleted" 	1 = All
	// 2 = Valid
	// 3 = Invalidated and deleted}
	VoucherStatus int `schema:"voucherstatus,omitempty"`
}

func (p AccountingQueryParams) ToURLValues() (url.Values, error) {
	encoder := NewSchemaEncoder()
	params := url.Values{}

	err := encoder.Encode(p, params)
	if err != nil {
		return params, err
	}

	return params, nil
}

func (r *AccountingRequest) QueryParams() *AccountingQueryParams {
	return r.queryParams
}

func (c *Client) NewAccountingPathParams() *AccountingPathParams {
	return &AccountingPathParams{}
}

type AccountingPathParams struct {
}

func (p *AccountingPathParams) Params() map[string]string {
	return map[string]string{}
}

func (r *AccountingRequest) PathParams() *AccountingPathParams {
	return r.pathParams
}

func (r *AccountingRequest) SetMethod(method string) {
	r.method = method
}

func (r *AccountingRequest) Method() string {
	return r.method
}

func (s *Client) NewAccountingRequestBody() AccountingRequestBody {
	return AccountingRequestBody{}
}

type AccountingRequestBody struct {
	XMLName xml.Name   `xml:"Root"`
	Voucher NewVoucher `xml:"Voucher"`
}

func (r *AccountingRequest) RequestBody() *AccountingRequestBody {
	return &r.requestBody
}

func (r *AccountingRequest) SetRequestBody(body AccountingRequestBody) {
	r.requestBody = body
}

func (r *AccountingRequest) NewResponseBody() *AccountingResponseBody {
	return &AccountingResponseBody{}
}

type AccountingResponseBody struct {
	ResponseStatus ResponseStatus
	Vouchers       Vouchers `xml:"Vouchers>Voucher"`
}

func (r *AccountingRequest) URL() url.URL {
	return r.client.GetEndpointURL("accounting.nv", r.PathParams())
}

func (r *AccountingRequest) Do() (AccountingResponseBody, error) {
	// Process query parameters
	u := r.URL()
	values, err := r.QueryParams().ToURLValues()
	if err != nil {
		return *r.NewResponseBody(), err
	}
	u = AddQueryParamsToURL(u, values)

	// Create http request
	req, err := r.client.NewRequest(nil, r.Method(), u, r.RequestBody())
	if err != nil {
		return *r.NewResponseBody(), err
	}

	responseBody := r.NewResponseBody()
	_, err = r.client.Do(req, responseBody)
	return *responseBody, err
}
