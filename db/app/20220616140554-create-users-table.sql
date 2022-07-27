
-- +migrate Up
CREATE TABLE `users` (
    `id` BIGINT(10) NOT NULL AUTO_INCREMENT,
    `email` VARCHAR(255) COLLATE utf8mb4_unicode_ci NOT NULL UNIQUE,
    `password` VARCHAR(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `username` VARCHAR(255) COLLATE utf8mb4_unicode_ci NOT NULL UNIQUE,
    `admin` BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (id, email, username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down

DROP TABLE `users`;
