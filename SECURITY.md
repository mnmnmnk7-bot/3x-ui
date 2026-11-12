# 3x-ui 安全維護說明

## 重要變更（2026-11）

- 本 Fork 使用 Go 的 HTTP Server / Gin，不使用 Apache httpd。若自行部署 Apache 反向代理，需獨立檢查 Apache 與 mod_http2 的安全更新。
- 新資料庫第一次啟動時，管理帳號為 `admin`，密碼以 `crypto/rand` 隨機產生，不再使用 `admin/admin`。**只在第一次初始化的伺服器日誌中輸出密碼**；請保護日誌存取權限。
- 已有資料庫的帳號與密碼不會被修改。安裝腳本若自行設定帳密，以安裝腳本結果為準。
- 登入失敗不再寫入輸入密碼，也不再將輸入密碼送進 Telegram 登入通知。
- CI 建置使用修補後的 Go 1.27.4，更新後的 Go module 最低需求為 1.27。
- 發布包所附 Xray 執行檔改為 v26.11.9；這與 `go.mod` 裡的 Xray library 更新是不同工作。更新後的 Xray Go module 與外部 Xray 執行檔需要在正式部署前進行設定相容性與入站／出站整合測試。

## 部署前檢查

1. 備份 `/etc/x-ui/` 中的資料庫與設定，保管憑證私鑰。
2. 對更新後的面板測試既有入站、出站、訂閱 URL、TLS、帳號權限與 Telegram 通知。
3. 只透過 HTTPS 或受信任的 VPN/反向代理開放管理介面；避免公開暴露管理 Port。
4. 避免保留初始密碼；可利用 `x-ui setting` 或選單重設。
5. 對每次依賴與 Xray 更新先查閱 Actions CI、govulncheck 及 Release 備註。

## 安全回報

請勿在公開 Issue 附上訂閱連結、資料庫、憑證、API Token 或真實登入資訊。
