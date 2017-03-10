package qingcloud

import "github.com/hashicorp/terraform/helper/schema"

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
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceQingcloudTagCreate(d *schema.ResourceData, meta interface{}) error {
	// clt := meta.(*QingCloudClient).tag
	// input := new(qc.CreateTagInput)
	// input.TagName = d.get("name")
	// input.
	// err := input.Validate()
	return nil
}
func resourceQingcloudTagRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceQingcloudTagUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceQingcloudTagDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
