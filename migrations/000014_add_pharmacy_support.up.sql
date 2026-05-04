-- is_pharmacy column already exists in items table (added previously).
-- Only create item_pharmacies table if not exists.

CREATE TABLE IF NOT EXISTS item_pharmacies (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    item_id BIGINT UNSIGNED NOT NULL UNIQUE,

    -- Pharmacy specific fields
    generic_name VARCHAR(200) NOT NULL,
    brand_name VARCHAR(200),
    strength VARCHAR(100),
    dosage_form VARCHAR(100),

    -- Regulatory & Control Information
    schedule_type VARCHAR(50),
    is_prescription_required BOOLEAN DEFAULT FALSE,
    is_controlled_drug BOOLEAN DEFAULT FALSE,

    -- Storage & Dosage Information
    storage_condition VARCHAR(200),
    max_daily_dose VARCHAR(100),

    -- Shelf Life & Reorder
    shelf_life_days INT,
    reorder_alert_days INT,

    -- Manufacturing & Registration
    manufacturer_name VARCHAR(200),
    drug_registration_no VARCHAR(100) UNIQUE,
    route_of_administration VARCHAR(100),
    therapeutic_class VARCHAR(150),

    -- Timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,

    -- Indexes
    INDEX idx_item_id (item_id),
    INDEX idx_generic_name (generic_name),
    INDEX idx_drug_registration_no (drug_registration_no),
    INDEX idx_therapeutic_class (therapeutic_class),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at),

    -- Foreign Key Constraint
    CONSTRAINT fk_item_pharmacies_item
        FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
