package qingcloud

import (
	"os"

	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["access_key"],
			},
			"secret_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["secret_key"],
			},
			"zone": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["zone"],
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"qingcloud_eip": resourceQingcloudEip(),
			// "qingcloud_eip_associate": resourceQingcloudEipAssociate(),
			"qingcloud_keypair":             resourceQingcloudKeypair(),
			"qingcloud_security_group":      resourceQingcloudSecurityGroup(),
			"qingcloud_security_group_rule": resourceQingcloudSecurityGroupRule(),
			"qingcloud_vxnet":               resourceQingcloudVxnet(),
			"qingcloud_router":              resourceQingcloudRouter(),
			// "qingcloud_router_static":      resourceQingcloudRouterStatic(),
			// "qingcloud_router_static_entry": resourceQingcloudRouterStaticEntry(),
			"qingcloud_instance":              resourceQingcloudInstance(),
			"qingcloud_volume":                resourceQingcloudVolume(),
			"qingcloud_volume_attachment":     resourceQingcloudVolumeAttachment(),
			"qingcloud_tag":                   resourceQingcloudTag(),
			"qingcloud_cache":                 resourceQingcloudCache(),
			"qingcloud_cache_parameter_group": resourceQingcloudCacheParameterGroup(),
			// "qingcloud_mongo": resourceQingcloudMongo(),
			// "qingcloud_loadbalancer":             resourceQingcloudLoadbalancer(),
			// "qingcloud_loadbalancer_listener":    resourceQingcloudLoadbalancerListener(),
			// "qingcloud_loadbalancer_backend":     resourceQingcloudLoadbalancerBackend(),
			// "qingcloud_loadbalancer_policy":      resourceQingcloudLoadbalancerPloicy(),
			// "qingcloud_loadbalancer_policy_rule": resourceQingcloudLoadbalancerPloicyRule(),
		},
		ConfigureFunc: providerConfigure,
	}
}

var qingcloudMutexKV = mutexkv.NewMutexKV()

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	accesskey, ok := d.GetOk("access_key")
	if !ok {
		accesskey = os.Getenv("QINGCLOUD_ACCESS_KEY")
	}
	secretkey, ok := d.GetOk("secret_key")
	if !ok {
		secretkey = os.Getenv("QINGCLOUD_SECRET_KEY")
	}
	zone, ok := d.GetOk("zone")
	if !ok {
		zone = os.Getenv("QINGCLOUD_ZONE")
		if zone == "" {
			zone = DEFAULT_ZONE
		}
	}
	config := Config{
		ID:     accesskey.(string),
		Secret: secretkey.(string),
		Zone:   zone.(string),
	}
	return config.Client()
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key": "qingcloud access key ID ",
		"secret_key": "qingcloud access key secret",
		"zone":   "qingcloud reigon zone",
	}
}
