package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"

	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifyCacheAttributes(d *schema.ResourceData, meta interface{}, create bool) error {
	clt := meta.(*QingCloudClient).cache
	input := new(qc.ModifyCacheAttributesInput)
	input.Cache = qc.String(d.Id())
	if create {
		if description := d.Get("description").(string); description == "" {
			return nil
		}
		input.Description = qc.String(d.Get("description").(string))
	} else {
		if !d.HasChange("description") && !d.HasChange("name") && !d.HasChange("auto_backup_time") {
			return nil
		}
		if d.HasChange("description") {
			input.Description = qc.String(d.Get("description").(string))
		}
		if d.HasChange("name") {
			input.CacheName = qc.String(d.Get("name").(string))
		}
		if d.HasChange("auto_backup_time") {
			input.AutoBackupTime = qc.Int(d.Get("auto_backup_time").(int))
		}
	}
	_, err := clt.ModifyCacheAttributes(input)
	return err
}

func resizeCache(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).cache
	cacheID := d.Id()
	// stop cache
	if _, err := CacheTransitionStateRefresh(clt, cacheID); err != nil {
		return err
	}
	stopInput := new(qc.StopCachesInput)
	stopInput.Caches = []*string{qc.String(cacheID)}
	if _, err := clt.StopCaches(stopInput); err != nil {
		return err
	}
	// resize cache
	if _, err := CacheTransitionStateRefresh(clt, cacheID); err != nil {
		return err
	}
	resizeInput := new(qc.ResizeCachesInput)
	resizeInput.Caches = []*string{qc.String(cacheID)}
	resizeInput.CacheSize = qc.Int(d.Get("size").(int))
	if _, err := clt.ResizeCaches(resizeInput); err != nil {
		return err
	}
	// start cache
	if _, err := CacheTransitionStateRefresh(clt, cacheID); err != nil {
		return err
	}
	startInput := new(qc.StartCachesInput)
	startInput.Caches = []*string{qc.String(cacheID)}
	if _, err := clt.StartCaches(startInput); err != nil {
		return err
	}
	if _, err := CacheTransitionStateRefresh(clt, cacheID); err != nil {
		return err
	}
	return nil
}
