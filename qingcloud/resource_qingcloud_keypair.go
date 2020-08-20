package qingcloud

import (
	"strings"

	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
	"time"
)

const (
	resourceKeyPairPublicKey = "public_key"
)

func resourceQingcloudKeypair() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudKeypairCreate,
		Read:   resourceQingcloudKeypairRead,
		Update: resourceQingcloudKeypairUpdate,
		Delete: resourceQingcluodKeypairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			resourceName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceKeyPairPublicKey: {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						keypair := strings.Split(strings.TrimSpace(v.(string)), " ")
						if len(keypair) >= 2 {
							return keypair[0] + " " + keypair[1]
						}
						return ""
					default:
						return ""
					}
				},
			},
			resourceDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceTagIds:   tagIdsSchema(),
			resourceTagNames: tagNamesSchema(),
		},
	}
}

func resourceQingcloudKeypairCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).keypair
	input := new(qc.CreateKeyPairInput)
	input.KeyPairName, _ = getNamePointer(d)
	input.Mode = qc.String("user")
	input.PublicKey = qc.String(d.Get(resourceKeyPairPublicKey).(string))
	var output *qc.CreateKeyPairOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.CreateKeyPair(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	d.SetId(qc.StringValue(output.KeyPairID))

	return resourceQingcloudKeypairUpdate(d, meta)
}

func resourceQingcloudKeypairRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).keypair
	input := new(qc.DescribeKeyPairsInput)
	input.KeyPairs = []*string{qc.String(d.Id())}
	var output *qc.DescribeKeyPairsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeKeyPairs(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if len(output.KeyPairSet) == 0 {
		d.SetId("")
		return nil
	}
	kp := output.KeyPairSet[0]
	d.Set(resourceName, qc.StringValue(kp.KeyPairName))
	d.Set(resourceDescription, qc.StringValue(kp.Description))
	d.Set(resourceKeyPairPublicKey, qc.StringValue(kp.EncryptMethod)+" "+qc.StringValue(kp.PubKey))
	if err := resourceSetTag(d, kp.Tags); err != nil {
		return err
	}
	return nil
}

func resourceQingcloudKeypairUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	if err := modifyKeypairAttributes(d, meta); err != nil {
		return err
	}
	d.SetPartial(resourceDescription)
	d.SetPartial(resourceName)
	if err := resourceUpdateTag(d, meta, qingcloudResourceTypeKeypair); err != nil {
		return err
	}
	d.SetPartial(resourceTagIds)
	d.Partial(false)
	return resourceQingcloudKeypairRead(d, meta)
}

func resourceQingcluodKeypairDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).keypair
	var err error
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		describeKeyPairsInput := new(qc.DescribeKeyPairsInput)
		describeKeyPairsInput.KeyPairs = []*string{qc.String(d.Id())}
		describeKeyPairsInput.Verbose = qc.Int(1)
		var describeKeyPairsOutput *qc.DescribeKeyPairsOutput
		describeKeyPairsOutput, err = clt.DescribeKeyPairs(describeKeyPairsInput)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		instanceIds := describeKeyPairsOutput.KeyPairSet[0].InstanceIDs

		if len(instanceIds) > 0 {
			detachKeyPairsInput := new(qc.DetachKeyPairsInput)
			detachKeyPairsInput.KeyPairs = describeKeyPairsInput.KeyPairs
			detachKeyPairsInput.Instances = instanceIds
			_, err = clt.DetachKeyPairs(detachKeyPairsInput)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			return resource.RetryableError(fmt.Errorf("there are still other resources [%v] depending on this resource[%v]", instanceIds, []*string{qc.String(d.Id())}))
		}
		input := new(qc.DeleteKeyPairsInput)
		input.KeyPairs = []*string{qc.String(d.Id())}
		_, err = clt.DeleteKeyPairs(input)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
}
