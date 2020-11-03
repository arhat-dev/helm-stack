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

package conf

import (
	"arhat.dev/pkg/log"
)

type Config struct {
	App AppConfig `json:"app" yaml:"app"`

	Repos        []RepoSpec    `json:"repos" yaml:"repos"`
	Charts       []ChartSpec   `json:"charts" yaml:"charts"`
	Environments []Environment `json:"environments" yaml:"environments"`
}

type AppConfig struct {
	Log log.ConfigSet `json:"log" yaml:"log"`

	DebugHelm bool `json:"debugHelm" yaml:"debugHelm"`

	// ChartsDir for chart cache
	ChartsDir string `json:"chartsDir" yaml:"chartsDir"`

	// EnvironmentsDir for environment values
	EnvironmentsDir string `json:"environmentsDir" yaml:"environmentsDir"`

	// LocalChartsDir for charts stored locally
	LocalChartsDir string `json:"localChartsDir" yaml:"localChartsDir"`
}

func (c *AppConfig) Override(o *AppConfig) *AppConfig {
	if o == nil {
		return c
	}

	result := &AppConfig{
		DebugHelm:       c.DebugHelm,
		ChartsDir:       c.ChartsDir,
		EnvironmentsDir: c.EnvironmentsDir,
		LocalChartsDir:  c.LocalChartsDir,
	}

	if o.DebugHelm {
		result.DebugHelm = o.DebugHelm
	}

	if o.ChartsDir != "" {
		result.ChartsDir = o.ChartsDir
	}

	if o.EnvironmentsDir != "" {
		result.EnvironmentsDir = o.EnvironmentsDir
	}

	if o.LocalChartsDir != "" {
		result.LocalChartsDir = o.LocalChartsDir
	}

	return result
}
