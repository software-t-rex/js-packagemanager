package packagemanager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/software-t-rex/monospace/packageJson"
	"sigs.k8s.io/yaml"
)

// PnpmWorkspaces is a representation of workspace package globs found
// in pnpm-workspace.yaml
type PnpmWorkspaces struct {
	Packages []string `yaml:"packages,omitempty"`
}

func readPnpmWorkspacePackages(workspaceFile string) ([]string, error) {
	bytes, err := os.ReadFile(workspaceFile)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", workspaceFile, err)
	}
	var pnpmWorkspaces PnpmWorkspaces
	if err := yaml.Unmarshal(bytes, &pnpmWorkspaces); err != nil {
		return nil, fmt.Errorf("%v: %w", workspaceFile, err)
	}
	return pnpmWorkspaces.Packages, nil
}

func getPnpmWorkspaceGlobs(rootpath string) ([]string, error) {
	pkgGlobs, err := readPnpmWorkspacePackages(filepath.Join(rootpath, "pnpm-workspace.yaml"))
	if err != nil {
		return nil, err
	}

	if len(pkgGlobs) == 0 {
		return nil, fmt.Errorf("pnpm-workspace.yaml: no packages found. packagemanager requires pnpm workspaces and thus packages to be defined in the root pnpm-workspace.yaml")
	}

	filteredPkgGlobs := []string{}
	for _, pkgGlob := range pkgGlobs {
		if !strings.HasPrefix(pkgGlob, "!") {
			filteredPkgGlobs = append(filteredPkgGlobs, pkgGlob)
		}
	}
	return filteredPkgGlobs, nil
}

func getPnpmWorkspaceIgnores(pm PackageManager, rootpath string) ([]string, error) {
	// Matches upstream values:
	// function: https://github.com/pnpm/pnpm/blob/d99daa902442e0c8ab945143ebaf5cdc691a91eb/packages/find-packages/src/index.ts#L27
	// key code: https://github.com/pnpm/pnpm/blob/d99daa902442e0c8ab945143ebaf5cdc691a91eb/packages/find-packages/src/index.ts#L30
	// call site: https://github.com/pnpm/pnpm/blob/d99daa902442e0c8ab945143ebaf5cdc691a91eb/packages/find-workspace-packages/src/index.ts#L32-L39
	ignores := []string{
		"**/node_modules/**",
		"**/bower_components/**",
	}
	pkgGlobs, err := readPnpmWorkspacePackages(filepath.Join(rootpath, "pnpm-workspace.yaml"))
	if err != nil {
		return nil, err
	}
	for _, pkgGlob := range pkgGlobs {
		if strings.HasPrefix(pkgGlob, "!") {
			ignores = append(ignores, pkgGlob[1:])
		}
	}
	return ignores, nil
}

var nodejsPnpm = PackageManager{
	Name:       "nodejs-pnpm",
	Slug:       "pnpm",
	Command:    "pnpm",
	Specfile:   "package.json",
	Lockfile:   "pnpm-lock.yaml",
	PackageDir: "node_modules",
	// pnpm v7+ changed their handling of '--'. We no longer need to pass it to pass args to
	// the script being run, and in fact doing so will cause the '--' to be passed through verbatim,
	// potentially breaking scripts that aren't expecting it.
	// We are allowed to use nil here because ArgSeparator already has a type, so it's a typed nil,
	// This could just as easily be []string{}, but the style guide says to prefer
	// nil for empty slices.
	ArgSeparator:               nil,
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
		c, err := semver.NewConstraint(">=7.0.0")
		if err != nil {
			return false, fmt.Errorf("could not create constraint: %w", err)
		}

		return c.Check(v), nil
	},

	detect: func(projectDirectory string, packageManager *PackageManager) (bool, error) {
		specfileExists := FileExists(filepath.Join(projectDirectory, packageManager.Specfile))
		lockfileExists := FileExists(filepath.Join(projectDirectory, packageManager.Lockfile))

		return (specfileExists && lockfileExists), nil
	},

	canPrune: func(cwd string) (bool, error) {
		return true, nil
	},

	// @FIXME unsuported lockfile
	// UnmarshalLockfile: func(contents []byte) (lockfile.Lockfile, error) {
	// 	return lockfile.DecodeNpmLockfile(contents)
	// },

	prunePatches: func(pkgJSON *packageJson.PackageJSON, patches []string) error {
		return pnpmPrunePatches(pkgJSON, patches)
	},
}

func pnpmPrunePatches(pkgJSON *packageJson.PackageJSON, patches []string) error {
	pkgJSON.Mu.Lock()
	defer pkgJSON.Mu.Unlock()

	keysToDelete := []string{}
	pnpmConfig, ok := pkgJSON.Pnpm.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid structure for pnpm field in package.json")
	}
	patchedDependencies, ok := pnpmConfig["patchedDependencies"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid structure for patchedDependencies field in package.json")
	}

	for dependency, untypedPatch := range patchedDependencies {
		patch, ok := untypedPatch.(string)
		if !ok {
			return fmt.Errorf("expected only strings in patchedDependencies. Got %v", untypedPatch)
		}

		inPatches := false

		for _, wantedPatch := range patches {
			if wantedPatch == patch {
				inPatches = true
				break
			}
		}

		if !inPatches {
			keysToDelete = append(keysToDelete, dependency)
		}
	}

	for _, key := range keysToDelete {
		delete(patchedDependencies, key)
	}

	return nil
}
