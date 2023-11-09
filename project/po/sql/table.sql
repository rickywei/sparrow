CREATE TABLE
    IF NOT EXISTS mydb.user(
        `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'id',
        `user_name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'user_name',
        `nick_name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'nick_name',
        `email` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'email',
        `password` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'password',
        `mobile` VARCHAR(16) NOT NULL DEFAULT '' COMMENT 'mobile',
        `created_at` BIGINT NOT NULL COMMENT 'create time',
        `updated_at` BIGINT NOT NULL COMMENT 'update time',
        `deleted_at` BIGINT NOT NULL COMMENT 'delete time',
        PRIMARY KEY(`id`),
        UNIQUE KEY `uk_user_name` (user_name(64)),
        UNIQUE KEY `uk_email` (email(64))
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

-- CREATE TABLE
--     IF NOT EXISTS mydb.todo(
--         `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'id',
--         `user_id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'user_id',
--         `name` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'todo task name',
--         `ddl` BIGINT NOT NULL COMMENT 'deadline',
--         `is_finish` TINYINT(1) NOT NULL COMMENT 'is todo finished',
--         `finish_at` BIGINT NOT NULL COMMENT 'finish time',
--         `created_at` BIGINT NOT NULL COMMENT 'create time',
--         `updated_at` BIGINT NOT NULL COMMENT 'update time',
--         `deleted_at` BIGINT NOT NULL COMMENT 'delete time',
--         PRIMARY KEY(`id`),
--         FOREIGN KEY (user_id) REFERENCES demo.user (id)
--     ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;