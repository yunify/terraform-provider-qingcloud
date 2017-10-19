# terraform-qingcloud 


Terraform-Qingcloud-Plugin [![Build Status](https://travis-ci.org/yunify/qingcloud-terraform-provider.svg?branch=master)](https://travis-ci.org/yunify/qingcloud-terraform-provider)

## Usage

### Install qingcloud-provider

#### On Linux
``` bash
go get github.com/yunify/qingcloud-terraform-provider
glide up 
make build
make test
cp ./terraform-provider-qingcloud $(dirname `which terraform`)/terraform-provider-qingcloud
```

#### On Mac
``` bash
go get github.com/yunify/qingcloud-terraform-provider
glide up 
make build
make test
cp ./terraform-provider-qingcloud $(dirname `which terraform`)/terraform-provider-qingcloud
```

## Finish Resource：
- [ ] Instance
- [ ] Volume
- [ ] Vxnet
- [ ] Router
- [x] Eip
- [ ] SecurityGroups
- [ ] Keypairs
- [ ] Image
- [ ] LoadBalancer
- [ ] Tag
- [ ] redis
- [ ] mongodb


## Contributing

1. Fork it ( https://github.com/yunify/qingcloud-terraform-provider/fork )
2. Create your feature branch (`git checkout -b new-feature`)
3. Commit your changes (`git commit -asm 'Add some feature'`)
4. Push to the branch (`git push origin new-feature`)
5. Create a new Pull Request    


## Special Thanks
[CuriosityChina](https://github.com/CuriosityChina)
