package netvisor

import (
	"net/http"
	"net/url"

	"github.com/omniboost/go-visma-netvisor/utils"
)

func (c *Client) NewDimensionlistRequest() DimensionlistRequest {
	return DimensionlistRequest{
		client:      c,
		queryParams: c.NewDimensionlistQueryParams(),
		pathParams:  c.NewDimensionlistPathParams(),
		method:      http.MethodGet,
		headers:     http.Header{},
		requestBody: c.NewDimensionlistRequestBody(),
	}
}

type DimensionlistRequest struct {
	client      *Client
	queryParams *DimensionlistQueryParams
	pathParams  *DimensionlistPathParams
	method      string
	headers     http.Header
	requestBody DimensionlistRequestBody
}

func (c *Client) NewDimensionlistQueryParams() *DimensionlistQueryParams {
	return &DimensionlistQueryParams{}
}

type DimensionlistQueryParams struct{}

func (p DimensionlistQueryParams) ToURLValues() (url.Values, error) {
	encoder := utils.NewSchemaEncoder()
	params := url.Values{}

	err := encoder.Encode(p, params)
	if err != nil {
		return params, err
	}

	return params, nil
}

func (r *DimensionlistRequest) QueryParams() *DimensionlistQueryParams {
	return r.queryParams
}

func (c *Client) NewDimensionlistPathParams() *DimensionlistPathParams {
	return &DimensionlistPathParams{}
}

type DimensionlistPathParams struct {
}

func (p *DimensionlistPathParams) Params() map[string]string {
	return map[string]string{}
}

func (r *DimensionlistRequest) PathParams() *DimensionlistPathParams {
	return r.pathParams
}

func (r *DimensionlistRequest) SetMethod(method string) {
	r.method = method
}

func (r *DimensionlistRequest) Method() string {
	return r.method
}

func (s *Client) NewDimensionlistRequestBody() DimensionlistRequestBody {
	return DimensionlistRequestBody{}
}

type DimensionlistRequestBody struct{}

func (r *DimensionlistRequest) RequestBody() *DimensionlistRequestBody {
	return &r.requestBody
}

func (r *DimensionlistRequest) SetRequestBody(body DimensionlistRequestBody) {
	r.requestBody = body
}

func (r *DimensionlistRequest) NewResponseBody() *DimensionlistResponseBody {
	return &DimensionlistResponseBody{}
}

type DimensionlistResponseBody struct {
}

func (r *DimensionlistRequest) URL() url.URL {
	return r.client.GetEndpointURL("dimensionlist.nv", r.PathParams())
}

func (r *DimensionlistRequest) Do() (DimensionlistResponseBody, error) {
	// Create http request
	req, err := r.client.NewRequest(nil, r.Method(), r.URL(), nil)
	if err != nil {
		return *r.NewResponseBody(), err
	}

	// Process query parameters
	err = utils.AddQueryParamsToRequest(r.QueryParams(), req, false)
	if err != nil {
		return *r.NewResponseBody(), err
	}

	responseBody := r.NewResponseBody()
	_, err = r.client.Do(req, responseBody)
	return *responseBody, err
}
