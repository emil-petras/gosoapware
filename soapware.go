package gosoapware

import (
	"net/http"
)

// SOAPContentKey is a type used for context keys in handling SOAP requests.
type SOAPContentKey string

// Soapware is a struct that holds the key for accessing SOAP content and a map of handlers for different SOAP actions.
type Soapware struct {
	contentKey SOAPContentKey
	handlers   map[string]http.Handler
}

// NewSoapware creates a new instance of Soapware with the specified SOAPContentKey.
func NewSoapware(SOAPContentKey SOAPContentKey) *Soapware {
	return &Soapware{
		contentKey: SOAPContentKey,
		handlers:   make(map[string]http.Handler),
	}
}
