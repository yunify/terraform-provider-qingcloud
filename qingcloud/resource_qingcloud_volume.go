package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"

	qc "github.com/lowstz/qingcloud-sdk-go/service"
)

func resourceQingcloudVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudVolumeCreate,
		Read:   resourceQingcloudVolumeRead,
		Update: resourceQingcloudVolumeUpdate,
		Delete: resourceQingcloudVolumeDelete,
		Schema: map[string]*schema.Schema{
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				Description: "硬盘容量，目前可创建最小 10G，最大 500G 的硬盘， 在此范围内的容量值必须是 10 的倍数	",
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ForceNew:     true,
				ValidateFunc: withinArrayInt(0, 1, 2, 3),
				Description: `性能型是 0
					超高性能型是 3 (只能与超高性能主机挂载，目前只支持北京2区)，
					容量型因技术升级过程中，在各区的 type 值略有不同:
					  北京1区，亚太1区：容量型是 1
					  北京2区，广东1区：容量型是 2`,
			},
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"device": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceQingcloudVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).volume
	input := new(qc.CreateVolumesInput)
	input.Count = qc.Int(1)
	input.Size = qc.Int(d.Get("size").(int))
	input.VolumeName = qc.String(d.Get("name").(string))
	input.VolumeType = qc.Int(d.Get("type").(int))
	output, err := clt.CreateVolumes(input)
	if err != nil {
		return err
	}
	d.SetId(qc.StringValue(output.Volumes[0]))
	if err := motifyVolumeAttributes(d, meta, true); err != nil {
		return err
	}
	if _, err = VolumeTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	return resourceQingcloudVolumeRead(d, meta)
}

func resourceQingcloudVolumeRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).volume
	input := new(qc.DescribeVolumesInput)
	input.Volumes = []*string{qc.String(d.Id())}
	output, err := clt.DescribeVolumes(input)
	if err != nil {
		return err
	}
	volume := output.VolumeSet[0]
	d.Set("name", qc.StringValue(volume.VolumeName))
	d.Set("description", qc.StringValue(volume.Description))
	d.Set("size", qc.IntValue(volume.Size))
	d.Set("type", qc.IntValue(volume.VolumeType))
	d.Set("status", qc.StringValue(volume.Status))
	if volume.Instance != nil {
		d.Set("instance_id", qc.StringValue(volume.Instance.InstanceID))
		d.Set("device", qc.StringValue(volume.Instance.Device))
	} else {
		d.Set("instance_id", "")
		d.Set("device", "")
	}
	return nil
}

func resourceQingcloudVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := motifyVolumeAttributes(d, meta, false); err != nil {
		return err
	}
	if err := changeVolumeSize(d, meta); err != nil {
		return err
	}
	return resourceQingcloudVolumeRead(d, meta)
}

func resourceQingcloudVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).volume
	if _, err := VolumeDeleteTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	input := new(qc.DeleteVolumesInput)
	input.Volumes = []*string{qc.String(d.Id())}
	if _, err := clt.DeleteVolumes(input); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
