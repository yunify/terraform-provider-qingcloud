package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func resourceQingcloudEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudEipCreate,
		Read:   resourceQingcloudEipRead,
		Update: resourceQingcloudEipUpdate,
		Delete: resourceQingcloudEipDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "公网 IP 的名称",
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"bandwidth": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "公网IP带宽上限，单位为Mbps",
			},
			"billing_mode": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "traffic",
				Description:  "公网IP计费模式：bandwidth 按带宽计费，traffic 按流量计费，默认是 bandwidth",
				ValidateFunc: withinArrayString("traffic", "bandwidth"),
			},
			"need_icp": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				Description:  "是否需要备案，1为需要，0为不需要，默认是0",
				ValidateFunc: withinArrayInt(0, 1),
			},
			// -------------------------------------------
			// ----------    如下是自动计算的     -----------
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
		},
	}
}

func resourceQingcloudEipCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip

	input := new(qc.AllocateEIPsInput)
	input.Bandwidth = qc.Int(d.Get("bandwidth").(int))
	input.BillingMode = qc.String(d.Get("billing_mode").(string))
	input.EIPName = qc.String(d.Get("name").(string))
	input.NeedICP = qc.Int(d.Get("need_icp").(int))
	input.Count = qc.Int(1)
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("Error create eip input validate: %s", err)
	}
	output, err := clt.AllocateEIPs(input)
	if err != nil {
		return fmt.Errorf("Error create eip: %s", err)
	}
	if *output.RetCode != 0 {
		return fmt.Errorf("Error create eip: %s", *output.Message)
	}
	d.SetId(qc.StringValue(output.EIPs[0]))
	// set eip description
	if err := modifyEipAttributes(d, meta, true); err != nil {
		return err
	}
	return resourceQingcloudEipRead(d, meta)
}

func resourceQingcloudEipRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip

	input := new(qc.DescribeEIPsInput)
	input.EIPs = []*string{qc.String(d.Id())}
	input.Verbose = qc.Int(1)
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("Error describe eip input validate: %s", err)
	}
	output, err := clt.DescribeEIPs(input)
	if err != nil {
		return fmt.Errorf("Error describe eip: %s", err)
	}
	if *output.RetCode != 0 {
		return fmt.Errorf("Error describe eip: %s", *output.Message)
	}
	ip := output.EIPSet[0]
	d.Set("name", ip.EIPName)
	d.Set("billing_mode", ip.BillingMode)
	d.Set("bandwidth", ip.Bandwidth)
	d.Set("need_icp", ip.NeedICP)
	d.Set("description", ip.Description)
	// 如下状态是稍等来获取的
	d.Set("addr", ip.EIPAddr)
	d.Set("status", ip.Status)
	d.Set("transition_status", ip.TransitionStatus)
	if err := d.Set("resource", getEIPResourceMap(ip)); err != nil {
		return fmt.Errorf("Error set eip resource %v", err)
	}
	return nil
}

func resourceQingcloudEipUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip

	if !d.HasChange("name") && !d.HasChange("description") && !d.HasChange("bandwidth") && !d.HasChange("billing_mode") {
		return nil
	}
	if d.HasChange("bandwidth") {
		input := new(qc.ChangeEIPsBandwidthInput)
		input.EIPs = []*string{qc.String(d.Id())}
		input.Bandwidth = qc.Int(d.Get("bandwidth").(int))
		err := input.Validate()
		if err != nil {
			return fmt.Errorf("Error Change EIP bandwidth input validate: %s", err)
		}
		output, err := clt.ChangeEIPsBandwidth(input)
		if err != nil {
			return fmt.Errorf("Errorf Change EIP bandwidth input: %s", err)
		}
		if *output.RetCode != 0 {
			return fmt.Errorf("Errorf Change EIP bandwidth input: %s", err)
		}
	}
	if d.HasChange("billing_mode") {
		input := new(qc.ChangeEIPsBillingModeInput)
		input.EIPs = []*string{qc.String(d.Id())}
		input.BillingMode = qc.String(d.Get("billing_mode").(string))
		err := input.Validate()
		if err != nil {
			return fmt.Errorf("Error Change EIPs billing_mode input validate: %s", err)
		}
		output, err := clt.ChangeEIPsBillingMode(input)
		if err != nil {
			return fmt.Errorf("Errorf Change EIPs billing_mode %s", err)
		}
		if *output.RetCode != 0 {
			return fmt.Errorf("Errorf Change EIP billing_mode %s", *output.Message)
		}
	}
	if err := modifyEipAttributes(d, meta, false); err != nil {
		return err
	}
	return resourceQingcloudEipRead(d, meta)
}

func resourceQingcloudEipDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip

	input := new(qc.ReleaseEIPsInput)
	input.EIPs = []*string{qc.String(d.Id())}
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("Error release eip input validate: %s", err)
	}
	output, err := clt.ReleaseEIPs(input)
	if err != nil {
		return fmt.Errorf("Error release eip: %s", err)
	}
	if *output.RetCode != 0 {
		return fmt.Errorf("Error describe eip: %s", *output.Message)
	}
	d.SetId("")
	return nil
}
