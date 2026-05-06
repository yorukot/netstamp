# Netstamp 計畫

Netstamp 是一套以 Go 建構的分散式網路量測系統，靈感來自 SmokePing、Globalping，以及 RIPE Atlas。第一個里程碑會聚焦在後端：中央控制器與分散式 probe。Probe 會執行已設定的 checks、儲存時間序列結果，並提供 API，供後續 UI、告警，以及公開或共享量測功能使用。

## 控制器

控制器是中央後端服務。它負責持有設定、將設定變更串流給 probe、接收 probe 結果、儲存歷史資料，並暴露 API。

核心職責：

- 註冊並驗證 probe。
- 儲存 probe 中繼資料，例如區域、供應商、主機名稱、標籤、IP family 支援，以及軟體版本。
- 定義 check，例如對某個網域、IP 位址、DNS 查詢目標，或 DNS resolver 執行 ping、traceroute，或 DNS 檢查。
- 定義哪些 probe 應該執行哪些 check。
- 儲存每個 check 的目標與參數，以及每個 probe assignment 的排程設定。第一版只支援 ping、traceroute，以及 DNS。
- 只將每個 probe 被指派的 check 設定傳給該 probe。
- 當某個 probe 被指派的 check 變更時，將更新後的設定串流給該 probe。
- 接收量測結果。
- 儲存原始與彙總後的量測資料。
- 暴露 REST API，供設定與查詢結果使用。
- 當封包遺失、延遲、路徑變化，或 DNS 查詢錯誤跨過閾值時，發出告警事件。

### Probe

Probe 是安裝在 VPS、裸機伺服器，或內部主機上的輕量 Go agent。它會連到控制器、接收被指派的 probe-check 設定、在本機排程這些檢查、執行檢查，並把結果送回。

核心職責：

- 使用 token 向控制器註冊。
- 維持安全的 gRPC 串流，並具備重連與重新同步行為。
- 套用 probe-check 設定快照與增量更新。
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

核心需求是控制器會把 probe 被指派的 `probe_check` + `check` 資料送給該 probe，並在資料變更時串流新的版本。Probe 在本機排程檢查，並透過同一條 gRPC 連線把觀測資料串流回來。

建議的第一版：

- REST 用於管理與面向 UI 的 API。
- gRPC 雙向串流用於 probe 控制與觀測資料。
- 若有幫助，使用 gRPC unary methods 處理註冊、bootstrap 設定抓取，以及明確的重新同步。
- 從第一天開始就對 protobuf 訊息做版本化。
- REST 結果匯入可保留為手動測試或 fallback，不作為主要 probe 路徑。

對這種 probe 設定模型來說，gRPC 比單純 polling 更適合，因為設定變更應該快速抵達 probe，而 probe 也需要頻繁把觀測資料送回。管理 API 保持 REST 仍然合理。務實的分工是：REST 給人與前端使用，gRPC streaming 給 probe 流量使用。

### 使用者與驗證

第一版使用簡單的 email + password 登入，並加入 team 作為資料隔離單位。MVP 只做 team scope 與簡單 membership，不做細緻 per-resource ACL。

- `users.email` 必須唯一。
- 密碼只儲存 password hash，不儲存明文。
- 登入成功後簽發 JWT access token。
- JWT 有效時間固定 12 小時。
- JWT claim 只需要 `sub`、`email`、`iat`、`exp`。
- 第一版不做 refresh token；token 過期後重新登入。
- `team_members.role` 先支援 `owner`、`admin`、`editor`、`viewer`。
- Team 與 team member 管理第一版先使用簡單角色規則：active member 可讀取 team 與 member；`owner`、`admin` 可更新 team 與管理 member；只有 `owner` 可刪除 team；任何人都不能把 member 設為 `owner`，`admin` 也不能把 member 設為 `admin`。
- MVP 可先讓 team 內成員都能操作該 team 的 probes、checks、probe-checks 與 results；之後再收斂 owner/admin/member 的差異。

### REST API

