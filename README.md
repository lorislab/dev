# dev

Dev is a tool for local k8s development 

[![License](https://img.shields.io/github/license/lorislab/dev?style=for-the-badge&logo=apache)](https://www.apache.org/licenses/LICENSE-2.0)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/lorislab/dev?sort=semver&logo=github&style=for-the-badge)](https://github.com/lorislab/dev/releases/latest)


## Commands

```shell
dev env help
```

Local environment configuration `env.yaml`
```
cluster:
  namespace: test
apps:
  ping-quarkus:
    tags:
      - tests
    helm:
      chart: helmrepo/ping-quarkus
      version: ">=0.0.0-0"
      values:
        test: value
        app:
          env:
            TEST: example-variable
      files:
        - values/ping-quarkus.yaml
```

The main commands:
* `dev env status` - status of the applications in the local environment
* `dev env sync` - synchronize applications in the local environment
* `dev env uninstall` - uninstall applications from the local environment


```
❯ dev env status
PRIO    CHART                           NAME            NAMESPACE       RULE            DEPLOY  VERSION                         ACTION
0       onecxportal/ping-quarkus        ping-quarkus    test            >=0.0.0-0       <nil>   0.18.0-rc009.g5f9df407609a      install
```