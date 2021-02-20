package cdn

var (
	// Const, unchanging
	cdn2Name = map[uint32]string{
		0:  "--",
		1:  "AFXcdn",
		2:  "Akamai",
		3:  "Alibaba Cloud",
		4:  "Amazon + AWS",
		5:  "Ananke",
		6:  "Apple",
		7:  "AT&T",
		8:  "Azion",
		9:  "BelugaCDN",
		10: "Blue Hat Network",
		11: "CacheFly",
		12: "CDN77",
		13: "CDNetworks",
		14: "CDNify",
		15: "ChinaCache",
		16: "Cloudflare",
		17: "EdgeCast Verizon",
		18: "Facebook",
		19: "Fastly",
		20: "G-Core CDN",
		21: "Google Youtube",
		22: "Hibernia",
		23: "Imperva Incapsula",
		24: "Instartlogic",
		25: "Internap",
		26: "KeyCDN",
		27: "LeaseWeb CDN",
		28: "Level3",
		29: "Limelight",
		30: "Medianova",
		31: "Netflix",
		32: "Mirror Image",
		33: "Ngenix",
		34: "QUANTIL ChinaNetCenter",
		35: "Reflected Networks - Swift",
		36: "Alibaba_Cloud",
		37: "Tata communications",
		38: "Telefonica",
		39: "TurboBytes",
		40: "Microsoft Azure",
		41: "Yahoo",
		42: "Zenedge",
		43: "Stackpath MaxCDN NetDNA Highwinds",
		44: "Dropbox",
		45: "Pandora",
		46: "Twitch.TV",
		47: "Comcast CDN",
		48: "BunnyCDN",
		49: "CDNvideo",
		50: "cdnnow",
		51: "SingularCDN",
		52: "Kingsoft Cloud",
		53: "Tencent Cloud",
	}
)

func NamesByCDNs(cdns []uint32) []string {
	ret := make([]string, len(cdns))
	for i := 0; i < len(cdns); i++ {
		ret[i] = cdn2Name[cdns[i]]
	}
	return ret
}

func NameByCDN(cdn uint32) string {
	return cdn2Name[cdn]
}
