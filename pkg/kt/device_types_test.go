package kt

import (
	"net"
	"testing"

	"github.com/kentik/api-schema-public/gen/go/kentik/device/v202504beta2"
	"google.golang.org/protobuf/types/known/durationpb"
)

// Helper to ptr-ify a bool, since MinimizeSnmp is *bool in the proto.
func boolPtr(b bool) *bool { return &b }

// --- happy path ---

func TestMapDeviceDetailedToDevice_Full(t *testing.T) {
	proto := &device.DeviceDetailed{
		Id:                  "42",
		CompanyId:           "7",
		DeviceName:          "core-router-01",
		DeviceType:          "router",
		DeviceSubtype:       "cisco_csr",
		DeviceDescription:   "Core router",
		DeviceSampleRate:    "1024",
		DeviceSnmpIp:        "10.0.0.1",
		DeviceSnmpCommunity: "public",
		DeviceBgpType:       "device",
		CdnAttr:             "N",
		MaxFlowRate:         500,
		CustomColumns:       "col1,col2",
		SendingIps:          []string{"192.168.1.1", "192.168.1.2"},
		MinimizeSnmp:        boolPtr(false),

		Site: &device.Site{
			Id:       "99",
			SiteName: "HQ",
		},
		Plan: &device.Plan{
			Id:   "12",
			Name: "Premium",
		},
		Labels: []*device.Label{
			{Id: "1", Name: "prod", Description: "Production"},
			{Id: "2", Name: "core", Description: "Core network"},
		},
		DeviceSnmpV3Conf: &device.DeviceSnmpV3Conf{
			Username:                 "snmpuser",
			AuthenticationProtocol:   "SHA",
			AuthenticationPassphrase: "authpass",
			PrivacyProtocol:          "AES",
			PrivacyPassphrase:        "privpass",
		},
		AllInterfaces: []*device.Interface{
			{
				SnmpId:               "101",
				SnmpAlias:            "eth0",
				SnmpSpeed:            "1000",
				InterfaceDescription: "Uplink",
				NetworkBoundary:      "external",
				ConnectivityType:     "transit",
				Provider:             "ISP-A",
			},
			{
				SnmpId:               "102",
				SnmpAlias:            "eth1",
				SnmpSpeed:            "10000",
				InterfaceDescription: "Peering",
				NetworkBoundary:      "external",
				ConnectivityType:     "peering",
				Provider:             "IX-B",
			},
		},
		CustomColumnData: []*device.CustomColumnData{
			{
				DeviceId:    "42",
				FieldId:     "f1",
				ColName:     "region",
				Description: "Region tag",
				ColType:     "string",
				DeviceType:  "router",
			},
		},
	}

	got, err := MapDeviceDetailedToDevice(proto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// --- scalar fields ---
	assertEqual(t, "ID", int64(got.ID), int64(42))
	assertEqual(t, "CompanyID", int64(got.CompanyID), int64(7))
	assertEqual(t, "Name", got.Name, "core-router-01")
	assertEqual(t, "DeviceType", got.DeviceType, "router")
	assertEqual(t, "DeviceSubtype", got.DeviceSubtype, "cisco_csr")
	assertEqual(t, "Description", got.Description, "Core router")
	assertEqual(t, "SampleRate", got.SampleRate, uint32(1024))
	assertEqual(t, "BgpType", got.BgpType, "device")
	assertEqual(t, "CdnAttr", got.CdnAttr, "N")
	assertEqual(t, "MaxFlowRate", got.MaxFlowRate, 500)
	assertEqual(t, "CustomStr", got.CustomStr, "col1,col2")
	assertEqual(t, "SnmpCommunity", got.SnmpCommunity, "public")
	assertEqual(t, "SnmpIp", got.SnmpIp, "10.0.0.1")

	// --- IP ---
	wantIP := net.ParseIP("10.0.0.1")
	if !got.IP.Equal(wantIP) {
		t.Errorf("IP: got %v, want %v", got.IP, wantIP)
	}

	// --- SendingIps ---
	if len(got.SendingIps) != 2 {
		t.Fatalf("SendingIps: got %d entries, want 2", len(got.SendingIps))
	}
	assertEqual(t, "SendingIps[0]", got.SendingIps[0].String(), "192.168.1.1")
	assertEqual(t, "SendingIps[1]", got.SendingIps[1].String(), "192.168.1.2")

	// --- Site ---
	assertEqual(t, "Site.ID", got.Site.ID, 99)
	assertEqual(t, "Site.SiteName", got.Site.SiteName, "HQ")

	// --- Plan ---
	assertEqual(t, "Plan.ID", got.Plan.ID, 12)
	assertEqual(t, "Plan.Name", got.Plan.Name, "Premium")

	// --- Labels ---
	if len(got.Labels) != 2 {
		t.Fatalf("Labels: got %d, want 2", len(got.Labels))
	}
	assertEqual(t, "Labels[0].ID", got.Labels[0].ID, 1)
	assertEqual(t, "Labels[0].Name", got.Labels[0].Name, "prod")
	assertEqual(t, "Labels[0].Desc", got.Labels[0].Desc, "Production")
	assertEqual(t, "Labels[1].ID", got.Labels[1].ID, 2)
	assertEqual(t, "Labels[1].Name", got.Labels[1].Name, "core")

	// --- SnmpV3 ---
	if got.SnmpV3 == nil {
		t.Fatal("SnmpV3: got nil, want non-nil")
	}
	assertEqual(t, "SnmpV3.UserName", got.SnmpV3.UserName, "snmpuser")
	assertEqual(t, "SnmpV3.AuthenticationProtocol", got.SnmpV3.AuthenticationProtocol, "SHA")
	assertEqual(t, "SnmpV3.AuthenticationPassphrase", got.SnmpV3.AuthenticationPassphrase, "authpass")
	assertEqual(t, "SnmpV3.PrivacyProtocol", got.SnmpV3.PrivacyProtocol, "AES")
	assertEqual(t, "SnmpV3.PrivacyPassphrase", got.SnmpV3.PrivacyPassphrase, "privpass")

	// --- AllInterfaces ---
	if len(got.AllInterfaces) != 2 {
		t.Fatalf("AllInterfaces: got %d, want 2", len(got.AllInterfaces))
	}
	iface0 := got.AllInterfaces[0]
	assertEqual(t, "AllInterfaces[0].DeviceID", int64(iface0.DeviceID), int64(42))
	assertEqual(t, "AllInterfaces[0].SnmpID", int64(iface0.SnmpID), int64(101))
	assertEqual(t, "AllInterfaces[0].Alias", iface0.Alias, "eth0")
	assertEqual(t, "AllInterfaces[0].SnmpSpeedMbps", iface0.SnmpSpeedMbps, int64(1000))
	assertEqual(t, "AllInterfaces[0].Description", iface0.Description, "Uplink")
	assertEqual(t, "AllInterfaces[0].NetworkBoundary", iface0.NetworkBoundary, "external")
	assertEqual(t, "AllInterfaces[0].ConnectivityType", iface0.ConnectivityType, "transit")
	assertEqual(t, "AllInterfaces[0].Provider", iface0.Provider, "ISP-A")

	// --- Interfaces map ---
	if len(got.Interfaces) != 2 {
		t.Fatalf("Interfaces map: got %d entries, want 2", len(got.Interfaces))
	}
	if _, ok := got.Interfaces[IfaceID(101)]; !ok {
		t.Error("Interfaces map: missing key 101")
	}
	if _, ok := got.Interfaces[IfaceID(102)]; !ok {
		t.Error("Interfaces map: missing key 102")
	}

	// --- Customs ---
	if len(got.Customs) != 1 {
		t.Fatalf("Customs: got %d, want 1", len(got.Customs))
	}
	col := got.Customs[0]
	assertEqual(t, "Customs[0].Name", col.Name, "region")
	assertEqual(t, "Customs[0].Description", col.Description, "Region tag")
	assertEqual(t, "Customs[0].Type", col.Type, "string")
}

// --- nil input ---

func TestMapDeviceDetailedToDevice_NilInput(t *testing.T) {
	_, err := MapDeviceDetailedToDevice(nil)
	if err == nil {
		t.Fatal("expected error for nil input, got nil")
	}
}

// --- invalid IDs ---

func TestMapDeviceDetailedToDevice_InvalidDeviceID(t *testing.T) {
	proto := &device.DeviceDetailed{
		Id:               "not-a-number",
		CompanyId:        "7",
		DeviceSampleRate: "1",
	}
	_, err := MapDeviceDetailedToDevice(proto)
	if err == nil {
		t.Fatal("expected error for non-numeric device ID, got nil")
	}
}

func TestMapDeviceDetailedToDevice_InvalidCompanyID(t *testing.T) {
	proto := &device.DeviceDetailed{
		Id:               "1",
		CompanyId:        "not-a-number",
		DeviceSampleRate: "1",
	}
	_, err := MapDeviceDetailedToDevice(proto)
	if err == nil {
		t.Fatal("expected error for non-numeric company ID, got nil")
	}
}

// --- graceful degradation ---

func TestMapDeviceDetailedToDevice_EmptySendingIps(t *testing.T) {
	proto := &device.DeviceDetailed{Id: "1", CompanyId: "2", DeviceSampleRate: "100"}
	got, err := MapDeviceDetailedToDevice(proto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.SendingIps) != 0 {
		t.Errorf("SendingIps: got %v, want empty", got.SendingIps)
	}
}

func TestMapDeviceDetailedToDevice_InvalidSendingIpSkipped(t *testing.T) {
	proto := &device.DeviceDetailed{
		Id:               "1",
		CompanyId:        "2",
		DeviceSampleRate: "100",
		SendingIps:       []string{"bad-ip", "10.1.2.3", "also-bad"},
	}
	got, err := MapDeviceDetailedToDevice(proto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.SendingIps) != 1 {
		t.Fatalf("SendingIps: got %d entries, want 1", len(got.SendingIps))
	}
	assertEqual(t, "SendingIps[0]", got.SendingIps[0].String(), "10.1.2.3")
}

func TestMapDeviceDetailedToDevice_InvalidSampleRateDefaultsToZero(t *testing.T) {
	proto := &device.DeviceDetailed{
		Id:               "1",
		CompanyId:        "2",
		DeviceSampleRate: "not-a-number",
	}
	got, err := MapDeviceDetailedToDevice(proto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "SampleRate", got.SampleRate, uint32(0))
}

func TestMapDeviceDetailedToDevice_NilSnmpV3ConfYieldsNilPointer(t *testing.T) {
	proto := &device.DeviceDetailed{
		Id:               "1",
		CompanyId:        "2",
		DeviceSampleRate: "100",
		DeviceSnmpV3Conf: nil,
	}
	got, err := MapDeviceDetailedToDevice(proto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.SnmpV3 != nil {
		t.Errorf("SnmpV3: expected nil, got %+v", got.SnmpV3)
	}
}

func TestMapDeviceDetailedToDevice_NilSiteAndPlan(t *testing.T) {
	proto := &device.DeviceDetailed{
		Id:               "1",
		CompanyId:        "2",
		DeviceSampleRate: "100",
		Site:             nil,
		Plan:             nil,
	}
	got, err := MapDeviceDetailedToDevice(proto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "Site.ID", got.Site.ID, 0)
	assertEqual(t, "Site.SiteName", got.Site.SiteName, "")
	assertEqual(t, "Plan.ID", got.Plan.ID, 0)
	assertEqual(t, "Plan.Name", got.Plan.Name, "")
}

func TestMapDeviceDetailedToDevice_LabelWithNonNumericID(t *testing.T) {
	proto := &device.DeviceDetailed{
		Id:               "1",
		CompanyId:        "2",
		DeviceSampleRate: "100",
		Labels: []*device.Label{
			{Id: "bad", Name: "whatever", Description: "desc"},
		},
	}
	got, err := MapDeviceDetailedToDevice(proto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Labels) != 1 {
		t.Fatalf("Labels: got %d, want 1", len(got.Labels))
	}
	// Non-numeric label ID should fall back to 0.
	assertEqual(t, "Labels[0].ID", got.Labels[0].ID, 0)
	assertEqual(t, "Labels[0].Name", got.Labels[0].Name, "whatever")
}

func TestMapDeviceDetailedToDevice_InterfaceWithBadSnmpIDDefaultsToZero(t *testing.T) {
	proto := &device.DeviceDetailed{
		Id:               "5",
		CompanyId:        "2",
		DeviceSampleRate: "100",
		AllInterfaces: []*device.Interface{
			{SnmpId: "not-a-number", SnmpAlias: "eth0", SnmpSpeed: "1000"},
		},
	}
	got, err := MapDeviceDetailedToDevice(proto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.AllInterfaces) != 1 {
		t.Fatalf("AllInterfaces: got %d, want 1", len(got.AllInterfaces))
	}
	assertEqual(t, "AllInterfaces[0].SnmpID", int64(got.AllInterfaces[0].SnmpID), int64(0))
}

func TestMapDeviceDetailedToDevice_EmptyDeviceHasEmptyInterfaceMap(t *testing.T) {
	proto := &device.DeviceDetailed{
		Id:               "1",
		CompanyId:        "2",
		DeviceSampleRate: "100",
	}
	got, err := MapDeviceDetailedToDevice(proto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Interfaces == nil {
		t.Error("Interfaces map should be non-nil even when empty")
	}
	if len(got.Interfaces) != 0 {
		t.Errorf("Interfaces: got %d entries, want 0", len(got.Interfaces))
	}
}

// Verify that NMS config (present in the proto but absent from Device) does not
// cause a panic — the mapping simply ignores it.
func TestMapDeviceDetailedToDevice_NmsConfigIgnored(t *testing.T) {
	proto := &device.DeviceDetailed{
		Id:               "1",
		CompanyId:        "2",
		DeviceSampleRate: "100",
		Nms: &device.DeviceNmsConfig{
			AgentId:   "agent-xyz",
			IpAddress: "10.0.0.5",
			Snmp: &device.DeviceNmsSnmpConfig{
				CredentialName: "my-cred",
				Port:           161,
				Timeout:        durationpb.New(2e9),
			},
		},
	}
	_, err := MapDeviceDetailedToDevice(proto)
	if err != nil {
		t.Fatalf("unexpected error with NMS config present: %v", err)
	}
}

// --- helper ---

// assertEqual is a typed-equality helper that avoids importing testify.
func assertEqual[T comparable](t *testing.T, field string, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %v, want %v", field, got, want)
	}
}
