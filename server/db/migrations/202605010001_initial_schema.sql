-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION netstamp_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email CITEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER set_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION netstamp_set_updated_at();

CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug CITEXT NOT NULL UNIQUE,
    created_by_user_id UUID NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE TRIGGER set_teams_updated_at
    BEFORE UPDATE ON teams
    FOR EACH ROW
    EXECUTE FUNCTION netstamp_set_updated_at();

CREATE TABLE team_members (
    team_id UUID NOT NULL REFERENCES teams(id),
    user_id UUID NOT NULL REFERENCES users(id),
    role TEXT NOT NULL CHECK (role IN ('owner', 'admin', 'member')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,
    PRIMARY KEY (team_id, user_id)
);

CREATE TRIGGER set_team_members_updated_at
    BEFORE UPDATE ON team_members
    FOR EACH ROW
    EXECUTE FUNCTION netstamp_set_updated_at();

CREATE INDEX idx_team_members_user_id ON team_members (user_id);

CREATE TABLE probes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL REFERENCES teams(id),
    name TEXT NOT NULL,
    provider TEXT NULL,
    region TEXT NULL,
    hostname TEXT NULL,
    labels JSONB NOT NULL DEFAULT '{}'::jsonb,
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_probes_team_id_id UNIQUE (team_id, id)
);

CREATE TRIGGER set_probes_updated_at
    BEFORE UPDATE ON probes
    FOR EACH ROW
    EXECUTE FUNCTION netstamp_set_updated_at();

CREATE INDEX idx_probes_team_id ON probes (team_id);

CREATE TABLE probe_credentials (
    probe_id UUID PRIMARY KEY REFERENCES probes(id),
    secret_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_rotated_at TIMESTAMPTZ NULL
);

CREATE TABLE probe_status (
    probe_id UUID PRIMARY KEY REFERENCES probes(id),
    status TEXT NOT NULL,
    last_seen_at TIMESTAMPTZ NULL,
    agent_version TEXT NULL,
    ip_families TEXT[] NOT NULL DEFAULT '{}'::text[],
    capabilities JSONB NOT NULL DEFAULT '{}'::jsonb,
    health JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE TABLE checks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL REFERENCES teams(id),
    name TEXT NOT NULL,
    check_type TEXT NOT NULL CHECK (check_type IN ('ping', 'traceroute', 'dns')),
    target TEXT NOT NULL,
    target_type TEXT NOT NULL CHECK (target_type IN ('host', 'ip', 'dns_query', 'dns_resolver')),
    description TEXT NULL,
    parameters JSONB NOT NULL DEFAULT '{}'::jsonb,
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_checks_team_id_id UNIQUE (team_id, id)
);

CREATE TRIGGER set_checks_updated_at
    BEFORE UPDATE ON checks
    FOR EACH ROW
    EXECUTE FUNCTION netstamp_set_updated_at();

CREATE INDEX idx_checks_team_id ON checks (team_id);
CREATE INDEX idx_checks_team_id_check_type ON checks (team_id, check_type);

CREATE TABLE probe_checks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL REFERENCES teams(id),
    probe_id UUID NOT NULL,
    check_id UUID NOT NULL,
    name TEXT NULL,
    enabled BOOLEAN NOT NULL DEFAULT true,
    interval_seconds INTEGER NOT NULL CHECK (interval_seconds > 0),
    jitter_seconds INTEGER NOT NULL DEFAULT 0 CHECK (jitter_seconds >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_probe_checks_team_id_id UNIQUE (team_id, id),
    CONSTRAINT fk_probe_checks_team_probe
        FOREIGN KEY (team_id, probe_id) REFERENCES probes(team_id, id),
    CONSTRAINT fk_probe_checks_team_check
        FOREIGN KEY (team_id, check_id) REFERENCES checks(team_id, id)
);

CREATE TRIGGER set_probe_checks_updated_at
    BEFORE UPDATE ON probe_checks
    FOR EACH ROW
    EXECUTE FUNCTION netstamp_set_updated_at();

