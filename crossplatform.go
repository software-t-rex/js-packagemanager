package packagemanager

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// CrossPlatform provides cross-platform utilities for package manager detection
type CrossPlatform struct{}

// GetExecutableName returns the executable name with proper extension for the current platform
func (cp CrossPlatform) GetExecutableName(command string) string {
	if runtime.GOOS == "windows" {
		// Check if it already has an extension
		if !strings.Contains(command, ".") {
			return command + ".exe"
		}
	}
	return command
}

// FindExecutable searches for an executable in PATH with cross-platform support
func (cp CrossPlatform) FindExecutable(command string) (string, error) {
	// Try with platform-specific extension first
	execName := cp.GetExecutableName(command)
	
	// Use exec.LookPath which handles PATH searching cross-platform
	path, err := exec.LookPath(execName)
	if err == nil {
		return path, nil
	}
	
	// If that fails and we're on Windows, try without extension
	if runtime.GOOS == "windows" && execName != command {
		return exec.LookPath(command)
	}
	
	return "", err
}

// IsExecutableAvailable checks if a command is available in the system PATH
func (cp CrossPlatform) IsExecutableAvailable(command string) bool {
	_, err := cp.FindExecutable(command)
	return err == nil
}

// GetPathSeparator returns the appropriate path separator for the current OS
func (cp CrossPlatform) GetPathSeparator() string {
	return string(os.PathSeparator)
}

// NormalizePath converts a path to use the correct separators for the current OS
func (cp CrossPlatform) NormalizePath(path string) string {
	return filepath.FromSlash(path)
}

// GetCommandWithArgs returns the command and args with proper cross-platform handling
func (cp CrossPlatform) GetCommandWithArgs(command string, args ...string) (string, []string) {
	execName := cp.GetExecutableName(command)
	return execName, args
}

// CreateCommand creates an exec.Cmd with proper cross-platform setup
func (cp CrossPlatform) CreateCommand(dir, command string, args ...string) *exec.Cmd {
	execName, cmdArgs := cp.GetCommandWithArgs(command, args...)
	cmd := exec.Command(execName, cmdArgs...)
	if dir != "" {
		cmd.Dir = dir
	}
	
	// Set environment variables for better cross-platform compatibility
	if runtime.GOOS == "windows" {
		// Ensure we have a proper environment on Windows
		env := os.Environ()
		cmd.Env = env
	}
	
	return cmd
}

// Global instance for easy access
var CrossPlatformUtils = CrossPlatform{}