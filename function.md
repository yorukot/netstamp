# Netstamp 功能清單

> `Probe` 節點代表安裝在 VPS、伺服器或內部主機上的量測 agent。

## 1. 前端功能

### 1.1 使用者帳號與登入

* 註冊帳號
  * 輸入 email
  * 輸入密碼
  * 確認密碼
  * 顯示註冊成功或錯誤訊息
  * email 重複時顯示錯誤
* 登入
  * 輸入 email
  * 輸入密碼
  * 登入成功後進入主要 dashboard
  * 登入失敗時顯示錯誤
  * token 過期後導回登入頁
* 登出
  * 使用者可從帳號選單登出
  * 登出後清除本機登入狀態
* 目前使用者資訊
  * 顯示目前登入 email
  * 顯示目前所屬組織 / team
  * 顯示目前角色，例如 owner、admin、member

---

### 1.2 初次使用與建立組織

* 初次登入 onboarding
  * 如果使用者尚未加入任何組織，導向建立組織流程
* 建立組織
  * 輸入組織名稱
  * 輸入或自動產生 slug
  * 建立後自動成為 owner
* 組織列表
  * 顯示使用者所屬的所有組織
  * 可切換目前操作的組織
* 組織詳細頁
  * 顯示組織名稱
  * 顯示 slug
  * 顯示建立時間
  * 顯示成員數量
  * 顯示節點數量
  * 顯示 check 數量
* 編輯組織
  * 修改組織名稱
  * 修改 slug
  * 顯示儲存成功 / 失敗狀態
* 組織成員管理
  * 查看成員列表
  * 顯示成員 email
  * 顯示成員角色
  * 新增成員
  * 修改成員角色
  * 支援角色：owner、admin、member
  * MVP 可先不做細緻權限差異，但 UI 仍需顯示角色

---

### 1.3 主控台 Dashboard

* 全域健康摘要
  * 總節點數
  * 在線節點數
  * 離線節點數
  * 啟用中的 checks 數量
  * 目前 active alerts 數量
* 最近異常
  * 最近封包遺失事件
  * 最近 latency 超過閾值事件
  * 最近 DNS 錯誤事件
  * 最近 traceroute path change 事件
* 節點狀態總覽
  * 節點名稱
  * 狀態：online、offline、draining
  * 最後連線時間
  * 區域
  * provider
  * agent version
* Check 狀態總覽
  * check 名稱
  * check 類型：ping、traceroute、dns
  * target
  * 啟用狀態
  * 最近一次執行結果
* 圖表摘要
  * Ping latency 趨勢
  * Ping packet loss 趨勢
  * DNS response time 趨勢
  * Traceroute path change 摘要
* 空狀態
  * 尚未建立節點時，提示「新增節點」
  * 尚未建立 check 時，提示「建立第一個 check」
  * 尚未指派 check 時，提示「將 check 指派給節點」

---

### 1.4 節點 / Probe 管理

### 節點列表

* 顯示所有節點
  * 節點名稱
  * 狀態
  * provider
  * region
  * hostname
  * IP family 支援：IPv4、IPv6
  * agent version
  * 最後 heartbeat 時間
  * 啟用 / 停用狀態
  * 可切換成地圖 View 顯示所有 Probe 與 Check 的位置
* 篩選節點
  * 依狀態篩選
  * 依 provider 篩選
  * 依 region 篩選
  * 依標籤篩選
* 搜尋節點
  * 依名稱
  * 依 hostname
  * 依 provider
  * 依 region
* 排序節點
  * 依名稱
  * 依狀態
  * 依最後連線時間
  * 依建立時間

### 新增節點

* 建立新節點流程
  * 輸入節點名稱
  * 選填 provider
  * 選填 region
  * 選填 hostname
  * 選填標籤
* 產生註冊 token
  * 顯示一次性 registration token
  * 顯示安裝 / 註冊指令
  * 提醒 token 僅顯示一次或有有效期限
* 節點安裝指引
  * 顯示 Linux 安裝指令
  * 顯示設定 controller URL 的方式
  * 顯示貼上 registration token 的方式
  * 顯示啟動 service 的方式
* 註冊完成狀態
  * 顯示等待節點連線
  * 節點成功連線後顯示 online
  * 節點註冊失敗時顯示 troubleshooting 提示

### 節點詳細頁

