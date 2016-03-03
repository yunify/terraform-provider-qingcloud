package qingcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/magicshui/qingcloud-go/eip"
	"github.com/magicshui/qingcloud-go/router"
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

// Waiting for no transition_status
func EipTransitionStateRefresh(clt *eip.EIP, id string) *resource.StateChangeConf {
	refreshFunc := func() (interface{}, string, error) {
		params := eip.DescribeEipsRequest{}
		params.EipsN.Add(id)
		params.Verbose.Set(1)

		resp, err := clt.DescribeEips(params)
		if err != nil {
			return nil, "", err
		}
		return resp.EipSet[0], resp.EipSet[0].TransitionStatus, nil
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"associating", "dissociating", "suspending", "resuming", "releasing"},
		Target:     []string{""},
		Refresh:    refreshFunc,
		Timeout:    10 * time.Minute,
		Delay:      2 * time.Second,
		MinTimeout: 1 * time.Second,
	}
	return stateConf
}

// resourceQingcloudEipCreate 创建一个 Eip
func resourceQingcloudEipCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip

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

	return resourceQingcloudEipRead(d, meta)
}

func getEipSourceMap(data eip.Eip) map[string]interface{} {
	var a = make(map[string]interface{}, 3)
	a["id"] = data.Resource.ResourceID
	a["name"] = data.Resource.ResourceName
	a["type"] = data.Resource.ResourceType
	return a
}

func resourceQingcloudEipRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip
	_, err := EipTransitionStateRefresh(clt, d.Id()).WaitForState()
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
	d.SetId("")
	return nil
}

func dissociateEipFromResource(meta interface{}, eipID, resourceType, resourceID string) error {
	switch resourceType {
	case "router":
		log.Printf("[Debug] dissociate eip form resource %s", "router")
		clt := meta.(*QingCloudClient).router
		params := router.ModifyRouterAttributesRequest{}
		params.Eip.Set("")
		params.Router.Set(resourceID)
		_, err := clt.ModifyRouterAttributes(params)
		if err != nil {
			return err
		}
		p2 := router.UpdateRoutersRequest{}
		p2.RoutersN.Add(resourceID)
		_, err = clt.UpdateRouters(p2)
		if err != nil {
			return err
		}
		_, err = RouterTransitionStateRefresh(clt, resourceID)
		return err
	default:
		clt := meta.(*QingCloudClient).eip
		params := eip.DissociateEipsRequest{}
		params.EipsN.Add(eipID)
		_, err := clt.DissociateEips(params)
		if err != nil {
			return err
		}
		_, err = InstanceTransitionStateRefresh(meta.(*QingCloudClient).instance, resourceID)
		return err
	}
}

func resourceQingcloudEipDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip

	// 如果其他的资源在使用，那么首先解绑
	resource := d.Get("resource").(map[string]interface{})
	log.Printf("[Debug] eip resource %#v", resource)
	if resource["id"].(string) != "" {
		fmt.Printf("[DEBUG] Current eip is in using by %s %s", resource["type"].(string), resource["id"].(string))
		if err := dissociateEipFromResource(meta, d.Id(), resource["type"].(string), resource["id"].(string)); err != nil {
			return err
		}
	}

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
			return fmt.Errorf(
				"Error change eip bandwidth %s", err)
		}
	}
	if d.HasChange("billing_mode") {
		params := eip.ChangeEipsBillingModeRequest{}
		params.EipsN.Add(d.Id())
		params.BillingMode.Set(d.Get("billing_mode").(string))
		_, err := clt.ChangeEipsBillingMode(params)
		if err != nil {
			return fmt.Errorf(
				"Error change billing mode %s", err)
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
