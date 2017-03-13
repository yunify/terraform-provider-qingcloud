package qingcloud

// import "github.com/hashicorp/terraform/helper/schema"

// func resourceQingcloudVolume() *schema.Resource {
// 	return &schema.Resource{
// 		Create: resourceQingcloudVolumeCreate,
// 		Read:   resourceQingcloudVolumeRead,
// 		Update: resourceQingcloudVolumeUpdate,
// 		Delete: resourceQingcloudVolumeDelete,
// 		Schema: map[string]*schema.Schema{
// 			"size": &schema.Schema{
// 				Type:     schema.TypeInt,
// 				Required: true,
// 				Description: "硬盘容量，目前可创建最小 10G，最大 500G 的硬盘， 在此范围内的容量值必须是 10 的倍数	",
// 			},
// 			"name": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Required: true,
// 			},
// 			"description": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Optional: true,
// 			},
// 			"type": &schema.Schema{
// 				Type:     schema.TypeInt,
// 				Optional: true,
// 				Default:  0,
// 				ForceNew: true,
// 				Description: `性能型是 0
// 					超高性能型是 3 (只能与超高性能主机挂载，目前只支持北京2区)，
// 					容量型因技术升级过程中，在各区的 type 值略有不同:
// 					  北京1区，亚太1区：容量型是 1
// 					  北京2区，广东1区：容量型是 2`,
// 			},
// 		},
// 	}
// }

// func motifyVolumeAttributes(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).volume
// 	modifyParams := volume.ModifyVolumeAttributesRequest{}
// 	modifyParams.Volume.Set(d.Id())
// 	modifyParams.Description.Set(d.Get("description").(string))
// 	modifyParams.VolumeName.Set(d.Get("name").(string))
// 	_, err := clt.ModifyVolumeAttributes(modifyParams)
// 	return err
// }

// func resourceQingcloudVolumeCreate(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).volume

// 	params := volume.CreateVolumesRequest{}
// 	params.Size.Set(d.Get("size").(int))
// 	params.VolumeName.Set(d.Get("name").(string))
// 	params.VolumeType.Set(d.Get("type").(int))

// 	resp, err := clt.CreateVolumes(params)
// 	if err != nil {
// 		return fmt.Errorf("Error creating volume: %s", err)
// 	}
// 	d.SetId(resp.Volumes[0])

// 	if err := changeQingcloudVolumeAttributes(d, meta); err != nil {
// 		return err
// 	}
// 	_, err = VolumeTransitionStateRefresh(clt, d.Id())
// 	return err
// }

// func resourceQingcloudVolumeRead(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).volume

// 	params := volume.DescribeVolumesRequest{}
// 	params.VolumesN.Add(d.Id())
// 	params.Verbose.Set(1)
// 	resp, err := clt.DescribeVolumes(params)
// 	if err != nil {
// 		return err
// 	}
// 	if len(resp.VolumeSet) == 0 {
// 		return nil
// 	}
// 	return nil
// }

// func resourceQingcloudVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).volume

// 	if !d.HasChange("size") && !d.HasChange("description") && d.HasChange("name") {
// 		return nil
// 	}

// 	if d.HasChange("size") {
// 		oldSize, newSize := d.GetChange("size")
// 		if oldSize.(int) > newSize.(int) {
// 			d.Set("size", oldSize.(int))
// 			return fmt.Errorf("Error you can only increase the size")
// 		}
// 		params := volume.ResizeVolumesRequest{}
// 		params.VolumesN.Add(d.Id())
// 		params.Size.Set(d.Get("size").(int))
// 		_, err := clt.ResizeVolumes(params)
// 		if err != nil {
// 			return fmt.Errorf("Error resize the volume: %s", err)
// 		}

// 	}

// 	if d.HasChange("description") || d.HasChange("name") {
// 		if err := changeQingcloudVolumeAttributes(d, meta); err != nil {
// 			return err
// 		}
// 	}

// 	return resourceQingcloudVolumeRead(d, meta)
// }

// func resourceQingcloudVolumeDelete(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).volume

// 	params := volume.DeleteVolumesRequest{}
// 	params.VolumesN.Add(d.Id())
// 	_, err := clt.DeleteVolumes(params)
// 	if err != nil {
// 		return fmt.Errorf(
// 			"Error deleting volume: %s", err)
// 	}

// 	_, err = VolumeTransitionStateRefresh(clt, d.Id())
// 	if err != nil {
// 		return fmt.Errorf(
// 			"Error waiting for volume (%s) to update: %s", d.Id(), err)
// 	}
// 	return nil
// }