* 基本資訊
  * 名稱
  * provider
  * region
  * hostname
  * enabled 狀態
  * 建立時間
  * 更新時間
  * 繪製地圖位置
* 連線狀態
  * online / offline / draining
  * 最後 heartbeat
  * 最近一次重新連線時間
  * agent version
  * uptime
* 能力資訊
  * 支援 IPv4
  * 支援 IPv6
  * 是否支援 raw ICMP
  * 是否支援 DNS TCP fallback
  * 作業系統資訊
* 健康資訊
  * CPU 使用率
  * 記憶體使用量
  * queue length
  * uptime
  * 最近錯誤事件
* 已指派 checks
  * check 名稱
  * check 類型
  * target
  * interval
  * jitter
  * 啟用狀態
  * 最近一次結果
* 節點操作
  * 編輯節點資訊
  * 啟用 / 停用節點
  * 重新產生安裝指令
  * 輪替 probe secret
  * 標記 draining
  * 移除節點

---

### 1.5 Check 管理

### Check 列表

* 顯示所有 checks
  * check 名稱
  * 類型：ping、traceroute、dns
  * target
  * target type
  * 啟用狀態
  * 被幾個節點指派
  * 最近一次執行狀態
* 篩選 checks
  * 依 check 類型
  * 依啟用狀態
  * 依 target type
* 搜尋 checks
  * 依名稱
  * 依 target
  * 依描述
* 排序 checks
  * 依名稱
  * 依建立時間
  * 依最近結果時間

### 新增 Check

* 選擇 check 類型
  * Ping
  * Traceroute
  * DNS
* 設定基本資訊
  * check 名稱
  * target
  * target type：host、ip、dns_query、dns_resolver
  * 描述
  * 是否啟用
* 設定共用參數
  * IP version：IPv4、IPv6、auto
  * timeout
* 儲存後可進入指派節點流程

### 編輯 Check

* 修改名稱
* 修改描述
* 修改參數
* 啟用 / 停用 check
* 修改 target
* 修改 target type
* 修改 check type
* 修改 target、target type、check type 時，前端應顯示警告：
  * 此操作可能影響歷史圖表語意
  * 舊結果會與新 target 顯示在同一個 check 下

---

### 1.6 Ping Check 設定

* Ping 設定表單
  * target host 或 IP
  * IP version：IPv4、IPv6、auto
  * packet count
  * packet interval
  * timeout
  * packet size
  * source interface，可選
  * raw ICMP / system ping fallback 顯示狀態
* Ping 結果顯示
  * sent count
  * received count
  * loss percentage
  * min latency
  * avg latency
  * median latency
  * max latency
  * standard deviation
  * RTT samples
  * resolved IP
  * error message
* Ping 圖表
  * latency time series
  * packet loss time series
  * SmokePing 風格 latency distribution，可作為後續設計方向

---

### 1.7 Traceroute Check 設定

* Traceroute 設定表單
  * target host 或 IP
  * IP version：IPv4、IPv6、auto
  * max hops
  * queries per hop
  * timeout per hop
  * protocol：第一版只顯示 ICMP
* Traceroute 結果顯示
  * resolved target IP
  * reached / not reached
  * hop count
  * path hash
  * 每一 hop 的 IP
  * 每一 hop 的 hostname
  * 每一 hop 的 RTT samples
  * 每一 hop 的 loss percentage
  * timeout / error
* Traceroute 視覺化
  * hop table
  * path change 標示
  * 最近幾次 traceroute 差異比較
  * route timeline，可作為後續功能

---

### 1.8 DNS Check 設定

* DNS 設定表單
  * query name
  * record type：A、AAAA、CNAME、MX、NS、TXT
  * resolver：system 或指定 resolver address
  * transport：UDP，必要時 fallback TCP
  * timeout
  * attempts
  * IP version：IPv4、IPv6、auto
* DNS 結果顯示
  * query name
  * record type
  * resolver address
  * transport
  * RCODE
  * success / failed
  * response time
  * answer count
  * answers
  * CNAME chain
  * error message
* DNS 圖表
  * response time time series
  * success rate
  * RCODE 分佈
  * answer 變化紀錄

---

### 1.9 Probe-Check 指派管理

* 指派 check 給節點
  * 選擇一個 check
  * 選擇一個或多個節點
  * 設定 interval seconds
  * 設定 jitter seconds
  * 設定 enabled 狀態