CREATE UNIQUE INDEX idx_probe_checks_active_team_probe_check
    ON probe_checks (team_id, probe_id, check_id)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_probe_checks_team_id ON probe_checks (team_id);
CREATE INDEX idx_probe_checks_probe_id ON probe_checks (probe_id);
CREATE INDEX idx_probe_checks_check_id ON probe_checks (check_id);

CREATE TABLE ping_results (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL,
    probe_check_id UUID NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ NOT NULL,
    duration_ms INTEGER NOT NULL CHECK (duration_ms >= 0),
    status TEXT NOT NULL CHECK (status IN ('success', 'partial', 'timeout', 'error')),
    sent_count INTEGER NOT NULL CHECK (sent_count >= 0),
    received_count INTEGER NOT NULL CHECK (received_count >= 0),
    loss_percent DOUBLE PRECISION NOT NULL CHECK (loss_percent >= 0 AND loss_percent <= 100),
    rtt_min_ms DOUBLE PRECISION NULL CHECK (rtt_min_ms >= 0),
    rtt_avg_ms DOUBLE PRECISION NULL CHECK (rtt_avg_ms >= 0),
    rtt_median_ms DOUBLE PRECISION NULL CHECK (rtt_median_ms >= 0),
    rtt_max_ms DOUBLE PRECISION NULL CHECK (rtt_max_ms >= 0),
    rtt_stddev_ms DOUBLE PRECISION NULL CHECK (rtt_stddev_ms >= 0),
    rtt_samples_ms DOUBLE PRECISION[] NOT NULL DEFAULT '{}'::double precision[],
    resolved_ip INET NULL,
    ip_version TEXT NULL CHECK (ip_version IN ('ipv4', 'ipv6', 'auto')),
    raw JSONB NOT NULL DEFAULT '{}'::jsonb,
    error_code TEXT NULL,
    error_message TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (id, started_at),
    CHECK (finished_at >= started_at),
    CHECK (received_count <= sent_count),
    CONSTRAINT fk_ping_results_team_probe_check
        FOREIGN KEY (team_id, probe_check_id) REFERENCES probe_checks(team_id, id)
);

SELECT create_hypertable('ping_results', 'started_at', if_not_exists => TRUE);

CREATE INDEX idx_ping_results_team_id_started_at ON ping_results (team_id, started_at DESC);
CREATE INDEX idx_ping_results_probe_check_id_started_at ON ping_results (probe_check_id, started_at DESC);
CREATE INDEX idx_ping_results_status_started_at ON ping_results (status, started_at DESC);

CREATE TABLE dns_results (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL,
    probe_check_id UUID NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ NOT NULL,
    duration_ms INTEGER NOT NULL CHECK (duration_ms >= 0),
    status TEXT NOT NULL CHECK (status IN ('success', 'partial', 'timeout', 'error')),
    query_name TEXT NOT NULL,
    record_type TEXT NOT NULL,
    resolver TEXT NOT NULL,
    transport TEXT NOT NULL CHECK (transport IN ('udp', 'tcp')),
    rcode TEXT NULL,
    success BOOLEAN NOT NULL,
    answer_count INTEGER NOT NULL DEFAULT 0 CHECK (answer_count >= 0),
    response_time_ms INTEGER NOT NULL CHECK (response_time_ms >= 0),
    answers JSONB NOT NULL DEFAULT '[]'::jsonb,
    raw JSONB NOT NULL DEFAULT '{}'::jsonb,
    error_code TEXT NULL,
    error_message TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (id, started_at),
    CHECK (finished_at >= started_at),
    CONSTRAINT fk_dns_results_team_probe_check
        FOREIGN KEY (team_id, probe_check_id) REFERENCES probe_checks(team_id, id)
);

SELECT create_hypertable('dns_results', 'started_at', if_not_exists => TRUE);

