# Netstamp 計畫

Netstamp 是一套以 Go 建構的分散式網路量測系統，靈感來自 SmokePing、Globalping，以及 RIPE Atlas。第一個里程碑會聚焦在後端：中央控制器與分散式 probe。Probe 會量測已設定的 target、儲存時間序列結果，並提供 API，供後續 UI、告警，以及公開或共享量測功能使用。

## 控制器

控制器是中央後端服務。它負責持有設定、將設定變更串流給 probe、接收 probe 結果、儲存歷史資料，並暴露 API。

核心職責：

- 註冊並驗證 probe。
- 儲存 probe 中繼資料，例如區域、供應商、主機名稱、標籤、IP family 支援，以及軟體版本。
- 定義 target，例如網域、IP 位址、DNS 查詢目標，或 DNS resolver。
- 定義哪些 probe 應該量測哪些 target。
- 儲存每個連線的量測設定。第一版只支援 ping、traceroute，以及 DNS。
- 只將每個 probe 被指派的 probe-to-target 檢查設定與 target 資料傳給該 probe。
- 當某個 probe 被指派的 target 或檢查變更時，將更新後的設定串流給該 probe。
- 接收量測結果。
- 儲存原始與彙總後的量測資料。
- 暴露 REST API，供設定與查詢結果使用。
- 當封包遺失、延遲、路徑變化，或 DNS 查詢錯誤跨過閾值時，發出告警事件。

### Probe

Probe 是安裝在 VPS、裸機伺服器，或內部主機上的輕量 Go agent。它會連到控制器、接收被指派的 probe-to-target 檢查設定、在本機排程這些檢查、執行檢查，並把結果送回。

核心職責：

- 使用 token 向控制器註冊。
- 維持安全的 gRPC 串流，並具備重連與重新同步行為。
- 套用 probe-to-target 設定快照與增量更新。
- 根據收到的設定在本機排程檢查。
- 執行 ping 檢查。
- 執行 traceroute 檢查。
- 執行 DNS 檢查。
- 強制執行本機安全限制，避免 probe 被濫用為掃描器。
- 回傳帶有時間、錯誤與中繼資料的結構化結果。
- 自我回報健康度量，例如 CPU、記憶體、版本、uptime、queue length，以及支援能力。

## 架構

```text
+-----------------+   gRPC probe stream      +-----------------+
|                 | <----------------------> |                 |
|   Controller    |                          |      Probe      |
|                 |                          |                 |
+--------+--------+                          +--------+--------+
         |                                            |
         |                                            |
         v                                            v
+-----------------+                          +----------------------+
|   PostgreSQL    |                          | ping/traceroute/dns |
|   Timeseries    |                          | OS networking        |
+-----------------+                          +----------------------+
         |
         v
+-----------------+
| Redis / Queue   |
| optional later  |
+-----------------+
```

## 控制器 API

管理與設定 API 使用 HTTP REST，因為它容易除錯，也方便未來前端使用。控制器與 probe 之間的流量使用 gRPC，因為這個系統需要一條長時間存在的雙向 probe 通道：

- 控制器到 probe：設定快照與增量更新。
- Probe 到控制器：觀測到的量測資料、設定確認、heartbeat，以及健康狀態。

核心需求是控制器會把 probe 被指派的 `probe_target` + `probe_target_check` + `target` 資料送給該 probe，並在資料變更時串流新的版本。Probe 在本機排程檢查，並透過同一條 gRPC 連線把觀測資料串流回來。

建議的第一版：

- REST 用於管理與面向 UI 的 API。
- gRPC 雙向串流用於 probe 控制與觀測資料。
- 若有幫助，使用 gRPC unary methods 處理註冊、bootstrap 設定抓取，以及明確的重新同步。
- 從第一天開始就對 protobuf 訊息做版本化。
- REST 結果匯入可保留為手動測試或 fallback，不作為主要 probe 路徑。

