package sshcommand

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider exports terraform-provider-sshcommand, which can be used in tests
// for other providers.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"sshcommand_command": resourceCommand(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sshcommand_command": dataSourceCommand(),
		},
	}
}
