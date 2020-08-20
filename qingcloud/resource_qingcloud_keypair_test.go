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

func TestAccQingcloudKeypair_basic(t *testing.T) {
	var keypair qc.DescribeKeyPairsOutput
	testTag := "terraform-test-keypair-basic" + os.Getenv("CIRCLE_BUILD_NUM")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "qingcloud_keypair.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccKeypairConfig, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeypairExists("qingcloud_keypair.foo", &keypair),
				),
			},
			{
				Config: fmt.Sprintf(testAccKeypairConfigTwo, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeypairExists("qingcloud_keypair.foo", &keypair),
					resource.TestCheckResourceAttr(
						"qingcloud_keypair.foo", resourceName, "keypair1"),
					resource.TestCheckResourceAttr(
						"qingcloud_keypair.foo", resourceDescription, "test"),
				),
			},
		},
	})
}
func TestAccQingcloudKeypair_tag(t *testing.T) {
	var keypair qc.DescribeKeyPairsOutput
	keypairTag1Name := "terraform-" + os.Getenv("CIRCLE_BUILD_NUM") + "-kp-tag1"
	keypairTag2Name := "terraform-" + os.Getenv("CIRCLE_BUILD_NUM") + "-kp-tag2"
	testTagNameValue := func(names ...string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			tags := keypair.KeyPairSet[0].Tags
			same_count := 0
			for _, tag := range tags {
				for _, name := range names {
					if qc.StringValue(tag.TagName) == name {
						same_count++
					}
					if same_count == len(keypair.KeyPairSet[0].Tags) {
						return nil
					}
				}
			}
			return fmt.Errorf("tag name error %#v", names)
		}
	}
	testTagDetach := func() resource.TestCheckFunc {
		return func(state *terraform.State) error {
			if len(keypair.KeyPairSet[0].Tags) != 0 {
				return fmt.Errorf("tag not detach ")
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "qingcloud_keypair.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccKeypairConfigTagTemplate, keypairTag1Name, keypairTag2Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeypairExists(
						"qingcloud_keypair.foo", &keypair),
					testTagNameValue(keypairTag1Name, keypairTag2Name),
				),
			},
			{
				Config: fmt.Sprintf(testAccKeypairConfigTagTwoTemplate, keypairTag1Name, keypairTag2Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeypairExists(
						"qingcloud_keypair.foo", &keypair),
					testTagDetach(),
				),
			},
		},
	})

}

func testAccCheckKeypairExists(n string, tag *qc.DescribeKeyPairsOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Keypair ID is set")
		}

		client := testAccProvider.Meta().(*QingCloudClient)
		input := new(qc.DescribeKeyPairsInput)
		input.KeyPairs = []*string{qc.String(rs.Primary.ID)}
		d, err := client.keypair.DescribeKeyPairs(input)

		log.Printf("[WARN] tag id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		if d == nil || qc.StringValue(d.KeyPairSet[0].KeyPairID) == "" {
			return fmt.Errorf("tag not found")
		}

		*tag = *d
		return nil
	}
}
func testAccCheckKeypairDestroy(s *terraform.State) error {
	return testAccCheckKeypairDestroyWithProvider(s, testAccProvider)
}
func testAccCheckKeypairDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*QingCloudClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "qingcloud_keypair" {
			continue
		}

		// Try to find the resource
		input := new(qc.DescribeKeyPairsInput)
		input.KeyPairs = []*string{qc.String(rs.Primary.ID)}
		output, err := client.keypair.DescribeKeyPairs(input)
		if err == nil {
			if len(output.KeyPairSet) != 0 {
				return fmt.Errorf("Found  keypair: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

const testAccKeypairConfig = `
resource "qingcloud_keypair" "foo"{
	public_key = "    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyLSPqVIdXGH0QlGeWcPwa1fjTRKl6WtMiaSsP8/GnwjakDSKILUCoNe1yIpiK8F0/gmL71xaDQyfl7k6aE+gn6lSLUjpDmucAF1luGg6l7CIN+6hCqY3YqlAI05Tqwu0PdLAwCbGwdHcaWfECcbROJk5D0zpCTHmissrrAxdOv72g9Ple8KJ6C7F1tz6wmG0zUeineguGjW/PvfZiBDWZ/CyXGPeMDJxv3lrIiLa/ShgnQOxFTdHJPCw+F0/XlSzlIzP3gfni1vXxJWvYjdE9ULo7Z1DLWgZ73FCbeAvX/0e9C9jwT21Qa5RUy4pSP8m4WXSJgw2f9IpY1vIJFSZP root@centos1    "
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}
`
const testAccKeypairConfigTwo = `
resource "qingcloud_keypair" "foo"{
	name="keypair1"
	description="test"
	public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCa7DDEP+CR9l26yJ5dAlmZvTwlZoxIZ2yD79dU/UNAlNzhdc+iqaC+aA2afyM4KCx8qtsawj5zzU0714fK+uA27MvQSwB0A25NqnJPgAw3v0WrOfFFG01Ecirc2MmMU2RHUk0cwZ5rVbcg8SUOwSs2tVKlWi98v1XcEw3vuM2ruPLkj8z9/Rf0o6FJ8vkpvsPXigFW82wkmI2WsgczvCbwApklaqdewC7Dxa0dFtA0gcqsgQzD4NR4glrHObyfxP3WRlPeyR7fFJRZFBqoLLELrqS5tYEpp6jVdzlHAf7WqOuLf0AoI+1Qsx57c92M0Rnj2MLs/6QNWKOVjzEfgXTD root@junwuhui.cn"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}
`
const testAccKeypairConfigTagTemplate = `

resource "qingcloud_keypair" "foo" {
	public_key = "    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyLSPqVIdXGH0QlGeWcPwa1fjTRKl6WtMiaSsP8/GnwjakDSKILUCoNe1yIpiK8F0/gmL71xaDQyfl7k6aE+gn6lSLUjpDmucAF1luGg6l7CIN+6hCqY3YqlAI05Tqwu0PdLAwCbGwdHcaWfECcbROJk5D0zpCTHmissrrAxdOv72g9Ple8KJ6C7F1tz6wmG0zUeineguGjW/PvfZiBDWZ/CyXGPeMDJxv3lrIiLa/ShgnQOxFTdHJPCw+F0/XlSzlIzP3gfni1vXxJWvYjdE9ULo7Z1DLWgZ73FCbeAvX/0e9C9jwT21Qa5RUy4pSP8m4WXSJgw2f9IpY1vIJFSZP root@centos1    "
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
const testAccKeypairConfigTagTwoTemplate = `

resource "qingcloud_keypair" "foo" {
	public_key = "    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyLSPqVIdXGH0QlGeWcPwa1fjTRKl6WtMiaSsP8/GnwjakDSKILUCoNe1yIpiK8F0/gmL71xaDQyfl7k6aE+gn6lSLUjpDmucAF1luGg6l7CIN+6hCqY3YqlAI05Tqwu0PdLAwCbGwdHcaWfECcbROJk5D0zpCTHmissrrAxdOv72g9Ple8KJ6C7F1tz6wmG0zUeineguGjW/PvfZiBDWZ/CyXGPeMDJxv3lrIiLa/ShgnQOxFTdHJPCw+F0/XlSzlIzP3gfni1vXxJWvYjdE9ULo7Z1DLWgZ73FCbeAvX/0e9C9jwT21Qa5RUy4pSP8m4WXSJgw2f9IpY1vIJFSZP root@centos1    "
}
resource "qingcloud_tag" "test"{
	name="%v"
}
resource "qingcloud_tag" "test2"{
	name="%v"
}
`
