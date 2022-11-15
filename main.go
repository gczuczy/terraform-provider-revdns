package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/gczuczy/terraform-provider-revdns/revdns"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: revdns.New,
	})
}
