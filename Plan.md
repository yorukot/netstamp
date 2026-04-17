# Netstamp Plan

Netstamp is a Go-based distributed network measurement system inspired by SmokePing and Globalping. The first milestone focuses on the backend: a central controller and distributed probe nodes that measure configured ends, store time-series results, and expose APIs for later UI, alerting, and public/shared measurement features.

## Goals

- Build a modern SmokePing-like backend with distributed nodes.
- Support `ping` as the first latency and packet-loss measurement.
- Support `traceroute` as a lighter diagnostic function similar in spirit to Globalping probes.
- Keep the architecture open for more measurement types later, such as DNS, HTTP, TCP connect, TLS, MTR-style path checks, and scheduled multi-node tests.
- Use Go for both the controller and node services.
- Start backend-first: controller API, node agent, node-to-end connection configuration, streamed node configuration, result ingestion, storage, and basic operational security.

## Non-Goals For The First Milestone

- Full frontend dashboard.
- Browser-based live map.
- Advanced SmokePing graph rendering.
- Replacing mature observability systems on day one.

These can be added after the controller and node protocol are stable.

## Product Shape

### Controller

The controller is the central backend service. It owns configuration, streams configuration changes to nodes, receives node results, stores historical data, and exposes APIs.

Core responsibilities:

- Register and authenticate nodes.
- Store node metadata such as region, provider, hostname, tags, IP family support, and software version.
- Define ends such as domains, IP addresses, URLs, DNS nameservers, or named services.
- Define which nodes should measure which ends.
- Store per-connection measurement configuration, for example ping, traceroute, DNS, HTTP, TCP connect, or TLS checks.
- Send each node only its assigned node-to-end check configuration and end data.
- Stream updated configuration to a node when its assigned ends or checks change.
- Receive measurement results.
- Store raw and aggregated measurement data.
- Expose REST APIs for configuration and querying results.
- Emit alert events when packet loss, latency, or path changes cross thresholds.

### Node

The node is a lightweight Go agent installed on VPS, bare-metal servers, or internal hosts. It connects to the controller, receives its assigned node-to-end check configuration, schedules those checks locally, executes them, and sends results back.

Core responsibilities:

- Register with the controller using a token.
- Maintain a secure gRPC stream with reconnect and resync behavior.
- Apply node-to-end configuration snapshots and incremental updates.
- Schedule checks locally from the received configuration.
- Execute ping checks.
- Execute traceroute checks.
- Enforce local safety limits so the node cannot be abused as a scanner.
- Return structured results with timing, errors, and metadata.
- Self-report health metrics such as CPU, memory, version, uptime, queue length, and supported capabilities.

## Architecture

```text
+-----------------+   gRPC node stream       +-----------------+
|                 | <----------------------> |                 |
|   Controller    |                          |      Node       |
|                 |                          |                 |
+--------+--------+                          +--------+--------+
         |                                            |
         |                                            |
         v                                            v
+-----------------+                          +-----------------+
|   PostgreSQL    |                          | ping/traceroute |
|   Timeseries    |                          | OS networking   |
+-----------------+                          +-----------------+
         |
         v
+-----------------+
| Redis / Queue   |
| optional later  |
+-----------------+
```

Recommended first implementation:

- Single Go module with separate binaries:
  - `cmd/controller`
  - `cmd/node`
- Shared packages:
  - `internal/api`
  - `internal/auth`
  - `internal/config`
  - `internal/measure`
  - `internal/store`
  - `internal/nodescheduler`
  - `internal/protocol`
- PostgreSQL for durable relational data.
- TimescaleDB extension if long-term time-series scale becomes important.
- Start with node-managed scheduling from controller-supplied configuration.
- Add Redis, NATS, or another queue only when config fan-out, event delivery, or ingestion buffering needs it.

## Controller API

Use REST over HTTP for admin/configuration APIs because it is simple to debug and easy for a future frontend to consume. Use gRPC for controller-node traffic because this system needs a long-lived node channel for two directions:

- Controller to node: configuration snapshots and incremental updates.
- Node to controller: observed measurement data, config acknowledgements, heartbeat, and health.

