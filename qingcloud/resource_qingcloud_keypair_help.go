/**
 * Copyright (c) 2016 Magicshui
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */
/**
 * Copyright (c) 2017 yunify
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifyKeypairAttributes(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).keypair
	input := new(qc.ModifyKeyPairAttributesInput)
	input.KeyPair = qc.String(d.Id())
	nameUpdate := false
	descriptionUpdate := false
	input.KeyPairName, nameUpdate = getNamePointer(d)
	input.Description, descriptionUpdate = getDescriptionPointer(d)
	if nameUpdate || descriptionUpdate {
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
