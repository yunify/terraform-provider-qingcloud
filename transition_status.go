package qingcloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/magicshui/qingcloud-go/eip"
	"github.com/magicshui/qingcloud-go/instance"
	"github.com/magicshui/qingcloud-go/loadbalancer"
	"github.com/magicshui/qingcloud-go/router"
	"github.com/magicshui/qingcloud-go/volume"
	"time"
)

func LoadbalancerTransitionStateRefresh(clt *loadbalancer.LOADBALANCER, id string) (interface{}, error) {
	refreshFunc := func() (interface{}, string, error) {
		params := loadbalancer.DescribeLoadBalancersRequest{}
		params.LoadbalancersN.Add(id)
		params.Verbose.Set(1)

		resp, err := clt.DescribeLoadBalancers(params)
		if err != nil {
			return nil, "", err
		}
		return resp.LoadbalancerSet[0], resp.LoadbalancerSet[0].TransitionStatus, nil
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "starting", "stopping", "updating", "suspending", "resuming", "deleting"},
		Target:     []string{""},
		Refresh:    refreshFunc,
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}

// Waiting for no transition_status
func EipTransitionStateRefresh(clt *eip.EIP, id string) (interface{}, error) {
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
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}
func VolumeTransitionStateRefresh(clt *volume.VOLUME, id string) (interface{}, error) {
	refreshFunc := func() (interface{}, string, error) {
		params := volume.DescribeVolumesRequest{}
		params.VolumesN.Add(id)
		params.Verbose.Set(1)

		resp, err := clt.DescribeVolumes(params)
		if err != nil {
			return nil, "", err
		}
		return resp.VolumeSet[0], resp.VolumeSet[0].TransitionStatus, nil
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "attaching", "detaching", "suspending", "suspending", "resuming", "deleting", "recovering"}, // creating, attaching, detaching, suspending，resuming，deleting，recovering
		Target:     []string{""},
		Refresh:    refreshFunc,
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()

}

func RouterTransitionStateRefresh(clt *router.ROUTER, id string) (interface{}, error) {
	refreshFunc := func() (interface{}, string, error) {
		params := router.DescribeRoutersRequest{}
		params.RoutersN.Add(id)
		params.Verbose.Set(1)
		resp, err := clt.DescribeRouters(params)
		if err != nil {
			return nil, "", err
		}
		return resp.RouterSet[0], resp.RouterSet[0].TransitionStatus, nil
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "updating", "suspending", "resuming", "poweroffing", "poweroning", "deleting"},
		Target:     []string{""},
		Refresh:    refreshFunc,
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}

func InstanceTransitionStateRefresh(clt *instance.INSTANCE, id string) (interface{}, error) {
	refreshFunc := func() (interface{}, string, error) {
		params := instance.DescribeInstanceRequest{}
		params.InstancesN.Add(id)
		params.Verbose.Set(1)
		resp, err := clt.DescribeInstances(params)
		if err != nil {
			return nil, "", err
		}
		return resp.InstanceSet[0], resp.InstanceSet[0].TransitionStatus, nil
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "updating", "suspending", "resuming", "poweroffing", "poweroning", "deleting"},
		Target:     []string{""},
		Refresh:    refreshFunc,
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}
