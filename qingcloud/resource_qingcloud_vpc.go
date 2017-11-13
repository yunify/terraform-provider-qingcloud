package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func resourceQingcloudVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudVpcCreate,
		Read:   resourceQingcloudVpcRead,
		Update: resourceQingcloudVpcUpdate,
		Delete: resourceQingcloudVpcDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of Vpc",
			},
			"type": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: withinArrayInt(0, 1, 2),
				Description: "Type of Vpc: 0 - medium,1 - small,2 - large,default 1	",
			},
			"vpc_network": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: withinArrayString("192.168.0.0/16", "172.16.0.0/16", "172.17.0.0/16",
					"172.18.0.0/16", "172.19.0.0/16", "172.20.0.0/16", "172.21.0.0/16", "172.22.0.0/16",
					"172.23.0.0/16", "172.24.0.0/16", "172.25.0.0/16"),
				Description: "Network address range of vpc.",
			},
			"tag_ids": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The tag ids' id used by the vpc",
			},
			"tag_names": &schema.Schema{
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "compute by tag ids ",
			},
			"eip_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The eip's id used by the vpc",
			},
			"security_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The security group's id used by the vpc",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of vpc",
			},
			"private_ip": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The private ip of vpc",
			},
			"public_ip": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public ip of vpc",
			},
		},
	}
}

// resourceQingcloudRouterCreate
func resourceQingcloudVpcCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.CreateRoutersInput)
	input.RouterName = qc.String(d.Get("name").(string))
	if d.Get("name").(string) != "" {
		input.RouterName = qc.String(d.Get("name").(string))
	}
	input.RouterType = qc.Int(d.Get("type").(int))
	if d.Get("vpc_network").(string) != "" {
		input.VpcNetwork = qc.String(d.Get("vpc_network").(string))
	}
	if d.Get("security_group_id").(string) != "" {
		input.SecurityGroup = qc.String(d.Get("security_group_id").(string))
	}
	input.Count = qc.Int(1)
	output, err := clt.CreateRouters(input)
	if err != nil {
		return fmt.Errorf("Error create router: %s", err.Error())
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error create router: %s", *output.Message)
	}
	d.SetId(qc.StringValue(output.Routers[0]))
	_, err = RouterTransitionStateRefresh(clt, d.Id())
	if err != nil {
		return fmt.Errorf("Error waiting for router (%s) to start: %s", d.Id(), err.Error())
	}
	return resourceQingcloudVpcUpdate(d, meta)
}

func resourceQingcloudVpcRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.DescribeRoutersInput)
	input.Routers = []*string{qc.String(d.Id())}
	input.Verbose = qc.Int(1)
	output, err := clt.DescribeRouters(input)
	if err != nil {
		return fmt.Errorf("Error describe router: %s", err)
	}
	if *output.RetCode != 0 {
		return fmt.Errorf("Error describe router: %s", *output.Message)
	}
	if len(output.RouterSet) == 0 {
		d.SetId("")
		return nil
	}
	rtr := output.RouterSet[0]
	d.Set("name", qc.StringValue(rtr.RouterName))
	d.Set("type", qc.IntValue(rtr.RouterType))
	d.Set("security_group_id", qc.StringValue(rtr.SecurityGroupID))
	d.Set("description", qc.StringValue(rtr.Description))
	d.Set("private_ip", qc.StringValue(rtr.PrivateIP))
	d.Set("eip_id", qc.StringValue(rtr.EIP.EIPID))
	d.Set("public_ip", qc.StringValue(rtr.EIP.EIPAddr))
	return nil
}

func resourceQingcloudVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	if _, err := RouterTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	waitRouterLease(d, meta)
	d.Partial(true)
	if err := modifyRouterAttributes(d, meta); err != nil {
		return err
	}
	d.SetPartial("name")
	d.SetPartial("description")
	if d.HasChange("eip_id") {
		input := new(qc.UpdateRoutersInput)
		input.Routers = []*string{qc.String(d.Id())}
		output, err := clt.UpdateRouters(input)
		if err != nil {
			return fmt.Errorf("Error update router: %s", err.Error())
		}
		if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
			return fmt.Errorf("Error update router: %s", *output.Message)
		}
		_, err = RouterTransitionStateRefresh(clt, d.Id())
		if err != nil {
			return fmt.Errorf("Error waiting for router (%s) to start: %s", d.Id(), err.Error())
		}
	}
	d.SetPartial("eip_id")
	if d.HasChange("security_group_id") && !d.IsNewResource() {
		sgClt := meta.(*QingCloudClient).securitygroup
		input := new(qc.ApplySecurityGroupInput)
		input.SecurityGroup = qc.String(d.Get("security_group_id").(string))
		output, err := sgClt.ApplySecurityGroup(input)
		if err != nil {
			return fmt.Errorf("Error apply security group (%s) update %s", *input.SecurityGroup, *output.Message)
		}
	}
	d.SetPartial("security_group_id")
	if err := resourceUpdateTag(d, meta, qingcloudResourceTypeRouter); err != nil {
		return err
	}
	d.SetPartial("tag_ids")
	d.Partial(false)
	return resourceQingcloudRouterRead(d, meta)
}

func resourceQingcloudVpcDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	if _, err := RouterTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	waitRouterLease(d, meta)
	input := new(qc.DeleteRoutersInput)
	input.Routers = []*string{qc.String(d.Id())}
	output, err := clt.DeleteRouters(input)
	if err != nil {
		return fmt.Errorf("Error delete router: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error delete router: %s", *output.Message)
	}
	if _, err := RouterTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
