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

package main

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/godog"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

// QingCloudServiceFeatureContext provides feature context for QingCloudService.
func QingCloudServiceFeatureContext(s *godog.Suite) {
	s.Step(`^initialize QingCloud service$`, initializeQingCloudService)
	s.Step(`^the QingCloud service is initialized$`, theQingCloudServiceIsInitialized)

	s.Step(`^describe zones$`, describeZones)
	s.Step(`^describe zones should get (\d+) zone at least$`, describeZonesShouldGetZoneAtLeast)
	s.Step(`^describe zones should have the zone I\'m using$`, describeZonesShouldHaveTheZoneIamUsing)
}

// --------------------------------------------------------------------------

func initializeQingCloudService() error {
	return nil
}

func theQingCloudServiceIsInitialized() error {
	if qcService == nil {
		return errors.New("QingCloud service is not initialized")
	}
	return nil
}

// --------------------------------------------------------------------------

var describeZonesOutput *qc.DescribeZonesOutput

func describeZones() error {
	describeZonesOutput, err = qcService.DescribeZones(nil)
	return err
}

func describeZonesShouldGetZoneAtLeast(count int) error {
	if len(describeZonesOutput.ZoneSet) >= count {
		return nil
	}
	return fmt.Errorf("DescribeZones only got \"%d\" zone(s)", count)
}

func describeZonesShouldHaveTheZoneIamUsing() error {
	for _, zone := range describeZonesOutput.ZoneSet {
		if qc.StringValue(zone.ZoneID) == tc.Zone {
			return nil
		}
	}

	return fmt.Errorf("DescribeZones dosen't have zone \"%s\"", tc.Zone)
}
