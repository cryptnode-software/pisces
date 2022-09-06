
-- +migrate Up
CREATE TABLE `inquiries` (
  `id` VARCHAR(36) NOT NULL DEFAULT (UUID()),
  `description` TEXT COLLATE utf8mb4_unicode_ci NOT NULL,
  `first_name` TEXT COLLATE utf8mb4_unicode_ci NULL,
  `last_name` TEXT COLLATE utf8mb4_unicode_ci NULL,
  `email` TEXT COLLATE utf8mb4_unicode_ci NULL,
  `number` TEXT COLLATE utf8mb4_unicode_ci NULL,
  `closed_at` TIMESTAMP DEFAULT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down
DROP TABLE `inquiries`;

