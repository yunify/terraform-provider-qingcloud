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

func TestAccQingcloudLoadBalancerBackend_basic(t *testing.T) {
	var lbb qc.DescribeLoadBalancerBackendsOutput
	testTag := "terraform-test-lb-backend-basic" + os.Getenv("CIRCLE_BUILD_NUM")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "qingcloud_loadbalancer_backend.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckLoadBalancerBackendDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccLBBConfigBasic, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerBackendExists(
						"qingcloud_loadbalancer_backend.foo", &lbb),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_backend.foo", "port", "80"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_backend.foo", "weight", "1"),
				),
			},
			{
				Config: fmt.Sprintf(testAccLBBConfigBasicTwo, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerBackendExists(
						"qingcloud_loadbalancer_backend.foo", &lbb),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_backend.foo", "port", "81"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_backend.foo", "weight", "2"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_backend.foo", "name", "test")),
			},
		},
	})
}

func testAccCheckLoadBalancerBackendDestroy(s *terraform.State) error {
	return testAccCheckLoadBalancerBackendDestroyWithProvider(s, testAccProvider)
}

func testAccCheckLoadBalancerBackendDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*QingCloudClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "qingcloud_loadbalancer_backend" {
			continue
		}
		input := new(qc.DescribeLoadBalancerBackendsInput)
		input.LoadBalancerBackends = []*string{qc.String(rs.Primary.ID)}
		output, err := client.loadbalancer.DescribeLoadBalancerBackends(input)
		if err == nil {
			if len(output.LoadBalancerBackendSet) != 0 {
				return fmt.Errorf("fount  loadbalancer backend: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckLoadBalancerBackendExists(n string, i *qc.DescribeLoadBalancerBackendsOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no loadbalancer backend ID is set")
		}

		client := testAccProvider.Meta().(*QingCloudClient)
		input := new(qc.DescribeLoadBalancerBackendsInput)
		input.Verbose = qc.Int(1)
		input.LoadBalancerBackends = []*string{qc.String(rs.Primary.ID)}
		d, err := client.loadbalancer.DescribeLoadBalancerBackends(input)

		log.Printf("[WARN] loadbalancer backend id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		if d == nil || len(d.LoadBalancerBackendSet) == 0 {
			return fmt.Errorf("Lb backend not found ")
		}

		*i = *d
		return nil
	}
}

const testAccLBBConfigBasic = `
resource "qingcloud_eip" "foo" {
  bandwidth = 2
  tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_loadbalancer" "foo" {
  eip_ids =["${qingcloud_eip.foo.id}"]
  tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
  name="%v"
}

resource "qingcloud_loadbalancer_listener" "foo"{
  load_balancer_id = "${qingcloud_loadbalancer.foo.id}"
  listener_port = "80"
  listener_protocol = "http"
}
resource "qingcloud_instance" "foo" {
	image_id = "centos7x64d"
	keypair_ids = ["${qingcloud_keypair.foo.id}"]
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_loadbalancer_backend" "foo" {
  loadbalancer_listener_id = "${qingcloud_loadbalancer_listener.foo.id}"
  port = 80
  resource_id = "${qingcloud_instance.foo.id}"
}
resource "qingcloud_keypair" "foo"{
	public_key = "    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyLSPqVIdXGH0QlGeWcPwa1fjTRKl6WtMiaSsP8/GnwjakDSKILUCoNe1yIpiK8F0/gmL71xaDQyfl7k6aE+gn6lSLUjpDmucAF1luGg6l7CIN+6hCqY3YqlAI05Tqwu0PdLAwCbGwdHcaWfECcbROJk5D0zpCTHmissrrAxdOv72g9Ple8KJ6C7F1tz6wmG0zUeineguGjW/PvfZiBDWZ/CyXGPeMDJxv3lrIiLa/ShgnQOxFTdHJPCw+F0/XlSzlIzP3gfni1vXxJWvYjdE9ULo7Z1DLWgZ73FCbeAvX/0e9C9jwT21Qa5RUy4pSP8m4WXSJgw2f9IpY1vIJFSZP root@centos1    "
	tag_ids = ["${qingcloud_tag.test.id}"]
}
`
const testAccLBBConfigBasicTwo = `
resource "qingcloud_eip" "foo" {
  bandwidth = 2
  tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_loadbalancer" "foo" {
  eip_ids =["${qingcloud_eip.foo.id}"]
  tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
  name="%v"
}

resource "qingcloud_loadbalancer_listener" "foo"{
  load_balancer_id = "${qingcloud_loadbalancer.foo.id}"
  listener_port = "80"
  listener_protocol = "http"
}
resource "qingcloud_instance" "foo" {
	image_id = "centos7x64d"
	keypair_ids = ["${qingcloud_keypair.foo.id}"]
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_loadbalancer_backend" "foo" {
  loadbalancer_listener_id = "${qingcloud_loadbalancer_listener.foo.id}"
  port = 81
  resource_id = "${qingcloud_instance.foo.id}"
  weight = 2 
  name = "test"
}
resource "qingcloud_keypair" "foo"{
	public_key = "    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyLSPqVIdXGH0QlGeWcPwa1fjTRKl6WtMiaSsP8/GnwjakDSKILUCoNe1yIpiK8F0/gmL71xaDQyfl7k6aE+gn6lSLUjpDmucAF1luGg6l7CIN+6hCqY3YqlAI05Tqwu0PdLAwCbGwdHcaWfECcbROJk5D0zpCTHmissrrAxdOv72g9Ple8KJ6C7F1tz6wmG0zUeineguGjW/PvfZiBDWZ/CyXGPeMDJxv3lrIiLa/ShgnQOxFTdHJPCw+F0/XlSzlIzP3gfni1vXxJWvYjdE9ULo7Z1DLWgZ73FCbeAvX/0e9C9jwT21Qa5RUy4pSP8m4WXSJgw2f9IpY1vIJFSZP root@centos1    "
	tag_ids = ["${qingcloud_tag.test.id}"]
}

`
