package netvisor

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	null "gopkg.in/guregu/null.v3"

	"github.com/gorilla/schema"
)

type SchemaMarshaler interface {
	MarshalSchema() string
}

type ToURLValues interface {
	ToURLValues() (url.Values, error)
}

func AddQueryParamsToRequest(requestParams interface{}, req *http.Request, skipEmpty bool) error {
	var err error
	params := url.Values{}

	to, ok := requestParams.(ToURLValues)
	if ok == true {
		params, err = to.ToURLValues()
		if err != nil {
			return err
		}
	} else {
		encoder := NewSchemaEncoder()
		err := encoder.Encode(requestParams, params)
		if err != nil {
			return err
		}
	}

	query := AddQueryParams(req.URL.Query(), params)
	// if skipEmpty && v == "" {
	// 	continue
	// }

	// if skipEmpty && v == "0" {
	// 	continue
	// }

	req.URL.RawQuery = query.Encode()

	// force $ in query parameters
	req.URL.RawQuery = strings.Replace(req.URL.RawQuery, "%24", "$", -1)
	return nil
}

func AddQueryParamsToURL(url url.URL, params url.Values) url.URL {
	query := url.Query()
	for k, vals := range params {
		for _, v := range vals {
			query.Add(k, v)
		}
	}
	url.RawQuery = query.Encode()
	return url
}

func AddQueryParams(params1, params2 url.Values) url.Values {
	for k, vals := range params2 {
		for _, v := range vals {
			params1.Add(k, v)
		}
	}
	return params1
}

func NewSchemaEncoder() *schema.Encoder {
	encoder := schema.NewEncoder()

	// register custom encoders
	encodeSchemaMarshaler := func(v reflect.Value) string {
		marshaler, ok := v.Interface().(SchemaMarshaler)
		if ok == true {
			return marshaler.MarshalSchema()
		}

		stringer, ok := v.Interface().(fmt.Stringer)
		if ok == true {
			return stringer.String()
		}

		return ""
	}

	// encoder.RegisterEncoder(&odata.Expand{}, encodeSchemaMarshaler)
	// encoder.RegisterEncoder(&odata.Filter{}, encodeSchemaMarshaler)
	// encoder.RegisterEncoder(&odata.Select{}, encodeSchemaMarshaler)
	// encoder.RegisterEncoder(&odata.Top{}, encodeSchemaMarshaler)
	// encoder.RegisterEncoder(&odata.OrderBy{}, encodeSchemaMarshaler)
	// encoder.RegisterEncoder(&odata.Skip{}, encodeSchemaMarshaler)
	encoder.RegisterEncoder(DateTime{}, encodeSchemaMarshaler)

	encodeNullFloat := func(v reflect.Value) string {
		nullFloat, _ := v.Interface().(null.Float)
		if nullFloat.IsZero() {
			return ""
		}
		return strconv.FormatFloat(nullFloat.Float64, 'f', 6, 64)
	}

	encodeNullBool := func(v reflect.Value) string {
		nullBool, _ := v.Interface().(null.Bool)
		if nullBool.IsZero() {
			return ""
		}
		return strconv.FormatBool(nullBool.Bool)
	}

	// encoder.RegisterEncoder(Date{}, encodeSchemaMarshaler)
	encoder.RegisterEncoder(null.Float{}, encodeNullFloat)
	encoder.RegisterEncoder(null.Bool{}, encodeNullBool)
	return encoder
}
