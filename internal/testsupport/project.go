package testsupport

import (
	"os"
	"path/filepath"
	"testing"
)

// ProjectRoot searches from the current working directory upwards until it
// finds a go.mod file. The directory containing the go.mod file is then
// assumed to be the root of this project.
func ProjectRoot(t *testing.T) string {
	t.Helper()

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	cur, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for cur != home {
		goModFile := filepath.Join(cur, "go.mod")
		s, err := os.Stat(goModFile)
		if os.IsNotExist(err) {
			cur = filepath.Dir(cur)
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		if !s.Mode().IsRegular() {
			t.Fatalf("not a file: %s", cur)
		}
		return cur
	}
	t.Fatalf("project root not found: reached: %s", home)
	return ""
}
