-- +goose Up
ALTER TABLE users
    ADD COLUMN display_name text,
    ADD CONSTRAINT users_display_name_not_empty CHECK (
        display_name IS NULL OR (
            length(btrim(display_name)) > 0 AND
            length(btrim(display_name)) <= 100
        )
    );

ALTER TABLE probe_statuses
    ADD COLUMN public_ip inet,
    ADD COLUMN asn bigint,
    ADD CONSTRAINT probe_statuses_asn_range CHECK (
        asn IS NULL OR (asn > 0 AND asn <= 4294967295)
    );

-- +goose Down
ALTER TABLE probe_statuses
    DROP CONSTRAINT probe_statuses_asn_range,
    DROP COLUMN asn,
    DROP COLUMN public_ip;

ALTER TABLE users
    DROP CONSTRAINT users_display_name_not_empty,
    DROP COLUMN display_name;
