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

package request

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yunify/qingcloud-sdk-go/config"
	"github.com/yunify/qingcloud-sdk-go/request/data"
)

type InstanceServiceProperties struct {
	Zone *string `json:"zone" name:"zone"` // Required
}

type DescribeInstancesInput struct {
	ImageID       []*string `json:"image_id" name:"image_id" location:"params"`
	InstanceClass *int      `json:"instance_class" name:"instance_class" location:"params" default:"0"` // Available values: 0, 1
	InstanceType  []*string `json:"instance_type" name:"instance_type" location:"params"`
	Instances     []*string `json:"instances" name:"instances" location:"params"`
	Limit         *int      `json:"limit" name:"limit" location:"params"`
	Offset        *int      `json:"offset" name:"offset" location:"params"`
	SearchWord    *string   `json:"search_word" name:"search_word" location:"params"`
	Status        []*string `json:"status" name:"status" location:"params"` // Available values: pending, running, stopped, suspended, terminated, ceased
	Tags          []*string `json:"tags" name:"tags" location:"params"`
	Verbose       *int      `json:"verbose" name:"verbose" location:"params"` // Available values: 0, 1
}

func (i *DescribeInstancesInput) Validate() error {
	return nil
}

func String(v string) *string {
	return &v
}

func StringSlice(src []string) []*string {
	dst := make([]*string, len(src))
	for i := 0; i < len(src); i++ {
		dst[i] = &(src[i])
	}
	return dst
}

func Int(v int) *int {
	return &v
}

func TestBuilder(t *testing.T) {

	conf, err := config.NewDefault()
	assert.Nil(t, err)
	conf.Host = "api.qc.dev"

	builder := &Builder{}
	operation := &data.Operation{
		Config: conf,
		Properties: &InstanceServiceProperties{
			Zone: String("beta"),
		},
		APIName:       "DescribeInstances",
		ServiceName:   "Instance",
		RequestMethod: "GET",
		RequestURI:    "/DescribeInstances",
		StatusCodes: []int{
			200,
		},
	}
	inputValue := reflect.ValueOf(&DescribeInstancesInput{
		ImageID:       StringSlice([]string{"img-xxxxxxxx", "img-zzzzzzzz"}),
		InstanceClass: Int(0),
		InstanceType:  StringSlice([]string{"type1", "type2"}),
		Instances:     StringSlice([]string{"i-xxxxxxxx", "i-zzzzzzzz"}),
		SearchWord:    String("search_word"),
		Status:        StringSlice([]string{"running"}),
		Tags:          StringSlice([]string{"tag1", "tag2"}),
		Verbose:       Int(1),
	})
	httpRequest, err := builder.BuildHTTPRequest(operation, &inputValue)
	assert.Nil(t, err)

	assert.True(t, strings.Contains(httpRequest.URL.String(), "https://api.qc.dev:443/iaas?"))
	assert.True(t, strings.Contains(httpRequest.URL.String(), "image_id.1=img-xxxxxxxx"))
	assert.True(t, strings.Contains(httpRequest.URL.String(), "instance_class=0"))
	assert.True(t, strings.Contains(httpRequest.URL.String(), "instance_type.1=type1"))
	assert.True(t, strings.Contains(httpRequest.URL.String(), "instances.2=i-zzzzzzzz"))
	assert.True(t, !strings.Contains(httpRequest.URL.String(), "limit"))
	assert.True(t, !strings.Contains(httpRequest.URL.String(), "offset"))
	assert.True(t, strings.Contains(httpRequest.URL.String(), "search_word=search_word"))
	assert.True(t, strings.Contains(httpRequest.URL.String(), "status.1=running"))
	assert.True(t, strings.Contains(httpRequest.URL.String(), "tags.1=tag1"))
	assert.True(t, strings.Contains(httpRequest.URL.String(), "verbose=1"))
	assert.True(t, strings.Contains(httpRequest.URL.String(), "zone=beta"))
}