* 從節點詳細頁新增 assignment
  * 選擇現有 check
  * 設定排程
* 從 check 詳細頁新增 assignment
  * 選擇要執行此 check 的節點
  * 設定每個節點的 interval / jitter
* 指派列表
  * 節點名稱
  * check 名稱
  * check 類型
  * target
  * interval
  * jitter
  * enabled
  * 最近一次執行結果
* 編輯 assignment
  * 啟用 / 停用
  * 修改 interval
  * 修改 jitter
  * 修改顯示名稱
* 移除 assignment
  * soft delete
  * 不影響歷史結果
* 重複指派防呆
  * 同一個節點不可重複指派同一個 active check
  * 前端需顯示清楚錯誤訊息

---

### 1.10 Results 查詢與圖表

* 全域 results 頁
  * 依時間範圍查詢
  * 依 check type 篩選
  * 依 probe 篩選
  * 依 check 篩選
  * 依 status 篩選
* Check results 頁
  * 顯示某個 check 在所有節點上的結果
  * 可比較不同節點的量測結果
  * 可切換時間範圍
* Probe-check results 頁
  * 顯示某個節點執行某個 check 的完整歷史
  * 適合單一圖表與詳細除錯
* 時間範圍選擇
  * 最近 1 小時
  * 最近 6 小時
  * 最近 24 小時
  * 最近 7 天
  * 最近 30 天
  * 自訂時間範圍
* 結果狀態
  * success
  * partial
  * timeout
  * error
* 詳細結果 drawer / modal
  * 顯示單次 run 的完整資料
  * 顯示 raw metadata
  * 顯示錯誤碼與錯誤訊息
  * 顯示 started_at / finished_at / duration

---

### 1.11 告警與事件

* Alerts 列表
  * 顯示告警類型
  * 顯示影響的 check
  * 顯示影響的節點
  * 顯示嚴重程度
  * 顯示觸發時間
  * 顯示目前狀態
* 告警類型
  * packet loss 超過閾值
  * latency 超過閾值
  * traceroute path change
  * DNS query error
  * DNS response code 異常
  * probe offline
* 告警詳細頁
  * 觸發原因
  * 相關結果
  * 相關節點
  * 相關 check
  * 時間線
* 告警狀態
  * active
  * resolved
  * acknowledged，後續可做
* MVP 可先做事件列表，不一定要做完整通知通道

---

### 1.12 系統與使用體驗

* Loading state
  * 頁面 loading
  * 表單送出中
  * 圖表資料載入中
* Empty state
  * 無節點
  * 無 checks
  * 無 assignments
  * 無 results
  * 無 alerts
* Error state
  * 登入失敗
  * 權限不足
  * 資料載入失敗
  * 表單驗證失敗
  * 節點註冊失敗
* 表單驗證
  * 必填欄位
  * 數字範圍
  * interval 合理範圍
  * timeout 合理範圍
  * target 格式提示
* 危險操作確認
  * 移除節點
  * 移除 check
  * 移除 assignment
  * 修改 check type / target
  * 停用節點
* 權限顯示
  * 依角色顯示功能入口
  * MVP 可先不限制細節，但設計上保留角色概念

---

## 2. 後端 / Controller 功能

### 2.1 使用者與驗證

* 使用者註冊
  * 建立使用者
  * 確認 email 唯一
  * 儲存 password hash
* 使用者登入
  * 驗證 email / password
  * 簽發 JWT access token
  * token 有效時間 12 小時
* 使用者身份解析
  * 根據 JWT 識別使用者
  * 回傳目前使用者資訊
* 帳號狀態
  * 支援 active / disabled
  * 停用帳號不可登入

---

### 2.2 Team / 組織管理

* 建立 team
  * 建立 team 基本資訊
  * 建立者自動成為 owner
* 查詢使用者所屬 teams
* 查詢 team 詳細資訊
* 更新 team
  * 名稱
  * slug
* Team membership
  * 新增成員
  * 查詢成員
  * 修改角色
  * soft delete 成員關係
* Team scope 隔離
  * probes 屬於 team
  * checks 屬於 team
  * probe-check assignments 屬於 team
  * results 屬於 team
  * alerts 屬於 team

---

### 2.3 Probe / 節點管理

* 建立 probe 記錄
  * 名稱
  * provider
  * region
  * hostname
  * enabled 狀態
