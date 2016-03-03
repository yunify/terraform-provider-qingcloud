package qingcloud

// import (
// 	"fmt"
// 	"log"
// 	// "github.com/hashicorp/terraform/helper/resource"
// 	"github.com/hashicorp/terraform/helper/schema"
// 	"github.com/magicshui/qingcloud-go/tag"
// )

// // resourceQingcloudTagCreate
// func resourceQingcloudTagCreate(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).tag

// 	// TODO: 这个地方以后需要判断错误
// 	tagName := d.Get("name").(string)
// 	tagType := d.Get("type").(int)
// 	tagVPCNetwork := d.Get("vpc_network").(string)
// 	tagSecurityGroupID := d.Get("security_group_id").(string)

// 	params := tag.CreatetagsRequest{}
// 	params.tagName.Set(tagName)
// 	params.tagType.Set(tagType)
// 	params.VpcNetwork.Set(tagVPCNetwork)
// 	params.SecurityGroup.Set(tagSecurityGroupID)

// 	resp, err := clt.Createtags(params)
// 	if err != nil {
// 		return fmt.Errorf("Error create tag ", err)
// 	}

// 	// description := d.Get("description").(string)
// 	// if description != "" {
// 	// 	modifyAtrributes := tag.ModifytagAttributesRequest{}

// 	// 	modifyAtrributes.tag.Set(resp.tags[0])
// 	// 	modifyAtrributes.Description.Set(description)
// 	// 	_, err := clt.ModifytagAttributes(modifyAtrributes)
// 	// 	if err != nil {
// 	// 		return fmt.Errorf("Error modify tag description: %s", err)
// 	// 	}
// 	// }

// 	d.SetId(resp.tags[0])
// 	return resourceQingcloudTagRead(d, meta)
// }

// func resourceQingcloudTagRead(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).tag

// 	// 设置请求参数
// 	params := tag.DescribetagsRequest{}
// 	params.tagsN.Add(d.Id())
// 	params.Verbose.Set(1)

// 	resp, err := clt.Describetags(params)
// 	if err != nil {
// 		return fmt.Errorf("Error retrieving tags: %s", err)
// 	}
// 	log.Printf("Fetch the tag information: %s", resp)
// 	for _, v := range resp.tagSet {
// 		if v.tagID == d.Id() {
// 			d.Set("name", v.tagName)
// 			d.Set("type", v.tagType)
// 			d.Set("vpc_network", v.Vxnets)
// 			d.Set("security_group_id", v.SecurityGroupID)
// 			d.Set("description", v.Description)

// 			// 如下状态是稍等来获取的
// 			d.Set("vxnets", v.Vxnets)
// 			d.Set("private_ip", v.PrivateIP)
// 			d.Set("is_applied", v.IsApplied)
// 			d.Set("eip", v.Eip)
// 			d.Set("status", v.Status)
// 			d.Set("transition_status", v.TransitionStatus)
// 			return nil
// 		}
// 	}
// 	d.SetId("")
// 	return nil
// }

// func resourceQingcloudTagDelete(d *schema.ResourceData, meta interface{}) error {

// 	return nil
// }

// func resourceQingcloudTagUpdate(d *schema.ResourceData, meta interface{}) error {
// 	clt := meta.(*QingCloudClient).tag

// 	params := tag.ModifytagAttributesRequest{}
// 	if !d.HasChange("description") && !d.HasChange("name") {
// 		return nil
// 	}
// 	params.tag.Set(d.Id())

// 	if d.HasChange("description") {
// 		params.Description.Set(d.Get("description").(string))
// 	}
// 	if d.HasChange("name") {
// 		params.tagName.Set(d.Get("name").(string))
// 	}

// 	_, err := clt.ModifytagAttributes(params)
// 	if err != nil {
// 		return fmt.Errorf("Error update tag: %s", err)
// 	}
// 	return nil
// }

// func resourceQingcloudTagSchema(computed bool) map[string]*schema.Schema {
// 	return map[string]*schema.Schema{
// 		"name": &schema.Schema{
// 			Type:     schema.TypeString,
// 			Required: true,
// 		},

// 		"type": &schema.Schema{
// 			Type:     schema.TypeInt,
// 			Optional: true,
// 		},
// 		"vpc_network": &schema.Schema{
// 			Type:     schema.TypeString,
// 			Optional: true,
// 		},
// 		"security_group_id": &schema.Schema{
// 			Type:     schema.TypeString,
// 			Optional: true,
// 		},

// 		"description": &schema.Schema{
// 			Type:     schema.TypeString,
// 			Optional: true,
// 		},
// 		"private_ip": &schema.Schema{
// 			Type:     schema.TypeString,
// 			Computed: true,
// 		},
// 		"is_applied": &schema.Schema{
// 			Type:     schema.TypeInt,
// 			Computed: true,
// 		},

// 		"status": &schema.Schema{
// 			Type:     schema.TypeString,
// 			Computed: true,
// 		},
// 		"transition_status": &schema.Schema{
// 			Type:     schema.TypeString,
// 			Computed: true,
// 		},
// 		"id": &schema.Schema{
// 			Type:     schema.TypeString,
// 			Optional: true,
// 			Computed: true,
// 		},
// 	}
// }