The core requirement is that the controller sends a node its assigned `node_end` + `node_end_check` + `end` data, then streams a new version whenever that data changes. The node schedules checks locally and streams observe data back over the same gRPC connection.

Recommended first version:

- REST for admin and UI-facing APIs.
- gRPC bidirectional stream for node control and observed data.
- gRPC unary methods for registration, bootstrap config fetch, and explicit resync if useful.
- Version protobuf messages from day one.
- Keep REST result ingestion optional for manual testing or fallback, not as the primary node path.

gRPC is a better fit than plain polling for this node configuration model because config changes should reach nodes quickly and nodes need to send frequent observe data back. It is still reasonable to keep admin APIs REST. A pragmatic split is REST for humans/frontends, and gRPC streaming for node traffic.

### Public/Admin APIs

Initial endpoints:

- `POST /api/v1/nodes/register`
- `GET /api/v1/nodes`
- `GET /api/v1/nodes/{id}`
- `PATCH /api/v1/nodes/{id}`
- `POST /api/v1/ends`
- `GET /api/v1/ends`
- `GET /api/v1/ends/{id}`
- `PATCH /api/v1/ends/{id}`
- `POST /api/v1/node-ends`
- `GET /api/v1/node-ends`
- `GET /api/v1/node-ends/{id}`
- `PATCH /api/v1/node-ends/{id}`
- `POST /api/v1/node-ends/{id}/checks`
- `GET /api/v1/node-ends/{id}/checks`
- `PATCH /api/v1/node-end-checks/{check_id}`
- `GET /api/v1/results`
- `GET /api/v1/node-ends/{id}/results`
- `GET /api/v1/alerts`

### Node gRPC APIs

Initial node-facing RPCs:

- `RegisterNode(RegisterNodeRequest) returns (RegisterNodeResponse)`
- `GetNodeConfig(GetNodeConfigRequest) returns (NodeConfigSnapshot)`
- `Observe(stream NodeObserveMessage) returns (stream ControllerMessage)`

The `Observe` stream should carry:

- Controller messages: config snapshot, config delta, resync request, shutdown/drain instruction.
- Node messages: heartbeat, config ack, ping result, traceroute result, DNS result, health metrics, log event.

Node authentication:

- Node receives a one-time registration token.
- Controller returns a node ID and long-lived node secret.
- Every node request uses an HMAC signature or mTLS.
- First version can use signed bearer tokens; mTLS can be added when operating many untrusted nodes.

## Measurement Types

### Ping

Ping should collect SmokePing-style latency and packet-loss data.

Required fields:

- End host or address.
- IP version: IPv4, IPv6, or auto.
- Packet count.
- Interval between packets.
- Timeout.
- Packet size.
- Source interface optional.
- Privileged raw ICMP mode or unprivileged system command fallback.

Result fields:

- Sent packet count.
- Received packet count.
- Packet loss percentage.
- Minimum latency.
- Average latency.
- Median latency.
- Maximum latency.
- Standard deviation.
- Per-packet latency samples.
- Resolved IP address.
- Error message if failed.

Implementation choices:

- Prefer a Go ICMP implementation using `golang.org/x/net/icmp`.
- Support fallback to system `ping` only where raw socket permissions are unavailable.
- Store both summary data and samples so SmokePing-like charts can be generated later.

### Traceroute

Traceroute should be a lightweight diagnostic function similar to Globalping-style checks, but integrated into the SmokePing-like historical platform.

Required fields:

- End host or address.
- IP version: IPv4, IPv6, or auto.
- Max hops.
- Queries per hop.
- Timeout per hop.
- Protocol: ICMP first, UDP/TCP later.
- Destination port for UDP/TCP modes later.

Result fields:

- Resolved end IP.
- Hop number.
- Hop IP.
- Hop hostname optional.
- RTT samples per hop.
- Packet loss per hop.
- Error or timeout per hop.
- Final reached flag.
- Path hash for detecting route changes.

Implementation choices:

- Start with system `traceroute`/`tracert` parser only if raw implementation is too slow to build.
- Preferred long-term approach is native Go traceroute for portability and structured results.
- Store full hop data as JSONB and also store a path hash for fast path-change detection.

