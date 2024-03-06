package gosoapware

import (
	"encoding/xml"
)

// Envelope represents the outer structure of a SOAP message.
type Envelope struct {
	Namespaces map[string]string `xml:"-"`
	Header     *Header
	Body       *Body
}

// Body represents the body of a SOAP message.
type Body struct {
	XMLName xml.Name `xml:"Body"`
	Content []byte   `xml:",innerxml"`
	Fault   *Fault   `xml:"Fault"`
}

// Header represents the header of a SOAP message.
type Header struct {
	XMLName    xml.Name    `xml:"Header"`
	Content    string      `xml:",chardata"`
	Attrs      []xml.Attr  `xml:",any,attr"`
	SubHeaders []SubHeader `xml:",any"`
}

// SubHeader represents a nested header within a SOAP message header.
type SubHeader struct {
	XMLName    xml.Name
	Content    string      `xml:",chardata"`
	Attrs      []xml.Attr  `xml:",any,attr"`
	SubHeaders []SubHeader `xml:",any"`
}

// Fault represents a SOAP fault structure.
type Fault struct {
	Code   string `xml:"faultcode"`
	String string `xml:"faultstring"`
	Actor  string `xml:"faultactor"`
	Detail Detail `xml:"detail"`
}

// Detail represents detailed information about a SOAP fault.
type Detail struct {
	XMLName   xml.Name
	SubFaults []SubDetail `xml:",any"`
	Content   string      `xml:",chardata"`
}

// SubDetail represents a nested detail within a SOAP fault detail.
type SubDetail struct {
	XMLName    xml.Name
	Content    string      `xml:",chardata"`
	SubDetails []SubDetail `xml:",any"`
}
