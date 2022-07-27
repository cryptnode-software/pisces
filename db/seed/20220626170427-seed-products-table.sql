
-- +migrate Up

-- INSERT INTO `products`(`id`, `description`, `name`, `cost`)
-- VALUES('seeded-budget-1', 1, 1, 1, 'monthly', 1000, 1000, 'completes', '[]', NOW(), NOW() + INTERVAL 1 MONTH, 1, NOW());

INSERT INTO `products` (`id`, `description`, `name`, `cost`)
VALUES 
    (1, 'Test One Product Description', 'Test One Product', 0.00);

INSERT INTO `products` (`id`, `description`, `name`, `cost`)
VALUES
    (2, 'Test Two Product Description', 'Test Two Product', 0.00);

INSERT INTO `products` (`id`, `description`, `name`, `cost`)
VALUES
    (3, 'Test Three Product Description', 'Test Three Product', 0.00);

-- +migrate Down
DELETE FROM products WHERE id = 0;
DELETE FROM products WHERE id = 1;
DELETE FROM products WHERE id = 2;

