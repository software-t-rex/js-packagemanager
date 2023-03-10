module github.com/software-t-rex/js-packagemanager

go 1.20

replace github.com/software-t-rex/packageJson => ../packageJson

require (
	github.com/Masterminds/semver v1.5.0
	github.com/bmatcuk/doublestar/v4 v4.6.0
	github.com/software-t-rex/packageJson v0.0.3
	gotest.tools/v3 v3.4.0
	sigs.k8s.io/yaml v1.3.0
)

require (
	github.com/google/go-cmp v0.5.5 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
