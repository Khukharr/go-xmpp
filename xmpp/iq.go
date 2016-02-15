package xmpp

import (
	"encoding/xml"

	"github.com/processone/gox/xmpp/iot"
)

// info/query
type ClientIQ struct {
	XMLName xml.Name `xml:"jabber:client iq"`
	Packet
	Payload IQPayload `xml:",omitempty"`
	RawXML  string    `xml:",innerxml"`
	// TODO We need to support detecting the IQ namespace / Query packet
	// 	Error   clientError
}

type IQPayload interface {
	IsIQPayload()
}

// UnmarshalXML implements custom parsing for IQs
func (iq *ClientIQ) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	iq.XMLName = start.Name
	// Extract IQ attributes
	for _, attr := range start.Attr {
		if attr.Name.Local == "id" {
			iq.Id = attr.Value
		}
		if attr.Name.Local == "to" {
			iq.To = attr.Value
		}
		if attr.Name.Local == "from" {
			iq.From = attr.Value
		}
		if attr.Name.Local == "lang" {
			iq.Lang = attr.Value
		}
	}

	// decode inner elements
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}

		var p IQPayload
		switch tt := t.(type) {

		case xml.StartElement:
			switch tt.Name.Space + " " + tt.Name.Local {
			case "urn:ietf:params:xml:ns:xmpp-bind bind":
				p = new(bindBind)
			case "urn:xmpp:iot:control set":
				p = new(iot.ControlSet)
				// TODO: Add a default Type that passes RawXML
			}
			if p != nil {
				err = d.DecodeElement(p, &tt)
				if err != nil {
					return err
				}
				iq.Payload = p
				p = nil
			}

		case xml.EndElement:
			if tt == start.End() {
				return nil
			}
		}
	}
}
