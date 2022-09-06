
-- +migrate Up
CREATE TABLE `users` (
    `id` VARCHAR(36) NOT NULL DEFAULT (UUID()),
    `email` VARCHAR(255) COLLATE utf8mb4_unicode_ci NOT NULL UNIQUE,
    `password` TEXT COLLATE utf8mb4_unicode_ci NOT NULL,
    `username` VARCHAR(255) COLLATE utf8mb4_unicode_ci NOT NULL UNIQUE,
    `admin` BOOLEAN DEFAULT FALSE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP,
    PRIMARY KEY (id, email, username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down
DROP TABLE `users`;
