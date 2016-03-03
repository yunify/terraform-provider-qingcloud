package qingcloud

import (
	"fmt"

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

	// 创建
	params := eip.AllocateEipsRequest{}
	params.Bandwidth.Set(d.Get("bandwidth").(int))
	params.BillingMode.Set(d.Get("billing_mode").(string))
	params.EipName.Set(d.Get("name").(string))
	params.NeedIcp.Set(d.Get("need_icp").(int))
	resp, err := clt.AllocateEips(params)
	if err != nil {
		return fmt.Errorf("Error create eip ", err)
	}
	d.SetId(resp.Eips[0])

	// 设置描述信息
	if description := d.Get("description").(string); description != "" {
		modifyAtrributes := eip.ModifyEipAttributesRequest{}
		modifyAtrributes.Eip.Set(d.Id())
		modifyAtrributes.Description.Set(description)
		_, err := clt.ModifyEipAttributes(modifyAtrributes)
		if err != nil {
			return fmt.Errorf("Error modify eip description: %s", err)
		}
	}

	// 配置一下
	return resourceQingcloudEipRead(d, meta)
}

func resourceQingcloudEipRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip
	_, err := EipTransitionStateRefresh(clt, d.Id())
	if err != nil {
		return fmt.Errorf(
			"Error waiting for the transition %s", err)
	}

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
			if err := d.Set("resource", getEipSourceMap(sg)); err != nil {
				return fmt.Errorf("Error set eip resource %v", err)
			}
			return nil
		}
	}
	d.SetId(d.Id())
	return nil
}

func resourceQingcloudEipDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip

	params := eip.ReleaseEipsRequest{}
	params.EipsN.Add(d.Id())
	_, err := clt.ReleaseEips(params)
	if err != nil {
		return fmt.Errorf("Error delete eip %s", err)
	}
	d.SetId("")
	return nil
}

func resourceQingcloudEipUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip

	if !d.HasChange("name") && !d.HasChange("description") && !d.HasChange("bandwidth") && !d.HasChange("billing_mode") {
		return nil
	}

	if d.HasChange("bandwidth") {
		params := eip.ChangeEipsBandwidthRequest{}
		params.EipsN.Add(d.Id())
		params.Bandwidth.Set(d.Get("bandwidth").(int))
		_, err := clt.ChangeEipsBandwidth(params)
		if err != nil {
			return err
		}
	}
	if d.HasChange("billing_mode") {
		params := eip.ChangeEipsBillingModeRequest{}
		params.EipsN.Add(d.Id())
		params.BillingMode.Set(d.Get("billing_mode").(string))
		_, err := clt.ChangeEipsBillingMode(params)
		if err != nil {
			return err
		}
	}
	if d.HasChange("name") || d.HasChange("description") {
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
	}

	return resourceQingcloudEipRead(d, meta)
}

func resourceQingcloudEipSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "公网 IP 的名称",
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"bandwidth": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			// TODO: only two types
		},
		"billing_mode": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "traffic",
		},
		"need_icp": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
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
		// 目前正在使用这个 IP 的资源
		"resource": &schema.Schema{
			Type:         schema.TypeMap,
			Computed:     true,
			ComputedWhen: []string{"id"},
		},
		"id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}
