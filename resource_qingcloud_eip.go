package qingcloud

import (
	"fmt"
	// "log"
	// "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/magicshui/qingcloud-go/eip"
)

func resourceQingcloudEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudEipCreate,
		Read:   resourceQingcloudEipRead,
		Update: resourceQingcloudEipUpdate,
		Delete: resourceQingcloudEipDelete,
		Schema: resourceQingcloudEipSchema(),
	}
}

func resourceQingcloudEipCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip

	// TODO: 这个地方以后需要判断错误
	eipBandwidth := d.Get("bandwidth").(int)
	eipBillingMode := d.Get("billing_mode").(string)
	eipName := d.Get("name").(string)
	eipNeedIcp := d.Get("need_icp").(int)

	params := eip.AllocateEipsRequest{}
	params.Bandwidth.Set(eipBandwidth)
	params.BillingMode.Set(eipBillingMode)
	params.EipName.Set(eipName)
	params.NeedIcp.Set(eipNeedIcp)

	resp, err := clt.AllocateEips(params)
	if err != nil {
		return fmt.Errorf("Error create eip ", err)
	}

	description := d.Get("description").(string)
	if description != "" {
		modifyAtrributes := eip.ModifyEipAttributesRequest{}

		modifyAtrributes.Eip.Set(resp.Eips[0])
		modifyAtrributes.Description.Set(description)
		_, err := clt.ModifyEipAttributes(modifyAtrributes)
		if err != nil {
			return fmt.Errorf("Error modify eip description: %s", err)
		}
	}

	d.SetId(resp.Eips[0])
	return resourceQingcloudEipRead(d, meta)
}

func resourceQingcloudEipRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip

	// 设置请求参数
	params := eip.DescribeEipsRequest{}
	params.EipsN.Add(d.Id())
	params.Verbose.Set(1)

	resp, err := clt.DescribeEips(params)
	if err != nil {
		return fmt.Errorf("Error retrieving eips: %s", err)
	}
	for _, sg := range resp.EipSet {
		if sg.EipID == d.Id() {
			d.Set("name", sg.EipName)
			d.Set("billing_mode", sg.BillingMode)
			d.Set("bandwidth", sg.Bandwidth)
			d.Set("need_icp", sg.NeedIcp)
			d.Set("description", sg.Description)
			// 如下状态是稍等来获取的
			d.Set("addr", sg.EipAddr)
			d.Set("status", sg.Status)
			d.Set("transition_status", sg.TransitionStatus)
			return nil
		}
	}
	d.SetId("")
	return nil
}

func resourceQingcloudEipDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip
	// 判断当前的防火墙是否有人在使用
	describeParams := eip.DescribeEipsRequest{}
	describeParams.EipsN.Add(d.Id())
	describeParams.Verbose.Set(1)

	resp, err := clt.DescribeEips(describeParams)
	if err != nil {
		return fmt.Errorf("Error retrieving eip: %s", err)
	}
	for _, sg := range resp.EipSet {
		if sg.EipID == d.Id() {
			if sg.Resource.ResourceID != "" {
				// 如果公网IP正与其他资源绑定，则需要先解绑，再释放， 保证被释放的IP处于“可用”（ available ）状态。
				return fmt.Errorf("Current eip is in using", nil)
			}
		}
	}

	params := eip.ReleaseEipsRequest{}
	params.EipsN.Add(d.Id())
	_, err = clt.ReleaseEips(params)
	if err != nil {
		return fmt.Errorf("Error delete eip %s", err)
	}
	d.SetId("")
	return nil
}

func resourceQingcloudEipUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip

	if !d.HasChange("name") && !d.HasChange("description") {
		return nil
	}

	params := eip.ModifyEipAttributesRequest{}
	if d.HasChange("description") {
		params.Description.Set(d.Get("description").(string))
	}
	if d.HasChange("name") {
		params.EipName.Set(d.Get("name").(string))
	}
	params.Eip.Set(d.Id())
	_, err := clt.ModifyEipAttributes(params)
	if err != nil {
		return fmt.Errorf("Error modify eip %s", d.Id())
	}
	return resourceQingcloudEipRead(d, meta)
}

func resourceQingcloudEipSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"bandwidth": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			// TODO: only two types
		},
		"billing_mode": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "traffic",
		},
		"need_icp": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			Default:  0,
		},

		"addr": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},

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
