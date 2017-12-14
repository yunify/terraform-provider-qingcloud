package qingcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccQingcloudVolume_importBasic(t *testing.T) {
	resourceName := "qingcloud_volume.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVolumeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVolumeConfigTwo,
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

