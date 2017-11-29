/**
 * Copyright (c) 2016 Magicshui
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

func resourceUpdateKeyPairs(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange("keypair_ids") {
		return nil
	}
	clt := meta.(*QingCloudClient).keypair
	oldV, newV := d.GetChange("keypair_ids")
	var oldKeyPairs []string
	var newKeyPairs []string
	for _, v := range oldV.(*schema.Set).List() {
		oldKeyPairs = append(oldKeyPairs, v.(string))
	}
	for _, v := range newV.(*schema.Set).List() {
		newKeyPairs = append(newKeyPairs, v.(string))
	}
	attachKeyPairs, detachKeyPairs := stringSliceDiff(newKeyPairs, oldKeyPairs)

	if len(detachKeyPairs) > 0 {
		input := new(qc.DetachKeyPairsInput)
		input.Instances = []*string{qc.String(d.Id())}
		for _, keypair := range detachKeyPairs {
			input.KeyPairs = append(input.KeyPairs, &keypair)
		}
		var output *qc.DetachKeyPairsOutput
		var err error
		simpleRetry(func() error {
			output, err = clt.DetachKeyPairs(input)
			return isServerBusy(err)
		})
		if err != nil {
			return err
		}
	}
	if len(attachKeyPairs) > 0 {
		input := new(qc.AttachKeyPairsInput)
		input.Instances = []*string{qc.String(d.Id())}
		for _, keypair := range attachKeyPairs {
			input.KeyPairs = append(input.KeyPairs, &keypair)
		}
		var output *qc.AttachKeyPairsOutput
		var err error
		simpleRetry(func() error {
			output, err = clt.AttachKeyPairs(input)
			return isServerBusy(err)
		})
		if err != nil {
			return err
		}
	}
	return nil
}
