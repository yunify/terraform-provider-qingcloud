package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yunify/qingcloud-sdk-go/client"
	qc "github.com/yunify/qingcloud-sdk-go/service"
	"time"
)

const (
	resourceEipBandwidth = "bandwidth"
	resourceEipBillMode  = "billing_mode"
	resourceEipNeedIcp   = "need_icp"
	resourceEipAddr      = "addr"
	resourceEipResource  = "resource"
)

func resourceQingcloudEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudEipCreate,
		Read:   resourceQingcloudEipRead,
		Update: resourceQingcloudEipUpdate,
		Delete: resourceQingcloudEipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			resourceName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceEipBandwidth: {
				Type:     schema.TypeInt,
				Required: true,
			},
			resourceEipBillMode: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "bandwidth",
				ValidateFunc: withinArrayString("traffic", "bandwidth"),
			},
			resourceEipNeedIcp: {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: withinArrayInt(0, 1),
			},
			resourceTagIds:   tagIdsSchema(),
			resourceTagNames: tagNamesSchema(),
			resourceEipAddr: {
				Type:     schema.TypeString,
				Computed: true,
			},
			resourceEipResource: {
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
	input.Bandwidth = qc.Int(d.Get(resourceEipBandwidth).(int))
	input.BillingMode = qc.String(d.Get(resourceEipBillMode).(string))
	input.NeedICP = qc.Int(d.Get(resourceEipNeedIcp).(int))
	input.Count = qc.Int(1)
	input.EIPName, _ = getNamePointer(d)
	var output *qc.AllocateEIPsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.AllocateEIPs(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
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
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if isEipDeleted(output.EIPSet) {
		d.SetId("")
		return nil
	}
	ip := output.EIPSet[0]
	d.Set(resourceName, qc.StringValue(ip.EIPName))
	d.Set(resourceEipBillMode, qc.StringValue(ip.BillingMode))
	d.Set(resourceEipBandwidth, qc.IntValue(ip.Bandwidth))
	d.Set(resourceEipNeedIcp, qc.IntValue(ip.NeedICP))
	d.Set(resourceDescription, qc.StringValue(ip.Description))
	d.Set(resourceEipAddr, qc.StringValue(ip.EIPAddr))
	if err := d.Set(resourceEipResource, getEIPResourceMap(ip)); err != nil {
		return fmt.Errorf("Error set eip resource %v", err)
	}
	if err := resourceSetTag(d, ip.Tags); err != nil {
		return err
	}
	return nil
}

func resourceQingcloudEipUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	if err := waitEipLease(d, meta); err != nil {
		return err
	}
	if d.HasChange(resourceEipNeedIcp) && !d.IsNewResource() {
		return fmt.Errorf("Errorf EIP need_icp could not be updated")
	}
	if err := changeEIPBandwidth(d, meta); err != nil {
		return err
	}
	d.SetPartial(resourceEipBandwidth)
	if err := changeEIPBillMode(d, meta); err != nil {
		return err
	}
	d.SetPartial(resourceEipBillMode)
	if err := modifyEipAttributes(d, meta); err != nil {
		return err
	}
	d.SetPartial(resourceDescription)
	d.SetPartial(resourceName)
	if err := resourceUpdateTag(d, meta, qingcloudResourceTypeEIP); err != nil {
		return err
	}
	d.SetPartial(resourceTagIds)
	d.Partial(false)
	return resourceQingcloudEipRead(d, meta)
}

func resourceQingcloudEipDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip
	_, refreshErr := EIPTransitionStateRefresh(clt, d.Id())
	if refreshErr != nil {
		return refreshErr
	}
	if err := waitEipLease(d, meta); err != nil {
		return err
	}
	input := new(qc.ReleaseEIPsInput)
	input.EIPs = []*string{qc.String(d.Id())}
	var output *qc.ReleaseEIPsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.ReleaseEIPs(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	client.WaitJob(meta.(*QingCloudClient).job,
		qc.StringValue(output.JobID),
		time.Duration(waitJobTimeOutDefault)*time.Second, time.Duration(waitJobIntervalDefault)*time.Second)
	d.SetId("")
	return nil
}
