package main

import (
	"github.com/CuriosityChina/terraform-qingcloud"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: qingcloud.Provider,
	})
}
