package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/magicshui/qingcloud-go/cache"
)

func resourceQingcloudCache() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudCacheCreate,
		Read:   resourceQingcloudCacheRead,
		Update: resourceQingcloudCacheUpdate,
		Delete: resourceQingcloudCacheDelete,
		Schema: map[string]*schema.Schema{
			"vxnet": &schema.Schema{
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
				ValidateFunc: withinArrayString("redis2.8.17", "memcached1.4.13"),
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "缓存名称",
			},
			"parameter_group": &schema.Schema{
				Type: schema.TypeString,
				Description: "缓存服务配置组ID，如果不指定则为默认配置组。	",
				Optional: true,
			},
			"auto_backup_time": &schema.Schema{
				Type: schema.TypeInt,
				Description: "自动备份时间(UTC 的 Hour 部分)，有效值0-23，任何大于23的整型值均表示关闭自动备份，默认值 99	",
				ValidateFunc: withinArrayIntRange(0, 23),
				Optional:     true,
			},
			"cache_class": &schema.Schema{
				Type: schema.TypeInt,
				Description: "性能型和高性能型缓存服务，性能型：0，高性能型：1	",
				ValidateFunc: withinArrayInt(0, 1),
				Required:     true,
			},
			// "private_ips": &schema.Schema{
			// 	Type: schema.TypeSet,
			// 	Elem: &schema.Resource{},
			// },
		},
	}
}

func resourceQingcloudCacheCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).cahce
	params := cache.CreateCacheRequest{}
	params.Vxnet.Set(d.Get("vxnet").(string))
	params.CacheSize.Set(d.Get("size").(int))
	params.CacheType.Set(d.Get("type").(string))
	params.CacheName.Set(d.Get("name").(string))
	params.CacheParameterGroup.Set(d.Get("parameter_group").(string))
	params.AutoBackupTime.Set(d.Get("auto_backup_time").(int))
	params.CacheClass.Set(d.Get("cache_class").(string))
	resp, err := clt.CreateCache(params)
	if err != nil {
		return err
	}
	d.SetId(resp.CacheId)
	return nil
}

func resourceQingcloudCacheRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).cahce
	params := cache.DescribeCachesRequest{}
	params.CachesN.Add(d.Id())
	_, err := clt.DescribeCaches(params)
	if err != nil {
		return err
	}
	return nil
}

func resourceQingcloudCacheUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceQingcloudCacheDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).cahce
	params := cache.DeleteCachesRequest{}
	params.CachesN.Add(d.Id())
	_, err := clt.DeleteCaches(params)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
