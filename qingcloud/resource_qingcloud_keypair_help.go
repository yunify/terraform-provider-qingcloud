package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifyKeypairAttributes(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).keypair
	input := new(qc.ModifyKeyPairAttributesInput)
	input.KeyPair = qc.String(d.Id())
	attributeUpdate := false
	attributeUpdate2 := false
	input.KeyPairName, attributeUpdate = getNamePointer(d)
	input.Description, attributeUpdate2 = getDescriptionPointer(d)
	if attributeUpdate || attributeUpdate2 {
		var output *qc.ModifyKeyPairAttributesOutput
		var err error
		simpleRetry(func() error {
			output, err = clt.ModifyKeyPairAttributes(input)
			return isServerBusy(err)
		})
		if err != nil {
			return err
		}
	}
	return nil
}