對這種 probe 設定模型來說，gRPC 比單純 polling 更適合，因為設定變更應該快速抵達 probe，而 probe 也需要頻繁把觀測資料送回。管理 API 保持 REST 仍然合理。務實的分工是：REST 給人與前端使用，gRPC streaming 給 probe 流量使用。

### 使用者與驗證

第一版使用簡單的 email + password 登入，不做角色、租戶，或細緻權限模型。所有登入使用者都可以操作同一組設定。

- `users.email` 必須唯一。
- 密碼只儲存 password hash，不儲存明文。
- 登入成功後簽發 JWT access token。
- JWT 有效時間固定 12 小時。
- JWT claim 只需要 `sub`、`email`、`iat`、`exp`。
- 第一版不做 refresh token；token 過期後重新登入。
- 第一版不做 permission table、role table，或 per-resource ACL。

### REST API

初始 endpoints：

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/auth/me`
- `POST /api/v1/probes/register`
- `GET /api/v1/probes`
- `GET /api/v1/probes/{id}`
- `PATCH /api/v1/probes/{id}`
- `POST /api/v1/targets`
- `GET /api/v1/targets`
- `GET /api/v1/targets/{id}`
- `PATCH /api/v1/targets/{id}`
- `POST /api/v1/probe-targets`
- `GET /api/v1/probe-targets`
- `GET /api/v1/probe-targets/{id}`
- `PATCH /api/v1/probe-targets/{id}`
- `POST /api/v1/probe-targets/{id}/checks`
- `GET /api/v1/probe-targets/{id}/checks`
- `PATCH /api/v1/probe-target-checks/{check_id}`
- `GET /api/v1/config-versions`
- `GET /api/v1/config-versions/{config_version}`
- `GET /api/v1/results`
- `GET /api/v1/probe-targets/{id}/results`
- `GET /api/v1/alerts`

### Probe gRPC API

初始面向 probe 的 RPC：

- `RegisterProbe(RegisterProbeRequest) returns (RegisterProbeResponse)`
- `GetProbeConfig(GetProbeConfigRequest) returns (ProbeConfigSnapshot)`
- `Observe(stream ProbeObserveMessage) returns (stream ControllerMessage)`

`Observe` stream 應該承載：

- 控制器訊息：設定快照、設定 delta、重新同步請求、shutdown/drain 指令。
- Probe 訊息：heartbeat、設定 ack、ping 結果、traceroute 結果、DNS 結果、健康度量、log event。

Probe 驗證：

- Probe 會收到一次性註冊 token。
- 控制器回傳 probe ID 與長期 probe secret。
- 每個 probe 請求都使用 HMAC signature 或 mTLS。
- 第一版可以使用 signed bearer tokens；等到營運許多不受信任的 probe 時，再加入 mTLS。

## 量測類型

第一版只支援 `ping`、`traceroute`，以及 `dns`。HTTP、TCP connect、TLS 檢查都先不列入第一版資料模型與 API。

### Ping

Ping 應收集 SmokePing 風格的延遲與封包遺失資料。

必要欄位：

- Target host 或 address。
- IP version：IPv4、IPv6，或 auto。
- 封包數量。
- 封包間隔。
- Timeout。
- 封包大小。
- 可選的來源 interface。
- 特權 raw ICMP 模式，或非特權 system command fallback。

結果欄位：

- 已送出封包數。
- 已收到封包數。
- 封包遺失百分比。
- 最小延遲。
- 平均延遲。
- 中位數延遲。
- 最大延遲。
- 標準差。
- 每個封包的延遲樣本。
- 解析後的 IP 位址。
- 失敗時的錯誤訊息。

實作選擇：

- 優先使用以 `golang.org/x/net/icmp` 建構的 Go ICMP 實作。
- 只有在 raw socket 權限不可用的地方，才支援 fallback 到系統 `ping`。
- 同時儲存摘要資料與樣本，讓之後可以產生類似 SmokePing 的圖表。

### Traceroute

Traceroute 應該是輕量診斷功能，類似 Globalping 風格的檢查，但整合進類似 SmokePing 的歷史平台。

必要欄位：

- Target host 或 address。
- IP version：IPv4、IPv6，或 auto。
- Max hops。
- 每一 hop 的 queries 數量。
- 每一 hop 的 timeout。
- Protocol：第一版只支援 ICMP。

結果欄位：

- 解析後的 target IP。
- Hop number。
- Hop IP。
- 可選的 hop hostname。
- 每一 hop 的 RTT 樣本。
- 每一 hop 的封包遺失。
- 每一 hop 的錯誤或 timeout。
- 是否抵達終點的 flag。
- 用於偵測路由變化的 path hash。

實作選擇：

- 只有在 raw implementation 建置太慢時，才先使用系統 `traceroute`/`tracert` parser。
- 長期偏好的做法是原生 Go traceroute，以取得可攜性與結構化結果。
- 將完整 hop data 儲存為 JSONB，並另外儲存 path hash，以快速偵測路徑變更。

### DNS

DNS 檢查用來追蹤解析是否成功、解析延遲、回應碼，以及答案是否變化。第一版先支援一般查詢，不處理 DNSSEC 驗證或 DoH/DoT。

必要欄位：

- Query name；預設可使用 `targets.address`。
- Record type：至少支援 `A`、`AAAA`、`CNAME`、`MX`、`NS`、`TXT`。
- Resolver：`system` 或指定 resolver address。
- Transport：第一版使用 UDP，必要時 fallback TCP。
- Timeout。
- Attempts。
- IP version：IPv4、IPv6，或 auto。

結果欄位：

- Query name。
- Record type。
- Resolver address。
- Transport。
- RCODE。
- 是否成功。
- 查詢耗時。
- Answer count。
- Answers，包含 name、type、ttl、data。
- CNAME chain。
- 失敗時的錯誤訊息。

實作選擇：

- 優先使用 `github.com/miekg/dns`，因為需要拿到結構化的 answer、rcode、ttl，以及 raw DNS metadata。
- `system` resolver 代表使用 probe 作業系統的 resolver 設定；指定 resolver 則由 probe 直接查詢該 resolver。
- 完整 DNS response 可放在 `results.raw`，穩定查詢欄位放在 `results.summary` 與 `results.samples`。

## 資料模型

初始資料表應該以 `probe -> target` 來建模產品，而不是以人工管理的 job 來建模。持久化的使用者設定是：

```text
probe + target + enabled checks + schedule + parameters
```

控制器儲存這份設定並送給 probe。Probe 使用收到的 schedule 與 parameters 在本機執行檢查。控制器不需要為每個檢查間隔建立或派送個別 job。

命名注意事項：本計畫使用 `probe` 表示執行量測的 agent，使用 `target` 表示被量測的目標。`anchor` 保留給未來「穩定、公開、可被其他 probe 量測的基準 probe」概念，不作為第一版主表命名。

建議 PostgreSQL extensions：

- `pgcrypto`：提供 `gen_random_uuid()`。
- `citext`：讓 email 唯一性不受大小寫影響。

### `config_version` 設計原則

`config_version` 是全域單調遞增的設定版本，用來追蹤所有使用者設定變化。每次設定變更都先建立一筆 `config_versions`，再把受影響的設定資料以新 `config_version` 插入版本化資料表。

- 版本化設定表不使用 `updated_at`。更新是一筆新的 version row。
- 每個版本都保留，不覆寫舊資料。
- 刪除不是 hard delete，而是在新版本插入 `deleted = true` 的 tombstone row。
- 查詢某個版本的有效設定時，對每個 logical `id` 取 `config_version <= requested_config_version` 的最新 row，並排除 `deleted = true`。
- 版本化表之間的 `probe_id`、`target_id`、`probe_target_id` 是 logical reference；service layer 需要驗證它們在同一個 effective config version 中存在。
- `config_version` 追蹤的是 desired configuration，不用來記錄每次 result、heartbeat，或 runtime health 變化。
- `config_versions.created_at` 是版本建立時間；版本化設定表本身不需要 `created_at` / `updated_at`。

### `users`

第一版只支援 email + password 登入。此表不處理權限，所有使用者平權。

| 欄位 | PostgreSQL type | 限制 | 說明 |
| --- | --- | --- | --- |
| `id` | `UUID` | `PRIMARY KEY DEFAULT gen_random_uuid()` | 使用者 ID |
| `email` | `CITEXT` | `NOT NULL UNIQUE` | 登入 email；需要 `citext` extension |
| `password_hash` | `TEXT` | `NOT NULL` | Argon2id 或 bcrypt 的 encoded hash |
| `is_active` | `BOOLEAN` | `NOT NULL DEFAULT true` | 停用帳號用，不代表權限 |
| `created_at` | `TIMESTAMPTZ` | `NOT NULL DEFAULT now()` | 建立時間 |

### `config_versions`

| 欄位 | PostgreSQL type | 限制 | 說明 |
| --- | --- | --- | --- |
| `config_version` | `BIGINT` | `PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY` | 全域設定版本 |
| `created_by_user_id` | `UUID` | `NULL REFERENCES users(id)` | 由哪個使用者造成；系統變更可為 NULL |
| `change_type` | `TEXT` | `NOT NULL` | 例如 `create_probe`、`update_check`、`delete_target` |
| `change_summary` | `JSONB` | `NOT NULL DEFAULT '{}'::jsonb` | 給 audit/API 顯示用的摘要 |
| `checksum` | `TEXT` | `NULL` | 可選的設定快照 checksum |
| `created_at` | `TIMESTAMPTZ` | `NOT NULL DEFAULT now()` | 版本建立時間 |

### `probes`

Probe 設定表是 append-only versioned table。`id` 是 logical probe ID，`config_version` 是該 logical probe 在某次設定版本的內容。

| 欄位 | PostgreSQL type | 限制 | 說明 |
| --- | --- | --- | --- |
| `id` | `UUID` | `NOT NULL` | Logical probe ID |
| `config_version` | `BIGINT` | `NOT NULL REFERENCES config_versions(config_version)` | 此 row 所屬設定版本 |
| `name` | `TEXT` | `NOT NULL` | 顯示名稱 |
| `provider` | `TEXT` | `NULL` | 例如 `hetzner`、`aws`、`home` |
| `region` | `TEXT` | `NULL` | Probe 區域 |
| `hostname` | `TEXT` | `NULL` | Probe 主機名稱 |
| `labels` | `JSONB` | `NOT NULL DEFAULT '{}'::jsonb` | 自訂標籤 |
| `enabled` | `BOOLEAN` | `NOT NULL DEFAULT true` | 是否允許 probe 執行檢查 |
| `deleted` | `BOOLEAN` | `NOT NULL DEFAULT false` | tombstone |

主鍵：`PRIMARY KEY (id, config_version)`。

### `probe_credentials`

Probe 憑證不放進 probe 設定快照，也不會被串流給 probe。第一版只需要保存目前有效的 probe secret hash。

| 欄位 | PostgreSQL type | 限制 | 說明 |
| --- | --- | --- | --- |
| `probe_id` | `UUID` | `PRIMARY KEY` | Logical probe ID |
| `secret_hash` | `TEXT` | `NOT NULL` | 長期 probe secret 的 hash |
| `created_at` | `TIMESTAMPTZ` | `NOT NULL DEFAULT now()` | 憑證建立時間 |
| `last_rotated_at` | `TIMESTAMPTZ` | `NULL` | 最近一次輪替時間 |

### `probe_status`

Probe runtime 狀態不是設定歷史，使用目前狀態表即可。後續若要長期保存 health history，再另外建立 time-series 表。

| 欄位 | PostgreSQL type | 限制 | 說明 |
| --- | --- | --- | --- |
| `probe_id` | `UUID` | `PRIMARY KEY` | Logical probe ID |
| `status` | `TEXT` | `NOT NULL` | 例如 `online`、`offline`、`draining` |
| `last_seen_at` | `TIMESTAMPTZ` | `NULL` | 最後 heartbeat 時間 |
| `agent_version` | `TEXT` | `NULL` | Probe 軟體版本 |
| `ip_families` | `TEXT[]` | `NOT NULL DEFAULT '{}'::text[]` | 例如 `{ipv4,ipv6}` |
| `capabilities` | `JSONB` | `NOT NULL DEFAULT '{}'::jsonb` | Probe 能力，例如 raw ICMP、DNS TCP fallback |
| `health` | `JSONB` | `NOT NULL DEFAULT '{}'::jsonb` | CPU、記憶體、queue length 等 |

### `targets`

`targets` 表示可被 probe 量測的目標。第一版的 `target_type` 可先限制為 `host`、`ip`、`dns_query`、`dns_resolver`。

| 欄位 | PostgreSQL type | 限制 | 說明 |
| --- | --- | --- | --- |
| `id` | `UUID` | `NOT NULL` | Logical target ID |
| `config_version` | `BIGINT` | `NOT NULL REFERENCES config_versions(config_version)` | 此 row 所屬設定版本 |
| `name` | `TEXT` | `NOT NULL` | 顯示名稱 |
| `address` | `TEXT` | `NOT NULL` | host、IP、query name，或 resolver address |
| `target_type` | `TEXT` | `NOT NULL` | `host`、`ip`、`dns_query`、`dns_resolver` |
| `description` | `TEXT` | `NULL` | 備註 |
| `deleted` | `BOOLEAN` | `NOT NULL DEFAULT false` | tombstone |

主鍵：`PRIMARY KEY (id, config_version)`。

範例：

- `1.1.1.1` 作為 `ip`。
- `example.com` 作為 `host` 或 `dns_query`。
- `8.8.8.8:53` 作為 `dns_resolver`。

### `probe_targets`

這是連線設定表，用來表示哪個 probe 量測哪個 target。

| 欄位 | PostgreSQL type | 限制 | 說明 |
| --- | --- | --- | --- |
| `id` | `UUID` | `NOT NULL` | Logical probe-target ID |
| `config_version` | `BIGINT` | `NOT NULL REFERENCES config_versions(config_version)` | 此 row 所屬設定版本 |
| `probe_id` | `UUID` | `NOT NULL` | Logical probe ID |
| `target_id` | `UUID` | `NOT NULL` | Logical target ID |
| `name` | `TEXT` | `NULL` | 覆寫顯示名稱 |
| `enabled` | `BOOLEAN` | `NOT NULL DEFAULT true` | 是否啟用此 probe-to-target |
| `interval_seconds` | `INTEGER` | `NOT NULL` | 預設排程間隔 |
| `jitter_seconds` | `INTEGER` | `NOT NULL DEFAULT 0` | Probe 本機 jitter 上限 |
| `deleted` | `BOOLEAN` | `NOT NULL DEFAULT false` | tombstone |

主鍵：`PRIMARY KEY (id, config_version)`。

範例資料列：

| id | probe_id | target_id | enabled | interval_seconds |
| --- | --- | --- | --- | --- |
| `pt_aaa` | `probe_1111` | `target_9999` | true | 60 |
| `pt_bbb` | `probe_1111` | `target_8888` | true | 300 |

### `probe_target_checks`

這張表儲存每個 probe-to-target 連線應該執行什麼。一個 `probe_target` 可以對同一個 target 執行多個檢查。第一版 `check_type` 只允許 `ping`、`traceroute`、`dns`。

| 欄位 | PostgreSQL type | 限制 | 說明 |
| --- | --- | --- | --- |
| `id` | `UUID` | `NOT NULL` | Logical check ID |
| `config_version` | `BIGINT` | `NOT NULL REFERENCES config_versions(config_version)` | 此 row 所屬設定版本 |
| `probe_target_id` | `UUID` | `NOT NULL` | Logical probe-target ID |
| `check_type` | `TEXT` | `NOT NULL` | `ping`、`traceroute`、`dns` |
| `enabled` | `BOOLEAN` | `NOT NULL DEFAULT true` | 是否啟用此檢查 |
| `parameters` | `JSONB` | `NOT NULL DEFAULT '{}'::jsonb` | 檢查參數 |
| `schedule_override_seconds` | `INTEGER` | `NULL` | 若存在，覆寫 `probe_targets.interval_seconds` |
| `deleted` | `BOOLEAN` | `NOT NULL DEFAULT false` | tombstone |

主鍵：`PRIMARY KEY (id, config_version)`。

範例資料列：

| probe_target_id | check_type | parameters |
| --- | --- | --- |
| `pt_aaa` | `ping` | `{"count":20,"timeout_ms":3000,"ip_version":"auto","packet_size":56}` |
| `pt_aaa` | `traceroute` | `{"max_hops":30,"queries_per_hop":3,"timeout_ms":3000,"protocol":"icmp"}` |
| `pt_bbb` | `dns` | `{"query_name":"example.com","record_type":"A","resolver":"system","timeout_ms":2000,"attempts":2}` |

### `probe_config_acks`

Probe 套用設定後回報 ack。這是 immutable event table，不需要 `updated_at`。

| 欄位 | PostgreSQL type | 限制 | 說明 |
| --- | --- | --- | --- |
| `id` | `UUID` | `PRIMARY KEY DEFAULT gen_random_uuid()` | Ack event ID |
| `probe_id` | `UUID` | `NOT NULL` | Logical probe ID |
| `config_version` | `BIGINT` | `NOT NULL REFERENCES config_versions(config_version)` | Probe 已處理的設定版本 |
| `status` | `TEXT` | `NOT NULL` | `applied`、`rejected` |
| `error_message` | `TEXT` | `NULL` | 套用失敗原因 |
| `acked_at` | `TIMESTAMPTZ` | `NOT NULL DEFAULT now()` | Ack 時間 |

### `results`

`results` 是 immutable time-series event table。每筆 row 表示某個 probe 根據某個 `config_version` 執行一次檢查的結果。因為 result 需要回溯到當時的設定，所以必須保存 `config_version`。

| 欄位 | PostgreSQL type | 限制 | 說明 |
| --- | --- | --- | --- |
| `id` | `UUID` | `PRIMARY KEY DEFAULT gen_random_uuid()` | Result ID |
| `probe_target_id` | `UUID` | `NOT NULL` | 當時執行的 logical probe-target ID |
| `check_id` | `UUID` | `NOT NULL` | 當時執行的 logical check ID |
| `probe_id` | `UUID` | `NOT NULL` | Logical probe ID |
| `target_id` | `UUID` | `NOT NULL` | Logical target ID |
| `check_type` | `TEXT` | `NOT NULL` | `ping`、`traceroute`、`dns` |
| `config_version` | `BIGINT` | `NOT NULL REFERENCES config_versions(config_version)` | 此結果對應的設定版本 |
| `started_at` | `TIMESTAMPTZ` | `NOT NULL` | Probe 開始執行時間 |
| `finished_at` | `TIMESTAMPTZ` | `NOT NULL` | Probe 完成時間 |
| `duration_ms` | `INTEGER` | `NOT NULL` | `finished_at - started_at` 的毫秒數 |
| `status` | `TEXT` | `NOT NULL` | `success`、`partial`、`timeout`、`error` |
| `summary` | `JSONB` | `NOT NULL DEFAULT '{}'::jsonb` | 穩定查詢與圖表欄位 |
| `samples` | `JSONB` | `NOT NULL DEFAULT '[]'::jsonb` | ping samples、traceroute hops，或 DNS answers |
| `raw` | `JSONB` | `NOT NULL DEFAULT '{}'::jsonb` | 原始輸出或完整 protocol metadata |
| `error_code` | `TEXT` | `NULL` | 可機器判讀的錯誤碼 |
| `error_message` | `TEXT` | `NULL` | 人類可讀錯誤 |
| `created_at` | `TIMESTAMPTZ` | `NOT NULL DEFAULT now()` | 控制器寫入時間 |

建議索引：

- `results (probe_id, started_at DESC)`
- `results (probe_target_id, started_at DESC)`
- `results (check_id, started_at DESC)`
- `results (check_type, started_at DESC)`
- `results (config_version)`

### Result JSON 欄位約定

`summary` 放固定且常用的查詢欄位，適合列表、圖表與告警。`samples` 放同一次檢查內的多筆觀測資料。`raw` 放完整原始資料或尚未穩定成 schema 的資料。

Ping result：

```json
{
  "summary": {
    "sent": 20,
    "received": 19,
    "loss_percent": 5,
    "rtt_min_ms": 12.4,
    "rtt_avg_ms": 16.8,
    "rtt_median_ms": 15.9,
    "rtt_max_ms": 31.2,
    "rtt_stddev_ms": 4.1,
    "resolved_ip": "1.1.1.1",
    "ip_version": "ipv4"
  },
  "samples": [
    {"seq": 1, "received": true, "rtt_ms": 13.2, "ttl": 57},
    {"seq": 2, "received": false, "error": "timeout"}
  ]
}
```

Traceroute result：

```json
{
  "summary": {
    "resolved_ip": "1.1.1.1",
    "reached": true,
    "hop_count": 8,
    "path_hash": "sha256:...",
    "protocol": "icmp"
  },
  "samples": [
    {"hop": 1, "ip": "192.168.1.1", "hostname": "router.local", "rtts_ms": [1.2, 1.4, 1.3], "loss_percent": 0},
    {"hop": 2, "ip": null, "hostname": null, "rtts_ms": [], "loss_percent": 100, "error": "timeout"}
  ]
}
```

DNS result：

```json
{
  "summary": {
    "query_name": "example.com",
    "record_type": "A",
    "resolver": "system",
    "transport": "udp",
    "rcode": "NOERROR",
    "answer_count": 2,
    "duration_ms": 23
  },
  "samples": [
    {"name": "example.com.", "type": "A", "ttl": 300, "data": "93.184.216.34"}
  ],
  "raw": {
    "cname_chain": [],
    "truncated": false
  }
}
```

## 設定傳遞與排程

控制器不應該向 probe 傳送個別 job。它應該傳送 desired configuration，probe 再依據 desired state 自行設定。

控制器行為：

- 每次 REST 設定 mutation 都建立新的 `config_versions` row，並 append 受影響的版本化設定 rows。
- 從指定 `config_version` 下有效的 `probes`、`probe_targets`、`probe_target_checks`，以及 `targets` 建立每個 probe 專屬的設定快照。
- 快照包含單調遞增的 `config_version` 與 checksum。
- 當 probe 註冊、重新連線，或要求重新同步時，透過 gRPC 傳送完整快照。
- 當被指派的 target 或檢查變更時，透過 gRPC 串流增量更新。
- 每個 probe 最新 desired config version 可由最新設定版本與該 probe 有效設定推導；probe 實際套用狀態由 `probe_config_acks` 追蹤。
- 追蹤 probe ack，讓操作人員可以看到 probe 是否已套用最新設定。
- 透過 probe gRPC stream 接收觀測到的量測資料。

Probe 行為：

- 在本機保留最新已套用設定。
- 當新的設定快照或更新抵達時，reconcile 本機 schedules。
- 依照每個 `probe_target.interval_seconds` 或 `probe_target_check.schedule_override_seconds` 執行已啟用的檢查。
- 在本機加入 jitter，避免同步量測。
- 串流包含 `probe_target_id`、`check_id`、`target_id`，以及 `config_version` 的結果。
- 若控制器暫時不可達，繼續執行上一份有效設定。

後續改進：

- 使用 Redis、NATS，或 Kafka 為設定更新提供持久化事件傳遞。
- 對大型 probe 設定使用 delta updates，取代完整設定快照。
- 控制器端偵測過期的設定 ack。
- 當 gRPC stream 暫時斷線時，probe 可緩衝結果上傳。
- 當 probe 很慢或離線時提供 backpressure。
