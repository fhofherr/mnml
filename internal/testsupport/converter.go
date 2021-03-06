package testsupport

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Converter defines a type that is capable from converting some input to
// some output.
type Converter func(in io.Reader, out io.Writer) error

// ConverterTest is a test case for a conversion from one file format
// to another.
//
// The default name of the converter test case is derived from the name of the
// input file. The default name can be overridden by setting the Name field.
//
// If the -update flag is passed to the test, the ConverterTest writes
// any output produced to the ExpectedFile before comparison.
type ConverterTest struct {
	Name         string // Name of the test case. Derived from InFile if empty.
	InputFile    string // Input file for the test case.
	ExpectedFile string // File containing expected output.

	Converter Converter // The converter to test.
}

// Run runs the converter test case.
func (tt *ConverterTest) Run(t *testing.T) {
	var (
		actual         bytes.Buffer
		expectedWriter io.WriteCloser
		w              io.Writer
	)
	t.Helper()

	w = &actual
	if IsUpdate() {
		var err error

		t.Logf("Update flag passed. Writing file %q", tt.ExpectedFile)
		expectedWriter, err = os.Create(tt.ExpectedFile)
		if !assert.NoError(t, err) {
			return
		}
		w = io.MultiWriter(&actual, expectedWriter)
	}

	r, err := os.Open(tt.InputFile)
	if !assert.NoError(t, err) {
		return
	}
	defer r.Close()

	// Perform test by calling Converter. Fail if converter returns an
	// error.
	if !assert.NoError(t, tt.Converter(r, w)) {
		return
	}

	if expectedWriter != nil {
		// Close to ensure the expected file is completely written before we
		// read from it.
		expectedWriter.Close()
	}

	expected, err := os.ReadFile(tt.ExpectedFile)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotEmpty(t, expected, "ExpectedFile was empty") {
		return
	}
	assert.Equal(t, string(expected), actual.String())
}

// FindConverterTests searches dir for files matching glob and creates
// a converter test from them.
//
// All expectation files created by converter tests have the name of the
// respective input file with the suffix .golden appended to them.
func FindConverterTests(t *testing.T, dir, glob string, c Converter) []*ConverterTest {
	var tests []*ConverterTest // nolint: prealloc

	entries, err := os.ReadDir(dir)
	if !assert.NoError(t, err) {
		return tests
	}

	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}
		matched, err := filepath.Match(glob, entry.Name())
		if !assert.NoError(t, err) {
			return tests
		}
		if !matched {
			continue
		}
		tt := &ConverterTest{
			InputFile:    filepath.Join(dir, entry.Name()),
			ExpectedFile: filepath.Join(dir, entry.Name()+".golden"),
			Converter:    c,
		}
		tt.Name = fmt.Sprintf("Convert %q to %q", tt.InputFile, tt.ExpectedFile)
		tests = append(tests, tt)
	}
	return tests
}
