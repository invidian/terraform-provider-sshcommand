package sshcommand_test

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/invidian/terraform-provider-sshcommand/sshcommand"
)

func TestProvider(t *testing.T) {
	if err := sshcommand.Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
