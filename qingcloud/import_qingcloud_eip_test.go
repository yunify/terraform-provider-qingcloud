package qingcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccQingcloudEip_importBasic(t *testing.T) {
	resourceName := "qingcloud_eip.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEIPConfigTwo,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
