
-- +migrate Up
CREATE TABLE `orders` (
  `id` BIGINT(10) NOT NULL AUTO_INCREMENT,
  `ext_id` VARCHAR(255) COLLATE utf8mb4_unicode_ci NULL, -- confirmation code
  `due` TIMESTAMP NOT NULL,
  `inquiry_id` BIGINT,
  INDEX inq_id(inquiry_id),
  `payment_method` VARCHAR(255) COLLATE utf8mb4_unicode_ci,
  `status` VARCHAR(255) COLLATE utf8mb4_unicode_ci,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (inquiry_id)
    REFERENCES inquires (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE `inquires` 
ADD COLUMN `order_id` BIGINT NULL,
ADD FOREIGN KEY fk_order_id (order_id)
  REFERENCES orders (id);

-- +migrate Down
ALTER TABLE `inquires`
DROP FOREIGN KEY `inquires_ibfk_1`,
DROP COLUMN `order_id`;
DROP TABLE `orders`;
