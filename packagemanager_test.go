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
			want:             "nodejs-npm",
			wantErr:          false,
		},
		{
			name:             "finds pnpm6 from a package manager string",
			projectDirectory: cwd,
			pkg:              &packageJson.PackageJSON{PackageManager: "pnpm@1.2.3"},
			want:             "nodejs-pnpm6",
			wantErr:          false,
		},
		{
			name:             "finds pnpm from a package manager string",
			projectDirectory: cwd,
			pkg:              &packageJson.PackageJSON{PackageManager: "pnpm@7.8.9"},
			want:             "nodejs-pnpm",
			wantErr:          false,
		},
		{
			name:             "finds yarn from a package manager string",
			projectDirectory: cwd,
			pkg:              &packageJson.PackageJSON{PackageManager: "yarn@1.2.3"},
			want:             "nodejs-yarn",
			wantErr:          false,
		},
		{
			name:             "finds berry from a package manager string",
			projectDirectory: cwd,
			pkg:              &packageJson.PackageJSON{PackageManager: "yarn@2.3.4"},
			want:             "nodejs-berry",
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
			want:       "nodejs-npm",
			wantErr:    false,
		},
		{
			name:       "finds pnpm6 from a package manager string",
			pkgMngrStr: "pnpm@1.2.3",
			want:       "nodejs-pnpm6",
			wantErr:    false,
		},
		{
			name:       "finds pnpm from a package manager string",
			pkgMngrStr: "pnpm@7.8.9",
			want:       "nodejs-pnpm",
			wantErr:    false,
		},
		{
			name:       "finds yarn from a package manager string",
			pkgMngrStr: "yarn@1.2.3",
			want:       "nodejs-yarn",
			wantErr:    false,
		},
		{
			name:       "finds berry from a package manager string",
			pkgMngrStr: "yarn@2.3.4",
			want:       "nodejs-berry",
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
		"nodejs-npm":   filepath.Join(cwd, "testdata/with-yarn"),
		"nodejs-berry": filepath.Join(cwd, "testdata/with-yarn"),
		"nodejs-yarn":  filepath.Join(cwd, "testdata/with-yarn"),
		"nodejs-pnpm":  filepath.Join(cwd, "testdata/basic"),
		"nodejs-pnpm6": filepath.Join(cwd, "testdata/basic"),
	}

	want := map[string][]string{
		"nodejs-npm": {
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/apps/docs/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/apps/web/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/eslint-config-custom/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/tsconfig/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/ui/package.json")),
		},
		"nodejs-berry": {
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/apps/docs/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/apps/web/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/eslint-config-custom/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/tsconfig/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/ui/package.json")),
		},
		"nodejs-yarn": {
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/apps/docs/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/apps/web/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/eslint-config-custom/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/tsconfig/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/with-yarn/packages/ui/package.json")),
		},
		"nodejs-pnpm": {
			filepath.ToSlash(filepath.Join(cwd, "testdata/basic/apps/docs/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/basic/apps/web/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/basic/packages/eslint-config-custom/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/basic/packages/tsconfig/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/basic/packages/ui/package.json")),
		},
		"nodejs-pnpm6": {
			filepath.ToSlash(filepath.Join(cwd, "testdata/basic/apps/docs/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/basic/apps/web/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/basic/packages/eslint-config-custom/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/basic/packages/tsconfig/package.json")),
			filepath.ToSlash(filepath.Join(cwd, "testdata/basic/packages/ui/package.json")),
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
		"nodejs-npm":   {"**/node_modules/**"},
		"nodejs-berry": {"**/node_modules", "**/.git", "**/.yarn"},
		"nodejs-yarn":  {"apps/*/node_modules/**", "packages/*/node_modules/**"},
		"nodejs-pnpm":  {"**/node_modules/**", "**/bower_components/**", "packages/skip"},
		"nodejs-pnpm6": {"**/node_modules/**", "**/bower_components/**", "packages/skip"},
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
		"nodejs-npm":   {true, false},
		"nodejs-berry": {false, true},
		"nodejs-yarn":  {true, false},
		"nodejs-pnpm":  {true, false},
		"nodejs-pnpm6": {true, false},
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
