package packagemanager

import (
	"runtime"
	"testing"

	"gotest.tools/v3/assert"
)

func TestCrossPlatformBasic(t *testing.T) {
	cp := CrossPlatformUtils

	t.Run("GetExecutableName basic test", func(t *testing.T) {
		name := cp.GetExecutableName("test")
		if runtime.GOOS == "windows" {
			assert.Equal(t, name, "test.exe")
		} else {
			assert.Equal(t, name, "test")
		}
	})

	t.Run("GetPathSeparator", func(t *testing.T) {
		sep := cp.GetPathSeparator()
		if runtime.GOOS == "windows" {
			assert.Equal(t, sep, "\\")
		} else {
			assert.Equal(t, sep, "/")
		}
	})

	t.Run("NormalizePath", func(t *testing.T) {
		// Test basic normalization
		result := cp.NormalizePath("test/path")
		// Should work on all platforms
		assert.Assert(t, len(result) > 0, "Normalized path should not be empty")
	})

	t.Run("IsExecutableAvailable - go should exist", func(t *testing.T) {
		// Since we're running tests, go should be available
		available := cp.IsExecutableAvailable("go")
		assert.Assert(t, available, "go command should be available during tests")
	})
}
