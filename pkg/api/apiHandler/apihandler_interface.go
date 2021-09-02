package apiHandler

import "net/http"

// Handler - interface with necessary method signatures
// this interface is better to be a dependency than an internal
// there must be an Authorize method to authorize caller in real applications
type Handler interface {
	URL() string                                      // URL returns the request URL to this handler
	Methods() []string                                // Methods returns allowed HTTP methods
	ParseArgs(r *http.Request) (*http.Request, error) // ParseArgs parses and validates request arguments
	Process(r *http.Request) *http.Response           // Process handles the request
}
