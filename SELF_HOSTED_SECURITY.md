# Self-hosted GitHub Actions runner policy

This public fork uses `runs-on: self-hosted` for **every job**.

## Required runner hardening

- Use a dedicated, disposable Linux **x86-64 Ubuntu 24.04** host / VM. Never attach a runner to a workstation containing credentials, personal files, or production VPS data.
- Provision Go 1.27.4, Docker + Buildx / QEMU, curl, unzip, tar, bash, sudo, Git and cross-compilers (or allow Ubuntu package installation). A runner without these tools will fail the jobs.
- The Release workflow runs **seven architecture jobs**. With one runner they will queue and run sequentially.
- Do **not** enable automatic fork-PR execution. Each PR job has an explicit same-repository source guard. This reduces exposure, but does **not** make persistent self-hosted runners safe against malicious changes made by trusted collaborators.
- Protect `main`, restrict write access to workflow files and runner registration, require reviews for workflow changes, and avoid exposing secrets / cloud credentials to build jobs.
- Docker uses host-level privileges. Assume a malicious build can compromise a persistent self-hosted runner. Prefer one-shot ephemeral VMs with snapshot reset or container isolation per job.
- The self-hosted runner may be offline. Jobs will remain queued until a compatible runner connects; this is not a source-code build failure.
- Do not attach runners to networks with management interfaces, SSH agents, API keys or unrestricted access to production machines.

## Static security audit — 2026-11-12

**Changes in this maintenance branch:**
- Signed cookie session no longer stores the full user record or password hash. It holds a user ID plus a credential-revision fingerprint; the user is fetched server-side and password resets invalidate old cookies.
- Session cookies are HttpOnly and SameSite=Lax, with Secure enabled for direct HTTPS requests; terminate TLS carefully if behind a reverse proxy.
- Cross-origin POST / other state-changing browser requests with an Origin header are rejected by authenticated panel routes.
- The default session signing key now comes from the cryptographic random source.
- An invalid configured TLS certificate/key now stops the HTTP or subscription listener instead of silently downgrading to plaintext HTTP.
- HTTP server sets ReadHeaderTimeout / IdleTimeout and additional conservative response headers.

**Unresolved / operational risks requiring review:**
- `install.sh` and menu actions execute remote downloaded code; packages depend on GitHub release publication. Verify repository identity, package checksums and supply-chain trust.
- The Xray update endpoint downloads and installs Xray executables from GitHub without detached-signature verification; only authenticated administrators may trigger it, but compromises to an admin session are severe.
- The panel login has no dedicated rate limiter; restrict management to VPN/private network or deploy a trusted reverse proxy with rate limiting.
- `X-Real-IP`, `X-Forwarded-For`, and `X-Forwarded-Host` are used in some paths without a configured trusted proxy list. Treat them as untrusted unless a reverse proxy sanitizes them.
- The database import endpoint accepts user-supplied SQLite and may consume substantial memory/disk; isolate the admin interface and back up DB files before import.
- The application's TLS-reverse-proxy configuration has not been verified on a live host, and third-party frontend assets are vendored.
- Govulncheck's prior scan reported **zero reachable vulnerabilities**, not an absence of architectural or configuration vulnerabilities.

Back up `/etc/x-ui` before deployment. Confirm login/logout, password changes, old-session invalidation, HTTPS, reverse proxy, subscription URLs and inbound/outbound Xray traffic in a staging instance.
