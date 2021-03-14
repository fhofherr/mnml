package version

import (
	"fmt"
	"time"
)

var (
	// Version contains the build version, or 'development' if this is a
	// development build.
	Version = "devlopment"

	// Commit contains the commit hash of the build, or 'development' if this
	// is a development build.
	Commit = "development"

	// Date contains the built date, or the current time if this is a
	// development build.
	Date = time.Now().UTC().Format(time.RFC3339)
)

// String returns the current version information formatted as a string.
func String() string {
	return fmt.Sprintf("%s-%s-%s", Version, Commit, Date)
}
