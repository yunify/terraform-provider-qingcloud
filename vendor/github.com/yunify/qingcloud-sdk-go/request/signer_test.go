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
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yunify/qingcloud-sdk-go/utils"
)

func TestSigner0(t *testing.T) {
	url := "https://api.qc.dev/iaas?instance.0=i-xxxxxxxx&action=DescribeInstance&verbose=1"
	httpRequest, err := http.NewRequest("GET", url, nil)
	httpRequest.Header.Set("Date", utils.TimeToString(time.Time{}, "RFC 822"))
	assert.Nil(t, err)

	s := Signer{
		AccessKeyID:     "ENV_ACCESS_KEY_ID",
		SecretAccessKey: "ENV_SECRET_ACCESS_KEY",
	}
	err = s.WriteSignature(httpRequest)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(httpRequest.URL.String(), "https://api.qc.dev/iaas?"))
	assert.True(t, strings.Contains(
		httpRequest.URL.String(), "time_stamp=0001-01-01T00%3A00%3A00Z"))
	assert.True(t, strings.Contains(
		httpRequest.URL.String(), "signature=ZHa2iQ8PeyP1ktMF9C%2BDjOQBl537ti9RnYZ1Qqr6KRg%3D"))
}

func TestSigner1(t *testing.T) {
	url := "https://api.qc.dev/iaas/?action=RunInstances&count=1&image_id=centos64x86a&instance_name=demo&instance_type=small_b&login_mode=passwd&login_passwd=QingCloud20130712&signature_method=HmacSHA256&signature_version=1&time_stamp=2013-08-27T14%3A30%3A10Z&version=1&vxnets.1=vxnet-0&zone=pek1"
	httpRequest, err := http.NewRequest("GET", url, nil)
	timeValue, err := utils.StringToTime("2013-08-27T14:30:10Z", "ISO 8601")
	assert.Nil(t, err)
	httpRequest.Header.Set("Date", utils.TimeToString(timeValue, "RFC 822"))
	assert.Nil(t, err)

	s := Signer{
		AccessKeyID:     "QYACCESSKEYIDEXAMPLE",
		SecretAccessKey: "SECRETACCESSKEY",
	}
	err = s.WriteSignature(httpRequest)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(httpRequest.URL.String(), "https://api.qc.dev/iaas/?"))
	assert.True(t, strings.Contains(
		httpRequest.URL.String(), "time_stamp=2013-08-27T14%3A30%3A10Z"))
	assert.True(t, strings.Contains(
		httpRequest.URL.String(), "signature=32bseYy39DOlatuewpeuW5vpmW51sD1A%2FJdGynqSpP8%3D"))
}
