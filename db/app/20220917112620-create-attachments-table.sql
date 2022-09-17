
-- +migrate Up
CREATE TABLE `attachments`(
    `id` VARCHAR(36) NOT NULL DEFAULT (UUID()),
    `type` VARCHAR(36) NOT NULL,
    `url` TEXT NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `inquiry_attachments` (
    `attachment_id` VARCHAR(36) NOT NULL,
    `inquiry_id` VARCHAR(36) NOT NULL,
    FOREIGN KEY (attachment_id) REFERENCES attachments (id),
    FOREIGN KEY (inquiry_id) REFERENCES inquiries (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down
DROP TABLE `inquiry_attachments`;
DROP TABLE `attachments`;
