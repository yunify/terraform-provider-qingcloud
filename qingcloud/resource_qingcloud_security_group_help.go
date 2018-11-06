package qingcloud

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yunify/qingcloud-sdk-go/client"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifySecurityGroupAttributes(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.ModifySecurityGroupAttributesInput)
	input.SecurityGroup = qc.String(d.Id())
	nameUpdate := false
	descriptionUpdate := false
	input.SecurityGroupName, nameUpdate = getNamePointer(d)
	input.Description, descriptionUpdate = getDescriptionPointer(d)
	if nameUpdate || descriptionUpdate {
		var err error
		simpleRetry(func() error {
			_, err = clt.ModifySecurityGroupAttributes(input)
			return isServerBusy(err)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func applySecurityGroupRule(sgID *string, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.ApplySecurityGroupInput)
	input.SecurityGroup = sgID
	var output *qc.ApplySecurityGroupOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.ApplySecurityGroup(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	client.WaitJob(meta.(*QingCloudClient).job,
		qc.StringValue(output.JobID),
		time.Duration(waitJobTimeOutDefault)*time.Second, time.Duration(waitJobIntervalDefault)*time.Second)
	if _, err := SecurityGroupApplyTransitionStateRefresh(clt, sgID); err != nil {
		return err
	}
	return nil
}
