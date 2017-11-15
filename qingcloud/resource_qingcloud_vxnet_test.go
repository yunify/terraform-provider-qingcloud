package qingcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func testAccCheckVxNetExists(n string, eip *qc.DescribeVxNetsOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VxNet ID is set")
		}

		client := testAccProvider.Meta().(*QingCloudClient)
		input := new(qc.DescribeVxNetsInput)
		input.VxNets = []*string{qc.String(rs.Primary.ID)}
		d, err := client.vxnet.DescribeVxNets(input)

		log.Printf("[WARN] eip id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		if d == nil || len(d.VxNetSet) == 0 {
			return fmt.Errorf("VxNet not found")
		}

		*eip = *d
		return nil
	}
}
