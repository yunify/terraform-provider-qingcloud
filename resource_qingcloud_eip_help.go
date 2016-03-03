package qingcloud

import (
	"github.com/magicshui/qingcloud-go/eip"
	"github.com/magicshui/qingcloud-go/router"
	"log"
)

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
