package qingcloud

import (
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
			resourceName: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceDescription: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"managed_vxnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "vxnet-0",
			},
			"keypair_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"security_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"eip_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"public_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			resourceTagIds:   tagIdsSchema(),
			resourceTagNames: tagNamesSchema(),
		},
	}
}

func resourceQingcloudInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).instance
	input := new(qc.RunInstancesInput)
	input.Count = qc.Int(1)
	input.InstanceName, _ = getNamePointer(d)
	input.ImageID = qc.String(d.Get("image_id").(string))
	input.CPU = qc.Int(d.Get("cpu").(int))
	input.Memory = qc.Int(d.Get("memory").(int))
	input.SecurityGroup = qc.String(d.Get("security_group_id").(string))
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
	return resourceQingcloudInstanceUpdate(d, meta)
}

func resourceQingcloudInstanceRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).instance
	input := new(qc.DescribeInstancesInput)
	input.Instances = []*string{qc.String(d.Id())}
	input.Verbose = qc.Int(1)
	output, err := clt.DescribeInstances(input)
	if err != nil {
		return err
	}
	if isInstanceDeleted(output.InstanceSet) {
		d.SetId("")
		return nil
	}
	instance := output.InstanceSet[0]
	d.Set(resourceName, qc.StringValue(instance.InstanceName))
	d.Set(resourceDescription, qc.StringValue(instance.Description))
	d.Set("image_id", qc.StringValue(instance.Image.ImageID))
	d.Set("cpu", qc.IntValue(instance.VCPUsCurrent))
	d.Set("memory", qc.IntValue(instance.MemoryCurrent))
	//set managed vxnet
	for _, vxnet := range instance.VxNets {
		if qc.IntValue(vxnet.VxNetType) != 0 {
			d.Set("managed_vxnet_id", qc.StringValue(vxnet.VxNetID))
			d.Set("private_ip", qc.StringValue(vxnet.PrivateIP))
			break
		}
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
	if instance.Volumes != nil {
		volumeIDs := make([]string, 0, len(instance.Volumes))
		for _, volume := range instance.Volumes {
			volumeIDs = append(volumeIDs, qc.StringValue(volume.VolumeID))
		}
		d.Set("volume_ids", volumeIDs)
	}
	resourceSetTag(d, instance.Tags)
	return nil
}

func resourceQingcloudInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := waitInstanceLease(d, meta); err != nil {
		return err
	}
	if err := modifyInstanceAttributes(d, meta); err != nil {
		return err
	}
	// change vxnet
	if err := instanceUpdateChangeManagedVxNet(d, meta); err != nil {
		return err
	}
	// change security_group
	if err := instanceUpdateChangeSecurityGroup(d, meta); err != nil {
		return err
	}
	// change eip
	if err := instanceUpdateChangeEip(d, meta); err != nil {
		return err
	}
	// change keypairs
	if err := instanceUpdateChangeKeyPairs(d, meta); err != nil {
		return err
	}
	// change volumes
	if err := updateInstanceVolume(d, meta); err != nil {
		return err
	}
	// resize instance
	if err := instanceUpdateResize(d, meta); err != nil {
		return err
	}
	if err := resourceUpdateTag(d, meta, qingcloudResourceTypeInstance); err != nil {
		return err
	}
	return resourceQingcloudInstanceRead(d, meta)
}

func resourceQingcloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	if err := waitInstanceLease(d, meta); err != nil {
		return err
	}
	clt := meta.(*QingCloudClient).instance
	input := new(qc.TerminateInstancesInput)
	input.Instances = []*string{qc.String(d.Id())}
	if _, err := clt.TerminateInstances(input); err != nil {
		return err
	}
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
