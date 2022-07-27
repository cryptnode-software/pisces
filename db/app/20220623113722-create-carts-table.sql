
-- +migrate Up
CREATE TABLE `carts`(
    `id` BIGINT(10) NOT NULL AUTO_INCREMENT,
    `quantity` BIGINT NOT NULL,
    `product_id` BIGINT,
    `order_id` BIGINT,
    INDEX (product_id, order_id),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (product_id) REFERENCES products (id),
    FOREIGN KEY (order_id) REFERENCES orders (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down
DROP TABLE `carts`;
