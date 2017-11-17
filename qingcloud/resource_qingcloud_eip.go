package qingcloud

import (
	"fmt"

	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yunify/qingcloud-sdk-go/client"
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
				Optional:    true,
				Description: "the name of eip",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the description of eip",
			},
			"bandwidth": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Maximum bandwidth to the elastic public network, measured in Mbps",
			},
			"billing_mode": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "bandwidth",
				Description:  "Internet charge type of the EIP : bandwidth , traffic ,default bandwidth",
				ValidateFunc: withinArrayString("traffic", "bandwidth"),
			},
			"need_icp": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				Description:  "need icp , 1 need , 0 no need ,default 0",
				ValidateFunc: withinArrayInt(0, 1),
			},
			"tag_ids": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "tag ids , eip wants to use",
			},
			"tag_names": &schema.Schema{
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "compute by tag ids",
			},
			"addr": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ip address of this eip",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the status of eip",
			},
			"resource": &schema.Schema{
				Type:         schema.TypeMap,
				Computed:     true,
				ComputedWhen: []string{"id"},
				Description:  "the resource who use this eip",
			},
		},
	}
}

func resourceQingcloudEipCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip
	input := new(qc.AllocateEIPsInput)
	input.Bandwidth = qc.Int(d.Get("bandwidth").(int))
	input.BillingMode = qc.String(d.Get("billing_mode").(string))
	input.NeedICP = qc.Int(d.Get("need_icp").(int))
	input.Count = qc.Int(1)
	if d.Get("name") != "" {
		input.EIPName = qc.String(d.Get("name").(string))
	} else {
		input.EIPName = nil
	}
	var output *qc.AllocateEIPsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.AllocateEIPs(input)
		return IsIaasAPIServerBusy(output.RetCode, err)
	})
	if err != nil {
		return fmt.Errorf("Error create eip: %s", err)
	}
	if *output.RetCode != 0 {
		return fmt.Errorf("Error create eip: %s", *output.Message)
	}
	d.SetId(qc.StringValue(output.EIPs[0]))

	if _, err := EIPTransitionStateRefresh(clt, d.Id()); err != nil {
		return nil
	}
	return resourceQingcloudEipUpdate(d, meta)
}

func resourceQingcloudEipRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip
	input := new(qc.DescribeEIPsInput)
	input.EIPs = []*string{qc.String(d.Id())}
	input.Verbose = qc.Int(1)
	var output *qc.DescribeEIPsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeEIPs(input)
		return IsIaasAPIServerBusy(output.RetCode, err)

	})
	if err != nil {
		return fmt.Errorf("Error describe eip: %s", err)
	}
	if *output.RetCode != 0 {
		return fmt.Errorf("Error describe eip: %s", *output.Message)
	}
	if len(output.EIPSet) == 0 || qc.StringValue(output.EIPSet[0].Status) == "ceased" || qc.StringValue(output.EIPSet[0].Status) == "released" {
		d.SetId("")
		return nil
	}
	ip := output.EIPSet[0]
	d.Set("name", qc.StringValue(ip.EIPName))
	d.Set("billing_mode", qc.StringValue(ip.BillingMode))
	d.Set("bandwidth", qc.IntValue(ip.Bandwidth))
	d.Set("need_icp", qc.IntValue(ip.NeedICP))
	d.Set("description", qc.StringValue(ip.Description))
	// 如下状态是稍等来获取的
	d.Set("addr", qc.StringValue(ip.EIPAddr))
	d.Set("status", qc.StringValue(ip.Status))
	if err := d.Set("resource", getEIPResourceMap(ip)); err != nil {
		return fmt.Errorf("Error set eip resource %v", err)
	}
	resourceSetTag(d, ip.Tags)
	return nil
}

func resourceQingcloudEipUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip
	d.Partial(true)
	waitEipLease(d, meta)
	if d.HasChange("need_icp") && !d.IsNewResource() {
		return fmt.Errorf("Errorf EIP need_icp could not be updated")
	}
	if d.HasChange("bandwidth") && !d.IsNewResource() {
		input := new(qc.ChangeEIPsBandwidthInput)
		input.EIPs = []*string{qc.String(d.Id())}
		input.Bandwidth = qc.Int(d.Get("bandwidth").(int))
		var output *qc.ChangeEIPsBandwidthOutput
		var err error
		simpleRetry(func() error {
			output, err = clt.ChangeEIPsBandwidth(input)
			return IsIaasAPIServerBusy(output.RetCode, err)
		})
		if err != nil {
			return fmt.Errorf("Errorf Change EIP bandwidth input: %s", err)
		}
		if *output.RetCode != 0 {
			return fmt.Errorf("Errorf Change EIP bandwidth input: %s", err)
		}
		if _, err := EIPTransitionStateRefresh(clt, d.Id()); err != nil {
			return nil
		}
		d.SetPartial("bandwidth")
	}
	if d.HasChange("billing_mode") && !d.IsNewResource() {
		input := new(qc.ChangeEIPsBillingModeInput)
		input.EIPs = []*string{qc.String(d.Id())}
		input.BillingMode = qc.String(d.Get("billing_mode").(string))
		var output *qc.ChangeEIPsBillingModeOutput
		var err error
		simpleRetry(func() error {
			output, err = clt.ChangeEIPsBillingMode(input)
			return IsIaasAPIServerBusy(output.RetCode, err)
		})
		if err != nil {
			return fmt.Errorf("Errorf Change EIPs billing_mode %s", err)
		}
		if *output.RetCode != 0 {
			return fmt.Errorf("Errorf Change EIP billing_mode %s", *output.Message)
		}
		if _, err := EIPTransitionStateRefresh(clt, d.Id()); err != nil {
			return nil
		}
		d.SetPartial("billing_mode")
	}
	if err := modifyEipAttributes(d, meta); err != nil {
		return err
	}
	d.SetPartial("description")
	d.SetPartial("name")
	if err := resourceUpdateTag(d, meta, qingcloudResourceTypeEIP); err != nil {
		return err
	}
	d.SetPartial("tag_ids")
	d.Partial(false)
	return resourceQingcloudEipRead(d, meta)
}

func resourceQingcloudEipDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip
	_, refreshErr := EIPTransitionStateRefresh(clt, d.Id())
	if refreshErr != nil {
		return refreshErr
	}
	waitEipLease(d, meta)

	input := new(qc.ReleaseEIPsInput)
	input.EIPs = []*string{qc.String(d.Id())}
	var output *qc.ReleaseEIPsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.ReleaseEIPs(input)
		return IsIaasAPIServerBusy(output.RetCode, err)
	})
	if err != nil {
		return fmt.Errorf("Error release eip: %s", err)
	}
	if *output.RetCode != 0 {
		return fmt.Errorf("Error describe eip: %s", *output.Message)
	}
	client.WaitJob(meta.(*QingCloudClient).job,
		qc.StringValue(output.JobID),
		time.Duration(10)*time.Second, time.Duration(1)*time.Second)
	d.SetId("")
	return nil
}
