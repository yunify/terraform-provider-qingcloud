package qingcloud

import (
	"fmt"
	"time"

	"errors"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/magicshui/qingcloud-go/volume"
)

func resourceQingcloudVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudVolumeCreate,
		Read:   resourceQingcloudVolumeRead,
		Update: resourceQingcloudVolumeUpdate,
		Delete: resourceQingcloudVolumeDelete,
		Schema: resouceQingcloudVolumeSchema(),
	}
}

// Waiting for no transition_status
func VolumeTransitionStateRefresh(clt *volume.VOLUME, id string) *resource.StateChangeConf {
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
		Delay:      2 * time.Second,
		MinTimeout: 1 * time.Second,
	}
	return stateConf
}

func resourceQingcloudVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).volume

	params := volume.CreateVolumesRequest{}
	params.Size.Set(d.Get("size").(int))
	params.VolumeName.Set(d.Get("name").(string))
	params.VolumeType.Set(d.Get("type").(int))

	resp, err := clt.CreateVolumes(params)
	if err != nil {
		return fmt.Errorf("Error creating volume: %s", err)
	}
	d.SetId(resp.Volumes[0])

	if d.Get("description") != nil {
		modifyParams := volume.ModifyVolumeAttributesRequest{}
		modifyParams.Volume.Set(d.Id())
		modifyParams.Description.Set(d.Get("description").(string))
		_, err = clt.ModifyVolumeAttributes(modifyParams)
		if err != nil {
			return fmt.Errorf("Error modify the volume attributes", err)
		}
	}

	return resourceQingcloudVolumeRead(d, meta)
}

func resourceQingcloudVolumeRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).volume

	_, err := VolumeTransitionStateRefresh(clt, d.Id()).WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for volume (%s) to update: %s", d.Id(), err)
	}

	params := volume.DescribeVolumesRequest{}
	params.VolumesN.Add(d.Id())
	params.Verbose.Set(1)
	resp, err := clt.DescribeVolumes(params)
	if err != nil {
		return fmt.Errorf("Error read volume %s", err)
	}

	for _, v := range resp.VolumeSet {
		if v.VolumeID == d.Id() {
			d.Set("id", v.VolumeID)
			d.Set("size", v.Size)
			d.Set("instance_id", v.Instance.InstanceID)
			d.Set("instance_name", v.Instance.InstanceName)
			d.Set("status", v.Status)
			return nil
		}
	}
	return nil
}

func resourceQingcloudVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).volume

	if !d.HasChange("size") && !d.HasChange("description") && d.HasChange("name") {
		return nil
	}

	//
	if d.HasChange("size") {
		oldSize, newSize := d.GetChange("size")
		if oldSize.(int) > newSize.(int) {
			d.Set("size", oldSize.(int))
			return fmt.Errorf("Error you can only increase the size", errors.New("INCREASE SIZE ONLY"))
		}
		params := volume.ResizeVolumesRequest{}
		params.VolumesN.Add(d.Id())
		params.Size.Set(d.Get("size").(int))
		_, err := clt.ResizeVolumes(params)
		if err != nil {
			return fmt.Errorf("Error resize the volume: %s", err)
		}

	}

	if d.HasChange("description") || d.HasChange("name") {
		params := volume.ModifyVolumeAttributesRequest{}
		params.Volume.Set(d.Id())
		params.VolumeName.Set(d.Get("name").(string))
		params.Description.Set(d.Get("description").(string))

		_, err := clt.ModifyVolumeAttributes(params)
		if err != nil {
			return fmt.Errorf("Error update the volume: %s", err)
		}
	}
	return resourceQingcloudVolumeRead(d, meta)
}

func resourceQingcloudVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).volume

	params := volume.DeleteVolumesRequest{}
	params.VolumesN.Add(d.Id())
	_, err := clt.DeleteVolumes(params)
	if err != nil {
		return fmt.Errorf(
			"Error deleting volume: %s", err)
	}

	_, err = VolumeTransitionStateRefresh(clt, d.Id()).WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for volume (%s) to update: %s", d.Id(), err)
	}
	return nil
}

func resouceQingcloudVolumeSchema() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"size": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
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
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
			ForceNew: true,
		},
		"instance_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"instance_name": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
	}
}
