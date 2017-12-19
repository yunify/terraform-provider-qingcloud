package qingcloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
	"os"
	"fmt"
)

func TestAccQingcloudInstance_importBasic(t *testing.T) {
	resourceName := "qingcloud_instance.foo"
	testTag := "terraform-test-instance-import-basic" + os.Getenv("TRAVIS_BUILD_ID") + "-" + os.Getenv("TRAVIS_JOB_NUMBER")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccInstanceConfig,testTag),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
