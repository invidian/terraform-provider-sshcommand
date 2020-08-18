package sshcommand

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceCommand() *schema.Resource {
	return &schema.Resource{
		Read:   resourceCommandCreate,
		Schema: resourceCommandSchema(),
	}
}
