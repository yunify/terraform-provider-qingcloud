package qingcloud

// import (
// 	"fmt"

// 	// "github.com/magicshui/qingcloud-go/eip"

// 	"github.com/hashicorp/terraform/helper/schema"
// )

// func resourceQingcloudEip() *schema.Resource {
// 	res := &schema.Resource{
// 		Create: resourceQingcloudEipCreate,
// 		Read:   resourceQingcloudEipRead,
// 		Delete: resourceQingcloudEipDelete,
// 		Update: resourceQingcloudEipUpdate,

// 		// Schema
// 		Schema: nil,
// 	}
// 	return res
// }

// // 创建Eip
// func resourceQingcloudEipCreate(d *schema.ResourceData, meta interface{}) error {
// 	return nil
// }

// func resourceQingcloudEipRead(d *schema.ResourceData, meta interface{}) error {
// 	return nil
// }

// func resourceQingcloudEipDelete(d *schema.ResourceData, meta interface{}) error {
// 	return nil
// }

// func resourceQingcloudEipUpdate(d *schema.ResourceData, meta interface{}) error {
// 	return nil
// }

// func resourceQingcloudEipSchema() map[string]*schema.Schema {
// 	resourceSchema := map[string]*schema.Schema{
// 		// 开通账号
// 		"bandwidth": &schema.Schema{
// 			Type:     schema.TypeInt,
// 			Optional: false,
// 		},
// 		"billing_mode": &schema.Schema{
// 			Type:     schema.TypeString,
// 			Optional: false,
// 			Default:  "bandwidth",
// 			ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
// 				value := v.(string)
// 				if value != "bandwidth" && value != "traffic" {
// 					errors = append(errors, fmt.Errorf(
// 						"%q 计费模式不正确，只能是 bandwidth 或者  traffic", k, nil))
// 				}
// 				return
// 			},
// 		},
// 		"eip_name": &schema.Schema{},
// 		"count": &schema.Schema{
// 			Type:     schema.TypeInt,
// 			Optional: false,
// 			Default:  1,
// 		},
// 		"need_icp": &schema.Schema{
// 			Type:     schema.TypeInt,
// 			Optional: false,
// 			Default:  0,
// 		},
// 	}
// 	return resourceSchema
// }
