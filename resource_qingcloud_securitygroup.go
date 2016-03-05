package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/magicshui/qingcloud-go/securitygroup"
)

func resourceQingcloudSecuritygroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudSecuritygroupCreate,
		Read:   resourceQingcloudSecuritygroupRead,
		Update: resourceQingcloudSecuritygroupUpdate,
		Delete: resourceQingcloudSecuritygroupDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "防火墙名称",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "防火墙介绍",
			},
		},
	}
}

func resourceQingcloudSecuritygroupCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup

	// TODO: 这个地方以后需要判断错误
	securityGroupName := d.Get("name").(string)

	// 开始创建 ssh 密钥
	params := securitygroup.CreateSecurityGroupRequest{}
	params.SecurityGroupName.Set(securityGroupName)

	resp, err := clt.CreateSecurityGroup(params)
	if err != nil {
		return fmt.Errorf("Error create security group", err)
	}

	description := d.Get("description").(string)
	if description != "" {
		modifyAtrributes := securitygroup.ModifySecurityGroupAttributesRequest{}
		modifyAtrributes.SecurityGroup.Set(resp.SecurityGroupId)
		modifyAtrributes.Description.Set(description)
		_, err := clt.ModifySecurityGroupAttributes(modifyAtrributes)
		if err != nil {
			// 这里可以不用返回错误
			return fmt.Errorf("Error modify security description: %s", err)
		}
	}

	d.SetId(resp.SecurityGroupId)
	return nil
}

func resourceQingcloudSecuritygroupRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup

	// 设置请求参数
	params := securitygroup.DescribeSecurityGroupsRequest{}
	params.SecurityGroupsN.Add(d.Id())
	params.Verbose.Set(1)

	resp, err := clt.DescribeSecurityGroups(params)
	if err != nil {
		return fmt.Errorf("Error retrieving Keypair: %s", err)
	}
	for _, sg := range resp.SecurityGroupSet {
		if sg.SecurityGroupID == d.Id() {
			d.Set("name", sg.SecurityGroupName)
			d.Set("description", sg.Description)
			return nil
		}
	}
	d.SetId("")
	return nil
}

func resourceQingcloudSecuritygroupDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	// 判断当前的防火墙是否有人在使用
	describeParams := securitygroup.DescribeSecurityGroupsRequest{}
	describeParams.SecurityGroupsN.Add(d.Id())
	describeParams.Verbose.Set(1)

	resp, err := clt.DescribeSecurityGroups(describeParams)
	if err != nil {
		return fmt.Errorf("Error retrieving Keypair: %s", err)
	}
	for _, sg := range resp.SecurityGroupSet {
		if sg.SecurityGroupID == d.Id() {
			if len(sg.Resources) > 0 {
				// 要删除的防火墙已加载规则到主机，则需要先调用 ApplySecurityGroup 将其他防火墙的规则应用到对应主机，之后才能被删除。
				// 要删除的防火墙已加载规则到路由器，则需要先调用 ModifyRouterAttributes 并 UpdateRouters 将其他防火墙的规则应用到对应路由器，之后才能被删除。
				return fmt.Errorf("Current security group is in using", nil)
			}
		}
	}

	params := securitygroup.DeleteSecurityGroupsRequest{}
	params.SecurityGroupsN.Add(d.Id())
	_, err = clt.DeleteSecurityGroups(params)
	if err != nil {
		return fmt.Errorf("Error delete security group %s", err)
	}
	d.SetId("")
	return nil
}

func resourceQingcloudSecuritygroupUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup

	if !d.HasChange("name") && !d.HasChange("description") {
		return nil
	}

	params := securitygroup.ModifySecurityGroupAttributesRequest{}
	if d.HasChange("description") {
		params.Description.Set(d.Get("description").(string))
	}
	if d.HasChange("name") {
		params.SecurityGroupName.Set(d.Get("name").(string))
	}
	params.SecurityGroup.Set(d.Id())
	_, err := clt.ModifySecurityGroupAttributes(params)
	if err != nil {
		return fmt.Errorf("Error modify security group %s", d.Id())
	}
	return nil
}
