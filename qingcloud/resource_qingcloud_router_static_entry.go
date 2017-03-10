package qingcloud

// import (
// 	"github.com/hashicorp/terraform/helper/schema"
// 	"github.com/magicshui/qingcloud-go/router"
// )

// func resourceQingcloudRouterStaticEntry() *schema.Resource {
// 	return &schema.Resource{
// 		Create: nil,
// 		Read:   nil,
// 		Update: nil,
// 		Delete: nil,
// 		Schema: map[string]*schema.Schema{
// 			"router_static": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Required: true,
// 				Description: "需要增加条目的路由器规则ID	",
// 			},
// 			"val1": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Required: true,
// 				Description: `PPTP 账户信息：val1 表示账户名
// 					三层 GRE 隧道：val1 表示目标网络
// 					三层 IPsec 隧道：val1 表示本地网络 (val2 可为空)`,
// 			},
// 			"val2": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Required: true,
// 				Description: `PPTP 账户信息：val2 表示密码
// 					三层 IPsec 隧道：val2 表示目标网络 (val1 可为空)`,
// 			},
// 		},
// 	}
// }

// func resourceQingcloudRouterStaticEntryCreate(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).router
// 	params := router.AddRouterStaticEntriesRequest{}
// 	params.RouterStatic.Set(d.Get("router_static").(string))
// 	params.EntriesNVal1.Add(d.Get("val1").(string))
// 	params.EntriesNVal2.Add(d.Get("val2").(string))
// 	resp, err := clt.AddRouterStaticEntries(params)
// 	if err != nil {
// 		return err
// 	}
// 	rs := resp.RouterStaticEntries[0]
// 	d.SetId(rs)
// 	return nil
// }

// func resourceQingcloudRouterStaticEntryRead(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).router
// 	params := router.DescribeRouterStaticEntriesRequest{}
// 	params.RouterStaticEntryId.Set(d.Id())
// 	resp, err := clt.DescribeRouterStaticEntries(params)
// 	if err != nil {
// 		return err
// 	}
// 	if len(resp.RouterStaticEntrySet) == 0 {
// 		d.SetId("")
// 		return nil
// 	}

// 	rs := resp.RouterStaticEntrySet[0]
// 	d.Set("router_static", rs.RouterStaticID)
// 	d.Set("val1", rs.Val1)
// 	d.Set("val2", rs.Val2)
// 	return nil
// }

// func resourceQingcloudRouterStaticEntryDelete(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).router
// 	params := router.DeleteRouterStaticEntriesReqeust{}
// 	params.RouterStaticEntriesN.Add(d.Id())
// 	_, err := clt.DeleteRouterStaticEntries(params)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId("")
// 	return nil
// }

// func resourceQingcloudRouterStaticEntryUpdate(d *schema.ResourceData, meta interface{}) error {
// 	return nil
// }
