
-- +migrate Up
CREATE TABLE `inquires` (
  `id` BIGINT(10) NOT NULL AUTO_INCREMENT,
  `description` TEXT COLLATE utf8mb4_unicode_ci NOT NULL,
  `first_name` VARCHAR (255) COLLATE utf8mb4_unicode_ci NULL,
  `last_name` VARCHAR (255) COLLATE utf8mb4_unicode_ci NULL,
  `email` VARCHAR (255) COLLATE utf8mb4_unicode_ci NULL,
  `number` VARCHAR (255) COLLATE utf8mb4_unicode_ci NULL,
  `closed_at` TIMESTAMP DEFAULT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `inquiry_attachments` (
    `id` BIGINT(10) NOT NULL AUTO_INCREMENT,
    `value` VARCHAR(255) COLLATE utf8mb4_unicode_ci,
    `inquiry_id` BIGINT,
    INDEX inq_id(inquiry_id),
    PRIMARY KEY (id),
    FOREIGN KEY (inquiry_id)
        REFERENCES inquires (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down

DROP TABLE `inquiry_attachments`;
DROP TABLE `inquires`;

