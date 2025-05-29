# js packagemanager

This is a go module to deal with javascript ecosystem packagemanager.

# features
- detection of common package managers such as yarn, npm, pnpm, bun, and deno.
- get workspaces package.json when dealing with mono[repo|space]
- load a package.json into a struct (provided by the packageJson module in case you only need this)
- support for Deno's deno.json configuration file (with fallback to package.json)
- cross-platform support for Windows, macOS, and Linux
- CLI tool for package manager detection and analysis

## CLI Tool

This package includes a command-line tool called `jspm` for package manager detection and analysis.

### Installation

```bash
# Build from source
go run build.go

# Build for development (with race detection)
go run build.go -target=dev

# Install globally
go run build.go -target=install

# Cross-compile for multiple platforms
go run build.go -target=cross-compile

# Run tests
go run build.go -target=test

# Clean build artifacts
go run build.go -target=clean

# Run demo
go run build.go -target=demo

# Build everything
go run build.go -target=all

# Verbose output
go run build.go -target=build -v
```

### Available Build Targets

- **`build`** - Build the CLI tool for current platform
- **`dev`** - Build with race detection for development
- **`test`** - Run all tests
- **`cross-compile`** - Build for multiple platforms (Windows, macOS, Linux)
- **`clean`** - Remove build artifacts from `dist/` directory
- **`demo`** - Build and run interactive CLI demonstration
- **`install`** - Install globally using `go install`
- **`all`** - Complete build pipeline (clean → test → build → cross-compile)

### Demo

The demo target builds the CLI and runs several commands to showcase functionality:

```bash
go run build.go -target=demo

### Usage

```bash
# Detect package manager in current directory
jspm detect

# Detect in specific directory
jspm -dir /path/to/project detect

# List workspaces
jspm workspaces

# Check which package managers are installed
jspm installed

# Show detailed information (for installed command)
jspm -detailed installed
```

### Commands

- `detect` - Detect the package manager in use and output packageManager string
- `workspaces` - List workspace packages in a monorepo
- `installed` - Check which package managers are available on the system

### Output Formats

The CLI outputs human-readable text by default, which is clean and focused for command-line usage.

## Cross-Platform Support

The library includes comprehensive cross-platform support:

- Automatic executable name resolution (`.exe` on Windows)
- Cross-platform path handling
- Proper environment setup for command execution
- Support for Windows, macOS, and Linux

Other features are planed like some common commands to launch on the sytem with thoose package managers,
a better documentation is also planed. All of this depending on the interest manifested by the module.

## history of this package
The code in this repository was originally inspired by https://github.com/replit/upm and largely modified by the turbo team at vercel
This derivative work is not affiliated to any of thoose, but it is important to know the origin of the code in this repository.

This package was then extracted from turbo and cleaned up to be usable in other contexts than turbo repo.
This we lost some capabilities proposed by the original code, like pruning lock files, that will perhaps be re-integrated in the future, but there's no particular plan on this. I hope this will be helpful for a bunch of crazy devs around and as always you can propose PR to make this a better tool.

## Fundings
If you want, you can sponsors my work on this project here: https://github.com/sponsors/malko