package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

const (
	resourceServerCertificateContent    = "certificate_content"
	resourceServerCertificatePrivateKey = "private_key"
)

func resourceQingcloudServerCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudServerCertificateCreate,
		Read:   resourceQingcloudServerCertificateRead,
		Update: resourceQingcloudServerCertificateUpdate,
		Delete: resourceQingcloudServerCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			resourceName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceServerCertificateContent: {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			resourceServerCertificatePrivateKey: {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				ForceNew:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceQingcloudServerCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.CreateServerCertificateInput)
	input.ServerCertificateName, _ = getNamePointer(d)
	input.PrivateKey = getSetStringPointer(d, resourceServerCertificatePrivateKey)
	input.CertificateContent = getSetStringPointer(d, resourceServerCertificateContent)
	var output *qc.CreateServerCertificateOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.CreateServerCertificate(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	d.SetId(qc.StringValue(output.ServerCertificateID))
	return resourceQingcloudServerCertificateUpdate(d, meta)
}

func resourceQingcloudServerCertificateRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.DescribeServerCertificatesInput)
	input.ServerCertificates = []*string{qc.String(d.Id())}
	var output *qc.DescribeServerCertificatesOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeServerCertificates(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if len(output.ServerCertificateSet) == 0 {
		d.SetId("")
		return nil
	}
	sg := output.ServerCertificateSet[0]
	d.Set(resourceName, qc.StringValue(sg.ServerCertificateName))
	d.Set(resourceDescription, qc.StringValue(sg.Description))
	return nil
}
func resourceQingcloudServerCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := modifyServerCertificateAttributes(d, meta); err != nil {
		return err
	}
	return resourceQingcloudServerCertificateRead(d, meta)
}

func resourceQingcloudServerCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.DeleteServerCertificatesInput)
	input.ServerCertificates = []*string{qc.String(d.Id())}
	var err error
	simpleRetry(func() error {
		_, err = clt.DeleteServerCertificates(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
