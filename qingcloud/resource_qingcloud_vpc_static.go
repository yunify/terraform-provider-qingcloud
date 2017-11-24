package qingcloud

import "github.com/hashicorp/terraform/helper/schema"

func resourceQingcloudVpcStatic() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudVpcCreate,
		Read:   resourceQingcloudVpcRead,
		Update: resourceQingcloudVpcUpdate,
		Delete: resourceQingcloudVpcDelete,
		Schema: map[string]*schema.Schema{
			resourceName: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of Vpc",
			},
			"type": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      1,
				ValidateFunc: withinArrayInt(0, 1, 2, 3),
				Description: "Type of Vpc: 0 - medium, 1 - small, 2 - large, 3 - ultra-large, default 1	",
			},
			"vpc_network": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: withinArrayString("192.168.0.0/16", "172.16.0.0/16", "172.17.0.0/16",
					"172.18.0.0/16", "172.19.0.0/16", "172.20.0.0/16", "172.21.0.0/16", "172.22.0.0/16",
					"172.23.0.0/16", "172.24.0.0/16", "172.25.0.0/16"),
				Description: "Network address range of vpc.",
			},
			resourceTagIds:   tagIdsSchema(),
			resourceTagNames: tagNamesSchema(),
			"eip_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The eip's id used by the vpc",
			},
			"security_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The security group's id used by the vpc",
			},
			resourceDescription: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of vpc",
			},
			"private_ip": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The private ip of vpc",
			},
			"public_ip": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public ip of vpc",
			},
		},
	}
}

func resourceQingcloudVpcStaticCreate(d *schema.ResourceData, meta interface{})  {

}

func resourceQingcloudVpcStaticRead(d *schema.ResourceData, meta interface{})  {

}

func resourceQingcloudVpcStaticUpdate(d *schema.ResourceData, meta interface{})  {

}

func resourceQingcloudVpcStaticDelete(d *schema.ResourceData, meta interface{})  {

}

