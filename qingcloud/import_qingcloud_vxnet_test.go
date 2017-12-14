package qingcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccQingcloudVxnet_importBasic(t *testing.T) {
	resourceName := "qingcloud_vxnet.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVxNetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVxNetConfigThree,
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
