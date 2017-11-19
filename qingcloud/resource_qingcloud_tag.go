package qingcloud

import (
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
	input.TagName, _ = getNamePointer(d)
	var output *qc.CreateTagOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.CreateTag(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
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
	simpleRetry(func() error {
		output, err = clt.DescribeTags(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
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
	if err := modifyTagAttributes(d, meta); err != nil {
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
	simpleRetry(func() error {
		output, err = clt.DeleteTags(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
