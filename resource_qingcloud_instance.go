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
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"image_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
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
			"vxnet_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"eip_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"eip_addr": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
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
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	return resourceQingcloudInstanceRead(d, meta)
}

func resourceQingcloudInstanceRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).instance

	params := instance.DescribeInstanceRequest{}
	params.InstancesN.Add(d.Id())
	params.Verbose.Set(1)
	resp, err := clt.DescribeInstances(params)
	if err != nil {
		return fmt.Errorf("Descirbe Instance :%s", err)
	}

	// TODO: if this is nil
	k := resp.InstanceSet[0]
	// TODO: not setting the default value
	d.Set("instance_type", k.InstanceType)
	if len(k.Vxnets) >= 1 {
		d.Set("vxnet_name", k.Vxnets[0].VxnetName)
		d.Set("vxnet_id", k.Vxnets[0].VxnetID)
		d.Set("private_ip", k.Vxnets[0].PrivateIP)
	}
	d.Set("eip_addr", k.Eip.EipAddr)

	// d.Set("eip_id", k.Eip.EipID)
	// instanceKeypairsId := make([]string, 0, len(k.KeypairIds))
	// for _, kp := range k.KeypairIds {
	// 	instanceKeypairsId = append(instanceKeypairsId, kp)
	// }
	// d.Set("keypair_ids", instanceKeypairsId)
	// d.Set("security_group_id", k.SecurityGroup.SecurityGroupID)
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
