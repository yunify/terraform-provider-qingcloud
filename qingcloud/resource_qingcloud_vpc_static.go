package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

const (
	resourceVpcStaticVpcid = "vpc_id"
	resourceVpcStaticType  = "type"
	resourceVpcVal1        = "val1"
	resourceVpcVal2        = "val2"
	resourceVpcVal3        = "val3"
	resourceVpcVal4        = "val4"
	resourceVpcVal5        = "val5"
)

func resourceQingcloudVpcStatic() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudVpcStaticCreate,
		Read:   resourceQingcloudVpcStaticRead,
		Update: resourceQingcloudVpcStaticUpdate,
		Delete: resourceQingcloudVpcStaticDelete,
		Schema: map[string]*schema.Schema{
			resourceName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceVpcStaticVpcid: {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			resourceVpcStaticType: {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: withinArrayInt(1, 2, 3, 4, 6, 7, 8),
			},
			resourceVpcVal1: {
				Type:     schema.TypeString,
				Required: true,
			},
			resourceVpcVal2: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			resourceVpcVal3: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			resourceVpcVal4: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			resourceVpcVal5: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceQingcloudVpcStaticCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.AddRouterStaticsInput)
	static := new(qc.RouterStatic)
	input.Router = getSetStringPointer(d, resourceVpcStaticVpcid)
	static.RouterID = getSetStringPointer(d, resourceVpcStaticVpcid)
	static.RouterStaticName, _ = getNamePointer(d)
	static.StaticType = qc.Int(d.Get(resourceVpcStaticType).(int))
	static.Val1 = getSetStringPointer(d, resourceVpcVal1)
	static.Val2 = getSetStringPointer(d, resourceVpcVal2)
	static.Val3 = getSetStringPointer(d, resourceVpcVal3)
	static.Val4 = getSetStringPointer(d, resourceVpcVal4)
	static.Val5 = getSetStringPointer(d, resourceVpcVal5)
	input.Statics = []*qc.RouterStatic{static}
	var output *qc.AddRouterStaticsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.AddRouterStatics(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	d.SetId(qc.StringValue(output.RouterStatics[0]))
	if err := applyRouterUpdate(qc.String(d.Get(resourceVpcStaticVpcid).(string)), meta); err != nil {
		return nil
	}
	return resourceQingcloudVpcStaticRead(d, meta)
}

func resourceQingcloudVpcStaticRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.DescribeRouterStaticsInput)
	input.Router = getSetStringPointer(d, resourceVpcStaticVpcid)
	input.RouterStatics = []*string{qc.String(d.Id())}
	var output *qc.DescribeRouterStaticsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeRouterStatics(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if len(output.RouterStaticSet) == 0 {
		d.SetId("")
		return nil
	}
	d.Set(resourceName, qc.StringValue(output.RouterStaticSet[0].RouterStaticName))
	d.Set(resourceVpcStaticType, qc.IntValue(output.RouterStaticSet[0].StaticType))
	d.Set(resourceVpcVal1, qc.StringValue(output.RouterStaticSet[0].Val1))
	d.Set(resourceVpcVal2, qc.StringValue(output.RouterStaticSet[0].Val2))
	d.Set(resourceVpcVal3, qc.StringValue(output.RouterStaticSet[0].Val3))
	d.Set(resourceVpcVal4, qc.StringValue(output.RouterStaticSet[0].Val4))
	d.Set(resourceVpcVal5, qc.StringValue(output.RouterStaticSet[0].Val5))
	return nil
}

func resourceQingcloudVpcStaticUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.ModifyRouterStaticAttributesInput)
	input.RouterStatic = qc.String(d.Id())
	input.RouterStaticName, _ = getNamePointer(d)
	input.Val1 = getUpdateStringPointer(d, resourceVpcVal1)
	input.Val2 = getUpdateStringPointer(d, resourceVpcVal2)
	input.Val3 = getUpdateStringPointer(d, resourceVpcVal3)
	input.Val4 = getUpdateStringPointer(d, resourceVpcVal4)
	input.Val5 = getUpdateStringPointer(d, resourceVpcVal5)
	var err error
	simpleRetry(func() error {
		_, err = clt.ModifyRouterStaticAttributes(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if err := applyRouterUpdate(qc.String(d.Get(resourceVpcStaticVpcid).(string)), meta); err != nil {
		return nil
	}
	return resourceQingcloudVpcStaticRead(d, meta)
}

func resourceQingcloudVpcStaticDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.DeleteRouterStaticsInput)
	input.RouterStatics = []*string{qc.String(d.Id())}
	var err error
	simpleRetry(func() error {
		_, err = clt.DeleteRouterStatics(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if err := applyRouterUpdate(qc.String(d.Get(resourceVpcStaticVpcid).(string)), meta); err != nil {
		return nil
	}
	d.SetId("")
	return nil
}
