package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
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
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "#9f9bb7",
				ValidateFunc: validateColorString,
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
	var output *qc.CreateTagOutput
	var err error
	retryServerBusy(func() (s *int, err error) {
		output, err = clt.CreateTag(input)
		return output.RetCode, err
	})
	if err != nil {
		return fmt.Errorf("Error create tag: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error create tag: %s", *output.Message)
	}
	d.SetId(qc.StringValue(output.TagID))
	return resourceQingcloudTagUpdate(d, meta)
}
func resourceQingcloudTagRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).tag
	input := new(qc.DescribeTagsInput)
	input.Tags = []*string{qc.String(d.Id())}
	var output *qc.DescribeTagsOutput
	var err error
	retryServerBusy(func() (s *int, err error) {
		output, err = clt.DescribeTags(input)
		return output.RetCode, err
	})
	if err != nil {
		return fmt.Errorf("Error describe tag: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error create tag: %s", *output.Message)
	}
	if len(output.TagSet) == 0 {
		d.SetId("")
		return nil
	}
	tag := output.TagSet[0]
	d.Set("name", qc.StringValue(tag.TagName))
	d.Set("description", qc.StringValue(tag.Description))
	d.Set("color", qc.StringValue(tag.Color))
	return nil
}
func resourceQingcloudTagUpdate(d *schema.ResourceData, meta interface{}) error {
	err := modifyTagAttributes(d, meta)
	if err != nil {
		return err
	}
	return nil
}
func resourceQingcloudTagDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).tag
	input := new(qc.DeleteTagsInput)
	input.Tags = []*string{qc.String(d.Id())}
	var output *qc.DeleteTagsOutput
	var err error
	retryServerBusy(func() (s *int, err error) {
		output, err = clt.DeleteTags(input)
		return output.RetCode, err
	})
	if err != nil {
		return fmt.Errorf("Error delete tag: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error delete tag: %s", *output.Message)
	}
	d.SetId("")
	return nil
}
