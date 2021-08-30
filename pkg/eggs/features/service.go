package features

import (
	"fmt"
	"github.com/kentik/ktranslate/pkg/eggs/properties"
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
		panic(fmt.Errorf("You used an unsupported feature name: %v.", featureName))
	}

	return dfsvc.Enabled(featureName)
}

func (dfsvc *stdFeatureService) Enabled(featureName string, sub ...interface{}) bool {

	if !IsValidFeatureName(featureName) {
		panic(fmt.Errorf("You used an unsupported feature name: %v.", featureName))
	}

	return dfsvc.propsvc.GetBoolSub(fmt.Sprintf("features.%s", featureName), false, sub...).Bool()
}
