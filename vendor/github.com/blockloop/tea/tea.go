package tea

import (
	"net/http"

	"github.com/go-chi/render"
)

// Responder is the default responder used to write the messages back to the
// client. It uses render.DefaultResponder by default which will respond using
// the appropriate data type based on the Accept header. You can set it to
// render.JSON, render.XML, render.Data, etc. See github.com/go-chi/render
// for more info
var Responder = render.DefaultResponder

// StatusHandlerFunc is a handler that returns a status code and a message body
type StatusHandlerFunc func(w http.ResponseWriter, r *http.Request) (int, interface{})

// Handler wraps a StatusHandlerFunc and returns a standard lib http.HandlerFunc
func Handler(h StatusHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, response := h(w, r)
		w.WriteHeader(status)
		if response == nil {
			return
		}

		Responder(w, r, response)
	}
}
