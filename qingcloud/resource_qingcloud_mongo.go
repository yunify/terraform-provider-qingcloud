package qingcloud

// import (
// 	"github.com/hashicorp/terraform/helper/schema"
// 	"github.com/magicshui/qingcloud-go/mongo"
// )

// func resourceQingcloudMongo() *schema.Resource {
// 	return &schema.Resource{
// 		Create: resourceQingcloudMongoCreate,
// 		Read:   resourceQingcloudMongoRead,
// 		Update: resourceQingcloudMongoUpdate,
// 		Delete: resourceQingcloudMongoDelete,
// 		Schema: map[string]*schema.Schema{
// 			"vxnet": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Required: true,
// 				Description: "私有网络 ID	",
// 			},
// 			"version": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Default:  "3.0",
// 			},
// 			"type": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Required: true,
// 				Description: "Mongo 配置型号，1 – 1核2G，2 – 2核4G，3 – 4核8G，4 – 8核16G，5 – 8核32G	",
// 				ValidateFunc: withinArrayInt(1, 2, 3, 4, 5),
// 			},
// 			"size": &schema.Schema{
// 				Type:     schema.TypeInt,
// 				Required: true,
// 				Description: "Mongo 存储容量(GB)，用于存放数据和日志，最小10G，最大1000G	",
// 			},
// 			resourceName: &schema.Schema{
// 				Type:        schema.TypeString,
// 				Optional:    true,
// 				Description: "名称",
// 			},
// 			resourceDescription: &schema.Schema{
// 				Type:     schema.TypeString,
// 				Optional: true,
// 			},
// 			"auto_backup_time": &schema.Schema{
// 				Type: schema.TypeInt,
// 				Description: "自动备份时间(UTC 的 Hour 部分)，有效值0-23，任何大于23的整型值均表示关闭自动备份，默认值 99	",
// 				ValidateFunc: withinArrayIntRange(0, 23),
// 				Optional:     true,
// 			},

// 			// "private_ips": &schema.Schema{
// 			// 	Type: schema.TypeSet,
// 			// 	Elem: &schema.Resource{},
// 			// },
// 		},
// 	}
// }

// func resourceQingcloudMongoCreate(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).mongo
// 	params := mongo.CreateMongoRequest{}
// 	params.Vxnet.Set(d.Get("vxnet").(string))
// 	params.MongoVersion.Set(d.Get("version").(string))
// 	params.MongoName.Set(d.Get(resourceName).(string))
// 	params.MongoType.Set(d.Get("type").(int))
// 	params.StorageSize.Set(d.Get("size").(int))
// 	params.Description.Set(d.Get(resourceDescription).(string))
// 	params.AutoBackupTime.Set(d.Get("auto_backup_time").(int))

// 	resp, err := clt.CreateMongo(params)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId(resp.Mongo)
// 	return nil
// }

// func resourceQingcloudMongoRead(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).mongo
// 	params := mongo.DescribeMongosRequest{}
// 	params.MongosN.Add(d.Id())
// 	_, err := clt.DescribeMongos(params)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func resourceQingcloudMongoUpdate(d *schema.ResourceData, meta interface{}) error {
// 	return nil
// }

// func resourceQingcloudMongoDelete(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).mongo
// 	params := mongo.DeleteMongosRequest{}
// 	params.MongosN.Add(d.Id())
// 	_, err := clt.DeleteMongos(params)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId("")
// 	return nil
// }
