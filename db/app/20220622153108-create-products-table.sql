
-- +migrate Up
CREATE TABLE `products` (
    `id` BIGINT(10) NOT NULL AUTO_INCREMENT,
    `description` TEXT,
    `name` VARCHAR(255),
    `cost` DECIMAL(13,2),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down

DROP TABLE `products`;
