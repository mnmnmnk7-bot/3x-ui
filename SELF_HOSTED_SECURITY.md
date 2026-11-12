# 自架 GitHub Actions Runner 安全與環境需求

本 Repository 的所有 GitHub Actions job 都使用 `runs-on: self-hosted`。

## Runner 環境

- 建議專用、可還原快照的 Ubuntu 24.04 x86-64 Linux VM。
- 必備 Git、Bash、Go 1.27.4（或允許 actions/setup-go 安裝）、Docker、Buildx、QEMU、curl、wget、tar、unzip、sudo。
- Release 建置使用 `apt-get` 安裝 gcc 系列跨平台編譯工具，因此不能直接將 Windows、macOS 或非 Debian 系 Runner 混用同一個 `self-hosted` 標籤。
- 7 個架構的 Release 使用 matrix。只有一台 Runner 時會依序排隊；Docker Buildx 會使用高資源。
- 如 Runner 未連線，工作會持續保持 queued，無法替代為 GitHub-hosted runner。

## 安全限制

- **禁止公開 Fork PR 直接執行於自架 Runner**：PR job 僅允許來自本 Repository、且由 Repository Owner 建立的 PR。其他來源會跳過。
- 僅依靠工作流程中的 `if` 無法防範擁有 Repository 寫入權限或篡改 workflow 權限的人。保護 `main` 並限制 workflow 檔案的修改權。
- 不應將自架 Runner 架在個人日常工作站、存有 SSH 私鑰的電腦或具生產網路直接存取權的機器。
- Docker socket、QEMU/binfmt、sudo 及安裝系統套件均屬高權限操作。建議每次執行後重建 ephemeral Runner VM，而不是共用長期運行的工作站。
- 不要讓 Runner 使用常駐的個人 PAT、雲端管理憑證或部署環境 Secret。
- 第三方 Actions、依賴套件和遠端下載都是供應鏈風險；應限制來源並定期檢查。

此文件只描述 Runner 的部署與安全政策，不代表其他應用程式安全修補已經合併。
