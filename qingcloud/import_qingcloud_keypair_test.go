package qingcloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccQingcloudKeyPair_importBasic(t *testing.T) {
	resourceName := "qingcloud_keypair.foo"
	testTag := "terraform-test-kepair-import-basic" + os.Getenv("CIRCLE_BUILD_NUM")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccKeypairConfigTwo, testTag),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
