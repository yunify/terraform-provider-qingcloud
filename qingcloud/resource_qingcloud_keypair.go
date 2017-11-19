package qingcloud

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func resourceQingcloudKeypair() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudKeypairCreate,
		Read:   resourceQingcloudKeypairRead,
		Update: resourceQingcloudKeypairUpdate,
		Delete: resourceQingcluodKeypairDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of keypair ",
			},
			"public_key": &schema.Schema{
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The SSH public key ",
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						return strings.TrimSpace(v.(string))
					default:
						return ""
					}
				},
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of keypair ",
			},
			"tag_ids":   tagIdsSchema(),
			"tag_names": tagNamesSchema(),
		},
	}
}

func resourceQingcloudKeypairCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).keypair
	input := new(qc.CreateKeyPairInput)
	input.KeyPairName, _ = getNamePointer(d)
	input.Mode = qc.String("user")
	input.PublicKey = qc.String(d.Get("public_key").(string))
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
	d.Set("name", qc.StringValue(kp.KeyPairName))
	d.Set("description", qc.StringValue(kp.Description))
	resourceSetTag(d, kp.Tags)
	return nil
}

func resourceQingcloudKeypairUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	if err := modifyKeypairAttributes(d, meta); err != nil {
		return err
	}
	d.SetPartial("description")
	d.SetPartial("name")
	if err := resourceUpdateTag(d, meta, qingcloudResourceTypeKeypair); err != nil {
		return err
	}
	d.SetPartial("tag_ids")
	d.Partial(false)
	return resourceQingcloudKeypairRead(d, meta)
}

func resourceQingcluodKeypairDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).keypair
	describeKeyPairsInput := new(qc.DescribeKeyPairsInput)
	describeKeyPairsInput.KeyPairs = []*string{qc.String(d.Id())}
	var describeKeyPairsOutput *qc.DescribeKeyPairsOutput
	var err error
	simpleRetry(func() error {
		describeKeyPairsOutput, err = clt.DescribeKeyPairs(describeKeyPairsInput)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	input := new(qc.DeleteKeyPairsInput)
	input.KeyPairs = []*string{qc.String(d.Id())}
	var output *qc.DeleteKeyPairsOutput
	simpleRetry(func() error {
		output, err = clt.DeleteKeyPairs(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	return nil
}
