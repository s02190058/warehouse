BEGIN;

CREATE TABLE warehouses
(
    id           INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name         TEXT    NOT NULL,
    is_available BOOLEAN NOT NULL
);

CREATE TABLE products
(
    id       INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name     TEXT        NOT NULL,
    size     TEXT,
    code     TEXT UNIQUE NOT NULL,
    quantity INTEGER     NOT NULL CHECK ( quantity >= 0 ) DEFAULT 0
);

CREATE TABLE warehouse_products
(
    warehouse_id      INTEGER NOT NULL REFERENCES warehouses ON DELETE RESTRICT,
    product_id        INTEGER NOT NULL REFERENCES products ON DELETE RESTRICT,
    quantity          INTEGER NOT NULL CHECK ( quantity >= 0 ),
    reserved_quantity INTEGER NOT NULL CHECK ( 0 <= reserved_quantity AND reserved_quantity <= quantity) DEFAULT 0,

    PRIMARY KEY (warehouse_id, product_id)
);

COMMIT;