# terraform-qingcloud


Terraform-Qingcloud-Plugin [![Build Status](https://travis-ci.org/yunify/terraform-provider-qingcloud.svg?branch=master)](https://travis-ci.org/yunify/terraform-provider-qingcloud)
[![codebeat badge](https://codebeat.co/badges/4559529b-cb96-4120-a489-30ca998c3790)](https://codebeat.co/projects/github-com-yunify-terraform-provider-qingcloud-master)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fyunify%2Fterraform-provider-qingcloud.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fyunify%2Fterraform-provider-qingcloud?ref=badge_shield)

## Usage

### Install terraform-provider-qingcloud

To install Terraform, find the [appropriate package](https://github.com/yunify/terraform-provider-qingcloud/releases) for your system and download it. Terraform is packaged as a tgz archive.  
After downloading Terraform, unzip the package.   
On Linux or Mac ,Rename the single binary to `terraform-provider-qingcloud` , and put it to in the sub-path .terraform.d/plugins in your user's home directory.
On Windows , Rename the single binary to `terraform-provider-qingcloud.exe` , and put it to in the sub-path terraform.d/plugins beneath your user's "Application Data" directory.
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
+ provider.qingcloud (unversioned)
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



## Contributing

1. Fork it ( https://github.com/yunify/terraform-provider-qingcloud/fork )
2. Create your feature branch (`git checkout -b new-feature`)
3. Commit your changes (`git commit -asm 'Add some feature'`)
4. Push to the branch (`git push origin new-feature`)
5. Create a new Pull Request    


## Special Thanks
[CuriosityChina](https://github.com/CuriosityChina)
