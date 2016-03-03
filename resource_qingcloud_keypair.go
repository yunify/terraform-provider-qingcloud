package qingcloud

import (
	"fmt"

	// "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/magicshui/qingcloud-go/keypair"
)

func resourceQingcloudKeypair() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudKeypairCreate,
		Read:   resourceQingcloudKeypairRead,
		Update: resourceQingcloudKeypairUpdate,
		Delete: resourceQingcluodKeypairDelete,
		Schema: resourceQingCloudKeypairSchema(),
	}
}

func resourceQingcloudKeypairCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).keypair

	params := keypair.CreateKeyPairRequest{}
	params.KeypairName.Set(d.Get("keypair_name").(string))
	params.PublicKey.Set(d.Get("public_key").(string))

	result, err := clt.CreateKeyPair(params)
	if err != nil {
		return fmt.Errorf("Error create Keypair: %s", err)
	}

	if description := d.Get("description").(string); description != "" {
		modifyAtrributes := keypair.ModifyKeyPairAttributesRequest{}
		modifyAtrributes.Keypair.Set(result.KeypairId)
		modifyAtrributes.Description.Set(description)
		_, err := clt.ModifyKeyPairAttributes(modifyAtrributes)
		if err != nil {
			return fmt.Errorf("Error modify keypair description: %s", err)
		}
	}

	d.SetId(result.KeypairId)
	return nil
}

func resourceQingcloudKeypairRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).keypair

	// 设置请求参数
	params := keypair.DescribeKeyPairsRequest{}
	params.KeypairsN.Add(d.Id())
	params.Verbose.Set(1)
	params.Limit.Set(10000)

	resp, err := clt.DescribeKeyPairs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving Keypair: %s", err)
	}
	for _, kp := range resp.KeypairSet {
		if kp.KeypairID == d.Id() {
			d.Set("keypair_name", kp.KeypairName)

			var instanceIDs = make([]string, 0)
			for _, o := range kp.InstanceIds {
				instanceIDs = append(instanceIDs, o)
			}
			d.Set("instance_ids", instanceIDs)
			return nil
		}
	}
	d.SetId("")
	return nil
}

// 如果要删除一个密钥，那么需要看一下这个密钥是否在其他的instance上是否有使用
func resourceQingcluodKeypairDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).keypair
	if err := deleteKeypairFromInstance(meta, d.Id(), d.Get("instance_ids").([]interface{})...); err != nil {
		return fmt.Errorf("Error %s", err)
	}

	params := keypair.DeleteKeyPairsRequest{}
	params.KeypairsN.Add(d.Id())
	_, deleteErr := clt.DeleteKeyPairs(params)
	if deleteErr != nil {
		return fmt.Errorf("Error %s", deleteErr)
	}

	return nil
}

func resourceQingcloudKeypairUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).keypair

	if !d.HasChange("description") && !d.HasChange("keypair_name") {
		return nil
	}
	params := keypair.ModifyKeyPairAttributesRequest{}
	params.Keypair.Set(d.Id())

	if d.HasChange("description") {
		params.Description.Set(d.Get("description").(string))
	}
	if d.HasChange("keypair_name") {
		params.KeypairName.Set(d.Get("keypair_name").(string))
	}

	_, err := clt.ModifyKeyPairAttributes(params)
	if err != nil {
		return fmt.Errorf("Error modify keypair description: %s", err)
	}
	return nil
}

func resourceQingCloudKeypairSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"keypair_name": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"public_key": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"instance_ids": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}

func deleteKeypairFromInstance(meta interface{}, keypairID string, instanceID ...interface{}) error {
	clt := meta.(*QingCloudClient).keypair
	params := keypair.DetachKeyPairsRequest{}
	var instances = make([]string, 0)
	for _, o := range instanceID {
		instances = append(instances, o.(string))
	}

	params.InstancesN.Add(instances...)
	params.KeypairsN.Add(keypairID)
	_, err := clt.DetachKeyPairs(params)

	for _, o := range instances {
		_, err := InstanceTransitionStateRefresh(meta.(*QingCloudClient).instance, o)
		if err != nil {
			return err
		}
	}

	return err
}
