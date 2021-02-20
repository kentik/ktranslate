package version

import "fmt"

type VersionInfo struct {
	Version  string
	Date     string
	Platform string
	Distro   string
}

func (v VersionInfo) String() string {
	return fmt.Sprintf("version %s built on %s", v.Version, v.Date)
}
