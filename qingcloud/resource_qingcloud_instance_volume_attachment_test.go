package qingcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccQingcloudInstanceVolumeAttachment_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "qingcloud_instance_volume_attachment.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testInstanceVolumeAttachmentConfig),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("qingcloud_instance_volume_attachment.foo", "volume_id", "vos-6byu0d7a"),
					resource.TestCheckResourceAttr("qingcloud_instance_volume_attachment.foo", "instance_id", "i-tyotg828"),
				),
			},
		},
	})
}

const testInstanceVolumeAttachmentConfig = `
resource "qingcloud_instance_volume_attachment" "foo"{
  volume_id = "vos-6byu0d7a"
  instance_id = "i-tyotg828"
}
`
