package qingcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifyInstanceAttributes(d *schema.ResourceData, meta interface{}, create bool) error {
	clt := meta.(*QingCloudClient).instance
	input := new(qc.ModifyInstanceAttributesInput)
	input.Instance = qc.String(d.Id())
	if create {
		if description := d.Get("description").(string); description == "" {
			return nil
		}
		input.Description = qc.String(d.Get("description").(string))
	} else {
		if !d.HasChange("description") && !d.HasChange("name") {
			return nil
		}
		if d.HasChange("description") {
			input.Description = qc.String(d.Get("description").(string))
		}
		if d.HasChange("name") {
			input.InstanceName = qc.String(d.Get("name").(string))
		}
	}
	output, err := clt.ModifyInstanceAttributes(input)
	if err != nil {
		return fmt.Errorf("Error modify instance attributes: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error modify instance attributes: %s", *output.Message)
	}
	return nil
}

func instanceUpdateChangeVxNet(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange("vxnet_id") {
		return nil
	}
	clt := meta.(*QingCloudClient).instance
	vxnetClt := meta.(*QingCloudClient).vxnet
	oldV, newV := d.GetChange("vxnet_id")
	// leave old vxnet
	if oldV.(string) != "" {
		if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
			return err
		}
		leaveVxnetInput := new(qc.LeaveVxNetInput)
		leaveVxnetInput.Instances = []*string{qc.String(d.Id())}
		leaveVxnetInput.VxNet = qc.String(oldV.(string))

		leaveVxnetOutput, err := vxnetClt.LeaveVxNet(leaveVxnetInput)
		if err != nil {
			return fmt.Errorf("Error leave vxnet: %s", err)
		}
		if leaveVxnetOutput.RetCode != nil && qc.IntValue(leaveVxnetOutput.RetCode) != 0 {
			return fmt.Errorf("Error leave vxnet: %s", err)
		}
		if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
			return err
		}
	}
	// join new vxnet
	if newV.(string) != "" {
		joinVxnetInput := new(qc.JoinVxNetInput)
		joinVxnetInput.Instances = []*string{qc.String(d.Id())}
		joinVxnetInput.VxNet = qc.String(newV.(string))
		joinVxnetOutput, err := vxnetClt.JoinVxNet(joinVxnetInput)
		if err != nil {
			return fmt.Errorf("Error leave vxnet: %s", err)
		}
		if joinVxnetOutput.RetCode != nil && qc.IntValue(joinVxnetOutput.RetCode) != 0 {
			return fmt.Errorf("Error join vxnet: %s", err)
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
	output, err := sgClt.ApplySecurityGroup(input)
	if err != nil {
		return fmt.Errorf("Error apply security group: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error apply security group: %s", *output.Message)
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
	describeEIPInput := new(qc.DescribeEIPsInput)
	describeEIPInput.EIPs = []*string{qc.String(d.Get("eip_id").(string))}
	describeEIPOutput, err := eipClt.DescribeEIPs(describeEIPInput)
	if err != nil {
		return fmt.Errorf("Error describe eip: %s", err)
	}
	if describeEIPOutput.RetCode != nil && qc.IntValue(describeEIPOutput.RetCode) != 0 {
		return fmt.Errorf("Error describe eip: %s", *describeEIPOutput.Message)
	}
	if qc.StringValue(describeEIPOutput.EIPSet[0].Status) != "available" {
		return fmt.Errorf("Error eip %s state is %s", d.Get("eip_id").(string), qc.StringValue(describeEIPOutput.EIPSet[0].Status))
	}
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
		dissociateEIPOutput, err := eipClt.DissociateEIPs(dissociateEIPInput)
		if err != nil {
			return fmt.Errorf("Error dissociate eip: %s", err)
		}
		if dissociateEIPOutput.RetCode != nil && qc.IntValue(dissociateEIPOutput.RetCode) != 0 {
			return fmt.Errorf("Error dissocidate eip: %s", *dissociateEIPOutput.Message)
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
		assoicateEIPOutput, err := eipClt.AssociateEIP(assoicateEIPInput)
		if err != nil {
			return fmt.Errorf("Error assoicate eip: %s", err)
		}
		if assoicateEIPOutput.RetCode != nil && qc.IntValue(assoicateEIPOutput.RetCode) != 0 {
			return fmt.Errorf("Error assoicate eip: %s", err)
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
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err == nil {
		return err
	}
	// attach new key_pair
	if len(additions) > 0 {
		attachInput := new(qc.AttachKeyPairsInput)
		attachInput.Instances = []*string{qc.String(d.Id())}
		attachInput.KeyPairs = qc.StringSlice(additions)
		attachOutput, err := kpClt.AttachKeyPairs(attachInput)
		if err != nil {
			return fmt.Errorf("Error attach keypairs: %s", err)
		}
		if attachOutput.RetCode != nil && qc.IntValue(attachOutput.RetCode) != 0 {
			return fmt.Errorf("Error attach keypairs: %s", *attachOutput.Message)
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
		detachOutput, err := kpClt.DetachKeyPairs(detachInput)
		if err != nil {
			return fmt.Errorf("Errorr detach keypairs: %s", err)
		}
		if detachOutput.RetCode != nil && qc.IntValue(detachOutput.RetCode) != 0 {
			return fmt.Errorf("Error detach keypairs: %s", *detachOutput.Message)
		}
		if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
			return err
		}
	}
	return nil
}

func instanceUpdateResize(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange("cpu") && !d.HasChange("memory") {
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
			log.Printf("InstanceTransitionStateRefresh err: %s", err)
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
	if d.HasChange("cpu") {
		input.CPU = qc.Int(d.Get("cpu").(int))
	}
	if d.HasChange("memory") {
		input.Memory = qc.Int(d.Get("memory").(int))
	}
	output, err := clt.ResizeInstances(input)
	if err != nil {
		return fmt.Errorf("Error resize instance: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error resize instance: %s", err)
	}
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
	output, err := clt.DescribeInstances(input)
	if err != nil {
		return nil, fmt.Errorf("Error describe instance: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
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
		return nil, fmt.Errorf("Error stop instance: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return nil, fmt.Errorf("Error stop instance: %s", *output.Message)
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
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return nil, fmt.Errorf("Error start instance: %s", *output.Message)
	}
	return output, nil
}

func deleteInstanceLeaveVxnet(d *schema.ResourceData, meta interface{}) (*qc.LeaveVxNetOutput, error) {
	vxnetID := d.Get("vxnet_id").(string)
	if vxnetID != "" {
		clt := meta.(*QingCloudClient).vxnet
		input := new(qc.LeaveVxNetInput)
		input.Instances = []*string{qc.String(d.Id())}
		input.VxNet = qc.String(vxnetID)
		output, err := clt.LeaveVxNet(input)
		if err != nil {
			return nil, fmt.Errorf("Error instance leave vxnet: %s", err)
		}
		if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
			return output, fmt.Errorf("Error instance leave vxnet: %s", qc.StringValue(output.Message))
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
		output, err := clt.DissociateEIPs(input)
		if err != nil {
			return nil, fmt.Errorf("Error dissciate eip: %s", err)
		}
		if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
			return nil, fmt.Errorf("Error dissciate eip input validate: %s", *output.Message)
		}
		if _, err := EIPTransitionStateRefresh(clt, eipID); err != nil {
			return nil, err
		}
	}
	return nil, nil
}
