package sshcommand

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider exports terraform-provider-sshcommand, which can be used in tests
// for other providers.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"sshcommand_command": resourceCommand(),
		},
	}
}
