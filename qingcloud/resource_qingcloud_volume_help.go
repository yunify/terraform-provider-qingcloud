package qingcloud

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"

	qc "github.com/lowstz/qingcloud-sdk-go/service"
)

func motifyVolumeAttributes(d *schema.ResourceData, meta interface{}, create bool) error {
	clt := meta.(*QingCloudClient).volume
	input := new(qc.ModifyVolumeAttributesInput)
	input.Volume = qc.String(d.Id())
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
			input.VolumeName = qc.String(d.Get("name").(string))
		}
	}
	_, err := clt.ModifyVolumeAttributes(input)
	return err
}

func changeVolumeSize(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange("size") {
		return nil
	}
	clt := meta.(*QingCloudClient).volume

	// new size must bigger than old size
	oldV, newV := d.GetChange("size")
	oldSize := oldV.(int)
	newSize := newV.(int)
	if oldSize >= newSize {
		return errors.New("volume size can't reduce")
	}
	if newSize%10 != 0 {
		return errors.New("volume size must be a multiple of 10")
	}
	// if disk is attached, shutdown instance, detach disk,
	if d.Get("status").(string) == "in-use" {
		instanceID := d.Get("instance_id").(string)

		instanceClt := meta.(*QingCloudClient).instance
		stopInstanceInput := new(qc.StopInstancesInput)
		qingcloudMutexKV.Lock(instanceID)
		defer qingcloudMutexKV.Unlock(instanceID)
		stopInstanceInput.Instances = []*string{qc.String(instanceID)}
		if _, err := instanceClt.StopInstances(stopInstanceInput); err != nil {
			return err
		}
		if _, err := InstanceTransitionStateRefresh(instanceClt, instanceID); err != nil {
			return err
		}
		detachVolumeInput := new(qc.DetachVolumesInput)
		detachVolumeInput.Instance = qc.String(instanceID)
		detachVolumeInput.Volumes = []*string{qc.String(d.Id())}
		if _, err := clt.DetachVolumes(detachVolumeInput); err != nil {
			return err
		}
		if _, err := VolumeTransitionStateRefresh(clt, d.Id()); err != nil {
			return err
		}
	}
	// increase disk size
	input := new(qc.ResizeVolumesInput)
	input.Volumes = []*string{qc.String(d.Id())}
	input.Size = qc.Int(newSize)
	if _, err := clt.ResizeVolumes(input); err != nil {
		return err
	}
	if _, err := VolumeTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}

	// attach disk, running instance
	if d.Get("status").(string) == "in-use" {
		instanceID := d.Get("instance_id").(string)
		attachVolumeInput := new(qc.AttachVolumesInput)
		attachVolumeInput.Instance = qc.String(instanceID)
		attachVolumeInput.Volumes = []*string{qc.String(d.Id())}
		if _, err := clt.AttachVolumes(attachVolumeInput); err != nil {
			return err
		}
		if _, err := VolumeTransitionStateRefresh(clt, d.Id()); err != nil {
			return err
		}
		instanceClt := meta.(*QingCloudClient).instance
		startInstanceInput := new(qc.StartInstancesInput)
		startInstanceInput.Instances = []*string{qc.String(instanceID)}
		if _, err := instanceClt.StartInstances(startInstanceInput); err != nil {
			return err
		}
		if _, err := InstanceTransitionStateRefresh(instanceClt, instanceID); err != nil {
			return err
		}
	}
	return nil
}