* 產生一次性 registration token
* Probe 註冊
  * 驗證 registration token
  * 建立 probe ID
  * 回傳長期 probe secret
  * 建立 probe credential
* Probe 驗證
  * 驗證 signed bearer token 或 HMAC signature
  * 確認 probe 是否啟用
  * 確認 probe 是否屬於正確 team
* Probe metadata 儲存
  * provider
  * region
  * hostname
  * labels
  * IP family 支援
  * agent version
  * capabilities
* Probe 狀態管理
  * online
  * offline
  * draining
  * last_seen_at
  * health
* Probe secret 輪替
* Probe 停用
  * 停用後不再下發檢查設定
* Probe soft delete
  * 移除節點但保留歷史結果 reference

---

### 2.4 Check 管理

* 建立 check
  * check name
  * check type
  * target
  * target type
  * description
  * parameters
  * enabled
* 查詢 check 列表
* 查詢 check 詳細資料
* 更新 check
  * 名稱
  * 描述
  * target
  * target type
  * check type
  * parameters
  * enabled
* 停用 check
  * 停用後不應再被 probe 執行
* soft delete check
  * 保留歷史結果 reference
* Check 參數驗證
  * ping 參數驗證
  * traceroute 參數驗證
  * DNS 參數驗證
  * timeout 合理範圍
  * interval 合理範圍
  * record type 合法性
  * IP version 合法性

---

### 2.5 Probe-check Assignment 管理

* 建立 assignment
  * 指定 probe
  * 指定 check
  * 設定 interval
  * 設定 jitter
  * 設定 enabled
* 查詢 assignments
  * 依 probe 查詢
  * 依 check 查詢
  * 依 team 查詢
* 更新 assignment
  * enabled
  * interval
  * jitter
  * 顯示名稱
* soft delete assignment
* 防止重複 active assignment
  * 同一 team 下，同一 probe + check 不可重複啟用
* Assignment 變更後觸發 probe 設定更新
  * check 被修改時通知相關 probes
  * assignment 被修改時通知對應 probe
  * probe 被停用時通知該 probe 停止執行 checks

---

### 2.6 Probe 設定產生與同步

* 產生 probe 專屬設定快照
  * 只包含該 probe 被指派的 checks
  * 只包含 enabled probe-checks
  * 只包含 enabled checks
  * 只包含目前有效設定
* 設定快照內容
  * probe_check ID
  * check ID
  * check type
  * target
  * target type
  * parameters
  * interval
  * jitter
  * enabled 狀態
* 產生 snapshot hash
  * 讓 probe 判斷是否需要重新同步
* 下發完整設定快照
  * probe 初次連線
  * probe 重新連線
  * probe 主動要求 resync
* 下發設定 delta
  * check 修改
  * assignment 修改
  * check 停用
  * assignment 停用
  * probe 停用
* 接收設定 ack
  * probe 回報已套用設定
  * controller 記錄目前 probe 套用版本
* 要求 probe resync
  * 當設定狀態不一致時要求完整重新同步

---

### 2.7 Probe gRPC Stream 管理

* 維持 probe 長連線
  * 雙向 streaming
  * controller 傳送設定
  * probe 傳送結果與健康狀態
* 接收 heartbeat
  * 更新 last_seen_at
  * 更新 online 狀態
* 接收 health metrics
  * CPU
  * memory
  * uptime
  * queue length
  * capabilities
  * version
* 接收 log event
  * probe 端錯誤
  * probe 端警告
  * 設定套用失敗
* 傳送 controller message
  * config snapshot
  * config delta
  * res 端警告
  * 設定套用失敗
* 傳送 controller message
  * config snapshot
  * config delta
    ync request
  * shutdown / drain 指令
* Stream 重連處理
  * probe 重新連線後重新驗證
  * 必要時傳送完整設定
  * 避免重複套用錯誤設定

---

### 2.8 結果接收與儲存

* 接收 ping 結果
  * 驗證 probe_check_id
  * 驗證 team scope
  * 儲存摘要
  * 儲存 RTT samples
  * 儲存錯誤資訊
* 接收 traceroute 結果
  * 儲存 run 摘要
  * 儲存每一 hop 資訊
  * 儲存 path hash
  * 儲存 reached 狀態
  * 儲存錯誤資訊
