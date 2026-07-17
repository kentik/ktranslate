package kt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// The Docker image ENTRYPOINT always passes -api_devices /etc/ktranslate/devices.json.
// After Site moved onto protobuf string IDs (#916), numeric site.id / site.company_id in
// that file caused every container to fail at startup with:
//
//	json: cannot unmarshal number into Go struct field Site.Site.id of type string
func TestShippedDevicesJSONUnmarshals(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	devicesFile := filepath.Join(filepath.Dir(thisFile), "..", "..", "config", "devices.json")

	raw, err := os.ReadFile(devicesFile)
	if err != nil {
		t.Fatalf("reading %s: %v", devicesFile, err)
	}

	var devices map[string]*Device
	if err := json.Unmarshal(raw, &devices); err != nil {
		t.Fatalf("unmarshal shipped config/devices.json: %v", err)
	}
	if len(devices) == 0 {
		t.Fatal("expected at least one device in shipped config/devices.json")
	}

	for id, d := range devices {
		if d == nil {
			t.Fatalf("device %q is nil", id)
		}
		if d.Site == nil {
			continue
		}
		if d.Site.GetId() == "" {
			t.Fatalf("device %q has site with empty id", id)
		}
	}
}
