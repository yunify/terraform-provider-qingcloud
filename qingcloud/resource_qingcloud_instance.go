package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
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
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_class": &schema.Schema{
				Type:         schema.TypeInt,
				Default:      0,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: withinArrayInt(0, 1),
			},
			"instance_state": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "running",
				ValidateFunc: withinArrayString("pending", "running", "stopped", "suspended", "terminated", "ceased"),
			},
			"cpu": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: withinArrayInt(1, 2, 4, 8, 16),
				Default:      1,
			},
			"memory": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: withinArrayInt(1024, 2048, 4096, 6144, 8192, 12288, 16384, 24576, 32768),
				Default:      1024,
			},
			"vxnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"static_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"keypair_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"security_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"eip_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume_device_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag_ids":   tagIdsSchema(),
			"tag_names": tagNamesSchema(),
		},
	}
}

func resourceQingcloudInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).instance
	input := new(qc.RunInstancesInput)
	input.Count = qc.Int(1)
	input.InstanceName = qc.String(d.Get("name").(string))
	input.ImageID = qc.String(d.Get("image_id").(string))
	input.InstanceClass = qc.Int(d.Get("instance_class").(int))
	// input.InstanceType = qc.String(d.Get("instance_type").(string))
	if d.Get("cpu").(int) != 0 && d.Get("memory").(int) != 0 {
		input.CPU = qc.Int(d.Get("cpu").(int))
		input.Memory = qc.Int(d.Get("memory").(int))
	}
	var vxnet string
	if d.Get("static_ip").(string) != "" {
		vxnet = fmt.Sprintf("%s|%s", d.Get("vxnet_id").(string), d.Get("static_ip").(string))
	} else {
		vxnet = d.Get("vxnet_id").(string)
	}
	input.VxNets = []*string{qc.String(vxnet)}
	if d.Get("security_group_id").(string) != "" {
		input.SecurityGroup = qc.String(d.Get("security_group_id").(string))
	}
	input.LoginMode = qc.String("keypair")
	kps := d.Get("keypair_ids").(*schema.Set).List()
	if len(kps) > 0 {
		kp := kps[0].(string)
		input.LoginKeyPair = qc.String(kp)
	}
	output, err := clt.RunInstances(input)
	if err != nil {
		return err
	}
	d.SetId(qc.StringValue(output.Instances[0]))
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	err = modifyInstanceAttributes(d, meta, true)
	if err != nil {
		return err
	}
	// associate eip to instance
	if eipID := d.Get("eip_id").(string); eipID != "" {
		eipClt := meta.(*QingCloudClient).eip
		if _, err := EIPTransitionStateRefresh(eipClt, eipID); err != nil {
			return err
		}
		associateEIPInput := new(qc.AssociateEIPInput)
		associateEIPInput.EIP = qc.String(eipID)
		associateEIPInput.Instance = qc.String(d.Id())
		associateEIPoutput, err := eipClt.AssociateEIP(associateEIPInput)
		if err != nil {
			return fmt.Errorf("Error associate eip: %s", err)
		}
		if associateEIPoutput.RetCode != nil && qc.IntValue(associateEIPoutput.RetCode) != 0 {
			return fmt.Errorf("Error associate eip: %s", *associateEIPoutput.Message)
		}
		if _, err := EIPTransitionStateRefresh(eipClt, eipID); err != nil {
			return err
		}
	}
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	if err := resourceUpdateTag(d, meta, qingcloudResourceTypeInstance); err != nil {
		return err
	}

	// update volume
	// volumeDS :=
	return resourceQingcloudInstanceRead(d, meta)
}

func resourceQingcloudInstanceRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).instance
	input := new(qc.DescribeInstancesInput)
	input.Instances = []*string{qc.String(d.Id())}
	input.Verbose = qc.Int(1)
	output, err := clt.DescribeInstances(input)
	if err != nil {
		return fmt.Errorf("Error describe instance: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error describe instance: %s", *output.Message)
	}
	if len(output.InstanceSet) == 0 {
		d.SetId("")
		return nil
	}

	instance := output.InstanceSet[0]
	d.Set("name", qc.StringValue(instance.InstanceName))
	d.Set("image_id", qc.StringValue(instance.Image.ImageID))
	d.Set("description", qc.StringValue(instance.Description))
	d.Set("instance_class", qc.IntValue(instance.InstanceClass))
	d.Set("instance_state", qc.StringValue(instance.Status))
	d.Set("cpu", qc.IntValue(instance.VCPUsCurrent))
	d.Set("memory", qc.IntValue(instance.MemoryCurrent))
	if instance.VxNets != nil && len(instance.VxNets) > 0 {
		vxnet := instance.VxNets[0]
		if qc.IntValue(vxnet.VxNetType) == 2 {
			d.Set("vxnet_id", "vxnet-0")
		} else {
			d.Set("vxnet_id", qc.StringValue(vxnet.VxNetID))
		}
		d.Set("private_ip", qc.StringValue(vxnet.PrivateIP))
		if d.Get("static_ip") != "" {
			d.Set("static_ip", qc.StringValue(vxnet.PrivateIP))
		}
	} else {
		d.Set("vxnet_id", "")
		d.Set("private_ip", "")
	}
	if instance.EIP != nil {
		d.Set("eip_id", qc.StringValue(instance.EIP.EIPID))
		d.Set("public_ip", qc.StringValue(instance.EIP.EIPAddr))
	}
	if instance.SecurityGroup != nil {
		d.Set("security_group_id", qc.StringValue(instance.SecurityGroup.SecurityGroupID))
	}
	if instance.KeyPairIDs != nil {
		keypairIDs := make([]string, 0, len(instance.KeyPairIDs))
		for _, kp := range instance.KeyPairIDs {
			keypairIDs = append(keypairIDs, qc.StringValue(kp))
		}
		d.Set("keypair_ids", keypairIDs)
	}
	resourceSetTag(d, instance.Tags)
	return nil
}

func resourceQingcloudInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	// clt := meta.(*QingCloudClient).instance
	err := modifyInstanceAttributes(d, meta, false)
	if err != nil {
		return err
	}
	// change vxnet
	err = instanceUpdateChangeVxNet(d, meta)
	if err != nil {
		return err
	}
	// change security_group
	err = instanceUpdateChangeSecurityGroup(d, meta)
	if err != nil {
		return err
	}
	// change eip
	err = instanceUpdateChangeEip(d, meta)
	if err != nil {
		return err
	}
	// change keypair
	err = instanceUpdateChangeKeyPairs(d, meta)
	if err != nil {
		return err
	}
	// resize instance
	err = instanceUpdateResize(d, meta)
	if err != nil {
		return err
	}
	if err := resourceUpdateTag(d, meta, qingcloudResourceTypeInstance); err != nil {
		return err
	}
	return resourceQingcloudInstanceRead(d, meta)
}

func resourceQingcloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).instance
	// dissociate eip before leave vxnet
	if _, err := deleteInstanceDissociateEip(d, meta); err != nil {
		return err
	}
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	_, err := deleteInstanceLeaveVxnet(d, meta)
	if err != nil {
		return err
	}
	if _, err := InstanceNetworkTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	input := new(qc.TerminateInstancesInput)
	input.Instances = []*string{qc.String(d.Id())}
	output, err := clt.TerminateInstances(input)
	if err != nil {
		return fmt.Errorf("Error terminate instance: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error terminate instance: %s", *output.Message)
	}
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
