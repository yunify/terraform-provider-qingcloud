package qingcloud

import (
	"os"
	"testing"

	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccQingcloudEip_importBasic(t *testing.T) {
	resourceName := "qingcloud_eip.foo"
	testTag := "terraform-test-eip-import-basic" + os.Getenv("TRAVIS_BUILD_ID") + "-" + os.Getenv("TRAVIS_JOB_NUMBER")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccEIPConfigTwo, testTag),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
