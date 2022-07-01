//nolint:gochecknoglobals
package config

import (
	"fmt"
	"os"
)

var (
	version   string
	buildDate string
	commit    string
	notes     string
)

func Version() string {
	return version
}

func VersionVerbose() string {
	return fmt.Sprintf("Version %s\nRevision %s\nBuild at %s\n\n%s", Version(), Commit(), BuildDate(), notes)
}

func BuildDate() string {
	return buildDate
}

func Commit() string {
	return commit
}

func Hostname() string {
	host, _ := os.Hostname()

	if host == "" {
		host = "unknown"
	}

	return host
}
