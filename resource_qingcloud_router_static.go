package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/magicshui/qingcloud-go/router"
)

func resourceQingcloudRouterStatic() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudRouterStaticCreate,
		Read:   resourceQingcloudRouterStaticRead,
		Update: resourceQingcloudRouterStaticUpdate,
		Delete: resourceQingcloudRouterStaticDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"router": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"val1": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"val2": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"val3": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"val4": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"val5": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceQingcloudRouterStaticCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	// 确保没有在更新
	if _, err := RouterTransitionStateRefresh(clt, d.Get("router").(string)); err != nil {
		return err
	}

	params := router.AddRouterStaticsRequest{}
	params.Router.Set(d.Get("router").(string))
	params.StaticsNRouterStaticName.Add(d.Get("name").(string))
	params.StaticsNStaticType.Add(int64(d.Get("type").(int)))
	params.StaticsNVal1.Add(d.Get("val1").(string))
	params.StaticsNVal2.Add(d.Get("val2").(string))
	params.StaticsNVal3.Add(d.Get("val3").(string))
	params.StaticsNVal4.Add(d.Get("val4").(string))
	params.StaticsNVal5.Add(d.Get("val5").(string))
	resp, err := clt.AddRouterStatics(params)
	if err != nil {
		return err
	}
	d.SetId(resp.RouterStatics[0])

	return applyRouterUpdates(meta, d.Get("router").(string))
}
func resourceQingcloudRouterStaticRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	params := router.DescribeRouterStaticsRequest{}
	params.RouterStaticsN.Add(d.Id())
	resp, err := clt.DescribeRouterStatics(params)
	if err != nil {
		return err
	}
	rS := resp.RouterStaticSet[0]
	d.Set("router", rS.RouterID)
	d.Set("type", int(rS.StaticType))
	d.Set("val1", rS.Val1)
	d.Set("val2", rS.Val2)
	d.Set("val3", rS.Val3)
	d.Set("val4", rS.Val4)
	d.Set("val5", rS.Val5)
	return nil
}
func resourceQingcloudRouterStaticUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	if _, err := RouterTransitionStateRefresh(clt, d.Get("router").(string)); err != nil {
		return err
	}
	params := router.ModifyRouterStaticAttributesRequest{}
	params.RouterStatic.Set(d.Id())
	params.RouterStaticName.Set(d.Get("name").(string))
	params.Val1.Set(d.Get("val1").(string))
	params.Val2.Set(d.Get("val2").(string))
	params.Val3.Set(d.Get("val3").(string))
	params.Val4.Set(d.Get("val4").(string))
	params.Val5.Set(d.Get("val5").(string))
	_, err := clt.ModifyRouterStaticAttributes(params)
	if err != nil {
		return err
	}
	return applyRouterUpdates(meta, d.Get("router").(string))
}
func resourceQingcloudRouterStaticDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	if _, err := RouterTransitionStateRefresh(clt, d.Get("router").(string)); err != nil {
		return err
	}
	params := router.DeleteRouterStaticsRequest{}
	params.RouterStaticsN.Add(d.Id())
	_, err := clt.DeleteRouterStatics(params)
	if err != nil {
		return err
	}
	if err = applyRouterUpdates(meta, d.Get("router").(string)); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
