/*
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
*/

package main_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"arhat.dev/helm-stack/pkg/cmd"
)

func TestSuite(t *testing.T) {
	if !assert.NoError(t, os.Chdir(os.Getenv("PROJECT_DIR")), "failed to go to project dir") {
		return
	}

	if !assert.NoError(t, os.RemoveAll("./build/envs"), "failed to cleanup cluster data") {
		return
	}

	args := [][]string{
		{"helm-stack", "ensure", "--force-pull"},
		{"helm-stack", "gen", "all"},
		{"helm-stack", "clean", "-y"},
	}
	for _, a := range args {
		os.Args = a
		rootCmd := cmd.NewHelmStackCmd()
		err := rootCmd.Execute()
		if !assert.NoError(t, err, strings.Join(a, " ")) {
			return
		}
	}

	actual, err := ioutil.ReadFile("./build/envs/bar/manifests/testing.foo[foo@latest].yaml")
	assert.NoError(t, err)

	expected, err := ioutil.ReadFile("./testdata/expected-manifests.yaml")
	assert.NoError(t, err)

	assert.EqualValues(t, expected, actual)
}
