package qingcloud

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yunify/qingcloud-sdk-go/client"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

const (
	resourceVolumeSize = "size"
	resourceVolumeType = "type"
)

func resourceQingcloudVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudVolumeCreate,
		Read:   resourceQingcloudVolumeRead,
		Update: resourceQingcloudVolumeUpdate,
		Delete: resourceQingcloudVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			resourceVolumeSize: {
				Type:     schema.TypeInt,
				Required: true,
			},
			resourceName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceVolumeType: {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ForceNew:     true,
				ValidateFunc: withinArrayInt(0, 1, 2, 3),
			},
			resourceTagIds:   tagIdsSchema(),
			resourceTagNames: tagNamesSchema(),
		},
	}
}

func resourceQingcloudVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).volume
	input := new(qc.CreateVolumesInput)
	input.Count = qc.Int(1)
	input.Size = qc.Int(d.Get(resourceVolumeSize).(int))
	input.VolumeName, _ = getNamePointer(d)
	input.VolumeType = qc.Int(d.Get(resourceVolumeType).(int))
	var output *qc.CreateVolumesOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.CreateVolumes(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	d.SetId(qc.StringValue(output.Volumes[0]))
	if _, err = VolumeTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	return resourceQingcloudVolumeUpdate(d, meta)
}

func resourceQingcloudVolumeRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).volume
	input := new(qc.DescribeVolumesInput)
	input.Volumes = []*string{qc.String(d.Id())}
	var output *qc.DescribeVolumesOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeVolumes(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if len(output.VolumeSet) == 0 {
		d.SetId("")
		return nil
	}
	volume := output.VolumeSet[0]
	d.Set(resourceName, qc.StringValue(volume.VolumeName))
	d.Set(resourceDescription, qc.StringValue(volume.Description))
	d.Set(resourceVolumeSize, qc.IntValue(volume.Size))
	d.Set(resourceVolumeType, qc.IntValue(volume.VolumeType))
	if err := resourceSetTag(d, volume.Tags); err != nil {
		return err
	}
	return nil
}

func resourceQingcloudVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	if err := waitVolumeLease(d, meta); err != nil {
		return err
	}
	if err := motifyVolumeAttributes(d, meta); err != nil {
		return err
	}
	d.SetPartial(resourceName)
	d.SetPartial(resourceDescription)
	if err := changeVolumeSize(d, meta); err != nil {
		return err
	}
	d.SetPartial(resourceVolumeSize)
	if err := resourceUpdateTag(d, meta, qingcloudResourceTypeVolume); err != nil {
		return err
	}
	d.SetPartial(resourceTagIds)
	d.Partial(false)
	return resourceQingcloudVolumeRead(d, meta)
}

func resourceQingcloudVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	if err := waitVolumeLease(d, meta); err != nil {
		return err
	}
	clt := meta.(*QingCloudClient).volume
	if _, err := VolumeDeleteTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	input := new(qc.DeleteVolumesInput)
	input.Volumes = []*string{qc.String(d.Id())}
	var output *qc.DeleteVolumesOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DeleteVolumes(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	client.WaitJob(meta.(*QingCloudClient).job,
		qc.StringValue(output.JobID),
		time.Duration(10)*time.Second, time.Duration(1)*time.Second)

	d.SetId("")
	return nil
}
