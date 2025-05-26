package entities

type contextKey string

const (
	LoggerKey = contextKey("logger")
	RequestId = contextKey("request_id")
)