CREATE INDEX idx_dns_results_team_id_started_at ON dns_results (team_id, started_at DESC);
CREATE INDEX idx_dns_results_probe_check_id_started_at ON dns_results (probe_check_id, started_at DESC);
CREATE INDEX idx_dns_results_status_started_at ON dns_results (status, started_at DESC);
CREATE INDEX idx_dns_results_team_query_started_at ON dns_results (team_id, query_name, started_at DESC);

CREATE TABLE traceroute_results (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL,
    probe_check_id UUID NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ NOT NULL,
    duration_ms INTEGER NOT NULL CHECK (duration_ms >= 0),
    status TEXT NOT NULL CHECK (status IN ('success', 'partial', 'timeout', 'error')),
    resolved_ip INET NULL,
    reached BOOLEAN NOT NULL DEFAULT false,
    hop_count INTEGER NOT NULL DEFAULT 0 CHECK (hop_count >= 0),
    path_hash TEXT NULL,
    protocol TEXT NOT NULL CHECK (protocol IN ('icmp')),
    raw JSONB NOT NULL DEFAULT '{}'::jsonb,
    error_code TEXT NULL,
    error_message TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (id, started_at),
    CHECK (finished_at >= started_at),
    CONSTRAINT fk_traceroute_results_team_probe_check
        FOREIGN KEY (team_id, probe_check_id) REFERENCES probe_checks(team_id, id)
);

SELECT create_hypertable('traceroute_results', 'started_at', if_not_exists => TRUE);

CREATE INDEX idx_traceroute_results_team_id_started_at ON traceroute_results (team_id, started_at DESC);
CREATE INDEX idx_traceroute_results_probe_check_id_started_at ON traceroute_results (probe_check_id, started_at DESC);
CREATE INDEX idx_traceroute_results_status_started_at ON traceroute_results (status, started_at DESC);
CREATE INDEX idx_traceroute_results_team_path_hash_started_at ON traceroute_results (team_id, path_hash, started_at DESC);

CREATE TABLE traceroute_hops (
    traceroute_result_id UUID NOT NULL,
    team_id UUID NOT NULL,
    probe_check_id UUID NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    hop_number INTEGER NOT NULL CHECK (hop_number > 0),
    hop_ip INET NULL,
    hostname TEXT NULL,
    rtts_ms DOUBLE PRECISION[] NOT NULL DEFAULT '{}'::double precision[],
    loss_percent DOUBLE PRECISION NOT NULL CHECK (loss_percent >= 0 AND loss_percent <= 100),
    error_code TEXT NULL,
    error_message TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (traceroute_result_id, started_at, hop_number),
    CONSTRAINT fk_traceroute_hops_team_probe_check
        FOREIGN KEY (team_id, probe_check_id) REFERENCES probe_checks(team_id, id)
);

SELECT create_hypertable('traceroute_hops', 'started_at', if_not_exists => TRUE);

CREATE INDEX idx_traceroute_hops_team_id_started_at ON traceroute_hops (team_id, started_at DESC);
CREATE INDEX idx_traceroute_hops_probe_check_id_started_at ON traceroute_hops (probe_check_id, started_at DESC);
CREATE INDEX idx_traceroute_hops_result_id_started_at ON traceroute_hops (
    traceroute_result_id,
    started_at DESC
);
CREATE INDEX idx_traceroute_hops_team_hop_ip_started_at ON traceroute_hops (team_id, hop_ip, started_at DESC);
CREATE INDEX idx_traceroute_hops_team_probe_check_hop_started_at ON traceroute_hops (
    team_id,
    probe_check_id,
    hop_number,
    started_at DESC
);

-- +goose Down
DROP TABLE IF EXISTS traceroute_hops;
DROP TABLE IF EXISTS traceroute_results;
DROP TABLE IF EXISTS dns_results;
DROP TABLE IF EXISTS ping_results;
DROP TABLE IF EXISTS probe_checks;
DROP TABLE IF EXISTS checks;
DROP TABLE IF EXISTS probe_status;
DROP TABLE IF EXISTS probe_credentials;
DROP TABLE IF EXISTS probes;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS users;
DROP FUNCTION IF EXISTS netstamp_set_updated_at();
