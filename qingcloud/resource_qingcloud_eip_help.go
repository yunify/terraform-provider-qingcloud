package qingcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/magicshui/qingcloud-go/eip"
	"github.com/magicshui/qingcloud-go/router"
	qc "github.com/yunify/qingcloud-sdk-go/service"
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
	input := new(qc.ModifyEIPAttributesInput)
	input.Description = qc.String(d.Get("description").(string))
	input.EIPName = qc.String(d.Get("name").(string))
	input.EIP = qc.String(d.Id())
	output, err := clt.ModifyEIPAttributes(input)
	if err != nil {
		return fmt.Errorf("Error modify eip attributes: %s", err)
	}
	if *output.RetCode == 0 {
		return fmt.Errorf("Error modify eip attributes: %s", output.Message)
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

func getEIPResourceMap(data *qc.EIP) map[string]interface{} {
	var a = make(map[string]interface{}, 3)
	a["resource_id"] = data.Resource.ResourceID
	a["resource_name"] = data.Resource.ResourceName
	a["resource_type"] = data.Resource.ResourceType
	return a
}
