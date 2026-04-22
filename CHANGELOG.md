# Changelog

All notable changes to `gxsrvs/dtx` are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.1] — 2026-04-22

### Added

- `Unix()` on `OffsetDateTime` and `NullOffsetDateTime`. Returns the
  underlying instant as a Unix timestamp (seconds since 1970-01-01
  UTC); the nullable variant returns `0` for NULL.

### Fixed

- `NullDecimal` and `NullUuid` now implement
  `database/sql/driver.Valuer`. Previously these two types satisfied
  only `sql.Scanner`, which caused drivers that rely on `Valuer` to
  serialise query parameters (notably `pgx` in simple-query mode) to
  fail when a `NullDecimal` or `NullUuid` was passed in. NULL emits
  `(nil, nil)`; valid values delegate to `decimal.Decimal.Value` and
  `uuid.UUID.Value` respectively.

## [0.1.0] — 2026-04-19

Initial public release.

[Unreleased]: https://github.com/gxsrvs/dtx/compare/v0.1.1...HEAD
[0.1.1]: https://github.com/gxsrvs/dtx/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/gxsrvs/dtx/releases/tag/v0.1.0
