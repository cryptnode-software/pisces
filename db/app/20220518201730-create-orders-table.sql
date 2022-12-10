
-- +migrate Up
CREATE TABLE `orders` (
  `id` VARCHAR(36) NOT NULL DEFAULT (UUID()),
  `ext_id` VARCHAR(255) COLLATE utf8mb4_unicode_ci NULL, -- confirmation code
  `due` TIMESTAMP NOT NULL,
  `inquiry_id` VARCHAR(36) NOT NULL,
  INDEX inq_id(inquiry_id),
  `payment_method` VARCHAR(255) COLLATE utf8mb4_unicode_ci,
  `status` VARCHAR(255) COLLATE utf8mb4_unicode_ci,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (inquiry_id)
    REFERENCES inquiries (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down
DROP TABLE `orders`;
