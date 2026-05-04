-- Partner Groups Table
CREATE TABLE IF NOT EXISTS partner_groups (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Partner Sub Groups Table
CREATE TABLE IF NOT EXISTS partner_sub_groups (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    partner_group_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(150) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_partner_group FOREIGN KEY (partner_group_id) REFERENCES partner_groups(id) ON DELETE CASCADE
);

-- Districts Table
CREATE TABLE IF NOT EXISTS districts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Thanas Table
CREATE TABLE IF NOT EXISTS thanas (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    district_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(150) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_district FOREIGN KEY (district_id) REFERENCES districts(id) ON DELETE CASCADE
);

-- Tax Brackets Table
CREATE TABLE IF NOT EXISTS tax_brackets (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE,
    percentage DECIMAL(5,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Suppliers Table
CREATE TABLE IF NOT EXISTS suppliers (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    office_id BIGINT UNSIGNED NOT NULL,
    code VARCHAR(50) UNIQUE,
    name VARCHAR(255) NOT NULL,
    vendor_name VARCHAR(255),
    mobile_no VARCHAR(20) NOT NULL,
    email VARCHAR(190),
    credit_days INT DEFAULT 0,
    partner_group_id BIGINT UNSIGNED,
    partner_sub_group_id BIGINT UNSIGNED,
    tolerance DECIMAL(5,2) DEFAULT 0.00,
    bin_number VARCHAR(100),
    tcs_percentage DECIMAL(5,2) DEFAULT 0.00,
    vat_reg_no_central VARCHAR(100),
    tin_number VARCHAR(100),
    trade_license_number VARCHAR(100),
    is_foreign_supplier BOOLEAN DEFAULT FALSE,
    vds_applicable BOOLEAN DEFAULT FALSE,
    is_anonymous BOOLEAN DEFAULT FALSE,
    is_vat_accounting_not_applicable BOOLEAN DEFAULT FALSE,
    billing_contact_person VARCHAR(150),
    billing_contact_no VARCHAR(20),
    district_id BIGINT UNSIGNED,
    thana_id BIGINT UNSIGNED,
    billing_address TEXT,
    tax_bracket_id BIGINT UNSIGNED,
    remarks TEXT,
    attachment VARCHAR(255),
    address TEXT,
    gps VARCHAR(255),
    contact_person VARCHAR(150),
    vat_reg_no VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    CONSTRAINT fk_supplier_office FOREIGN KEY (office_id) REFERENCES offices(id),
    CONSTRAINT fk_supplier_group FOREIGN KEY (partner_group_id) REFERENCES partner_groups(id),
    CONSTRAINT fk_supplier_sub_group FOREIGN KEY (partner_sub_group_id) REFERENCES partner_sub_groups(id),
    CONSTRAINT fk_supplier_district FOREIGN KEY (district_id) REFERENCES districts(id),
    CONSTRAINT fk_supplier_thana FOREIGN KEY (thana_id) REFERENCES thanas(id),
    CONSTRAINT fk_supplier_tax_bracket FOREIGN KEY (tax_bracket_id) REFERENCES tax_brackets(id)
);
