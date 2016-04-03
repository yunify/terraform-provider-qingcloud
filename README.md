# terraform-qingcloud
Terraform 的 QingCloud 插件

# 使用方式

## 安装 qingcloud-provider
```
go install -v github.com/CuriosityChina/terraform-qingcloud/provider-qingcloud
```

## 设置 terraform 的插件路径

```
# 启动编辑器
subl ~/.terraformrc

# 修改如下qingcloud 到你本地的路径
providers {
	qingcloud = "/Users/YOUR/GO/PATH/bin/provider-qingcloud"
}
```

## 目前我们会用到的资源：

+ Instance
+ Volume
+ Vxnet
+ Routers
+ Eip
+ SecurityGroups
+ Keypairs
+ Image
+ LoadBalancer
+ Tag

其他资源欢迎提交 PR
