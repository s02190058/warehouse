BEGIN;

INSERT INTO warehouses (name, is_available)
VALUES ('warehouse_A', true),
       ('warehouse_B', false),
       ('warehouse_C', true);

INSERT INTO products (name, size, code)
VALUES ('product_A', 'A', 'AAA'),
       ('product_B', 'B', 'BBB');

INSERT INTO warehouse_products (warehouse_id, product_id, quantity)
VALUES (1, 1, 4),
       (1, 2, 3),
       (2, 2, 5),
       (3, 1, 10);

COMMIT;