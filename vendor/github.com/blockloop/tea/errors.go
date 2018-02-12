package tea

import "net/http"

// StatusError creates a response message consisting of the status code
// and the http.StatusText which applies to that code.
//
// This is useful within StatusHandlerFuncs to quickly break out of the
// normal flow of code
//
// Example:
//   func(w http.ResponseWriter, r *http.Request) (int, interface{}) {
//           u, p, ok := r.BasicAuth()
//           if !ok {
//                   return StatusError(StatusUnauthorized)
//	     }
//   }
func StatusError(status int) (int, *ErrorResponse) {
	return Error(status, http.StatusText(status))
}

// Error creates a response message consisting of the status code
// and the error string provided. The error string will be rendered back
// to the client. This is ideal for client errors where the client should
// be informed of the specific error message.
//
// This is useful within StatusHandlerFuncs to quickly break out of the
// normal flow of code. It also renders the status code in the body of
// the response which is often very helpful for clients.
//
// Example:
//   func CreateUser(w http.ResponseWriter, r *http.Request) (int, interface{}) {
//           // parse request
//           if req.Name == "" {
//                   return Error(400, "name is required")
//	     }
//   }
func Error(status int, err string) (int, *ErrorResponse) {
	return status, &ErrorResponse{
		Code:  status,
		Error: err,
	}
}

// ErrorResponse is a generic response object
type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

// NotFound is a handler for 404 requests
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	_, msg := StatusError(404)
	Responder(w, r, msg)
}

// MethodNotAllowed is a handler for 405 requests
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(405)
	_, msg := StatusError(405)
	Responder(w, r, msg)
}
