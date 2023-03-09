package packagemanager

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"gotest.tools/v3/assert"
)

func TestInferRoot(t *testing.T) {
	type file struct {
		path    string
		content []byte
	}

	tests := []struct {
		name               string
		fs                 []file
		executionDirectory string
		rootPath           string
		packageMode        PackageType
	}{
		// Scenario 0
		{
			name: "turbo.json at current dir, no package.json",
			fs: []file{
				{path: "turbo.json"},
			},
			executionDirectory: "",
			rootPath:           "",
			packageMode:        Multi,
		},
		{
			name: "turbo.json at parent dir, no package.json",
			fs: []file{
				{path: "execution/path/subdir/.file"},
				{path: "turbo.json"},
			},
			executionDirectory: "execution/path/subdir",
			// This is "no inference"
			rootPath:    "execution/path/subdir",
			packageMode: Multi,
		},
		// Scenario 1A
		{
			name: "turbo.json at current dir, has package.json, has workspaces key",
			fs: []file{
				{path: "turbo.json"},
				{
					path:    "package.json",
					content: []byte("{ \"workspaces\": [ \"exists\" ] }"),
				},
			},
			executionDirectory: "",
			rootPath:           "",
			packageMode:        Multi,
		},
		{
			name: "turbo.json at parent dir, has package.json, has workspaces key",
			fs: []file{
				{path: "execution/path/subdir/.file"},
				{path: "turbo.json"},
				{
					path:    "package.json",
					content: []byte("{ \"workspaces\": [ \"exists\" ] }"),
				},
			},
			executionDirectory: "execution/path/subdir",
			rootPath:           "",
			packageMode:        Multi,
		},
		{
			name: "turbo.json at parent dir, has package.json, has pnpm workspaces",
			fs: []file{
				{path: "execution/path/subdir/.file"},
				{path: "turbo.json"},
				{
					path:    "package.json",
					content: []byte("{}"),
				},
				{
					path:    "pnpm-workspace.yaml",
					content: []byte("packages:\n  - docs"),
				},
			},
			executionDirectory: "execution/path/subdir",
			rootPath:           "",
			packageMode:        Multi,
		},
		// Scenario 1A aware of the weird thing we do for packages.
		{
			name: "turbo.json at current dir, has package.json, has packages key",
			fs: []file{
				{path: "turbo.json"},
				{
					path:    "package.json",
					content: []byte("{ \"packages\": [ \"exists\" ] }"),
				},
			},
			executionDirectory: "",
			rootPath:           "",
			packageMode:        Single,
		},
		{
			name: "turbo.json at parent dir, has package.json, has packages key",
			fs: []file{
				{path: "execution/path/subdir/.file"},
				{path: "turbo.json"},
				{
					path:    "package.json",
					content: []byte("{ \"packages\": [ \"exists\" ] }"),
				},
			},
			executionDirectory: "execution/path/subdir",
			rootPath:           "",
			packageMode:        Single,
		},
		// Scenario 1A aware of the the weird thing we do for packages when both methods of specification exist.
		{
			name: "turbo.json at current dir, has package.json, has workspace and packages key",
			fs: []file{
				{path: "turbo.json"},
				{
					path:    "package.json",
					content: []byte("{ \"workspaces\": [ \"clobbered\" ], \"packages\": [ \"exists\" ] }"),
				},
			},
			executionDirectory: "",
			rootPath:           "",
			packageMode:        Multi,
		},
		{
			name: "turbo.json at parent dir, has package.json, has workspace and packages key",
			fs: []file{
				{path: "execution/path/subdir/.file"},
				{path: "turbo.json"},
				{
					path:    "package.json",
					content: []byte("{ \"workspaces\": [ \"clobbered\" ], \"packages\": [ \"exists\" ] }"),
				},
			},
			executionDirectory: "execution/path/subdir",
			rootPath:           "",
			packageMode:        Multi,
		},
		// Scenario 1B
		{
			name: "turbo.json at current dir, has package.json, no workspaces",
			fs: []file{
				{path: "turbo.json"},
				{
					path:    "package.json",
					content: []byte("{}"),
				},
			},
			executionDirectory: "",
			rootPath:           "",
			packageMode:        Single,
		},
		{
			name: "turbo.json at parent dir, has package.json, no workspaces",
			fs: []file{
				{path: "execution/path/subdir/.file"},
				{path: "turbo.json"},
				{
					path:    "package.json",
					content: []byte("{}"),
				},
			},
			executionDirectory: "execution/path/subdir",
			rootPath:           "",
			packageMode:        Single,
		},
		{
			name: "turbo.json at parent dir, has package.json, no workspaces, includes pnpm",
			fs: []file{
				{path: "execution/path/subdir/.file"},
				{path: "turbo.json"},
				{
					path:    "package.json",
					content: []byte("{}"),
				},
				{
					path:    "pnpm-workspace.yaml",
					content: []byte(""),
				},
			},
			executionDirectory: "execution/path/subdir",
			rootPath:           "",
			packageMode:        Single,
		},
		// Scenario 2A
		{
			name:               "no turbo.json, no package.json at current",
			fs:                 []file{},
			executionDirectory: "",
			rootPath:           "",
			packageMode:        Multi,
		},
		{
			name: "no turbo.json, no package.json at parent",
			fs: []file{
				{path: "execution/path/subdir/.file"},
			},
			executionDirectory: "execution/path/subdir",
			rootPath:           "execution/path/subdir",
			packageMode:        Multi,
		},
		// Scenario 2B
		{
			name: "no turbo.json, has package.json with workspaces at current",
			fs: []file{
				{
					path:    "package.json",
					content: []byte("{ \"workspaces\": [ \"exists\" ] }"),
				},
			},
			executionDirectory: "",
			rootPath:           "",
			packageMode:        Multi,
		},
		{
			name: "no turbo.json, has package.json with workspaces at parent",
			fs: []file{
				{path: "execution/path/subdir/.file"},
				{
					path:    "package.json",
					content: []byte("{ \"workspaces\": [ \"exists\" ] }"),
				},
			},
			executionDirectory: "execution/path/subdir",
			rootPath:           "execution/path/subdir",
			packageMode:        Multi,
		},
		{
			name: "no turbo.json, has package.json with pnpm workspaces at parent",
			fs: []file{
				{path: "execution/path/subdir/.file"},
				{
					path:    "package.json",
					content: []byte("{ \"workspaces\": [ \"exists\" ] }"),
				},
				{
					path:    "pnpm-workspace.yaml",
					content: []byte("packages:\n  - docs"),
				},
			},
			executionDirectory: "execution/path/subdir",
			rootPath:           "execution/path/subdir",
			packageMode:        Multi,
		},
		// Scenario 3A
		{
			name: "no turbo.json, lots of package.json files but no workspaces",
			fs: []file{
				{
					path:    "package.json",
					content: []byte("{}"),
				},
				{
					path:    "one/package.json",
					content: []byte("{}"),
				},
				{
					path:    "one/two/package.json",
					content: []byte("{}"),
				},
				{
					path:    "one/two/three/package.json",
					content: []byte("{}"),
				},
			},
			executionDirectory: "one/two/three",
			rootPath:           "one/two/three",
			packageMode:        Single,
		},
		// Scenario 3BI
		{
			name: "no turbo.json, lots of package.json files, and a workspace at the root that matches execution directory",
			fs: []file{
				{
					path:    "package.json",
					content: []byte("{ \"workspaces\": [ \"one/two/three\" ] }"),
				},
				{
					path:    "one/package.json",
					content: []byte("{}"),
				},
				{
					path:    "one/two/package.json",
					content: []byte("{}"),
				},
				{
					path:    "one/two/three/package.json",
					content: []byte("{}"),
				},
			},
			executionDirectory: "one/two/three",
			rootPath:           "one/two/three",
			packageMode:        Multi,
		},
		// Scenario 3BII
		{
			name: "no turbo.json, lots of package.json files, and a workspace at the root that matches execution directory",
			fs: []file{
				{
					path:    "package.json",
					content: []byte("{ \"workspaces\": [ \"does-not-exist\" ] }"),
				},
				{
					path:    "one/package.json",
					content: []byte("{}"),
				},
				{
					path:    "one/two/package.json",
					content: []byte("{}"),
				},
				{
					path:    "one/two/three/package.json",
					content: []byte("{}"),
				},
			},
			executionDirectory: "one/two/three",
			rootPath:           "one/two/three",
			packageMode:        Single,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsRoot := t.TempDir()
			for _, file := range tt.fs {
				path := filepath.Join(fsRoot, file.path)
				assert.NilError(t, os.MkdirAll(filepath.Dir(path), 0777))
				assert.NilError(t, os.WriteFile(path, file.content, 0777))
			}

			turboRoot, packageMode := InferRoot(filepath.Join(fsRoot, tt.executionDirectory))
			if !reflect.DeepEqual(turboRoot, filepath.Join(fsRoot, tt.rootPath)) {
				t.Errorf("InferRoot() turboRoot = %v, want %v", turboRoot, filepath.Join(fsRoot, tt.rootPath))
			}
			if packageMode != tt.packageMode {
				t.Errorf("InferRoot() packageMode = %v, want %v", packageMode, tt.packageMode)
			}
		})
	}
}
