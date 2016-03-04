package qingcloud

import (
	// "github.com/hashicorp/terraform/helper/schema"
	"github.com/magicshui/qingcloud-go/loadbalancer"
)

func updateLoadBalancer(meta interface{}, id string) error {
	clt := meta.(*QingCloudClient).loadbalancer
	params := loadbalancer.UpdateLoadBalancersRequest{}
	params.LoadbalancersN.Add(id)
	_, err := clt.UpdateLoadBalancers(params)
	if err != nil {
		return err
	}

	_, err = LoadbalancerTransitionStateRefresh(clt, id)
	return err
}
