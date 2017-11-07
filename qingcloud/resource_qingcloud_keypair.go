package qingcloud

import (
	"fmt"
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
			"tag_ids": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "tag ids , keypair wants to use ",
			},
			"tag_names": &schema.Schema{
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "compute by tag ids ",
			},
		},
	}
}

func resourceQingcloudKeypairCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).keypair
	input := new(qc.CreateKeyPairInput)
	if d.Get("name").(string) != "" {
		input.KeyPairName = qc.String(d.Get("name").(string))
	} else {
		input.KeyPairName = nil
	}
	input.Mode = qc.String("user")
	input.PublicKey = qc.String(d.Get("public_key").(string))
	output, err := clt.CreateKeyPair(input)
	if err != nil {
		return fmt.Errorf("Error create keypair: %s", err)
	}
	if qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error create keypair: %s", *output.Message)
	}
	d.SetId(qc.StringValue(output.KeyPairID))

	return resourceQingcloudKeypairUpdate(d, meta)
}

func resourceQingcloudKeypairRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).keypair
	input := new(qc.DescribeKeyPairsInput)
	input.KeyPairs = []*string{qc.String(d.Id())}
	output, err := clt.DescribeKeyPairs(input)
	if err != nil {
		return fmt.Errorf("Error describe keypair: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error describe keypair: %s", *output.Message)
	}
	kp := output.KeyPairSet[0]
	d.Set("name", qc.StringValue(kp.KeyPairName))
	d.Set("description", qc.StringValue(kp.Description))
	resourceSetTag(d, kp.Tags)
	return nil
}

func resourceQingcloudKeypairUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	err := modifyKeypairAttributes(d, meta)
	if err != nil {
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
	describeKeyPairsOutput, err := clt.DescribeKeyPairs(describeKeyPairsInput)
	if err != nil {
		return fmt.Errorf("Error describe keypair: %s", err)
	}
	if describeKeyPairsOutput.RetCode != nil && qc.IntValue(describeKeyPairsOutput.RetCode) != 0 {
		return fmt.Errorf("Error describe keypair: %s", *describeKeyPairsOutput.Message)
	}
	input := new(qc.DeleteKeyPairsInput)
	input.KeyPairs = []*string{qc.String(d.Id())}
	output, err := clt.DeleteKeyPairs(input)
	if err != nil {
		return fmt.Errorf("Error delete keypairs: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error delete keypairs: %s", *output.Message)
	}
	return nil
}
