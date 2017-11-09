# Change Log
All notable changes to QingCloud SDK for Go will be documented in this file.

## [v2.0.0-alpha.10] - 2017-08-28

### Fixed

- Fixed loadbalancers section in request of StopLoadBalancers

## [v2.0.0-alpha.9] - 2017-08-23

### Added

- Add vxnetid and loadbalancertype parms for load balancer

### Fixed

- Fixed vxnet section in response of DescribeInstances

## [v2.0.0-alpha.8] - 2017-08-13

### Added

- Add missing parameter for describenic
- Add vxnet parms for instances

## [v2.0.0-alpha.7] - 2017-08-02

### Added

- Add timeout parameter for http client
- Add missing parameters in nic, router and security groups

## [v2.0.0-alpha.6] - 2017-07-17

### Added

- Update advanced client. [@jolestar]
- Fix several APIs. [@jolestar]

## [v2.0.0-alpha.5] - 2017-03-27

### Added

- Add advanced client for simple instance management. [@jolestar]
- Add wait utils for waiting a job to finish. [@jolestar]

## [v2.0.0-alpha.4] - 2017-03-14

### Fixed

- Fix Features type in RouterVxNet.

## [v2.0.0-alpha.3] - 2017-01-15

### Changed

- Fix request signer.

## [v2.0.0-alpha.2] - 2017-01-05

### Changed

- Fix logger output format, don't parse special characters.
- Rename package "errs" to "errors".

### Added

- Add type converters.

### BREAKING CHANGES

- Change value type in input and output to pointer.

## v2.0.0-alpha.1 - 2016-12-03

### Added

- QingCloud SDK for the Go programming language.
[v2.0.0-alpha.10]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.9...v2.0.0-alpha.10
[v2.0.0-alpha.9]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.8...v2.0.0-alpha.9
[v2.0.0-alpha.8]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.7...v2.0.0-alpha.8
[v2.0.0-alpha.7]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.6...v2.0.0-alpha.7
[v2.0.0-alpha.6]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.5...v2.0.0-alpha.6
[v2.0.0-alpha.5]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.4...v2.0.0-alpha.5
[v2.0.0-alpha.4]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.3...v2.0.0-alpha.4
[v2.0.0-alpha.3]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.2...v2.0.0-alpha.3
[v2.0.0-alpha.2]: https://github.com/yunify/qingcloud-sdk-go/compare/v2.0.0-alpha.1...v2.0.0-alpha.2

[@jolestar]: https://github.com/jolestar
