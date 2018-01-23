/**
 * Copyright (c) 2016 Magicshui
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */
/**
 * Copyright (c) 2017 yunify
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

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
			resource.TestStep{
				Config: testAccServerCertificateConfig,
			},

			resource.TestStep{
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
