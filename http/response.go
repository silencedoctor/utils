package http

import (
	"net/http"
)

// Response ...
type Response http.Response

// IsOK ...
func (r *Response) IsOK() bool {
	return r.StatusCode >= 100 && r.StatusCode < 300
}
