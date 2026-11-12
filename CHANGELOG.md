# Changelog

Notable changes are documented below. Earlier releases are available in the repository tags and GitHub Releases.

## [1.9.4] - 2026-11-13

### Changed

- Require Go 1.27.4, which carries security fixes to cgo, the go command, the compiler, the linker and `net/http`.
- Update `github.com/spf13/pflag` from v1.0.10 to v1.1.0.

The module path and the public API are unchanged.

## [1.9.3] - 2026-11-01

### Changed

- Split URL scheme and DNS checks, IPv4 special-use policy, path containment, email-domain policy and username-style validation into smaller helpers.
- Preserve public APIs, validation order, errors and SSRF/path policy boundaries.
- Extract existing test cases into named helpers and simplify slice assertions; all six Go Report Card checks now score 100%.
- Pin Go Test Report v1.1.3 and Go Report Card v1.1.4; retain Go 1.27.3 and pflag v1.0.10.

