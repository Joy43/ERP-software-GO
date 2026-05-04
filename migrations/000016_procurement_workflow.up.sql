-- =============================================
-- TABLE: purchase_orders
-- =============================================
CREATE TABLE IF NOT EXISTS purchase_orders (
    id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    po_number     VARCHAR(50)  NOT NULL UNIQUE,
    po_date       DATE         NOT NULL,
    delivery_date DATE         NULL,

    requisition_id BIGINT UNSIGNED NULL,
    order_type     ENUM('DIRECT','REQUISITION_BASED','CONTRACT') NOT NULL DEFAULT 'DIRECT',

    office_id   BIGINT UNSIGNED NOT NULL,
    location_id BIGINT UNSIGNED NOT NULL,
    supplier_id BIGINT UNSIGNED NOT NULL,

    payment_terms    VARCHAR(100) NULL,
    general_remarks  TEXT         NULL,
    shipping_address TEXT         NULL,

    subtotal        DECIMAL(15,2) NOT NULL DEFAULT 0,
    vat_amount      DECIMAL(15,2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_amount    DECIMAL(15,2) NOT NULL DEFAULT 0,

    status ENUM('DRAFT','ISSUED','CONFIRMED','PARTIALLY_RECEIVED','FULLY_RECEIVED','CANCELLED') NOT NULL DEFAULT 'DRAFT',

    created_by_id  BIGINT UNSIGNED NULL,
    approved_by_id BIGINT UNSIGNED NULL,
    approved_at    TIMESTAMP NULL DEFAULT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    INDEX idx_po_status     (status),
    INDEX idx_po_supplier   (supplier_id),
    INDEX idx_po_req        (requisition_id),
    INDEX idx_po_created_by (created_by_id),
    INDEX idx_po_deleted_at (deleted_at),

    CONSTRAINT fk_po_requisition FOREIGN KEY (requisition_id) REFERENCES requisitions(id)  ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_po_office      FOREIGN KEY (office_id)      REFERENCES offices(id)        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_po_location    FOREIGN KEY (location_id)    REFERENCES locations(id)      ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_po_supplier    FOREIGN KEY (supplier_id)    REFERENCES suppliers(id)      ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_po_created_by  FOREIGN KEY (created_by_id)  REFERENCES users(id)          ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_po_approved_by FOREIGN KEY (approved_by_id) REFERENCES users(id)          ON UPDATE CASCADE ON DELETE SET NULL

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- =============================================
-- TABLE: purchase_order_items
-- =============================================
CREATE TABLE IF NOT EXISTS purchase_order_items (
    id    BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    po_id BIGINT UNSIGNED NOT NULL,

    requisition_item_id BIGINT UNSIGNED NULL,
    item_id             BIGINT UNSIGNED NOT NULL,

    order_quantity    DECIMAL(15,2) NOT NULL,
    received_quantity DECIMAL(15,2) NOT NULL DEFAULT 0,

    uom_id     BIGINT UNSIGNED NULL,
    unit_price DECIMAL(15,2)   NOT NULL,

    vat_percentage      DECIMAL(5,2)  NOT NULL DEFAULT 0,
    vat_amount          DECIMAL(15,2) NOT NULL DEFAULT 0,
    discount_percentage DECIMAL(5,2)  NOT NULL DEFAULT 0,
    discount_amount     DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_amount        DECIMAL(15,2) NOT NULL DEFAULT 0,

    remarks TEXT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_poi_po_id            (po_id),
    INDEX idx_poi_item_id          (item_id),
    INDEX idx_poi_requisition_item (requisition_item_id),

    CONSTRAINT fk_poi_po               FOREIGN KEY (po_id)               REFERENCES purchase_orders(id)   ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_poi_item             FOREIGN KEY (item_id)             REFERENCES items(id)             ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_poi_requisition_item FOREIGN KEY (requisition_item_id) REFERENCES requisition_items(id) ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_poi_uom              FOREIGN KEY (uom_id)              REFERENCES uoms(id)              ON UPDATE CASCADE ON DELETE SET NULL

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- =============================================
-- TABLE: goods_receipt_notes
-- =============================================
CREATE TABLE IF NOT EXISTS goods_receipt_notes (
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    grn_number VARCHAR(50) NOT NULL UNIQUE,
    grn_date   DATE        NOT NULL,

    receive_type ENUM('DIRECT','AGAINST_PO')              NOT NULL,
    status       ENUM('DRAFT','CONFIRMED','CANCELLED') NOT NULL DEFAULT 'DRAFT',

    po_id          BIGINT UNSIGNED NULL,
    requisition_id BIGINT UNSIGNED NULL,

    office_id   BIGINT UNSIGNED NOT NULL,
    location_id BIGINT UNSIGNED NOT NULL,
    supplier_id BIGINT UNSIGNED NOT NULL,

    challan_no               VARCHAR(100) NULL,
    challan_date             DATE         NULL,
    sales_invoice_number     VARCHAR(100) NULL,
    vat_challan_number       VARCHAR(100) NULL,
    delivery_number          VARCHAR(100) NULL,
    shipping_address         TEXT         NULL,
    shipment_document_number VARCHAR(100) NULL,

    payment_method_id BIGINT UNSIGNED NULL,
    remarks           TEXT            NULL,
    attachment_path   VARCHAR(500)    NULL,

    created_by_id  BIGINT UNSIGNED NULL,
    received_by_id BIGINT UNSIGNED NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_grn_status     (status),
    INDEX idx_grn_po_id      (po_id),
    INDEX idx_grn_supplier   (supplier_id),
    INDEX idx_grn_created_by (created_by_id),

    CONSTRAINT fk_grn_po             FOREIGN KEY (po_id)            REFERENCES purchase_orders(id) ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_grn_requisition    FOREIGN KEY (requisition_id)   REFERENCES requisitions(id)    ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_grn_office         FOREIGN KEY (office_id)        REFERENCES offices(id)         ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_grn_location       FOREIGN KEY (location_id)      REFERENCES locations(id)       ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_grn_supplier       FOREIGN KEY (supplier_id)      REFERENCES suppliers(id)       ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_grn_payment_method FOREIGN KEY (payment_method_id) REFERENCES payment_modes(id)  ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_grn_created_by     FOREIGN KEY (created_by_id)    REFERENCES users(id)           ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_grn_received_by    FOREIGN KEY (received_by_id)   REFERENCES users(id)           ON UPDATE CASCADE ON DELETE SET NULL

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- =============================================
-- TABLE: grn_items
-- =============================================
CREATE TABLE IF NOT EXISTS grn_items (
    id     BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    grn_id BIGINT UNSIGNED NOT NULL,

    po_item_id BIGINT UNSIGNED NULL,
    item_id    BIGINT UNSIGNED NOT NULL,

    received_quantity   DECIMAL(15,2) NOT NULL,
    uom_id              BIGINT UNSIGNED NULL,
    purchase_price      DECIMAL(15,2) NOT NULL,

    vat_percentage      DECIMAL(5,2)  NOT NULL DEFAULT 0,
    vat_amount          DECIMAL(15,2) NOT NULL DEFAULT 0,
    discount_percentage DECIMAL(5,2)  NOT NULL DEFAULT 0,
    discount_amount     DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_amount        DECIMAL(15,2) NOT NULL DEFAULT 0,

    category_id       BIGINT UNSIGNED NULL,
    sub_category_id   BIGINT UNSIGNED NULL,
    minor_category_id BIGINT UNSIGNED NULL,

    remarks TEXT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_grni_grn_id  (grn_id),
    INDEX idx_grni_item_id (item_id),
    INDEX idx_grni_po_item (po_item_id),

    CONSTRAINT fk_grni_grn          FOREIGN KEY (grn_id)           REFERENCES goods_receipt_notes(id)  ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_grni_po_item      FOREIGN KEY (po_item_id)       REFERENCES purchase_order_items(id) ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_grni_item         FOREIGN KEY (item_id)          REFERENCES items(id)                ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_grni_uom          FOREIGN KEY (uom_id)           REFERENCES uoms(id)                 ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_grni_category     FOREIGN KEY (category_id)      REFERENCES categories(id)           ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_grni_sub_category FOREIGN KEY (sub_category_id)  REFERENCES sub_categories(id)       ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_grni_minor_cat    FOREIGN KEY (minor_category_id) REFERENCES minor_categories(id)    ON UPDATE CASCADE ON DELETE SET NULL

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- =============================================
-- TABLE: stock_transactions
-- =============================================
CREATE TABLE IF NOT EXISTS stock_transactions (
    id                 BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    transaction_number VARCHAR(50) NOT NULL UNIQUE,
    transaction_type   ENUM('GRN','ISSUE','TRANSFER','ADJUSTMENT','RETURN','DAMAGE','EXPIRED') NOT NULL,

    item_id     BIGINT UNSIGNED NOT NULL,
    location_id BIGINT UNSIGNED NOT NULL,

    quantity_change DECIMAL(15,2) NOT NULL,
    before_quantity DECIMAL(15,2) NOT NULL DEFAULT 0,
    after_quantity  DECIMAL(15,2) NOT NULL DEFAULT 0,
    unit_cost       DECIMAL(15,2) NOT NULL DEFAULT 0,

    reference_type VARCHAR(50)     NULL,
    reference_id   BIGINT UNSIGNED NULL,
    grn_item_id    BIGINT UNSIGNED NULL,

    remarks       TEXT            NULL,
    created_by_id BIGINT UNSIGNED NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_st_item_location (item_id, location_id),
    INDEX idx_st_reference     (reference_type, reference_id),
    INDEX idx_st_type          (transaction_type),
    INDEX idx_st_created_at    (created_at),
    INDEX idx_st_grn_item      (grn_item_id),

    CONSTRAINT fk_st_item       FOREIGN KEY (item_id)       REFERENCES items(id)     ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_st_location   FOREIGN KEY (location_id)   REFERENCES locations(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_st_grn_item   FOREIGN KEY (grn_item_id)   REFERENCES grn_items(id) ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_st_created_by FOREIGN KEY (created_by_id) REFERENCES users(id)     ON UPDATE CASCADE ON DELETE SET NULL

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
