package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func resourceQingcloudVpcStatic() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudVpcStaticCreate,
		Read:   resourceQingcloudVpcStaticRead,
		Update: resourceQingcloudVpcStaticUpdate,
		Delete: resourceQingcloudVpcStaticDelete,
		Schema: map[string]*schema.Schema{
			resourceName: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of Vpc Static",
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"static_type": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: withinArrayInt(1, 2, 3, 4, 6, 7, 8),
				Description: "1 : port_forwarding , 2 " +
					"2 : VPN rule " +
					"3 : DHCP " +
					"4 :  Two layers GRE " +
					"6 :  Three layers GRE " +
					"7 :  Three layers IPsec " +
					"8 :  Private DNS ",
			},
			"val1": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "port_forwarding : source port " +
					"VPN : type of vpn ,'openvpn', 'pptp','l2tp', default 'openvpn' " +
					"DHCP : id of DHCP host" +
					"Two layers GRE : remote ip , secret key, example : gre|1.2.3.4|888	" +
					"Three layers GRE: remote ip , secret key, local p2p ip , opposite end p2p ip , example : 6.6.6.6|key|1.2.3.4|4.3.2.1 " +
					"Three layers IPsec : remote ip(support 0.0.0.0 for any) ; encryption method :phase2alg&ike , default aes ; secret key & remote device id" +
					"Private DNS , private domain name",
			},
			"val2": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "port_forwarding : destination ip " +
					"OpenVPN : VPN Server Port , default 1194" +
					"PPTP/L2TP : username & password , format (user:password)" +
					"DHCP : DHCP  configuration content " +
					"Three layers GRE: destination network , multiple networks are separated by '|' " +
					"Three layers IPsec : local network , multiple networks are separated by '|' " +
					"Private DNS , IP address ,192.168.1.2;192.168.1.3",
			},
			"val3": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "port_forwarding : destination port " +
					"OpenVPN : VPN protocol , default udp" +
					"PPTP VPN : Max Connections , 1-253" +
					"L2TP VPN :(PSK, preshared secrets) " +
					"Three layers IPsec : destination network , multiple networks are separated by '|' ",
			},
			"val4": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "port_forwarding : protocol , default tcp , support udp & tcp " +
					"VPN : client CIDR ,support 10.255.x.0/24 , default auto allocation" +
					"Three layers IPsec : tunnel pattern . default main , support main & aggrmode ",
			},
			"val5": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "OpenVPN :  verification mode " +
					"L2TP VPN : L2TP server port 1701",
			},
		},
	}
}

func resourceQingcloudVpcStaticCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.AddRouterStaticsInput)
	static := new(qc.RouterStatic)
	input.Router = qc.String(d.Get("vpc_id").(string))
	static.RouterID = qc.String(d.Get("vpc_id").(string))
	static.RouterStaticName, _ = getNamePointer(d)
	static.StaticType = qc.Int(d.Get("static_type").(int))
	static.Val1 = getResourceString(d, "val1")
	static.Val2 = getResourceString(d, "val2")
	static.Val3 = getResourceString(d, "val3")
	static.Val4 = getResourceString(d, "val4")
	input.Statics = []*qc.RouterStatic{static}
	var output *qc.AddRouterStaticsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.AddRouterStatics(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	d.SetId(qc.StringValue(output.RouterStatics[0]))
	if err := applyRouterUpdate(qc.String(d.Get("vpc_id").(string)), meta); err != nil {
		return nil
	}
	if _, err := RouterTransitionStateRefresh(clt, d.Get("vpc_id").(string)); err != nil {
		return err
	}
	return resourceQingcloudVpcStaticRead(d, meta)
}

func resourceQingcloudVpcStaticRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.DescribeRouterStaticsInput)
	input.RouterStatics = []*string{qc.String(d.Id())}
	var output *qc.DescribeRouterStaticsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeRouterStatics(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if len(output.RouterStaticSet)==0{
		d.SetId("")
		return nil
	}
	d.Set(resourceName,qc.StringValue(output.RouterStaticSet[0].RouterStaticName))
	d.Set("static_type",qc.IntValue(output.RouterStaticSet[0].StaticType))
	d.Set("val1",qc.StringValue(output.RouterStaticSet[0].Val1))
	d.Set("val2",qc.StringValue(output.RouterStaticSet[0].Val2))
	d.Set("val3",qc.StringValue(output.RouterStaticSet[0].Val3))
	d.Set("val4",qc.StringValue(output.RouterStaticSet[0].Val4))
	return nil
}

func resourceQingcloudVpcStaticUpdate(d *schema.ResourceData, meta interface{}) error {

	return resourceQingcloudVpcStaticRead(d,meta)
}

func resourceQingcloudVpcStaticDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}
