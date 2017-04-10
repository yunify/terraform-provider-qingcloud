package qingcloud

import (

	// "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	qc "github.com/lowstz/qingcloud-sdk-go/service"
)

// func deleteKeypairFromInstance(meta interface{}, keypairID string, instanceID ...interface{}) error {
// 	clt := meta.(*QingCloudClient).keypair
// 	params := keypair.DetachKeyPairsRequest{}
// 	var instances = make([]string, 0)
// 	for _, o := range instanceID {
// 		instances = append(instances, o.(string))
// 	}

// 	params.InstancesN.Add(instances...)
// 	params.KeypairsN.Add(keypairID)
// 	_, err := clt.DetachKeyPairs(params)

// 	for _, o := range instances {
// 		_, err := InstanceTransitionStateRefresh(meta.(*QingCloudClient).instance, o)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return err
// }

func modifyKeypairAttributes(d *schema.ResourceData, meta interface{}, create bool) error {
	clt := meta.(*QingCloudClient).keypair
	input := new(qc.ModifyKeyPairAttributesInput)
	input.KeyPair = qc.String(d.Id())

	if create {
		if description := d.Get("description").(string); description == "" {
			return nil
		}
		input.Description = qc.String(d.Get("description").(string))
	} else {
		if !d.HasChange("description") && !d.HasChange("name") {
			return nil
		}
		if d.HasChange("description") {
			input.Description = qc.String(d.Get("description").(string))
		}
		if d.HasChange("name") {
			input.KeyPairName = qc.String(d.Get("name").(string))
		}
	}
	_, err := clt.ModifyKeyPairAttributes(input)
	return err
}
