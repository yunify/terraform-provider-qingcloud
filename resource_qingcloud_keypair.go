package qingcloud

import (
	"fmt"
	"log"

	// "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/magicshui/qingcloud-go/keypair"
)

// TODO: attach keypair to instances

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

	// TODO: 这个地方以后需要判断错误
	keypairName := d.Get("keypair_name").(string)
	publicKey := d.Get("public_key").(string)

	// 开始创建 ssh 密钥
	params := keypair.CreateKeyPairRequest{}
	params.KeypairName.Set(keypairName)
	params.PublicKey.Set(publicKey)

	result, err := clt.CreateKeyPair(params)
	if err != nil {
		return fmt.Errorf("Error create Keypair: %s", err)
	}

	// 如果名称没有变
	description := d.Get("description").(string)
	if description != "" {
		modifyAtrributes := keypair.ModifyKeyPairAttributesRequest{}
		modifyAtrributes.Keypair.Set(result.KeypairId)
		modifyAtrributes.Description.Set(description)

		_, err := clt.ModifyKeyPairAttributes(modifyAtrributes)
		if err != nil {
			// 这里可以不用返回错误
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

	resp, err := clt.DescribeKeyPairs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving Keypair: %s", err)
	}
	for _, kp := range resp.KeypairSet {
		if kp.KeypairID == d.Id() {
			d.Set("keypair_name", kp.KeypairName)
			d.Set("instance_ids", kp.InstanceIds)
			return nil
		}
	}
	log.Printf("Unable to find key pair %#v within: %#v", d.Id(), resp.KeypairSet)
	d.SetId("")
	return nil
}

// 如果要删除一个密钥，那么需要看一下这个密钥是否在其他的instance上是否有使用
func resourceQingcluodKeypairDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).keypair

	log.Printf("[DEBUG] Delete the keypair id:%v", d.Id())

	// 是否有主机绑定了密钥
	describeParams := keypair.DescribeKeyPairsRequest{}
	describeParams.KeypairsN.Add(d.Id())
	// TODO: 这个应该是自动的配置过程
	describeParams.Limit.Set(1000)
	describeParams.Verbose.Set(1)
	resp, err := clt.DescribeKeyPairs(describeParams)
	if err != nil {
		return fmt.Errorf("Error retrieving Keypair: %s", err)
	}

	log.Printf("Instance attached count is: %d", resp.TotalCount)

	for _, kp := range resp.KeypairSet {
		if kp.KeypairID == d.Id() {
			detachRequest := keypair.DetachKeyPairsRequest{}
			detachRequest.KeypairsN.Add(d.Id())
			detachRequest.InstancesN.Add(kp.InstanceIds...)

			_, err := clt.DetachKeyPairs(detachRequest)
			if err != nil {
				log.Printf("[ERROR] Detach key pair %s error from instance %s,error is : %s", d.Id(), kp.InstanceIds, err)
				continue
			}
			log.Printf("[DEBUG] Detach key pair %s from instances %s", d.Id(), kp.InstanceIds)
		}
	}

	params := keypair.DeleteKeyPairsRequest{}
	params.KeypairsN.Add(d.Id())
	_, deleteErr := clt.DeleteKeyPairs(params)
	if deleteErr != nil {
		return fmt.Errorf(
			"Error delete keypair: %s", deleteErr)
	}

	return nil
}

func resourceQingcloudKeypairUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).keypair

	if !d.HasChange("description") && !d.HasChange("keypair_name") {
		return nil
	}

	params := keypair.ModifyKeyPairAttributesRequest{}
	if d.HasChange("description") {
		params.Description.Set(d.Get("description").(string))
	}
	if d.HasChange("keypair_name") {
		params.KeypairName.Set(d.Get("keypair_name").(string))
	}
	params.Keypair.Set(d.Id())
	_, err := clt.ModifyKeyPairAttributes(params)
	if err != nil {
		// 这里可以不用返回错误
		return fmt.Errorf("Error modify keypair description: %s", err)
	}
	return nil
	// return resourceQingcloudKeypairRead(d, meta)

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
		"id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}

func deleteKeypairFromInstance(meta interface{}, keypairID string, instanceID ...string) error {
	clt := meta.(*QingCloudClient).keypair
	params := keypair.DetachKeyPairsRequest{}
	params.InstancesN.Add(instanceID...)
	params.KeypairsN.Add(keypairID)
	_, err := clt.DetachKeyPairs(params)

	for _, o := range instanceID {
		_, err := InstanceTransitionStateRefresh(meta.(*QingCloudClient).instance, o)
		if err != nil {
			return err
		}
	}

	return err
}
