--changeset racerxdl:create_gpgkey_table
CREATE TABLE chevron_gpg_key
(
    gpg_key_id               uuid      NOT NULL PRIMARY KEY,
    gpg_key_full_fingerprint varchar   NOT NULL,
    gpg_key_fingerprint16    varchar   NOT NULL,
    gpg_key_keybits          int       NOT NULL,
    gpg_key_parent           uuid      REFERENCES chevron_gpg_key (gpg_key_id) ON DELETE CASCADE,

    gpg_key_public_key       text,
    gpg_key_private_key      text,

    gpg_key_created_at       timestamp NOT NULL DEFAULT now(),
    gpg_key_updated_at       timestamp NOT NULL DEFAULT now(),
    gpg_key_deleted_at       timestamp NULL
);

CREATE INDEX chevron_gpg_key_full_fingerprint_idx ON chevron_gpg_key (gpg_key_full_fingerprint);
CREATE INDEX chevron_gpg_key_fingerprint16_idx ON chevron_gpg_key (gpg_key_fingerprint16);
