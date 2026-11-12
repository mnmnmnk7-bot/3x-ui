# 3x-ui security re-review (2026-11-12)

This is a static code review, not a completed penetration test. The existing README remains unchanged.

## Hardening in this draft branch, not yet deployed
- Cookie session no longer contains the full user record / bcrypt hash; uses a minimal ID and credential-revision fingerprint. A password change invalidates other sessions; existing cookies require re-login.
- HttpOnly + SameSite=Lax cookie flags; Secure on direct HTTPS. Behind a reverse proxy, enforce HTTPS and verify cookie security.
- Authenticated routes check cross-origin request Origin when supplied.
- Cryptographically secure default session signing key via `crypto/rand.Text()`.
- Both panel and subscription listeners fail closed when configured TLS credentials cannot be loaded, instead of falling back to plaintext HTTP.
- Response headers (nosniff, deny iframe, no-referrer) and HTTP ReadHeaderTimeout / IdleTimeout limits.
- Password/username updates preserve the existing 2FA requirement.
- Regression test for minimal session tokens and password-change revocation.

## Open risks not fixed in this draft
- **High:** self-hosted runners on a *public* repo can be compromised by untrusted code. Workflows now have NO pull_request triggers, but use branch protections, restrict write permissions, and use disposable isolation.
- **Medium-high:** `install.sh` downloads release tarballs as root without checking published SHA256SUMS or signed provenance before replacing the installation. Verify and back up prior to install.
- **Medium-high:** `web/service/server.go` can download and replace Xray executables via authenticated admin endpoints without a verified signature; administrators should not accept arbitrary unreviewed versions.
- **Medium:** optional menu actions in `x-ui.sh` use downloaded shell code (`curl | bash`, etc.). Do not run without auditing those upstream scripts.
- **Medium:** panel login lacks a dedicated rate limiter; isolate access by VPN/ACL or trusted reverse-proxy rate limits.
- **Medium:** `getRemoteIp` and subscription host derivation trust forwarded headers without an explicit trusted-proxy allowlist; external callers can spoof header values.
- **Medium:** authenticated SQLite DB import lacks a clear upload size bound; may consume excessive memory/storage and needs a robust rollback test.
- **Medium:** vendored JS and unpinned third-party GitHub Actions require separate supply-chain review.
- A previous `govulncheck` reported zero *reachable* vulnerabilities, not proof that no security bugs exist.

## Acceptance
Run Go tests, vet, build and govulncheck on an isolated self-hosted runner; test login, logout, existing cookies, password reset, 2FA, TLS, reverse proxy, subscriptions and Xray traffic on a staging VPS. Back up `/etc/x-ui` before deployment.

Official self-hosted security guidance:
https://docs.github.com/en/actions/reference/security/secure-use
