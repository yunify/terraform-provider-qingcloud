package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	qc "github.com/lowstz/qingcloud-sdk-go/service"
)

func resourceQingcloudCache() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudCacheCreate,
		Read:   resourceQingcloudCacheRead,
		Update: resourceQingcloudCacheUpdate,
		Delete: resourceQingcloudCacheDelete,
		Schema: map[string]*schema.Schema{
			"vxnet_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "缓存服务运行的私有网络ID",
			},
			"size": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "缓存服务节点内存大小，单位 GB",
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: withinArrayString("redis3.0.5", "redis2.8.17", "memcached1.4.13"),
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "缓存名称",
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"cache_parameter_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "缓存服务配置组ID，如果不指定则为默认配置组.",
				Optional:    true,
			},
			"auto_backup_time": &schema.Schema{
				Type: schema.TypeInt,
				Description: "自动备份时间(UTC 的 Hour 部分)，有效值0-23，任何大于23的整型值均表示关闭自动备份，默认值 99	",
				ValidateFunc: withinArrayIntRange(0, 99),
				Optional:     true,
			},
			"cache_class": &schema.Schema{
				Type: schema.TypeInt,
				Description: "性能型和高性能型缓存服务，性能型：0，高性能型：1	",
				ValidateFunc: withinArrayInt(0, 1),
				Required:     true,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceQingcloudCacheCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).cache
	input := new(qc.CreateCacheInput)
	input.VxNet = qc.String(d.Get("vxnet_id").(string))
	input.CacheSize = qc.Int(d.Get("size").(int))
	input.CacheType = qc.String(d.Get("type").(string))
	input.CacheName = qc.String(d.Get("name").(string))
	input.CacheParameterGroup = qc.String(d.Get("cache_parameter_group_id").(string))
	input.AutoBackupTime = qc.Int(d.Get("auto_backup_time").(int))
	input.CacheClass = qc.Int(d.Get("cache_class").(int))
	output, err := clt.CreateCache(input)
	if err != nil {
		return err
	}
	d.SetId(qc.StringValue(output.CacheID))
	return resourceQingcloudCacheRead(d, meta)
}

func resourceQingcloudCacheRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).cache
	input := new(qc.DescribeCachesInput)
	input.Caches = []*string{qc.String(d.Id())}
	input.Verbose = qc.Int(1)
	output, err := clt.DescribeCaches(input)
	if err != nil {
		return err
	}
	cache := output.CacheSet[0]
	d.Set("vxnet_id", qc.StringValue(cache.VxNet.VxNetID))
	d.Set("size", qc.IntValue(cache.CacheSize))
	d.Set("type", qc.StringValue(cache.CacheType))
	d.Set("name", qc.StringValue(cache.CacheName))
	d.Set("description", qc.StringValue(cache.Description))
	d.Set("cache_parameter_group_id", qc.StringValue(cache.CacheParameterGroupID))
	d.Set("auto_backup_time", qc.IntValue(cache.AutoBackupTime))
	d.Set("cache_class", qc.IntValue(cache.CacheClass))
	for _, v := range cache.Nodes {
		if qc.StringValue(v.CacheRole) == "master" {
			d.Set("private_ip", qc.StringValue(v.PrivateIP))
			break
		}
	}
	return nil
}

func resourceQingcloudCacheUpdate(d *schema.ResourceData, meta interface{}) error {
	// modify cache attributes
	if err := modifyCacheAttributes(d, meta, false); err != nil {
		return err
	}
	// resize cache
	if d.HasChange("size") {
		oldV, newV := d.GetChange("size")
		if newV.(int) <= oldV.(int) {
			return fmt.Errorf("newsize %d must bigger than oldsize %d", newV.(int), oldV.(int))
		}
		if err := resizeCache(d, meta); err != nil {
			return err
		}
	}
	return resourceQingcloudCacheRead(d, meta)
}

func resourceQingcloudCacheDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).cache
	if _, err := CacheTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	input := new(qc.DeleteCachesInput)
	input.Caches = []*string{qc.String(d.Id())}
	if _, err := clt.DeleteCaches(input); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
