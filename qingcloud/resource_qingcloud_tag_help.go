package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifyTagAttributes(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).tag
	input := new(qc.ModifyTagAttributesInput)
	input.Tag = qc.String(d.Id())
	attributeUpdate := false
	if d.HasChange("color") {
		input.Color = qc.String(d.Get("color").(string))
		attributeUpdate = true
	}
	if d.HasChange("description") {
		if d.Get("description").(string) == "" {
			input.Description = qc.String(" ")
		} else {
			input.Description = qc.String(d.Get("description").(string))
		}
		attributeUpdate = true
	}
	if d.HasChange("name") && !d.IsNewResource() {
		input.TagName = qc.String(d.Get("name").(string))
		attributeUpdate = true
	}
	if attributeUpdate {
		var output *qc.ModifyTagAttributesOutput
		var err error
		retryServerBusy(func() (s *int, err error) {
			output, err = clt.ModifyTagAttributes(input)
			return output.RetCode, err
		})
		if err := getQingCloudErr("modify tag attributes", output.RetCode, output.Message, err); err != nil {
			return err
		}
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
		var output *qc.DetachTagsOutput
		var err error
		retryServerBusy(func() (s *int, err error) {
			output, err = clt.DetachTags(input)
			return output.RetCode, err
		})
		if err := getQingCloudErr("detach tag", output.RetCode, output.Message, err); err != nil {
			return err
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
		var output *qc.AttachTagsOutput
		var err error
		retryServerBusy(func() (s *int, err error) {
			output, err = clt.AttachTags(input)
			return output.RetCode, err
		})
		if err := getQingCloudErr("attach tag", output.RetCode, output.Message, err); err != nil {
			return err
		}
	}
	return nil
}
