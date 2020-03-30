package netvisor

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/gofrs/uuid"
)

const (
	libraryVersion    = "0.0.1"
	userAgent         = "go-visma-netvisor/" + libraryVersion
	mediaType         = "application/xml"
	charset           = "utf-8"
	sender            = "omniboost.io"
	interfaceLanguage = "EN"
)

var (
	baseURL = url.URL{
		Scheme: "https",
		Host:   "isvapi.netvisor.fi",
		Path:   "/",
	}
)

// NewClient returns a new Visma Netvisor client
func NewClient(httpClient *http.Client, customerID, partnerID, organisationID, partnerKey, privateKey string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	client := &Client{}

	client.SetHTTPClient(httpClient)
	client.SetCustomerID(customerID)
	client.SetPartnerID(partnerID)
	client.SetOrganisationID(organisationID)
	client.SetPartnerKey(partnerKey)
	client.SetPrivateKey(privateKey)
	client.SetBaseURL(baseURL)
	client.SetSender(sender)
	client.SetInterfaceLanguage(interfaceLanguage)
	client.SetDebug(false)
	client.SetUserAgent(userAgent)
	client.SetMediaType(mediaType)
	client.SetCharset(charset)

	return client
}

// Client manages communication with Visma Netvisor
type Client struct {
	// HTTP client used to communicate with the Client.
	http *http.Client

	debug   bool
	baseURL url.URL

	// credentials
	customerID     string
	partnerID      string
	organisationID string
	privateKey     string
	partnerKey     string

	// headers
	sender            string
	interfaceLanguage string

	// User agent for client
	userAgent string

	mediaType string
	charset   string

	// Optional function called after every successful request made to the DO Clients
	beforeRequestDo    BeforeRequestDoCallback
	onRequestCompleted RequestCompletionCallback
}

type BeforeRequestDoCallback func(*http.Client, *http.Request, interface{})

// RequestCompletionCallback defines the type of the request callback function
type RequestCompletionCallback func(*http.Request, *http.Response)

func (c *Client) SetHTTPClient(client *http.Client) {
	c.http = client
}

func (c *Client) Debug() bool {
	return c.debug
}

func (c *Client) SetDebug(debug bool) {
	c.debug = debug
}

func (c *Client) CustomerID() string {
	return c.customerID
}

func (c *Client) SetCustomerID(customerID string) {
	c.customerID = customerID
}

func (c *Client) PartnerID() string {
	return c.partnerID
}

func (c *Client) SetPartnerID(partnerID string) {
	c.partnerID = partnerID
}

func (c *Client) OrganisationID() string {
	return c.organisationID
}

func (c *Client) SetOrganisationID(organisationID string) {
	c.organisationID = organisationID
}

func (c *Client) PartnerKey() string {
	return c.partnerKey
}

func (c *Client) SetPartnerKey(partnerKey string) {
	c.partnerKey = partnerKey
}

func (c *Client) PrivateKey() string {
	return c.privateKey
}

func (c *Client) SetPrivateKey(privateKey string) {
	c.privateKey = privateKey
}

func (c *Client) BaseURL() url.URL {
	return c.baseURL
}

func (c *Client) SetBaseURL(baseURL url.URL) {
	c.baseURL = baseURL
}

func (c *Client) Sender() string {
	return c.sender
}

func (c *Client) SetSender(sender string) {
	c.sender = sender
}

func (c *Client) InterfaceLanguage() string {
	return c.interfaceLanguage
}

func (c *Client) SetInterfaceLanguage(interfaceLanguage string) {
	c.interfaceLanguage = interfaceLanguage
}

func (c *Client) SetMediaType(mediaType string) {
	c.mediaType = mediaType
}

func (c *Client) MediaType() string {
	return mediaType
}

func (c *Client) SetCharset(charset string) {
	c.charset = charset
}

func (c *Client) Charset() string {
	return charset
}

func (c *Client) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
}

func (c *Client) UserAgent() string {
	return userAgent
}

func (c *Client) SetBeforeRequestDo(fun BeforeRequestDoCallback) {
	c.beforeRequestDo = fun
}

func (c *Client) GetEndpointURL(path string, pathParams PathParams) url.URL {
	clientURL := c.BaseURL()
	clientURL.Path = clientURL.Path + path

	tmpl, err := template.New("endpoint_url").Parse(clientURL.Path)
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	params := pathParams.Params()
	err = tmpl.Execute(buf, params)
	if err != nil {
		log.Fatal(err)
	}

	clientURL.Path = buf.String()
	clientURL.RawPath = buf.String()
	return clientURL
}

