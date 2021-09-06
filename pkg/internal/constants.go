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

func getPoseMap() map[string]bool {
	posMap := make(map[string]bool, 0)
	posMap[POSAlpha] = true
	posMap[POSBeta] = true
	return posMap
}

func ValidatePoses(reqPos string) bool {
	if _, ok := getPoseMap()[reqPos]; ok {
		return true
	}
	return false
}
