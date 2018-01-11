# terraform-qingcloud


Terraform-Qingcloud-Plugin [![CircleCI](https://circleci.com/gh/yunify/terraform-provider-qingcloud/tree/master.svg?style=svg)](https://circleci.com/gh/yunify/terraform-provider-qingcloud/tree/master)
[![codebeat badge](https://codebeat.co/badges/d6cc83ea-779f-4fce-8091-abc0b719d271)](https://codebeat.co/projects/github-com-yunify-qingcloud-terraform-provider-master-3c5cd450-e81b-4eb1-aaf6-aa9b76158d6f)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fyunify%2Fterraform-provider-qingcloud.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fyunify%2Fterraform-provider-qingcloud?ref=badge_shield)

## Usage

### Install terraform-provider-qingcloud

To install Terraform, find the [appropriate package](https://github.com/yunify/terraform-provider-qingcloud/releases) for your system and download it. Terraform is packaged as a tgz archive.  
After downloading Terraform, unzip the package.   
On Linux or Mac , put the binary file to in the sub-path .terraform.d/plugins in your user's home directory.
On Windows , put the binary file to in the sub-path terraform.d/plugins beneath your user's "Application Data" directory.
Then put the binary file into terraform 's PATH.

### Verifying the Installation

```shell
git clone https://github.com/yunify/terraform-provider-qingcloud.git
cd ./terraform-provider-qingcloud/terraform/example/init
terraform init
terraform -v
```
You can execute the above script . If you installed the provider correctly, you should see output similar to the one below .  
```shell
Terraform v0.11.1
+ provider.qingcloud (v1.1)
```

## Finish Resource:
- [x] Instance
- [x] Volume
- [x] Vxnet
- [ ] Router(Deprecated,Use Vpc in SDN2.0)
- [x] Eip
- [x] SecurityGroups
- [x] SecurityGroupRules
- [x] Keypairs
- [x] Vpc
- [x] Tag
- [x] VpcStatic
- [x] LoadBalancer
- [x] LoadBalancerListener
- [x] LoadBalancerBackend
- [x] Server Certificate
- [x] VPN Cert



## Contributing

1. Fork it ( https://github.com/yunify/terraform-provider-qingcloud/fork )
2. Create your feature branch (`git checkout -b new-feature`)
3. Commit your changes (`git commit -asm 'Add some feature'`)
4. Push to the branch (`git push origin new-feature`)
5. Create a new Pull Request    


## Special Thanks
[CuriosityChina](https://github.com/CuriosityChina)
