//go:build ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	binaryName = "jspm"
	buildDir   = "dist"
	cmdPath    = "cmd/js-packagemanager"
)

func main() {
	var (
		target  = flag.String("target", "build", "Target to run: build, test, cross-compile, clean, demo, install, dev, all")
		verbose = flag.Bool("v", false, "Verbose output")
	)
	flag.Parse()

	// Change working directory to the CLI module directory
	originalDir, err := os.Getwd()
	if err != nil {
		fatal("Failed to get current directory: %v", err)
	}

	if err := os.Chdir(cmdPath); err != nil {
		fatal("Failed to change to CLI directory: %v", err)
	}

	if *verbose {
		fmt.Printf("Running target: %s from %s\n", *target, cmdPath)
	}

	// Store the absolute path to the original directory for building outputs
	buildDirAbs := filepath.Join(originalDir, buildDir)

	switch *target {
	case "build":
		build(false, *verbose, buildDirAbs)
	case "dev":
		build(true, *verbose, buildDirAbs)
	case "test":
		test(*verbose)
	case "cross-compile":
		crossCompile(*verbose, buildDirAbs)
	case "clean":
		clean(*verbose, originalDir)
	case "demo":
		demo(*verbose, buildDirAbs)
	case "install":
		install(*verbose)
	case "all":
		all(*verbose, buildDirAbs, originalDir)
	default:
		fmt.Printf("Unknown target: %s\n", *target)
		fmt.Println("Available targets: build, dev, test, cross-compile, clean, demo, install, all")
		os.Exit(1)
	}
}

func build(dev bool, verbose bool, buildDirAbs string) {
	if verbose {
		fmt.Println("Creating build directory...")
	}
	if err := os.MkdirAll(buildDirAbs, 0755); err != nil {
		fatal("Failed to create build directory: %v", err)
	}

	// First ensure dependencies are up to date in the main module (from CLI directory, go up two levels)
	if verbose {
		fmt.Println("Updating main module dependencies...")
	}
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = "../.."
	if err := tidyCmd.Run(); err != nil {
		if verbose {
			fmt.Printf("Warning: go mod tidy failed in main module: %v\n", err)
		}
	}

	// Then update CLI module dependencies (current directory)
	if verbose {
		fmt.Println("Updating CLI module dependencies...")
	}
	tidyCmd = exec.Command("go", "mod", "tidy")
	if err := tidyCmd.Run(); err != nil {
		if verbose {
			fmt.Printf("Warning: go mod tidy failed in CLI module: %v\n", err)
		}
	}

	args := []string{"build"}
	if dev {
		args = append(args, "-race")
		if verbose {
			fmt.Println("Building with race detection...")
		}
	}

	output := filepath.Join(buildDirAbs, binaryName)
	if dev {
		output += "-dev"
	}
	if runtime.GOOS == "windows" {
		output += ".exe"
	}

	args = append(args, "-o", output, "main.go")

	if verbose {
		fmt.Printf("Running: go %s (from CLI directory)\n", strings.Join(args, " "))
	}

	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fatal("Build failed: %v", err)
	}

	fmt.Printf("✅ Built %s\n", filepath.Join(buildDir, filepath.Base(output)))
}

func test(verbose bool) {
	args := []string{"test"}
	if verbose {
		args = append(args, "-v")
	}
	args = append(args, "./...")

	if verbose {
		fmt.Printf("Running: go %s\n", strings.Join(args, " "))
	}

	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fatal("Tests failed: %v", err)
	}

	fmt.Println("✅ Tests passed")
}

func crossCompile(verbose bool, buildDirAbs string) {
	if verbose {
		fmt.Println("Creating build directory...")
	}
	if err := os.MkdirAll(buildDirAbs, 0755); err != nil {
		fatal("Failed to create build directory: %v", err)
	}

	platforms := []struct {
		goos   string
		goarch string
		name   string
	}{
		{"windows", "amd64", "windows-amd64"},
		{"windows", "386", "windows-386"},
		{"linux", "amd64", "linux-amd64"},
		{"linux", "386", "linux-386"},
		{"linux", "arm64", "linux-arm64"},
		{"darwin", "amd64", "darwin-amd64"},
		{"darwin", "arm64", "darwin-arm64"},
	}

	fmt.Println("Cross-compiling for multiple platforms...")
	for _, platform := range platforms {
		output := filepath.Join(buildDirAbs, fmt.Sprintf("%s-%s", binaryName, platform.name))
		if platform.goos == "windows" {
			output += ".exe"
		}

		env := append(os.Environ(),
			fmt.Sprintf("GOOS=%s", platform.goos),
			fmt.Sprintf("GOARCH=%s", platform.goarch),
		)

		if verbose {
			fmt.Printf("Building for %s/%s...\n", platform.goos, platform.goarch)
		}

		cmd := exec.Command("go", "build", "-o", output, "main.go")
		cmd.Env = env
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fatal("Cross-compile failed for %s/%s: %v", platform.goos, platform.goarch, err)
		}

		fmt.Printf("✅ Built %s\n", filepath.Join(buildDir, filepath.Base(output)))
	}
}

func clean(verbose bool, originalDir string) {
	buildDirPath := filepath.Join(originalDir, buildDir)
	if verbose {
		fmt.Printf("Removing %s directory...\n", buildDirPath)
	}

	if err := os.RemoveAll(buildDirPath); err != nil {
		fatal("Failed to clean: %v", err)
	}

	fmt.Printf("✅ Cleaned %s directory\n", buildDir)
}

func demo(verbose bool, buildDirAbs string) {
	// First build the project
	build(false, verbose, buildDirAbs)

	binaryPath := filepath.Join(buildDirAbs, binaryName)
	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}

	fmt.Println("\n=== JSPM CLI Demo ===")

	demos := []struct {
		name string
		args []string
	}{
		{"Detect package manager", []string{"detect"}},
		{"List installed package managers", []string{"installed"}},
		{"Show workspaces (if any)", []string{"workspaces"}},
	}

	for _, demo := range demos {
		fmt.Printf("\n%s:\n", demo.name)
		if verbose {
			fmt.Printf("Running: %s %s\n", binaryPath, strings.Join(demo.args, " "))
		}

		cmd := exec.Command(binaryPath, demo.args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("⚠️  Demo command failed: %v\n", err)
		}
	}

	fmt.Println("\n✅ Demo complete!")
}

func install(verbose bool) {
	if verbose {
		fmt.Println("Installing globally...")
	}

	// Install needs to run from the current directory (CLI module)
	cmd := exec.Command("go", "install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fatal("Install failed: %v", err)
	}

	fmt.Println("✅ Installed globally")
}

func all(verbose bool, buildDirAbs string, originalDir string) {
	fmt.Println("Running full build pipeline...")
	clean(verbose, originalDir)
	test(verbose)
	build(false, verbose, buildDirAbs)
	crossCompile(verbose, buildDirAbs)
	fmt.Println("✅ All tasks completed!")
}

func fatal(format string, args ...interface{}) {
	fmt.Printf("❌ "+format+"\n", args...)
	os.Exit(1)
}
