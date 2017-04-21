package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/lowstz/qingcloud-sdk-go/service"
)

func resourceQingcloudTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudTagCreate,
		Read:   resourceQingcloudTagRead,
		Update: resourceQingcloudTagUpdate,
		Delete: resourceQingcloudTagDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"color": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceQingcloudTagCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).tag
	input := new(qc.CreateTagInput)
	input.TagName = qc.String(d.Get("name").(string))
	// input.Color = qc.String(d.Get("color").(string))
	output, err := clt.CreateTag(input)
	if err != nil {
		return fmt.Errorf("Error create tag: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error create tag: %s", *output.Message)
	}
	d.SetId(qc.StringValue(output.TagID))
	err = modifyTagAttributes(d, meta, true)
	if err != nil {
		return err
	}
	return nil
}
func resourceQingcloudTagRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).tag
	input := new(qc.DescribeTagsInput)
	input.Tags = []*string{qc.String(d.Id())}
	output, err := clt.DescribeTags(input)
	if err != nil {
		return fmt.Errorf("Error describe tag: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error create tag: %s", *output.Message)
	}
	if len(output.TagSet) == 0 {
		return fmt.Errorf("Error tag not found")
	}
	tag := output.TagSet[0]
	d.Set("name", tag.TagName)
	d.Set("description", tag.Description)
	if qc.StringValue(tag.Color) != "default" {
		d.Set("color", tag.Color)
	}
	return nil
}
func resourceQingcloudTagUpdate(d *schema.ResourceData, meta interface{}) error {
	err := modifyTagAttributes(d, meta, false)
	if err != nil {
		return err
	}
	return nil
}
func resourceQingcloudTagDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).tag
	input := new(qc.DeleteTagsInput)
	input.Tags = []*string{qc.String(d.Id())}
	output, err := clt.DeleteTags(input)
	if err != nil {
		return fmt.Errorf("Error delete tag: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error delete tag: %s", *output.Message)
	}
	d.SetId("")
	return nil
}
