package gosoapware

import (
	"context"
	"encoding/xml"
	"log"
	"net/http"
	"regexp"
)

// Add registers a handler for a specific SOAP action.
func (s *Soapware) Add(action string, next http.Handler) {
	s.handlers[action] = next
}

// Handlers returns a http.Handler that routes SOAP requests to the appropriate handler based on the SOAP action.
func (s *Soapware) Handlers() http.Handler {
	contentTypeRegex := regexp.MustCompile(`(?i)SOAPAction\s*=\s*"(.*?)"`)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			defer func() {
				err := r.Body.Close()
				if err != nil {
					log.Printf("failed to close request body: %v\n", err)
				}
			}()
		}

		var envelope Envelope
		if err := xml.NewDecoder(r.Body).Decode(&envelope); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		SOAPAction := r.Header.Get("SOAPAction")
		if SOAPAction == "" {
			contentType := r.Header.Get("Content-Type")
			matches := contentTypeRegex.FindStringSubmatch(contentType)

			if len(matches) > 1 {
				SOAPAction = matches[1]
			} else {
				http.Error(w, "unsupported SOAP body type", http.StatusNotFound)
				return
			}
		}

		ctx := context.WithValue(r.Context(), s.contentKey, envelope)
		handler, ok := s.handlers[SOAPAction]
		if !ok {
			http.Error(w, "handler for this SOAP action is missing", http.StatusNotImplemented)
			return
		}

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
