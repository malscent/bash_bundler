// white box unit testing
//nolint:testpackage
package bundler

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsShebangReturnsTrueForShebang(t *testing.T) {
	t.Parallel()

	val := isShebang("#!/usr/bin/bash")
	assert.True(t, val)
}

func TestIsShebangReturnsFalseForNotShebang(t *testing.T) {
	t.Parallel()

	val := isShebang("# This is not a shebang")
	assert.False(t, val)
}

func TestDeleteEmptyDoesNotDeleteNonEmpty(t *testing.T) {
	t.Parallel()

	values := []string{"one", "two", "three"}
	newValues := deleteEmpty(values)

	assert.Len(t, newValues, 3)
}

func TestDeleteEmptyRemovesSingleEmptyEntry(t *testing.T) {
	t.Parallel()

	values := []string{"one", "   ", "three"}
	newValues := deleteEmpty(values)

	assert.Len(t, newValues, 2)
}

func TestDeleteEmptyRemovesMultipleEmptyEntries(t *testing.T) {
	t.Parallel()

	values := []string{"    ", "   ", "three"}
	newValues := deleteEmpty(values)

	assert.Len(t, newValues, 1)
}

func TestDeleteEmptyRemovesMultipleEmptyWithTabs(t *testing.T) {
	t.Parallel()

	values := []string{"    ", "\t ", "three"}
	newValues := deleteEmpty(values)

	assert.Len(t, newValues, 1)
}

func TestTrimQuotesRemovesDoubleQuotes(t *testing.T) {
	t.Parallel()

	value := "\"This is a string with double quotes\""
	newValue := trimQuotes(value)

	assert.Equal(t, "This is a string with double quotes", newValue)
}

func TestTrimQuotesRemovesSingleQuotes(t *testing.T) {
	t.Parallel()

	value := "'This is a string with single quotes'"
	newValue := trimQuotes(value)

	assert.Equal(t, "This is a string with single quotes", newValue)
}

func TestTrimQuotesDoesNotRemovesPreceedingSingleQuote(t *testing.T) {
	t.Parallel()

	value := "\"This is a string with double quotes"
	newValue := trimQuotes(value)

	assert.Equal(t, "\"This is a string with double quotes", newValue)
}

func TestTrimQuotesDoesNotRemovesTrailingSingleQuote(t *testing.T) {
	t.Parallel()

	value := "This is a string with double quotes\""
	newValue := trimQuotes(value)

	assert.Equal(t, "This is a string with double quotes\"", newValue)
}

func getExpected(path string) string {
	in, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}

	return string(in)
}

func compareValues(expected string, generated string) bool {
	expectedSplit := strings.Split(expected, "\n")
	generatedSplit := strings.Split(generated, "\n")

	if len(expectedSplit) != len(generatedSplit) {
		return false
	}

	for i := 0; i < len(expectedSplit); i++ {
		if strings.HasPrefix(expectedSplit[i], "#") && strings.HasPrefix(generatedSplit[i], "#") {
			continue
		}

		if expectedSplit[i] != generatedSplit[i] {
			return false
		}
	}

	return true
}

func TestGenerateSameLevelSources(t *testing.T) {
	t.Parallel()

	_, filename, _, ok := runtime.Caller(0)

	assert.True(t, ok)

	sourcePath := path.Dir(filename)

	expectedPath := filepath.Join(sourcePath, "testdata", "/expected/same_level.golden")
	generatedPath := filepath.Join(sourcePath, "testdata", "same_level.sh")

	expected := getExpected(expectedPath)
	generated, err := Bundle(generatedPath, true)

	assert.Nil(t, err)
	assert.True(t, compareValues(expected, generated))
}

func TestGenerateSameLevelSourcesMinified(t *testing.T) {
	t.Parallel()

	_, filename, _, ok := runtime.Caller(0)

	assert.True(t, ok)

	sourcePath := path.Dir(filename)

	expectedPath := filepath.Join(sourcePath, "testdata", "/expected/same_level_min.golden")
	generatedPath := filepath.Join(sourcePath, "testdata", "same_level.sh")

	expected := getExpected(expectedPath)
	generated, err := Bundle(generatedPath, true)
	assert.Nil(t, err)

	generated, err = Minify(generated)
	assert.Nil(t, err)
	assert.True(t, compareValues(expected, generated))
}

