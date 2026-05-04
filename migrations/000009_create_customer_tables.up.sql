-- Customers Table
CREATE TABLE IF NOT EXISTS customers (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    office_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    billing_name VARCHAR(255),
    mobile_no VARCHAR(20) NOT NULL,
    email VARCHAR(190),
    credit_days INT DEFAULT 0,
    marital_status VARCHAR(50),
    marriage_date DATE,
    birthday_date DATE,
    gender VARCHAR(20),
    nid VARCHAR(50),
    partner_group_id BIGINT UNSIGNED,
    partner_sub_group_id BIGINT UNSIGNED,
    credit_limit DECIMAL(15,2) DEFAULT 0.00,
    default_sales_rep_id BIGINT UNSIGNED,
    tolerance DECIMAL(5,2) DEFAULT 0.00,
    bin_number VARCHAR(100),
    price_type VARCHAR(100),
    tcs_percentage DECIMAL(5,2) DEFAULT 0.00,
    vat_reg_no_central VARCHAR(100),
    tin_number VARCHAR(100),
    trade_license_number VARCHAR(100),
    is_foreign_customer BOOLEAN DEFAULT FALSE,
    is_distributor BOOLEAN DEFAULT FALSE,
    is_employee_customer BOOLEAN DEFAULT FALSE,
    user_id BIGINT UNSIGNED,
    billing_contact_person VARCHAR(150),
    billing_contact_no VARCHAR(20),
    district_id BIGINT UNSIGNED,
    thana_id BIGINT UNSIGNED,
    billing_address TEXT,
    point_expiry_days INT DEFAULT 0,
    tax_bracket_id BIGINT UNSIGNED,
    remarks TEXT,
    attachment VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    CONSTRAINT fk_customer_office FOREIGN KEY (office_id) REFERENCES offices(id),
    CONSTRAINT fk_customer_group FOREIGN KEY (partner_group_id) REFERENCES partner_groups(id),
    CONSTRAINT fk_customer_sub_group FOREIGN KEY (partner_sub_group_id) REFERENCES partner_sub_groups(id),
    CONSTRAINT fk_customer_sales_rep FOREIGN KEY (default_sales_rep_id) REFERENCES users(id),
    CONSTRAINT fk_customer_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_customer_district FOREIGN KEY (district_id) REFERENCES districts(id),
    CONSTRAINT fk_customer_thana FOREIGN KEY (thana_id) REFERENCES thanas(id),
    CONSTRAINT fk_customer_tax_bracket FOREIGN KEY (tax_bracket_id) REFERENCES tax_brackets(id)
);

-- Customer Shipping Addresses Table
CREATE TABLE IF NOT EXISTS customer_shipping_addresses (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    customer_id BIGINT UNSIGNED NOT NULL,
    code VARCHAR(50),
    ref_outlet VARCHAR(150),
    address TEXT NOT NULL,
    gps VARCHAR(255),
    contact_person VARCHAR(150),
    contact_no VARCHAR(20),
    vat_reg_no VARCHAR(100),
    district_id BIGINT UNSIGNED,
    thana_id BIGINT UNSIGNED,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_shipping_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    CONSTRAINT fk_shipping_district FOREIGN KEY (district_id) REFERENCES districts(id),
    CONSTRAINT fk_shipping_thana FOREIGN KEY (thana_id) REFERENCES thanas(id)
);

-- Customer Bank Information Table
CREATE TABLE IF NOT EXISTS customer_bank_infos (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    customer_id BIGINT UNSIGNED NOT NULL,
    bank_name VARCHAR(255),
    account_no VARCHAR(100),
    account_name VARCHAR(255),
    branch_name VARCHAR(255),
    routing_no VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_bank_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);
