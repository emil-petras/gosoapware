[![Language Go](https://img.shields.io/badge/Language-Go-orange.svg?style=shields)](https://golang.org/)

# **gosoapware**

## Description
**gosoapware** is a lightweight library that makes parsing SOAP messages easier. It relies on the standard net\http.

## Installation
To install use 

```go
go get -u github.com/emil-petras/gosoapware
```

## Example SOAP messages

Header for both requests need to have SOAPAction: MyType

SOAP request message

```xml
<?xml version="1.0"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:web="http://example.com/webservice">
   <soapenv:Header>
      <web:AuthHeader>
         <web:Username>exampleUser</web:Username>
         <web:Password>examplePassword</web:Password>
      </web:AuthHeader>
   </soapenv:Header>
   <soapenv:Body>
      <MyType xmlns="http://example.com/yournamespace">
         <Field1>Value1</Field1>
         <Field2>Value2</Field2>
      </MyType>
   </soapenv:Body>
</soapenv:Envelope>
```

SOAP request message containing a fault

```xml
<?xml version="1.0"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:web="http://example.com/webservice">
   <soapenv:Body>
    <soap:Fault>
    <faultcode>soap:Client</faultcode>
    <faultstring>Invalid or missing parameter</faultstring>
    <detail>
        <error xmlns="http://www.example.com/">Detail</error>
                        <detail>
        <error xmlns="http://www.example.com/">Subdetail</error>
            </detail>
    </detail>
</soap:Fault>
   </soapenv:Body>
</soapenv:Envelope>
```

## How to use

Example how to parse it:

```go
package main

import (
	"encoding/xml"
	"fmt"
	"main/gosoapware" // Importing your custom gosoapware package
	"net/http"
)

// MyType represents the structure of your expected SOAP request body.
type MyType struct {
	XMLName xml.Name
	Field1  string `xml:"http://example.com/yournamespace Field1"` // XML tag for Field1
	Field2  string `xml:"http://example.com/yournamespace Field2"` // XML tag for Field2
}

func main() {
	// Initialize the SOAP middleware.
	soap := gosoapware.NewSoapware("envelope")
	
	// Register the handler for 'MyType' SOAP action.
	soap.Add("MyType", http.HandlerFunc(MyHandler))

	// Set up the HTTP server mux and handle requests using SOAP middleware.
	mux := http.NewServeMux()
	mux.Handle("/", soap.Handlers())

	// Start the server on port 8080.
	http.ListenAndServe(":8080", mux)
}

// envelopeKey is the context key used for retrieving the SOAP envelope.
var envelopeKey = gosoapware.SOAPContentKey("envelope")

// MyHandler is the function that handles incoming SOAP requests.
func MyHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the SOAP envelope from the request context.
	envelope, ok := r.Context().Value(envelopeKey).(gosoapware.Envelope)
	if !ok {
		http.Error(w, "invalid request data", http.StatusBadRequest)
		return
	}

	// Check if there is a SOAP fault and handle it.
	if envelope.Body.Fault != nil {
		faultCode := envelope.Body.Fault.Code
		faultString := envelope.Body.Fault.String
		// Handle the detail field of the fault. Assuming Detail is a struct, modify as needed.
		detail := envelope.Body.Fault.Detail

		// Correctly format the fault detail for printing.
		fmt.Fprintf(w, "Code: %s, String: %s, Detail: %+v", faultCode, faultString, detail)
		return
	}

	// Extract user and secret from the SOAP header for authentication or other purposes.
	user := envelope.Header.SubHeaders[0].SubHeaders[0].Content
	secret := envelope.Header.SubHeaders[0].SubHeaders[1].Content

	var data MyType

	// Unmarshal the XML body into the MyType struct.
	if err := xml.Unmarshal(envelope.Body.Content, &data); err != nil {
		http.Error(w, "Error unmarshalling request data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the response back with the extracted data.
	fmt.Fprintf(w, "User: %s, Secret: %s, Received Field1: %s, Field2: %s", user, secret, data.Field1, data.Field2)
}
```

