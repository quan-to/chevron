--changeset racerxdl:create_gpgkeyuid_table
CREATE TABLE chevron_gpg_key_uid
(
    gpg_key_uid_id          uuid      NOT null PRIMARY KEY,
    gpg_key_uid_name        varchar,
    gpg_key_uid_email       varchar,
    gpg_key_uid_description varchar,
    gpg_key_uid_parent      uuid      NOT NULL REFERENCES chevron_gpg_key (gpg_key_id) ON DELETE CASCADE,

    gpg_key_uid_created_at  timestamp NOT NULL DEFAULT now(),
    gpg_key_uid_updated_at  timestamp NOT NULL DEFAULT now(),
    gpg_key_uid_deleted_at  timestamp NULL
);

CREATE INDEX chevron_gpg_key_uid_name_idx ON chevron_gpg_key_uid (gpg_key_uid_name);
CREATE INDEX chevron_gpg_key_uid_email_idx ON chevron_gpg_key_uid (gpg_key_uid_email);
