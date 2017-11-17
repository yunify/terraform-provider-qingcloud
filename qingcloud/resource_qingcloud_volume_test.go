package qingcloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func TestAccQingcloudVolume_basic(t *testing.T) {
	var volume qc.DescribeVolumesOutput
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "qingcloud_volume.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVolumeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVolumeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("qingcloud_volume.foo", &volume),
					resource.TestCheckResourceAttr(
						"qingcloud_volume.foo", "size", "10"),
				),
			},
			resource.TestStep{
				Config: testAccVolumeConfigTwo,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("qingcloud_volume.foo", &volume),
					resource.TestCheckResourceAttr(
						"qingcloud_volume.foo", "name", "volume"),
					resource.TestCheckResourceAttr(
						"qingcloud_volume.foo", "description", "volume"),
					resource.TestCheckResourceAttr(
						"qingcloud_volume.foo", "size", "20"),
				),
			},
		},
	})
}
func testAccCheckVolumeExists(n string, tag *qc.DescribeVolumesOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Volume ID is set")
		}

		client := testAccProvider.Meta().(*QingCloudClient)
		input := new(qc.DescribeVolumesInput)
		input.Volumes = []*string{qc.String(rs.Primary.ID)}
		d, err := client.volume.DescribeVolumes(input)

		log.Printf("[WARN] volume id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		if d == nil || len(d.VolumeSet) == 0 {
			return fmt.Errorf("volume not found")
		}

		*tag = *d
		return nil
	}
}

func testAccCheckVolumeDestroy(s *terraform.State) error {
	return testAccCheckVolumeDestroyWithProvider(s, testAccProvider)
}
func testAccCheckVolumeDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*QingCloudClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "qingcloud_volume" {
			continue
		}
		// Try to find the resource
		input := new(qc.DescribeVolumesInput)
		input.Volumes = []*string{qc.String(rs.Primary.ID)}
		output, err := client.volume.DescribeVolumes(input)
		if err == nil {
			if len(output.VolumeSet) != 0 && qc.StringValue(output.VolumeSet[0].Status)!="deleted" {
				return fmt.Errorf("Found  volume: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

const testAccVolumeConfig = `
resource "qingcloud_volume" "foo"{
	size = 10
}
`
const testAccVolumeConfigTwo = `
resource "qingcloud_volume" "foo"{
	size = 20
	name = "volume"
	description = "volume"
}
`
