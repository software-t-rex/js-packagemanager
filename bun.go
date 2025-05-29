package packagemanager

import (
	"fmt"
	"path/filepath"

	"github.com/software-t-rex/packageJson"
)

var bun = PackageManager{
	Name:         "bun",
	Slug:         "bun",
	Command:      "bun",
	Specfile:     "package.json",
	Lockfile:     "bun.lockb",
	PackageDir:   "node_modules",
	ArgSeparator: []string{"--"},

	getWorkspaceGlobs: func(rootpath string) ([]string, error) {
		pkg, err := packageJson.Read(filepath.Join(rootpath, "package.json"))
		if err != nil {
			return nil, fmt.Errorf("package.json: %w", err)
		}
		if len(pkg.Workspaces) == 0 {
			return nil, fmt.Errorf("package.json: no workspaces found. packagemanager requires Bun workspaces to be defined in the root package.json")
		}
		return pkg.Workspaces, nil
	},

	getWorkspaceIgnores: func(pm PackageManager, rootpath string) ([]string, error) {
		// Bun follows similar patterns to npm for workspace ignores
		return []string{
			"**/node_modules/**",
		}, nil
	},

	Matches: func(manager string, version string) (bool, error) {
		return manager == "bun", nil
	},

	detect: func(projectDirectory string, packageManager *PackageManager) (bool, error) {
		specfileExists := FileExists(filepath.Join(projectDirectory, packageManager.Specfile))
		lockfileExists := FileExists(filepath.Join(projectDirectory, packageManager.Lockfile))

		return (specfileExists && lockfileExists), nil
	},

	canPrune: func(cwd string) (bool, error) {
		return true, nil
	},

	prunePatches: func(pkgJSON *packageJson.PackageJSON, patches []string) error {
		// Bun doesn't have built-in patch support like pnpm, but we can implement
		// basic patch pruning if needed in the future
		return nil
	},
}
