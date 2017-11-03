package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifyTagAttributes(d *schema.ResourceData, meta interface{}, create bool) error {
	clt := meta.(*QingCloudClient).tag
	input := new(qc.ModifyTagAttributesInput)
	input.Tag = qc.String(d.Id())
	if create {
		if !d.HasChange("description") && !d.HasChange("color") {
			return nil
		}
		if d.HasChange("description") {
			input.Description = qc.String(d.Get("description").(string))
		}
		if d.HasChange("color") {
			input.Color = qc.String(d.Get("color").(string))
		}
	} else {
		if !d.HasChange("description") && !d.HasChange("name") && !d.HasChange("color") {
			return nil
		}
		if d.HasChange("description") {
			if d.Get("description").(string) == "" {
				input.Description = qc.String(" ")
			} else {
				input.Description = qc.String(d.Get("description").(string))
			}
		}
		if d.HasChange("name") {
			input.TagName = qc.String(d.Get("name").(string))
		}
		if d.HasChange("color") {
			input.Color = qc.String(d.Get("color").(string))
		}
	}
	output, err := clt.ModifyTagAttributes(input)
	if err != nil {
		return fmt.Errorf("Error modify tag attributes: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error modify tag attributes: %s", *output.Message)
	}
	return nil
}

func resourceSetTag(d *schema.ResourceData, tags []*qc.Tag) {
	tagIDs := make([]string, 0, len(tags))
	tagNames := make([]string, 0, len(tags))
	for _, tag := range tags {
		tagIDs = append(tagIDs, qc.StringValue(tag.TagID))
		tagNames = append(tagNames, qc.StringValue(tag.TagName))
	}
	d.Set("tag_ids", tagIDs)
	d.Set("tag_names", tagNames)
}

func resourceUpdateTag(d *schema.ResourceData, meta interface{}, resourceType string) error {
	if !d.HasChange("tag_ids") {
		return nil
	}
	clt := meta.(*QingCloudClient).tag
	oldV, newV := d.GetChange("tag_ids")
	var oldTags []string
	var newTags []string
	for _, v := range oldV.(*schema.Set).List() {
		oldTags = append(oldTags, v.(string))
	}
	for _, v := range newV.(*schema.Set).List() {
		newTags = append(newTags, v.(string))
	}
	attachTags, detachTags := stringSliceDiff(newTags, oldTags)

	if len(detachTags) > 0 {
		input := new(qc.DetachTagsInput)
		for _, tag := range detachTags {
			rtp := qc.ResourceTagPair{
				TagID:        qc.String(tag),
				ResourceID:   qc.String(d.Id()),
				ResourceType: qc.String(resourceType),
			}
			input.ResourceTagPairs = append(input.ResourceTagPairs, &rtp)
		}
		_, err := clt.DetachTags(input)
		if err != nil {
			return fmt.Errorf("Error detach tag: %s", err)
		}
	}
	if len(attachTags) > 0 {
		input := new(qc.AttachTagsInput)
		for _, tag := range attachTags {
			rtp := qc.ResourceTagPair{
				TagID:        qc.String(tag),
				ResourceID:   qc.String(d.Id()),
				ResourceType: qc.String(resourceType),
			}
			input.ResourceTagPairs = append(input.ResourceTagPairs, &rtp)
		}
		_, err := clt.AttachTags(input)
		if err != nil {
			return fmt.Errorf("Error detach tag: %s", err)
		}
	}
	return nil
}
