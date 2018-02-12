package tea

import (
	"encoding/json"
	"net/http"
)

// Body parses a JSON request body and validates it's contents
func Body(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return Validate.StructCtx(r.Context(), v)
}
