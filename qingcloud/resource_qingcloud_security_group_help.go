package qingcloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yunify/qingcloud-sdk-go/client"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifySecurityGroupAttributes(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.ModifySecurityGroupAttributesInput)
	input.SecurityGroup = qc.String(d.Id())
	attributeUpdate := false
	if d.HasChange("description") {
		if d.Get("description") == "" {
			input.Description = qc.String(" ")
		} else {
			input.Description = qc.String(d.Get("description").(string))
		}
		attributeUpdate = true
	}
	if d.HasChange("name") && !d.IsNewResource() {
		input.SecurityGroupName = qc.String(d.Get("name").(string))
		attributeUpdate = true
	}
	if attributeUpdate {
		_, err := clt.ModifySecurityGroupAttributes(input)
		if err != nil {
			return fmt.Errorf("Error modify security group attributes: %s ", err)
		}
	}
	return nil
}

func applySecurityGroupRule(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	sgID := d.Get("security_group_id").(string)
	input := new(qc.ApplySecurityGroupInput)
	input.SecurityGroup = qc.String(sgID)
	output, err := clt.ApplySecurityGroup(input)
	if err != nil {
		return err
	}
	client.WaitJob(meta.(*QingCloudClient).job,
		qc.StringValue(output.JobID),
		time.Duration(10)*time.Second, time.Duration(1)*time.Second)
	if _, err := SecurityGroupApplyTransitionStateRefresh(clt, sgID); err != nil {
		return err
	}
	return nil
}
