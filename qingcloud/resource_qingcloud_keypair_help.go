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
