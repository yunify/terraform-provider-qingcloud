package qingcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

// EipTransitionStateRefresh Waiting for no transition_status
func EIPTransitionStateRefresh(clt *qc.EIPService, id string) (interface{}, error) {
	refreshFunc := func() (interface{}, string, error) {
		input := new(qc.DescribeEIPsInput)
		input.EIPs = []*string{qc.String(id)}
		var output *qc.DescribeEIPsOutput
		var err error
		simpleRetry(func() error {
			output, err = clt.DescribeEIPs(input)
			return isServerBusy(err)
		})
		if err != nil {
			return nil, "", err
		}
		if len(output.EIPSet) == 0 {
			return nil, "", fmt.Errorf("Error eip set is empty, request id %s", id)
		}
		return output.EIPSet[0], qc.StringValue(output.EIPSet[0].TransitionStatus), nil
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"associating", "dissociating", "suspending", "resuming", "releasing"},
		Target:     []string{""},
		Refresh:    refreshFunc,
		Timeout:    waitJobTimeOutDefault * time.Second,
		Delay:      waitJobIntervalDefault * time.Second,
		MinTimeout: waitJobIntervalDefault * time.Second,
	}
	return stateConf.WaitForState()
}

func VolumeTransitionStateRefresh(clt *qc.VolumeService, id string) (interface{}, error) {
	refreshFunc := func() (interface{}, string, error) {
		input := new(qc.DescribeVolumesInput)
		input.Volumes = []*string{qc.String(id)}
		var output *qc.DescribeVolumesOutput
		var err error
		simpleRetry(func() error {
			output, err = clt.DescribeVolumes(input)
			return isServerBusy(err)
		})
		if err != nil {
			return nil, "", err
		}
		if len(output.VolumeSet) != 1 {
			return output, "creating", nil
		}
		volume := output.VolumeSet[0]
		return volume, qc.StringValue(volume.TransitionStatus), nil
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "attaching", "detaching", "suspending", "suspending", "resuming", "deleting", "recovering"}, // creating, attaching, detaching, suspending，resuming，deleting，recovering
		Target:     []string{""},
		Refresh:    refreshFunc,
		Timeout:    waitJobTimeOutDefault * time.Second,
		Delay:      waitJobIntervalDefault * time.Second,
		MinTimeout: waitJobIntervalDefault * time.Second,
	}
	return stateConf.WaitForState()
}

func VolumeDeleteTransitionStateRefresh(clt *qc.VolumeService, id string) (interface{}, error) {
	refreshFunc := func() (interface{}, string, error) {
		input := new(qc.DescribeVolumesInput)
		input.Volumes = []*string{qc.String(id)}
		var output *qc.DescribeVolumesOutput
		var err error
		simpleRetry(func() error {
			output, err = clt.DescribeVolumes(input)
			return isServerBusy(err)
		})
		if err != nil {
			return nil, "", err
		}
		volume := output.VolumeSet[0]
		return volume, qc.StringValue(volume.Status), nil
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "in-use", "suspended", "deleted", "ceased"},
		Target:     []string{"available"},
		Refresh:    refreshFunc,
		Timeout:    waitJobTimeOutDefault * time.Second,
		Delay:      waitJobIntervalDefault * time.Second,
		MinTimeout: waitJobIntervalDefault * time.Second,
	}
	return stateConf.WaitForState()
}

// RouterTransitionStateRefresh Waiting for no transition_status
func RouterTransitionStateRefresh(clt *qc.RouterService, id string) (interface{}, error) {
	if id == "" {
		return nil, nil
	}
	refreshFunc := func() (interface{}, string, error) {
		input := new(qc.DescribeRoutersInput)
		input.Routers = []*string{qc.String(id)}
		input.Verbose = qc.Int(1)
		var output *qc.DescribeRoutersOutput
		var err error
		simpleRetry(func() error {
			output, err = clt.DescribeRouters(input)
			return isServerBusy(err)
		})
		if err != nil {
			return nil, "", fmt.Errorf("Errorf describe router: %s", err)
		}
		if len(output.RouterSet) == 0 {
			return nil, "", fmt.Errorf("Error router set is empty, request id %s", id)
		}
		return output.RouterSet[0], qc.StringValue(output.RouterSet[0].TransitionStatus), nil
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "updating", "suspending", "resuming", "poweroffing", "poweroning", "deleting"},
		Target:     []string{""},
		Refresh:    refreshFunc,
		Timeout:    waitJobTimeOutDefault * time.Second,
		Delay:      waitJobIntervalDefault * time.Second,
		MinTimeout: waitJobIntervalDefault * time.Second,
	}
	return stateConf.WaitForState()
}

