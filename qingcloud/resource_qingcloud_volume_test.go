package qingcloud

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func TestAccQingcloudVolume_basic(t *testing.T) {
	var volume qc.DescribeVolumesOutput
	testTag := "terraform-test-volume-basic" + os.Getenv("CIRCLE_BUILD_NUM")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "qingcloud_volume.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccVolumeConfig, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("qingcloud_volume.foo", &volume),
					resource.TestCheckResourceAttr(
						"qingcloud_volume.foo", "size", "10"),
				),
			},
			{
				Config: fmt.Sprintf(testAccVolumeConfigTwo, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("qingcloud_volume.foo", &volume),
					resource.TestCheckResourceAttr(
						"qingcloud_volume.foo", resourceName, "volume"),
					resource.TestCheckResourceAttr(
						"qingcloud_volume.foo", resourceDescription, "volume"),
					resource.TestCheckResourceAttr(
						"qingcloud_volume.foo", "size", "20"),
				),
			},
		},
	})
}

func TestAccQingcloudVolume_tag(t *testing.T) {
	var volume qc.DescribeVolumesOutput
	volumeTag1Name := "terraform-" + os.Getenv("CIRCLE_BUILD_NUM") + "-volume-tag1"
	volumeTag2Name := "terraform-" + os.Getenv("CIRCLE_BUILD_NUM") + "-volume-tag2"

	testTagNameValue := func(names ...string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			tags := volume.VolumeSet[0].Tags
			same_count := 0
			for _, tag := range tags {
				for _, name := range names {
					if qc.StringValue(tag.TagName) == name {
						same_count++
					}
					if same_count == len(volume.VolumeSet[0].Tags) {
						return nil
					}
				}
			}
			return fmt.Errorf("tag name error %#v", names)
		}
	}

	testTagDetach := func() resource.TestCheckFunc {
		return func(state *terraform.State) error {
			if len(volume.VolumeSet[0].Tags) != 0 {
				return fmt.Errorf("tag not detach ")
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "qingcloud_volume.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccVolumeConfigTagTemplate, volumeTag1Name, volumeTag2Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("qingcloud_volume.foo", &volume),
					testTagNameValue(volumeTag1Name, volumeTag1Name),
				),
			},
			{
				Config: fmt.Sprintf(testAccVolmeConfigTagTwoTemplate, volumeTag1Name, volumeTag2Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("qingcloud_volume.foo", &volume),
					testTagDetach(),
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
			if len(output.VolumeSet) != 0 && qc.StringValue(output.VolumeSet[0].Status) != "deleted" {
				return fmt.Errorf("Found  volume: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

const testAccVolumeConfig = `
resource "qingcloud_volume" "foo"{
	size = 10
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}
`
const testAccVolumeConfigTwo = `
resource "qingcloud_volume" "foo"{
	size = 20
	name = "volume"
	description = "volume"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}
`
const testAccVolumeConfigTagTemplate = `
resource "qingcloud_volume" "foo"{
	size = 10
	tag_ids = ["${qingcloud_tag.test.id}",
				"${qingcloud_tag.test2.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}
resource "qingcloud_tag" "test2"{
	name="%v"
}
`
const testAccVolmeConfigTagTwoTemplate = `
resource "qingcloud_volume" "foo"{
	size = 10
}
resource "qingcloud_tag" "test"{
	name="%v"
}
resource "qingcloud_tag" "test2"{
	name="%v"
}
`
