package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

const (
	resourceLoadBalancerBackendResrourceId = "resource_id"
	resourceLoadBalancerBackendPort        = "port"
	resourceLoadBalancerBackendWeight      = "weight"
	resourceLoadBalancerBackendListenerId  = "loadbalancer_listener_id"
)

func resourceQingcloudLoadBalancerBackend() *schema.Resource {

	return &schema.Resource{
		Create: resourceQingcloudLoadBalancerBackendCreate,
		Read:   resourceQingcloudLoadBalancerBackendRead,
		Update: resourceQingcloudLoadBalancerBackendUpdate,
		Delete: resourceQingcloudLoadBalancerBackendDelete,
		Schema: map[string]*schema.Schema{
			resourceName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceLoadBalancerBackendResrourceId: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			resourceLoadBalancerBackendPort: {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: withinArrayIntRange(1, 65535),
			},
			resourceLoadBalancerBackendListenerId: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			resourceLoadBalancerBackendWeight: {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: withinArrayIntRange(1, 100),
			},
		},
	}
}
func resourceQingcloudLoadBalancerBackendCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.AddLoadBalancerBackendsInput)
	lbe := new(qc.LoadBalancerBackend)
	lbe.LoadBalancerListenerID = getSetStringPointer(d, resourceLoadBalancerBackendListenerId)
	lbe.Port = qc.Int(d.Get(resourceLoadBalancerBackendPort).(int))
	lbe.LoadBalancerBackendName = getSetStringPointer(d, resourceName)
	lbe.ResourceID = getSetStringPointer(d, resourceLoadBalancerBackendResrourceId)
	lbe.Weight = qc.Int(d.Get(resourceLoadBalancerBackendWeight).(int))
	input.Backends = []*qc.LoadBalancerBackend{lbe}
	input.LoadBalancerListener = getSetStringPointer(d, resourceLoadBalancerBackendListenerId)
	var output *qc.AddLoadBalancerBackendsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.AddLoadBalancerBackends(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	lbId, err := getLBIdFromLBB(output.LoadBalancerBackends[0], meta)
	if err != nil {
		return err
	}
	if err := updateLoadBalancer(lbId, meta); err != nil {
		return nil
	}
	d.SetId(qc.StringValue(output.LoadBalancerBackends[0]))
	return resourceQingcloudLoadBalancerBackendRead(d, meta)
}
func resourceQingcloudLoadBalancerBackendDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	var err error
	lbId, err := getLBIdFromLBB(qc.String(d.Id()), meta)
	if err != nil {
		return err
	}
	input := new(qc.DeleteLoadBalancerBackendsInput)
	input.LoadBalancerBackends = []*string{qc.String(d.Id())}
	simpleRetry(func() error {
		_, err = clt.DeleteLoadBalancerBackends(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if err := updateLoadBalancer(lbId, meta); err != nil {
		return nil
	}
	d.SetId("")
	return nil

}
func resourceQingcloudLoadBalancerBackendRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.DescribeLoadBalancerBackendsInput)
	input.LoadBalancerBackends = []*string{qc.String(d.Id())}
	input.LoadBalancerListener = getSetStringPointer(d, resourceLoadBalancerBackendListenerId)
	var output *qc.DescribeLoadBalancerBackendsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeLoadBalancerBackends(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if len(output.LoadBalancerBackendSet) == 0 {
		d.SetId("")
		return nil
	}
	d.Set(resourceName, qc.StringValue(output.LoadBalancerBackendSet[0].LoadBalancerBackendName))
	d.Set(resourceLoadBalancerBackendResrourceId, qc.StringValue(output.LoadBalancerBackendSet[0].ResourceID))
	d.Set(resourceLoadBalancerBackendListenerId, qc.StringValue(output.LoadBalancerBackendSet[0].LoadBalancerListenerID))
	d.Set(resourceLoadBalancerBackendPort, qc.IntValue(output.LoadBalancerBackendSet[0].Port))
	d.Set(resourceLoadBalancerBackendWeight, qc.IntValue(output.LoadBalancerBackendSet[0].Weight))
	return nil
}

func resourceQingcloudLoadBalancerBackendUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.ModifyLoadBalancerBackendAttributesInput)
	input.LoadBalancerBackend = qc.String(d.Id())
	input.Weight = qc.Int(d.Get(resourceLoadBalancerBackendWeight).(int))
	input.Port = qc.Int(d.Get(resourceLoadBalancerBackendPort).(int))
	input.LoadBalancerBackendName = getUpdateStringPointer(d, resourceName)

	var err error
	simpleRetry(func() error {
		_, err = clt.ModifyLoadBalancerBackendAttributes(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if err := updateLoadBalancer(qc.String(d.Get(resourceLoadBalancerBackendListenerId).(string)), meta); err != nil {
		return nil
	}
	return resourceQingcloudLoadBalancerListenerRead(d, meta)

}

func getLBIdFromLBB(lbbId *string, meta interface{}) (*string, error) {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.DescribeLoadBalancerBackendsInput)
	input.Verbose = qc.Int(1)
	input.LoadBalancerBackends = []*string{lbbId}
	var output *qc.DescribeLoadBalancerBackendsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeLoadBalancerBackends(input)
		return isServerBusy(err)
	})
	if err != nil {
		return nil, err
	}
	return output.LoadBalancerBackendSet[0].LoadBalancerID, nil
}
