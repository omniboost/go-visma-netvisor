package netvisor

import (
	"net/http"
	"net/url"
)

func (c *Client) NewAccountingledgerRequest() AccountingledgerRequest {
	return AccountingledgerRequest{
		client:      c,
		queryParams: c.NewAccountingledgerQueryParams(),
		pathParams:  c.NewAccountingledgerPathParams(),
		method:      http.MethodGet,
		headers:     http.Header{},
		requestBody: c.NewAccountingledgerRequestBody(),
	}
}

type AccountingledgerRequest struct {
	client      *Client
	queryParams *AccountingledgerQueryParams
	pathParams  *AccountingledgerPathParams
	method      string
	headers     http.Header
	requestBody AccountingledgerRequestBody
}

func (c *Client) NewAccountingledgerQueryParams() *AccountingledgerQueryParams {
	return &AccountingledgerQueryParams{}
}

type AccountingledgerQueryParams struct {
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

func (p AccountingledgerQueryParams) ToURLValues() (url.Values, error) {
	encoder := NewSchemaEncoder()
	params := url.Values{}

	err := encoder.Encode(p, params)
	if err != nil {
		return params, err
	}

	return params, nil
}

func (r *AccountingledgerRequest) QueryParams() *AccountingledgerQueryParams {
	return r.queryParams
}

func (c *Client) NewAccountingledgerPathParams() *AccountingledgerPathParams {
	return &AccountingledgerPathParams{}
}

type AccountingledgerPathParams struct {
}

func (p *AccountingledgerPathParams) Params() map[string]string {
	return map[string]string{}
}

func (r *AccountingledgerRequest) PathParams() *AccountingledgerPathParams {
	return r.pathParams
}

func (r *AccountingledgerRequest) SetMethod(method string) {
	r.method = method
}

func (r *AccountingledgerRequest) Method() string {
	return r.method
}

func (s *Client) NewAccountingledgerRequestBody() AccountingledgerRequestBody {
	return AccountingledgerRequestBody{}
}

type AccountingledgerRequestBody struct{}

func (r *AccountingledgerRequest) RequestBody() *AccountingledgerRequestBody {
	return &r.requestBody
}

func (r *AccountingledgerRequest) SetRequestBody(body AccountingledgerRequestBody) {
	r.requestBody = body
}

func (r *AccountingledgerRequest) NewResponseBody() *AccountingledgerResponseBody {
	return &AccountingledgerResponseBody{}
}

type AccountingledgerResponseBody struct {
	ResponseStatus ResponseStatus
	Vouchers       Vouchers `xml:"Vouchers>Voucher"`
}

func (r *AccountingledgerRequest) URL() url.URL {
	return r.client.GetEndpointURL("accountingledger.nv", r.PathParams())
}

func (r *AccountingledgerRequest) Do() (AccountingledgerResponseBody, error) {
	u := r.URL()

	// Process query parameters
	values, err := r.QueryParams().ToURLValues()
	if err != nil {
		return *r.NewResponseBody(), err
	}

	u = AddQueryParamsToURL(u, values)

	// Create http request
	req, err := r.client.NewRequest(nil, r.Method(), u, nil)
	if err != nil {
		return *r.NewResponseBody(), err
	}

	responseBody := r.NewResponseBody()
	_, err = r.client.Do(req, responseBody)
	return *responseBody, err
}
