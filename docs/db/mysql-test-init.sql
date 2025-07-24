DROP
DATABASE IF EXISTS fresh_db_test;
CREATE
DATABASE fresh_db_test
    DEFAULT CHARACTER SET = 'utf8mb4';

USE
fresh_db_test;
DROP TABLE IF EXISTS buyers;
CREATE TABLE buyers
(
    id             INT AUTO_INCREMENT PRIMARY KEY,
    card_number_id VARCHAR(10)  NOT NULL UNIQUE,
    first_name     VARCHAR(255) NOT NULL,
    last_name      VARCHAR(255) NOT NULL
);

DROP TABLE IF EXISTS product_types;
CREATE TABLE product_types
(
    id          INT PRIMARY KEY AUTO_INCREMENT,
    name        VARCHAR(255) NOT NULL,
    description TEXT
);

DROP TABLE IF EXISTS localities;
CREATE TABLE localities
(
    id            VARCHAR(255) PRIMARY KEY,
    locality_name VARCHAR(255) NOT NULL,
    province_name VARCHAR(255) NOT NULL,
    country_name  VARCHAR(255) NOT NULL
);

DROP TABLE IF EXISTS sellers;
CREATE TABLE sellers
(
    id           INT AUTO_INCREMENT PRIMARY KEY,
    cid          INT          NOT NULL UNIQUE,
    company_name VARCHAR(255) NOT NULL,
    address      VARCHAR(255) NOT NULL,
    telephone    VARCHAR(255) NOT NULL,
    locality_id  VARCHAR(255) NOT NULL,
    FOREIGN KEY (locality_id) REFERENCES localities (id)
);

DROP TABLE IF EXISTS products;
CREATE TABLE products
(
    id                               INT PRIMARY KEY AUTO_INCREMENT,
    product_code                     VARCHAR(255) NOT NULL,
    description                      TEXT,
    width                            DECIMAL(10, 2),
    height                           DECIMAL(10, 2),
    length                           DECIMAL(10, 2),
    net_weight                       DECIMAL(10, 2),
    expiration_rate                  INT,
    recommended_freezing_temperature DECIMAL(10, 2),
    freezing_rate                    INT,
    product_type_id                  INT,
    seller_id                        INT,
    FOREIGN KEY (product_type_id) REFERENCES product_types (id),
    FOREIGN KEY (seller_id) REFERENCES sellers (id),
    UNIQUE (product_code)
);

DROP TABLE IF EXISTS carries;
CREATE TABLE carries
(
    id           INT AUTO_INCREMENT PRIMARY KEY,
    cid          VARCHAR(255) NOT NULL UNIQUE,
    company_name VARCHAR(255) NOT NULL,
    address      VARCHAR(255) NOT NULL,
    telephone    VARCHAR(255) NOT NULL,
    locality_id  VARCHAR(255) NOT NULL,
    FOREIGN KEY (locality_id) REFERENCES localities (id)
);

DROP TABLE IF EXISTS warehouses;
CREATE TABLE warehouses
(
    id                  INT AUTO_INCREMENT PRIMARY KEY,
    warehouse_code      VARCHAR(255) NOT NULL UNIQUE,
    address             VARCHAR(255) NOT NULL,
    telephone           VARCHAR(255) NOT NULL,
    minimum_capacity    INT          NOT NULL,
    minimum_temperature FLOAT
);

DROP TABLE IF EXISTS employees;
CREATE TABLE employees
(
    id             INT PRIMARY KEY AUTO_INCREMENT,
    card_number_id VARCHAR(10)  NOT NULL UNIQUE,
    first_name     VARCHAR(255) NOT NULL,
    last_name      VARCHAR(255) NOT NULL,
    warehouse_id   INT,
    FOREIGN KEY (warehouse_id) REFERENCES warehouses (id)
);

