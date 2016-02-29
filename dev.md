## 每添加一个配置

provider.go 中 添加：

```
ResourcesMap: map[string]*schema.Resource{
			// "qingcloud_eip":     resourceQingcloudEip(),
			"qingcloud_keypair":       resourceQingcloudKeypair(),
			"qingcloud_securitygroup": resourceQingcloudSecuritygroup(),
		},
		ConfigureFunc: providerConfigure,
	}
```
-------------------------
config.go 中 添加：
```
type QingCloudClient struct {
	eip           *eip.EIP
	keypair       *keypair.KEYPAIR
	securitygroup *securitygroup.SECURITYGROUP
}
```

和

```
func (c *Config) Client() (*QingCloudClient, error) {
	clt := qingcloud.NewClient()
	clt.ConnectToZone(c.Zone, c.ID, c.Secret)

	return &QingCloudClient{
		eip:           eip.NewClient(clt),
		keypair:       keypair.NewClient(clt),
		securitygroup: securitygroup.NewClient(clt),
	}, nil
}

```

## 在创建一个新的 resource 时

1. describeXXXX 需要设置请求的 params 的 verbose 为 1
2. Read 必须要在拿到结果以后，设置 read