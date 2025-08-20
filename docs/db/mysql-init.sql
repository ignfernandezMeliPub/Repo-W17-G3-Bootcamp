
DROP DATABASE IF EXISTS fresh_db;
CREATE DATABASE fresh_db
    DEFAULT CHARACTER SET = 'utf8mb4';

USE fresh_db;

DROP TABLE IF EXISTS logs;

CREATE TABLE logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    source VARCHAR(255) NOT NULL,
    endpoint VARCHAR(255) NOT NULL,
    http_method VARCHAR(255) NOT NULL,
    layer VARCHAR(255) NOT NULL,
    action VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    level VARCHAR(255) NOT NULL,
    time DATETIME NOT NULL
);

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
    `id`                  INT          NOT NULL AUTO_INCREMENT,
    `section_number`      INT          NOT NULL UNIQUE,
    `current_temperature` FLOAT        NOT NULL,
    `minimum_temperature` FLOAT        NOT NULL,
    `current_capacity`    INT          NOT NULL,
    `minimum_capacity`    INT          NOT NULL,
    `maximum_capacity`    INT          NOT NULL,
    `warehouse_id`        INT          NOT NULL,
    `product_type_id`     INT          NOT NULL,

    PRIMARY KEY (`id`),
    KEY `idx_sections_warehouse_id` (`warehouse_id`),
    KEY `idx_sections_product_type_id` (`product_type_id`),
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
    KEY `idx_product_batches_product_id` (`product_id`),
    KEY `idx_product_batches_section_id` (`section_id`),
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


INSERT INTO buyers (card_number_id, first_name, last_name)
VALUES
    ('1234567890', 'Juan', 'Pérez'),
    ('0987654321', 'María', 'González'),
    ('1122334455', 'Pedro', 'Martínez'),
    ('2233445566', 'Laura', 'Fernández'),
    ('3344556677', 'Ana', 'López');
-- Inserciones para la tabla product_types
INSERT INTO product_types (name, description)
VALUES
    ('Frutas', 'Productos frescos de origen frutal'),
    ('Verduras', 'Productos frescos de origen vegetal'),
    ('Carnes', 'Productos frescos de origen animal'),
    ('Pescados', 'Productos frescos del mar'),
    ('Lácteos', 'Productos derivados de la leche');
-- Inserciones para la tabla localities
INSERT INTO localities (id, locality_name, province_name, country_name)
VALUES
    ('LOC001', 'Buenos Aires', 'Buenos Aires', 'Argentina'),
    ('LOC002', 'Córdoba', 'Córdoba', 'Argentina'),
    ('LOC003', 'Rosario', 'Santa Fe', 'Argentina'),
    ('LOC004', 'Mendoza', 'Mendoza', 'Argentina'),
    ('LOC005', 'La Plata', 'Buenos Aires', 'Argentina');
-- Inserciones para la tabla sellers
INSERT INTO sellers (cid, company_name, address, telephone, locality_id)
VALUES
    (1001, 'Frutas del Sur', 'Calle Falsa 123', '123456789', 'LOC001'),
    (1002, 'Verduras Frescas', 'Avenida Siempreviva 742', '987654321', 'LOC002'),
    (1003, 'Carnes Argentinas', 'Boulevard de los Sueños Rotos 45', '456123789', 'LOC003'),
    (1004, 'Pescados del Atlántico', 'Ruta 40 km 2020', '789456123', 'LOC004'),
    (1005, 'Lácteos La Campiña', 'Callejón sin salida 37', '321789456', 'LOC005');
-- Inserciones para la tabla products
INSERT INTO products (product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id)
VALUES
    ('FRU001', 'Manzanas Red Delicious', 10.5, 10.0, 15.0, 1.2, 7, -1.0, 0, 1, 1),
    ('VER002', 'Zanahorias Frescas', 5.0, 15.0, 5.0, 0.8, 5, 0.0, 0, 2, 2),
    ('CAR003', 'Bife de Chorizo', 20.0, 2.5, 10.0, 0.5, 14, -18.0, 3, 3, 3),
    ('PES004', 'Merluza Congelada', 30.0, 2.0, 8.0, 0.9, 21, -20.0, 7, 4, 4),
    ('LAC005', 'Queso Parmesano', 10.0, 5.0, 5.0, 0.3, 30, 4.0, 0, 5, 5);
