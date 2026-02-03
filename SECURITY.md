# Security

This document describes security-relevant behavior and recommendations when using cli-kit.

## Path validation (validator)

- **Path traversal**: `ValidatePath` rejects paths containing `..` both in the raw input and after normalization (`filepath.Clean`), to reduce bypass via encoding or platform quirks.
- **AllowedDirs**: When `AllowedDirs` is set, the resolved path must be exactly an allowed directory or under it; prefix tricks (e.g. `/tmpfoo` when `/tmp` is allowed) are rejected. The error message does not include the list of allowed directories to avoid information disclosure if the error is surfaced to untrusted parties.
- **Symlinks**: Symlinks are not resolved. Paths under an allowed directory may point outside it via symlinks. For strict containment, resolve symlinks at the call site or use OS-specific checks.

## URL validation (validator)

- **SSRF**: `ValidateURL` blocks private IPs and localhost by default. Use `URLOptions.AllowLocalhost` / `AllowPrivateIP` only when intentional.
- **Resolution**: When `ResolveHostTimeout` is set, hostnames are resolved and all resolved IPs are checked against the same rules.
- **Redirects**: Validation applies only to the URL as given. It does not protect against HTTP redirects (e.g. to private IPs) when the application later performs the request. Configure the HTTP client to restrict redirects or re-validate the resolved URL if needed.

## Configuration and environment (env, configutil)

- **Empty vs unset**: `env.Get` and `env.GetTrimmed` return the default when the variable is **not set or set to empty**. To tell “not set” from “set to empty”, use `env.Lookup` or `env.Has`. This matters for “must be set” or “empty means disable” semantics.
- **Sensitive values**: Avoid logging or error messages that include resolved config (e.g. URLs, paths, tokens). Prefer redaction in logs.
- **Environment variable keys**: Keys containing NUL (`\x00`) can have undefined or unsafe behavior on some systems. The `env` package does not validate keys; avoid passing user-controlled or untrusted strings as keys. Use testutil's `EnvManager` in tests, which rejects empty and NUL keys.

## Passwords (flagutil)

- **ReadPasswordFromFile**: The path is validated with traversal checks. The password is returned as a string and will remain in process memory; minimize copies and lifetime where possible.
- **TOCTOU**: There is a short window between path validation and reading the file. If an attacker can replace the path with a symlink in between (e.g. to another file), the read may target the new target. For highest assurance, use a dedicated directory with restricted permissions and no symlinks, or open files under a locked working directory.

## Validators (validator)

- **UsernameOptions.CustomPattern**: If the regex is built from user or external input, it can be vulnerable to ReDoS (catastrophic backtracking). Use only fixed, well-tested patterns or restrict pattern complexity when the source is untrusted.

## Test utilities (testutil)

- **EnvManager**: `Set`, `Unset`, and `SetMultiple` reject empty or NUL-containing environment variable keys to avoid unsafe or platform-dependent behavior.

## Reporting

If you believe you’ve found a security issue, please report it privately (e.g. via the maintainers or a private security channel) rather than in a public issue.
