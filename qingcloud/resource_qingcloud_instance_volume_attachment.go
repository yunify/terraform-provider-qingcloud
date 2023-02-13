package qingcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
	"log"
	"time"
)

const (
	resourceVolumes    = "volume_id"
	resourceInstanceId = "instance_id"
)

func resourceQingcloudInstanceVolumeAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudInstanceVolumeAttachmentCreate,
		Read:   resourceQingcloudInstanceVolumeAttachmentRead,
		Update: resourceQingcloudInstanceVolumeAttachmentUpdate,
		Delete: resourceQingcloudInstanceVolumeAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			resourceVolumes: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "绑定的磁盘id",
			},
			resourceInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "绑定的机器id",
			},
		},
	}
}

func resourceQingcloudInstanceVolumeAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).volume
	input := new(qc.AttachVolumesInput)
	volumeId := []string{d.Get(resourceVolumes).(string)}
	input.Volumes = qc.StringSlice(volumeId)
	input.Instance = qc.String(d.Get(resourceInstanceId).(string))
	var err error
	simpleRetry(func() error {
		_, err = clt.AttachVolumes(input)
		return WrapError(isServerBusy(err))
	})
	if err != nil {
		return WrapError(err)
	}
	d.SetId(fmt.Sprint(volumeId[0], ":", *input.Instance))
	if _, err := VolumeTransitionStateRefresh(clt, volumeId[0]); err != nil {
		return WrapError(err)
	}
	return resourceQingcloudInstanceVolumeAttachmentRead(d, meta)
}

func resourceQingcloudInstanceVolumeAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).volume
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	input := new(qc.DescribeVolumesInput)
	input.Volumes = []*string{qc.String(parts[0])}
	var output *qc.DescribeVolumesOutput
	simpleRetry(func() error {
		output, err = clt.DescribeVolumes(input)
		return WrapError(isServerBusy(err))
	})
	if err != nil {
		return WrapError(err)
	}
	if len(output.VolumeSet) == 0 {
		d.SetId("")
		return nil
	}
	volume := output.VolumeSet[0]
	if *volume.Status != string(InUse) && *volume.Instances[0].InstanceID != parts[1] {
		return WrapError(fmt.Errorf("the specified %s %s is not found", "VolumeAttach", d.Id()))
	}
	d.Set("volume_id", parts[0])
	d.Set("instance_id", parts[1])
	return nil
}

func resourceQingcloudInstanceVolumeAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceQingcloudInstanceVolumeAttachmentRead(d, meta)
}

func resourceQingcloudInstanceVolumeAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).volume
	input := new(qc.DetachVolumesInput)
	volumeId := []string{d.Get(resourceVolumes).(string)}
	input.Volumes = qc.StringSlice(volumeId)
	input.Instance = qc.String(d.Get(resourceInstanceId).(string))
	var err error
	simpleRetry(func() error {
		_, err = clt.DetachVolumes(input)
		return WrapError(isServerBusy(err))
	})
	if err != nil {
		return WrapError(err)
	}
	d.SetId("")
	if _, err := VolumeTransitionStateRefresh(clt, volumeId[0]); err != nil {
		return WrapError(err)
	}
	return nil
}
