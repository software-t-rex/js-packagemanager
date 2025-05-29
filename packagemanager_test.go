package packagemanager

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/software-t-rex/packageJson"

	"gotest.tools/v3/assert"
)

func TestParsePackageManagerString(t *testing.T) {
	tests := []struct {
		name           string
		packageManager string
		wantManager    string
		wantVersion    string
		wantErr        bool
	}{
		{
			name:           "errors with a tag version",
			packageManager: "npm@latest",
			wantManager:    "",
			wantVersion:    "",
			wantErr:        true,
		},
		{
			name:           "errors with no version",
			packageManager: "npm",
			wantManager:    "",
			wantVersion:    "",
			wantErr:        true,
		},
		{
			name:           "requires fully-qualified semver versions (one digit)",
			packageManager: "npm@1",
			wantManager:    "",
			wantVersion:    "",
			wantErr:        true,
		},
		{
			name:           "requires fully-qualified semver versions (two digits)",
			packageManager: "npm@1.2",
			wantManager:    "",
			wantVersion:    "",
			wantErr:        true,
		},
		{
			name:           "supports custom labels",
			packageManager: "npm@1.2.3-alpha.1",
			wantManager:    "npm",
			wantVersion:    "1.2.3-alpha.1",
			wantErr:        false,
		},
		{
			name:           "only supports specified package managers",
			packageManager: "pip@1.2.3",
			wantManager:    "",
			wantVersion:    "",
			wantErr:        true,
		},
		{
			name:           "supports npm",
			packageManager: "npm@0.0.1",
			wantManager:    "npm",
			wantVersion:    "0.0.1",
			wantErr:        false,
		},
		{
			name:           "supports pnpm",
			packageManager: "pnpm@0.0.1",
			wantManager:    "pnpm",
			wantVersion:    "0.0.1",
			wantErr:        false,
		},
		{
			name:           "supports yarn",
			packageManager: "yarn@111.0.1",
			wantManager:    "yarn",
			wantVersion:    "111.0.1",
			wantErr:        false,
		},
		{
			name:           "supports bun",
			packageManager: "bun@1.0.0",
			wantManager:    "bun",
			wantVersion:    "1.0.0",
			wantErr:        false,
		},
		{
			name:           "supports deno",
			packageManager: "deno@1.36.0",
			wantManager:    "deno",
			wantVersion:    "1.36.0",
			wantErr:        false,
		},
		{
			name:           "supports deno 2.x",
			packageManager: "deno@2.0.0",
			wantManager:    "deno",
			wantVersion:    "2.0.0",
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotManager, gotVersion, err := ParsePackageManagerString(tt.packageManager)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePackageManagerString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotManager != tt.wantManager {
				t.Errorf("ParsePackageManagerString() got manager = %v, want manager %v", gotManager, tt.wantManager)
			}
			if gotVersion != tt.wantVersion {
				t.Errorf("ParsePackageManagerString() got version = %v, want version %v", gotVersion, tt.wantVersion)
			}
		})
	}
}

func TestGetPackageManager(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NilError(t, err, "os.Getwd")
	tests := []struct {
		name             string
		projectDirectory string
		pkg              *packageJson.PackageJSON
		want             string
		wantErr          bool
	}{
		{
			name:             "finds npm from a package manager string",
			projectDirectory: cwd,
			pkg:              &packageJson.PackageJSON{PackageManager: "npm@1.2.3"},
			want:             "npm",
			wantErr:          false,
		},
		{
			name:             "finds pnpm6 from a package manager string",
			projectDirectory: cwd,
			pkg:              &packageJson.PackageJSON{PackageManager: "pnpm@1.2.3"},
			want:             "pnpm6",
			wantErr:          false,
		},
		{
			name:             "finds pnpm from a package manager string",
			projectDirectory: cwd,
			pkg:              &packageJson.PackageJSON{PackageManager: "pnpm@7.8.9"},
			want:             "pnpm",
			wantErr:          false,
		},
		{
			name:             "finds yarn from a package manager string",
			projectDirectory: cwd,
			pkg:              &packageJson.PackageJSON{PackageManager: "yarn@1.2.3"},
			want:             "yarn",
			wantErr:          false,
		},
		{
			name:             "finds berry from a package manager string",
			projectDirectory: cwd,
			pkg:              &packageJson.PackageJSON{PackageManager: "yarn@2.3.4"},
			want:             "berry",
			wantErr:          false,
		},
		{
			name:             "finds bun from a package manager string",
			projectDirectory: cwd,
			pkg:              &packageJson.PackageJSON{PackageManager: "bun@1.0.0"},
			want:             "bun",
			wantErr:          false,
		},
		{
			name:             "finds deno from a package manager string",
			projectDirectory: cwd,
			pkg:              &packageJson.PackageJSON{PackageManager: "deno@1.36.0"},
			want:             "deno",
			wantErr:          false,
		},
		{
			name:             "finds deno 2.x from a package manager string",
			projectDirectory: cwd,
			pkg:              &packageJson.PackageJSON{PackageManager: "deno@2.0.0"},
			want:             "deno",
			wantErr:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPackageManager, err := GetPackageManager(tt.projectDirectory, tt.pkg)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPackageManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPackageManager.Name != tt.want {
				t.Errorf("GetPackageManager() = %v, want %v", gotPackageManager.Name, tt.want)
			}
		})
	}
}

