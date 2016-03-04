package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/magicshui/qingcloud-go/loadbalancer"
)

func resourceQingcloudLoadbalancerBackend() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudLoadbalancerBackendCreate,
		Read:   resourceQingcloudLoadbalancerBackendRead,
		Update: resourceQingcloudLoadbalancerBackendUpdate,
		Delete: resourceQingcloudLoadbalancerBackendDelete,
		Schema: resourceQingcloudLoadbalancerBackendSchema(),
	}
}

func resourceQingcloudLoadbalancerBackendCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceQingcloudLoadbalancerBackendRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceQingcloudLoadbalancerBackendUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceQingcloudLoadbalancerBackendDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceQingcloudLoadbalancerBackendSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"listener": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"resource_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"policy": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
		},
		"port": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"weight": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
}
