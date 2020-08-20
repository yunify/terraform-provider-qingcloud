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

func TestAccQingcloudLoadBalancer_basic(t *testing.T) {
	var lb qc.DescribeLoadBalancersOutput
	testTag := "terraform-test-lb-basic" + os.Getenv("CIRCLE_BUILD_NUM")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "qingcloud_loadbalancer.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccLBConfigBasic, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(
						"qingcloud_loadbalancer.foo", &lb),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer.foo", "type", "0"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer.foo", "name", "test"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer.foo", "description", "test"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer.foo", "node_count", "1"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer.foo", "vxnet_id", "vxnet-0"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer.foo", "http_header_size", "15"),
				),
			},
			{
				Config: fmt.Sprintf(testAccLBConfigBasicTwo, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(
						"qingcloud_loadbalancer.foo", &lb),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer.foo", "type", "1"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer.foo", "name", "test1"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer.foo", "description", "test1"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer.foo", "node_count", "2"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer.foo", "vxnet_id", "vxnet-0"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer.foo", "http_header_size", "10"),
				),
			},
		},
	})
}

func TestAccQingcloudLoadBalancer_tag(t *testing.T) {
	var lb qc.DescribeLoadBalancersOutput
	lbTag1Name := "terraform-" + os.Getenv("CIRCLE_BUILD_NUM") + "-lb-tag1"
	lbTag2Name := "terraform-" + os.Getenv("CIRCLE_BUILD_NUM") + "-lb-tag2"
	testTagNameValue := func(names ...string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			tags := lb.LoadBalancerSet[0].Tags
			same_count := 0
			for _, tag := range tags {
				for _, name := range names {
					if qc.StringValue(tag.TagName) == name {
						same_count++
					}
					if same_count == len(lb.LoadBalancerSet[0].Tags) {
						return nil
					}
				}
			}
			return fmt.Errorf("tag name error %#v", names)
		}
	}
	testTagDetach := func() resource.TestCheckFunc {
		return func(state *terraform.State) error {
			if len(lb.LoadBalancerSet[0].Tags) != 0 {
				return fmt.Errorf("tag not detach ")
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "qingcloud_loadbalancer.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccLBConfigTagTemplate, lbTag1Name, lbTag2Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(
						"qingcloud_loadbalancer.foo", &lb),
					testTagNameValue(lbTag1Name, lbTag2Name),
				),
			},
			{
				Config: fmt.Sprintf(testAccLBConfigTagTwoTemplate, lbTag1Name, lbTag2Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(
						"qingcloud_loadbalancer.foo", &lb),
					testTagDetach(),
				),
			},
		},
	})
}

func TestAccQingcloudLoadBalancer_mutiEipsByCount(t *testing.T) {
	var lb qc.DescribeLoadBalancersOutput
	testTag := "terraform-test-lb-mutiEips" + os.Getenv("CIRCLE_BUILD_NUM")

	testCheck := func(eipCount int) resource.TestCheckFunc {
		return func(*terraform.State) error {
			if len(lb.LoadBalancerSet[0].Cluster) < 0 {
				return fmt.Errorf("no eips: %#v", lb.LoadBalancerSet[0].Cluster)
			}

			if len(lb.LoadBalancerSet[0].Cluster) != eipCount {
				return fmt.Errorf("eip count inconformity : %#v", lb.LoadBalancerSet[0].Cluster)
			}

			return nil
		}
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "qingcloud_loadbalancer.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccLBConfigMutiEips, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(
						"qingcloud_loadbalancer.foo", &lb),
					testCheck(1),
				),
			},
			{
				Config: fmt.Sprintf(testAccLBConfigMutiEipsTwo, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(
						"qingcloud_loadbalancer.foo", &lb),
					testCheck(3),
				),
			},
			{
				Config: fmt.Sprintf(testAccLBConfigMutiEipsThree, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(
						"qingcloud_loadbalancer.foo", &lb),
					testCheck(2),
				),
			},
		},
	})
}

func TestAccQingcloudLoadBalancer_inter_private_ip(t *testing.T) {
	var lb qc.DescribeLoadBalancersOutput
	testTag := "terraform-test-lb-inter-privateip-basic" + os.Getenv("CIRCLE_BUILD_NUM")

	testCheck := func(privateIp string) resource.TestCheckFunc {
		return func(*terraform.State) error {
			if qc.StringValue(lb.LoadBalancerSet[0].PrivateIPs[0]) != privateIp {
				return fmt.Errorf("private ip error: %#v", lb.LoadBalancerSet[0].PrivateIPs)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "qingcloud_loadbalancer.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccLBConfigInternal, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(
						"qingcloud_loadbalancer.foo", &lb),
					testCheck("192.168.0.3"),
				),
			},
			{
				Config: fmt.Sprintf(testAccLBConfigInternalTwo, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(
						"qingcloud_loadbalancer.foo", &lb),
					testCheck("192.168.0.4"),
				),
			},
		},
	})
}

