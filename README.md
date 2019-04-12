# SSH Command Terraform Provider [![Build Status](https://travis-ci.com/invidian/terraform-provider-sshcommand.svg?branch=master)](https://travis-ci.com/invidian/terraform-provider-sshcommand)

This provider allow to execute commands remotely via SSH and capture the output from them.

This provider uses built-in [crypto/ssh](https://godoc.org/golang.org/x/crypto/ssh) Golang library to act as a SSH Client. Currently the implementation is very limited, but can be easily extended.

## Table of contents
* [Requirements](#requirements)
* [Building](#building)
* [Installing the provider](#installing-the-provider)
* [Resources](#resources)
* [Example usage](#example-usage)
* [Authors](#authors)

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
-	[Go](https://golang.org/doc/install) 1.11.x

## Building

For testing builds, simply run `docker build .`, which will download all dependencies, run build, test and linter.

For local builds, simply follow the steps from [Dockerfile](https://github.com/invidian/terraform-provider-sshcommand/blob/master/Dockerfile).

## Installing the provider

After building the provider, install it using the Terraform instructions for [installing a third party provider](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).

## Resources

### sshcommand_command

This resource executes given command on remote system and stores it's output in Terraform.

#### Parameters
  - `host` - Hostname to connect.
  - `private_key` - SSH private key used for authentication (SSH Agent support is not implemented).
  - `command` - Command to execute.
  - `user` - User used for SSH log in. Default value is `root`.
  - `port` - Port to open SSH connection. Default is 22.
  - `connection_timeout` - Timeout for opening TCP connection. Default is `5m`.

#### Attributes
  - `result` - Output of executed command.

## Example usage
```hcl
provider "sshcommand" {
  version = "~> 0.1.0"
}

output "example" {
  value = "\n${sshcommand_command.ssh_host_fingerprints.result}"
}

resource "sshcommand_command" "ssh_host_fingerprints" {
  host               = "example"
  command            = "ssh-keygen -r $(hostname -f) | cut -d' ' -f4-6"
  private_key        = "${file(".ssh/id_rsa")}"
}
```

## Authors
* **Mateusz Gozdek** - *Initial work* - [invidian](https://github.com/invidian)
