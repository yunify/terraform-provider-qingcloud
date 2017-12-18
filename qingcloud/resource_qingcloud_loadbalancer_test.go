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

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "qingcloud_loadbalancer.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckLoadBalancerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccLBConfigBasic,
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
			resource.TestStep{
				Config: testAccLBConfigBasicTwo,
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
	lbTag1Name := os.Getenv("TRAVIS_BUILD_ID") + "-" + os.Getenv("TRAVIS_JOB_NUMBER") + "-lb-tag1"
	lbTag2Name := os.Getenv("TRAVIS_BUILD_ID") + "-" + os.Getenv("TRAVIS_JOB_NUMBER") + "-lb-tag2"
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
			resource.TestStep{
				Config: fmt.Sprintf(testAccLBConfigTagTemplate, lbTag1Name, lbTag2Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(
						"qingcloud_loadbalancer.foo", &lb),
					testTagNameValue(lbTag1Name, lbTag2Name),
				),
			},
			resource.TestStep{
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
	testCheck := func(eipCount int) resource.TestCheckFunc {
		return func(*terraform.State) error {
			if len(lb.LoadBalancerSet[0].Cluster) < 0 {
				return fmt.Errorf("no eips: %#v",lb.LoadBalancerSet[0].Cluster)
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
			resource.TestStep{
				Config: testAccLBConfigMutiEips,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(
						"qingcloud_loadbalancer.foo", &lb),
					testCheck(1),
				),
			},
			resource.TestStep{
				Config: testAccLBConfigMutiEipsTwo,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(
						"qingcloud_loadbalancer.foo", &lb),
					testCheck(3),
				),
			},
			resource.TestStep{
				Config: testAccLBConfigMutiEipsThree,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerExists(
						"qingcloud_loadbalancer.foo", &lb),
					testCheck(2),
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
}
resource "qingcloud_loadbalancer" "foo" {
	eip_ids =["${qingcloud_eip.foo.id}"]
    name = "test"
    description = "test"
}
`
const testAccLBConfigBasicTwo = `
resource "qingcloud_eip" "foo" {
    bandwidth = 2
}
resource "qingcloud_loadbalancer" "foo" {
	eip_ids =["${qingcloud_eip.foo.id}"]
    name = "test1"
    description = "test1"
    http_header_size = 10
    node_count = 2
    type = 1
}
`

const testAccLBConfigMutiEips = `
resource "qingcloud_eip" "foo1" {
    bandwidth = 2
}
resource "qingcloud_eip" "foo2" {
    bandwidth = 2
}
resource "qingcloud_eip" "foo3" {
    bandwidth = 2
}
resource "qingcloud_eip" "foo4" {
    bandwidth = 2
}
resource "qingcloud_loadbalancer" "foo" {
	eip_ids =["${qingcloud_eip.foo1.id}"]
}
`
const testAccLBConfigMutiEipsTwo = `
resource "qingcloud_eip" "foo1" {
    bandwidth = 2
}
resource "qingcloud_eip" "foo2" {
    bandwidth = 2
}
resource "qingcloud_eip" "foo3" {
    bandwidth = 2
}
resource "qingcloud_loadbalancer" "foo" {
	eip_ids =["${qingcloud_eip.foo1.id}","${qingcloud_eip.foo2.id}","${qingcloud_eip.foo3.id}"]
}
`
const testAccLBConfigMutiEipsThree = `
resource "qingcloud_eip" "foo1" {
    bandwidth = 2
}
resource "qingcloud_eip" "foo2" {
    bandwidth = 2
}
resource "qingcloud_eip" "foo3" {
    bandwidth = 2
}
resource "qingcloud_loadbalancer" "foo" {
	eip_ids =["${qingcloud_eip.foo1.id}","${qingcloud_eip.foo2.id}"]
}
`
