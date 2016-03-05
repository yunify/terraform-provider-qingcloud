package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/magicshui/qingcloud-go/cache"
)

func resourceQingcloudCacheParameterGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudCacheParameterGroupCreate,
		Read:   resourceQingcloudCacheParameterGroupRead,
		Update: resourceQingcloudCacheParameterGroupUpdate,
		Delete: resourceQingcloudCacheParameterGroupDelete,
		Schema: nil, // map[string]*schema.Schema{}
	}
}

func resourceQingcloudCacheParameterGroupCreate(d *schema.ResourceData, meta interface{}) error {
	// clt := meta.(*QingCloudClient).cahce
	// params:=cache.CreateCacheParameterGroupRequest{}
	// params.CacheParameterGroupName
	return nil
}

func resourceQingcloudCacheParameterGroupRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceQingcloudCacheParameterGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceQingcloudCacheParameterGroupDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
