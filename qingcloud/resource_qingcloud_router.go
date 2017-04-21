package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/lowstz/qingcloud-sdk-go/service"
)

func resourceQingcloudRouter() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudRouterCreate,
		Read:   resourceQingcloudRouterRead,
		Update: resourceQingcloudRouterUpdate,
		Delete: resourceQingcloudRouterDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "路由器名称",
			},
			"type": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: withinArrayInt(0, 1, 2),
				Description: "路由器类型: 0 - 中型，1 - 小型，2 - 大型，默认为 1	",
			},
			"vpc_network": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: withinArrayString("192.168.0.0/16", "172.16.0.0/16", "172.17.0.0/16",
					"172.18.0.0/16", "172.19.0.0/16", "172.20.0.0/16", "172.21.0.0/16", "172.22.0.0/16",
					"172.23.0.0/16", "172.24.0.0/16", "172.25.0.0/16"),
				Description: "VPC 网络地址范围，目前支持 192.168.0.0/16 或 172.16.0.0/16 。 注：此参数只在北京3区需要且是必填参数。",
			},
			"tag_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"tag_names": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"eip_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SecurityGroup ID",
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// resourceQingcloudRouterCreate
func resourceQingcloudRouterCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.CreateRoutersInput)
	input.RouterName = qc.String(d.Get("name").(string))
	input.RouterType = qc.Int(d.Get("type").(int))
	input.VpcNetwork = qc.String(d.Get("vpc_network").(string))
	input.SecurityGroup = qc.String(d.Get("security_group_id").(string))
	input.Count = qc.Int(1)
	output, err := clt.CreateRouters(input)
	if err != nil {
		return fmt.Errorf("Error create router: %s", err.Error())
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error create router: %s", *output.Message)
	}
	d.SetId(qc.StringValue(output.Routers[0]))
	d.Set("vpc_network", qc.String(d.Get("vpc_network").(string)))

	// qingcloudMutexKV.Lock(d.Id())
	// defer qingcloudMutexKV.Unlock(d.Id())

	_, err = RouterTransitionStateRefresh(clt, d.Id())
	if err != nil {
		return fmt.Errorf("Error waiting for router (%s) to start: %s", d.Id(), err.Error())
	}

	if err := modifyRouterAttributes(d, meta, true); err != nil {
		return err
	}
	if d.HasChange("eip_id") {
		input := new(qc.UpdateRoutersInput)
		input.Routers = []*string{qc.String(d.Id())}
		output, err := clt.UpdateRouters(input)
		if err != nil {
			return fmt.Errorf("Error update router: %s", err.Error())
		}
		if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
			return fmt.Errorf("Error update router input validate: %s", err)
		}
		_, err = RouterTransitionStateRefresh(clt, d.Id())
		if err != nil {
			return fmt.Errorf("Error waiting for router (%s) to start: %s", d.Id(), err.Error())
		}
	}
	if d.HasChange("security_group_id") {
		sgClt := meta.(*QingCloudClient).securitygroup
		input := new(qc.ApplySecurityGroupInput)
		input.SecurityGroup = qc.String(d.Get("security_group_id").(string))
		output, err := sgClt.ApplySecurityGroup(input)
		if err != nil {
			return fmt.Errorf("Error apply security group (%s) update %s", *input.SecurityGroup, *output.Message)
		}
	}
	// if err := resourceUpdateTag(d, meta, qingcloudResourceTypeRouter); err != nil {
	// 	return err
	// }
	return resourceQingcloudRouterRead(d, meta)
}

func resourceQingcloudRouterRead(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error router not found")
	}
	rtr := output.RouterSet[0]
	if rtr == nil {
		d.SetId("")
		return nil
	}
	d.Set("name", qc.StringValue(rtr.RouterName))
	d.Set("type", qc.IntValue(rtr.RouterType))
	d.Set("security_group_id", qc.StringValue(rtr.SecurityGroupID))
	d.Set("description", qc.StringValue(rtr.Description))
	d.Set("private_ip", qc.StringValue(rtr.PrivateIP))
	d.Set("eip_id", qc.StringValue(rtr.EIP.EIPID))
	d.Set("public_ip", qc.StringValue(rtr.EIP.EIPAddr))
	return nil
}

func resourceQingcloudRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	qingcloudMutexKV.Lock(d.Id())
	defer qingcloudMutexKV.Unlock(d.Id())
	if _, err := RouterTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	if err := modifyRouterAttributes(d, meta, false); err != nil {
		return err
	}
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
	if d.HasChange("security_group_id") {
		sgClt := meta.(*QingCloudClient).securitygroup
		input := new(qc.ApplySecurityGroupInput)
		input.SecurityGroup = qc.String(d.Get("security_group_id").(string))
		output, err := sgClt.ApplySecurityGroup(input)
		if err != nil {
			return fmt.Errorf("Error apply security group (%s) update %s", *input.SecurityGroup, *output.Message)
		}
	}
	if err := resourceUpdateTag(d, meta, qingcloudResourceTypeRouter); err != nil {
		return err
	}
	return resourceQingcloudRouterRead(d, meta)
}

func resourceQingcloudRouterDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	if d.Get("eip_id").(string) != "" {
		qingcloudMutexKV.Lock(d.Get("eip_id").(string))
		defer qingcloudMutexKV.Unlock(d.Get("eip_id").(string))
	}
	if _, err := RouterTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
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
