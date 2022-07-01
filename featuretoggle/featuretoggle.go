package featuretoggle

// values represents the feature toggle key-value pairs
type values map[string]interface{}

// fts its the aplication feature toggles
var fts values

// Init instantiates the feature toggles of a project, given the key-value pairs.
func Init(v values) {
	fts = v
}

// IsEnabled searches for a feature toggle key, and checks if its value its true.
// Returns false if the feature toggle is false, does not exist, or its not a bool value.
func IsEnabled(key string) bool {
	val, ok := fts[key]
	if !ok {
		return false
	}

	valBool, ok := val.(bool)
	if !ok {
		return false
	}

	return valBool
}
