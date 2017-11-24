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
				Description: "The name of Vpc Static",
			},
			"static_type": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: withinArrayInt(1, 2, 3, 4, 5, 6, 7, 8),
				Description: "1 : port_forwarding , 2" +
							 "2 : VPN rule" +
							 "3 : DHCP" +
							 "4 :  Two layers GRE"+
							 "4 :  Two layers GRE",

			},
			"val1": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network address range of vpc.",
			},
			"val2": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The eip's id used by the vpc",
			},
			"val3": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The security group's id used by the vpc",
			},
			"val4": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of vpc",
			},
			"val5": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The private ip of vpc",
			},
			"val6": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public ip of vpc",
			},
		},
	}
}

func resourceQingcloudVpcStaticCreate(d *schema.ResourceData, meta interface{}) {

}

func resourceQingcloudVpcStaticRead(d *schema.ResourceData, meta interface{}) {

}

func resourceQingcloudVpcStaticUpdate(d *schema.ResourceData, meta interface{}) {

}

func resourceQingcloudVpcStaticDelete(d *schema.ResourceData, meta interface{}) {

}
