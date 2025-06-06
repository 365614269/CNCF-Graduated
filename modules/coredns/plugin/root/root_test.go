package root

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
)

func TestRoot(t *testing.T) {
	// Predefined error substrings
	parseErrContent := "Error during parsing:"
	unableToAccessErrContent := "unable to access root path"

	existingDirPath, err := getTempDirPath()
	if err != nil {
		t.Fatalf("BeforeTest: Failed to find an existing directory for testing! Error was: %v", err)
	}

	nonExistingDir := filepath.Join(existingDirPath, "highly_unlikely_to_exist_dir")

	existingFile, err := os.CreateTemp(t.TempDir(), "root_test")
	if err != nil {
		t.Fatalf("BeforeTest: Failed to create temp file for testing! Error was: %v", err)
	}
	defer func() {
		existingFile.Close()
		os.Remove(existingFile.Name())
	}()

	inaccessiblePath := getInaccessiblePath(existingFile.Name())

	tests := []struct {
		input              string
		shouldErr          bool
		expectedRoot       string // expected root, set to the controller. Empty for negative cases.
		expectedErrContent string // substring from the expected error. Empty for positive cases.
	}{
		// positive
		{
			fmt.Sprintf(`root %s`, nonExistingDir), false, nonExistingDir, "",
		},
		{
			fmt.Sprintf(`root %s`, existingDirPath), false, existingDirPath, "",
		},
		// negative
		{
			`root `, true, "", parseErrContent,
		},
		{
			fmt.Sprintf(`root %s`, inaccessiblePath), true, "", unableToAccessErrContent,
		},
		{
			fmt.Sprintf(`root {
				%s
			}`, existingDirPath), true, "", parseErrContent,
		},
	}

	for i, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		err := setup(c)
		cfg := dnsserver.GetConfig(c)

		if test.shouldErr && err == nil {
			t.Errorf("Test %d: Expected error but found %s for input %s", i, err, test.input)
		}

		if err != nil {
			if !test.shouldErr {
				t.Errorf("Test %d: Expected no error but found one for input %s. Error was: %v", i, test.input, err)
			}

			if !strings.Contains(err.Error(), test.expectedErrContent) {
				t.Errorf("Test %d: Expected error to contain: %v, found error: %v, input: %s", i, test.expectedErrContent, err, test.input)
			}
		}

		// check root only if we are in a positive test.
		if !test.shouldErr && test.expectedRoot != cfg.Root {
			t.Errorf("Root not correctly set for input %s. Expected: %s, actual: %s", test.input, test.expectedRoot, cfg.Root)
		}
	}
}

// getTempDirPath returns the path to the system temp directory. If it does not exist - an error is returned.
func getTempDirPath() (string, error) {
	tempDir := os.TempDir()
	_, err := os.Stat(tempDir)
	if err != nil {
		return "", err
	}
	return tempDir, nil
}

func getInaccessiblePath(file string) string {
	return filepath.Join("C:", "file\x00name") // null byte in filename is not allowed on Windows AND unix
}
