package testsupport

import "flag"

var flagUpdate = flag.Bool("update", false, "Update test golden files")

// IsUpdate returns true if the -update flag was passed to the test.
func IsUpdate() bool {
	return *flagUpdate
}
