package qingcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/magicshui/qingcloud-go/eip"
	"github.com/magicshui/qingcloud-go/router"
	"log"
)

func modifyEipAttributes(d *schema.ResourceData, meta interface{}, create bool) error {
	clt := meta.(*QingCloudClient).eip
	modifyAtrributes := eip.ModifyEipAttributesRequest{}
	if create {
		if description := d.Get("description").(string); description == "" {
			return nil
		}
	} else {
		if !d.HasChange("description") && !d.HasChange("name") {
			return nil
		}
	}

	modifyAtrributes.Eip.Set(d.Id())
	modifyAtrributes.Description.Set(d.Get("description").(string))
	modifyAtrributes.EipName.Set(d.Get("name").(string))
	_, err := clt.ModifyEipAttributes(modifyAtrributes)
	if err != nil {
		return fmt.Errorf("Error modify eip description: %s", err)
	}
	return nil
}

func dissociateEipFromResource(meta interface{}, eipID, resourceType, resourceID string) error {
	switch resourceType {
	case "router":
		log.Printf("[Debug] dissociate eip form resource %s", "router")
		clt := meta.(*QingCloudClient).router
		params := router.ModifyRouterAttributesRequest{}
		params.Eip.Set("")
		params.Router.Set(resourceID)
		_, err := clt.ModifyRouterAttributes(params)
		if err != nil {
			return err
		}
		p2 := router.UpdateRoutersRequest{}
		p2.RoutersN.Add(resourceID)
		_, err = clt.UpdateRouters(p2)
		if err != nil {
			return err
		}
		_, err = RouterTransitionStateRefresh(clt, resourceID)
		return err
	default:
		clt := meta.(*QingCloudClient).eip
		params := eip.DissociateEipsRequest{}
		params.EipsN.Add(eipID)
		_, err := clt.DissociateEips(params)
		if err != nil {
			return err
		}
		_, err = InstanceTransitionStateRefresh(meta.(*QingCloudClient).instance, resourceID)
		return err
	}
}

func getEipSourceMap(data eip.Eip) map[string]interface{} {
	var a = make(map[string]interface{}, 3)
	a["id"] = data.Resource.ResourceID
	a["name"] = data.Resource.ResourceName
	a["type"] = data.Resource.ResourceType
	return a
}
