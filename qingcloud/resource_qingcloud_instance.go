package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

const (
	resourceInstanceImageID         = "image_id"
	resourceInstanceCPU             = "cpu"
	resourceInstanceHostName        = "hostname"
	resourceInstanceMemory          = "memory"
	resourceInstanceClass           = "instance_class"
	resourceInstanceManagedVxnetID  = "managed_vxnet_id"
	resourceInstancePrivateIP       = "private_ip"
	resourceInstanceKeyPairIDs      = "keypair_ids"
	resourceInstanceSecurityGroupId = "security_group_id"
	resourceInstanceEipID           = "eip_id"
	resourceInstanceVolumeIDs       = "volume_ids"
	resourceInstancePublicIP        = "public_ip"
	resourceInstanceUserData        = "userdata"
	resourceInstanceLoginPassword   = "login_passwd"
	resourceInstanceOsDiskSize      = "os_disk_size"
)

func resourceQingcloudInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudInstanceCreate,
		Read:   resourceQingcloudInstanceRead,
		Update: resourceQingcloudInstanceUpdate,
		Delete: resourceQingcloudInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			resourceName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceInstanceImageID: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			resourceInstanceHostName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceInstanceCPU: {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: withinArrayInt(1, 2, 4, 8, 16),
				Default:      1,
			},
			resourceInstanceMemory: {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: withinArrayInt(1024, 2048, 4096, 6144, 8192, 12288, 16384, 24576, 32768),
				Default:      1024,
			},
			resourceInstanceClass: {
				Type:         schema.TypeInt,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: withinArrayInt(0, 1, 2, 3, 4, 5, 6, 100, 101, 200, 201, 300, 301),
				Default:      0,
			},
			resourceInstanceManagedVxnetID: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  BasicNetworkID,
			},
			resourceInstancePrivateIP: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			resourceInstanceOsDiskSize: {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: withinArrayIntRange(20, 100),
			},
			resourceInstanceKeyPairIDs: {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			resourceInstanceSecurityGroupId: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			resourceInstanceEipID: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceInstanceVolumeIDs: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			resourceInstancePublicIP: {
				Type:     schema.TypeString,
				Computed: true,
			},
			resourceInstanceUserData: {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
					s := v.(string)
					if !isBase64Encoded([]byte(s)) {
						errs = append(errs, fmt.Errorf(
							"%s: must be base64-encoded", name,
						))
					}
					return
				},
			},
			resourceTagIds:   tagIdsSchema(),
			resourceTagNames: tagNamesSchema(),
			resourceInstanceLoginPassword: {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceQingcloudInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).instance
	input := new(qc.RunInstancesInput)
	input.Count = qc.Int(1)
	input.Hostname = getSetStringPointer(d, resourceInstanceHostName)
	input.InstanceName, _ = getNamePointer(d)
	input.ImageID = getSetStringPointer(d, resourceInstanceImageID)
	input.CPU = qc.Int(d.Get(resourceInstanceCPU).(int))
	input.Memory = qc.Int(d.Get(resourceInstanceMemory).(int))
	input.InstanceClass = qc.Int(d.Get(resourceInstanceClass).(int))
	input.SecurityGroup = getSetStringPointer(d, resourceInstanceSecurityGroupId)
	if d.Get(resourceInstanceOsDiskSize).(int) != 0 {
		input.OSDiskSize = qc.Int(d.Get(resourceInstanceOsDiskSize).(int))
	}

	kps := d.Get(resourceInstanceKeyPairIDs).(*schema.Set).List()
	if len(kps) > 0 {
		kp := kps[0].(string)
		input.LoginMode = qc.String("keypair")
		input.LoginKeyPair = qc.String(kp)
	} else if d.Get(resourceInstanceLoginPassword).(string) != "" {
		input.LoginMode = qc.String("passwd")
		input.LoginPasswd = qc.String(d.Get(resourceInstanceLoginPassword).(string))
	} else {
		return fmt.Errorf("loginMode is Required!")
	}

	if d.Get(resourceInstanceUserData).(string) != "" {
		if err := setInstanceUserData(d, meta, input); err != nil {
			return err
		}
	}
	var output *qc.RunInstancesOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.RunInstances(input)
		return isServerBusy(err)
	})
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
	var output *qc.DescribeInstancesOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeInstances(input)
		return isServerBusy(err)
	})
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
	d.Set(resourceInstanceImageID, qc.StringValue(instance.Image.ImageID))
	d.Set(resourceInstanceCPU, qc.IntValue(instance.VCPUsCurrent))
	d.Set(resourceInstanceMemory, qc.IntValue(instance.MemoryCurrent))
	d.Set(resourceInstanceClass, qc.IntValue(instance.InstanceClass))
	d.Set(resourceInstanceOsDiskSize, qc.IntValue(instance.Extra.OSDiskSize))
	//set managed vxnet
	for _, vxnet := range instance.VxNets {
		if qc.IntValue(vxnet.VxNetType) != 0 {
			if qc.IntValue(vxnet.VxNetType) == 1 {
				d.Set(resourceInstanceManagedVxnetID, qc.StringValue(vxnet.VxNetID))
				d.Set(resourceInstancePrivateIP, qc.StringValue(vxnet.PrivateIP))
			} else {
				d.Set(resourceInstanceManagedVxnetID, BasicNetworkID)
				d.Set(resourceInstancePrivateIP, qc.StringValue(vxnet.PrivateIP))
			}
		}
	}
	if instance.EIP != nil {
		d.Set(resourceInstanceEipID, qc.StringValue(instance.EIP.EIPID))
		d.Set(resourceInstancePublicIP, qc.StringValue(instance.EIP.EIPAddr))
	}
	if instance.SecurityGroup != nil {
		d.Set(resourceInstanceSecurityGroupId, qc.StringValue(instance.SecurityGroup.SecurityGroupID))
	}
	if instance.KeyPairIDs != nil {
		d.Set(resourceInstanceKeyPairIDs, qc.StringValueSlice(instance.KeyPairIDs))
	}
	if instance.Volumes != nil {
		volumeIDs := make([]string, 0, len(instance.Volumes))
		for _, volume := range instance.Volumes {
			volumeIDs = append(volumeIDs, qc.StringValue(volume.VolumeID))
		}
		if err := d.Set(resourceInstanceVolumeIDs, volumeIDs); err != nil {
			return err
		}
	}
	if err := resourceSetTag(d, instance.Tags); err != nil {
		return err
	}
	return nil
}

func resourceQingcloudInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	if err := waitInstanceLease(d, meta); err != nil {
		return err
	}
	if err := modifyInstanceAttributes(d, meta); err != nil {
		return err
	}
	d.SetPartial(resourceName)
	d.SetPartial(resourceDescription)
	// change vxnet
	if err := instanceUpdateChangeManagedVxNet(d, meta); err != nil {
		return err
	}
	d.SetPartial(resourceInstanceManagedVxnetID)
	d.SetPartial(resourceInstancePrivateIP)
	// change security_group
	if err := instanceUpdateChangeSecurityGroup(d, meta); err != nil {
		return err
	}
	d.SetPartial(resourceInstanceSecurityGroupId)
	// change eip
	if err := instanceUpdateChangeEip(d, meta); err != nil {
		return err
	}
	d.SetPartial(resourceInstanceEipID)
	// change keypairs
	if err := instanceUpdateChangeKeyPairs(d, meta); err != nil {
		return err
	}
	d.SetPartial(resourceInstanceKeyPairIDs)
	// change volumes
	if err := updateInstanceVolume(d, meta); err != nil {
		return err
	}
	d.SetPartial(resourceInstanceVolumeIDs)
	// resize instance
	if err := instanceUpdateResize(d, meta); err != nil {
		return err
	}
	d.SetPartial(resourceInstanceCPU)
	d.SetPartial(resourceInstanceMemory)
	if err := resourceUpdateTag(d, meta, qingcloudResourceTypeInstance); err != nil {
		return err
	}
	d.Partial(false)
	return resourceQingcloudInstanceRead(d, meta)
}

func resourceQingcloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	if err := waitInstanceLease(d, meta); err != nil {
		return err
	}
	clt := meta.(*QingCloudClient).instance
	input := new(qc.TerminateInstancesInput)
	input.Instances = []*string{qc.String(d.Id())}
	var err error
	simpleRetry(func() error {
		_, err = clt.TerminateInstances(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if _, err := InstanceTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