DROP TABLE IF EXISTS sections;
CREATE TABLE `sections`
(
    `id`                  INT   NOT NULL AUTO_INCREMENT,
    `section_number`      INT   NOT NULL UNIQUE,
    `current_temperature` FLOAT NOT NULL,
    `minimum_temperature` FLOAT NOT NULL,
    `current_capacity`    INT   NOT NULL,
    `minimum_capacity`    INT   NOT NULL,
    `maximum_capacity`    INT   NOT NULL,
    `warehouse_id`        INT   NOT NULL,
    `product_type_id`     INT   NOT NULL,

    PRIMARY KEY (`id`),
    KEY                   `idx_sections_warehouse_id` (`warehouse_id`),
    KEY                   `idx_sections_product_type_id` (`product_type_id`),
    CONSTRAINT `fk_sections_warehouse_id` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouses` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_sections_product_type_id` FOREIGN KEY (`product_type_id`) REFERENCES `product_types` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

DROP TABLE IF EXISTS product_batches;
CREATE TABLE `product_batches`
(
    `id`                  INT  NOT NULL AUTO_INCREMENT,
    `batch_number`        INT  NOT NULL UNIQUE,
    `current_quantity`    INT  NOT NULL,
    `current_temperature` INT  NOT NULL,
    `due_date`            DATE NOT NULL,
    `initial_quantity`    INT  NOT NULL,
    `manufacturing_date`  DATE NOT NULL,
    `manufacturing_hour`  INT  NOT NULL,
    `minimum_temperature` INT  NOT NULL,
    `product_id`          INT  NOT NULL,
    `section_id`          INT  NOT NULL,

    PRIMARY KEY (`id`),
    KEY                   `idx_product_batches_product_id` (`product_id`),
    KEY                   `idx_product_batches_section_id` (`section_id`),
    CONSTRAINT `fk_product_batches_product_id` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_product_batches_section_id` FOREIGN KEY (`section_id`) REFERENCES `sections` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);
DROP TABLE IF EXISTS inbound_orders;
CREATE TABLE inbound_orders
(
    id               INT PRIMARY KEY AUTO_INCREMENT,
    order_date       DATE         NOT NULL,
    order_number     VARCHAR(255) NOT NULL UNIQUE,
    employee_id      INT,
    product_batch_id INT,
    warehouse_id     INT,
    FOREIGN KEY (employee_id) REFERENCES employees (id) ON DELETE CASCADE,
    FOREIGN KEY (product_batch_id) REFERENCES product_batches (id) ON DELETE CASCADE,
    FOREIGN KEY (warehouse_id) REFERENCES warehouses (id) ON DELETE CASCADE
);

DROP TABLE IF EXISTS purchase_orders;
CREATE TABLE purchase_orders
(
    id            INT AUTO_INCREMENT PRIMARY KEY,
    order_number  VARCHAR(255) NOT NULL UNIQUE,
    order_date    DATETIME     NOT NULL,
    tracking_code VARCHAR(255) NOT NULL,
    buyer_id      INT          NOT NULL,
    FOREIGN KEY (buyer_id) REFERENCES buyers (id)
);

DROP TABLE IF EXISTS product_records;
CREATE TABLE product_records
(
    id               INT AUTO_INCREMENT PRIMARY KEY,
    last_update_date DATE           NOT NULL,
    purchase_price   DECIMAL(10, 2) NOT NULL,
    sale_price       DECIMAL(10, 2) NOT NULL,
    product_id       INT            NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products (id)
);

DROP TABLE IF EXISTS purchase_order_details;
CREATE TABLE purchase_order_details
(
    id                INT AUTO_INCREMENT PRIMARY KEY,
    order_id          INT NOT NULL,
    product_record_id INT NOT NULL,
    quantity          INT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES purchase_orders (id),
    FOREIGN KEY (product_record_id) REFERENCES product_records (id)
);

INSERT INTO buyers (card_number_id, first_name, last_name)
VALUES
    ('1234567890', 'Juan', 'Pérez');
-- Inserciones para la tabla product_types
INSERT INTO product_types (name, description)
VALUES
    ('Frutas', 'Productos frescos de origen frutal');
-- Inserciones para la tabla localities
INSERT INTO localities (id, locality_name, province_name, country_name)
VALUES
    ('LOC001', 'Buenos Aires', 'Buenos Aires', 'Argentina');
-- Inserciones para la tabla sellers
INSERT INTO sellers (cid, company_name, address, telephone, locality_id)
VALUES
    (1001, 'Frutas del Sur', 'Calle Falsa 123', '123456789', 'LOC001');
-- Inserciones para la tabla products
INSERT INTO products (product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id)
VALUES
    ('FRU001', 'Manzanas Red Delicious', 10.5, 10.0, 15.0, 1.2, 7, -1.0, 0, 1, 1);
-- Inserciones para la tabla carries
INSERT INTO carries (cid, company_name, address, telephone, locality_id)
VALUES
    ('CARR001', 'Logística Rápida', 'Camino al Aeropuerto 1234', '123123123', 'LOC001');
-- Inserciones para la tabla warehouses
INSERT INTO warehouses (warehouse_code, address, telephone, minimum_capacity, minimum_temperature)
VALUES
    ('WH001', 'Parque Industrial 1', '555-5551', 1000, -10.0);
-- Inserciones para la tabla employees
INSERT INTO employees (card_number_id, first_name, last_name, warehouse_id)
VALUES
    ('EMP001', 'Raul', 'García', 1);
-- Inserciones para la tabla sections
INSERT INTO sections (section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, product_type_id)
VALUES
    (1, -10.0, -20.0, 300, 100, 500, 1, 1);
-- Inserciones para la tabla product_batches
INSERT INTO product_batches (batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minimum_temperature, product_id, section_id)
VALUES
    (1001, 500, -18, '2023-12-31', 1000, '2023-10-10', 10, -20, 1, 1);
-- Inserciones para la tabla inbound_orders
INSERT INTO inbound_orders (order_date, order_number, employee_id, product_batch_id, warehouse_id)
VALUES
    ('2023-10-15', 'ORD001', 1, 1, 1);
-- Inserciones para la tabla purchase_orders
INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id)
VALUES
    ('PO001', '2023-10-14 10:00:00', 'TRK001', 1);
-- Inserciones para la tabla product_records
INSERT INTO product_records (last_update_date, purchase_price, sale_price, product_id)
VALUES
    ('2023-10-01', 10.0, 15.0, 1);
-- Inserciones para la tabla purchase_order_details
INSERT INTO purchase_order_details (order_id, product_record_id, quantity)
VALUES
    (1, 1, 100);