func InstanceTransitionStateRefresh(clt *qc.InstanceService, id string) (interface{}, error) {
	if id == "" {
		return nil, nil
	}
	refreshFunc := func() (interface{}, string, error) {
		input := new(qc.DescribeInstancesInput)
		input.Instances = []*string{qc.String(id)}
		var output *qc.DescribeInstancesOutput
		var err error
		simpleRetry(func() error {
			output, err = clt.DescribeInstances(input)
			return isServerBusy(err)
		})
		if err != nil {
			return nil, "", err
		}
		if len(output.InstanceSet) == 0 {
			return nil, "", fmt.Errorf("Error instance set is empty, request id %s", id)
		}
		if isInstanceDeleted(output.InstanceSet) {
			return output.InstanceSet[0], "", nil
		}
		if len(output.InstanceSet[0].VxNets) != 0 {
			// return output.InstanceSet[0], "creating", nil
			if qc.StringValue(output.InstanceSet[0].VxNets[0].PrivateIP) == "" {
				return output.InstanceSet[0], "creating", nil
			}
		}
		return output.InstanceSet[0], qc.StringValue(output.InstanceSet[0].TransitionStatus), nil
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "updating", "suspending", "resuming", "poweroffing", "poweroning", "deleting", "stopping", "starting", "terminating"},
		Target:     []string{""},
		Refresh:    refreshFunc,
		Timeout:    waitJobTimeOutDefault * time.Second,
		Delay:      waitJobIntervalDefault * time.Second,
		MinTimeout: waitJobIntervalDefault * time.Second,
	}
	return stateConf.WaitForState()
}

func VxnetLeaveRouterTransitionStateRefresh(clt *qc.VxNetService, id string) (interface{}, error) {
	if id == "" {
		return nil, nil
	}
	refreshFunc := func() (interface{}, string, error) {
		input := new(qc.DescribeVxNetsInput)
		input.VxNets = []*string{qc.String(id)}
		var output *qc.DescribeVxNetsOutput
		var err error
		simpleRetry(func() error {
			output, err = clt.DescribeVxNets(input)
			return isServerBusy(err)
		})
		if err != nil {
			return nil, "", err
		}
		if len(output.VxNetSet) == 0 {
			return nil, "", nil
		}
		vxnet := output.VxNetSet[0]
		log.Printf("VxnetLeaveRouterTransitionStateRefresh vpc id: %s", *vxnet.VpcRouterID)
		if qc.StringValue(vxnet.VpcRouterID) != "" {
			return vxnet, "vxnet_not_leave_router", nil
		}
		log.Printf("skip if VxnetLeaveRouterTransitionStateRefresh vpc id: %s", *vxnet.VpcRouterID)
		return vxnet, "", nil
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"vxnet_not_leave_router"},
		Target:     []string{""},
		Refresh:    refreshFunc,
		Timeout:    waitJobTimeOutDefault * time.Second,
		Delay:      waitJobIntervalDefault * time.Second,
		MinTimeout: waitJobIntervalDefault * time.Second,
	}
	return stateConf.WaitForState()
}

func SecurityGroupApplyTransitionStateRefresh(clt *qc.SecurityGroupService, id *string) (interface{}, error) {
	refreshFunc := func() (interface{}, string, error) {
		input := new(qc.DescribeSecurityGroupsInput)
		input.SecurityGroups = []*string{id}
		var output *qc.DescribeSecurityGroupsOutput
		var err error
		simpleRetry(func() error {
			output, err = clt.DescribeSecurityGroups(input)
			return isServerBusy(err)
		})
		if err != nil {
			return nil, "not_updated", err
		}
		sg := output.SecurityGroupSet[0]
		if sg.IsApplied != nil && qc.IntValue(sg.IsApplied) == 1 {
			return sg, "updated", nil
		}
		return sg, "not_updated", nil
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"not_updated"},
		Target:     []string{"updated"},
		Refresh:    refreshFunc,
		Timeout:    waitJobTimeOutDefault * time.Second,
		Delay:      waitJobIntervalDefault * time.Second,
		MinTimeout: waitJobIntervalDefault * time.Second,
	}
	return stateConf.WaitForState()
}

func WaitForLease(CreateTime *time.Time) {
	now := time.Now()
	subS := now.Sub(qc.TimeValue(CreateTime)).Seconds()
	if subS < float64(waitLeaseSecond) {
		time.Sleep(time.Second * time.Duration(waitLeaseSecond))
	}
}
func LoadBalancerTransitionStateRefresh(clt *qc.LoadBalancerService, id *string) (interface{}, error) {
	refreshFunc := func() (interface{}, string, error) {
		input := new(qc.DescribeLoadBalancersInput)
		input.LoadBalancers = []*string{id}
		var output *qc.DescribeLoadBalancersOutput
		var err error
		simpleRetry(func() error {
			output, err = clt.DescribeLoadBalancers(input)
			return isServerBusy(err)
		})
		if err != nil {
			return nil, "", err
		}
		if len(output.LoadBalancerSet) == 0 {
			return nil, "", fmt.Errorf("error lb set is empty, request id %s", *id)
		}
		return output.LoadBalancerSet[0], qc.StringValue(output.LoadBalancerSet[0].TransitionStatus), nil
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "starting", "stopping", "updating", "suspending", "resuming", "deleting"},
		Target:     []string{""},
		Refresh:    refreshFunc,
		Timeout:    waitJobTimeOutDefault * time.Second,
		Delay:      waitJobIntervalDefault * time.Second,
		MinTimeout: waitJobIntervalDefault * time.Second,
	}
	return stateConf.WaitForState()
}
