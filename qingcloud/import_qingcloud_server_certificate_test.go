package qingcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccQingcloudServerCertificate_importBasic(t *testing.T) {
	resourceName := "qingcloud_server_certificate.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServerCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerCertificateConfig,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"private_key",
					"certificate_content"},
			},
		},
	})
}
