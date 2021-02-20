package features

import (
	"github.com/kentik/ktranslate/pkg/eggs/properties"
	"fmt"
)

type stdFeatureService struct {
	propsvc properties.PropertyService
}

func NewFeatureService(propsvc properties.PropertyService) *stdFeatureService {
	return &stdFeatureService{
		propsvc: propsvc,
	}
}

func (dfsvc *stdFeatureService) EnabledGlobally(featureName string) bool {
	if !IsValidFeatureName(featureName) {
		panic(fmt.Errorf("Invalid feature name '%s'", featureName))
	}

	return dfsvc.Enabled(featureName)
}

func (dfsvc *stdFeatureService) Enabled(featureName string, sub ...interface{}) bool {

	if !IsValidFeatureName(featureName) {
		panic(fmt.Errorf("Invalid feature name '%s'", featureName))
	}

	return dfsvc.propsvc.GetBoolSub(fmt.Sprintf("features.%s", featureName), false, sub...).Bool()
}
