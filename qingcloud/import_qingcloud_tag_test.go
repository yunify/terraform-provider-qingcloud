package qingcloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccQingcloudTag_importBasic(t *testing.T) {
	resourceName := "qingcloud_tag.foo"
	tagName := os.Getenv("CIRCLE_BUILD_NUM") + "-tag-import"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTagConfigTempalte, tagName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
