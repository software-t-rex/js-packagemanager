package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	packagemanager "github.com/software-t-rex/js-packagemanager"
	"github.com/software-t-rex/packageJson"
)

func main() {
	var (
		detailed = flag.Bool("detailed", false, "Show detailed information (for installed command)")
	)

	// Add shorthand flag
	flag.BoolVar(detailed, "D", false, "Show detailed information (for installed command) (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <command> [directory]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  detect [dir]         Detect package manager and get packageManager string\n")
		fmt.Fprintf(os.Stderr, "  workspaces [dir]     List workspace packages\n")
		fmt.Fprintf(os.Stderr, "  installed            Check if package managers are installed\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	command := flag.Arg(0)

	// Get directory from second argument or default to current directory
	projectDir := "."
	if flag.NArg() > 1 {
		projectDir = flag.Arg(1)
	}

	// Resolve absolute path for commands that need it
	var absDir string
	if command == "detect" || command == "workspaces" {
		var err error
		absDir, err = filepath.Abs(projectDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid directory: %v\n", err)
			os.Exit(1)
		}
	}

	switch command {
	case "detect":
		detectCommand(absDir)
	case "workspaces":
		workspacesCommand(absDir)
	case "installed":
		installedCommand(*detailed)
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func detectCommand(projectDir string) {
	// Try to read package.json
	pkgJSONPath := filepath.Join(projectDir, "package.json")
	var pkg *packageJson.PackageJSON

	if packagemanager.FileExists(pkgJSONPath) {
		var err error
		pkg, err = packageJson.Read(pkgJSONPath)
		if err != nil {
			// Silently ignore package.json read errors
		}
	} else {
		pkg = &packageJson.PackageJSON{}
	}

	pm, err := packagemanager.GetPackageManager(projectDir, pkg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	if !pm.IsInstalled() {
		fmt.Fprintf(os.Stderr, "Error: %s is not installed\n", pm.Name)
		os.Exit(1)
	}

	result, err := pm.DetectPackageManagerString(projectDir, pkg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(result.PackageManagerString)
}

func workspacesCommand(projectDir string) {
	// Try to read package.json
	pkgJSONPath := filepath.Join(projectDir, "package.json")
	var pkg *packageJson.PackageJSON

	if packagemanager.FileExists(pkgJSONPath) {
		var err error
		pkg, err = packageJson.Read(pkgJSONPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not read package.json: %v\n", err)
			os.Exit(1)
		}
	} else {
		pkg = &packageJson.PackageJSON{}
	}

	pm, err := packagemanager.GetPackageManager(projectDir, pkg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	workspaces, err := pm.GetWorkspaces(projectDir, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	ignores, err := pm.GetWorkspaceIgnores(projectDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	// Build package manager string with version using DetectPackageManagerString
	pmString := pm.Slug
	if result, err := pm.DetectPackageManagerString(projectDir, pkg); err == nil {
		pmString = result.PackageManagerString
	}

	// Output text format
	fmt.Printf("PackageManager: %s\n", pmString)
	if len(workspaces) > 0 {
		fmt.Printf("Workspaces (%d):\n", len(workspaces))
		for _, workspace := range workspaces {
			fmt.Printf("  - %s\n", workspace)
		}
	}
	if len(ignores) > 0 {
		fmt.Printf("Ignored:\n")
		for _, ignore := range ignores {
			fmt.Printf("  - %s\n", ignore)
		}
	}
}

func installedCommand(detailed bool) {
	for _, pm := range getAllPackageManagers() {
		if pm.IsInstalled() {
			// Get version to verify installation and check if it matches this specific package manager variant
			if version, err := pm.GetStandardVersion(); err == nil {
				// Use the Matches function to verify this version actually corresponds to this package manager variant
				if matches, matchErr := pm.Matches(pm.Slug, version); matchErr == nil && matches {
					if detailed {
						// Detailed output: add path in parentheses
						if path, err := pm.GetInstallationPath(); err == nil {
							fmt.Printf("%s@%s (%s)\n", pm.Slug, version, path)
						} else {
							fmt.Printf("%s@%s\n", pm.Slug, version)
						}
					} else {
						// Simple output: just slug@version
						fmt.Printf("%s@%s\n", pm.Slug, version)
					}
				}
			}
		}
	}
}

func getAllPackageManagers() []packagemanager.PackageManager {
	return packagemanager.GetAllPackageManagers()
}