func testAccCheckLoadBalancerDestroy(s *terraform.State) error {
	return testAccCheckLoadBalancerDestroyWithProvider(s, testAccProvider)
}

func testAccCheckLoadBalancerDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*QingCloudClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "qingcloud_loadbalancer" {
			continue
		}
		input := new(qc.DescribeLoadBalancersInput)
		input.LoadBalancers = []*string{qc.String(rs.Primary.ID)}
		output, err := client.loadbalancer.DescribeLoadBalancers(input)
		if err == nil {
			if !isLoadBalancerDeleted(output.LoadBalancerSet) {
				return fmt.Errorf("fount  loadbalancer: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckLoadBalancerExists(n string, i *qc.DescribeLoadBalancersOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LoadBalancer ID is set")
		}

		client := testAccProvider.Meta().(*QingCloudClient)
		input := new(qc.DescribeLoadBalancersInput)
		input.Verbose = qc.Int(1)
		input.LoadBalancers = []*string{qc.String(rs.Primary.ID)}
		d, err := client.loadbalancer.DescribeLoadBalancers(input)

		log.Printf("[WARN] loadbalancer id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		if d == nil || len(d.LoadBalancerSet) == 0 {
			return fmt.Errorf("Lb not found ")
		}

		*i = *d
		return nil
	}
}

const testAccLBConfigTagTemplate = `
resource "qingcloud_eip" "foo" {
    bandwidth = 2
}
resource "qingcloud_loadbalancer" "foo" {
	eip_ids =["${qingcloud_eip.foo.id}"]
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
const testAccLBConfigTagTwoTemplate = `
resource "qingcloud_eip" "foo" {
    bandwidth = 2
}
resource "qingcloud_loadbalancer" "foo" {
	eip_ids =["${qingcloud_eip.foo.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}
resource "qingcloud_tag" "test2"{
	name="%v"
}
`
const testAccLBConfigBasic = `
resource "qingcloud_eip" "foo" {
    bandwidth = 2
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_loadbalancer" "foo" {
	eip_ids =["${qingcloud_eip.foo.id}"]
    name = "test"
    description = "test"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}
`
const testAccLBConfigBasicTwo = `
resource "qingcloud_eip" "foo" {
    bandwidth = 2
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_loadbalancer" "foo" {
	eip_ids =["${qingcloud_eip.foo.id}"]
    name = "test1"
    description = "test1"
    http_header_size = 10
    node_count = 2
    type = 1
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}
`

const testAccLBConfigMutiEips = `
resource "qingcloud_eip" "foo1" {
    bandwidth = 2
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_eip" "foo2" {
    bandwidth = 2
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_eip" "foo3" {
    bandwidth = 2
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_loadbalancer" "foo" {
	eip_ids =["${qingcloud_eip.foo1.id}"]
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}
`
const testAccLBConfigMutiEipsTwo = `
resource "qingcloud_eip" "foo1" {
    bandwidth = 2
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_eip" "foo2" {
    bandwidth = 2
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_eip" "foo3" {
    bandwidth = 2
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_loadbalancer" "foo" {
	eip_ids =["${qingcloud_eip.foo1.id}","${qingcloud_eip.foo2.id}","${qingcloud_eip.foo3.id}"]
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}
`
const testAccLBConfigMutiEipsThree = `
resource "qingcloud_eip" "foo1" {
    bandwidth = 2
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_eip" "foo2" {
    bandwidth = 2
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_eip" "foo3" {
    bandwidth = 2
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_loadbalancer" "foo" {
	eip_ids =["${qingcloud_eip.foo1.id}","${qingcloud_eip.foo2.id}"]
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}
`

const testAccLBConfigInternal = `

resource "qingcloud_security_group" "foo" {
    name = "first_sg"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
	vpc_network = "192.168.0.0/16"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_vxnet" "foo" {
    type = 1
	vpc_id = "${qingcloud_vpc.foo.id}"
	ip_network = "192.168.0.0/24"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_loadbalancer" "foo" {
	vxnet_id = "${qingcloud_vxnet.foo.id}"
	private_ips = ["192.168.0.3"]
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}
`
const testAccLBConfigInternalTwo = `

resource "qingcloud_security_group" "foo" {
    name = "first_sg"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
	vpc_network = "192.168.0.0/16"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_vxnet" "foo" {
    type = 1
	vpc_id = "${qingcloud_vpc.foo.id}"
	ip_network = "192.168.0.0/24"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_loadbalancer" "foo" {
	vxnet_id = "${qingcloud_vxnet.foo.id}"
	private_ips = ["192.168.0.4"]
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}
`
