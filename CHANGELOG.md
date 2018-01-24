## 1.0.0 (November 29, 2017)

FEATURES:

* **New Resource**: `qingcloud_eip` ([#21](https://github.com/yunify/terraform-provider-qingcloud/issues/21))
* **New Resource**: `qingcloud_tag` ([#24](https://github.com/yunify/terraform-provider-qingcloud/issues/24))
* **New Resource**: `qingcloud_security_group` ([#27](https://github.com/yunify/terraform-provider-qingcloud/issues/27))
* **New Resource**: `qingcloud_security_group_rule` ([#30](https://github.com/yunify/terraform-provider-qingcloud/issues/30))
* **New Resource**: `qingcloud_keypair` ([#31](https://github.com/yunify/terraform-provider-qingcloud/issues/31))
* **New Resource**: `qingcloud_vpc` ([#67](https://github.com/yunify/terraform-provider-qingcloud/issues/67))
* **New Resource**: `qingcloud_instance` ([#70](https://github.com/yunify/terraform-provider-qingcloud/issues/70))
* **New Resource**: `qingcloud_vxnet` ([#72](https://github.com/yunify/terraform-provider-qingcloud/issues/72))
* **New Resource**: `qingcloud_volume` ([#75](https://github.com/yunify/terraform-provider-qingcloud/issues/75))
* **New Resource**: `qingcloud_vpc_static` ([#98](https://github.com/yunify/terraform-provider-qingcloud/issues/98))

## 1.0.1 (December 14,, 2017)

IMPROVEMENTS:

* provider : Add paramertes to set alternative api endpoint ([#106](https://github.com/yunify/terraform-provider-qingcloud/issues/106))
* resource/qingcloud_eip add importer ([#105](https://github.com/yunify/terraform-provider-qingcloud/issues/105))
* resource/qingcloud_instance add importer ([#105](https://github.com/yunify/terraform-provider-qingcloud/issues/105))
* resource/qingcloud_keypair add importer ([#105](https://github.com/yunify/terraform-provider-qingcloud/issues/105))
* resource/qingcloud_tag add importer ([#105](https://github.com/yunify/terraform-provider-qingcloud/issues/105))
* resource/qingcloud_volume add importer ([#105](https://github.com/yunify/terraform-provider-qingcloud/issues/105))
* resource/qingcloud_vxnet add importer ([#105](https://github.com/yunify/terraform-provider-qingcloud/issues/105))

## 1.2.0 (January 11, 2018)

FEATURES:

* **New Resource**: `qingcloud_loadbalancer`([#115](https://github.com/yunify/terraform-provider-qingcloud/pull/115))
* **New Resource**: `qingcloud_server_certificate`([#128](https://github.com/yunify/terraform-provider-qingcloud/pull/128))
* **New Resource**: `qingcloud_loadbalancer_listener`([#129](https://github.com/yunify/terraform-provider-qingcloud/pull/129))
* **New Resource**: `qingcloud_loadbalancer_backend`([#132](https://github.com/yunify/terraform-provider-qingcloud/pull/132))
* **New Data Source**: `qingcloud_vpn_cert`([#134](https://github.com/yunify/terraform-provider-qingcloud/pull/134))

IMPROVEMENTS:

* provider: Add version info to help debug.
* resource/qingcloud_instance Add instance userdata ([#133](https://github.com/yunify/terraform-provider-qingcloud/pull/115))
* provider: wait resource create finshed time ([#126](https://github.com/yunify/terraform-provider-qingcloud/pull/126))
* provider: Add ACC test resource tag ([#125](https://github.com/yunify/terraform-provider-qingcloud/pull/125))
* provider: Use circle ci to run test and release

## 1.2.1 (January 24, 2018)

IMPROVEMENTS:

* release: Optimize packaging type ([#145](https://github.com/yunify/terraform-provider-qingcloud/pull/145))
* docs: Add docs to use resource keypair ([#147](https://github.com/yunify/terraform-provider-qingcloud/pull/147))
* resource/qingcloud_instance Add wait job for resize instance ([#148](https://github.com/yunify/terraform-provider-qingcloud/pull/148))
* resource/qingcloud_vpc Add has been deleted status ([#149](https://github.com/yunify/terraform-provider-qingcloud/pull/149))

