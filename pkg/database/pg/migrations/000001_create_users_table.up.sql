--changeset racerxdl:create_users_table
CREATE TABLE chevron_user (
    user_id uuid NOT null PRIMARY KEY,
    user_fingerprint varchar NOT NULL,
    user_password bytea NOT NULL,
    user_full_name varchar NOT NULL,
    user_created_at timestamp NOT NULL DEFAULT now(),
    user_updated_at timestamp NOT NULL DEFAULT now(),
    user_deleted_at timestamp NULL
);

CREATE INDEX chevron_user_fingerprint_idx ON chevron_user (user_fingerprint);
