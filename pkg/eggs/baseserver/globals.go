package baseserver

import (
	"sync"

	"github.com/kentik/ktranslate/pkg/eggs/features"
	"github.com/kentik/ktranslate/pkg/eggs/preconditions"
	"github.com/kentik/ktranslate/pkg/eggs/properties"
)

var globalBaseServer *BaseServer

var globalMutex sync.Mutex

func GetGlobalBaseServer() *BaseServer {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	preconditions.AssertNonNil(globalBaseServer, "globalBaseServer has not been set")
	return globalBaseServer
}

func setGlobalBaseServer(bs *BaseServer) {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	preconditions.AssertNil(globalBaseServer, "globalBaseServer has already been set")
	globalBaseServer = bs
}

// for testing
func resetGlobalBaseServer() {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	globalBaseServer = nil
}

// GetGlobalPropertyService looks up the property service on globalBaseServer.
// If it doesn't exist (i.e. during tests), that's ok, give us an empty one.
func GetGlobalPropertyService() properties.PropertyService {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	if globalBaseServer == nil {
		return properties.NewPropertyService()
	}

	return globalBaseServer.GetPropertyService()
}

// GetGlobalFeatureService looks up the feature service on globalBaseServer.
// If it doesn't exist (i.e. during tests), that's ok, give us an empty one.
func GetGlobalFeatureService() features.FeatureService {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	if globalBaseServer == nil {
		return features.NewFeatureService(properties.NewPropertyService())
	}

	return globalBaseServer.GetFeatureService()
}

// InitGlobalBaseServerTestingOnly is for tests to create a BaseServer and
// play with its PropertyService and FeatureService.
func InitGlobalBaseServerTestingOnly(propertyMap map[string]string, defaultPropertyBacking properties.PropertyBacking) {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	props := properties.NewPropertyService(
		properties.NewStaticMapPropertyBacking(propertyMap),
		defaultPropertyBacking,
	)

	globalBaseServer = &BaseServer{
		propertyService: props,
		featureService:  features.NewFeatureService(props),
	}
}

// ResetGlobalBaseServerTestingOnly should be defered after a call to
// InitGlobalBaseServerTestingOnly.
func ResetGlobalBaseServerTestingOnly() {
	resetGlobalBaseServer()
}
