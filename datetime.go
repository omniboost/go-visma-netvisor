package netvisor

import (
	"encoding/json"
	"encoding/xml"
	"time"
)

type DateTime struct {
	time.Time
	Layout string
}

func (d DateTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if d.Time.IsZero() {
		return e.Encode(nil)
	}

	return e.EncodeElement(d.Time.Format(d.Layout), start)
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return json.Marshal(nil)
	}

	return json.Marshal(d.Time.Format(d.Layout))
}

func (d DateTime) MarshalSchema() string {
	layout := "2006-01-02"
	return d.Time.Format(layout)
}

func (d DateTime) IsEmpty() bool {
	return d.Time.IsZero()
}

func (d *DateTime) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var value string
	err := dec.DecodeElement(&value, &start)
	if err != nil {
		return err
	}

	if value == "" {
		return nil
	}

	// first try standard date
	layout := time.RFC3339
	d.Time, err = time.Parse(layout, value)
	if err == nil {
		d.Layout = layout
		return nil
	}

	layout = "2006-01-02"
	d.Time, err = time.Parse(layout, value)
	if err == nil {
		d.Layout = layout
		return err
	}

	layout = "2.1.2006"
	d.Time, err = time.Parse(layout, value)
	if err == nil {
		d.Layout = layout
		return err
	}

	layout = "2.1.2006 15:04:05"
	d.Time, err = time.Parse(layout, value)
	if err == nil {
		d.Layout = layout
		return err
	}

	return nil
}
