# 上游同步策略 / Upstream sync

此 Fork 源自 `xyftkf/3x-ui`（當前基底 v2.6.2，2025-07-14），其上游源頭為 `xeefei/X-Panel`。

2026-11-12 比對發現，基底落後 X-Panel 約 968 個 commits。不要直接以 X-Panel 的 `main` 覆蓋本 Fork：中文介面、設定資料庫結構與自訂功能的差異需要逐項測試。

## 現階段完成

- 移除登入失敗密碼外洩。
- 新安裝使用隨機初始管理密碼。
- 建置 Go toolchain 更新至 1.27.4，Go module 最低需求同步更新為 1.27。
- Release / Docker 打包的 **外部** Xray 執行檔更新至 v26.11.9。
- 重建 GitHub Actions、GHCR 發布、Dependabot 與漏洞掃描。

## 尚待整合測試的上游同步

- `go.mod` 中的 `github.com/xtls/xray-core` 已升級至 2026-11-04 的上游提交，`go.sum` 已重新產生；Go CI 與 govulncheck 通過，但尚需對設定產生、JSON 格式、資料庫相容性、入站／出站進行實機整合測試。
- Gin、gRPC、HTTP middleware 與安全掃描指出的 8 個模組已更新並通過 build/test/vet/govulncheck；後續版本仍由 Dependabot 提出 PR，需通過 CI 與手動測試才合併。
- 上游 2026 版本引入的新協定、API/DB migrations、訂閱格式、安全功能需獨立移植，並對原本中文優化版逐項比對。
- 因缺少真實部署環境，不能宣稱新 Xray binary 與舊版 Go 模組的設定輸出已全面相容。

## 合併前驗收

1. GitHub Actions Go CI 與 release build 對 amd64 / arm64 都成功。
2. 新安裝可用初始帳密登入；舊資料庫可登入且設定原樣保存。
3. 入站/出站、新增/編輯節點、啟停 Xray、流量統計、訂閱頁面與 TLS 驗證成功。
4. 不應把真實密碼、Token、訂閱 URL 或 DB 檔案放進公開 CI logs。
