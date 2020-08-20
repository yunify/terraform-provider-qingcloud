package qingcloud

import (
	"testing"

	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"os"
)

func TestAccQingcloudVxnet_importBasic(t *testing.T) {
	resourceName := "qingcloud_vxnet.foo"
	testTag := "terraform-test-vxnet-import-basic" + os.Getenv("CIRCLE_BUILD_NUM")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVxNetDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccVxNetConfigThree, testTag),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
