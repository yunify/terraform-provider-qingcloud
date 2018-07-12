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
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yunify/qingcloud-sdk-go/request/data"
	"github.com/yunify/qingcloud-sdk-go/request/errors"
)

func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

func IntValue(v *int) int {
	if v != nil {
		return *v
	}
	return 0
}

func TimeValue(v *time.Time) time.Time {
	if v != nil {
		return *v
	}
	return time.Time{}
}

func TestUnpackerUnpackHTTPRequest(t *testing.T) {
	type Instance struct {
		Device       *string `json:"device" name:"device"`
		InstanceID   *string `json:"instance_id" name:"instance_id"`
		InstanceName *string `json:"instance_name" name:"instance_name"`
	}

	type Volume struct {
		CreateTime       *time.Time `json:"create_time" name:"create_time" format:"ISO 8601"`
		Description      *string    `json:"description" name:"description"`
		Instance         *Instance  `json:"instance" name:"instance"`
		Size             *int       `json:"size" name:"size"`
		Status           *string    `json:"status" name:"status"`
		StatusTime       *time.Time `json:"status_time" name:"status_time" format:"ISO 8601"`
		SubCode          *int       `json:"sub_code" name:"sub_code"`
		TransitionStatus *string    `json:"transition_status" name:"transition_status"`
		VolumeID         *string    `json:"volume_id" name:"volume_id"`
		VolumeName       *string    `json:"volume_name" name:"volume_name"`
	}

	type DescribeVolumesOutput struct {
		StatusCode int `location:"statusCode"`
		Error      *errors.QingCloudError

		Action     *string   `json:"action" name:"action"`
		RetCode    *int      `json:"ret_code" name:"ret_code"`
		TotalCount *int      `json:"total_count" name:"total_count"`
		VolumeSet  []*Volume `json:"volume_set" name:"volume_set"`
		Message    *string   `json:"message" name:"message"`
	}

	httpResponse := &http.Response{Header: http.Header{}}
	httpResponse.StatusCode = 200
	httpResponse.Header.Set("Content-Type", "application/json")
	responseString := `{
	  "action": "DescribeVolumesResponse",
	  "total_count": 1024,
	  "volume_set": [
	    {
	      "status": "in-use",
	      "description": null,
	      "volume_name": "vol name",
	      "sub_code": 0,
	      "transition_status": "",
	      "instance": {
	        "instance_id": "i-xxxxxxxx",
	        "instance_name": "",
	        "device": "/dev/sdb"
	      },
	      "create_time": "2013-08-30T05:13:25Z",
	      "volume_id": "vol-xxxxxxxx",
	      "status_time": "2013-08-30T05:13:32Z",
	      "size": 10
	    }
	  ],
	  "ret_code": 0
	}`
	responseString = strings.Replace(responseString, ": ", ":", -1)
	httpResponse.Body = ioutil.NopCloser(bytes.NewReader([]byte(responseString)))

	output := &DescribeVolumesOutput{}
	outputValue := reflect.ValueOf(output)
	unpacker := Unpacker{}
	err := unpacker.UnpackHTTPRequest(&data.Operation{}, httpResponse, &outputValue)
	assert.Nil(t, err)
	assert.Equal(t, "i-xxxxxxxx", StringValue(output.VolumeSet[0].Instance.InstanceID))
	assert.Equal(t, "vol-xxxxxxxx", StringValue(output.VolumeSet[0].VolumeID))
	assert.Equal(t, "vol name", StringValue(output.VolumeSet[0].VolumeName))
	assert.Equal(t, 1024, IntValue(output.TotalCount))
	statusTime := time.Date(2013, 8, 30, 5, 13, 32, 0, time.UTC)
	assert.Equal(t, statusTime, TimeValue(output.VolumeSet[0].StatusTime))
}

func TestUnpacker_UnpackHTTPRequestWithError(t *testing.T) {
	type DescribeInstanceTypesOutput struct {
		RetCode *int    `json:"ret_code" name:"ret_code"`
		Message *string `json:"message" name:"message"`
	}

	httpResponse := &http.Response{Header: http.Header{}}
	httpResponse.StatusCode = 200
	httpResponse.Header.Set("Content-Type", "application/json")
	responseString := `{
  	  "message":"PermissionDenied, instance [i-xxxxxxxx] is not running, can not be stopped",
  	  "ret_code":1400
	}`
	httpResponse.Body = ioutil.NopCloser(bytes.NewReader([]byte(responseString)))

	output := &DescribeInstanceTypesOutput{}
	outputValue := reflect.ValueOf(output)
	unpacker := Unpacker{}
	err := unpacker.UnpackHTTPRequest(&data.Operation{}, httpResponse, &outputValue)
	assert.NotNil(t, err)
	switch e := err.(type) {
	case errors.QingCloudError:
		assert.Equal(t, e.RetCode, 1400)
	}
}
