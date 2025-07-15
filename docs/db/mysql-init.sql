DROP DATABASE IF EXISTS fresh_db;
CREATE DATABASE fresh_db
    DEFAULT CHARACTER SET = 'utf8mb4';

USE fresh_db;
DROP TABLE IF EXISTS buyers;
CREATE TABLE buyers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    card_number_id VARCHAR(10) NOT NULL UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL
);

DROP TABLE IF EXISTS product_types;
CREATE TABLE product_types (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

DROP TABLE IF EXISTS localities;
CREATE TABLE localities (
    id             VARCHAR(255) PRIMARY KEY,
    locality_name  VARCHAR(255) NOT NULL,
    province_name  VARCHAR(255) NOT NULL,
    country_name   VARCHAR(255) NOT NULL
);

DROP TABLE IF EXISTS sellers;
CREATE TABLE sellers (
    id              INT AUTO_INCREMENT PRIMARY KEY,
    cid             INT NOT NULL UNIQUE,
    company_name    VARCHAR(255) NOT NULL,
    address         VARCHAR(255) NOT NULL,
    telephone       VARCHAR(255) NOT NULL,
    locality_id     VARCHAR(255) NOT NULL,
    FOREIGN KEY (locality_id) REFERENCES localities(id)
);

DROP TABLE IF EXISTS products;
CREATE TABLE products (
    id INT PRIMARY KEY AUTO_INCREMENT,
    product_code VARCHAR(255) NOT NULL,
    description TEXT,
    width DECIMAL(10, 2),
    height DECIMAL(10, 2),
    length DECIMAL(10, 2),
    net_weight DECIMAL(10, 2),
    expiration_rate INT,
    recommended_freezing_temperature DECIMAL(10, 2),
    freezing_rate INT,
    product_type_id INT,
    seller_id INT,
    FOREIGN KEY (product_type_id) REFERENCES product_types(id),
    FOREIGN KEY (seller_id) REFERENCES sellers(id),
    UNIQUE (product_code)
);

DROP TABLE IF EXISTS carries;
CREATE TABLE carries (
    id INT AUTO_INCREMENT PRIMARY KEY,
    cid VARCHAR(255) NOT NULL UNIQUE,
    company_name VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    telephone VARCHAR(255) NOT NULL,
    locality_id VARCHAR(255) NOT NULL,
    FOREIGN KEY (locality_id) REFERENCES localities(id)
);

DROP TABLE IF EXISTS warehouses;
CREATE TABLE warehouses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    warehouse_code VARCHAR(255) NOT NULL UNIQUE,
    address VARCHAR(255) NOT NULL,
    telephone VARCHAR(255) NOT NULL,
    minimum_capacity INT NOT NULL,
    minimum_temperature FLOAT
);

DROP TABLE IF EXISTS employees;
CREATE TABLE employees (
    id INT PRIMARY KEY AUTO_INCREMENT,
    card_number_id VARCHAR(10) NOT NULL UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    warehouse_id INT,
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
); 

DROP TABLE IF EXISTS sections;
CREATE TABLE `sections`
(
    `id`                  int          NOT NULL AUTO_INCREMENT,
    `section_number`      varchar(255) NOT NULL UNIQUE,
    `current_temperature` float        NOT NULL,
    `minimum_temperature` float        NOT NULL,
    `current_capacity`    int          NOT NULL,
    `minimum_capacity`    int          NOT NULL,
    `maximum_capacity`    int          NOT NULL,
    `warehouse_id`        int          NOT NULL,
    `product_type_id`     int          NOT NULL,

    PRIMARY KEY (`id`),
    KEY                   `idx_sections_warehouse_id` (`warehouse_id`),
    KEY                   `idx_sections_product_type_id` (`product_type_id`),
    CONSTRAINT `fk_sections_warehouse_id` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouses` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_sections_product_type_id` FOREIGN KEY (`product_type_id`) REFERENCES `product_types` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

DROP TABLE IF EXISTS product_batches;
CREATE TABLE `product_batches`
(
    `id`                  int          NOT NULL AUTO_INCREMENT,
    `batch_number`        int          NOT NULL UNIQUE,
    `current_quantity`    int          NOT NULL,
    `current_temperature` int          NOT NULL,
    `due_date`            varchar(255) NOT NULL,
    `initial_quantity`    int          NOT NULL,
    `manufacturing_date`  varchar(255) NOT NULL,
    `manufacturing_hour`  int          NOT NULL,
    `minumum_temperature` int          NOT NULL,
    `product_id`          int          NOT NULL,
    `section_id`          int          NOT NULL,

    PRIMARY KEY (`id`),
    KEY                   `idx_product_batches_product_id` (`product_id`),
    KEY                   `idx_product_batches_section_id` (`section_id`),
    CONSTRAINT `fk_product_batches_product_id` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_product_batches_section_id` FOREIGN KEY (`section_id`) REFERENCES `sections` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);
DROP TABLE IF EXISTS inbound_orders;
CREATE TABLE inbound_orders (
    id INT PRIMARY KEY AUTO_INCREMENT,
    order_date DATE NOT NULL,
    order_number VARCHAR(255) NOT NULL UNIQUE,
    employee_id INT,
    product_batch_id INT,
    warehouse_id INT,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    FOREIGN KEY (product_batch_id) REFERENCES product_batches(id) ON DELETE CASCADE,
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE
);

DROP TABLE IF EXISTS purchase_orders;
CREATE TABLE purchase_orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_number VARCHAR(255) NOT NULL UNIQUE,
    order_date DATETIME NOT NULL,
    tracking_code VARCHAR(255) NOT NULL,
    buyer_id INT NOT NULL,
    FOREIGN KEY (buyer_id) REFERENCES buyers(id)
);

DROP TABLE IF EXISTS product_records;
CREATE TABLE product_records (
    id INT AUTO_INCREMENT PRIMARY KEY,
    last_update_date DATE NOT NULL,
    purchase_price DECIMAL(10, 2) NOT NULL,
    sale_price DECIMAL(10, 2) NOT NULL,
    product_id INT NOT NULL,    
    FOREIGN KEY (product_id) REFERENCES products(id)
);

DROP TABLE IF EXISTS purchase_order_details;
CREATE TABLE purchase_order_details (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT NOT NULL,
    product_record_id INT NOT NULL,
    quantity INT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES purchase_orders(id),
    FOREIGN KEY (product_record_id) REFERENCES product_records(id)
);
