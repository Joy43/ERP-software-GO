-- =============================================
-- TABLE: locations
-- =============================================
CREATE TABLE IF NOT EXISTS locations (
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    code       VARCHAR(50)  NOT NULL UNIQUE,
    name       VARCHAR(200) NOT NULL,
    type       ENUM('warehouse','store','showroom','outlet') NOT NULL,

    office_id  BIGINT UNSIGNED NOT NULL,
    parent_id  BIGINT UNSIGNED NULL,
    manager_id BIGINT UNSIGNED NULL,

    location   TEXT NULL,
    is_active  BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    INDEX idx_locations_office_id  (office_id),
    INDEX idx_locations_parent_id  (parent_id),
    INDEX idx_locations_manager_id (manager_id),
    INDEX idx_locations_deleted_at (deleted_at),

    CONSTRAINT fk_locations_office
        FOREIGN KEY (office_id)  REFERENCES offices(id)    ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_locations_parent
        FOREIGN KEY (parent_id)  REFERENCES locations(id)  ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_locations_manager
        FOREIGN KEY (manager_id) REFERENCES users(id)      ON UPDATE CASCADE ON DELETE SET NULL

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- =============================================
-- TABLE: inventory_types
-- =============================================
CREATE TABLE IF NOT EXISTS inventory_types (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    type_code   VARCHAR(150) NOT NULL UNIQUE,
    type_name   VARCHAR(100) NOT NULL,
    description TEXT NULL,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- =============================================
-- TABLE: projects
-- =============================================
CREATE TABLE IF NOT EXISTS projects (
    id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_name VARCHAR(255) NOT NULL,
    description  LONGTEXT     NULL,
    project_code VARCHAR(50)  NOT NULL UNIQUE,

    start_date DATE           NULL,
    end_date   DATE           NULL,
    budget     DECIMAL(15,2)  NULL,

    status    ENUM('PLANNING','ACTIVE','ON_HOLD','COMPLETED','CANCELLED') NOT NULL DEFAULT 'PLANNING',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    -- NULL-able FKs: SET NULL is valid
    manager_id BIGINT UNSIGNED NULL,
    office_id  BIGINT UNSIGNED NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    INDEX idx_projects_project_code (project_code),
    INDEX idx_projects_status       (status),
    INDEX idx_projects_office_id    (office_id),
    INDEX idx_projects_manager_id   (manager_id),
    INDEX idx_projects_is_active    (is_active),
    INDEX idx_projects_deleted_at   (deleted_at),

    CONSTRAINT fk_projects_manager
        FOREIGN KEY (manager_id) REFERENCES users(id)   ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_projects_office
        FOREIGN KEY (office_id)  REFERENCES offices(id) ON UPDATE CASCADE ON DELETE CASCADE

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- =============================================
-- TABLE: location_stocks
-- =============================================
CREATE TABLE IF NOT EXISTS location_stocks (
    id                BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    quantity          DECIMAL(15,3) NOT NULL DEFAULT 0,
    reserved_quantity DECIMAL(15,3) NOT NULL DEFAULT 0,
    last_cost         DECIMAL(15,4) NOT NULL DEFAULT 0,
    average_cost      DECIMAL(15,4) NOT NULL DEFAULT 0,

    item_id     BIGINT UNSIGNED NOT NULL,
    location_id BIGINT UNSIGNED NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    UNIQUE KEY idx_item_location (item_id, location_id),
    INDEX idx_location_stocks_item_id     (item_id),
    INDEX idx_location_stocks_location_id (location_id),
    INDEX idx_location_stocks_deleted_at  (deleted_at),

    CONSTRAINT fk_location_stocks_item
        FOREIGN KEY (item_id)     REFERENCES items(id)     ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_location_stocks_location
        FOREIGN KEY (location_id) REFERENCES locations(id) ON UPDATE CASCADE ON DELETE CASCADE,

    CONSTRAINT chk_reserved_le_quantity
        CHECK (reserved_quantity <= quantity)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- =============================================
-- TABLE: requisitions
-- =============================================
CREATE TABLE IF NOT EXISTS requisitions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    requisition_number VARCHAR(50)  NOT NULL UNIQUE,
    requisition_type   ENUM('EMPLOYEE','DEPARTMENT','PROJECT') NOT NULL DEFAULT 'EMPLOYEE',
    status             ENUM('DRAFT','PENDING','DEPARTMENT_APPROVED','FINANCE_APPROVED','APPROVED','REJECTED','CANCELLED','ORDERED') NOT NULL DEFAULT 'DRAFT',
    rejection_reason   LONGTEXT NULL,

    expected_date DATE     NOT NULL,
    created_date  DATE     NOT NULL DEFAULT (CURDATE()),
    remarks       LONGTEXT NULL,
    description   LONGTEXT NULL,

    -- NULL-able party FKs → SET NULL is valid
    employee_id   BIGINT UNSIGNED NULL,
    department_id BIGINT UNSIGNED NULL,
    project_id    BIGINT UNSIGNED NULL,
    buyer_id      BIGINT UNSIGNED NULL,

    -- NULL-able location FKs → SET NULL is valid
    office_id         BIGINT UNSIGNED NULL,
    location_id       BIGINT UNSIGNED NULL,
    inventory_type_id BIGINT UNSIGNED NULL,
    supplier_id       BIGINT UNSIGNED NULL,

    -- NULL-able audit FKs → SET NULL is valid
    created_by_id BIGINT UNSIGNED NULL,
    updated_by_id BIGINT UNSIGNED NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    INDEX idx_req_requisition_type  (requisition_type),
    INDEX idx_req_status            (status),
    INDEX idx_req_employee_id       (employee_id),
    INDEX idx_req_department_id     (department_id),
    INDEX idx_req_project_id        (project_id),
    INDEX idx_req_buyer_id          (buyer_id),
    INDEX idx_req_office_id         (office_id),
    INDEX idx_req_location_id       (location_id),
    INDEX idx_req_inventory_type_id (inventory_type_id),
    INDEX idx_req_supplier_id       (supplier_id),
    INDEX idx_req_created_by_id     (created_by_id),
    INDEX idx_req_deleted_at        (deleted_at),

    CONSTRAINT fk_req_employee    FOREIGN KEY (employee_id)       REFERENCES users(id)           ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_req_department  FOREIGN KEY (department_id)     REFERENCES departments(id)     ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_req_project     FOREIGN KEY (project_id)        REFERENCES projects(id)        ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_req_buyer       FOREIGN KEY (buyer_id)          REFERENCES users(id)           ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_req_office      FOREIGN KEY (office_id)         REFERENCES offices(id)         ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_req_location    FOREIGN KEY (location_id)       REFERENCES locations(id)       ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_req_inv_type    FOREIGN KEY (inventory_type_id) REFERENCES inventory_types(id) ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_req_supplier    FOREIGN KEY (supplier_id)       REFERENCES suppliers(id)       ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_req_created_by  FOREIGN KEY (created_by_id)     REFERENCES users(id)           ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_req_updated_by  FOREIGN KEY (updated_by_id)     REFERENCES users(id)           ON UPDATE CASCADE ON DELETE SET NULL

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- =============================================
-- TABLE: requisition_items
-- =============================================
CREATE TABLE IF NOT EXISTS requisition_items (
    id             BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    requisition_id BIGINT UNSIGNED NOT NULL,
    item_id        BIGINT UNSIGNED NOT NULL,

    request_quantity  DECIMAL(15,2) NOT NULL,
    approved_quantity DECIMAL(15,2) NULL,
    current_stock     DECIMAL(15,2) NULL,

    last_cost    DECIMAL(15,2) NULL,
    average_cost DECIMAL(15,2) NULL,
    description  LONGTEXT      NULL,

    -- NULL-able classification FKs → SET NULL is valid
    item_type_id      BIGINT UNSIGNED NULL,
    category_id       BIGINT UNSIGNED NULL,
    sub_category_id   BIGINT UNSIGNED NULL,
    minor_category_id BIGINT UNSIGNED NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_ri_requisition_id (requisition_id),
    INDEX idx_ri_item_id        (item_id),

    CONSTRAINT fk_ri_requisition    FOREIGN KEY (requisition_id)    REFERENCES requisitions(id)      ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_ri_item           FOREIGN KEY (item_id)           REFERENCES items(id)             ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_ri_item_type      FOREIGN KEY (item_type_id)      REFERENCES item_types(id)        ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_ri_category       FOREIGN KEY (category_id)       REFERENCES categories(id)        ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_ri_sub_category   FOREIGN KEY (sub_category_id)   REFERENCES sub_categories(id)    ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_ri_minor_category FOREIGN KEY (minor_category_id) REFERENCES minor_categories(id)  ON UPDATE CASCADE ON DELETE SET NULL

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- =============================================
-- TABLE: requisition_status_history
-- =============================================
CREATE TABLE IF NOT EXISTS requisition_status_history (
    id             BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    requisition_id BIGINT UNSIGNED NOT NULL,

    -- NULL-able FK → SET NULL is valid
    user_id     BIGINT UNSIGNED NULL,
    from_status VARCHAR(50) NOT NULL DEFAULT '',
    to_status   VARCHAR(50) NOT NULL,
    remarks     LONGTEXT    NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_rsh_requisition_id (requisition_id),
    INDEX idx_rsh_user_id        (user_id),

    CONSTRAINT fk_rsh_requisition FOREIGN KEY (requisition_id) REFERENCES requisitions(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_rsh_user        FOREIGN KEY (user_id)        REFERENCES users(id)        ON UPDATE CASCADE ON DELETE SET NULL

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- =============================================
-- TABLE: payment_by_grns
-- FIX: grn_id, office_id, supplier_id, office_head_id, payment_mode_id
--      were NOT NULL but had ON DELETE SET NULL → changed to NULL
-- =============================================
CREATE TABLE IF NOT EXISTS payment_by_grns (
    payment_by_grn_id BIGINT UNSIGNED   NOT NULL AUTO_INCREMENT PRIMARY KEY,
    payment_date      TIMESTAMP       NOT NULL,
    money_receipt_no  VARCHAR(100)    NOT NULL UNIQUE,

    -- NULL-able FKs → SET NULL is valid
    grn_id          BIGINT UNSIGNED NULL,   -- references goods_receipt_notes(id)
    office_id       BIGINT UNSIGNED NULL,
    supplier_id     BIGINT UNSIGNED NULL,
    office_head_id  BIGINT UNSIGNED NULL,
    payment_mode_id BIGINT UNSIGNED NULL,

    payable_amount    DECIMAL(18,4) NOT NULL DEFAULT 0,
    paying_amount     DECIMAL(18,4) NOT NULL DEFAULT 0,
    adjustment_amount DECIMAL(18,4) NOT NULL DEFAULT 0,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_pbg_grn_id          (grn_id),
    INDEX idx_pbg_office_id       (office_id),
    INDEX idx_pbg_supplier_id     (supplier_id),
    INDEX idx_pbg_office_head_id  (office_head_id),
    INDEX idx_pbg_payment_mode_id (payment_mode_id),

    CONSTRAINT fk_pbg_grn
        FOREIGN KEY (grn_id)          REFERENCES goods_receipt_notes(id) ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_pbg_office
        FOREIGN KEY (office_id)       REFERENCES offices(id)         ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_pbg_supplier
        FOREIGN KEY (supplier_id)     REFERENCES suppliers(id)       ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_pbg_office_head
        FOREIGN KEY (office_head_id)  REFERENCES users(id)           ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_pbg_payment_mode
        FOREIGN KEY (payment_mode_id) REFERENCES payment_modes(id)   ON UPDATE CASCADE ON DELETE SET NULL

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- =============================================
-- TABLE: advance_payments
-- FIX: office_id, account_head_id, supplier_head_id, po_id, payment_mode_id
--      were NOT NULL but had ON DELETE SET NULL → changed to NULL
-- =============================================
CREATE TABLE IF NOT EXISTS advance_payments (
    advance_payment_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    payment_date       TIMESTAMP     NOT NULL,
    narration          TEXT          NULL,
    lc_no              VARCHAR(100)  NULL,

    -- NULL-able FKs → SET NULL is valid
    office_id        BIGINT UNSIGNED NULL,
    account_head_id  BIGINT UNSIGNED NULL,
    supplier_head_id BIGINT UNSIGNED NULL,
    po_id            BIGINT UNSIGNED NULL,
    payment_mode_id  BIGINT UNSIGNED NULL,

    cash_amount DECIMAL(18,4) NOT NULL DEFAULT 0,
    amount      DECIMAL(18,4) NOT NULL DEFAULT 0,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_ap_lc_no            (lc_no),
    INDEX idx_ap_office_id        (office_id),
    INDEX idx_ap_account_head_id  (account_head_id),
    INDEX idx_ap_supplier_head_id (supplier_head_id),
    INDEX idx_ap_po_id            (po_id),
    INDEX idx_ap_payment_mode_id  (payment_mode_id),

    CONSTRAINT fk_ap_office
        FOREIGN KEY (office_id)        REFERENCES offices(id)         ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_ap_account_head
        FOREIGN KEY (account_head_id)  REFERENCES users(id)           ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_ap_supplier_head
        FOREIGN KEY (supplier_head_id) REFERENCES users(id)           ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_ap_purchase_order
        FOREIGN KEY (po_id)            REFERENCES purchase_orders(id) ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_ap_payment_mode
        FOREIGN KEY (payment_mode_id)  REFERENCES payment_modes(id)   ON UPDATE CASCADE ON DELETE SET NULL

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- =============================================
-- TABLE: supplier_bills
-- FIX: office_id, supplier_id were NOT NULL but had ON DELETE RESTRICT — these
--      are fine as RESTRICT. file_id was NULL already — also fine.
--      No changes needed here, included for completeness.
-- =============================================
CREATE TABLE IF NOT EXISTS supplier_bills (
    supplier_bill_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    create_date      TIMESTAMP     NOT NULL,
    tent_pay_date    TIMESTAMP     NULL,
    bill_no          VARCHAR(100)  NOT NULL UNIQUE,
    vat_challan_no   VARCHAR(100)  NULL,
    remarks          TEXT          NULL,

    bill_amount DECIMAL(18,4) NOT NULL DEFAULT 0,
    discount    DECIMAL(18,4) NOT NULL DEFAULT 0,
    advance     DECIMAL(18,4) NOT NULL DEFAULT 0,
    net_pay     DECIMAL(18,4) NOT NULL DEFAULT 0,
    vat         DECIMAL(18,4) NOT NULL DEFAULT 0,
    sd          DECIMAL(18,4) NOT NULL DEFAULT 0,

    -- NOT NULL FKs with RESTRICT → valid
    office_id   BIGINT UNSIGNED NOT NULL,
    supplier_id BIGINT UNSIGNED NOT NULL,
    -- NULL FK with SET NULL → valid
    file_id     BIGINT UNSIGNED NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_sb_office_id   (office_id),
    INDEX idx_sb_supplier_id (supplier_id),
    INDEX idx_sb_file_id     (file_id),

    CONSTRAINT fk_sb_office
        FOREIGN KEY (office_id)   REFERENCES offices(id)   ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_sb_supplier
        FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_sb_file
        FOREIGN KEY (file_id)     REFERENCES files(id)     ON UPDATE CASCADE ON DELETE SET NULL

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;