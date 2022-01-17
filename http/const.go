package http

const (
	ApplicationJSON        = "application/json"
	ApplicationUrlencoded  = "application/x-www-form-urlencoded"
	ApplicationOctetStream = "application/octet-stream"
	ApplicationZIP         = "application/zip"

	MultipartFormdata = "multipart/form-data"
)

const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH" // RFC 5789
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"
)
