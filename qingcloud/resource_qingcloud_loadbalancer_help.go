package qingcloud

// import (
// 	// "github.com/hashicorp/terraform/helper/schema"
// 	"github.com/magicshui/qingcloud-go/loadbalancer"
// )

// func updateLoadBalancer(meta interface{}, id string) error {
// 	clt := meta.(*QingCloudClient).loadbalancer
// 	params := loadbalancer.UpdateLoadBalancersRequest{}
// 	params.LoadbalancersN.Add(id)
// 	_, err := clt.UpdateLoadBalancers(params)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = LoadbalancerTransitionStateRefresh(clt, id)
// 	return err
// }

// func applyLoadBalancerPolicy(meta interface{}, id string) error {
// 	clt := meta.(*QingCloudClient).loadbalancer
// 	params := loadbalancer.ApplyLoadBalancerPolicyRequest{}
// 	params.LoadbalancerPolicy.Set(id)
// 	_, err := clt.ApplyLoadBalancerPolicy(params)
// 	return err
// }