## Data Model

Initial tables should model the product as `node -> end`, not as manually managed jobs. The durable user configuration is:

```text
node + end + enabled checks + schedule + parameters
```

The controller stores this configuration and sends it to the node. The node uses the received schedule and parameters to run checks locally. The controller does not need to create or dispatch individual jobs for every check interval.

Naming note: avoid a table literally named `end` because `END` is a SQL keyword in many contexts. Use `ends`, `measurement_ends`, or `destinations`. This plan uses `ends`.

### `nodes`

- `id`
- `name`
- `region`
- `country`
- `provider`
- `tags`
- `capabilities`
- `version`
- `status`
- `last_seen_at`
- `config_version`
- `created_at`
- `updated_at`

### `ends`

- `id`
- `name`
- `address`
- `end_type`
- `tags`
- `created_at`
- `updated_at`

Examples:

- `example.com` as `dns_name`
- `1.1.1.1` as `ip`
- `https://example.com` as `url`
- `tcp://example.com:443` as `tcp_service`

### `node_ends`

This is the connection table that says which node measures which end.

- `id`
- `node_id`
- `end_id`
- `name`
- `enabled`
- `schedule`
- `jitter`
- `config_version`
- `tags`
- `created_at`
- `updated_at`

Example rows:

| id | node_id | end_id | enabled | schedule |
| --- | --- | --- | --- | --- |
| `ne_aaa` | `node_1111` | `end_9999` | true | `30s` |
| `ne_bbb` | `node_1111` | `end_8888` | true | `60s` |

### `node_end_checks`

This table stores what each node-to-end connection should run. One `node_end` can run several checks against the same end.

- `id`
- `node_end_id`
- `check_type`
- `enabled`
- `parameters`
- `schedule_override`
- `config_version`
- `created_at`
- `updated_at`

Example rows:

| node_end_id | check_type | parameters |
| --- | --- | --- |
| `ne_aaa` | `ping` | `{"count":20,"timeout_ms":3000,"ip_version":"auto"}` |
| `ne_aaa` | `traceroute` | `{"max_hops":30,"queries_per_hop":3,"timeout_ms":3000}` |
| `ne_bbb` | `dns` | `{"record_type":"A","resolver":"system","timeout_ms":2000}` |

### `node_config_versions`

This table tracks what configuration version a node should be running.

- `id`
- `node_id`
- `version`
- `checksum`
- `created_at`
- `applied_at`
- `acked_at`
- `error`

### `measurement_executions`

Executions are reported by nodes after they run a check locally. This table is optional for the first version if `results` is enough, but it is useful when tracking retries, local scheduler drift, and execution errors separately from measurement payloads.

- `id`
- `node_end_id`
- `check_id`
- `node_id`
- `end_id`
- `check_type`
- `status`
- `config_version`
- `scheduled_for`
- `started_at`
- `finished_at`
- `attempt`
- `error`

### `results`

- `id`
- `execution_id`
- `node_end_id`
- `check_id`
- `node_id`
- `end_id`
- `check_type`
- `config_version`
- `started_at`
- `finished_at`
- `summary`
- `samples`
- `raw`
- `created_at`

### `alerts`

- `id`
- `node_end_id`
- `check_id`
- `node_id`
- `end_id`
- `severity`
- `state`
- `reason`
- `opened_at`
- `closed_at`

## Configuration Delivery And Scheduling

The controller should not send individual jobs to nodes. It should send desired configuration, and the node should self-configure from that desired state.

Controller behavior:

- Build a per-node configuration snapshot from `node_ends`, `node_end_checks`, and `ends`.
- Include a monotonically increasing `config_version` and checksum.
- Send the full snapshot over gRPC when a node registers, reconnects, or requests resync.
- Stream incremental updates over gRPC when assigned ends or checks change.
- Store the latest desired config version for each node.
- Track node acknowledgements so operators can see whether a node has applied the latest config.
- Receive observed measurement data over the node gRPC stream.

Node behavior:

