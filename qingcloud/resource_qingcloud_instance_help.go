package qingcloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yunify/qingcloud-sdk-go/client"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifyInstanceAttributes(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).instance
	input := new(qc.ModifyInstanceAttributesInput)
	input.Instance = qc.String(d.Id())
	nameUpdate := false
	descriptionUpdate := false
	input.InstanceName, nameUpdate = getNamePointer(d)
	input.Description, descriptionUpdate = getDescriptionPointer(d)
	if nameUpdate || descriptionUpdate {
		var err error
		simpleRetry(func() error {
			_, err = clt.ModifyInstanceAttributes(input)
			return isServerBusy(err)
		})
		if err != nil {
			return err
		}
	}
	return nil
}
func setInstanceUserData(d *schema.ResourceData, meta interface{}, runInstanceInput *qc.RunInstancesInput) error {
	clt := meta.(*QingCloudClient).userdata
	input := new(qc.UploadUserDataAttachmentInput)
	input.AttachmentName = qc.String(resourceInstanceUserData)
	input.AttachmentContent = qc.String(d.Get(resourceInstanceUserData).(string))
	var output *qc.UploadUserDataAttachmentOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.UploadUserDataAttachment(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	runInstanceInput.NeedUserdata = qc.Int(1)
	runInstanceInput.UserdataType = qc.String("tar")
	runInstanceInput.UserdataValue = output.AttachmentID
	runInstanceInput.UserdataFile = qc.String("userdata")
	runInstanceInput.UserdataPath = qc.String("/etc/qingcloud/userdata")
	return nil
}

func instanceUpdateChangeManagedVxNet(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange(resourceInstanceManagedVxnetID) {
		return nil
	}
	clt := meta.(*QingCloudClient).instance
	vxnetClt := meta.(*QingCloudClient).vxnet
	oldV, newV := d.GetChange(resourceInstanceManagedVxnetID)
	// leave old vxnet
	if oldV.(string) != "" {
		if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
			return err
		}
		leaveVxnetInput := new(qc.LeaveVxNetInput)
		leaveVxnetInput.Instances = []*string{qc.String(d.Id())}
		leaveVxnetInput.VxNet = qc.String(oldV.(string))
		var err error
		simpleRetry(func() error {
			_, err = vxnetClt.LeaveVxNet(leaveVxnetInput)
			return isServerBusy(err)
		})
		if err != nil {
			return err
		}
		if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
			return err
		}
	}
	if newV.(string) != "" {
		selfManaged, err := isVxnetSelfManaged(newV.(string), vxnetClt)
		if err != nil {
			return err
		}
		if selfManaged {
			return fmt.Errorf("can not use selfManaged ip as Managed ip")
		}
		joinVxnetInput := new(qc.JoinVxNetInput)
		if newV.(string) != BasicNetworkID && d.HasChange(resourceInstancePrivateIP) && d.Get(resourceInstancePrivateIP).(string) != "" {
			newV = fmt.Sprintf("%s|%s", newV.(string), d.Get(resourceInstancePrivateIP).(string))
		}
		joinVxnetInput.Instances = []*string{qc.String(d.Id())}
		joinVxnetInput.VxNet = qc.String(newV.(string))
		simpleRetry(func() error {
			_, err = vxnetClt.JoinVxNet(joinVxnetInput)
			return isServerBusy(err)
		})
		if err != nil {
			return err
		}
		if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
			return err
		}
	}
	return nil
}

