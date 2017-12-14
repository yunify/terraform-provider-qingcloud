package qingcloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yunify/qingcloud-sdk-go/client"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func updateLoadBalancer(lbID *string, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.UpdateLoadBalancersInput)
	input.LoadBalancers = []*string{lbID}
	var output *qc.UpdateLoadBalancersOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.UpdateLoadBalancers(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	client.WaitJob(meta.(*QingCloudClient).job,
		qc.StringValue(output.JobID),
		time.Duration(10)*time.Second, time.Duration(1)*time.Second)
	if _, err := LoadBalancerTransitionStateRefresh(clt, lbID); err != nil {
		return err
	}
	return nil
}

func resizeLoadBalancer(lbID *string, lbType *int, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.ResizeLoadBalancersInput)
	input.LoadBalancers = []*string{lbID}
	input.LoadBalancerType = lbType
	var output *qc.ResizeLoadBalancersOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.ResizeLoadBalancers(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	client.WaitJob(meta.(*QingCloudClient).job,
		qc.StringValue(output.JobID),
		time.Duration(10)*time.Second, time.Duration(1)*time.Second)
	if _, err := LoadBalancerTransitionStateRefresh(clt, lbID); err != nil {
		return err
	}
	return nil
}

func associateEipsToLoadBalancer(lbID *string, eips []*string, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.AssociateEIPsToLoadBalancerInput)
	input.LoadBalancer = lbID
	input.EIPs = eips
	var output *qc.AssociateEIPsToLoadBalancerOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.AssociateEIPsToLoadBalancer(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	client.WaitJob(meta.(*QingCloudClient).job,
		qc.StringValue(output.JobID),
		time.Duration(10)*time.Second, time.Duration(1)*time.Second)
	if _, err := LoadBalancerTransitionStateRefresh(clt, lbID); err != nil {
		return err
	}
	return nil
}

func dissociateEipsToLoadBalancer(lbID *string, eips []*string, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.DissociateEIPsFromLoadBalancerInput)
	input.LoadBalancer = lbID
	input.EIPs = eips
	var output *qc.DissociateEIPsFromLoadBalancerOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DissociateEIPsFromLoadBalancer(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	client.WaitJob(meta.(*QingCloudClient).job,
		qc.StringValue(output.JobID),
		time.Duration(10)*time.Second, time.Duration(1)*time.Second)
	if _, err := LoadBalancerTransitionStateRefresh(clt, lbID); err != nil {
		return err
	}
	return nil
}

func modifyLoadBalancerAttributes(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.ModifyLoadBalancerAttributesInput)
	input.LoadBalancer = qc.String(d.Id())
	nameUpdate := false
	descriptionUpdate := false
	sgUpdate := false
	ncUpdate := false
	privateIPUpdate := false
	input.LoadBalancerName, nameUpdate = getNamePointer(d)
	input.Description, descriptionUpdate = getDescriptionPointer(d)
	input.SecurityGroup, sgUpdate = getUpdateStringPointerInfo(d, resourceLoadBalancerSecurityGroupID)
	input.NodeCount, ncUpdate = getUpdateIntPointerInfo(d, resourceLoadBalancerNodeCount)
	if d.HasChange(resourceLoadBalancerPrivateIPs) {
		privateIPs := d.Get(resourceLoadBalancerPrivateIPs).(*schema.Set).List()
		if len(privateIPs) != 1 || d.Get(resourceLoadBalancerVxnetID).(string) == "vxnet-0" {
			return fmt.Errorf("error private_ips info")
		}
		input.PrivateIP = qc.String(privateIPs[0].(string))
		privateIPUpdate = true
	}
	if nameUpdate || descriptionUpdate || sgUpdate || ncUpdate || privateIPUpdate {
		var err error
		simpleRetry(func() error {
			_, err = clt.ModifyLoadBalancerAttributes(input)
			return isServerBusy(err)
		})
		if err != nil {
			return err
		}
	}
	if sgUpdate || ncUpdate || privateIPUpdate {
		updateLoadBalancer(qc.String(d.Id()), meta)
	}
	return nil
}
