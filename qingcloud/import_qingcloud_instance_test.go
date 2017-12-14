package qingcloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccQingcloudInstance_importBasic(t *testing.T) {
	resourceName := "qingcloud_instance.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccInstanceConfigKeyPair,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
