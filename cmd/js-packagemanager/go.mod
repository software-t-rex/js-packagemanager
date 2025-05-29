module github.com/software-t-rex/js-packagemanager/cmd/js-packagemanager

go 1.23

toolchain go1.23.9

replace github.com/software-t-rex/js-packagemanager => ../..

replace github.com/software-t-rex/packageJson => ../../../packageJson

require (
	github.com/software-t-rex/js-packagemanager v0.0.5
	github.com/software-t-rex/packageJson v0.0.3
)

require (
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/bmatcuk/doublestar/v4 v4.6.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
