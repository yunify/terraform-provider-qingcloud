package qingcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/magicshui/qingcloud-go/router"
)

func applyRouterUpdates(meta interface{}, routerID string) error {
	clt := meta.(*QingCloudClient).router
	params := router.UpdateRoutersRequest{}
	params.RoutersN.Add(routerID)
	if _, err := clt.UpdateRouters(params); err != nil {
		return err
	}
	if _, err := RouterTransitionStateRefresh(clt, routerID); err != nil {
		return err
	}
	return nil
}

func modifyRouterAttributes(d *schema.ResourceData, meta interface{}, create bool) error {
	clt := meta.(*QingCloudClient).router
	params := router.ModifyRouterAttributesRequest{}
	params.Router.Set(d.Id())

	if create {
		if description := d.Get("description").(string); description != "" {
			params.Description.Set(description)
		}
	} else {
		if d.HasChange("description") {
			params.Description.Set(d.Get("description").(string))
		}
		if d.HasChange("name") {
			params.RouterName.Set(d.Get("name").(string))
		}
	}
	_, err := clt.ModifyRouterAttributes(params)
	if err != nil {
		return fmt.Errorf("Error modify router description: %s", err)
	}
	return nil
}
