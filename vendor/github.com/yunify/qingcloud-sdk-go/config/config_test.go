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

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yunify/qingcloud-sdk-go/logger"
)

func TestConfig(t *testing.T) {
	config := Config{
		AccessKeyID:       "AccessKeyID",
		SecretAccessKey:   "SecretAccessKey",
		Protocol:          "http",
		ConnectionRetries: 10,
		LogLevel:          "warn",
	}

	assert.Equal(t, "AccessKeyID", config.AccessKeyID)
	assert.Equal(t, "SecretAccessKey", config.SecretAccessKey)
	assert.Equal(t, "http", config.Protocol)
	assert.Equal(t, 10, config.ConnectionRetries)
	assert.Equal(t, "warn", config.LogLevel)
}

func TestConfig_LoadDefaultConfig(t *testing.T) {
	config := Config{}
	config.LoadDefaultConfig()

	assert.Equal(t, "", config.AccessKeyID)
	assert.Equal(t, "", config.SecretAccessKey)
	assert.Equal(t, "https", config.Protocol)
	assert.Equal(t, "warning", logger.GetLevel())
}

func TestConfig_LoadUserConfig(t *testing.T) {
	config := Config{}
	config.LoadUserConfig()

	assert.Equal(t, "https", config.Protocol)
}

func TestConfig_LoadConfigFromContent(t *testing.T) {
	fileContent := `
qy_access_key_id: 'access_key_id'
qy_secret_access_key: 'secret_access_key'

log_level: 'debug'
`

	config := Config{}
	config.LoadConfigFromContent([]byte(fileContent))

	assert.Equal(t, "access_key_id", config.AccessKeyID)
	assert.Equal(t, "secret_access_key", config.SecretAccessKey)
	assert.Equal(t, "https", config.Protocol)
	assert.Equal(t, "debug", logger.GetLevel())
}

func TestNewDefault(t *testing.T) {
	config, err := NewDefault()
	assert.Nil(t, err)

	assert.Equal(t, "", config.AccessKeyID)
	assert.Equal(t, "", config.SecretAccessKey)
	assert.Equal(t, "https", config.Protocol)
	assert.Equal(t, 3, config.ConnectionRetries)
}

func TestNew(t *testing.T) {
	config, err := New("AccessKeyID", "SecretAccessKey")
	assert.Nil(t, err)

	assert.Equal(t, "AccessKeyID", config.AccessKeyID)
	assert.Equal(t, "SecretAccessKey", config.SecretAccessKey)
	assert.Equal(t, "https", config.Protocol)
}

func TestNewWithEndpoint(t *testing.T) {
	config, err := NewWithEndpoint("AccessKeyID", "SecretAccessKey", "https://api.qingcloud.com/iaas")
	assert.Nil(t, err)

	assert.Equal(t, "AccessKeyID", config.AccessKeyID)
	assert.Equal(t, "SecretAccessKey", config.SecretAccessKey)
	assert.Equal(t, "https", config.Protocol)
	assert.Equal(t, "api.qingcloud.com", config.Host)
	assert.Equal(t, 443, config.Port)
	assert.Equal(t, "/iaas", config.URI)
}