- Keep the latest applied configuration locally.
- Reconcile local schedules when a new config snapshot or update arrives.
- Run enabled checks according to each `node_end` schedule or `node_end_check.schedule_override`.
- Add jitter locally to avoid synchronized measurements.
- Stream results with `node_end_id`, `check_id`, `end_id`, and `config_version`.
- Continue running the last valid config if the controller is temporarily unreachable.

Later improvements:

- Durable event delivery for config updates with Redis, NATS, or Kafka.
- Delta updates instead of full config snapshots for large node configs.
- Controller-side detection for stale config acknowledgement.
- Buffered result upload from nodes when the gRPC stream is temporarily disconnected.
- Backpressure when nodes are slow or offline.

## Result Storage And Aggregation

Store raw results first. Add aggregation once the first measurements are stable.

Aggregation levels:

- Raw samples for recent data.
- 5-minute rollups.
- 1-hour rollups.
- 1-day rollups.

Ping aggregation:

- Loss percentage.
- Median latency.
- P95 latency.
- Minimum and maximum latency.
- Jitter or standard deviation.

Traceroute aggregation:

- Most common path hash.
- Path-change count.
- Hop count.
- Per-hop median RTT.
- Per-hop timeout percentage.

## Security And Abuse Prevention

Network measurement systems can be abused if not constrained. The first version should include basic safety controls.

Controller controls:

- Require authentication for admin APIs.
- Require node authentication for node APIs.
- Validate ends before scheduling.
- Block private ranges by default for public nodes.
- Allow private ranges only for explicitly trusted/internal nodes.
- Rate-limit end creation, node-to-end connection changes, config updates, and result ingestion.
- Keep an audit log for end creation, node-to-end configuration changes, and node registration.

Node controls:

- Only apply configuration from its configured controller.
- Enforce maximum packet count, max hops, timeout, and frequency locally.
- Refuse unsupported measurement types.
- Optionally restrict allowed end ranges.
- Run with least privilege where possible.

## Go Implementation Plan

### Phase 1: Project Foundation

- Create Go module.
- Add `cmd/controller` and `cmd/node`.
- Add config loading from environment variables and optional config files.
- Add structured logging.
- Add graceful shutdown.
- Add health endpoints.
- Add protobuf definitions for node protocol.
- Add Dockerfiles for controller and node.
- Add local `docker-compose.yml` with PostgreSQL.

Deliverable:

- Controller and node binaries start cleanly.
- Controller can connect to PostgreSQL.
- Node can load config and call controller health endpoint.

### Phase 2: Controller Core

- Implement database migrations.
- Implement node registration.
- Implement node heartbeat.
- Implement end CRUD.
- Implement node-to-end connection CRUD.
- Implement per-connection check configuration.
- Implement per-node config snapshot generation.
- Implement config version tracking and node acknowledgement.
- Implement gRPC node stream for config updates and observed data.
- Implement result ingestion from node observe messages.

Deliverable:

- A registered node can receive its assigned config over gRPC and stream a fake result back.

### Phase 3: Node Core

- Implement node registration flow.
- Implement heartbeat loop.
- Implement gRPC observe stream.
- Implement config snapshot fetch and streamed config updates.
- Implement local schedule reconciliation.
- Implement worker pool.
- Implement observed data streaming with retry/buffering.
- Add local measurement timeout handling.

Deliverable:

- Node receives controller config, schedules a placeholder measurement locally, and streams a result back.

### Phase 4: Ping Measurement

- Implement ICMP ping runner.
- Add IPv4 and IPv6 support.
- Add packet count, timeout, interval, and packet size parameters.
- Convert packet samples into summary stats.
- Store ping results.
- Add unit tests for summary calculations.
- Add integration test with localhost or a controlled end.

Deliverable:

- Controller sends ping check config for enabled node-to-end connections.
- Node executes ping and streams real packet-loss and latency results.

### Phase 5: Traceroute Measurement

- Implement traceroute runner.
- Add max hops, queries per hop, timeout, and IP version parameters.
- Generate structured hop results.
- Generate path hash.
- Store traceroute results.
- Add unit tests for path hashing and result normalization.
- Add integration test where supported by the host OS.

Deliverable:

