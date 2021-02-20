package ic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrafficNameFromNumbers(t *testing.T) {
	assert.Equal(t, "internal", TrafficNameFromNumbers(NETWORK_SRC_INTERNAL, NETWORK_SRC_INTERNAL))
	assert.Equal(t, "through", TrafficNameFromNumbers(NETWORK_SRC_EXTERNAL, NETWORK_SRC_EXTERNAL))
	assert.Equal(t, "from outside, terminated inside", TrafficNameFromNumbers(NETWORK_SRC_EXTERNAL, NETWORK_SRC_INTERNAL))
	assert.Equal(t, "originated inside, to outside", TrafficNameFromNumbers(NETWORK_SRC_INTERNAL, NETWORK_SRC_EXTERNAL))
	assert.Equal(t, "from outside to cloud", TrafficNameFromNumbers(NETWORK_SRC_EXTERNAL, NETWORK_SRC_CLOUD_AWS))
	assert.Equal(t, "from outside to cloud", TrafficNameFromNumbers(NETWORK_SRC_EXTERNAL, NETWORK_SRC_CLOUD_AZURE))
	assert.Equal(t, "from outside to cloud", TrafficNameFromNumbers(NETWORK_SRC_EXTERNAL, NETWORK_SRC_CLOUD_GCP))
	assert.Equal(t, "from outside to cloud", TrafficNameFromNumbers(NETWORK_SRC_EXTERNAL, NETWORK_SRC_CLOUD_IBM))
	assert.Equal(t, "from cloud to outside", TrafficNameFromNumbers(NETWORK_SRC_CLOUD_AWS, NETWORK_SRC_EXTERNAL))
	assert.Equal(t, "from cloud to outside", TrafficNameFromNumbers(NETWORK_SRC_CLOUD_AZURE, NETWORK_SRC_EXTERNAL))
	assert.Equal(t, "from cloud to outside", TrafficNameFromNumbers(NETWORK_SRC_CLOUD_GCP, NETWORK_SRC_EXTERNAL))
	assert.Equal(t, "from cloud to inside", TrafficNameFromNumbers(NETWORK_SRC_CLOUD_AWS, NETWORK_SRC_INTERNAL))
	assert.Equal(t, "from cloud to inside", TrafficNameFromNumbers(NETWORK_SRC_CLOUD_AZURE, NETWORK_SRC_INTERNAL))
	assert.Equal(t, "from cloud to inside", TrafficNameFromNumbers(NETWORK_SRC_CLOUD_GCP, NETWORK_SRC_INTERNAL))
	assert.Equal(t, "from inside to cloud", TrafficNameFromNumbers(NETWORK_SRC_INTERNAL, NETWORK_SRC_CLOUD_AWS))
	assert.Equal(t, "from inside to cloud", TrafficNameFromNumbers(NETWORK_SRC_INTERNAL, NETWORK_SRC_CLOUD_AZURE))
	assert.Equal(t, "from inside to cloud", TrafficNameFromNumbers(NETWORK_SRC_INTERNAL, NETWORK_SRC_CLOUD_GCP))
	assert.Equal(t, "multi-cloud", TrafficNameFromNumbers(NETWORK_SRC_CLOUD_AWS, NETWORK_SRC_CLOUD_AZURE))
	assert.Equal(t, "multi-cloud", TrafficNameFromNumbers(NETWORK_SRC_CLOUD_AWS, NETWORK_SRC_CLOUD_GCP))
	assert.Equal(t, "multi-cloud", TrafficNameFromNumbers(NETWORK_SRC_CLOUD_AZURE, NETWORK_SRC_CLOUD_AWS))
	assert.Equal(t, "multi-cloud", TrafficNameFromNumbers(NETWORK_SRC_CLOUD_AZURE, NETWORK_SRC_CLOUD_GCP))
	assert.Equal(t, "multi-cloud", TrafficNameFromNumbers(NETWORK_SRC_CLOUD_GCP, NETWORK_SRC_CLOUD_AZURE))
	assert.Equal(t, "multi-cloud", TrafficNameFromNumbers(NETWORK_SRC_CLOUD_GCP, NETWORK_SRC_CLOUD_AWS))
	assert.Equal(t, "multi-cloud", TrafficNameFromNumbers(NETWORK_SRC_CLOUD_IBM, NETWORK_SRC_CLOUD_AZURE))
}
