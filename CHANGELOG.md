# Changelog

## [v0.7.4](https://github.com/k1LoW/calver/compare/v0.7.3...v0.7.4) - 2025-10-23
- chore: setup tagpr labels by @k1LoW in https://github.com/k1LoW/calver/pull/31

## [v0.7.3](https://github.com/k1LoW/calver/compare/v0.7.2...v0.7.3) - 2024-06-03
- Add Linux arm64 by @k1LoW in https://github.com/k1LoW/calver/pull/29

## [v0.7.2](https://github.com/k1LoW/calver/compare/v0.7.1...v0.7.2) - 2023-08-01
- If the time version is updated in the next version and the first token in the layout is the time version, reset the major/minor/micro version. by @k1LoW in https://github.com/k1LoW/calver/pull/26

## [v0.7.1](https://github.com/k1LoW/calver/compare/v0.7.0...v0.7.1) - 2023-05-15

## [v0.7.0](https://github.com/k1LoW/calver/compare/v0.6.0...v0.7.0) - 2023-05-15
- Add `--modifier` by @k1LoW in https://github.com/k1LoW/calver/pull/20
- Trim trailing zero value version even with modifier set. by @k1LoW in https://github.com/k1LoW/calver/pull/22
- Support for parsing value with `--trim-suffix` by @k1LoW in https://github.com/k1LoW/calver/pull/23
- Fix for modifier by @k1LoW in https://github.com/k1LoW/calver/pull/24

## [v0.6.0](https://github.com/k1LoW/calver/compare/v0.5.1...v0.6.0) - 2023-05-13
- Support for trimming the trailing version of a zero value or an empty string. by @k1LoW in https://github.com/k1LoW/calver/pull/18

## [v0.5.1](https://github.com/k1LoW/calver/compare/v0.5.0...v0.5.1) - 2023-05-13
- Fix sort logic by @k1LoW in https://github.com/k1LoW/calver/pull/15
- Fix error handling of layout by @k1LoW in https://github.com/k1LoW/calver/pull/17

## [v0.5.0](https://github.com/k1LoW/calver/compare/v0.4.0...v0.5.0) - 2023-05-12
- Add Sort() and Latest() by @k1LoW in https://github.com/k1LoW/calver/pull/12
- Parse multi versions and use latest by @k1LoW in https://github.com/k1LoW/calver/pull/13
- Add `--major`, `--minor`, and `--micro` option by @k1LoW in https://github.com/k1LoW/calver/pull/14

## [v0.4.0](https://github.com/k1LoW/calver/compare/v0.3.0...v0.4.0) - 2023-05-10
- The `--next` option is only active if there is a version to parse by @k1LoW in https://github.com/k1LoW/calver/pull/10

## [v0.3.0](https://github.com/k1LoW/calver/compare/v0.2.0...v0.3.0) - 2023-05-10
- Unexport some functions by @k1LoW in https://github.com/k1LoW/calver/pull/8

## [v0.2.0](https://github.com/k1LoW/calver/compare/v0.1.0...v0.2.0) - 2023-05-09
- Add Layout by @k1LoW in https://github.com/k1LoW/calver/pull/4
- Add command `calver` by @k1LoW in https://github.com/k1LoW/calver/pull/5
- Fix release settings by @k1LoW in https://github.com/k1LoW/calver/pull/6

## [v0.0.1](https://github.com/k1LoW/calver/commits/v0.0.1) - 2023-05-09
- Add NewWithTime by @k1LoW in https://github.com/k1LoW/calver/pull/2
