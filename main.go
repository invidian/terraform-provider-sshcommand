package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/invidian/terraform-provider-sshcommand/sshcommand"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sshcommand.Provider})
}
