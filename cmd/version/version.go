package version

//go:generate bash ../../scripts/version.sh

/*
//BEGIN-TEMPLATE
package version

import "github.com/kentik/ktranslate/pkg/eggs/version"

var Version = version.VersionInfo{
	Version:  "@VERSION@",
	Date:     "@DATE@",
	Platform: "@PLATFORM@",
	Distro:   "@GOLANG@",
}
//END-TEMPLATE
*/
