package auth

import (
	"net"

	"github.com/kentik/ktranslate/pkg/kt"
)

type AuthConfig struct {
	DevicesFile string
}

type DeviceWrapper struct {
	Device *kt.Device `json:"device"`
}

type AllDeviceWrapper struct {
	Devices []*kt.Device `json:"devices"`
}

type DeviceCreate struct {
	Name        string   `json:"device_name"`
	Type        string   `json:"device_type"`
	Description string   `json:"device_description"`
	SampleRate  uint32   `json:"device_sample_rate,string"`
	BgpType     string   `json:"device_bgp_type"`
	PlanID      int      `json:"plan_id,omitempty"`
	SiteID      int      `json:"site_id,omitempty"`
	IPs         []net.IP `json:"sending_ips"`
	CdnAttr     string   `json:"cdn_attr"`
	ExportId    int      `json:"cloud_export_id,omitempty"`
	Subtype     string   `json:"device_subtype"`
	Region      string   `json:"cloud_region,omitempty"`
	Zone        string   `json:"cloud_zone,omitempty"`
}
