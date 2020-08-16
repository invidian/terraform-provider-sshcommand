module github.com/invidian/terraform-provider-sshcommand

go 1.12

// Workaround as this module was renamed
replace (
	github.com/golang/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422
	// https://github.com/sourcegraph/go-diff/issues/34
	sourcegraph.com/sourcegraph/go-diff v0.5.1 => github.com/sourcegraph/go-diff v0.5.1
)

require (
	cloud.google.com/go v0.63.0 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-cidr v1.1.0 // indirect
	github.com/aws/aws-sdk-go v1.34.5 // indirect
	github.com/bmatcuk/doublestar v1.3.2 // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/hashicorp/go-hclog v0.14.1 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/hashicorp/go-retryablehttp v0.6.7 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/go-version v1.2.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/hil v0.0.0-20200423225030-a18a1cd20038 // indirect
	github.com/hashicorp/terraform v0.13.0
	github.com/hashicorp/terraform-svchost v0.0.0-20200729002733-f050f53b9734 // indirect
	github.com/hashicorp/yamux v0.0.0-20200609203250-aecfd211c9ce // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/spf13/afero v1.3.4 // indirect
	github.com/ulikunitz/xz v0.5.7 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de
	golang.org/x/net v0.0.0-20200813134508-3edf25e44fcc // indirect
	golang.org/x/sys v0.0.0-20200814200057-3d37ad5750ed // indirect
	google.golang.org/genproto v0.0.0-20200815001618-f69a88009b70 // indirect
)
