package qingcloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccQingcloudKeyPair_importBasic(t *testing.T) {
	resourceName := "qingcloud_keypair.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKeypairDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKeypairConfigTwo,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
