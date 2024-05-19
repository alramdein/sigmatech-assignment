CREATE TABLE customer (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nik VARCHAR(20) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    legal_name VARCHAR(100) NOT NULL,
    birth_place VARCHAR(100),
    birth_date DATE,
    salary DECIMAL(15, 2),
    ktp_photo BLOB,
    selfie_photo BLOB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE `limit` (
    id INT AUTO_INCREMENT PRIMARY KEY,
    customer_id INT NOT NULL,
    tenor_1 DECIMAL(15, 2),
    tenor_2 DECIMAL(15, 2),
    tenor_3 DECIMAL(15, 2),
    tenor_4 DECIMAL(15, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customer(id)
);

CREATE TABLE transaction (
    id INT AUTO_INCREMENT PRIMARY KEY,
    customer_id INT NOT NULL,
    contract_number VARCHAR(100) NOT NULL,
    otr DECIMAL(15, 2),
    admin_fee DECIMAL(15, 2),
    installment_amount DECIMAL(15, 2),
    interest_amount DECIMAL(15, 2),
    asset_name VARCHAR(100),
    tenor INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customer(id)
);
