package qingcloud

// import (
// 	"github.com/hashicorp/terraform/helper/schema"
// 	"github.com/magicshui/qingcloud-go/loadbalancer"
// )

// func resourceQingcloudLoadbalancerBackend() *schema.Resource {
// 	return &schema.Resource{
// 		Create: resourceQingcloudLoadbalancerBackendCreate,
// 		Read:   resourceQingcloudLoadbalancerBackendRead,
// 		Update: resourceQingcloudLoadbalancerBackendUpdate,
// 		Delete: resourceQingcloudLoadbalancerBackendDelete,
// 		Schema: map[string]*schema.Schema{
// 			"name": &schema.Schema{
// 				Type:        schema.TypeString,
// 				Optional:    true,
// 				Description: "后端服务名称",
// 			},
// 			"listener": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Optional: true,
// 				Description: "要添加后端服务的监听器ID	",
// 			},
// 			"resource": &schema.Schema{
// 				Type:        schema.TypeString,
// 				Optional:    true,
// 				Description: "后端服务资源ID",
// 			},
// 			"policy": &schema.Schema{
// 				Type:        schema.TypeString,
// 				Required:    true,
// 				Description: "转发策略ID",
// 			},
// 			"port": &schema.Schema{
// 				Type:        schema.TypeInt,
// 				Optional:    true,
// 				Description: "后端服务端口",
// 			},
// 			"weight": &schema.Schema{
// 				Type:        schema.TypeInt,
// 				Optional:    true,
// 				Description: "后端服务权重",
// 			},
// 		},
// 	}
// }

// func resourceQingcloudLoadbalancerBackendCreate(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).loadbalancer
// 	params := loadbalancer.AddLoadBalancerBackendsRequest{}
// 	params.BackendsNLoadbalancerBackendName.Add(d.Get("name").(string))
// 	params.LoadbalancerListener.Set(d.Get("listener").(string))
// 	params.BackendsNResourceId.Add(d.Get("resource").(string))
// 	params.BackendsNLoadbalancerPolicyId.Add(d.Get("policy").(string))
// 	params.BackendsNPort.Add(int64(d.Get("port").(int)))
// 	params.BackendsNWeight.Add(int64(d.Get("weight").(int)))
// 	resp, err := clt.AddLoadBalancerBackends(params)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId(resp.LoadbalancerBackends[0])
// 	return nil
// }
// func resourceQingcloudLoadbalancerBackendRead(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).loadbalancer
// 	params := loadbalancer.DescribeLoadBalancerBackendsRequest{}
// 	params.LoadbalancerBackendsN.Add(d.Id())
// 	resp, err := clt.DescribeLoadBalancerBackends(params)
// 	if err != nil {
// 		return err
// 	}
// 	lb := resp.LoadbalancerBackendSet[0]
// 	d.Set("name", lb.LoadbalancerBackendName)
// 	d.Set("listener", lb.LoadbalancerListenerID)
// 	d.Set("resource", lb.ResourceID)
// 	d.Set("port", lb.Port)
// 	d.Set("weight", lb.Weight)
// 	return nil
// }

// func resourceQingcloudLoadbalancerBackendUpdate(d *schema.ResourceData, meta interface{}) error {
// 	return nil
// }

// func resourceQingcloudLoadbalancerBackendDelete(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).loadbalancer

// 	params := loadbalancer.DeleteLoadBalancerBackendsRequest{}
// 	params.LoadbalancerBackendsN.Add(d.Id())
// 	_, err := clt.DeleteLoadBalancerBackends(params)
// 	if err != nil {
// 		return err
// 	}
// 	d.SetId("")
// 	return nil
// }
