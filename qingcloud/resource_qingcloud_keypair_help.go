package qingcloud

import (
	"fmt"

	// "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/magicshui/qingcloud-go/keypair"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

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

func modifyKeypairAttributes(d *schema.ResourceData, meta interface{}, create bool) error {
	clt := meta.(*QingCloudClient).keypair
	input := new(qc.ModifyKeyPairAttributesInput)
	input.KeyPair = qc.String(d.Id())

	if create {
		if description := d.Get("description").(string); description != "" {
			input.Description = qc.String(d.Get("description").(string))
			params.Description.Set(description)
		}
	} else {
		if d.HasChange("description") {
			input.Description = qc.String(d.Get("description").(string))
		}
		if d.HasChange("keypair_name") {
			input.KeyPairName = qc.String(d.Get("keypair_name").(string))
		}
	}
	output, err := clt.ModifyKeyPairAttributes(input)
	if err != nil {
		return fmt.Errorf("Error modify keypair attributes: %s", err)
	}
	if output.RetCode != 0 {
		return fmt.Errorf("Error modify keypair attributes: %s", output.Message)
	}
	return nil
}