* 接收 DNS 結果
  * 儲存 query 資訊
  * 儲存 resolver
  * 儲存 RCODE
  * 儲存 answers
  * 儲存 response time
  * 儲存錯誤資訊
* 結果驗證
  * 確認結果來自正確 probe
  * 確認 probe_check 仍屬於同一 team
  * 確認時間欄位合理
  * 確認 payload 與 check type 一致
* 結果查詢
  * 查詢全域結果
  * 查詢某 check 的結果
  * 查詢某 probe-check assignment 的結果
  * 支援時間範圍
  * 支援 check type
  * 支援 status

---

### 2.9 彙總與圖表資料

* Ping 彙總
  * 平均 latency
  * median latency
  * max latency
  * packet loss
  * success rate
* DNS 彙總
  * response time
  * success rate
  * RCODE 統計
  * answer count
* Traceroute 彙總
  * hop count
  * path hash 變化
  * reached rate
* 時間區間彙總
  * 1 分鐘 bucket
  * 5 分鐘 bucket
  * 1 小時 bucket
  * 後續可依資料量調整
* 圖表資料輸出
  * latency series
  * loss series
  * DNS response time series
  * traceroute path change series

---

### 2.10 告警事件

* Ping 告警條件
  * packet loss 超過閾值
  * latency 超過閾值
  * timeout 過多
* DNS 告警條件
  * query failed
  * RCODE 異常
  * response time 超過閾值
  * answer count 異常
  * answer 變化
* Traceroute 告警條件
  * path hash 改變
  * hop count 大幅變化
  * 無法抵達終點
* Probe 告警條件
  * probe offline
  * heartbeat 過期
  * queue length 過高
  * agent version 過舊，後續可做
* 告警事件管理
  * 建立 alert event
  * 標記 resolved
  * 查詢 active alerts
  * 查詢歷史 alerts
* MVP 可以先做事件記錄與列表，不一定要做通知整合

---

### 2.11 安全與防濫用

* 使用者密碼安全
  * 只儲存 hash
  * 不記錄明文密碼
* JWT 驗證
  * 驗證簽章
  * 驗證過期時間
  * 驗證使用者狀態
* Team scope 驗證
  * 使用者只能讀寫自己 team 的資源
* Probe 驗證
  * 每個 probe request 都需要驗證
  * probe secret 不回傳給前端
  * probe credential 不進入設定快照
* Probe 權限限制
  * probe 只能上傳自己被指派的結果
  * probe 不能讀取其他 probe 設定
* 設定下發限制
  * controller 只下發該 probe 被指派的 checks
* 輸入驗證
  * target 格式
  * resolver 格式
  * timeout
  * interval
  * packet size
  * DNS record type

---

## 3. Probe / 節點 Agent 功能

### 3.1 安裝與初始化

* 支援安裝在
  * VPS
  * bare metal server
  * internal host
* 提供設定檔
  * controller URL
  * registration token
  * probe name，可選
  * local limits，可選
* 支援以 service 方式執行
  * systemd
  * container，後續可做
* 啟動時讀取設定
* 啟動時檢查本機能力
  * raw ICMP 是否可用
  * system ping 是否可用
  * traceroute 是否可用
  * IPv4 是否可用
  * IPv6 是否可用
  * DNS resolver 是否可用

---

### 3.2 Probe 註冊

* 使用一次性 registration token 向 controller 註冊
* 註冊成功後取得
  * probe ID
  * long-term probe secret
* 本機安全儲存 probe secret
* 註冊後回報基本 metadata
  * hostname
  * agent version
  * OS
  * architecture
  * IP family support
  * capabilities
* 註冊失敗時顯示明確錯誤
  * token 無效
  * token 過期
  * controller 不可達
  * TLS / 連線錯誤

---

### 3.3 與 Controller 連線

* 建立安全 gRPC stream
* 每次連線都進行 probe 驗證
* 維持長連線
* 自動重連
  * controller 暫時不可達時 retry
  * 使用 backoff 避免過度重試
* 重新連線後重新同步設定
* 支援 controller 要求 resync
* 支援 controller 要求 shutdown / drain
* stream 中斷時繼續使用最後有效設定執行 checks

---

### 3.4 設定接收與套用

