# helm-stack

[![CI](https://github.com/arhat-dev/helm-stack/workflows/CI/badge.svg)](https://github.com/arhat-dev/helm-stack/actions?query=workflow%3ACI)
[![Build](https://github.com/arhat-dev/helm-stack/workflows/Build/badge.svg)](https://github.com/arhat-dev/helm-stack/actions?query=workflow%3ABuild)
[![PkgGoDev](https://pkg.go.dev/badge/arhat.dev/helm-stack)](https://pkg.go.dev/arhat.dev/helm-stack)
[![GoReportCard](https://goreportcard.com/badge/arhat.dev/helm-stack)](https://goreportcard.com/report/arhat.dev/helm-stack)
[![codecov](https://codecov.io/gh/arhat-dev/helm-stack/branch/master/graph/badge.svg)](https://codecov.io/gh/arhat-dev/helm-stack)

Stack your cluster deployments the easy way

## Build

```bash
make helm-stack
```

## Install

Install helm-stack to `${GOPATH}/bin/helm-stack`

```bash
GOOS=$(go env GOHOSTOS) GOARCH=$(go env GOHOSTARCH) go install ./cmd/helm-stack
```

## Config

All configration files provided to helm-stack will be merged, please make sure there are no duplicate items in your configuration files

- For file based config: Please refer to [`.helm-stack.yaml`](./.helm-stack.yaml) for example
- For directory based config: Please refer to [`.helm-stack`](./.helm-stack) for example

**NOTE:** helm-stack by default will try to read configuration files in `.helm-stack` and `helm-stack.yaml`, but if you have provided any `-c` or `--config` flag, helm-stack will not use these default config files.

## Workflow

TL;DR: you can run `make test.cmd` to take a walkthrough of the workflow

1. Define your charts and deployment environments in a yaml/json config file or using multiple yaml/json config files (in the same parent directory)
2. Run `helm-stack ensure` to ensure charts and values files
3. Update yaml values files in `<environments-dir>/<environment-name>` according to your deployments requirements
4. After several updates, there may be some charts unused, you can remove these charts and related values file with `helm-stack clean`
5. Run `helm-stack gen` to generate kubernetes manifests
6. Run `helm-stack apply` to deploy manifests to your environment

Please refer to [`.helm-stack`](./.helm-stack/) for config structure

## LICENSE

```text
Copyright 2020 The arhat.dev Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
