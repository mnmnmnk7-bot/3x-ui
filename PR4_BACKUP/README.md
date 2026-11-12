# Pull Request #4 研究 / 還原備份

原 PR：https://github.com/mnmnmnk7-bot/3x-ui/pull/4

**Title:** security: harden 3x-ui sessions, TLS downgrade, 2FA and HTTP protections

**Status on backup:** open; Draft=true; Merged=false

**Source branch:** `hardening/self-hosted-security-audit-2026-11`

**Original head commit (exact):** `1c656e3cf0ace03902454dae00daeec3584d9206`

**Base commit at backup:** `4f92836f134f566af493926c671735f76cf706f3`

**Commit history:** 5 commits

**Changed files:** 14 files

## 保存項目

- `PR4_BACKUP/pr-4.diff`：GitHub 當前提供的完整 unified diff。
- `PR4_BACKUP/manifest.json`：PR 標題、描述、全部 Commit SHA、檔案清單及討論紀錄。
- `backup/pr-4-security-hardening-2026-11-12`：**原始 Head Commit 的獨立固定分支**（沒有任何額外檔案）。
- 本 `backup/pr-4-archive-2026-11-12` 分支含原始 PR4 Head 的完整 Git 歷史，以及本備份說明與 diff。

## 還原

若 Repo 脫離 Fork Network，PR #4 的 PR 頁面可能無法保留，但兩個備份分支可繼續存取。務必保留分支，且**不要合併未驗證修補**。

還原原始安全修補分支：
```bash
git fetch origin backup/pr-4-security-hardening-2026-11-12
git switch -c restored-pr4 FETCH_HEAD
```

如果已有本地乾淨的基底，也可以先下載 `pr-4.diff`，用：
```bash
git apply --check PR4_BACKUP/pr-4.diff
git apply PR4_BACKUP/pr-4.diff
```
注意：Patch 是相對於備份時 GitHub PR diff 的比較基底；之後主分支若改動，可能發生衝突，優先還原完整 Git 分支。

## 檢驗事項

此 PR 仍處於 Draft，尚未完成 self-hosted Runner 的編譯及實際 VPS 的 Session / TLS / Xray / DB 整合測試。備份不代表已通過安全驗證。

PR 原說明：

---
## Static security audit + application hardening (DRAFT; not merged into main)

Self-hosted CI migration and safe trigger policy have already merged independently in PR #5 and PR #6; this PR only needs application changes and code review.

### Proposed changes
- Minimal signed session cookie (user ID + SHA256 credential-revision fingerprint, **no user record or bcrypt hash in cookie**).
- Old cookies require re-login; changing password invalidates previous sessions. Add regression test.
- HttpOnly / SameSite=Lax, Secure if direct TLS; cross-origin unsafe browser requests blocked when Origin is supplied.
- Default session signing key uses crypto/rand.
- Fail closed on configured panel/subscription TLS cert errors; header/idle timeouts and response security headers.
- Preserve 2FA when username/password is updated rather than silently disabling it.
- SECURITY_REVIEW_2026.md records prioritized remaining application/installer/runner risks.

### Required validation before merging
- [ ] Self-hosted runner online: `go test ./...`, `go vet ./...`, `go build ./...`, `govulncheck ./...` all pass
- [ ] Test login, cookie, password change, 2FA, direct HTTPS and TLS reverse proxy on a staging VPS
- [ ] Check Xray subscriptions and existing SQLite DB migration
- [ ] Back up deployment state; verify rollback
- README.md unchanged

**Keep as Draft until CI is runnable and validation is complete.**
