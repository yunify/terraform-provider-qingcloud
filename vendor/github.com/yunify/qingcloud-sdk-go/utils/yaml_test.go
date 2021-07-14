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

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYAMLDecode_Unknown(t *testing.T) {
	yamlString := `
key1: "This is a string." # Single Line Comment
key2: 10.50
key3:
  - null
  - nestedKey1: Anothor string
`

	anyData, err := YAMLDecode([]byte(yamlString))
	assert.Nil(t, err)
	data := anyData.(map[interface{}]interface{})
	assert.Equal(t, 10.50, data["key2"])
}

func TestYAMLDecode_Known(t *testing.T) {
	type SampleYAML struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
	}
	sampleYAMLString := `name: "NAME"`

	sample := SampleYAML{Name: "NaMe", Description: "DeScRiPtIoN"}
	anyDataPointer, err := YAMLDecode([]byte(sampleYAMLString), &sample)
	assert.Nil(t, err)
	data := anyDataPointer.(*SampleYAML)
	assert.Equal(t, "NAME", sample.Name)
	assert.Equal(t, "DeScRiPtIoN", sample.Description)
	assert.Equal(t, "NAME", (*data).Name)
	assert.Equal(t, "DeScRiPtIoN", (*data).Description)
}

func TestYAMLDecode_Empty(t *testing.T) {
	yamlString := ""

	anyData, err := YAMLDecode([]byte(yamlString))
	assert.Nil(t, err)
	assert.Nil(t, anyData)
}

func TestYAMLEncode(t *testing.T) {
	type SampleYAML struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
	}
	sample := SampleYAML{Name: "NaMe", Description: "DeScRiPtIoN"}

	yamlBytes, err := YAMLEncode(sample)
	assert.Nil(t, err)
	assert.Equal(t, "name: NaMe\ndescription: DeScRiPtIoN\n", string(yamlBytes))
}
