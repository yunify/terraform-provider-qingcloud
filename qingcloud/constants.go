package qingcloud

const (
	qingcloudResourceTypeInstance      = "instance"
	qingcloudResourceTypeVolume        = "volume"
	qingcloudResourceTypeKeypair       = "keypair"
	qingcloudResourceTypeSecurityGroup = "security_group"
	qingcloudResourceTypeVxNet         = "vxnet"
	qingcloudResourceTypeEIP           = "eip"
	qingcloudResourceTypeRouter        = "router"
	qingcloudResourceTypeLoadBalancer  = "loadbalancer"

	StatusActive = "active"

	DEFAULT_ZONE           = "pek3a"
	DEFAULT_ENDPOINT       = "https://api.qingcloud.com:443/iaas"
	waitJobTimeOutDefault  = 240
	waitJobIntervalDefault = 5
	waitLeaseSecond        = 30

	resourceName        = "name"
	resourceDescription = "description"
	resourceTagIds      = "tag_ids"
	resourceTagNames    = "tag_names"
	DEFAULT_TAG_COLOR   = "#9f9bb7"
	BasicNetworkID      = "vxnet-0"
)
