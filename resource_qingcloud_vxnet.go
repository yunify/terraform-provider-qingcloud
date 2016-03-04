package qingcloud

import (
	"fmt"
	// "log"

	// "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/magicshui/qingcloud-go/router"
	"github.com/magicshui/qingcloud-go/vxnet"
)

func resourceQingcloudVxnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudVxnetCreate,
		Read:   resourceQingcloudVxnetRead,
		Update: resourceQingcloudVxnetUpdate,
		Delete: resourceQingcloudVxnetDelete,
		Schema: resourceQingcloudVxnetSchema(),
	}
}

func resourceQingcloudVxnetCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).vxnet

	params := vxnet.CreateVxnetsRequest{}
	params.VxnetName.Set(d.Get("name").(string))
	params.VxnetType.Set(d.Get("type").(int))

	resp, err := clt.CreateVxnets(params)
	if err != nil {
		return fmt.Errorf("Error create security group", err)
	}

	// waiting until state refresh
	if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, d.Get("router_id").(string)); err != nil {
		return err
	}

	if description := d.Get("description").(string); description != "" {
		modifyAtrributes := vxnet.ModifyVxnetAttributesRequest{}
		// 对于私有网络，一个定义文件只创建一个比较方便
		modifyAtrributes.Vxnet.Set(resp.Vxnets[0])
		modifyAtrributes.Description.Set(description)
		_, err := clt.ModifyVxnetAttributes(modifyAtrributes)
		if err != nil {
			// 这里可以不用返回错误
			return fmt.Errorf("Error modify vxnet description: %s", err)
		}
	}

	// waiting until state refresh
	if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, d.Get("router_id").(string)); err != nil {
		return err
	}

	// join the router
	joinRouterParams := router.JoinRouterRequest{}
	joinRouterParams.Vxnet.Set(resp.Vxnets[0])
	joinRouterParams.Router.Set(d.Get("router_id").(string))
	joinRouterParams.IpNetwork.Set(d.Get("ip_network").(string))

	clt2 := meta.(*QingCloudClient).router
	_, err = clt2.JoinRouter(joinRouterParams)
	if err != nil {
		return fmt.Errorf("Error modify vxnet description: %s", err)
	}
	d.SetId(resp.Vxnets[0])

	return nil
}

func resourceQingcloudVxnetRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).vxnet

	// 设置请求参数
	params := vxnet.DescribeVxnetsRequest{}
	params.VxnetsN.Add(d.Id())
	params.Verbose.Set(1)

	resp, err := clt.DescribeVxnets(params)
	if err != nil {
		return fmt.Errorf("Error retrieving vxnets: %s", err)
	}
	for _, sg := range resp.VxnetSet {
		if sg.VxnetID == d.Id() {
			d.Set("name", sg.VxnetName)
			d.Set("description", sg.Description)
			d.Set("router_id", sg.Router.RouterID)
			d.Set("ip_network", sg.Router.IPNetwork)
			return nil
		}
	}
	d.SetId("")
	return nil
}

func resourceQingcloudVxnetDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).vxnet
	// 判断当前的防火墙是否有人在使用
	describeParams := vxnet.DescribeVxnetsRequest{}
	describeParams.VxnetsN.Add(d.Id())
	describeParams.Verbose.Set(1)

	resp, err := clt.DescribeVxnets(describeParams)
	if err != nil {
		return fmt.Errorf("Error retrieving vxnet: %s", err)
	}
	for _, sg := range resp.VxnetSet {
		if sg.VxnetID == d.Id() {
			if len(sg.InstanceIds) > 0 {
				// 只能删除没有主机的私有网络，若删除时仍然有主机在此网络中，会返回错误信息。 可通过 LeaveVxnet 移出主机。
				return fmt.Errorf("Current vxnet is in using", nil)
			}
		}
	}

	params := vxnet.DeleteVxnetsRequest{}
	params.VxnetsN.Add(d.Id())
	_, err = clt.DeleteVxnets(params)
	if err != nil {
		return fmt.Errorf("Error delete vxnet %s", err)
	}
	d.SetId("")
	return nil
}

func resourceQingcloudVxnetUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).vxnet

	if !d.HasChange("name") && !d.HasChange("description") {
		return nil
	}

	params := vxnet.ModifyVxnetAttributesRequest{}
	if d.HasChange("description") {
		params.Description.Set(d.Get("description").(string))
	}
	if d.HasChange("name") {
		params.VxnetName.Set(d.Get("name").(string))
	}
	params.Vxnet.Set(d.Id())
	_, err := clt.ModifyVxnetAttributes(params)
	if err != nil {
		return fmt.Errorf("Error modify vxnet %s", d.Id())
	}
	return nil
}

func resourceQingcloudVxnetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"type": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			// TODO: only two types
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"router_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"ip_network": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},

		"id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}
