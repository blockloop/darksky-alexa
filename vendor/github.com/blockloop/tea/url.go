package tea

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

// ErrNoURLParam means you tried to get a URL param that does not exist
var ErrNoURLParam = fmt.Errorf("url param does not exist")

// URLInt gets a URL parameter from the request and parses it as an int. If
// the value of the URL parameter is empty then ErrNoURLParam is returned.
// If parsing passes, then the validation key will be used to validate the
// value of the int using Validate
func URLInt(r *http.Request, key, validation string) (int, error) {
	p := strings.TrimSpace(chi.URLParam(r, key))
	if len(p) == 0 {
		return 0, ErrNoURLParam
	}

	i, err := strconv.Atoi(p)
	if err != nil {
		return 0, err
	}

	if len(validation) == 0 {
		return i, nil
	}
	return i, Validate.VarCtx(r.Context(), &i, validation)
}

// URLInt64 gets a URL parameter from the request and parses it as an int64. If
// the value of the URL parameter is empty then ErrNoURLParam is returned.
// If parsing passes, then the validation key will be used to validate the
// value of the int64 using Validate
func URLInt64(r *http.Request, key, validation string) (int64, error) {
	p := strings.TrimSpace(chi.URLParam(r, key))
	if len(p) == 0 {
		return 0, ErrNoURLParam
	}

	i, err := strconv.ParseInt(p, 0, 64)
	if err != nil {
		return 0, err
	}

	if len(validation) == 0 {
		return i, nil
	}
	return i, Validate.VarCtx(r.Context(), &i, validation)
}

// URLFloat gets a URL parameter from the request and parses it as an float. If
// the value of the URL parameter is empty then ErrNoURLParam is returned.
// If parsing passes, then the validation key will be used to validate the
// value of the float using Validate
func URLFloat(r *http.Request, key, validation string) (float64, error) {
	p := strings.TrimSpace(chi.URLParam(r, key))
	if len(p) == 0 {
		return 0, ErrNoURLParam
	}

	i, err := strconv.ParseFloat(p, 64)
	if err != nil {
		return 0, err
	}

	if len(validation) == 0 {
		return i, nil
	}
	return i, Validate.VarCtx(r.Context(), &i, validation)
}

// URLUint gets a URL parameter from the request and parses it as an uint64. If
// the value of the URL parameter is empty then ErrNoURLParam is returned.
// If parsing passes, then the validation key will be used to validate the
// value of the uint64 using Validate
func URLUint(r *http.Request, key, validation string) (uint64, error) {
	p := strings.TrimSpace(chi.URLParam(r, key))
	if len(p) == 0 {
		return 0, ErrNoURLParam
	}

	i, err := strconv.ParseUint(p, 0, 64)
	if err != nil {
		return 0, err
	}

	if len(validation) == 0 {
		return i, nil
	}
	return i, Validate.VarCtx(r.Context(), &i, validation)
}
