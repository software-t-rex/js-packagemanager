package packagemanager

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestCrossPlatform_GetExecutableName(t *testing.T) {
	cp := CrossPlatform{}

	tests := []struct {
		name     string
		command  string
		expected string
	}{
		{
			name:     "simple command",
			command:  "npm",
			expected: getExpectedExecutableName("npm"),
		},
		{
			name:     "command with extension",
			command:  "yarn.cmd",
			expected: "yarn.cmd",
		},
		{
			name:     "command with path",
			command:  "/usr/bin/node",
			expected: getExpectedExecutableName("/usr/bin/node"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cp.GetExecutableName(tt.command)
			assert.Equal(t, result, tt.expected)
		})
	}
}

func TestCrossPlatform_IsExecutableAvailable(t *testing.T) {
	cp := CrossPlatform{}

	// Test with a command that should be available on most systems
	// Use 'go' since we're running go tests
	available := cp.IsExecutableAvailable("go")
	assert.Assert(t, available, "go command should be available")

	// Test with a command that definitely doesn't exist
	notAvailable := cp.IsExecutableAvailable("definitely-not-a-real-command-12345")
	assert.Assert(t, !notAvailable, "fake command should not be available")
}

func TestCrossPlatform_GetPathSeparator(t *testing.T) {
	cp := CrossPlatform{}
	separator := cp.GetPathSeparator()

	if runtime.GOOS == "windows" {
		assert.Equal(t, separator, "\\")
	} else {
		assert.Equal(t, separator, "/")
	}
}

func TestCrossPlatform_NormalizePath(t *testing.T) {
	cp := CrossPlatform{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "unix style path",
			input:    "path/to/file",
			expected: getExpectedPath("path/to/file"),
		},
		{
			name:     "windows style path",
			input:    "path\\to\\file",
			expected: getExpectedPath("path\\to\\file"),
		},
		{
			name:     "mixed style path",
			input:    "path/to\\file",
			expected: getExpectedPath("path/to\\file"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cp.NormalizePath(tt.input)
			assert.Equal(t, result, tt.expected)
		})
	}
}

func TestCrossPlatform_CreateCommand(t *testing.T) {
	cp := CrossPlatform{}

	cmd := cp.CreateCommand("/test/dir", "npm", "install")

	assert.Equal(t, cmd.Dir, "/test/dir")
	assert.Assert(t, len(cmd.Args) >= 2, "Command should have at least 2 args")

	// On exec.Command, args[0] is the command name/path
	expectedCmd := getExpectedExecutableName("npm")
	actualCmd := filepath.Base(cmd.Args[0])

	assert.Equal(t, actualCmd, expectedCmd, "Command executable should match expected")
	assert.Equal(t, cmd.Args[1], "install", "First argument should be 'install'")
}

// Helper functions for platform-specific expectations
func getExpectedExecutableName(command string) string {
	if runtime.GOOS == "windows" && !contains(command, ".") {
		return command + ".exe"
	}
	return command
}

func getExpectedPath(input string) string {
	// Use filepath.FromSlash which is what NormalizePath uses
	return filepath.FromSlash(input)
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