func Test_GetPackageManagerFromString(t *testing.T) {
	tests := []struct {
		name       string
		pkgMngrStr string
		want       string
		wantErr    bool
	}{
		{
			name:       "finds npm from a package manager string",
			pkgMngrStr: "npm@1.2.3",
			want:       "npm",
			wantErr:    false,
		},
		{
			name:       "finds pnpm6 from a package manager string",
			pkgMngrStr: "pnpm@1.2.3",
			want:       "pnpm6",
			wantErr:    false,
		},
		{
			name:       "finds pnpm from a package manager string",
			pkgMngrStr: "pnpm@7.8.9",
			want:       "pnpm",
			wantErr:    false,
		},
		{
			name:       "finds yarn from a package manager string",
			pkgMngrStr: "yarn@1.2.3",
			want:       "yarn",
			wantErr:    false,
		},
		{
			name:       "finds berry from a package manager string",
			pkgMngrStr: "yarn@2.3.4",
			want:       "berry",
			wantErr:    false,
		},
		{
			name:       "finds bun from a package manager string",
			pkgMngrStr: "bun@1.0.0",
			want:       "bun",
			wantErr:    false,
		},
		{
			name:       "finds deno from a package manager string",
			pkgMngrStr: "deno@1.36.0",
			want:       "deno",
			wantErr:    false,
		},
		{
			name:       "finds deno 2.x from a package manager string",
			pkgMngrStr: "deno@2.0.0",
			want:       "deno",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPackageManager, err := GetPackageManagerFromString(tt.pkgMngrStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("readPackageManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPackageManager.Name != tt.want {
				t.Errorf("readPackageManager() = %v, want %v", gotPackageManager.Name, tt.want)
			}
		})
	}
}