- Controller sends traceroute check config for enabled node-to-end connections.
- Node executes traceroute and streams structured hop results.

### Phase 6: Query APIs

- Implement result listing by node-end connection, end, node, time range, and check type.
- Implement latest status endpoint for each node-to-end check.
- Implement basic rollup query endpoint.
- Add pagination.
- Add indexes for common result queries.

Deliverable:

- A future frontend can query recent and historical ping/traceroute results.

### Phase 7: Alerting Foundation

- Add simple threshold rules:
  - Packet loss above percentage.
  - Median latency above threshold.
  - End fully down.
  - Traceroute path changed.
- Add alert state machine.
- Add webhook notification output.

Deliverable:

- Controller can open and close alerts based on stored results.

## Suggested Repository Layout

```text
.
├── cmd
│   ├── controller
│   │   └── main.go
│   └── node
│       └── main.go
├── internal
│   ├── api
│   ├── auth
│   ├── config
│   ├── measure
│   │   ├── ping
│   │   └── traceroute
│   ├── protocol
│   ├── nodescheduler
│   └── store
├── proto
│   └── netstamp
│       └── node
│           └── v1
│               └── node.proto
├── migrations
├── deployments
│   ├── controller.Dockerfile
│   ├── node.Dockerfile
│   └── docker-compose.yml
├── docs
│   ├── api.md
│   ├── node-protocol.md
│   └── security.md
├── go.mod
├── go.sum
├── LICENSE
└── Plan.md
```

## Suggested Libraries

- HTTP router: `chi` or standard library `net/http`.
- gRPC: `google.golang.org/grpc`.
- Protobuf: `google.golang.org/protobuf`.
- Database: `pgx`.
- Migrations: `goose` or `golang-migrate`.
- Logging: `slog`.
- Config: environment variables first, optional `koanf` or `viper` later.
- ICMP: `golang.org/x/net/icmp`.
- Testing: standard `testing`, plus `testcontainers-go` later if needed.

Use the Go standard library where it is enough. Add dependencies only when they remove meaningful complexity.

## First Node gRPC Protocol Draft

Sketch:

```protobuf
service NodeService {
  rpc RegisterNode(RegisterNodeRequest) returns (RegisterNodeResponse);
  rpc GetNodeConfig(GetNodeConfigRequest) returns (NodeConfigSnapshot);
  rpc Observe(stream NodeObserveMessage) returns (stream ControllerMessage);
}
```

`Observe` is the main long-lived channel. The node sends heartbeat, config acknowledgements, health, logs, and observed measurement results. The controller sends config snapshots, config deltas, and resync commands.

### Config Snapshot Request

```json
{
  "node_id": "node_123",
  "capabilities": ["ping", "traceroute"],
  "current_config_version": 41
}
```

### Config Snapshot Response

```json
{
  "config_version": 42,
  "checksum": "sha256:example",
  "node_ends": [
    {
      "id": "ne_aaa",
      "schedule": "30s",
      "jitter_ms": 5000,
      "enabled": true,
      "end": {
        "id": "end_9999",
        "name": "Example",
        "address": "example.com",
        "end_type": "dns_name"
      },
      "checks": [
        {
          "id": "check_ping_123",
          "check_type": "ping",
          "enabled": true,
          "schedule_override": null,
          "parameters": {
            "count": 20,
            "interval_ms": 1000,
            "timeout_ms": 3000,
            "ip_version": "auto"
          }
        },
        {
          "id": "check_trace_123",
          "check_type": "traceroute",
          "enabled": true,
          "schedule_override": "5m",
          "parameters": {
            "max_hops": 30,
            "queries_per_hop": 3,
            "timeout_ms": 3000,
            "ip_version": "auto"
          }
        }
      ]
    }
  ]
}
```

### Config Ack

```json
{
  "node_id": "node_123",
  "config_version": 42,
  "applied_at": "2026-04-17T10:00:00Z",
  "status": "applied"
}
```

### Observe Stream Messages

Node to controller:

