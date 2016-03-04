package qingcloud

import (
	"fmt"

	// "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/magicshui/qingcloud-go/keypair"
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
	params := keypair.ModifyKeyPairAttributesRequest{}
	params.Keypair.Set(d.Id())

	if create {
		if description := d.Get("description").(string); description != "" {
			params.Description.Set(description)
		}
	} else {
		if d.HasChange("description") {
			params.Description.Set(d.Get("description").(string))
		}
		if d.HasChange("keypair_name") {
			params.KeypairName.Set(d.Get("keypair_name").(string))
		}
	}
	_, err := clt.ModifyKeyPairAttributes(params)
	if err != nil {
		return fmt.Errorf("Error modify keypair description: %s", err)
	}
	return nil
}
