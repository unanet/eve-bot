package middleware

import "net/http"

// Adapter is the http.Handler middleware type
type Adapter func(http.Handler) http.Handler

// Adapt is an Adapter wrapper for wrapping middleware handlers
// See api.go for useage
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}
