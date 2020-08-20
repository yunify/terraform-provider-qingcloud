package qingcloud

import (
	"os"
	"testing"

	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccQingcloudEip_importBasic(t *testing.T) {
	resourceName := "qingcloud_eip.foo"
	testTag := "terraform-test-eip-import-basic" + os.Getenv("CIRCLE_BUILD_NUM")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccEIPConfigTwo, testTag),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
