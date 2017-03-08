package qingcloud

// import (
// 	"github.com/hashicorp/terraform/helper/schema"
// 	"github.com/magicshui/qingcloud-go/eip"
// 	"github.com/magicshui/qingcloud-go/router"
// )

// func resourceQingcloudEipAssociate() *schema.Resource {
// 	return &schema.Resource{
// 		Create: resourceQingcloudEipAssociateCreate,
// 		Read:   resourceQingcloudEipAssociateRead,
// 		Update: resourceQingcloudEipAssociateUpdate,
// 		Delete: resourceQingcloudEipAssociateDelete,
// 		Schema: map[string]*schema.Schema{
// 			"eip": &schema.Schema{
// 				Type:        schema.TypeString,
// 				Required:    true,
// 				Description: "公网IP",
// 			},
// 			"resource_type": &schema.Schema{
// 				Type:        schema.TypeString,
// 				Required:    true,
// 				Description: "资源类型，目前支持 router 和 instance",
// 			},
// 			"resource_id": &schema.Schema{
// 				Type:     schema.TypeString,
// 				Optional: true,
// 			},
// 		},
// 	}
// }

// func resourceQingcloudEipAssociateCreate(d *schema.ResourceData, meta interface{}) error {
// 	eipID := d.Get("eip").(string)
// 	d.SetId(eipID)

// 	resourceType := d.Get("resource_type").(string)
// 	resourceID := d.Get("resource_id").(string)
// 	switch resourceType {
// 	case "router":
// 		// Router
// 		clt := meta.(*QingCloudClient).router
// 		input :=
// 		params := router.ModifyRouterAttributesRequest{}
// 		params.Eip.Set(eipID)
// 		params.Router.Set(resourceID)
// 		if _, err := clt.ModifyRouterAttributes(params); err != nil {
// 			return err
// 		}
// 		return applyRouterUpdates(meta, resourceID)

// 	default:
// 		// instance
// 		clt := meta.(*QingCloudClient).eip
// 		params := eip.AssociateEipRequest{}
// 		params.Instance.Set(resourceID)
// 		params.Eip.Set(d.Id())
// 		if _, err := clt.AssociateEip(params); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func resourceQingcloudEipAssociateRead(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).eip
// 	params := eip.DescribeEipsRequest{}
// 	params.EipsN.Add(d.Id())
// 	if resp, err := clt.DescribeEips(params); err != nil {
// 		return err
// 	} else {
// 		d.Set("resource_type", resp.EipSet[0].Resource.ResourceType)
// 		d.Set("resource_id", resp.EipSet[0].Resource.ResourceID)
// 	}

// 	return nil
// }

// func resourceQingcloudEipAssociateUpdate(d *schema.ResourceData, meta interface{}) error {
// 	return nil
// }

// func resourceQingcloudEipAssociateDelete(d *schema.ResourceData, meta interface{}) error {
// 	resourceType := d.Get("resource_type").(string)
// 	resourceID := d.Get("resource_id").(string)
// 	switch resourceType {
// 	case "router":
// 		clt := meta.(*QingCloudClient).router
// 		params := router.ModifyRouterAttributesRequest{}
// 		params.Eip.Set(d.Id())
// 		params.Router.Set(resourceID)
// 		if _, err := clt.ModifyRouterAttributes(params); err != nil {
// 			return err
// 		}
// 		err := applyRouterUpdates(meta, resourceID)
// 		if err != nil {
// 			return err
// 		}
// 	default:
// 		clt := meta.(*QingCloudClient).eip
// 		params := eip.DissociateEipsRequest{}
// 		params.EipsN.Add(d.Id())
// 		if _, err := clt.DissociateEips(params); err != nil {
// 			return err
// 		}
// 	}
// 	d.SetId("")
// 	return nil
// }
