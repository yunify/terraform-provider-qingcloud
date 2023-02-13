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
					resource.TestCheckResourceAttr("qingcloud_instance_volume_attachment.foo", "volume_id", "vos-60xcbw8r"),
					resource.TestCheckResourceAttr("qingcloud_instance_volume_attachment.foo", "instance_id", "i-o1gm1smr"),
				),
			},
		},
	})
}

const testInstanceVolumeAttachmentConfig = `
resource "qingcloud_instance_volume_attachment" "foo"{
	volume_id = "vos-60xcbw8r"
	instance_id = "i-o1gm1smr"
}
`