func (c *Client) NewRequest(ctx context.Context, method string, URL url.URL, body interface{}) (*http.Request, error) {
	// convert body struct to xml
	buf := new(bytes.Buffer)
	if body != nil {
		err := xml.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	// create new http request
	req, err := http.NewRequest(method, URL.String(), buf)
	// req.Host = URL.Hostname()
	if err != nil {
		return nil, err
	}

	// optionally pass along context
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	// set other headers
	req.Header.Add("X-Netvisor-Authentication-Sender", c.Sender())
	req.Header.Add("X-Netvisor-Authentication-CustomerId", c.CustomerID())
	req.Header.Add("X-Netvisor-Authentication-PartnerId", c.PartnerID())
	req.Header.Add("X-Netvisor-Authentication-Timestamp", c.NewTimestamp())
	req.Header.Add("X-Netvisor-Interface-Language", c.InterfaceLanguage())
	req.Header.Add("X-Netvisor-Organisation-ID", c.OrganisationID())
	req.Header.Add("X-Netvisor-Authentication-TransactionId", c.NewTransactionID())
	req.Header.Add("X-Netvisor-Authentication-MAC", c.NewMAC(req))
	req.Header.Add("X-Netvisor-Authentication-MACHashCalculationAlgorithm", "SHA256")
	req.Header.Add("Content-Type", fmt.Sprintf("%s; charset=%s", c.MediaType(), c.Charset()))
	req.Header.Add("Accept", c.MediaType())
	req.Header.Add("User-Agent", c.UserAgent())

	return req, nil
}

// Do sends an Client request and returns the Client response. The Client response is xml decoded and stored in the value
// pointed to by v, or returned as an error if an Client error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(req *http.Request, responseBody interface{}) (*http.Response, error) {
	if c.beforeRequestDo != nil {
		c.beforeRequestDo(c.http, req, responseBody)
	}

	if c.debug == true {
		dump, _ := httputil.DumpRequestOut(req, true)
		log.Println(string(dump))
	}

	httpResp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, httpResp)
	}

	// close body io.Reader
	defer func() {
		if rerr := httpResp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	if c.debug == true {
		dump, _ := httputil.DumpResponse(httpResp, true)
		log.Println(string(dump))
	}

	// check if the response isn't an error
	err = CheckResponse(httpResp)
	if err != nil {
		return httpResp, err
	}

	// check the provided interface parameter
	if httpResp == nil {
		return httpResp, nil
	}

	if responseBody == nil {
		return httpResp, err
	}

	// interface implements io.Writer: write Body to it
	// if w, ok := response.Envelope.(io.Writer); ok {
	// 	_, err := io.Copy(w, httpResp.Body)
	// 	return httpResp, err
	// }

	// try to decode body into interface parameter
	// w := &Wrapper{}
	dec := xml.NewDecoder(httpResp.Body)
	err = dec.Decode(responseBody)
	if err != nil && err != io.EOF {
		// create a simple error response
		errorResponse := &ErrorResponse{Response: httpResp}
		errorResponse.Errors = append(errorResponse.Errors, err)
		return httpResp, errorResponse
	}

	// err = xml.Unmarshal(w.D.Results, responseBody)
	// if err != nil && err != io.EOF {
	// 	// @TODO: fix this
	// 	log.Fatal(err)
	// }

	return httpResp, nil
}

func (c *Client) NewTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05") + ".000"
}

func (c *Client) NewTransactionID() string {
	return uuid.Must(uuid.NewV4()).String()
}

func (c *Client) NewMAC(req *http.Request) string {
	pieces := []string{
		req.URL.String(),
		req.Header.Get("X-Netvisor-Authentication-Sender"),
		req.Header.Get("X-Netvisor-Authentication-CustomerId"),
		req.Header.Get("X-Netvisor-Authentication-Timestamp"),
		req.Header.Get("X-Netvisor-Interface-Language"),
		req.Header.Get("X-Netvisor-Organisation-ID"),
		req.Header.Get("X-Netvisor-Authentication-TransactionId"),
		c.PrivateKey(),
		c.PartnerKey(),
	}

	concat := strings.Join(pieces, "&")
	return fmt.Sprintf("%x", sha256.Sum256([]byte(concat)))

}

// CheckResponse checks the Client response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range. Client error responses are expected to have either no response
// body, or a xml response body that maps to ErrorResponse. Any other response
// body will be silently ignored.
func CheckResponse(r *http.Response) error {
	errorResponse := &ErrorResponse{Response: r}

	// Don't check content-lenght: a created response, for example, has no body
	// if r.Header.Get("Content-Length") == "0" {
	// 	errorResponse.Errors.Message = r.Status
	// 	return errorResponse
	// }

	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	err := checkContentType(r)
	if err != nil {
		errorResponse.Errors = append(errorResponse.Errors, errors.New(r.Status))
		return errorResponse
	}

	// read data and copy it back
	data, err := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(data))
	if err != nil {
		return errorResponse
	}

	if len(data) == 0 {
		return errorResponse
	}

	// convert xml to struct
	err = xml.Unmarshal(data, errorResponse)
	if err != nil {
		errorResponse.Errors = append(errorResponse.Errors, err)
		return errorResponse
	}

	return errorResponse
}

type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response `xml:"-"`

	Errors []error
}

func (e ErrorResponse) Error() string {
	return "TEST"
}

func checkContentType(response *http.Response) error {
	header := response.Header.Get("Content-Type")
	contentType := strings.Split(header, ";")[0]
	if contentType != mediaType {
		return fmt.Errorf("Expected Content-Type \"%s\", got \"%s\"", mediaType, contentType)
	}

	return nil
}

type PathParams interface {
	Params() map[string]string
}
