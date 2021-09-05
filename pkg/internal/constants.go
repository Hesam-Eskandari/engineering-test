package internal

type ContextKey string

const (
	BodyContextKey       ContextKey = "body"
	MenuContextKey       ContextKey = "menu"
	MethodContextKey     ContextKey = "method"
	POSAddressContextKey ContextKey = "destination"
	RequestContextKey    ContextKey = "pos"

	POSAlpha string = "alpha"
	POSBeta  string = "beta"
)