func instanceUpdateChangeSecurityGroup(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange(resourceInstanceSecurityGroupId) {
		return nil
	}
	clt := meta.(*QingCloudClient).instance
	sgClt := meta.(*QingCloudClient).securitygroup
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	input := new(qc.ApplySecurityGroupInput)
	input.SecurityGroup = getUpdateStringPointer(d, resourceInstanceSecurityGroupId)
	input.Instances = []*string{qc.String(d.Id())}
	var err error
	simpleRetry(func() error {
		_, err = sgClt.ApplySecurityGroup(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	return nil
}

func instanceUpdateChangeEip(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange(resourceInstanceEipID) {
		return nil
	}
	clt := meta.(*QingCloudClient).instance
	eipClt := meta.(*QingCloudClient).eip
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	oldV, newV := d.GetChange(resourceInstanceEipID)
	// dissociate old eip
	if oldV.(string) != "" {
		if _, err := EIPTransitionStateRefresh(eipClt, oldV.(string)); err != nil {
			return err
		}
		dissociateEIPInput := new(qc.DissociateEIPsInput)
		dissociateEIPInput.EIPs = []*string{qc.String(oldV.(string))}
		var err error
		simpleRetry(func() error {
			_, err = eipClt.DissociateEIPs(dissociateEIPInput)
			return isServerBusy(err)
		})
		if err != nil {
			return err
		}
		if _, err := EIPTransitionStateRefresh(eipClt, oldV.(string)); err != nil {
			return err
		}
	}
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	// associate new eip
	if newV.(string) != "" {
		if _, err := EIPTransitionStateRefresh(eipClt, newV.(string)); err != nil {
			return err
		}
		assoicateEIPInput := new(qc.AssociateEIPInput)
		assoicateEIPInput.EIP = qc.String(newV.(string))
		assoicateEIPInput.Instance = qc.String(d.Id())
		var err error
		simpleRetry(func() error {
			_, err = eipClt.AssociateEIP(assoicateEIPInput)
			return isServerBusy(err)
		})
		if err != nil {
			return err
		}
		if _, err := EIPTransitionStateRefresh(eipClt, newV.(string)); err != nil {
			return err
		}
	}
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	return nil
}

func instanceUpdateChangeKeyPairs(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange(resourceInstanceKeyPairIDs) {
		return nil
	}
	clt := meta.(*QingCloudClient).instance
	kpClt := meta.(*QingCloudClient).keypair

	oldV, newV := d.GetChange(resourceInstanceKeyPairIDs)
	var nkps []string
	var okps []string
	for _, v := range oldV.(*schema.Set).List() {
		okps = append(okps, v.(string))
	}
	for _, v := range newV.(*schema.Set).List() {
		nkps = append(nkps, v.(string))
	}
	additions, deletions := stringSliceDiff(nkps, okps)
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	// attach new key_pair
	if len(additions) > 0 {
		attachInput := new(qc.AttachKeyPairsInput)
		attachInput.Instances = []*string{qc.String(d.Id())}
		attachInput.KeyPairs = qc.StringSlice(additions)
		var err error
		simpleRetry(func() error {
			_, err = kpClt.AttachKeyPairs(attachInput)
			return isServerBusy(err)
		})
		if err != nil {
			return err
		}
		if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
			return err
		}
	}
	// dettach old key_pair
	if len(deletions) > 0 {
		detachInput := new(qc.DetachKeyPairsInput)
		detachInput.Instances = []*string{qc.String(d.Id())}
		detachInput.KeyPairs = qc.StringSlice(deletions)
		var err error
		simpleRetry(func() error {
			_, err = kpClt.DetachKeyPairs(detachInput)
			return isServerBusy(err)
		})
		if err != nil {
			return err
		}
		if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
			return err
		}
	}
	return nil
}

func instanceUpdateResize(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange(resourceInstanceCPU) && !d.HasChange(resourceInstanceMemory) || d.IsNewResource() {
		return nil
	}
	clt := meta.(*QingCloudClient).instance
	// check instance state
	describeInstanceOutput, err := describeInstance(d, meta)
	if err != nil {
		return err
	}
	instance := describeInstanceOutput.InstanceSet[0]
	// stop instance
	if *instance.Status == "running" {
		if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
			return err
		}
		_, err := stopInstance(d, meta)
		if err != nil {
			return err
		}
		if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
			return err
		}
	}
	//  resize instance
	var output *qc.ResizeInstancesOutput
	input := new(qc.ResizeInstancesInput)
	input.Instances = []*string{qc.String(d.Id())}
	input.CPU = qc.Int(d.Get(resourceInstanceCPU).(int))
	input.Memory = qc.Int(d.Get(resourceInstanceMemory).(int))
	simpleRetry(func() error {
		output, err = clt.ResizeInstances(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	client.WaitJob(meta.(*QingCloudClient).job,
		qc.StringValue(output.JobID),
		time.Duration(waitJobTimeOutDefault)*time.Second, time.Duration(waitJobIntervalDefault)*time.Second)
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	// start instance
	_, err = startInstance(d, meta)
	if err != nil {
		return err
	}
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	return nil
}

func describeInstance(d *schema.ResourceData, meta interface{}) (*qc.DescribeInstancesOutput, error) {
	clt := meta.(*QingCloudClient).instance
	input := new(qc.DescribeInstancesInput)
	input.Instances = []*string{qc.String(d.Id())}
	input.Verbose = qc.Int(1)
	var output *qc.DescribeInstancesOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeInstances(input)
		return isServerBusy(err)
	})
	if err != nil {
		return nil, err
	}
	return output, nil
}