func Test_GetWorkspaces(t *testing.T) {
	type test struct {
		name     string
		pm       PackageManager
		rootPath string
		want     []string
		wantErr  bool
	}

	cwd, err := os.Getwd()
	assert.NilError(t, err, "os.Getcwd")
	rootPath := map[string]string{
		"npm":   filepath.Join(cwd, "testdata/with-yarn"),
		"berry": filepath.Join(cwd, "testdata/with-yarn"),
		"yarn":  filepath.Join(cwd, "testdata/with-yarn"),
		"pnpm":  filepath.Join(cwd, "testdata/basic"),
		"pnpm6": filepath.Join(cwd, "testdata/basic"),
		"bun":   filepath.Join(cwd, "testdata/with-yarn"),
		"deno":  filepath.Join(cwd, "testdata/with-yarn"),
	}

	want := map[string][]string{
		"npm": {
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/apps/docs/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/apps/web/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/eslint-config-custom/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/tsconfig/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/ui/package.json")),
		},
		"berry": {
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/apps/docs/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/apps/web/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/eslint-config-custom/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/tsconfig/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/ui/package.json")),
		},
		"yarn": {
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/apps/docs/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/apps/web/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/eslint-config-custom/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/tsconfig/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/ui/package.json")),
		},
		"pnpm": {
			filepath.ToSlash(filepath.Join(cwd, "testdata/basic/apps/docs/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/basic/apps/web/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/basic/packages/eslint-config-custom/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/basic/packages/tsconfig/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/basic/packages/ui/package.json")),
		},
		"pnpm6": {
			filepath.ToSlash(filepath.Join(cwd, "testdata/basic/apps/docs/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/basic/apps/web/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/basic/packages/eslint-config-custom/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/basic/packages/tsconfig/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/basic/packages/ui/package.json")),
		},
		"bun": {
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/apps/docs/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/apps/web/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/eslint-config-custom/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/tsconfig/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/ui/package.json")),
		},
		"deno": {
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/apps/docs/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/apps/web/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/eslint-config-custom/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/tsconfig/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/ui/package.json")),
		},
	}

	tests := make([]test, len(packageManagers))
	for i, packageManager := range packageManagers {
		tests[i] = test{
			name:     packageManager.Name,
			pm:       packageManager,
			rootPath: rootPath[packageManager.Name],
			want:     want[packageManager.Name],
			wantErr:  false,
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWorkspaces, err := tt.pm.GetWorkspaces(tt.rootPath, false)
			gotToSlash := make([]string, len(gotWorkspaces))
			for index, workspace := range gotWorkspaces {
				gotToSlash[index] = filepath.ToSlash(workspace)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("GetWorkspaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			sort.Strings(gotToSlash)
			if !reflect.DeepEqual(gotToSlash, tt.want) {
				t.Errorf("GetWorkspaces() = %v, want %v", gotToSlash, tt.want)
			}
		})
	}
}

func Test_GetWorkspaceIgnores(t *testing.T) {
	type test struct {
		name     string
		pm       PackageManager
		rootPath string
		want     []string
		wantErr  bool
	}

	cwd, err := os.Getwd()
	assert.NilError(t, err, "os.Getwd")
	want := map[string][]string{
		"npm":   {"**/node_modules/**"},
		"berry": {"**/node_modules", "**/.git", "**/.yarn"},
		"yarn":  {"apps/*/node_modules/**", "packages/*/node_modules/**"},
		"pnpm":  {"**/node_modules/**", "**/bower_components/**", "packages/skip"},
		"pnpm6": {"**/node_modules/**", "**/bower_components/**", "packages/skip"},
		"bun":   {"**/node_modules/**"},
		"deno":  {"**/node_modules/**"},
	}

	tests := make([]test, len(packageManagers))
	for i, packageManager := range packageManagers {
		tests[i] = test{
			name:     packageManager.Name,
			pm:       packageManager,
			rootPath: filepath.Join(cwd, "testdata"),
			want:     want[packageManager.Name],
			wantErr:  false,
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWorkspaceIgnores, err := tt.pm.GetWorkspaceIgnores(tt.rootPath)

			gotToSlash := make([]string, len(gotWorkspaceIgnores))
			for index, ignore := range gotWorkspaceIgnores {
				gotToSlash[index] = filepath.ToSlash(ignore)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("GetWorkspaceIgnores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotToSlash, tt.want) {
				t.Errorf("GetWorkspaceIgnores() = %v, want %v", gotToSlash, tt.want)
			}
		})
	}
}

func Test_CanPrune(t *testing.T) {
	type test struct {
		name     string
		pm       PackageManager
		rootPath string
		want     bool
		wantErr  bool
	}

	type want struct {
		want    bool
		wantErr bool
	}

	cwd, err := os.Getwd()
	assert.NilError(t, err, "os.Getwd")
	wants := map[string]want{
		"npm":   {true, false},
		"berry": {false, true},
		"yarn":  {true, false},
		"pnpm":  {true, false},
		"pnpm6": {true, false},
		"bun":   {true, false},
		"deno":  {true, false},
	}

	tests := make([]test, len(packageManagers))
	for i, packageManager := range packageManagers {
		tests[i] = test{
			name:     packageManager.Name,
			pm:       packageManager,
			rootPath: filepath.Join(cwd, "testdata/with-yarn"),
			want:     wants[packageManager.Name].want,
			wantErr:  wants[packageManager.Name].wantErr,
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			canPrune, err := tt.pm.CanPrune(tt.rootPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("CanPrune() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if canPrune != tt.want {
				t.Errorf("CanPrune() = %v, want %v", canPrune, tt.want)
			}
		})
	}
}

func TestPackageManager_InstallationDetection(t *testing.T) {
	for _, pm := range packageManagers {
		t.Run(pm.Name, func(t *testing.T) {
			// Test IsInstalled method
			installed := pm.IsInstalled()

			// Test GetInstallationPath method
			path, err := pm.GetInstallationPath()

			if installed {
				// If installed, path should be found without error
				assert.NilError(t, err, "GetInstallationPath should not error when package manager is installed")
				assert.Assert(t, path != "", "Installation path should not be empty when package manager is installed")
			} else {
				// If not installed, path should error or be empty
				assert.Assert(t, err != nil || path == "", "GetInstallationPath should error or return empty path when package manager is not installed")
			}
		})
	}
}

func TestGetAllPackageManagers(t *testing.T) {
	allPMs := GetAllPackageManagers()

	// Should return all available package managers
	assert.Assert(t, len(allPMs) > 0, "Should return at least one package manager")

	// Should include npm
	found := false
	for _, pm := range allPMs {
		if pm.Slug == "npm" {
			found = true
			break
		}
	}
	assert.Assert(t, found, "Should include npm package manager")
}

func TestPackageManager_ValidateProject(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NilError(t, err, "os.Getwd")

	// Test with a valid project (testdata/with-yarn)
	pm := yarn
	issues := pm.ValidateProject(filepath.Join(cwd, "testdata/with-yarn"))

	// The issues depend on whether yarn is installed and whether node_modules exists
	// We just check that the method doesn't panic and returns a slice
	assert.Assert(t, issues != nil, "ValidateProject should return a non-nil slice")
}

func TestPackageManager_GetEffectiveSpecFile(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NilError(t, err, "os.Getwd")

	tests := []struct {
		name string
		pm   PackageManager
		dir  string
		want string
	}{
		{
			name: "npm uses package.json",
			pm:   npm,
			dir:  cwd,
			want: "package.json",
		},
		{
			name: "deno prefers deno.json when available",
			pm:   deno,
			dir:  filepath.Join(cwd, "testdata/deno-with-json"),
			want: "deno.json",
		},
		{
			name: "deno falls back to package.json",
			pm:   deno,
			dir:  filepath.Join(cwd, "testdata/deno-package-only"),
			want: "package.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pm.GetEffectiveSpecFile(tt.dir)
			assert.Equal(t, got, tt.want)
		})
	}
}
