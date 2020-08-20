package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/yunify/terraform-provider-qingcloud/qingcloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: qingcloud.Provider,
	})
}
