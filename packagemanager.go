// Adapted from https://github.com/replit/upm
// Copyright (c) 2019 Neoreason d/b/a Repl.it. All rights reserved.
// SPDX-License-Identifier: MIT
// By the turbo team for turborepo which is licensed the MPL v2.0 license
// https://github.com/vercel/turbo/tree/368b715/cli/internal/packagemanager
// This version remove any dependency on turborepo internals and refactor
// a bit of code to make this more generic package to re-use in other application
// As the License in the original file seems to be MIT all modification made will
// be under the MIT license too.

package packagemanager

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/software-t-rex/packageJson"
)

// PackageManager is an abstraction across package managers
type PackageManager struct {
	// The descriptive name of the Package Manager.
	Name string

	// The unique identifier of the Package Manager.
	Slug string

	// The command used to invoke the Package Manager.
	Command string

	// The location of the package spec file used by the Package Manager.
	Specfile string

	// The location of the package lock file used by the Package Manager.
	Lockfile string

	// The directory in which package assets are stored by the Package Manager.
	PackageDir string

	// The location of the file that defines the workspace. Empty if workspaces defined in package.json
	WorkspaceConfigurationPath string

	// The separator that the Package Manger uses to identify arguments that
	// should be passed through to the underlying script.
	ArgSeparator []string

	// Return the list of workspace glob
	getWorkspaceGlobs func(rootpath string) ([]string, error)

	// Return the list of workspace ignore globs
	getWorkspaceIgnores func(pm PackageManager, rootpath string) ([]string, error)

	// Detect if Turbo knows how to produce a pruned workspace for the project
	canPrune func(cwd string) (bool, error)

	// Test a manager and version tuple to see if it is the Package Manager.
	Matches func(manager string, version string) (bool, error)

	// Detect if the project is using the Package Manager by inspecting the system.
	detect func(projectDirectory string, packageManager *PackageManager) (bool, error)

	// @FIXME missing Lockfile support
	// Read a lockfile for a given package manager
	// UnmarshalLockfile func(contents []byte) (lockfile.Lockfile, error)

	// Prune the given pkgJSON to only include references to the given patches
	prunePatches func(pkgJSON *packageJson.PackageJSON, patches []string) error
}

var packageManagers = []PackageManager{
	nodejsYarn,
	nodejsBerry,
	nodejsNpm,
	nodejsPnpm,
	nodejsPnpm6,
}

var (
	packageManagerPattern = `(npm|pnpm|yarn)@(\d+)\.\d+\.\d+(-.+)?`
	packageManagerRegex   = regexp.MustCompile(packageManagerPattern)
)

// ParsePackageManagerString takes a package manager version string parses it into consituent components
func ParsePackageManagerString(packageManager string) (manager string, version string, err error) {
	match := packageManagerRegex.FindString(packageManager)
	if len(match) == 0 {
		return "", "", fmt.Errorf("we could not parse packageManager field in package.json, expected: %s, received: %s", packageManagerPattern, packageManager)
	}

	return strings.Split(match, "@")[0], strings.Split(match, "@")[1], nil
}

// GetPackageManager attempts all methods for identifying the package manager in use.
func GetPackageManager(projectDirectory string, pkg *packageJson.PackageJSON) (packageManager *PackageManager, err error) {
	result, _ := GetPackageManagerFromString(pkg.PackageManager)
	if result != nil {
		return result, nil
	}

	return detectPackageManager(projectDirectory)
}

func GetPackageManagerFromString(packageManagerStr string) (packageManager *PackageManager, err error) {
	if packageManagerStr == "" {
		return nil, fmt.Errorf("no package manager specified")
	}
	manager, version, err := ParsePackageManagerString(packageManagerStr)
	if err != nil {
		return nil, err
	}
	for _, packageManager := range packageManagers {
		isResponsible, err := packageManager.Matches(manager, version)
		if isResponsible && (err == nil) {
			return &packageManager, nil
		}
	}
	return nil, fmt.Errorf("we didn't find a matching package manager for '%s'", packageManagerStr)
}

