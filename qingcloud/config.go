package qingcloud

import (
	"github.com/yunify/qingcloud-sdk-go/config"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

type Config struct {
	ID     string
	Secret string
	Zone   string
}

type QingCloudClient struct {
	zone          string
	eip           *qc.EIPService
	keypair       *qc.KeyPairService
	securitygroup *qc.SecurityGroupService
	vxnet         *qc.VxNetService
	router        *qc.RouterService
	instance      *qc.InstanceService
	volume        *qc.VolumeService
	loadbalancer  *qc.LoadBalancerService
	tag           *qc.TagService
	cahce         *qc.CacheService
	mongo         *qc.MongoService
}

func (c *Config) Client() (*QingCloudClient, error) {
	cfg, err := config.New(c.ID, c.Secret)
	if err != nil {
		return nil, err
	}
	clt, err := qc.Init(cfg)
	if err != nil {
		return nil, err
	}
	eip, err := clt.EIP(c.Zone)
	if err != nil {
		return nil, err
	}
	keypair, err := clt.KeyPair(c.Zone)
	if err != nil {
		return nil, err
	}
	securitygroup, err := clt.SecurityGroup(c.Zone)
	if err != nil {
		return nil, err
	}
	vxnet, err := clt.VxNet(c.Zone)
	if err != nil {
		return nil, err
	}
	router, err := clt.Router(c.Zone)
	if err != nil {
		return nil, err
	}
	instance, err := clt.Instance(c.Zone)
	if err != nil {
		return nil, err
	}
	volume, err := clt.Volume(c.Zone)
	if err != nil {
		return nil, err
	}
	tag, err := clt.Tag(c.Zone)
	if err != nil {
		return nil, err
	}
	loadbalancer, err := clt.LoadBalancer(c.Zone)
	if err != nil {
		return nil, err
	}
	cache, err := clt.Cache(c.Zone)
	if err != nil {
		return nil, err
	}
	mongo, err := clt.Mongo(c.Zone)
	if err != nil {
		return nil, err
	}

	return &QingCloudClient{
		zone:          c.Zone,
		eip:           eip,
		keypair:       keypair,
		securitygroup: securitygroup,
		vxnet:         vxnet,
		router:        router,
		instance:      instance,
		volume:        volume,
		loadbalancer:  loadbalancer,
		tag:           tag,
		cahce:         cache,
		mongo:         mongo,
	}, nil
}
