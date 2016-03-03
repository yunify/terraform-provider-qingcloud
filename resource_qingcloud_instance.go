package qingcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/magicshui/qingcloud-go/instance"
)

func resourceQingcloudInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudInstanceCreate,
		Read:   resourceQingcloudInstanceRead,
		Update: resourceQingcloudInstanceUpdate,
		Delete: resourceQingcloudInstanceDelete,
		Schema: resourceQingcloudInstanceSchema(),
	}
}

func resourceQingcloudInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).instance

	params := instance.RunInstancesRequest{}
	params.InstanceName.Set(d.Get("name").(string))
	params.ImageId.Set(d.Get("image_id").(string))
	params.InstanceType.Set(d.Get("instance_type").(string))
	params.LoginMode.Set("keypair")
	params.VxnetsN.Add(d.Get("vxnet_id").(string))
	params.SecurityGroup.Set(d.Get("security_group_id").(string))
	for _, kp := range d.Get("keypair_ids").(*schema.Set).List() {
		params.LoginKeypair.Set(kp.(string))
	}
	params.InstanceClass.Set(d.Get("instance_class").(string))

	resp, err := clt.RunInstances(params)
	if err != nil {
		return fmt.Errorf("Error run instance :%s", err)
	}
	d.SetId(resp.Instances[0])
	return nil
}

func resourceQingcloudInstanceRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).instance

	params := instance.DescribeInstanceRequest{}
	params.InstancesN.Add(d.Id())
	params.Verbose.Set(1)

	resp, _ := clt.DescribeInstances(params)
	for _, k := range resp.InstanceSet {
		if d.Id() == k.InstanceID {
			d.Set("id", k.InstanceID)
			d.Set("name", k.InstanceName)
			d.Set("image_id", k.Image.ImageID)
			d.Set("instance_type", k.InstanceType)
			d.Set("vxnet_id", k.Vxnets[0].VxnetID)
			// keypair ids
			instanceKeypairsId := make([]string, 0, len(k.KeypairIds))
			for _, kp := range k.KeypairIds {
				instanceKeypairsId = append(instanceKeypairsId, kp)
			}
			d.Set("keypair_ids", instanceKeypairsId)
			d.Set("security_group_id", k.SecurityGroup.SecurityGroupID)
			return nil
		}
	}
	return nil
}

func resourceQingcloudInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceQingcloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).instance

	params := instance.StopInstancesRequest{}
	params.InstancesN.Add(d.Id())
	params.Force.Set(1)

	_, err := clt.StopInstances(params)
	if err != nil {
		return fmt.Errorf("Error run instance :%s", err)
	}
	return nil
}

func resourceQingcloudInstanceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		// 镜像类型
		"image_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		// 主机类型
		"instance_type": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		// cpu
		// memory
		// count
		// - login_mode 这个不实用
		//
		// 主机类别
		"instance_class": &schema.Schema{
			Type:     schema.TypeString,
			Default:  "0",
			Optional: true,
			ForceNew: true,
		},
		"vxnet_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"keypair_ids": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Set:      schema.HashString,
			Computed: true,
		},
		"security_group_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		// need_newsid
		// need_userdata
		// userdata_type
		// userdata_value
		// userdata_path
		// userdata_file
		"vxnets": &schema.Schema{
			Type:     schema.TypeSet,
			Computed: true,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"vxnet_name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"vxnet_type": &schema.Schema{
						Type:     schema.TypeInt,
						Optional: true,
					},
					"private_ip": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},

		"id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
	}

}
