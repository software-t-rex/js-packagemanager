package packagemanager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/software-t-rex/packageJson"
)

// DenoConfig represents the structure of deno.json configuration file
type DenoConfig struct {
	Workspaces []string `json:"workspaces,omitempty"`
}

// parseDenoVersion extracts the version number from Deno's --version output
func parseDenoVersion(output string) string {
	// Deno version output format: "deno 2.3.3 (stable, release, x86_64-unknown-linux-gnu)"
	// We want to extract just the version number "2.3.3"
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > 0 {
		firstLine := lines[0]
		// Use regex to extract version number from "deno X.Y.Z (...)"
		re := regexp.MustCompile(`deno\s+(\d+\.\d+\.\d+(?:-[^\s\)]+)?)`)
		matches := re.FindStringSubmatch(firstLine)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return output // fallback to original output if parsing fails
}

func getDenoWorkspaceGlobs(rootpath string) ([]string, error) {
	// First try to read from deno.json
	denoConfigPath := filepath.Join(rootpath, "deno.json")
	if FileExists(denoConfigPath) {
		bytes, err := os.ReadFile(denoConfigPath)
		if err != nil {
			return nil, fmt.Errorf("deno.json: %w", err)
		}
		var denoConfig DenoConfig
		if err := json.Unmarshal(bytes, &denoConfig); err != nil {
			return nil, fmt.Errorf("deno.json: %w", err)
		}
		if len(denoConfig.Workspaces) > 0 {
			return denoConfig.Workspaces, nil
		}
	}

	// Fallback to package.json
	pkg, err := packageJson.Read(filepath.Join(rootpath, "package.json"))
	if err != nil {
		return nil, fmt.Errorf("package.json: %w", err)
	}
	if len(pkg.Workspaces) == 0 {
		return nil, fmt.Errorf("no workspaces found. packagemanager requires Deno workspaces to be defined in deno.json or package.json")
	}
	return pkg.Workspaces, nil
}

// Custom GetVersion for Deno that parses its multi-line output
func getDenoVersion() (string, error) {
	cmd := CrossPlatformUtils.CreateCommand("", "deno", "--version")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("could not detect deno version: %w", err)
	}
	return parseDenoVersion(string(out)), nil
}

// Custom GetStandardVersion for Deno
func getDenoStandardVersion() (string, error) {
	version, err := getDenoVersion()
	if err != nil {
		return "", err
	}
	return strings.Split(version, "+")[0], nil
}

var deno = PackageManager{
	Name:         "deno",
	Slug:         "deno",
	Command:      "deno",
	Specfile:     "deno.json", // Primary spec file for Deno
	Lockfile:     "deno.lock",
	PackageDir:   "node_modules",
	ArgSeparator: []string{"--"},

	getWorkspaceGlobs: getDenoWorkspaceGlobs,

	getWorkspaceIgnores: func(pm PackageManager, rootpath string) ([]string, error) {
		// Deno follows similar patterns to npm for workspace ignores
		return []string{
			"**/node_modules/**",
		}, nil
	},

	Matches: func(manager string, version string) (bool, error) {
		if manager != "deno" {
			return false, nil
		}
		// Deno 1.x and 2.x are both supported
		// No specific version constraints needed as Deno maintains good backwards compatibility
		return true, nil
	},

	detect: func(projectDirectory string, packageManager *PackageManager) (bool, error) {
		// Check for deno.json first (primary spec file)
		denoJsonExists := FileExists(filepath.Join(projectDirectory, "deno.json"))
		// Also check for package.json as fallback
		packageJsonExists := FileExists(filepath.Join(projectDirectory, "package.json"))
		lockfileExists := FileExists(filepath.Join(projectDirectory, packageManager.Lockfile))

		// Deno project can have either deno.json or package.json, and should have deno.lock
		return (denoJsonExists || packageJsonExists) && lockfileExists, nil
	},

	canPrune: func(cwd string) (bool, error) {
		return true, nil
	},

	// Custom version parsing for Deno's multi-line output
	getVersion:         getDenoVersion,
	getStandardVersion: getDenoStandardVersion,

	prunePatches: func(pkgJSON *packageJson.PackageJSON, patches []string) error {
		// Deno doesn't have built-in patch support like pnpm, but we can implement
		// basic patch pruning if needed in the future
		return nil
	},
}
