package qingcloud

import (
	"fmt"

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
				Description: "密钥名称",
			},
			"public_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Computed: true,
			},
			"tag_names": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceQingcloudKeypairCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).keypair
	input := new(qc.CreateKeyPairInput)
	input.KeyPairName = qc.String(d.Get("name").(string))
	input.Mode = qc.String("user")
	input.PublicKey = qc.String(d.Get("public_key").(string))
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("Error create keypair input validate: %s", err)
	}
	output, err := clt.CreateKeyPair(input)
	if err != nil {
		return fmt.Errorf("Error create keypair: %s", err)
	}
	if qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error create keypair: %s", *output.Message)
	}
	d.SetId(qc.StringValue(output.KeyPairID))
	if err := resourceUpdateTag(d, meta, qingcloudResourceTypeKeypair); err != nil {
		return err
	}
	return modifyKeypairAttributes(d, meta, false)
}

func resourceQingcloudKeypairRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).keypair
	input := new(qc.DescribeKeyPairsInput)
	input.KeyPairs = []*string{qc.String(d.Id())}
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("Error describe keypair input validate: %s", err)
	}
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
	err := modifyKeypairAttributes(d, meta, false)
	if err != nil {
		return err
	}
	if err := resourceUpdateTag(d, meta, qingcloudResourceTypeKeypair); err != nil {
		return err
	}
	return resourceQingcloudKeypairRead(d, meta)
}

// 如果要删除一个密钥，那么需要看一下这个密钥是否在其他的instance上是否有使用
func resourceQingcluodKeypairDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).keypair
	describeKeyPairsInput := new(qc.DescribeKeyPairsInput)
	describeKeyPairsInput.KeyPairs = []*string{qc.String(d.Id())}
	err := describeKeyPairsInput.Validate()
	if err != nil {
		return fmt.Errorf("Error describe keypair input validate: %s", err)
	}
	describeKeyPairsOutput, err := clt.DescribeKeyPairs(describeKeyPairsInput)
	if err != nil {
		return fmt.Errorf("Error describe keypair: %s", err)
	}

	if describeKeyPairsOutput.RetCode != nil && qc.IntValue(describeKeyPairsOutput.RetCode) != 0 {
		return fmt.Errorf("Error describe keypair: %s", *describeKeyPairsOutput.Message)
	}
	if len(describeKeyPairsOutput.KeyPairSet[0].InstanceIDs) > 0 {
		return fmt.Errorf("Error keypair %s is using, can't delete", d.Id())
	}
	input := new(qc.DeleteKeyPairsInput)
	input.KeyPairs = []*string{qc.String(d.Id())}
	err = input.Validate()
	if err != nil {
		return fmt.Errorf("Error delete keypairs input validate: %s", err)
	}
	output, err := clt.DeleteKeyPairs(input)
	if err != nil {
		return fmt.Errorf("Error delete keypairs: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error delete keypairs: %s", *output.Message)
	}
	return nil
}
