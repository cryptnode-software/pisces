
-- +migrate Up
INSERT INTO `products` (`id`, `description`, `name`, `cost`)
VALUES 
    ("5915fa01-034a-4b60-8598-3dee4f4e4869", 'Test One Product Description', 'Test One Product', 0.00);

INSERT INTO `products` (`id`, `description`, `name`, `cost`)
VALUES
    ("bf8a8d16-5233-46c3-b3f6-bc7e098760cd", 'Test Two Product Description', 'Test Two Product', 0.00);

INSERT INTO `products` (`id`, `description`, `name`, `cost`)
VALUES
    ("0462870a-e0e8-4921-ba13-e1319fd56c0b", 'Test Three Product Description', 'Test Three Product', 0.00);

-- +migrate Down
DELETE FROM products WHERE id = "5915fa01-034a-4b60-8598-3dee4f4e4869";
DELETE FROM products WHERE id = "bf8a8d16-5233-46c3-b3f6-bc7e098760cd";
DELETE FROM products WHERE id = "0462870a-e0e8-4921-ba13-e1319fd56c0b";

