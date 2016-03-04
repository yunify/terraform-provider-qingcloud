package qingcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/magicshui/qingcloud-go/router"
)

func resourceQingcloudRouter() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudRouterCreate,
		Read:   resourceQingcloudRouterRead,
		Update: resourceQingcloudRouterUpdate,
		Delete: nil,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vpc_network": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_applied": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

// resourceQingcloudRouterCreate
func resourceQingcloudRouterCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router

	params := router.CreateRoutersRequest{}
	params.RouterName.Set(d.Get("name").(string))
	params.RouterType.Set(d.Get("type").(int))
	params.VpcNetwork.Set(d.Get("vpc_network").(string))
	params.SecurityGroup.Set(d.Get("security_group_id").(string))

	resp, err := clt.CreateRouters(params)
	if err != nil {
		return fmt.Errorf("Error create Router ", err)
	}
	d.SetId(resp.Routers[0])

	_, err = RouterTransitionStateRefresh(clt, d.Id())
	if err != nil {
		return fmt.Errorf("Error waiting for router (%s) to start: %s", d.Id(), err)
	}

	if err := modifyRouterAttributes(d, meta, false); err != nil {
		return err
	}

	return resourceQingcloudRouterRead(d, meta)
}

func resourceQingcloudRouterRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router

	if _, err := RouterTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	// 设置请求参数
	params := router.DescribeRoutersRequest{}
	params.RoutersN.Add(d.Id())
	params.Verbose.Set(1)

	resp, err := clt.DescribeRouters(params)
	if err != nil {
		return fmt.Errorf("Error retrieving Routers: %s", err)
	}

	for _, v := range resp.RouterSet {
		if v.RouterID == d.Id() {
			d.Set("name", v.RouterName)
			d.Set("type", v.RouterType)
			d.Set("vpc_network", v.Vxnets)
			d.Set("security_group_id", v.SecurityGroupID)
			d.Set("description", v.Description)

			// 如下状态是稍等来获取的
			// d.Set("vxnets", v.Vxnets)
			d.Set("private_ip", v.PrivateIP)
			d.Set("is_applied", v.IsApplied)
			d.Set("eip", v.Eip)
			d.Set("status", v.Status)
			return nil
		}
	}
	d.SetId("")
	return nil
}

func resourceQingcloudRouterDelete(d *schema.ResourceData, meta interface{}) error {

	clt := meta.(*QingCloudClient).router

	if _, err := RouterTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}

	params := router.DeleteRoutersRequest{}
	params.RoutersN.Add(d.Id())

	_, err := clt.DeleteRouters(params)
	if err != nil {
		return err
	}

	// 等待状态变化
	_, err = RouterTransitionStateRefresh(clt, d.Id())
	if err != nil {
		return err
	}

	return nil
}

func resourceQingcloudRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router

	if !d.HasChange("description") && !d.HasChange("name") {
		return nil
	}

	if _, err := RouterTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}

	return modifyRouterAttributes(d, meta, false)
}