// detectPackageManager attempts to detect the package manager by inspecting the project directory state.
func detectPackageManager(projectDirectory string) (packageManager *PackageManager, err error) {
	for _, packageManager := range packageManagers {
		isResponsible, err := packageManager.detect(projectDirectory, &packageManager)
		if err != nil {
			return nil, err
		}
		if isResponsible {
			return &packageManager, nil
		}
	}

	return nil, fmt.Errorf("we did not detect an in-use package manager for your project. Please set the \"packageManager\" property in your root package.json (https://nodejs.org/api/packages.html#packagemanager)")
}

// GetWorkspaces returns the list of package.json files for the current mono[space|repo].
func (pm PackageManager) GetWorkspaces(rootpath string, relativePath bool) ([]string, error) {
	globs, err := pm.getWorkspaceGlobs(rootpath)
	if err != nil {
		return nil, err
	}

	justJsons := make([]string, len(globs))
	for i, space := range globs {
		justJsons[i] = filepath.Join(space, "package.json")
	}

	ignores, err := pm.getWorkspaceIgnores(pm, rootpath)
	if err != nil {
		return nil, err
	}

	// f, err := globby.GlobFiles(rootpath, justJsons, ignores)
	fs := os.DirFS(rootpath)
	var res []string
	for _, glob := range justJsons {
		founds, err := doublestar.Glob(fs, glob)
		if err != nil {
			return nil, err
		}
		res = append(res, founds...)
	}
	for _, glob := range ignores {
		for i, path := range res {
			match, err := doublestar.Match(path, glob)
			if err != nil {
				return nil, err
			}
			if match {
				res = append(res[:i], res[i+1:]...)
			}
		}
	}

	// make res fullpath
	if !relativePath {
		for i, path := range res {
			res[i] = filepath.Join(rootpath, path)
		}
	}

	return res, nil
}

// GetWorkspaceIgnores returns an array of globs not to search for workspaces.
func (pm PackageManager) GetWorkspaceIgnores(rootpath string) ([]string, error) {
	return pm.getWorkspaceIgnores(pm, rootpath)
}

// CanPrune returns if we can produce a pruned workspace. Can error if fs issues occur
func (pm PackageManager) CanPrune(projectDirectory string) (bool, error) {
	if pm.canPrune != nil {
		return pm.canPrune(projectDirectory)
	}
	return false, nil
}

// @FIXME missing lockfile support
// ReadLockfile will read the applicable lockfile into memory
// func (pm PackageManager) ReadLockfile(projectDirectory string) (lockfile.Lockfile, error) {
// 	if pm.UnmarshalLockfile == nil {
// 		return nil, nil
// 	}
// 	contents, err := os.ReadFile(filepath.Join(projectDirectory, pm.Lockfile))
// 	if err != nil {
// 		return nil, fmt.Errorf("reading %s: %w", pm.Lockfile, err)
// 	}
// 	return pm.UnmarshalLockfile(contents)
// }

// PrunePatchedPackages will alter the provided pkgJSON to only reference the provided patches
func (pm PackageManager) PrunePatchedPackages(pkgJSON *packageJson.PackageJSON, patches []string) error {
	if pm.prunePatches != nil {
		return pm.prunePatches(pkgJSON, patches)
	}
	return nil
}

// YarnRC Represents contents of .yarnrc.yml
type YarnRC struct {
	NodeLinker string `yaml:"nodeLinker"`
}

func FileExists(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && !info.IsDir()
}

func PathExists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

func hasFile(name, dir string) (bool, error) {
	files, err := os.ReadDir(dir)

	if err != nil {
		return false, err
	}

	for _, f := range files {
		if name == f.Name() {
			return true, nil
		}
	}

	return false, nil
}

func FindupFrom(name, dir string) (string, error) {
	for {
		found, err := hasFile(name, dir)

		if err != nil {
			return "", err
		}

		if found {
			return filepath.Join(dir, name), nil
		}

		parent := filepath.Dir(dir)

		if parent == dir {
			return "", nil
		}

		dir = parent
	}
}