* 接收完整設定快照
* 接收設定 delta
* 驗證設定格式
* 驗證 check type 是否支援
* 驗證參數是否合法
* 套用 probe-check assignments
* 停用已移除或 disabled 的 assignments
* 回報 config ack
* 保存最後成功套用的設定
* 使用 snapshot hash 判斷是否需要重新同步
* 設定套用失敗時回報錯誤事件

---

### 3.5 本機排程器

* 根據 probe_check 設定排程
  * interval_seconds
  * jitter_seconds
  * enabled
* 每個 assignment 獨立排程
* 套用 jitter 避免多個 probe 同步量測
* 支援動態新增 schedule
* 支援動態移除 schedule
* 支援動態更新 interval / jitter
* 避免同一個 assignment 重疊執行
* 支援 timeout
* 支援 queue length 回報
* controller 不可達時仍使用最後設定繼續執行
* controller 恢復後補送或繼續送新結果，MVP 可先不保證完整補送

---

### 3.6 Ping 執行

* 執行 ping check
* 支援 target host
* 支援 target IP
* 支援 IPv4
* 支援 IPv6
* 支援 auto IP version
* 支援 packet count
* 支援 packet interval
* 支援 timeout
* 支援 packet size
* 支援 source interface，可選
* 優先使用 Go ICMP 實作
* raw socket 不可用時 fallback 到 system ping
* 收集結果
  * sent count
  * received count
  * loss percent
  * min RTT
  * avg RTT
  * median RTT
  * max RTT
  * stddev RTT
  * RTT samples
  * resolved IP
  * error code
  * error message
* 將結果格式化後送回 controller

---

### 3.7 Traceroute 執行

* 執行 traceroute check
* 支援 target host
* 支援 target IP
* 支援 IPv4
* 支援 IPv6
* 支援 auto IP version
* 支援 max hops
* 支援 queries per hop
* 支援 timeout per hop
* 第一版支援 ICMP protocol
* 收集結果
  * resolved target IP
  * hop number
  * hop IP
  * hop hostname
  * 每 hop RTT samples
  * 每 hop loss percent
  * 每 hop timeout / error
  * reached flag
  * hop count
  * path hash
* 如果原生 Go traceroute 來不及，可先使用 system traceroute parser
* 將結果格式化後送回 controller

---

### 3.8 DNS 執行

* 執行 DNS check
* 支援 query name
* 支援 record type
  * A
  * AAAA
  * CNAME
  * MX
  * NS
  * TXT
* 支援 system resolver
* 支援指定 resolver address
* 第一版使用 UDP
* UDP 需要時 fallback TCP
* 支援 timeout
* 支援 attempts
* 支援 IPv4 / IPv6 / auto
* 收集結果
  * query name
  * record type
  * resolver address
  * transport
  * RCODE
  * success
  * response time
  * answer count
  * answers
  * CNAME chain
  * error code
  * error message
* 使用結構化 DNS library 取得 answer、rcode、ttl 與 metadata
* 將結果格式化後送回 controller

---

### 3.9 Probe 本機安全限制

* 避免 probe 被濫用為掃描器
* 限制最大執行頻率
* 限制同時執行數
* 限制 packet count
* 限制 packet size
* 限制 traceroute max hops
* 限制 DNS attempts
* 限制 timeout 最大值
* 可禁止特定 target range，後續可做
* 可限制 private IP / loopback target，依部署需求決定
* 設定不合法時拒絕執行並回報錯誤
* controller 下發 disabled check 時停止執行

---

### 3.10 結果上傳

* 將每次 check run 結果送回 controller
* 結果包含
  * team ID
  * probe_check ID
  * started_at
  * finished_at
  * status
  * measurement payload
  * error code
  * error message
  * metadata
* 支援結果 streaming
* controller 暫時不可達時
  * MVP 可先丟棄或有限度緩衝
  * 後續支援本機 buffer 與補送
* 上傳失敗時記錄本機 log
* 上傳成功後可清除本機 queue

---

### 3.11 健康回報

* 定期傳送 heartbeat
* 回報 agent version
* 回報 uptime
* 回報 CPU 使用率
* 回報 memory 使用量
* 回報 queue length
* 回報目前 active schedules 數量
* 回報支援能力
  * raw ICMP
  * system ping fallback
  * traceroute support
  * DNS UDP
  * DNS TCP fallback
  * IPv4
  * IPv6
* 回報本機錯誤事件
  * config apply failed
  * check execution failed
  * controller connection failed
  * permission issue
  * raw socket unavailable

---
