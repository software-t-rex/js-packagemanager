package packagemanager

import (
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/v3/assert"
)

func TestDenoSpecFile(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NilError(t, err, "os.Getwd")

	tests := []struct {
		name      string
		directory string
		want      string
	}{
		{
			name:      "prefers deno.json when available",
			directory: filepath.Join(cwd, "testdata/deno-with-json"),
			want:      "deno.json",
		},
		{
			name:      "falls back to package.json when deno.json not available",
			directory: filepath.Join(cwd, "testdata/deno-package-only"),
			want:      "package.json",
		},
	}

	pm := deno
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pm.GetEffectiveSpecFile(tt.directory)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestDenoWorkspaceGlobs(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NilError(t, err, "os.Getwd")

	tests := []struct {
		name     string
		rootPath string
		want     []string
		wantErr  bool
	}{
		{
			name:     "reads workspaces from deno.json",
			rootPath: filepath.Join(cwd, "testdata/deno-with-json"),
			want:     []string{"apps/*", "packages/*"},
			wantErr:  false,
		},
		{
			name:     "reads workspaces from package.json when deno.json missing",
			rootPath: filepath.Join(cwd, "testdata/deno-package-only"),
			want:     []string{"apps/*", "packages/*"},
			wantErr:  false,
		},
		{
			name:     "errors when no workspaces found",
			rootPath: filepath.Join(cwd, "testdata/no-workspaces"),
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDenoWorkspaceGlobs(tt.rootPath)
			if tt.wantErr {
				assert.Assert(t, err != nil)
			} else {
				assert.NilError(t, err)
				assert.DeepEqual(t, got, tt.want)
			}
		})
	}
}

func TestDenoDetection(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NilError(t, err, "os.Getwd")

	tests := []struct {
		name      string
		directory string
		want      bool
		wantErr   bool
	}{
		{
			name:      "detects deno project with deno.json and deno.lock",
			directory: filepath.Join(cwd, "testdata/deno-with-json"),
			want:      true,
			wantErr:   false,
		},
		{
			name:      "detects deno project with package.json and deno.lock",
			directory: filepath.Join(cwd, "testdata/deno-package-only"),
			want:      true,
			wantErr:   false,
		},
		{
			name:      "does not detect when no lock file present",
			directory: filepath.Join(cwd, "testdata/no-lock"),
			want:      false,
			wantErr:   false,
		},
	}

	pm := deno
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip test if test directory doesn't exist
			if !PathExists(tt.directory) {
				t.Skipf("Test directory %s does not exist", tt.directory)
				return
			}

			got, err := pm.detect(tt.directory, &pm)
			if tt.wantErr {
				assert.Assert(t, err != nil)
			} else {
				assert.NilError(t, err)
				assert.Equal(t, got, tt.want)
			}
		})
	}
}
