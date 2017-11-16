package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
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

func modifyKeypairAttributes(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).keypair
	input := new(qc.ModifyKeyPairAttributesInput)
	input.KeyPair = qc.String(d.Id())
	attributeUpdate := false
	if d.HasChange("description") {
		if d.Get("description").(string) == "" {
			input.Description = qc.String(" ")
		} else {
			input.Description = qc.String(d.Get("description").(string))
		}
		attributeUpdate = true
	}
	if d.HasChange("name") && !d.IsNewResource() {
		if d.Get("name").(string) == "" {
			input.KeyPairName = qc.String(" ")
		} else {
			input.KeyPairName = qc.String(d.Get("name").(string))
		}
		attributeUpdate = true
	}
	if attributeUpdate {
		var output *qc.ModifyKeyPairAttributesOutput
		var err error
		retryServerBusy(func() (s *int, err error) {
			output, err = clt.ModifyKeyPairAttributes(input)
			return output.RetCode, err
		})
		if err := getQingCloudErr("modify keypair attributes", output.RetCode, output.Message, err); err != nil {
			return err
		}
	}
	return nil
}
