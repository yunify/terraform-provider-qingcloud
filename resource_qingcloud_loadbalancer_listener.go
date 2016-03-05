package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/magicshui/qingcloud-go/loadbalancer"
)

func resourceQingcloudLoadbalancerListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudLoadbalancerListenerCreate,
		Read:   resourceQingcloudLoadbalancerListenerRead,
		Update: resourceQingcloudLoadbalancerListenerUpdate,
		Delete: resourceQingcloudLoadbalancerListenerDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"loadbalancer": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"protocol": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: withinArrayString("tcp", "http", "https"),
			},
			"certificate": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"mode": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: withinArrayString("roundrobin", "leastconn", "source"),
				Description:  "监听器负载均衡方式：支持 roundrobin (轮询)， leastconn (最小连接)和 source (源地址) 三种。",
			},
			"session_sticky": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"forwardfor": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				Description: `转发请求时需要附加的 HTTP Header。此值是由当前支持的3个附加头字段以“按位与”的方式得到的十进制数：
						1. X-Forwarded-For: bit 位是1 (二进制的1)，表示是否将真实的客户端IP传递给后端。 附加选项“获取客户端IP”关闭时，后端 server 得到的 client IP 是负载均衡器本身的 IP 地址。 在开启本功能之后，后端服务器可以通过请求中的 X-Forwarded-For 字段来获取真实的用户IP。
						2. QC-LBID: bit 位是2 (二进制的10)，表示 Header 中是否包含 LoadBalancer 的 ID
						3. QC-LBIP: bit 位是3 (二进制的100)，表示 Header 中是否包含 LoadBalancer 的公网IP
						
						例如 Header 中包含 X-Forwarded-For 和 QC-LBIP 的话，forwarfor 的值则为:
						“X-Forwarded-For | QC-LBIP”，二进制结果为101，最后转换成十进制得到5。`,
			},
			"health_check_method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"health_check_option": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "10|5|2|5",
				Description: "inter | timeout | fall | rise",
			},
			"listener_option": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Description: `附加选项。此值是由当前支持的2个附加选项以“按位与”的方式得到的十进制数：
							1. 取消URL校验: bit 位是1 (二进制的1)，表示是否可以让负载均衡器接受不符合编码规范的 URL，例如包含未编码中文字符的URL等
							2. 获取客户端IP: bit 位是2 (二进制的10)，表示是否将客户端的IP直接传递给后端。 开启本功能后，负载均衡器对与后端是完全透明的。
							后端主机TCP连接得到的源地址是客户端的IP， 而不是负载均衡器的IP。注意：仅支持受管网络中的后端。使用基础网络后端时，此功能无效。
							3. 数据压缩: bit 位是4 (二进制的100)， 表示是否使用gzip算法压缩文本数据，以减少网络流量。
							4. 禁用不安全的加密方式: bit 位是8 (二进制的1000), 禁用存在安全隐患的加密方式， 可能会不兼容低版本的客户端。
							`,
			},
		},
	}
}

func resourceQingcloudLoadbalancerListenerCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	params := loadbalancer.AddLoadBalancerListenersRequest{}
	params.Loadbalancer.Set(d.Get("loadbalancer").(string))
	params.ListenersNListenerPort.Add(int64(d.Get("port").(int)))
	params.ListenersNListenerProtocol.Add(d.Get("protocol").(string))
	params.ListenersNServerCertificateId.Add(d.Get("certificate").(string))
	params.ListenersNBackendProtocol.Add(d.Get("protocol").(string))
	params.ListenersNLoadbalancerListenerName.Add(d.Get("name").(string))
	params.ListenersNBalanceMode.Add(d.Get("mode").(string))
	params.ListenersNSessionSticky.Add(d.Get("session_sticky").(string))
	params.ListenersNForwardfor.Add(int64(d.Get("forwardfor").(int)))
	params.ListenersNHealthyCheckMethod.Add(d.Get("health_check_method").(string))
	params.ListenersNHealthyCheckOption.Add(d.Get("health_check_option").(string))
	resp, err := clt.AddLoadBalancerListeners(params)
	if err != nil {
		return err
	}
	lb := resp.LoadbalancerListeners[0]
	d.SetId(lb)
	return nil
}

func resourceQingcloudLoadbalancerListenerRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	params := loadbalancer.DescribeLoadBalancerListenersRequest{}
	params.LoadbalancerListenersN.Add(d.Id())
	resp, err := clt.DescribeLoadBalancerListeners(params)
	if err != nil {
		return err
	}
	lb := resp.LoadbalancerListenerSet[0]
	d.Set("loadbalancer", lb.LoadbalancerID)
	d.Set("port", lb.ListenerPort)
	d.Set("protocol", lb.ListenerProtocol)
	d.Set("name", lb.LoadbalancerListenerName)
	d.Set("mode", lb.BalanceMode)
	d.Set("session_sticky", lb.SessionSticky)
	d.Set("forwardfor", lb.Forwardfor)
	d.Set("health_check_method", lb.HealthyCheckMethod)
	d.Set("health_check_option", lb.HealthyCheckOption)
	return nil
}

func resourceQingcloudLoadbalancerListenerDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	params := loadbalancer.DeleteLoadBalancerListenersRequest{}
	params.LoadbalancerListenersN.Add(d.Id())
	_, err := clt.DeleteLoadBalancerListeners(params)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceQingcloudLoadbalancerListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}
