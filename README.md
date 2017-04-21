# terraform-qingcloud [![Build Status](https://travis-ci.org/CuriosityChina/terraform-qingcloud.svg?branch=master)](https://travis-ci.org/CuriosityChina/terraform-qingcloud)

Old version can found on [v1](https://github.com/CuriosityChina/terraform-qingcloud/tree/v1) branch

Terraform-Qingcloud-Plugin

## Usage

### Install qingcloud-provider

#### On Linux
``` bash
wget -c https://github.com/CuriosityChina/terraform-qingcloud/releases/download/v2.0.0/terraform-provider-qingcloud_linux-amd64.tgz
tar -zxvf terraform-provider-qingcloud_linux-amd64.tgz
cp terraform-provider-qingcloud_linux-amd64 $(dirname `which terraform`)/terraform-provider-qingcloud
```

#### On Mac
``` bash
wget -c https://github.com/CuriosityChina/terraform-qingcloud/releases/download/v2.0.0/terraform-provider-qingcloud_darwin-amd64.tgz
tar -zxvf terraform-provider-qingcloud_darwin-amd64.tgz
cp terraform-provider-qingcloud_darwin-amd64 $(dirname `which terraform`)/terraform-provider-qingcloud
```

## Finish Resource：
- [x] Instance
- [x] Volume
- [x] Vxnet
- [x] Router
- [x] Eip
- [x] SecurityGroups
- [x] Keypairs
- [ ] Image
- [ ] LoadBalancer
- [x] Tag
- [x] redis
- [ ] mongodb


## Contributing

1. Fork it ( https://github.com/CuriosityChina/terraform-qingcloud/fork )
2. Create your feature branch (`git checkout -b new-feature`)
3. Commit your changes (`git commit -asm 'Add some feature'`)
4. Push to the branch (`git push origin new-feature`)
5. Create a new Pull Request
