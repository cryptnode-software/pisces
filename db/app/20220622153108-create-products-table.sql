
-- +migrate Up
CREATE TABLE `products` (
    `id` VARCHAR(36) NOT NULL DEFAULT (UUID()),
    `description` TEXT,
    `inventory` INT,
    `name` TEXT,
    `cost` DECIMAL(13,2),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down

DROP TABLE `products`;