初始 endpoints：

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/auth/me`
- `POST /api/v1/teams`
- `GET /api/v1/teams`
- `GET /api/v1/teams/{ref}`
- `PATCH /api/v1/teams/{ref}`
- `DELETE /api/v1/teams/{ref}`
- `GET /api/v1/teams/{ref}/members`
- `POST /api/v1/teams/{ref}/members`
- `PATCH /api/v1/teams/{ref}/members/{user_id}`
- `POST /api/v1/probes/register`
- `GET /api/v1/probes`
- `GET /api/v1/probes/{id}`
- `PATCH /api/v1/probes/{id}`
- `POST /api/v1/checks`
- `GET /api/v1/checks`
- `GET /api/v1/checks/{id}`
- `PATCH /api/v1/checks/{id}`
- `POST /api/v1/probe-checks`
- `GET /api/v1/probe-checks`
- `GET /api/v1/probe-checks/{id}`
- `PATCH /api/v1/probe-checks/{id}`
- `GET /api/v1/results`
- `GET /api/v1/checks/{id}/results`
- `GET /api/v1/probe-checks/{id}/results`
- `GET /api/v1/alerts`

Team route `{ref}` 可接受 team UUID 或 slug；response 仍同時回傳 `id` 與 `slug`。

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

- Query name；預設可使用 `checks.target`。
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

初始資料表應該以 `probe -> check` 來建模產品，而不是以人工管理的 job 來建模。第一版不建立獨立 target inventory；target 是 check 的欄位。持久化的使用者設定是：

```text
probe + check(target + type + parameters) + schedule
```

控制器儲存這份設定並送給 probe。Probe 使用收到的 schedule 與 parameters 在本機執行檢查。控制器不需要為每個檢查間隔建立或派送個別 job。

命名注意事項：本計畫使用 `probe` 表示執行量測的 agent，使用 `check` 表示一次可被排程的量測定義。`target` 在第一版只是 `checks.target` 欄位，表示被量測的 host、IP、DNS query name，或 resolver address，不作為獨立主表。`anchor` 保留給未來「穩定、公開、可被其他 probe 量測的基準 probe」概念，不作為第一版主表命名。

建議 PostgreSQL / TimescaleDB extensions：

- `pgcrypto`：提供 `gen_random_uuid()`。
- `citext`：讓 email 唯一性不受大小寫影響。
- `timescaledb`：將各類 result tables 建成 hypertables，專門處理時間序列結果。

### MVP 設定歷史原則

第一版不保留設定版本歷史。設定表代表目前 desired state，更新時直接修改 row，刪除時使用 `deleted_at` soft delete。這讓資料模型保持直覺，也讓 `probe_credentials`、`probe_status`、`probe_checks` 等資料表可以使用正常 foreign key。

- `users`、`teams`、`probes`、`checks`、`probe_checks` 使用一般 primary key。
- `target`、`target_type`、`check_type` 第一版允許修改；前端應在修改時明確警告這會影響歷史結果的語意解讀。
- Result tables 不保存設定版本，也不保存 target snapshot；查詢時若 join 到 `checks`，看到的是目前 check 設定。
- 若未來需要 audit、diff、rollback，再新增獨立 history/audit tables，不把 MVP schema 做成 append-only versioned tables。

### `users`

第一版只支援 email + password 登入。`users` 只儲存帳號本身；team membership 與簡單角色放在 `team_members`。

| 欄位            | PostgreSQL type | 限制                                    | 說明                                |
| --------------- | --------------- | --------------------------------------- | ----------------------------------- |
| `id`            | `UUID`          | `PRIMARY KEY DEFAULT gen_random_uuid()` | 使用者 ID                           |
| `email`         | `CITEXT`        | `NOT NULL UNIQUE`                       | 登入 email；需要 `citext` extension |
| `password_hash` | `TEXT`          | `NOT NULL`                              | Argon2id 或 bcrypt 的 encoded hash  |
| `is_active`     | `BOOLEAN`       | `NOT NULL DEFAULT true`                 | 停用帳號用，不代表權限              |
| `created_at`    | `TIMESTAMPTZ`   | `NOT NULL DEFAULT now()`                | 建立時間                            |
| `updated_at`    | `TIMESTAMPTZ`   | `NOT NULL DEFAULT now()`                | 更新時間                            |

### `teams`

Team 是 MVP 的資料隔離單位。所有 probes、checks、probe-checks 與 results 都屬於某個 team。

| 欄位                 | PostgreSQL type | 限制                                    | 說明                                                               |
| -------------------- | --------------- | --------------------------------------- | ------------------------------------------------------------------ |
| `id`                 | `UUID`          | `PRIMARY KEY DEFAULT gen_random_uuid()` | Team ID                                                            |
| `name`               | `TEXT`          | `NOT NULL`                              | 顯示名稱                                                           |
| `slug`               | `CITEXT`        | `NOT NULL UNIQUE`                       | URL/API 可用的穩定識別字；只允許小寫 `a-z`、數字 `0-9` 與 dash `-` |
| `created_by_user_id` | `UUID`          | `NULL REFERENCES users(id)`             | 建立者                                                             |
| `created_at`         | `TIMESTAMPTZ`   | `NOT NULL DEFAULT now()`                | 建立時間                                                           |
| `updated_at`         | `TIMESTAMPTZ`   | `NOT NULL DEFAULT now()`                | 更新時間                                                           |
| `deleted_at`         | `TIMESTAMPTZ`   | `NULL`                                  | soft delete 時間                                                   |

### `team_members`

Team membership 只先做簡單角色，不做細緻權限。MVP 可先讓 team 內成員都能操作 team 內資源。

| 欄位         | PostgreSQL type | 限制                            | 說明                                 |
| ------------ | --------------- | ------------------------------- | ------------------------------------ |
| `team_id`    | `UUID`          | `NOT NULL REFERENCES teams(id)` | Team ID                              |
| `user_id`    | `UUID`          | `NOT NULL REFERENCES users(id)` | User ID                              |
| `role`       | `TEXT`          | `NOT NULL`                      | `owner`、`admin`、`editor`、`viewer` |
| `created_at` | `TIMESTAMPTZ`   | `NOT NULL DEFAULT now()`        | 建立時間                             |
| `updated_at` | `TIMESTAMPTZ`   | `NOT NULL DEFAULT now()`        | 更新時間                             |
| `deleted_at` | `TIMESTAMPTZ`   | `NULL`                          | soft delete 時間                     |

主鍵：`PRIMARY KEY (team_id, user_id)`。

### `probes`

Probe 設定表保存目前設定。若 probe 被移除，使用 `deleted_at` soft delete，而不是 hard delete，避免破壞歷史結果的 reference。

| 欄位         | PostgreSQL type | 限制                                    | 說明                          |
| ------------ | --------------- | --------------------------------------- | ----------------------------- |
| `id`         | `UUID`          | `PRIMARY KEY DEFAULT gen_random_uuid()` | Probe ID                      |
| `team_id`    | `UUID`          | `NOT NULL REFERENCES teams(id)`         | 所屬 team                     |
| `name`       | `TEXT`          | `NOT NULL`                              | 顯示名稱                      |
| `provider`   | `TEXT`          | `NULL`                                  | 例如 `hetzner`、`aws`、`home` |
| `region`     | `TEXT`          | `NULL`                                  | Probe 區域                    |
| `hostname`   | `TEXT`          | `NULL`                                  | Probe 主機名稱                |
| `enabled`    | `BOOLEAN`       | `NOT NULL DEFAULT true`                 | 是否允許 probe 執行檢查       |
| `created_at` | `TIMESTAMPTZ`   | `NOT NULL DEFAULT now()`                | 建立時間                      |
| `updated_at` | `TIMESTAMPTZ`   | `NOT NULL DEFAULT now()`                | 更新時間                      |
| `deleted_at` | `TIMESTAMPTZ`   | `NULL`                                  | soft delete 時間              |

### `probe_credentials`

Probe 憑證不放進 probe 設定快照，也不會被串流給 probe。第一版只需要保存目前有效的 probe secret hash，並透過 FK 綁定到 `probes`。

| 欄位              | PostgreSQL type | 限制                                | 說明                      |
| ----------------- | --------------- | ----------------------------------- | ------------------------- |
| `probe_id`        | `UUID`          | `PRIMARY KEY REFERENCES probes(id)` | Probe ID                  |
| `secret_hash`     | `TEXT`          | `NOT NULL`                          | 長期 probe secret 的 hash |
| `created_at`      | `TIMESTAMPTZ`   | `NOT NULL DEFAULT now()`            | 憑證建立時間              |
| `last_rotated_at` | `TIMESTAMPTZ`   | `NULL`                              | 最近一次輪替時間          |

### `probe_status`

Probe runtime 狀態不是設定歷史，使用目前狀態表即可。後續若要長期保存 health history，再另外建立 time-series 表。

| 欄位            | PostgreSQL type | 限制                                | 說明                                        |
| --------------- | --------------- | ----------------------------------- | ------------------------------------------- |
| `probe_id`      | `UUID`          | `PRIMARY KEY REFERENCES probes(id)` | Probe ID                                    |
| `status`        | `TEXT`          | `NOT NULL`                          | 例如 `online`、`offline`、`draining`        |
| `last_seen_at`  | `TIMESTAMPTZ`   | `NULL`                              | 最後 heartbeat 時間                         |
| `agent_version` | `TEXT`          | `NULL`                              | Probe 軟體版本                              |
| `ip_families`   | `TEXT[]`        | `NOT NULL DEFAULT '{}'::text[]`     | 例如 `{ipv4,ipv6}`                          |
| `capabilities`  | `JSONB`         | `NOT NULL DEFAULT '{}'::jsonb`      | Probe 能力，例如 raw ICMP、DNS TCP fallback |
| `health`        | `JSONB`         | `NOT NULL DEFAULT '{}'::jsonb`      | CPU、記憶體、queue length 等                |

### `checks`

`checks` 表示可被 probe 執行的一個量測定義。第一版直接把 target 與 check 合併在這張表，避免為了簡單檢查維護過多中介表。`target` 是這個 check 的主要目標，不是 foreign key。

| 欄位          | PostgreSQL type | 限制                                    | 說明                                          |
| ------------- | --------------- | --------------------------------------- | --------------------------------------------- |
| `id`          | `UUID`          | `PRIMARY KEY DEFAULT gen_random_uuid()` | Check ID                                      |
| `team_id`     | `UUID`          | `NOT NULL REFERENCES teams(id)`         | 所屬 team                                     |
| `name`        | `TEXT`          | `NOT NULL`                              | 顯示名稱                                      |
| `check_type`  | `TEXT`          | `NOT NULL`                              | `ping`、`traceroute`、`dns`                   |
| `target`      | `TEXT`          | `NOT NULL`                              | host、IP、DNS query name，或 resolver address |
| `target_type` | `TEXT`          | `NOT NULL`                              | `host`、`ip`、`dns_query`、`dns_resolver`     |
| `description` | `TEXT`          | `NULL`                                  | 備註                                          |
| `parameters`  | `JSONB`         | `NOT NULL DEFAULT '{}'::jsonb`          | 檢查參數                                      |
| `enabled`     | `BOOLEAN`       | `NOT NULL DEFAULT true`                 | 是否允許此 check 被執行                       |
| `created_at`  | `TIMESTAMPTZ`   | `NOT NULL DEFAULT now()`                | 建立時間                                      |
| `updated_at`  | `TIMESTAMPTZ`   | `NOT NULL DEFAULT now()`                | 更新時間                                      |
| `deleted_at`  | `TIMESTAMPTZ`   | `NULL`                                  | soft delete 時間                              |

`target`、`target_type`、`check_type` 可以被更新，但前端應把這個操作標示為不建議，因為同一個 check 的歷史圖表可能會混入不同語意的量測對象。

範例資料列：

| id                | check_type   | target        | target_type | parameters                                                                |
| ----------------- | ------------ | ------------- | ----------- | ------------------------------------------------------------------------- |
| `check_ping_aaa`  | `ping`       | `1.1.1.1`     | `ip`        | `{"count":20,"timeout_ms":3000,"ip_version":"auto","packet_size":56}`     |
| `check_trace_bbb` | `traceroute` | `example.com` | `host`      | `{"max_hops":30,"queries_per_hop":3,"timeout_ms":3000,"protocol":"icmp"}` |
| `check_dns_ccc`   | `dns`        | `example.com` | `dns_query` | `{"record_type":"A","resolver":"system","timeout_ms":2000,"attempts":2}`  |

### `probe_checks`

`probe_checks` 表示哪個 probe 應該執行哪個 check，以及該 probe 上的排程設定。這張表保留 probe assignment，讓同一個 check 可以被多個 probe 執行。

| 欄位               | PostgreSQL type | 限制                                    | 說明                              |
| ------------------ | --------------- | --------------------------------------- | --------------------------------- |
| `id`               | `UUID`          | `PRIMARY KEY DEFAULT gen_random_uuid()` | Probe-check assignment ID         |
| `team_id`          | `UUID`          | `NOT NULL REFERENCES teams(id)`         | 所屬 team                         |
| `probe_id`         | `UUID`          | `NOT NULL REFERENCES probes(id)`        | Probe ID                          |
| `check_id`         | `UUID`          | `NOT NULL REFERENCES checks(id)`        | Check ID                          |
| `name`             | `TEXT`          | `NULL`                                  | 覆寫顯示名稱                      |
| `enabled`          | `BOOLEAN`       | `NOT NULL DEFAULT true`                 | 是否啟用此 probe-check assignment |
| `interval_seconds` | `INTEGER`       | `NOT NULL`                              | 排程間隔                          |
| `jitter_seconds`   | `INTEGER`       | `NOT NULL DEFAULT 0`                    | Probe 本機 jitter 上限            |
| `created_at`       | `TIMESTAMPTZ`   | `NOT NULL DEFAULT now()`                | 建立時間                          |
| `updated_at`       | `TIMESTAMPTZ`   | `NOT NULL DEFAULT now()`                | 更新時間                          |
| `deleted_at`       | `TIMESTAMPTZ`   | `NULL`                                  | soft delete 時間                  |

`checks` 與 `probe_checks` 分開是刻意的：`checks` 是「量測什麼」與參數，`probe_checks` 是「哪個 probe 何時執行」。同一個 check 可以被多個 probe 重用，也可以讓不同 probe 對同一個 check 使用不同 interval、jitter、enabled 狀態，並讓查詢可以自然地看「某個 check 在所有 probes 的結果」或「某個 probe-check assignment 的結果」。

建議限制同一個 active probe-check assignment 不重複：`UNIQUE (team_id, probe_id, check_id) WHERE deleted_at IS NULL`。

範例資料列：

| id       | probe_id     | check_id         | enabled | interval_seconds |
| -------- | ------------ | ---------------- | ------- | ---------------- |
| `pc_aaa` | `probe_1111` | `check_ping_aaa` | true    | 60               |
| `pc_bbb` | `probe_1111` | `check_dns_ccc`  | true    | 300              |

### Result tables

量測結果是主要長期資料，所以第一版不使用單一 generic `results` table。每種 check type 使用自己的 TimescaleDB hypertable，讓常用欄位成為一等 SQL 欄位，方便 ChartDB、API、告警與後續 Timescale aggregation 查詢。

共同設計原則：

- 所有 result tables 都是 immutable time-series tables，以 `started_at` 作為 hypertable 時間分區欄位。
- 所有 result tables 都只保存 `team_id` 與 `probe_check_id` 作為設定 reference，並使用 `(team_id, probe_check_id)` foreign key 指到 `probe_checks(team_id, id)`。
- 需要 probe、check、target 或 schedule 時，從 `probe_checks -> probes/checks` join 回去，不在每張 result table 重複保存 `probe_id` 與 `check_id`。
- MVP 不保存設定版本，也不保存 target snapshot；查詢時若 join 到 `checks`，看到的是目前 check 設定。
- `raw JSONB` 只作為 debug 或原始 protocol metadata 補充，不作為主要查詢模型。

### `ping_results`

`ping_results` 每筆 row 表示一次 ping check run 的摘要與 RTT samples。

| 欄位             | PostgreSQL type      | 限制                                        | 說明                                     |
| ---------------- | -------------------- | ------------------------------------------- | ---------------------------------------- |
| `id`             | `UUID`               | `NOT NULL DEFAULT gen_random_uuid()`        | Result ID                                |
| `team_id`        | `UUID`               | `NOT NULL`                                  | 所屬 team；與 `probe_check_id` 組成 FK   |
| `probe_check_id` | `UUID`               | `NOT NULL`                                  | 執行的 probe-check assignment            |
| `started_at`     | `TIMESTAMPTZ`        | `NOT NULL`                                  | Probe 開始執行時間                       |
| `finished_at`    | `TIMESTAMPTZ`        | `NOT NULL`                                  | Probe 完成時間                           |
| `duration_ms`    | `INTEGER`            | `NOT NULL`                                  | `finished_at - started_at` 的毫秒數      |
| `status`         | `TEXT`               | `NOT NULL`                                  | `success`、`partial`、`timeout`、`error` |
| `sent_count`     | `INTEGER`            | `NOT NULL`                                  | 已送出封包數                             |
| `received_count` | `INTEGER`            | `NOT NULL`                                  | 已收到封包數                             |
| `loss_percent`   | `DOUBLE PRECISION`   | `NOT NULL`                                  | 封包遺失百分比                           |
| `rtt_min_ms`     | `DOUBLE PRECISION`   | `NULL`                                      | 最小 RTT                                 |
| `rtt_avg_ms`     | `DOUBLE PRECISION`   | `NULL`                                      | 平均 RTT                                 |
| `rtt_median_ms`  | `DOUBLE PRECISION`   | `NULL`                                      | 中位數 RTT                               |
| `rtt_max_ms`     | `DOUBLE PRECISION`   | `NULL`                                      | 最大 RTT                                 |
| `rtt_stddev_ms`  | `DOUBLE PRECISION`   | `NULL`                                      | RTT 標準差                               |
| `rtt_samples_ms` | `DOUBLE PRECISION[]` | `NOT NULL DEFAULT '{}'::double precision[]` | 同一次 ping 的 RTT samples               |
| `resolved_ip`    | `INET`               | `NULL`                                      | 解析後目標 IP                            |
| `ip_version`     | `TEXT`               | `NULL`                                      | `ipv4`、`ipv6`、`auto`                   |
| `raw`            | `JSONB`              | `NOT NULL DEFAULT '{}'::jsonb`              | 原始或補充資料                           |
| `error_code`     | `TEXT`               | `NULL`                                      | 可機器判讀的錯誤碼                       |
| `error_message`  | `TEXT`               | `NULL`                                      | 人類可讀錯誤                             |
| `created_at`     | `TIMESTAMPTZ`        | `NOT NULL DEFAULT now()`                    | 控制器寫入時間                           |

### `dns_results`

`dns_results` 每筆 row 表示一次 DNS query check run。DNS answers 第一版保留在 JSONB，因為不同 record type 的資料形狀差異較大。

| 欄位               | PostgreSQL type | 限制                                 | 說明                                       |
| ------------------ | --------------- | ------------------------------------ | ------------------------------------------ |
| `id`               | `UUID`          | `NOT NULL DEFAULT gen_random_uuid()` | Result ID                                  |
| `team_id`          | `UUID`          | `NOT NULL`                           | 所屬 team；與 `probe_check_id` 組成 FK     |
| `probe_check_id`   | `UUID`          | `NOT NULL`                           | 執行的 probe-check assignment              |
| `started_at`       | `TIMESTAMPTZ`   | `NOT NULL`                           | Probe 開始執行時間                         |
| `finished_at`      | `TIMESTAMPTZ`   | `NOT NULL`                           | Probe 完成時間                             |
| `duration_ms`      | `INTEGER`       | `NOT NULL`                           | `finished_at - started_at` 的毫秒數        |
| `status`           | `TEXT`          | `NOT NULL`                           | `success`、`partial`、`timeout`、`error`   |
| `query_name`       | `TEXT`          | `NOT NULL`                           | DNS query name                             |
| `record_type`      | `TEXT`          | `NOT NULL`                           | `A`、`AAAA`、`CNAME`、`MX`、`NS`、`TXT` 等 |
| `resolver`         | `TEXT`          | `NOT NULL`                           | `system` 或 resolver address               |
| `transport`        | `TEXT`          | `NOT NULL`                           | `udp` 或 `tcp`                             |
| `rcode`            | `TEXT`          | `NULL`                               | DNS response code                          |
| `success`          | `BOOLEAN`       | `NOT NULL`                           | DNS query 是否成功                         |
| `answer_count`     | `INTEGER`       | `NOT NULL DEFAULT 0`                 | answer 數量                                |
| `response_time_ms` | `INTEGER`       | `NOT NULL`                           | DNS query 耗時                             |
| `answers`          | `JSONB`         | `NOT NULL DEFAULT '[]'::jsonb`       | DNS answers                                |
| `raw`              | `JSONB`         | `NOT NULL DEFAULT '{}'::jsonb`       | 原始 DNS metadata                          |
| `error_code`       | `TEXT`          | `NULL`                               | 可機器判讀的錯誤碼                         |
| `error_message`    | `TEXT`          | `NULL`                               | 人類可讀錯誤                               |
| `created_at`       | `TIMESTAMPTZ`   | `NOT NULL DEFAULT now()`             | 控制器寫入時間                             |

### `traceroute_results`

`traceroute_results` 每筆 row 表示一次 traceroute run 的摘要。路由明細拆到 `traceroute_hops`，讓 hop-level history 可以被 SQL、ChartDB、告警與後續 path diff 直接查詢。

| 欄位             | PostgreSQL type | 限制                                 | 說明                                     |
| ---------------- | --------------- | ------------------------------------ | ---------------------------------------- |
| `id`             | `UUID`          | `NOT NULL DEFAULT gen_random_uuid()` | Result ID                                |
| `team_id`        | `UUID`          | `NOT NULL`                           | 所屬 team；與 `probe_check_id` 組成 FK   |
| `probe_check_id` | `UUID`          | `NOT NULL`                           | 執行的 probe-check assignment            |
| `started_at`     | `TIMESTAMPTZ`   | `NOT NULL`                           | Probe 開始執行時間                       |
| `finished_at`    | `TIMESTAMPTZ`   | `NOT NULL`                           | Probe 完成時間                           |
| `duration_ms`    | `INTEGER`       | `NOT NULL`                           | `finished_at - started_at` 的毫秒數      |
| `status`         | `TEXT`          | `NOT NULL`                           | `success`、`partial`、`timeout`、`error` |
| `resolved_ip`    | `INET`          | `NULL`                               | 解析後目標 IP                            |
| `reached`        | `BOOLEAN`       | `NOT NULL DEFAULT false`             | 是否抵達終點                             |
| `hop_count`      | `INTEGER`       | `NOT NULL DEFAULT 0`                 | hop 數量                                 |
| `path_hash`      | `TEXT`          | `NULL`                               | 用於偵測路由變化                         |
| `protocol`       | `TEXT`          | `NOT NULL`                           | 第一版為 `icmp`                          |
| `raw`            | `JSONB`         | `NOT NULL DEFAULT '{}'::jsonb`       | 原始或補充資料                           |
| `error_code`     | `TEXT`          | `NULL`                               | 可機器判讀的錯誤碼                       |
| `error_message`  | `TEXT`          | `NULL`                               | 人類可讀錯誤                             |
| `created_at`     | `TIMESTAMPTZ`   | `NOT NULL DEFAULT now()`             | 控制器寫入時間                           |

### `traceroute_hops`

`traceroute_hops` 每筆 row 表示某次 traceroute run 的一個 hop。這張表保留 hop-level RTT、loss、timeout 與 IP/hostname，避免把路由分析埋在 JSONB 裡。

TimescaleDB 不支援 hypertable 對 hypertable 建立 foreign key，因此 `traceroute_hops.traceroute_result_id` 與 `started_at` 會對應 `traceroute_results.id` 與 `started_at`，但 parent/hops 的存在性由 controller 在同一個 transaction 寫入時保證。`traceroute_hops` 仍使用 `(team_id, probe_check_id)` 指到 `probe_checks`，需要 probe 或 target 時再 join 回設定表。

| 欄位                   | PostgreSQL type      | 限制                                        | 說明                                   |
| ---------------------- | -------------------- | ------------------------------------------- | -------------------------------------- |
| `traceroute_result_id` | `UUID`               | `NOT NULL`                                  | 對應的 traceroute result               |
| `team_id`              | `UUID`               | `NOT NULL`                                  | 所屬 team；與 `probe_check_id` 組成 FK |
| `probe_check_id`       | `UUID`               | `NOT NULL`                                  | 執行的 probe-check assignment          |
| `started_at`           | `TIMESTAMPTZ`        | `NOT NULL`                                  | 所屬 traceroute run 的開始時間         |
| `hop_number`           | `INTEGER`            | `NOT NULL`                                  | hop 序號                               |
| `hop_ip`               | `INET`               | `NULL`                                      | hop IP；timeout 時可為 NULL            |
| `hostname`             | `TEXT`               | `NULL`                                      | hop hostname                           |
| `rtts_ms`              | `DOUBLE PRECISION[]` | `NOT NULL DEFAULT '{}'::double precision[]` | 此 hop 的 RTT samples                  |
| `loss_percent`         | `DOUBLE PRECISION`   | `NOT NULL`                                  | 此 hop 的封包遺失                      |
| `error_code`           | `TEXT`               | `NULL`                                      | 可機器判讀的錯誤碼                     |
| `error_message`        | `TEXT`               | `NULL`                                      | 人類可讀錯誤                           |
| `created_at`           | `TIMESTAMPTZ`        | `NOT NULL DEFAULT now()`                    | 控制器寫入時間                         |

建議 result indexes：

- 每張 result table 都建立 `(team_id, started_at DESC)` 與 `(probe_check_id, started_at DESC)`。
- `traceroute_results` 另外建立 `(team_id, path_hash, started_at DESC)`，支援路由變化查詢。
- `traceroute_hops` 另外建立 `(traceroute_result_id, started_at DESC)`、`(team_id, hop_ip, started_at DESC)` 與 `(team_id, probe_check_id, hop_number, started_at DESC)`，支援 parent lookup 與 hop-level history。

## 設定傳遞與排程

控制器不應該向 probe 傳送個別 job。它應該傳送 desired configuration，probe 再依據 desired state 自行設定。

控制器行為：

- 每次 REST 設定 mutation 直接更新 team scope 內的 `probes`、`checks`，或 `probe_checks` 目前 row。
- 從目前有效的 `probes`、`probe_checks`，以及 `checks` 建立每個 probe 專屬的設定快照。
- 快照可包含 controller 產生的 `snapshot_hash`，讓 probe 判斷本機設定是否需要重新同步。
- 當 probe 註冊、重新連線，或要求重新同步時，透過 gRPC 傳送完整快照。
- 當被指派的 check 或排程變更時，透過 gRPC 串流增量更新。
- 透過 probe gRPC stream 接收觀測到的量測資料。

Probe 行為：

- 在本機保留最新已套用設定。
- 當新的設定快照或更新抵達時，reconcile 本機 schedules。
- 依照每個 `probe_check.interval_seconds` 執行已啟用的檢查。
- 在本機加入 jitter，避免同步量測。
- 串流包含 `team_id`、`probe_check_id`、執行時間、狀態，以及量測 payload 的結果；controller 可由 `probe_check_id` 找回 probe 與 check/target。
- 若控制器暫時不可達，繼續執行上一份有效設定。

後續改進：

- 使用 Redis、NATS，或 Kafka 為設定更新提供持久化事件傳遞。
- 對大型 probe 設定使用 delta updates，取代完整設定快照。
- 若之後需要操作可觀測性，再加入設定 ack 或 probe config apply event。
- 當 gRPC stream 暫時斷線時，probe 可緩衝結果上傳。
- 當 probe 很慢或離線時提供 backpressure。
