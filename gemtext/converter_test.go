package gemtext_test

import (
	"path/filepath"
	"testing"

	"github.com/fhofherr/mnml/gemtext"
	"github.com/fhofherr/mnml/internal/testsupport"
)

func TestFromAlmostGemtext(t *testing.T) {
	testdataDir := filepath.Join("testdata", t.Name())
	tests := testsupport.FindConverterTests(t, testdataDir, "*.agmi", gemtext.FromAlmostGemtext)
	tests = append(tests, &testsupport.ConverterTest{
		Name:         "Convert the Almost Gemtext spec to Gemtext",
		InputFile:    filepath.Join(testsupport.ProjectRoot(t), "docs", "almost_gemtext.agmi"),
		ExpectedFile: filepath.Join(testdataDir, "almost_gemtext.agmi.golden"),
		Converter:    gemtext.FromAlmostGemtext,
	})

	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, tt.Run)
	}
}
