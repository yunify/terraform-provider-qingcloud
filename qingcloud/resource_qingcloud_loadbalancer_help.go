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
		time.Duration(waitJobTimeOutDefault)*time.Second, time.Duration(waitJobIntervalDefault)*time.Second)
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
		time.Duration(10)*time.Minute, time.Duration(5)*time.Second)
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
		time.Duration(waitJobTimeOutDefault)*time.Minute, time.Duration(waitJobIntervalDefault)*time.Second)
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
	httpHeaderSizeUpdate := false
	input.LoadBalancerName, nameUpdate = getNamePointer(d)
	input.Description, descriptionUpdate = getDescriptionPointer(d)
	input.SecurityGroup, sgUpdate = getUpdateStringPointerInfo(d, resourceLoadBalancerSecurityGroupID)
	if d.Get(resourceLoadBalancerVxnetID).(string) == BasicNetworkID {
		input.NodeCount, ncUpdate = getUpdateIntPointerInfo(d, resourceLoadBalancerNodeCount)
	}
	input.HTTPHeaderSize, httpHeaderSizeUpdate = getUpdateIntPointerInfo(d, resourceLoadBalancerHttpHeaderSize)
	if d.HasChange(resourceLoadBalancerPrivateIPs) && !d.IsNewResource() {
		privateIPs := d.Get(resourceLoadBalancerPrivateIPs).(*schema.Set).List()
		if len(privateIPs) != 1 || d.Get(resourceLoadBalancerVxnetID).(string) == BasicNetworkID {
			return fmt.Errorf("error private_ips info")
		}
		input.PrivateIP = qc.String(privateIPs[0].(string))
		privateIPUpdate = true
	}
	if nameUpdate || descriptionUpdate || sgUpdate || ncUpdate || privateIPUpdate || httpHeaderSizeUpdate {
		var err error
		simpleRetry(func() error {
			_, err = clt.ModifyLoadBalancerAttributes(input)
			return isServerBusy(err)
		})
		if err != nil {
			return err
		}
	}
	if sgUpdate || ncUpdate || privateIPUpdate || httpHeaderSizeUpdate {
		updateLoadBalancer(qc.String(d.Id()), meta)
	}
	return nil
}

func updateLoadbalancerEips(d *schema.ResourceData, meta interface{}) error {
	if d.IsNewResource() && d.Get(resourceLoadBalancerVxnetID) == BasicNetworkID {
		return nil
	}
	oldV, newV := d.GetChange(resourceLoadBalancerEipIDs)
	var newEips []string
	var oldEips []string
	if oldV != nil {
		for _, v := range oldV.(*schema.Set).List() {
			oldEips = append(oldEips, v.(string))
		}
	}
	if newV != nil {
		for _, v := range newV.(*schema.Set).List() {
			newEips = append(newEips, v.(string))
		}
	}
	additions, deletions := stringSliceDiff(newEips, oldEips)
	if len(deletions) > 0 {
		if err := dissociateEipsToLoadBalancer(qc.String(d.Id()), qc.StringSlice(deletions), meta); err != nil {
			return err
		}
	}
	if len(additions) > 0 {
		if err := associateEipsToLoadBalancer(qc.String(d.Id()), qc.StringSlice(additions), meta); err != nil {
			return err
		}
	}
	return nil
}

func isLoadBalancerDeleted(lbSet []*qc.LoadBalancer) bool {
	if len(lbSet) == 0 || qc.StringValue(lbSet[0].Status) == "deleted" || qc.StringValue(lbSet[0].Status) == "ceased" {
		return true
	}
	return false
}

func waitLoadBalancerLease(d *schema.ResourceData, meta interface{}) error {
	if !d.IsNewResource() {
		return nil
	}
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.DescribeLoadBalancersInput)
	input.LoadBalancers = []*string{qc.String(d.Id())}
	input.Verbose = qc.Int(1)
	var output *qc.DescribeLoadBalancersOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeLoadBalancers(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	//wait for lease info
	WaitForLease(output.LoadBalancerSet[0].StatusTime)
	return nil
}
