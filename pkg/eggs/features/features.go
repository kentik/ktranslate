package features

import (
	"regexp"
)

// We're being more restrictive with feature names than property names to create some breathing room as feature
// flags are backed by properties; if feature names can't use _ or - we are free to use them internally as field
// separators
var validFeatureName = regexp.MustCompile(`^[a-zA-Z0-9.]+$`)

func IsValidFeatureName(featureName string) bool {
	return validFeatureName.MatchString(featureName)
}

type FeatureService interface {
	EnabledGlobally(featureName string) bool
	Enabled(featureName string, sub ...interface{}) bool
}
