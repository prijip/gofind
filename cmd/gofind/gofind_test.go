package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlags(t *testing.T) {
	flag.Set("config", "testdata/config.json")
	parseFlags()

	assert.Equal(t, "./testdata/output", config.OutputDirectory)
}

func getFileList(rootdir string) (fileList []string, err error) {
	err = filepath.Walk(rootdir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fileName, err := filepath.Rel(rootdir, path)
		if err != nil {
			return err
		}
		fileList = append(fileList, fileName)
		return nil
	})

	return
}

func normalizeNL(data []byte) []byte {
	// replace CR LF \r\n (windows) with LF \n (unix)
	data = bytes.Replace(data, []byte{'\r', '\n'}, []byte{'\n'}, -1)
	// replace CF \r (mac) with LF \n (unix)
	data = bytes.Replace(data, []byte{'\r'}, []byte{'\n'}, -1)

	return data
}

// assertFilesEqualIgnoreNL compares the contents of two text files
// Ignores how new line is represented
func assertFilesEqualIgnoreNL(t *testing.T, path1, path2 string) {
	info1, err := os.Stat(path1)
	assert.NoError(t, err)

	info2, err := os.Stat(path2)
	assert.NoError(t, err)

	// If both are directories, just return
	// Make sure if one name corresponds to a directory, the other one is a directory too
	if info1.IsDir() {
		assert.True(t, info2.IsDir())
		return
	}
	// First file is not a directory, so the second one should not be too
	assert.False(t, info2.IsDir())

	content1, err := ioutil.ReadFile(path1)
	assert.NoError(t, err)

	content1 = normalizeNL(content1)

	content2, err := ioutil.ReadFile(path2)
	assert.NoError(t, err)

	content2 = normalizeNL(content2)

	r := bytes.Compare(content1, content2)
	assert.Equal(t, 0, r, "No Match - '%s', '%s'", path1, path2)
}

func TestSearchReplace(t *testing.T) {
	flag.Set("config", "testdata/config.yaml")
	parseFlags()

	doFind()

	expectedFiles, err := getFileList("./testdata/expected_output")
	assert.NoError(t, err)

	actualFiles, err := getFileList("./testdata/output")
	assert.NoError(t, err)

	assert.ElementsMatch(t, expectedFiles, actualFiles)

	for _, fileName := range actualFiles {
		path1 := filepath.Join("./testdata/expected_output", fileName)
		path2 := filepath.Join("./testdata/output", fileName)

		assertFilesEqualIgnoreNL(t, path1, path2)
	}
}
