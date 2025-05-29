package packagemanager

import (
	"fmt"
	"path/filepath"

	"github.com/Masterminds/semver"
)

// Pnpm6Workspaces is a representation of workspace package globs found
// in pnpm-workspace.yaml
type Pnpm6Workspaces struct {
	Packages []string `yaml:"packages,omitempty"`
}

var pnpm6 = PackageManager{
	Name:                       "pnpm6",
	Slug:                       "pnpm",
	Command:                    "pnpm",
	Specfile:                   "package.json",
	Lockfile:                   "pnpm-lock.yaml",
	PackageDir:                 "node_modules",
	ArgSeparator:               []string{"--"},
	WorkspaceConfigurationPath: "pnpm-workspace.yaml",

	getWorkspaceGlobs: getPnpmWorkspaceGlobs,

	getWorkspaceIgnores: getPnpmWorkspaceIgnores,

	Matches: func(manager string, version string) (bool, error) {
		if manager != "pnpm" {
			return false, nil
		}

		v, err := semver.NewVersion(version)
		if err != nil {
			return false, fmt.Errorf("could not parse pnpm version: %w", err)
		}
		c, err := semver.NewConstraint("<7.0.0")
		if err != nil {
			return false, fmt.Errorf("could not create constraint: %w", err)
		}

		return c.Check(v), nil
	},

	detect: func(projectDirectory string, packageManager *PackageManager) (bool, error) {
		specfileExists := FileExists(filepath.Join(projectDirectory, packageManager.Specfile))
		lockfileExists := FileExists(filepath.Join(projectDirectory, packageManager.Lockfile))

		if !(specfileExists && lockfileExists) {
			return false, nil
		}

		// If pnpm is not installed, we can't determine the version
		// In this case, let pnpm handle the detection (modern default)
		if !packageManager.IsInstalled() {
			return false, nil
		}

		version, err := packageManager.GetStandardVersion()
		if err != nil {
			return false, nil // If we can't get version, let pnpm handle it
		}

		// Use the same logic as Matches to determine if this is pnpm6
		v, err := semver.NewVersion(version)
		if err != nil {
			return false, nil
		}
		c, err := semver.NewConstraint("<7.0.0")
		if err != nil {
			return false, nil
		}

		return c.Check(v), nil
	},

	canPrune: func(cwd string) (bool, error) {
		return true, nil
	},

	// @FIXME unsuported lockfile
	// UnmarshalLockfile: func(contents []byte) (lockfile.Lockfile, error) {
	// 	return lockfile.DecodeNpmLockfile(contents)
	// },

}
