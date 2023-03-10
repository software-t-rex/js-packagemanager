package packagemanager

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/software-t-rex/packageJson"
	"gotest.tools/v3/assert"
)

func pnpmPatchesSection(t *testing.T, pkgJSON *packageJson.PackageJSON) map[string]interface{} {
	t.Helper()
	pnpmSection, ok := pkgJSON.Pnpm.(map[string]interface{})
	assert.Assert(t, ok)
	patchesSection, ok := pnpmSection["patchedDependencies"].(map[string]interface{})
	assert.Assert(t, ok)
	return patchesSection
}

func getPnpmPackageJSON(t *testing.T) *packageJson.PackageJSON {
	t.Helper()
	rawCwd, err := os.Getwd()
	assert.NilError(t, err)
	assert.Assert(t, filepath.IsAbs(rawCwd), rawCwd+" is not absolute path")
	cwd := rawCwd
	assert.NilError(t, err)
	pkgJSONPath := filepath.Join(cwd, "testdata", "pnpm-patches.json")
	pkgJSON, err := packageJson.Read(pkgJSONPath)
	assert.NilError(t, err)
	return pkgJSON
}

func Test_PnpmPrunePatches_KeepsNecessary(t *testing.T) {
	pkgJSON := getPnpmPackageJSON(t)
	initialPatches := pnpmPatchesSection(t, pkgJSON)

	assert.DeepEqual(t, initialPatches, map[string]interface{}{"is-odd@3.0.1": "patches/is-odd@3.0.1.patch"})

	err := pnpmPrunePatches(pkgJSON, []string{"patches/is-odd@3.0.1.patch"})
	assert.NilError(t, err)

	newPatches := pnpmPatchesSection(t, pkgJSON)
	assert.DeepEqual(t, newPatches, map[string]interface{}{"is-odd@3.0.1": "patches/is-odd@3.0.1.patch"})
}

func Test_PnpmPrunePatches_RemovesExtra(t *testing.T) {
	pkgJSON := getPnpmPackageJSON(t)
	initialPatches := pnpmPatchesSection(t, pkgJSON)

	assert.DeepEqual(t, initialPatches, map[string]interface{}{"is-odd@3.0.1": "patches/is-odd@3.0.1.patch"})

	err := pnpmPrunePatches(pkgJSON, nil)
	assert.NilError(t, err)

	newPatches := pnpmPatchesSection(t, pkgJSON)
	assert.DeepEqual(t, newPatches, map[string]interface{}{})
}
