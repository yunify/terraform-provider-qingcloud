package qingcloud

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func getInstanceIDAndVolumeID(d *schema.ResourceData) (string, string, error) {
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid resource id %s", d.Id())
	}
	return parts[0], parts[1], nil
}

func genVolumeAttachmentID(d *schema.ResourceData) string {
	return fmt.Sprintf("%s:%s", d.Get("volume_id").(string), d.Get("instance_id").(string))
}
