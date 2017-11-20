package qingcloud

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifyCacheParameterGroupAttributes(d *schema.ResourceData, meta interface{}, create bool) error {
	clt := meta.(*QingCloudClient).cache
	input := new(qc.ModifyCacheParameterGroupAttributesInput)
	input.CacheParameterGroup = qc.String(d.Id())
	if create {
		if description := d.Get(resourceDescription).(string); description == "" {
			return nil
		}
		input.Description = qc.String(d.Get(resourceDescription).(string))
	} else {
		if !d.HasChange(resourceDescription) && !d.HasChange(resourceName) {
			return nil
		}
		if d.HasChange(resourceDescription) {
			input.Description = qc.String(d.Get(resourceDescription).(string))
		}
		if d.HasChange(resourceName) {
			input.CacheParameterGroupName = qc.String(d.Get(resourceName).(string))
		}
	}
	_, err := clt.ModifyCacheParameterGroupAttributes(input)
	return err
}

func cacheParameterGroupSetPassword(d *schema.ResourceData, meta interface{}, create bool) error {
	cacheType := d.Get("type").(string)
	if strings.HasPrefix(cacheType, "redis") {
		clt := meta.(*QingCloudClient).cache
		input := new(qc.UpdateCacheParametersInput)
		input.CacheParameterGroup = qc.String(d.Id())
		parameter := new(qc.CacheParameter)
		parameter.CacheParameterName = qc.String("requirepass")
		if create {
			if password := d.Get("password").(string); password == "" {
				return nil
			}
			parameter.CacheParameterValue = qc.String(d.Get("password").(string))
		} else {
			if !d.HasChange("password") {
				return nil
			}
			parameter.CacheParameterValue = qc.String(d.Get("password").(string))
		}
		input.Parameters = parameter
		if _, err := clt.UpdateCacheParameters(input); err != nil {
			return err
		}
	}
	return nil
}

func applyCacheParameterGroup(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).cache
	input := new(qc.ApplyCacheParameterGroupInput)
	input.CacheParameterGroup = qc.String(d.Id())
	if _, err := clt.ApplyCacheParameterGroup(input); err != nil {
		return err
	}
	if _, err := CacheParameterGroupTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	return nil
}
