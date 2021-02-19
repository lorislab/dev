module github.com/lorislab/dev

go 1.15

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/gosuri/uitable v0.0.4
	github.com/mattn/go-runewidth v0.0.10 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.20.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	go.hein.dev/go-version v0.1.0
	gopkg.in/yaml.v2 v2.4.0
	helm.sh/helm/v3 v3.5.2
	k8s.io/cli-runtime v0.20.2 // indirect
	k8s.io/client-go v0.20.2
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
)
