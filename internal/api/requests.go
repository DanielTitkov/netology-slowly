package api

// SlowRequestBody is used for deserializing slow request
type SlowRequestBody struct {
	Timeout int `json:"timeout"`
}
