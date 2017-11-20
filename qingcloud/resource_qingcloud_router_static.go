package qingcloud

// import (
// 	"github.com/hashicorp/terraform/helper/schema"
// 	"github.com/magicshui/qingcloud-go/router"
// )

// func resourceQingcloudRouterStatic() *schema.Resource {
// 	return &schema.Resource{
// 		Create: resourceQingcloudRouterStaticCreate,
// 		Read:   resourceQingcloudRouterStaticRead,
// 		Update: resourceQingcloudRouterStaticUpdate,
// 		Delete: resourceQingcloudRouterStaticDelete,
// 		Schema: map[string]*schema.Schema{
// 			resourceName: &schema.Schema{
// 				Type:     schema.TypeString,
// 				Required: true,
// 			},
// 			"router": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Required: true,
// 				Description: `支持的规则类型有：
// 					static_type=1：端口转发规则
// 					static_type=2：VPN 规则
// 					static_type=3：DHCP 选项
// 					static_type=4：二层 GRE 隧道
// 					static_type=5：过滤控制
// 					static_type=6：三层 GRE 隧道
// 					static_type=7：三层 IPsec 隧道
// 					static_type=8：私网DNS`,
// 				ValidateFunc: withinArrayIntRange(1, 8),
// 			},
// 			"type": &schema.Schema{
// 				Type:     schema.TypeInt,
// 				Required: true,
// 				ForceNew: true,
// 			},
// 			"val1": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Required: true,
// 				Description: `端口转发：val1 表示源端口。
// 					VPN：val1 表示 VPN 类型，目前支持 “openvpn” 和 “pptp”，默认值为 “openvpn”。
// 					DHCP 选项：val1 表示 DHCP 主机ID。
// 					二层 GRE 隧道：val1 表示二层隧道的远端 IP 和密钥，如：gre|1.2.3.4|888。
// 					过滤控制：val1 表示『源 IP』
// 					三层 GRE 隧道：val1 表示远端 IP 、密钥、本地点对点IP、对端点对点IP，格式如：6.6.6.6|key|1.2.3.4|4.3.2.1。
// 					三层 IPsec 隧道：val1 表示远端IP（支持接受任意对端，可填 0.0.0.0） 、加密算法(phase2alg&ike，可为空，默认aes)、密钥和远端设备ID（支持接受任意对端设备ID，可填 %any），格式如：1.2.3.4||passw0rd|device-id
// 					私网DNS：val1 表示私网域名，比如node1`,
// 			},
// 			"val2": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Description: `端口转发规则：val2 表示目标 IP 。
// 					OpenVPN 规则：val2 表示 VPN 服务端口号，默认为1194。
// 					PPTP VPN 规则：val2 表示用户名和密码，格式为 user:password
// 					DHCP 选项：val2 表示 DHCP 配置内容，格式为key1=value1;key2=value2，例如：”domain-name-servers=8.8.8.8;fixed-address=192.168.1.2”。
// 					过滤控制：val2 表示『源端口』
// 					三层 GRE 隧道：val2 表示目标网络，多个网络间以 “|” 分隔。注意目标网络不能和路由器已有的私有网络重复。
// 					三层 IPsec 隧道：val2 表示本地网络，多个网络间以 “|” 分隔。
// 					私网DNS：val2 表示IP地址，格式为ip1;ip2，例如：”192.168.1.2;192.168.1.3”`,
// 			},
// 			"val3": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Description: `端口转发规则：val3 表示目标端口号。
// 					OpenVPN 规则：val3 表示 VPN 协议，默认为 “udp”。
// 					PPTP VPN 规则：val3 表示最大连接数，连接数范围是 1-253
// 					过滤控制：val3 表示『目标 IP』
// 					三层 IPsec 隧道：val3 表示目标网络，多个网络间以 “|” 分隔。`,
// 			},
// 			"val4": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Description: `端口转发规则：val4 表示端口转发协议，默认为 “tcp” ，目前支持 “tcp” 和 “udp” 两种协议。
// 					VPN 规则(包括 OpenVPN 和 PPTP)：val4 表示 VPN 客户端的网络地址段，目前支持10.255.x.0/24，x的范围是[0-255]，默认为自动分配。
// 					过滤控制：val4 表示『目标端口』
// 					三层 IPsec 隧道：val4 表示IPsec隧道模式，默认为”main”，支持 主模式（main） 和 野蛮模式（aggrmode）。`,
// 			},
// 			"val5": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Description: `VPN 规则(OpenVPN)：val5 表示 OpenVPN 的验证方式，目前支持 1: 证书验证, 2: 用户名/密码验证, 3: 证书+用户名/密码验证，默认为 “证书验证” 方式。
// 					过滤控制：val5 表示『行为』，包括： “accept” 和 “drop”`,
// 			},
// 		},
// 	}
// }

