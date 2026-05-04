-- ITeam Profile: Item Categories table
CREATE TABLE IF NOT EXISTS categories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE,
    department_id BIGINT UNSIGNED,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_department_id (department_id),
    CONSTRAINT fk_categories_department
        FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ITeam Profile: Item Sub-Categories table
CREATE TABLE IF NOT EXISTS sub_categories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    category_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(150) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_sub_categories_category
        FOREIGN KEY (category_id)
        REFERENCES categories(id)
        ON DELETE CASCADE,
    INDEX idx_category_id (category_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ITeam Profile: Item Minor Categories table
CREATE TABLE IF NOT EXISTS minor_categories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    sub_category_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(150) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_minor_categories_sub_category
        FOREIGN KEY (sub_category_id)
        REFERENCES sub_categories(id)
        ON DELETE CASCADE,
    INDEX idx_sub_category_id (sub_category_id),
    UNIQUE INDEX unique_sub_category_name (sub_category_id, name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ITeam Profile: Item Types table
CREATE TABLE IF NOT EXISTS item_types (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ITeam Profile: UOM table
CREATE TABLE IF NOT EXISTS uoms (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ITeam Profile: Sales Setups table
CREATE TABLE IF NOT EXISTS sales_setups (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ITeam Profile: Tags table
CREATE TABLE IF NOT EXISTS tags (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ITeam Profile: Sales Supply Types table------
CREATE TABLE IF NOT EXISTS sales_supply_types (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ITeam Profile: Sales Tax Setups table-------
CREATE TABLE IF NOT EXISTS sales_tax_setups (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ITeam Profile: Create Items Table
CREATE TABLE IF NOT EXISTS items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    category_id          BIGINT UNSIGNED NOT NULL,
    sub_category_id      BIGINT UNSIGNED NOT NULL,
    minor_category_id    BIGINT UNSIGNED NOT NULL,
    uom_id               BIGINT UNSIGNED,
    sales_tax_setup_id   BIGINT UNSIGNED,
    sales_supply_type_id BIGINT UNSIGNED,
    supplier_id          BIGINT UNSIGNED,
    tag_id               BIGINT UNSIGNED,
    item_type_id         BIGINT UNSIGNED,
    department_id        BIGINT UNSIGNED,
    file_id              BIGINT UNSIGNED,
    name                 VARCHAR(150) NOT NULL,
    barcode              VARCHAR(400),
    sku                  VARCHAR(100) UNIQUE,
    calculate_base_price  DECIMAL(12, 2) DEFAULT 0,
    cost_price            DECIMAL(12, 2) DEFAULT 0,
    standard_sales_price  DECIMAL(12, 2) DEFAULT 0,
    cost_on_gp            DECIMAL(12, 2) DEFAULT 0,
    last_cost             DECIMAL(12, 2) DEFAULT 0,
    avg_cost              DECIMAL(12, 2) DEFAULT 0,
    price_increasing      DECIMAL(12, 2) DEFAULT 0,
    alt_umo               DECIMAL(10, 2) DEFAULT 0,
    alt_unit              VARCHAR(50),
    max_discount          DECIMAL(5, 2)  DEFAULT 0,
    reorder_minimum_qty   DECIMAL(10, 2) DEFAULT 0,
    is_child_barcode      BOOLEAN DEFAULT FALSE,
    auto_barcode          BOOLEAN DEFAULT FALSE,
    can_be_sold           BOOLEAN DEFAULT TRUE,
    can_be_produced       BOOLEAN DEFAULT FALSE,
    can_be_rented         BOOLEAN DEFAULT FALSE,
    can_be_purchased      BOOLEAN DEFAULT TRUE,
    is_vat_rebatable      BOOLEAN DEFAULT FALSE,
    is_not_allow_decimal  BOOLEAN DEFAULT FALSE,
    is_active             BOOLEAN DEFAULT TRUE,
    is_style              BOOLEAN DEFAULT FALSE,
    is_percentage         BOOLEAN DEFAULT FALSE,
    created_at            DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at            DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at            DATETIME NULL,
    INDEX idx_category_id       (category_id),
    INDEX idx_sub_category_id   (sub_category_id),
    INDEX idx_minor_category_id (minor_category_id),
    INDEX idx_item_type_id      (item_type_id),
    INDEX idx_tag_id            (tag_id),
    INDEX idx_department_id     (department_id),
    INDEX idx_file_id           (file_id),
    INDEX idx_name              (name),
    INDEX idx_barcode           (barcode),
    INDEX idx_is_active         (is_active),
    INDEX idx_created_at        (created_at),
    INDEX idx_deleted_at        (deleted_at), 
    CONSTRAINT fk_items_category
        FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT,
    CONSTRAINT fk_items_sub_category
        FOREIGN KEY (sub_category_id) REFERENCES sub_categories(id) ON DELETE RESTRICT,
    CONSTRAINT fk_items_minor_category
        FOREIGN KEY (minor_category_id) REFERENCES minor_categories(id) ON DELETE RESTRICT,
    CONSTRAINT fk_items_uom
        FOREIGN KEY (uom_id) REFERENCES uoms(id) ON DELETE SET NULL,
    CONSTRAINT fk_items_sales_tax_setup
        FOREIGN KEY (sales_tax_setup_id) REFERENCES sales_tax_setups(id) ON DELETE SET NULL,
    CONSTRAINT fk_items_sales_supply_type
        FOREIGN KEY (sales_supply_type_id) REFERENCES sales_supply_types(id) ON DELETE SET NULL,
    CONSTRAINT fk_items_supplier
        FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE SET NULL,
    CONSTRAINT fk_items_tag
        FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE SET NULL,
    CONSTRAINT fk_items_item_type
        FOREIGN KEY (item_type_id) REFERENCES item_types(id) ON DELETE SET NULL,
    CONSTRAINT fk_items_department
        FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL,
    CONSTRAINT fk_items_file
        FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;