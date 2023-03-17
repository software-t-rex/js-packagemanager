module github.com/software-t-rex/js-packagemanager

go 1.20

replace github.com/software-t-rex/packageJson => ../packageJson

require (
	github.com/Masterminds/semver v1.5.0
	github.com/bmatcuk/doublestar/v4 v4.6.0
	github.com/software-t-rex/packageJson v0.0.3
	gopkg.in/yaml.v3 v3.0.1
	gotest.tools/v3 v3.4.0
)

require (
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
