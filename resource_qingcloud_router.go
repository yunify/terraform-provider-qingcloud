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
		Schema: resourceQingcloudRouterSchema(false),
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

	// description := d.Get("description").(string)
	// if description != "" {
	// 	modifyAtrributes := Router.ModifyRouterAttributesRequest{}

	// 	modifyAtrributes.Router.Set(resp.Routers[0])
	// 	modifyAtrributes.Description.Set(description)
	// 	_, err := clt.ModifyRouterAttributes(modifyAtrributes)
	// 	if err != nil {
	// 		return fmt.Errorf("Error modify Router description: %s", err)
	// 	}
	// }

	return resourceQingcloudRouterRead(d, meta)
}

func resourceQingcloudRouterRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router

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

	if !d.HasChange("description") && !d.HasChange("name") {
		return nil
	}

	clt := meta.(*QingCloudClient).router
	params := router.ModifyRouterAttributesRequest{}
	params.Router.Set(d.Id())

	if d.HasChange("description") {
		params.Description.Set(d.Get("description").(string))
	}
	if d.HasChange("name") {
		params.RouterName.Set(d.Get("name").(string))
	}

	_, err := clt.ModifyRouterAttributes(params)
	if err != nil {
		return fmt.Errorf("Error update router: %s", err)
	}
	return nil
}

func resourceQingcloudRouterSchema(computed bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
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

		// "vxnets": &schema.Schema{
		// 	Type:     schema.TypeList,
		// 	Computed: true,
		// 	ForceNew: true,
		// 	Elem: []*schema.Schema{
		// 		"nic_id": &schema.Schema{
		// 			Type:     schema.TypeString,
		// 			Required: !computed,
		// 			Computed: computed,
		// 		},
		// 		"vxnet_id": &schema.Schema{
		// 			Type:     schema.TypeString,
		// 			Required: !computed,
		// 			Computed: computed,
		// 		},
		// 	},
		// },

		// "eip": &schema.Schema{
		// 	Type:     schema.TypeList,
		// 	Computed: true,
		// 	ForceNew: true,
		// 	Elem: []*schema.Schema{
		// 		"name": &schema.Schema{
		// 			Type:     schema.TypeString,
		// 			Required: !computed,
		// 			Computed: computed,
		// 		},
		// 		"ip": &schema.Schema{
		// 			Type:     schema.TypeString,
		// 			Required: !computed,
		// 			Computed: computed,
		// 		},
		// 		"addr": &schema.Schema{
		// 			Type:     schema.TypeString,
		// 			Required: !computed,
		// 			Computed: computed,
		// 		},
		// 	},
		// 	Set: schema.HashString,
		// },

		"status": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}
