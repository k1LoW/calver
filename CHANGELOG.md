# Changelog

## [v1.0.2](https://github.com/k1LoW/calver/compare/v1.0.1...v1.0.2) - 2026-02-03
- chore(deps): bump Songmu/tagpr from 1.12.1 to 1.15.0 in the dependencies group by @dependabot[bot] in https://github.com/k1LoW/calver/pull/40

## [v1.0.1](https://github.com/k1LoW/calver/compare/v1.0.0...v1.0.1) - 2026-01-29
- chore(deps): bump the dependencies group with 2 updates by @dependabot[bot] in https://github.com/k1LoW/calver/pull/37
- feat: support parsing variable-length tokens followed by fixed-length  tokens by @k1LoW in https://github.com/k1LoW/calver/pull/39

## [v1.0.0](https://github.com/k1LoW/calver/compare/v0.7.4...v1.0.0) - 2026-01-13

## [v0.7.4](https://github.com/k1LoW/calver/compare/v0.7.3...v0.7.4) - 2026-01-13
- chore: setup tagpr labels by @k1LoW in https://github.com/k1LoW/calver/pull/31
- chore(deps): bump the dependencies group with 3 updates by @dependabot[bot] in https://github.com/k1LoW/calver/pull/33
- chore(deps): bump the dependencies group across 1 directory with 3 updates by @dependabot[bot] in https://github.com/k1LoW/calver/pull/35

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
