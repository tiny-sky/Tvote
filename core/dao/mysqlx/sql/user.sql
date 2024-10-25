CREATE TABLE `User` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `votes` INT DEFAULT 0,
    PRIMARY KEY (`id`)
);
