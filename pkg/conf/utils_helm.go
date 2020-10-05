package conf

import (
	"bytes"
	"io/ioutil"
	"strings"

	"arhat.dev/pkg/exechelper"
	"github.com/rogpeppe/go-internal/semver"
)

/*
Copyright The Helm Authors.
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

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

func isHelmV2() bool {
	buf := new(bytes.Buffer)
	_, err := exechelper.Do(exechelper.Spec{
		Command: []string{"helm", "version", "--client", "--short"},
		Stdout:  buf,
		Stderr:  ioutil.Discard,
	})

	if err != nil {
		return false
	}

	// default to helm3
	ver := strings.TrimSpace(strings.TrimPrefix(buf.String(), "Client:"))
	return semver.Compare(ver, "v3") < 0
}
