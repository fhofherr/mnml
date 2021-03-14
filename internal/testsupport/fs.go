package testsupport

import (
	"os"
	"strings"
	"testing"
)

// MkdirTemp creates a temporary directory.
//
// It returns the path to the created temporary directory and a clean-up
// function which allows to delete it.
func MkdirTemp(t *testing.T) (string, func()) {
	t.Helper()

	prefix := strings.ReplaceAll(t.Name(), string(os.PathSeparator), "_")
	tempDir, err := os.MkdirTemp("", prefix)
	if err != nil {
		t.Fatal(err)
	}
	return tempDir, func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Error(err)
		}
	}
}
