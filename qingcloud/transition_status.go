package qingcloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/magicshui/qingcloud-go/eip"
	"github.com/magicshui/qingcloud-go/instance"
	"github.com/magicshui/qingcloud-go/loadbalancer"
	"github.com/magicshui/qingcloud-go/volume"
	qc "github.com/yunify/qingcloud-sdk-go/service"
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
func EipTransitionStateRefresh(clt *qc.EIPService, id string) (interface{}, error) {
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
func VolumeTransitionStateRefresh(clt *qc.VolumeService, id string) (interface{}, error) {
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

func RouterTransitionStateRefresh(clt *qc.RouterService, id string) (interface{}, error) {
	refreshFunc := func() (interface{}, string, error) {
		input := new(qc.DescribeRoutersInput)
		input.Routers = []*string{qc.String(id)}
		input.Verbose = qc.Int(1)
		err := input.Validate()
		if err != nil {
			return nil, "", fmt.Errorf("Error describe router validate input: %s", err)
		}
		output, err := clt.DescribeRouters(input)
		if err != nil {
			return nil, "", fmt.Errorf("Errorf describe router: %s", err)
		}
		return output.RouterSet[0], qc.StringValue(resp.RouterSet[0].TransitionStatus), nil
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

func InstanceTransitionStateRefresh(clt *qc.InstanceService, id string) (interface{}, error) {
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

func SecurityGroupTransitionStateRefresh(clt *qc.SecurityGroupService, id string) (interface{}, error) {
	refreshFunc := func() (interface{}, string, error) {
		input := new(qc.DescribeSecurityGroupsInput)
		input.SecurityGroups = []*string{qc.String(id)}
		err := input.Validate()
		if err != nil {
			return fmt.Errorf("Error describe securitygroup input validate: %s", err.Error())
		}
		output, err := clt.DescribeSecurityGroups(input)
		if err != nil {
			return fmt.Errorf("Error describe securitygroup input validate: %s", err.Error())
		}
		if output.RetCode != 0 {
			return fmt.Errorf("Error describe securitygroup input validate: %s", err.Error())
		}
		sg := output.SecurityGroupSet[0]
		if *sg.IsApplied == 1 {
			return nil, "updated", nil
		}
		return nil, "not_updated", nil
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"not_updated"},
		Target:     []string{"updated"},
		Refresh:    refreshFunc,
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}
