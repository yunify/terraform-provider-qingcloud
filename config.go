package qingcloud

import (
	"github.com/magicshui/qingcloud-go"
	"github.com/magicshui/qingcloud-go/eip"
	"github.com/magicshui/qingcloud-go/keypair"
	"github.com/magicshui/qingcloud-go/securitygroup"
	"github.com/magicshui/qingcloud-go/vxnet"
)

type Config struct {
	ID     string
	Secret string
	Zone   string
}

type QingCloudClient struct {
	eip           *eip.EIP
	keypair       *keypair.KEYPAIR
	securitygroup *securitygroup.SECURITYGROUP
	vxnet         *vxnet.VXNET
}

func (c *Config) Client() (*QingCloudClient, error) {
	clt := qingcloud.NewClient()
	clt.ConnectToZone(c.Zone, c.ID, c.Secret)

	return &QingCloudClient{
		eip:           eip.NewClient(clt),
		keypair:       keypair.NewClient(clt),
		securitygroup: securitygroup.NewClient(clt),
		vxnet:         vxnet.NewClient(clt),
	}, nil
}
