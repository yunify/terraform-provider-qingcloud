package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
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
		},
	}
}

func resourceQingcloudTagCreate(d *schema.ResourceData, meta interface{}) error {
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
