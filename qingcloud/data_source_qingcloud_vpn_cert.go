package qingcloud

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

const (
	resourceVPNCertRouterId   = "router_id"
	resourceVPNCertPlatform   = "platform"
	resourceVPNCertClientCrt  = "client_crt"
	resourceVPNCertClientKey  = "client_key"
	resourceVPNCertStaticKey  = "static_key"
	resourceVPNCertCaCert     = "ca_cert"
	resourceVPNCertConfSample = "conf_sample"
)

func dataSourceQingcloudVpnCert() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpnCertRead,

		Schema: map[string]*schema.Schema{
			resourceVPNCertRouterId: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			resourceVPNCertPlatform: &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "linux",
				ValidateFunc: withinArrayString("linux", "windows", "mac"),
			},
			resourceVPNCertClientCrt: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			resourceVPNCertClientKey: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			resourceVPNCertStaticKey: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			resourceVPNCertCaCert: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			resourceVPNCertConfSample: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVpnCertRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.GetVPNCertsInput)
	input.Router = getSetStringPointer(d, resourceVPNCertRouterId)
	input.Platform = getSetStringPointer(d, resourceVPNCertPlatform)
	var output *qc.GetVPNCertsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.GetVPNCerts(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	d.Set(resourceVPNCertClientCrt, qc.StringValue(output.ClientCrt))
	d.Set(resourceVPNCertClientKey, qc.StringValue(output.ClientKey))
	d.Set(resourceVPNCertStaticKey, qc.StringValue(output.StaticKey))
	d.Set(resourceVPNCertCaCert, qc.StringValue(output.CaCert))
	if d.Get(resourceVPNCertPlatform) == "linux" {
		d.Set(resourceVPNCertConfSample, qc.StringValue(output.LinuxConfSample))
	} else if d.Get(resourceVPNCertPlatform) == "mac" {
		d.Set(resourceVPNCertConfSample, qc.StringValue(output.MacConfSample))
	} else {
		d.Set(resourceVPNCertConfSample, qc.StringValue(output.WindowsConfSample))
	}
	d.SetId(time.Now().UTC().String())
	return nil
}
