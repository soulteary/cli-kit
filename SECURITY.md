# Security

This document describes security-relevant behavior and recommendations when using cli-kit.

## Path validation (validator)

- **Path traversal**: `ValidatePath` rejects paths containing `..` both in the raw input and after normalization (`filepath.Clean`), to reduce bypass via encoding or platform quirks.
- **AllowedDirs**: When `AllowedDirs` is set, the resolved path must be exactly an allowed directory or under it; prefix tricks (e.g. `/tmpfoo` when `/tmp` is allowed) are rejected. The error message does not include the list of allowed directories to avoid information disclosure if the error is surfaced to untrusted parties.
- **Symlinks**: Existing paths are resolved with `EvalSymlinks` before policy checks. This blocks “allowed-dir symlink escape” cases where a path appears to be under an allowed directory but points elsewhere. For non-existent paths, symlink resolution cannot be fully determined yet.

## URL validation (validator)

- **SSRF**: `ValidateURL` blocks private/internal IPs and localhost by default, including metadata/link-local ranges such as `169.254.0.0/16`.
- **Always-blocked addresses**: Unspecified and multicast targets are rejected even if `AllowPrivateIP` is enabled.
- **Userinfo**: URLs containing userinfo (for example `http://user:pass@host`) are rejected to reduce confusion and credential-leak risk.
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
