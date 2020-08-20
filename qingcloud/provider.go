package qingcloud

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["access_key"],
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["secret_key"],
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["zone"],
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["endpoint"],
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"qingcloud_vpn_cert": dataSourceQingcloudVpnCert(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"qingcloud_eip":                   resourceQingcloudEip(),
			"qingcloud_keypair":               resourceQingcloudKeypair(),
			"qingcloud_security_group":        resourceQingcloudSecurityGroup(),
			"qingcloud_security_group_rule":   resourceQingcloudSecurityGroupRule(),
			"qingcloud_vxnet":                 resourceQingcloudVxnet(),
			"qingcloud_vpc":                   resourceQingcloudVpc(),
			"qingcloud_instance":              resourceQingcloudInstance(),
			"qingcloud_volume":                resourceQingcloudVolume(),
			"qingcloud_tag":                   resourceQingcloudTag(),
			"qingcloud_vpc_static":            resourceQingcloudVpcStatic(),
			"qingcloud_loadbalancer":          resourceQingcloudLoadBalancer(),
			"qingcloud_loadbalancer_listener": resourceQingcloudLoadBalancerListener(),
			"qingcloud_loadbalancer_backend":  resourceQingcloudLoadBalancerBackend(),
			"qingcloud_server_certificate":    resourceQingcloudServerCertificate(),
		},
		ConfigureFunc: providerConfigure,
	}
}

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
	endpoint, ok := d.GetOk("endpoint")
	if !ok {
		endpoint = os.Getenv("QINGCLOUD_ENDPOINT")
		if endpoint == "" {
			endpoint = DEFAULT_ENDPOINT
		}
	}

	config := Config{
		ID:       accesskey.(string),
		Secret:   secretkey.(string),
		Zone:     zone.(string),
		EndPoint: endpoint.(string),
	}
	client, err := config.Client()
	if err != nil {
		return nil, err
	}
	// check zone & endpoint
	describeZonesOutput, err := client.qingcloud.DescribeZones(nil)
	if err != nil {
		return nil, err
	}
	if len(describeZonesOutput.ZoneSet) == 0 {
		return nil, fmt.Errorf("can not get zone info")
	}
	i := 0
	for _, az := range describeZonesOutput.ZoneSet {
		if qc.StringValue(az.ZoneID) == zone {
			if qc.StringValue(az.Status) != StatusActive {
				return nil, fmt.Errorf(" zone: %s", qc.StringValue(az.Status))
			}
			break
		}
		i++
	}
	if i == len(describeZonesOutput.ZoneSet) {
		return nil, fmt.Errorf("can not find zone: %s", zone)
	}
	return config.Client()
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key": "qingcloud access key ID ",
		"secret_key": "qingcloud access key secret",
		"zone":       "qingcloud reigon zone",
		"endpoint":   "qingcloud api endpoint",
	}
}
