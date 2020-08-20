package qingcloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	qc "github.com/yunify/qingcloud-sdk-go/service"
	"os"
)

func TestAccQingcloudTag_basic(t *testing.T) {
	var tag qc.DescribeTagsOutput
	Tag1Name := os.Getenv("CIRCLE_BUILD_NUM") + "-tag-create"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "qingcloud_tag.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTagConfigTempalte, Tag1Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagExists("qingcloud_tag.foo", &tag),
					resource.TestCheckResourceAttr(
						"qingcloud_tag.foo", resourceName, Tag1Name),
					resource.TestCheckResourceAttr(
						"qingcloud_tag.foo", "color", "#9f9bb7"),
				),
			},
			{
				Config: fmt.Sprintf(testAccTagConfigTwoTemplate, Tag1Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagExists("qingcloud_tag.foo", &tag),
					resource.TestCheckResourceAttr(
						"qingcloud_tag.foo", resourceName, Tag1Name),
					resource.TestCheckResourceAttr(
						"qingcloud_tag.foo", resourceDescription, "test"),
					resource.TestCheckResourceAttr(
						"qingcloud_tag.foo", "color", "#fff"),
				),
			},
		},
	})
}

func testAccCheckTagExists(n string, tag *qc.DescribeTagsOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Tag ID is set")
		}

		client := testAccProvider.Meta().(*QingCloudClient)
		input := new(qc.DescribeTagsInput)
		input.Tags = []*string{qc.String(rs.Primary.ID)}
		d, err := client.tag.DescribeTags(input)

		log.Printf("[WARN] tag id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		if d == nil || qc.StringValue(d.TagSet[0].TagID) == "" {
			return fmt.Errorf("tag not found")
		}

		*tag = *d
		return nil
	}
}

func testAccCheckTagDestroy(s *terraform.State) error {
	return testAccCheckTagDestroyWithProvider(s, testAccProvider)
}

func testAccCheckTagDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*QingCloudClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "qingcloud_tag" {
			continue
		}

		// Try to find the resource
		input := new(qc.DescribeTagsInput)
		input.Tags = []*string{qc.String(rs.Primary.ID)}
		output, err := client.tag.DescribeTags(input)
		if err == nil {
			if len(output.TagSet) != 0 {
				return fmt.Errorf("Found  tag: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

const testAccTagConfigTempalte = `
resource "qingcloud_tag" "foo"{
	name="%v"
}
`
const testAccTagConfigTwoTemplate = `
resource "qingcloud_tag" "foo"{
	name="%v"
	description="test"
	color = "#fff"
}
`
