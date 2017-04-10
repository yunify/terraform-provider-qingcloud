package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/lowstz/qingcloud-sdk-go/service"
)

func resourceQingcloudVolumeAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudVolumeAttachmentCreate,
		Read:   resourceQingcloudVolumeAttachmentRead,
		Update: nil,
		Delete: resourceQingcloudVolumeAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"volume_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"device_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceQingcloudVolumeAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	instanceClt := meta.(*QingCloudClient).instance
	volumeClt := meta.(*QingCloudClient).volume
	if _, err := InstanceTransitionStateRefresh(instanceClt, d.Get("instance_id").(string)); err != nil {
		return err
	}
	if _, err := VolumeTransitionStateRefresh(volumeClt, d.Get("volume_id").(string)); err != nil {
		return err
	}
	input := new(qc.AttachVolumesInput)
	input.Instance = qc.String(d.Get("instance_id").(string))
	input.Volumes = []*string{qc.String(d.Get("volume_id").(string))}
	_, err := volumeClt.AttachVolumes(input)
	if err != nil {
		return err
	}
	if _, err := InstanceTransitionStateRefresh(instanceClt, d.Get("instance_id").(string)); err != nil {
		return err
	}
	if _, err := VolumeTransitionStateRefresh(volumeClt, d.Get("volume_id").(string)); err != nil {
		return err
	}
	d.SetId(genVolumeAttachmentID(d))
	return resourceQingcloudVolumeAttachmentRead(d, meta)
}

func resourceQingcloudVolumeAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	volumeID, _, err := getInstanceIDAndVolumeID(d)
	if err != nil {
		return err
	}
	volumeClt := meta.(*QingCloudClient).volume
	input := new(qc.DescribeVolumesInput)
	input.Volumes = []*string{qc.String(volumeID)}
	output, err := volumeClt.DescribeVolumes(input)
	if err != nil {
		return err
	}
	if len(output.VolumeSet) == 0 {
		d.SetId("")
		return nil
	}
	volume := output.VolumeSet[0]
	d.Set("instance_id", qc.StringValue(volume.Instance.InstanceID))
	d.Set("volume_id", qc.StringValue(volume.VolumeID))
	d.Set("device_name", qc.StringValue(volume.Device))
	return nil
}

func resourceQingcloudVolumeAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	volumeID, instanceID, err := getInstanceIDAndVolumeID(d)
	if err != nil {
		return err
	}
	instanceClt := meta.(*QingCloudClient).instance
	volumeClt := meta.(*QingCloudClient).volume
	if _, err := InstanceTransitionStateRefresh(instanceClt, instanceID); err != nil {
		return err
	}
	if _, err := VolumeTransitionStateRefresh(volumeClt, volumeID); err != nil {
		return err
	}
	input := new(qc.DetachVolumesInput)
	input.Instance = qc.String(d.Get("instance_id").(string))
	input.Volumes = []*string{qc.String(d.Get("volume_id").(string))}
	if _, err = volumeClt.DetachVolumes(input); err != nil {
		return err
	}
	d.SetId("")
	if _, err := InstanceTransitionStateRefresh(instanceClt, instanceID); err != nil {
		return err
	}
	if _, err := VolumeTransitionStateRefresh(volumeClt, volumeID); err != nil {
		return err
	}
	return nil
}
