package main

import (
	"github.com/pvotal-tech/terraform-provider-data-utils/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {

	opts := &plugin.ServeOpts{ProviderFunc: provider.New()}

	plugin.Serve(opts)
}