func TestGenerateNestedSources(t *testing.T) {
	t.Parallel()

	_, filename, _, ok := runtime.Caller(0)

	assert.True(t, ok)

	sourcePath := path.Dir(filename)

	expectedPath := filepath.Join(sourcePath, "testdata", "/expected/nested.golden")
	generatedPath := filepath.Join(sourcePath, "testdata", "main_nested.sh")

	expected := getExpected(expectedPath)
	generated, err := Bundle(generatedPath, true)

	assert.Nil(t, err)
	assert.True(t, compareValues(expected, generated))
}

func TestGenerateNestedSourcesMinified(t *testing.T) {
	t.Parallel()

	_, filename, _, ok := runtime.Caller(0)

	assert.True(t, ok)

	sourcePath := path.Dir(filename)

	expectedPath := filepath.Join(sourcePath, "testdata", "/expected/nested_min.golden")
	generatedPath := filepath.Join(sourcePath, "testdata", "main_nested.sh")

	expected := getExpected(expectedPath)
	generated, err := Bundle(generatedPath, true)
	assert.Nil(t, err)

	generated, err = Minify(generated)
	assert.Nil(t, err)
	assert.True(t, compareValues(expected, generated))
}

func TestGenerateEmbeddedSources(t *testing.T) {
	t.Parallel()

	_, filename, _, ok := runtime.Caller(0)

	assert.True(t, ok)

	sourcePath := path.Dir(filename)

	expectedPath := filepath.Join(sourcePath, "testdata", "/expected/embedded.golden")
	generatedPath := filepath.Join(sourcePath, "testdata", "main_embedded.sh")

	expected := getExpected(expectedPath)
	generated, err := Bundle(generatedPath, true)

	assert.Nil(t, err)
	assert.True(t, compareValues(expected, generated))
}

func TestGenerateEmbeddedSourcesMinified(t *testing.T) {
	t.Parallel()

	_, filename, _, ok := runtime.Caller(0)

	assert.True(t, ok)

	sourcePath := path.Dir(filename)

	expectedPath := filepath.Join(sourcePath, "testdata", "/expected/embedded_min.golden")
	generatedPath := filepath.Join(sourcePath, "testdata", "main_embedded.sh")

	expected := getExpected(expectedPath)
	generated, err := Bundle(generatedPath, true)
	assert.Nil(t, err)

	generated, err = Minify(generated)
	assert.Nil(t, err)
	assert.True(t, compareValues(expected, generated))
}

func TestGenerateParentSources(t *testing.T) {
	t.Parallel()

	_, filename, _, ok := runtime.Caller(0)

	assert.True(t, ok)

	sourcePath := path.Dir(filename)

	expectedPath := filepath.Join(sourcePath, "testdata", "/expected/parent.golden")
	generatedPath := filepath.Join(sourcePath, "testdata", "nested", "main_parent.sh")

	expected := getExpected(expectedPath)
	generated, err := Bundle(generatedPath, true)

	assert.Nil(t, err)
	assert.True(t, compareValues(expected, generated))
}

func TestGenerateParentSourcesMinified(t *testing.T) {
	t.Parallel()

	_, filename, _, ok := runtime.Caller(0)

	assert.True(t, ok)

	sourcePath := path.Dir(filename)

	expectedPath := filepath.Join(sourcePath, "testdata", "/expected/parent_min.golden")
	generatedPath := filepath.Join(sourcePath, "testdata", "nested", "main_parent.sh")

	expected := getExpected(expectedPath)
	generated, err := Bundle(generatedPath, true)
	assert.Nil(t, err)

	generated, err = Minify(generated)
	assert.Nil(t, err)
	assert.True(t, compareValues(expected, generated))
}
