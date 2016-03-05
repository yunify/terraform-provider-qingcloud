package qingcloud

import (
	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["token"],
			},
			"secret": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["secret"],
			},
			"zone": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["zone"],
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"qingcloud_eip":           resourceQingcloudEip(),
			"qingcloud_eip_associate": resourceQingcloudEipAssociate(),

			"qingcloud_keypair": resourceQingcloudKeypair(),

			"qingcloud_securitygroup":      resourceQingcloudSecuritygroup(),
			"qingcloud_securitygroup_rule": resourceQingcloudSecuritygroupRule(),

			"qingcloud_vxnet": resourceQingcloudVxnet(),

			"qingcloud_router":              resourceQingcloudRouter(),
			"qingcloud_router_static":       resourceQingcloudRouterStatic(),
			"qingcloud_router_static_entry": resourceQingcloudRouterStaticEntry(),

			"qingcloud_instance": resourceQingcloudInstance(),

			"qingcloud_cache": resourceQingcloudCache(),

			"qingcloud_mongo": resourceQingcloudMongo(),

			"resource_qingcloud_cache_parametergroup": resourceQingcloudCacheParameterGroup(),

			"qingcloud_volume":            resourceQingcloudVolume(),
			"qingcloud_volume_attachment": resourceQingcloudVolumeAttachment(),

			"qingcloud_loadbalancer":             resourceQingcloudLoadbalancer(),
			"qingcloud_loadbalancer_listener":    resourceQingcloudLoadbalancerListener(),
			"qingcloud_loadbalancer_backend":     resourceQingcloudLoadbalancerBackend(),
			"qingcloud_loadbalancer_policy":      resourceQingcloudLoadbalancerPloicy(),
			"qingcloud_loadbalancer_policy_rule": resourceQingcloudLoadbalancerPloicyRule(),
		},
		ConfigureFunc: providerConfigure,
	}
}

var qingcloudMutexKV = mutexkv.NewMutexKV()

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		ID:     d.Get("id").(string),
		Secret: d.Get("secret").(string),
		Zone:   d.Get("zone").(string),
	}
	return config.Client()
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"id":     "青云的 ID ",
		"secret": "青云的密钥",
		"zone":   "青云的 zone ",
	}
}
