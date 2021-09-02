package internal

type ContextKey string

const (
	MethodContextKey     ContextKey = "method"
	POSAddressContextKey ContextKey = "destination"
	BodyContextKey       ContextKey = "body"
)
