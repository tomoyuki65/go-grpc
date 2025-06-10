package context

type contextKey string

const (
	XRequestId     contextKey = "x-request-id"
	XRequestSource contextKey = "x-request-source"
	XUid           contextKey = "x-uid"
	Status         contextKey = "status"
	StatusCode     contextKey = "statusCode"
)
