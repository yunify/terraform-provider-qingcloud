// +-------------------------------------------------------------------------
// | Copyright (C) 2016 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

package test

import (
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/DATA-DOG/godog"

	"github.com/yunify/qingcloud-sdk-go/config"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func TestMain(m *testing.M) {
	setUp()

	context := func(s *godog.Suite) {
		QingCloudServiceFeatureContext(s)
	}
	options := godog.Options{
		Format: "pretty",
		Paths:  []string{"./features"},
		Tags:   "",
	}
	status := godog.RunWithOptions("*", context, options)
	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func setUp() {
	loadTestConfig()
	loadConfig()
	initQingStorService()
}

var err error
var tc *testConfig
var c *config.Config
var qcService *qc.QingCloudService

type testConfig struct {
	Zone string `json:"zone" yaml:"zone"`

	RetryWaitTime int `json:"retry_wait_time" yaml:"retry_wait_time"`
	MaxRetries    int `json:"max_retries" yaml:"max_retries"`
}

func loadTestConfig() {
	if tc == nil {
		configYAML, err := ioutil.ReadFile("./test_config.yaml")
		checkErrorForExit(err)

		tc = &testConfig{}
		err = yaml.Unmarshal(configYAML, tc)
		checkErrorForExit(err)
	}
}

func loadConfig() {
	if c == nil {
		c, err = config.NewDefault()
		checkErrorForExit(err)

		err = c.LoadConfigFromFilepath("./config.yaml")
		checkErrorForExit(err)
	}
}

func initQingStorService() {
	if qcService == nil {
		qcService, err = qc.Init(c)
		checkErrorForExit(err)
	}
}