func stopInstance(d *schema.ResourceData, meta interface{}) (*qc.StopInstancesOutput, error) {
	clt := meta.(*QingCloudClient).instance
	input := new(qc.StopInstancesInput)
	input.Instances = []*string{qc.String(d.Id())}
	var output *qc.StopInstancesOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.StopInstances(input)
		return isServerBusy(err)
	})
	if err != nil {
		return nil, err
	}
	return output, nil
}

func startInstance(d *schema.ResourceData, meta interface{}) (*qc.StartInstancesOutput, error) {
	clt := meta.(*QingCloudClient).instance
	input := new(qc.StartInstancesInput)
	input.Instances = []*string{qc.String(d.Id())}
	var output *qc.StartInstancesOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.StartInstances(input)
		return isServerBusy(err)
	})
	if err != nil {
		return nil, err
	}
	return output, nil
}

func updateInstanceVolume(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).instance
	if !d.HasChange(resourceInstanceVolumeIDs) {
		return nil
	}
	volumeClt := meta.(*QingCloudClient).volume
	oldV, newV := d.GetChange(resourceInstanceVolumeIDs)
	var newVolumes []string
	var oldVolumes []string
	for _, v := range oldV.(*schema.Set).List() {
		oldVolumes = append(oldVolumes, v.(string))
	}
	for _, v := range newV.(*schema.Set).List() {
		newVolumes = append(newVolumes, v.(string))
	}
	additions, deletions := stringSliceDiff(newVolumes, oldVolumes)
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	// attach new key_pair
	if len(additions) > 0 {
		attachInput := new(qc.AttachVolumesInput)
		attachInput.Instance = qc.String(d.Id())
		attachInput.Volumes = qc.StringSlice(additions)
		var err error
		simpleRetry(func() error {
			_, err = volumeClt.AttachVolumes(attachInput)
			return isServerBusy(err)
		})
		if err != nil {
			return err
		}
		for _, volumeID := range additions {
			VolumeTransitionStateRefresh(volumeClt, volumeID)
		}
		if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
			return err
		}
	}
	// dettach old key_pair
	if len(deletions) > 0 {
		detachInput := new(qc.DetachVolumesInput)
		detachInput.Instance = qc.String(d.Id())
		detachInput.Volumes = qc.StringSlice(deletions)
		var err error
		simpleRetry(func() error {
			_, err = volumeClt.DetachVolumes(detachInput)
			return isServerBusy(err)
		})
		if err != nil {
			return err
		}
		for _, volumeID := range deletions {
			VolumeTransitionStateRefresh(volumeClt, volumeID)
		}
		if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
			return err
		}
	}
	return nil
}

func waitInstanceLease(d *schema.ResourceData, meta interface{}) error {
	if !d.IsNewResource() {
		return nil
	}
	clt := meta.(*QingCloudClient).instance
	input := new(qc.DescribeInstancesInput)
	input.Instances = []*string{qc.String(d.Id())}
	input.Verbose = qc.Int(1)
	var output *qc.DescribeInstancesOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeInstances(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	//wait for lease info
	WaitForLease(output.InstanceSet[0].StatusTime)
	return nil
}

func isInstanceDeleted(instanceSet []*qc.Instance) bool {
	if len(instanceSet) == 0 || qc.StringValue(instanceSet[0].Status) == "terminated" || qc.StringValue(instanceSet[0].Status) == "ceased" {
		return true
	}
	return false
}