-- Inserciones para la tabla carries
INSERT INTO carries (cid, company_name, address, telephone, locality_id)
VALUES
    ('CARR001', 'Logística Rápida', 'Camino al Aeropuerto 1234', '123123123', 'LOC001'),
    ('CARR002', 'Transportes del Norte', 'Ruta Nacional 9 km 200', '321321321', 'LOC002'),
    ('CARR003', 'Carga Segura', 'Avenida de los Ganaderos 789', '456456456', 'LOC003'),
    ('CARR004', 'Movilidad Expresa', 'Boulevard Marítimo 101', '654654654', 'LOC004'),
    ('CARR005', 'Envía Ya', 'Calle Comercio 148', '789789789', 'LOC005');
-- Inserciones para la tabla warehouses
INSERT INTO warehouses (warehouse_code, address, telephone, minimum_capacity, minimum_temperature)
VALUES
    ('WH001', 'Parque Industrial 1', '555-5551', 1000, -10.0),
    ('WH002', 'Centro Logístico 5', '555-5552', 1500, -5.0),
    ('WH003', 'Depósito Central', '555-5553', 2000, 0.0),
    ('WH004', 'Zona Franca Norte', '555-5554', 2500, -18.0),
    ('WH005', 'Ubicación Este', '555-5555', 1200, 4.0);
-- Inserciones para la tabla employees
INSERT INTO employees (card_number_id, first_name, last_name, warehouse_id)
VALUES
    ('EMP001', 'Raul', 'García', 1),
    ('EMP002', 'Sandra', 'Rojas', 2),
    ('EMP003', 'Luis', 'Molina', 3),
    ('EMP004', 'Marta', 'Díaz', 4),
    ('EMP005', 'Jorge', 'San Martín', 5);
-- Inserciones para la tabla sections
INSERT INTO sections (section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, product_type_id)
VALUES
    (1, -10.0, -20.0, 300, 100, 500, 1, 4),
    (2, -5.0, -15.0, 500, 200, 750, 2, 3),
    (3, 0.0, -2.0, 400, 100, 600, 3, 2),
    (4, 5.0, 2.0, 200, 50, 400, 4, 5),
    (5, 4.0, 1.0, 250, 75, 450, 5, 1);
-- Inserciones para la tabla product_batches
INSERT INTO product_batches (batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minimum_temperature, product_id, section_id)
VALUES
    (1001, 500, -18, '2023-12-31', 1000, '2023-10-10', 10, -20, 3, 1),
    (1002, 300, 3, '2023-10-20', 500, '2023-09-01', 9, 2, 1, 2),
    (1003, 400, -2, '2023-11-15', 600, '2023-09-15', 10, -5, 2, 3),
    (1004, 150, 2, '2023-09-30', 250, '2023-07-20', 11, 0, 5, 4),
    (1005, 600, -20, '2024-01-01', 750, '2023-10-02', 8, -22, 4, 5);
-- Inserciones para la tabla inbound_orders
INSERT INTO inbound_orders (order_date, order_number, employee_id, product_batch_id, warehouse_id)
VALUES
    ('2023-10-15', 'ORD001', 1, 1, 1),
    ('2023-10-10', 'ORD002', 2, 2, 2),
    ('2023-10-12', 'ORD003', 3, 3, 3),
    ('2023-10-11', 'ORD004', 4, 4, 4),
    ('2023-10-13', 'ORD005', 5, 5, 5);
-- Inserciones para la tabla purchase_orders
INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id)
VALUES
    ('PO001', '2023-10-14 10:00:00', 'TRK001', 1),
    ('PO002', '2023-10-15 11:00:00', 'TRK002', 2),
    ('PO003', '2023-10-16 12:00:00', 'TRK003', 3),
    ('PO004', '2023-10-17 13:00:00', 'TRK004', 4),
    ('PO005', '2023-10-18 14:00:00', 'TRK005', 5);
-- Inserciones para la tabla product_records
INSERT INTO product_records (last_update_date, purchase_price, sale_price, product_id)
VALUES
    ('2023-10-01', 10.0, 15.0, 1),
    ('2023-10-01', 5.0, 8.0, 2),
    ('2023-10-01', 12.0, 20.0, 3),
    ('2023-10-01', 9.0, 14.0, 4),
    ('2023-10-01', 7.0, 11.0, 5);
-- Inserciones para la tabla purchase_order_details
INSERT INTO purchase_order_details (order_id, product_record_id, quantity)
VALUES
    (1, 1, 100),
    (1, 2, 100),
    (2, 2, 150),
    (3, 3, 200),
    (4, 4, 250),
    (5, 5, 300);