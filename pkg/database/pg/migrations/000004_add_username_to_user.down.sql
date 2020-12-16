--changeset racerxdl:add_username_to_user

ALTER TABLE chevron_user
    DROP COLUMN user_username;

DROP INDEX chevron_user_username_idx;