// func resourceQingcloudRouterStaticCreate(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).router
// 	// 确保没有在更新
// 	if _, err := RouterTransitionStateRefresh(clt, d.Get("router").(string)); err != nil {
// 		return err
// 	}

// 	params := router.AddRouterStaticsRequest{}
// 	params.Router.Set(d.Get("router").(string))
// 	params.StaticsNRouterStaticName.Add(d.Get(resourceName).(string))
// 	params.StaticsNStaticType.Add(int64(d.Get("type").(int)))
// 	params.StaticsNVal1.Add(d.Get("val1").(string))
// 	params.StaticsNVal2.Add(d.Get("val2").(string))
// 	params.StaticsNVal3.Add(d.Get("val3").(string))
// 	params.StaticsNVal4.Add(d.Get("val4").(string))
// 	params.StaticsNVal5.Add(d.Get("val5").(string))
// 	resp, err := clt.AddRouterStatics(params)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId(resp.RouterStatics[0])

// 	return applyRouterUpdates(meta, d.Get("router").(string))
// }
// func resourceQingcloudRouterStaticRead(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).router
// 	params := router.DescribeRouterStaticsRequest{}
// 	params.RouterStaticsN.Add(d.Id())
// 	resp, err := clt.DescribeRouterStatics(params)
// 	if err != nil {
// 		return err
// 	}
// 	rS := resp.RouterStaticSet[0]
// 	d.Set("router", rS.RouterID)
// 	d.Set("type", int(rS.StaticType))
// 	d.Set("val1", rS.Val1)
// 	d.Set("val2", rS.Val2)
// 	d.Set("val3", rS.Val3)
// 	d.Set("val4", rS.Val4)
// 	d.Set("val5", rS.Val5)
// 	return nil
// }
// func resourceQingcloudRouterStaticUpdate(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).router
// 	if _, err := RouterTransitionStateRefresh(clt, d.Get("router").(string)); err != nil {
// 		return err
// 	}
// 	params := router.ModifyRouterStaticAttributesRequest{}
// 	params.RouterStatic.Set(d.Id())
// 	params.RouterStaticName.Set(d.Get(resourceName).(string))
// 	params.Val1.Set(d.Get("val1").(string))
// 	params.Val2.Set(d.Get("val2").(string))
// 	params.Val3.Set(d.Get("val3").(string))
// 	params.Val4.Set(d.Get("val4").(string))
// 	params.Val5.Set(d.Get("val5").(string))
// 	_, err := clt.ModifyRouterStaticAttributes(params)
// 	if err != nil {
// 		return err
// 	}
// 	return applyRouterUpdates(meta, d.Get("router").(string))
// }
// func resourceQingcloudRouterStaticDelete(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).router
// 	if _, err := RouterTransitionStateRefresh(clt, d.Get("router").(string)); err != nil {
// 		return err
// 	}
// 	params := router.DeleteRouterStaticsRequest{}
// 	params.RouterStaticsN.Add(d.Id())
// 	_, err := clt.DeleteRouterStatics(params)
// 	if err != nil {
// 		return err
// 	}
// 	if err = applyRouterUpdates(meta, d.Get("router").(string)); err != nil {
// 		return err
// 	}
// 	d.SetId("")
// 	return nil
// }
