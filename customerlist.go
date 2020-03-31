package netvisor

import (
	"net/http"
	"net/url"
)

func (c *Client) NewCustomerlistRequest() CustomerlistRequest {
	return CustomerlistRequest{
		client:      c,
		queryParams: c.NewCustomerlistQueryParams(),
		pathParams:  c.NewCustomerlistPathParams(),
		method:      http.MethodGet,
		headers:     http.Header{},
		requestBody: c.NewCustomerlistRequestBody(),
	}
}

type CustomerlistRequest struct {
	client      *Client
	queryParams *CustomerlistQueryParams
	pathParams  *CustomerlistPathParams
	method      string
	headers     http.Header
	requestBody CustomerlistRequestBody
}

func (c *Client) NewCustomerlistQueryParams() *CustomerlistQueryParams {
	return &CustomerlistQueryParams{}
}

type CustomerlistQueryParams struct {
	// Filters result list with given keyword. Match is searched from following
	// fields: Name, Customer Code, Organization identifier, CoName
	Keyword string `schema:"keyword,omitempty"`
	// Filters result to contain only customers having change after given date, date in format YYYY-MM-DD
	ChangedSince DateTime `schema:"changedsince,omitempty"`
}

func (p CustomerlistQueryParams) ToURLValues() (url.Values, error) {
	encoder := NewSchemaEncoder()
	params := url.Values{}

	err := encoder.Encode(p, params)
	if err != nil {
		return params, err
	}

	return params, nil
}

func (r *CustomerlistRequest) QueryParams() *CustomerlistQueryParams {
	return r.queryParams
}

func (c *Client) NewCustomerlistPathParams() *CustomerlistPathParams {
	return &CustomerlistPathParams{}
}

type CustomerlistPathParams struct {
}

func (p *CustomerlistPathParams) Params() map[string]string {
	return map[string]string{}
}

func (r *CustomerlistRequest) PathParams() *CustomerlistPathParams {
	return r.pathParams
}

func (r *CustomerlistRequest) SetMethod(method string) {
	r.method = method
}

func (r *CustomerlistRequest) Method() string {
	return r.method
}

func (s *Client) NewCustomerlistRequestBody() CustomerlistRequestBody {
	return CustomerlistRequestBody{}
}

type CustomerlistRequestBody struct{}

func (r *CustomerlistRequest) RequestBody() *CustomerlistRequestBody {
	return &r.requestBody
}

func (r *CustomerlistRequest) SetRequestBody(body CustomerlistRequestBody) {
	r.requestBody = body
}

func (r *CustomerlistRequest) NewResponseBody() *CustomerlistResponseBody {
	return &CustomerlistResponseBody{}
}

type CustomerlistResponseBody struct {
	ResponseStatus ResponseStatus
	Customers      Customers `xml:"Customerlist>Customer"`
}

func (r *CustomerlistRequest) URL() url.URL {
	return r.client.GetEndpointURL("customerlist.nv", r.PathParams())
}

func (r *CustomerlistRequest) Do() (CustomerlistResponseBody, error) {
	// Process query parameters
	u := r.URL()
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
