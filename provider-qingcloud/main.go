package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/magicshui/terraform-qingcloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: qingcloud.Provider,
	})
}
