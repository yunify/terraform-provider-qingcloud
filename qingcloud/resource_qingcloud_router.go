package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
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
				ValidateFunc: withinArrayString("192.168.0.0/16", "172.16.0.0/16", "172.17.0.0/16",
					"172.18.0.0/16", "172.19.0.0/16", "172.20.0.0/16", "172.21.0.0/16", "172.22.0.0/16",
					"172.23.0.0/16", "172.24.0.0/16", "172.25.0.0/16"),
				Description: "VPC 网络地址范围，目前支持 192.168.0.0/16 或 172.16.0.0/16 。 注：此参数只在北京3区需要且是必填参数。",
			},
			"eip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"securitygroup": &schema.Schema{
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
	input.VpcNetwork = qc.String(d.Get("vpc_network").(int))
	input.SecurityGroup = qc.String(d.get("securitygroup").(int))
	input.Count = 1
	err := input.Validate
	if err != nil {
		return fmt.Errorf("Error create router input validate: %s", err.Error())
	}
	output, err := clt.CreateRouters(input)
	if err != nil {
		return fmt.Errorf("Error create router: %s", err.Error())
	}
	if output.RetCode != 0 {
		return fmt.Errorf("Error create router: %s", output.Message)
	}
	d.SetId(qc.StringValue(output.Routers[0]))
	d.Set("vpc_network", qc.String(d.Get("vpc_network").(string)))

	qingcloudMutexKV.Lock(d.Id())
	defer qingcloudMutexKV.Unlock(d.Id())

	_, err = RouterTransitionStateRefresh(clt, d.Id())
	if err != nil {
		return fmt.Errorf("Error waiting for router (%s) to start: %s", d.Id(), err.Error())
	}

	if err := modifyRouterAttributes(d, meta, true); err != nil {
		return err
	}
	if d.HasChange("eip") {
		input := new(qc.UpdateRoutersInput)
		input.Routers = []string{qc.String(d.ID())}
		err := input.Validate()
		if err != nil {
			return fmt.Errorf("Error update router input validate: %s", err.Error())
		}
		output, err = clt.UpdateRouters(input)
		if err != nil {
			return fmt.Errorf("Error update router: %s", err.Error())
		}
		_, err = RouterTransitionStateRefresh(clt, d.Id())
		if err != nil {
			return fmt.Errorf("Error waiting for router (%s) to start: %s", d.Id(), err.Error())
		}
	}
	if d.HasChange("securitygroup") {
		sgClt := meta.(*QingCloudClient).securitygroup
		input := new(qc.ApplySecurityGroupInput)
		input.SecurityGroup = qc.String(d.Get("securitygroup").(string))
		err := input.Validate()
		if err != nil {
			return fmt.Errorf("Error apply securitygroup (%s) update input validate: %s", *input.SecurityGroup, err.Error())
		}
		output, err := sgClt.ApplySecurityGroup(input)
		if err != nil {
			return fmt.Errorf("Error apply securitygroup (%s) update %s", *input.SecurityGroup, output.Message)
		}
		_, err = SecurityGroupTransitionStateRefresh(clt, *input.SecurityGroup)
		if err != nil {
			return fmt.Errorf("Error waiting for apply securitygroup (%s) to updated: %s", *input.SecurityGroup, err.Error())
		}
	}
	return resourceQingcloudRouterRead(d, meta)
}

func resourceQingcloudRouterRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.DescribeRoutersInput)
	input.Routers = []string{qc.String(d.Id())}
	input.Verbose = 1
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("Error describe router: %s", err)
	}
	output, err := clt.DescribeRouters(input)
	if err != nil {
		return fmt.Errorf("Error describe router: %s", err)
	}
	if *output.RetCode != 0 {
		return fmt.Errorf("Error describe router: %s", output.Message)
	}
	rtr := output.RouterSet[0]
	if rtr == nil {
		d.SetId("")
		return nil
	}
	d.SetId(qc.String(rtr.RouterID))
	d.Set("name", qc.String(rtr.RouterName))
	d.Set("type", qc.String(rtr.RouterType))
	d.Set("securitygroup", qc.String(rtr.SecurityGroupID))
	d.Set("description", qc.String(rtr.Description))
	d.Set("private_ip", qc.String(rtr.PrivateIP))
	d.Set("eip", qc.String(rtr.EIP.EIPID))
	d.Set("public_ip", qc.String(rtr.EIP.EIPAddr))
	return nil
}

func resourceQingcloudRouterDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	qingcloudMutexKV.Lock(d.Id())
	defer qingcloudMutexKV.Unlock(d.Id())
	if _, err := RouterTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	input := new(qc.DeleteRoutersInput)
	input.Routers = []string{qc.String(d.Id())}
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("Error delete router input validate: %s", err)
	}
	output, err := clt.DeleteRouters(input)
	if err != nil {
		return fmt.Errorf("Error delete router: %s", err)
	}
	if output.RetCode != 0 {
		return fmt.Errorf("Error delete router: %s", output.Message)
	}
	d.SetId("")
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
	if d.HasChange("eip") {
		input := new(qc.UpdateRoutersInput)
		input.Routers = []string{qc.String(d.ID())}
		err := input.Validate()
		if err != nil {
			return fmt.Errorf("Error update router input validate: %s", err.Error())
		}
		output, err = clt.UpdateRouters(input)
		if err != nil {
			return fmt.Errorf("Error update router: %s", err.Error())
		}
		_, err = RouterTransitionStateRefresh(clt, d.Id())
		if err != nil {
			return fmt.Errorf("Error waiting for router (%s) to start: %s", d.Id(), err.Error())
		}
	}
	if d.HasChange("securitygroup") {
		sgClt := meta.(*QingCloudClient).securitygroup
		input := new(qc.ApplySecurityGroupInput)
		input.SecurityGroup = qc.String(d.Get("securitygroup").(string))
		err := input.Validate()
		if err != nil {
			return fmt.Errorf("Error apply securitygroup (%s) update input validate: %s", *input.SecurityGroup, err.Error())
		}
		output, err := sgClt.ApplySecurityGroup(input)
		if err != nil {
			return fmt.Errorf("Error apply securitygroup (%s) update %s", *input.SecurityGroup, output.Message)
		}
		_, err = SecurityGroupTransitionStateRefresh(clt, *input.SecurityGroup)
		if err != nil {
			return fmt.Errorf("Error waiting for apply securitygroup (%s) to updated: %s", *input.SecurityGroup, err.Error())
		}
	}
	return resourceQingcloudRouterRead(d, meta)
}
