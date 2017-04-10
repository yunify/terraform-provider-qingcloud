package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/magicshui/qingcloud-go/cache"
	qc "github.com/lowstz/qingcloud-sdk-go/service"
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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceQingcloudCacheParameterGroupCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).cache
	input := new(qc.CreateCacheParameterGroupInput)
	input.CacheParameterGroupName = qc.String(d.Get("name").(string))
	input.CacheType = qc.String(d.Get("type").(string))
	output, err := clt.CreateCacheParameterGroup(input)
	if err != nil {
		return err
	}
	d.SetId(qc.StringValue(output.CacheParameterGroupID))
	if err := modifyCacheParameterGroupAttributes(d, meta, true); err != nil {
		return err
	}
	if err := cacheParameterGroupSetPassword(d, meta, true); err != nil {
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
	group := output.CacheParameterGroupSet[0]
	d.Set("type", qc.StringValue(group.CacheType))
	d.Set("name", qc.StringValue(group.CacheParameterGroupName))
	d.Set("description", qc.StringValue(group.Description))

	describeCacheParameterInput := new(qc.DescribeCacheParametersInput)
	describeCacheParameterInput.CacheParameterGroup = qc.String(d.Id())
	describeCacheParameterOutput, err := clt.DescribeCacheParameters(describeCacheParameterInput)
	if err != nil {
		return err
	}
	for _, v := range describeCacheParameterOutput.CacheParameterSet {
		if qc.StringValue(v.CacheParameterName) == "requirepass" {
			d.Set("password", qc.StringValue(v.CacheParameterValue))
		}
	}
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
