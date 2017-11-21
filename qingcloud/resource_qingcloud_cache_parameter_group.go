package qingcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func resourceQingcloudCacheParameterGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudCacheParameterGroupCreate,
		Read:   resourceQingcloudCacheParameterGroupRead,
		Update: resourceQingcloudCacheParameterGroupUpdate,
		Delete: resourceQingcloudCacheParameterGroupDelete,
		Schema: map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: withinArrayString("redis3.0.5", "redis2.8.17", "memcached1.4.13"),
			},
			resourceName: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceDescription: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceQingcloudCacheParameterGroupCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).cache
	input := new(qc.CreateCacheParameterGroupInput)
	input.CacheParameterGroupName = qc.String(d.Get(resourceName).(string))
	input.CacheType = qc.String(d.Get("type").(string))
	output, err := clt.CreateCacheParameterGroup(input)
	if err != nil {
		return err
	}
	d.SetId(qc.StringValue(output.CacheParameterGroupID))
	if err := modifyCacheParameterGroupAttributes(d, meta, true); err != nil {
		return err
	}
	return resourceQingcloudCacheParameterGroupRead(d, meta)
}

func resourceQingcloudCacheParameterGroupRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).cache
	input := new(qc.DescribeCacheParameterGroupsInput)
	input.CacheParameterGroups = []*string{qc.String(d.Id())}
	output, err := clt.DescribeCacheParameterGroups(input)
	if err != nil {
		return err
	}
	if *output.RetCode != 0 {
		return fmt.Errorf("Error describe cache: %s ", *output.Message)
	}
	if len(output.CacheParameterGroupSet) == 0 {
		d.SetId("")
		return nil
	}
	group := output.CacheParameterGroupSet[0]
	d.Set("type", qc.StringValue(group.CacheType))
	d.Set(resourceName, qc.StringValue(group.CacheParameterGroupName))
	d.Set(resourceDescription, qc.StringValue(group.Description))
	return nil
}

func resourceQingcloudCacheParameterGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := modifyCacheParameterGroupAttributes(d, meta, false); err != nil {
		return err
	}
	if err := cacheParameterGroupSetPassword(d, meta, false); err != nil {
		return err
	}
	if err := applyCacheParameterGroup(d, meta); err != nil {
		return err
	}
	return resourceQingcloudCacheParameterGroupRead(d, meta)
}

func resourceQingcloudCacheParameterGroupDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).cache
	input := new(qc.DeleteCacheParameterGroupsInput)
	input.CacheParameterGroups = []*string{qc.String(d.Id())}
	if _, err := clt.DeleteCacheParameterGroups(input); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
