package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/nazarewk-iac/terraform-provider-custom/internal/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: provider.New})
}
