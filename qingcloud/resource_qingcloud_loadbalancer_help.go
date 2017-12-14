package qingcloud

import (
	"github.com/yunify/qingcloud-sdk-go/client"
	qc "github.com/yunify/qingcloud-sdk-go/service"
	"time"
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
