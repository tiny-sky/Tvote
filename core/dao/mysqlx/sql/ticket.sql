CREATE TABLE `Ticket` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `ticket` VARCHAR(255) NOT NULL UNIQUE,
    `created_at` BIGINT NOT NULL DEFAULT (UNIX_TIMESTAMP()),
    `expires_at` INT NOT NULL,
    `max_usage` INT NOT NULL,
    `used_count` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`)
);
