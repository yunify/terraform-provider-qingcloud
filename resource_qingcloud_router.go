package qingcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/magicshui/qingcloud-go/router"
	"log"
	"time"
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

func RouterStateRefreshFunc(clt *router.ROUTER, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		params := router.DescribeRoutersRequest{}
		params.RoutersN.Add(id)
		params.Verbose.Set(1)

		resp, err := clt.DescribeRouters(params)
		if err != nil {
			return nil, "", err
		}
		return resp.RouterSet[0], resp.RouterSet[0].Status, nil
	}
}

// resourceQingcloudRouterCreate
func resourceQingcloudRouterCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router

	// TODO: 这个地方以后需要判断错误
	routerName := d.Get("name").(string)
	routerType := d.Get("type").(int)
	routerVPCNetwork := d.Get("vpc_network").(string)
	routerSecurityGroupID := d.Get("security_group_id").(string)

	params := router.CreateRoutersRequest{}
	params.RouterName.Set(routerName)
	params.RouterType.Set(routerType)
	params.VpcNetwork.Set(routerVPCNetwork)
	params.SecurityGroup.Set(routerSecurityGroupID)

	resp, err := clt.CreateRouters(params)
	if err != nil {
		return fmt.Errorf("Error create Router ", err)
	}
	d.SetId(resp.Routers[0])
	// 等待路由器创建成功
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"active"},
		Refresh:    RouterStateRefreshFunc(clt, resp.Routers[0]),
		Timeout:    10 * time.Minute,
		Delay:      20 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for router (%s) to start: %s", d.Id(), err)
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
	log.Printf("Fetch the router information: %s", resp)
	for _, v := range resp.RouterSet {
		if v.RouterID == d.Id() {
			d.Set("name", v.RouterName)
			d.Set("type", v.RouterType)
			d.Set("vpc_network", v.Vxnets)
			d.Set("security_group_id", v.SecurityGroupID)
			d.Set("description", v.Description)

			// 如下状态是稍等来获取的
			d.Set("vxnets", v.Vxnets)
			d.Set("private_ip", v.PrivateIP)
			d.Set("is_applied", v.IsApplied)
			d.Set("eip", v.Eip)
			d.Set("status", v.Status)
			d.Set("transition_status", v.TransitionStatus)
			return nil
		}
	}
	d.SetId("")
	return nil
}

func resourceQingcloudRouterDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceQingcloudRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router

	params := router.ModifyRouterAttributesRequest{}
	if !d.HasChange("description") && !d.HasChange("name") {
		return nil
	}
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
		"transition_status": &schema.Schema{
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
