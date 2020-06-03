package api

// API messages
const (
	OkStatus            = "ok"
	ErrorTooLongTimeout = "timeout too long"
	ErrorInvalidRequest = "request is invalid"
)

// OkResponseBody is used for serializing ok response
type OkResponseBody struct {
	Status string `json:"status"`
}

// ErrorResponseBody is used for serializing error response
type ErrorResponseBody struct {
	Error string `json:"error"`
}
