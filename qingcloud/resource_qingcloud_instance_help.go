package qingcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
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
		_, err := clt.ModifyInstanceAttributes(input)
		if err != nil {
			return err
		}
	}
	return nil
}

func instanceUpdateChangeManagedVxNet(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange("managed_vxnet_id") {
		return nil
	}
	clt := meta.(*QingCloudClient).instance
	vxnetClt := meta.(*QingCloudClient).vxnet
	oldV, newV := d.GetChange("managed_vxnet_id")
	// leave old vxnet
	if oldV.(string) != "" {
		if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
			return err
		}
		leaveVxnetInput := new(qc.LeaveVxNetInput)
		leaveVxnetInput.Instances = []*string{qc.String(d.Id())}
		leaveVxnetInput.VxNet = qc.String(oldV.(string))

		_, err := vxnetClt.LeaveVxNet(leaveVxnetInput)
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
		if d.Get("static_ip").(string) != "" {
			newV = fmt.Sprintf("%s|%s", newV.(string), d.Get("static_ip").(string))
		}
		joinVxnetInput.Instances = []*string{qc.String(d.Id())}
		joinVxnetInput.VxNet = qc.String(newV.(string))
		_, err = vxnetClt.JoinVxNet(joinVxnetInput)
		if err != nil {
			return fmt.Errorf("Error leave vxnet: %s", err)
		}
		if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
			return err
		}
	}
	return nil
}

func instanceUpdateChangeSecurityGroup(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange("security_group_id") {
		return nil
	}
	clt := meta.(*QingCloudClient).instance
	sgClt := meta.(*QingCloudClient).securitygroup
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	input := new(qc.ApplySecurityGroupInput)
	input.SecurityGroup = qc.String(d.Get("security_group_id").(string))
	input.Instances = []*string{qc.String(d.Id())}
	_, err := sgClt.ApplySecurityGroup(input)
	if err != nil {
		return err
	}
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	return nil
}

func instanceUpdateChangeEip(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange("eip_id") {
		return nil
	}
	clt := meta.(*QingCloudClient).instance
	eipClt := meta.(*QingCloudClient).eip
	if _, err := EIPTransitionStateRefresh(eipClt, d.Get("eip_id").(string)); err != nil {
		return err
	}
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	oldV, newV := d.GetChange("eip_id")
	// dissociate old eip
	if oldV.(string) != "" {
		dissociateEIPInput := new(qc.DissociateEIPsInput)
		dissociateEIPInput.EIPs = []*string{qc.String(oldV.(string))}
		_, err := eipClt.DissociateEIPs(dissociateEIPInput)
		if err != nil {
			return err
		}
	}

	if _, err := EIPTransitionStateRefresh(eipClt, d.Get("eip_id").(string)); err != nil {
		return err
	}
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	// associate new eip
	if newV.(string) != "" {
		assoicateEIPInput := new(qc.AssociateEIPInput)
		assoicateEIPInput.EIP = qc.String(newV.(string))
		assoicateEIPInput.Instance = qc.String(d.Id())
		_, err := eipClt.AssociateEIP(assoicateEIPInput)
		if err != nil {
			return err
		}
	}
	if _, err := EIPTransitionStateRefresh(eipClt, d.Get("eip_id").(string)); err != nil {
		return err
	}
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	return nil
}

func instanceUpdateChangeKeyPairs(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange("keypair_ids") {
		return nil
	}
	clt := meta.(*QingCloudClient).instance
	kpClt := meta.(*QingCloudClient).keypair

	oldV, newV := d.GetChange("keypair_ids")
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
		_, err := kpClt.AttachKeyPairs(attachInput)
		if err != nil {
			return err
		}

	}
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	// dettach old key_pair
	if len(deletions) > 0 {
		detachInput := new(qc.DetachKeyPairsInput)
		detachInput.Instances = []*string{qc.String(d.Id())}
		detachInput.KeyPairs = qc.StringSlice(deletions)
		_, err := kpClt.DetachKeyPairs(detachInput)
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
	if !d.HasChange("cpu") && !d.HasChange("memory") || d.IsNewResource() {
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
	input := new(qc.ResizeInstancesInput)
	input.Instances = []*string{qc.String(d.Id())}
	input.CPU = qc.Int(d.Get("cpu").(int))
	input.Memory = qc.Int(d.Get("memory").(int))
	_, err = clt.ResizeInstances(input)
	if err != nil {
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
	output, err := clt.DescribeInstances(input)
	if err != nil {
		return nil, fmt.Errorf("Error describe instance: %s", err)
	}
	return output, nil
}

func stopInstance(d *schema.ResourceData, meta interface{}) (*qc.StopInstancesOutput, error) {
	clt := meta.(*QingCloudClient).instance
	input := new(qc.StopInstancesInput)
	input.Instances = []*string{qc.String(d.Id())}
	output, err := clt.StopInstances(input)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func startInstance(d *schema.ResourceData, meta interface{}) (*qc.StartInstancesOutput, error) {
	clt := meta.(*QingCloudClient).instance
	input := new(qc.StartInstancesInput)
	input.Instances = []*string{qc.String(d.Id())}
	output, err := clt.StartInstances(input)
	if err != nil {
		return nil, fmt.Errorf("Error start instance: %s", err)
	}
	return output, nil
}

func deleteInstanceLeaveVxnet(d *schema.ResourceData, meta interface{}) (*qc.LeaveVxNetOutput, error) {
	vxnetID := d.Get("managed_vxnet_id").(string)
	if vxnetID != "" {
		clt := meta.(*QingCloudClient).vxnet
		input := new(qc.LeaveVxNetInput)
		input.Instances = []*string{qc.String(d.Id())}
		input.VxNet = qc.String(vxnetID)
		_, err := clt.LeaveVxNet(input)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func deleteInstanceDissociateEip(d *schema.ResourceData, meta interface{}) (*qc.LeaveVxNetOutput, error) {
	eipID := d.Get("eip_id").(string)
	if eipID != "" {
		clt := meta.(*QingCloudClient).eip
		if _, err := EIPTransitionStateRefresh(clt, eipID); err != nil {
			return nil, err
		}
		input := new(qc.DissociateEIPsInput)
		input.EIPs = []*string{qc.String(eipID)}
		_, err := clt.DissociateEIPs(input)
		if err != nil {
			return nil, fmt.Errorf("Error dissciate eip: %s", err)
		}
		if _, err := EIPTransitionStateRefresh(clt, eipID); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func updateInstanceVolume(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).instance
	if !d.HasChange("volume_ids") {
		return nil
	}
	volumeClt := meta.(*QingCloudClient).volume
	oldV, newV := d.GetChange("volume_ids")
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
		_, err := volumeClt.AttachVolumes(attachInput)
		if err != nil {
			return err
		}
	}
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	// dettach old key_pair
	if len(deletions) > 0 {
		detachInput := new(qc.DetachVolumesInput)
		detachInput.Instance = qc.String(d.Id())
		detachInput.Volumes = qc.StringSlice(deletions)
		_, err := volumeClt.DetachVolumes(detachInput)
		if err != nil {
			return err
		}
		if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
			return err
		}
	}
	return nil
}
