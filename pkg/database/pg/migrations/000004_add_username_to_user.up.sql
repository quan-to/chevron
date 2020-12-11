--changeset racerxdl:add_username_to_user

ALTER TABLE chevron_user
    ADD COLUMN user_username VARCHAR NOT NULL DEFAULT '';

CREATE INDEX chevron_user_username_idx ON chevron_user (user_username);
