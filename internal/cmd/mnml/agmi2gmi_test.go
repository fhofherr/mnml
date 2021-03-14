package mnml_test

import (
	"path/filepath"
	"testing"

	"github.com/fhofherr/mnml/internal/cmd/mnml"
	"github.com/fhofherr/mnml/internal/testsupport"
	"github.com/stretchr/testify/assert"
)

func TestAGMI2GMICmd(t *testing.T) {
	tempDir, cleanUp := testsupport.MkdirTemp(t)
	defer cleanUp()

	srcFile := filepath.Join(testsupport.ProjectRoot(t), "docs", "almost_gemtext.agmi")
	destFile := filepath.Join(tempDir, "almost_gemtext.gmi")

	cmd := mnml.New()
	cmd.SetArgs([]string{"agmi2gmi", "--output", destFile, srcFile})
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.FileExists(t, destFile)
}