```json
{
  "message_type": "ping_result",
  "node_id": "node_123",
  "sent_at": "2026-04-17T10:00:21Z",
  "ping_result": {
    "execution_id": "exec_123",
    "node_end_id": "ne_aaa",
    "check_id": "check_ping_123",
    "check_type": "ping",
    "config_version": 42
  }
}
```

Controller to node:

```json
{
  "message_type": "config_delta",
  "config_version": 43,
  "changed_node_ends": ["ne_aaa"],
  "removed_node_ends": []
}
```

### Ping Result

```json
{
  "execution_id": "exec_123",
  "node_end_id": "ne_aaa",
  "check_id": "check_ping_123",
  "check_type": "ping",
  "config_version": 42,
  "started_at": "2026-04-17T10:00:00Z",
  "finished_at": "2026-04-17T10:00:20Z",
  "summary": {
    "sent": 20,
    "received": 20,
    "loss_percent": 0,
    "min_ms": 12.4,
    "avg_ms": 15.8,
    "median_ms": 15.1,
    "max_ms": 22.7,
    "stddev_ms": 2.1
  },
  "samples": [
    {"seq": 1, "rtt_ms": 14.2},
    {"seq": 2, "rtt_ms": 15.1}
  ]
}
```

### Traceroute Result

```json
{
  "execution_id": "exec_456",
  "node_end_id": "ne_aaa",
  "check_id": "check_trace_123",
  "check_type": "traceroute",
  "config_version": 42,
  "started_at": "2026-04-17T10:01:00Z",
  "finished_at": "2026-04-17T10:01:08Z",
  "summary": {
    "end_ip": "93.184.216.34",
    "reached": true,
    "hop_count": 12,
    "path_hash": "sha256:example"
  },
  "hops": [
    {
      "hop": 1,
      "ip": "192.0.2.1",
      "hostname": "router.local",
      "rtts_ms": [1.2, 1.4, 1.3],
      "loss_percent": 0
    }
  ]
}
```

## Milestone Checklist

- [ ] Create Go module and service skeletons.
- [ ] Add controller health endpoint.
- [ ] Add node health command.
- [ ] Add protobuf definitions for node gRPC protocol.
- [ ] Add PostgreSQL migrations.
- [ ] Add node registration.
- [ ] Add node heartbeat.
- [ ] Add end CRUD.
- [ ] Add node-to-end connection CRUD.
- [ ] Add per-connection check configuration.
- [ ] Add node config snapshot generation.
- [ ] Add gRPC observe stream for config updates and observed data.
- [ ] Add node config acknowledgement.
- [ ] Add result ingestion.
- [ ] Add ping runner.
- [ ] Add traceroute runner.
- [ ] Add result query APIs.
- [ ] Add basic alerting.
- [ ] Add Docker deployment files.
- [ ] Add minimal documentation for operating controller and node.

## Open Questions

- Should the first gRPC version use only bidirectional `Observe`, or keep unary `GetNodeConfig` for simpler reconnect/resync behavior?
- Should `node_end_checks` allow multiple schedules per node-to-end pair, or should all checks on one connection share a schedule in the first version?
- Should nodes persist the last valid config to disk so measurements continue across node restarts when the controller is unreachable?
- Should public nodes be allowed, or is this initially for trusted private infrastructure only?
- Should the controller support multi-tenant users from the beginning?
- Should ping use raw ICMP only, or should system `ping` fallback be mandatory for easier deployment?
- Should traceroute initially use native Go packets or shell out to system traceroute for faster delivery?
- What database retention policy is required for raw samples?
- Should the project prioritize a SmokePing-compatible import/export path?

## Recommended Next Step

Start with Phase 1 and Phase 2. The most important early milestone is a full controller-node loop:

1. Controller stores a node.
2. Controller stores an end.
3. Controller stores a `node_end` connection.
4. Controller stores one enabled `node_end_check`, for example ping.
5. Controller builds a config snapshot for the node.
6. Node fetches or receives the config snapshot over gRPC.
7. Node applies the config and acknowledges the version.
8. Node schedules a placeholder measurement locally.
9. Node streams a result with the applied `config_version`.
10. Controller stores and exposes the result.

After that loop works, ping, traceroute, DNS, and other checks can be added as real measurement runners without changing the overall architecture.
