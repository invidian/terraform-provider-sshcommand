# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.2] - 2020-11-22

### Added

- `sshcommand_command` resource and data source now supports `password` parameter for authentication.

## [0.2.1] - 2020-08-18

### Changed

- Migrated to Terraform Plugin SDK

## [0.2.0] - 2020-08-18

### Added

- Added `sshcommand_command` data source, which acts the same as `sshcommand_command resource.
- Codecov and Code Climate checks will now run as part of CI process for master branch and pull requests.
- Provider is now available in [Terraform Registy](https://registry.terraform.io/).

### Changed

- `goreleaser` will be now used for creating the releases
- Updated all dependencies
- Improved error phrasing

### Fixed

- Removing the resource should now be correctly removed from Terraform state.

## [0.1.2] - 2019-05-24

### Added

- Added 'retry', 'retry_timeout' and 'retry_interval' parameters and logic.
- Added `Makefile` with common tasks

### Changed

- Updated dependencies

## [0.1.1] - 2019-05-11

### Added

- Added 'ignore_execute_errors' parameter. If this is set to 'true' and executed command returns error, it will be ignored.

## [0.1.0] - 2019-05-09

### Added

- Initial release

[0.1.2]: https://github.com/invidian/terraform-provider-sshcommand/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/invidian/terraform-provider-sshcommand/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/invidian/terraform-provider-sshcommand/releases/tag/v0.1.0
