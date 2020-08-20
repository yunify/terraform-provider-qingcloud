package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func resourceQingcloudSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudSecurityGroupCreate,
		Read:   resourceQingcloudSecurityGroupRead,
		Update: resourceQingcloudSecurityGroupUpdate,
		Delete: resourceQingcloudSecurityGroupDelete,
		Schema: map[string]*schema.Schema{
			resourceName: {
				Type:     schema.TypeString,
				Required: true,
			},
			resourceDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceTagIds:   tagIdsSchema(),
			resourceTagNames: tagNamesSchema(),
		},
	}
}

func resourceQingcloudSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.CreateSecurityGroupInput)
	input.SecurityGroupName, _ = getNamePointer(d)
	var output *qc.CreateSecurityGroupOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.CreateSecurityGroup(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	d.SetId(qc.StringValue(output.SecurityGroupID))
	return resourceQingcloudSecurityGroupUpdate(d, meta)
}

func resourceQingcloudSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.DescribeSecurityGroupsInput)
	input.SecurityGroups = []*string{qc.String(d.Id())}
	var output *qc.DescribeSecurityGroupsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeSecurityGroups(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if len(output.SecurityGroupSet) == 0 {
		d.SetId("")
		return nil
	}
	sg := output.SecurityGroupSet[0]
	d.Set(resourceName, qc.StringValue(sg.SecurityGroupName))
	d.Set(resourceDescription, qc.StringValue(sg.Description))
	if err := resourceSetTag(d, sg.Tags); err != nil {
		return err
	}
	return nil
}
func resourceQingcloudSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	if err := modifySecurityGroupAttributes(d, meta); err != nil {
		return err
	}
	d.SetPartial(resourceDescription)
	d.SetPartial(resourceName)
	if err := resourceUpdateTag(d, meta, qingcloudResourceTypeSecurityGroup); err != nil {
		return err
	}
	d.SetPartial(resourceTagIds)
	d.Partial(false)
	return resourceQingcloudSecurityGroupRead(d, meta)
}

func resourceQingcloudSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.DeleteSecurityGroupsInput)
	input.SecurityGroups = []*string{qc.String(d.Id())}
	var err error
	simpleRetry(func() error {
		_, err = clt.DeleteSecurityGroups(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
