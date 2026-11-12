# Self-hosted GitHub Actions runners

All workflow jobs use `runs-on: self-hosted`.

**This is a public repository.** Do not execute untrusted pull-request code on a persistent self-hosted runner. GitHub and OWASP both warn that a PR can modify the workflow itself, bypassing job-level guards.

## Trigger policy

- No workflow has the `pull_request` or `pull_request_target` trigger.
- Only pushes to protected `main`, matching release tags, manually approved `workflow_dispatch`, and the scheduled security scan can run.
- The Go security scan also runs for `main` pushes that modify Go dependencies.
- To validate an owner-created feature branch before merge, explicitly dispatch a workflow on that trusted branch. **Do not** dispatch unreviewed external PR code.
- This workflow policy does not replace repository Actions settings, branch protection, contributor approval, or runner-group access controls.

## Runner requirements

- Dedicated, disposable Ubuntu 24.04 x86-64 Linux VM; do not reuse a workstation or production host.
- Git, Bash, Go (or setup-go), Docker/Buildx/QEMU, curl, wget, tar, unzip, sudo and apt-get / cross-compilers.
- A single runner queues the seven architecture Release matrix jobs sequentially; the Docker build consumes significant disk and RAM.
- Release packaging uses Ubuntu `apt-get`, and Docker workflow needs access to the Docker daemon. Treat these as elevated privileges.
- Protect the default branch, tags and workflow file changes. Never share runner host SSH keys, unrestricted production networking, credentials, PATs, or SSH agents with jobs.
- Prefer ephemeral one-job runners in isolated VMs; destroy state, including Docker layers and volumes, after each job.
- If no self-hosted runner is online, workflows remain queued. Do not interpret this as a successful build or security check.

Official GitHub guidance: https://docs.github.com/en/actions/reference/security/secure-use
OWASP checklist: https://cheatsheetseries.owasp.org/cheatsheets/GitHub_Actions_Security_Cheat_Sheet.html
