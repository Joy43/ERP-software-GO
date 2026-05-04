CREATE TABLE IF NOT EXISTS `purchase_receives` (
    `purchase_receive_id`   BIGINT UNSIGNED     NOT NULL AUTO_INCREMENT,
    `created_date`          DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `challan_no`            VARCHAR(255)        NOT NULL,
    `remarks`               TEXT                NULL,
    `shipping_address`      TEXT                NULL,
    `shipment_document`     VARCHAR(500)        NULL,
    `sales_invoice_no`      VARCHAR(255)        NULL,
    `vat_challan_no`        VARCHAR(255)        NULL,
    `challan_date`          DATETIME            NOT NULL,
    `delivery_number`       VARCHAR(255)        NULL,
    `total_amount`          DECIMAL(18, 4)      NOT NULL DEFAULT 0,
    `specification`         TEXT                NULL,
    `brand`                 VARCHAR(255)        NULL,
    `selling_price`         DECIMAL(18, 4)      NOT NULL DEFAULT 0,
    `rate`                  DECIMAL(18, 4)      NOT NULL DEFAULT 0,
    `quantity`              DECIMAL(18, 4)      NOT NULL DEFAULT 0,
    `unit`                  VARCHAR(100)        NOT NULL,

    -- ---------- Foreign Keys ----------
    `office_id`             BIGINT UNSIGNED     NOT NULL,        -- RESTRICT → NOT NULL ✓
    `category_id`           BIGINT UNSIGNED     NULL,            -- SET NULL  → NULL    ✓
    `sub_category_id`       BIGINT UNSIGNED     NULL,            -- SET NULL  → NULL    ✓
    `minor_category_id`     BIGINT UNSIGNED     NULL,            -- SET NULL  → NULL    ✓
    `location_id`           BIGINT UNSIGNED     NOT NULL,        -- RESTRICT → NOT NULL ✓
    `payment_mode_id`       BIGINT UNSIGNED     NULL,            -- SET NULL  → NULL    ✓ (was NOT NULL, caused error)
    `supplier_id`           BIGINT UNSIGNED     NULL,            -- SET NULL  → NULL    ✓ (was NOT NULL)
    `item_id`               BIGINT UNSIGNED     NOT NULL,        -- RESTRICT → NOT NULL ✓
    `file_id`               BIGINT UNSIGNED     NULL,            -- SET NULL  → NULL    ✓

    -- ---------- Primary Key ----------
    PRIMARY KEY (`purchase_receive_id`),

    -- ---------- Unique Index ----------
    UNIQUE INDEX `idx_challan_no_supplier_id` (`challan_no`, `supplier_id`),

    -- ---------- Indexes ----------
    INDEX `idx_purchase_receives_office_id`       (`office_id`),
    INDEX `idx_purchase_receives_category_id`     (`category_id`),
    INDEX `idx_purchase_receives_sub_category_id` (`sub_category_id`),
    INDEX `idx_purchase_receives_minor_cat_id`    (`minor_category_id`),
    INDEX `idx_purchase_receives_location_id`     (`location_id`),
    INDEX `idx_purchase_receives_payment_mode_id` (`payment_mode_id`),
    INDEX `idx_purchase_receives_supplier_id`     (`supplier_id`),
    INDEX `idx_purchase_receives_item_id`         (`item_id`),
    INDEX `idx_purchase_receives_file_id`         (`file_id`),

    -- ---------- Foreign Key Constraints ----------
    CONSTRAINT `fk_pr_office`
        FOREIGN KEY (`office_id`)
        REFERENCES `offices` (`id`)
        ON UPDATE CASCADE ON DELETE RESTRICT,

    CONSTRAINT `fk_pr_category`
        FOREIGN KEY (`category_id`)
        REFERENCES `categories` (`id`)
        ON UPDATE CASCADE ON DELETE SET NULL,

    CONSTRAINT `fk_pr_sub_category`
        FOREIGN KEY (`sub_category_id`)
        REFERENCES `sub_categories` (`id`)
        ON UPDATE CASCADE ON DELETE SET NULL,

    CONSTRAINT `fk_pr_minor_category`
        FOREIGN KEY (`minor_category_id`)
        REFERENCES `minor_categories` (`id`)
        ON UPDATE CASCADE ON DELETE SET NULL,

    CONSTRAINT `fk_pr_location`
        FOREIGN KEY (`location_id`)
        REFERENCES `locations` (`id`)
        ON UPDATE CASCADE ON DELETE RESTRICT,

    CONSTRAINT `fk_pr_payment_mode`
        FOREIGN KEY (`payment_mode_id`)
        REFERENCES `payment_modes` (`id`)
        ON UPDATE CASCADE ON DELETE SET NULL,

    CONSTRAINT `fk_pr_supplier`
        FOREIGN KEY (`supplier_id`)
        REFERENCES `suppliers` (`id`)
        ON UPDATE CASCADE ON DELETE SET NULL,

    CONSTRAINT `fk_pr_item`
        FOREIGN KEY (`item_id`)
        REFERENCES `items` (`id`)
        ON UPDATE CASCADE ON DELETE RESTRICT,

    CONSTRAINT `fk_pr_file`
        FOREIGN KEY (`file_id`)
        REFERENCES `files` (`id`)
        ON UPDATE CASCADE ON DELETE SET NULL

) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS `purchase_returns` (
    `id`                    BIGINT UNSIGNED     NOT NULL AUTO_INCREMENT,
    `return_number`         VARCHAR(50)         NOT NULL,
    `quantity`              DECIMAL(15, 2)      NOT NULL DEFAULT 0,
    `selling_price`         DECIMAL(15, 2)      NOT NULL DEFAULT 0,
    `remarks`               TEXT                NULL,
    `created_date`          DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- ---------- Foreign Keys ----------
    `office_id`             BIGINT UNSIGNED     NOT NULL,
    `location_id`           BIGINT UNSIGNED     NOT NULL,
    `supplier_id`           BIGINT UNSIGNED     NOT NULL,
    `item_id`               BIGINT UNSIGNED     NOT NULL,

    -- ---------- Soft Delete ----------
    `deleted_at`            DATETIME            NULL,

    -- ---------- Timestamps ----------
    `created_at`            DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`            DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- ---------- Primary Key ----------
    PRIMARY KEY (`id`),

    -- ---------- Unique Index ----------
    UNIQUE INDEX `idx_purchase_returns_return_number` (`return_number`),

    -- ---------- Indexes ----------
    INDEX `idx_purchase_returns_office_id`      (`office_id`),
    INDEX `idx_purchase_returns_location_id`    (`location_id`),
    INDEX `idx_purchase_returns_supplier_id`    (`supplier_id`),
    INDEX `idx_purchase_returns_item_id`        (`item_id`),
    INDEX `idx_purchase_returns_deleted_at`     (`deleted_at`),

    -- ---------- Foreign Key Constraints ----------
    CONSTRAINT `fk_purchase_returns_office`
        FOREIGN KEY (`office_id`)
        REFERENCES `offices` (`id`)
        ON UPDATE CASCADE ON DELETE RESTRICT,

    CONSTRAINT `fk_purchase_returns_location`
        FOREIGN KEY (`location_id`)
        REFERENCES `locations` (`id`)
        ON UPDATE CASCADE ON DELETE RESTRICT,

    CONSTRAINT `fk_purchase_returns_supplier`
        FOREIGN KEY (`supplier_id`)
        REFERENCES `suppliers` (`id`)
        ON UPDATE CASCADE ON DELETE RESTRICT,

    CONSTRAINT `fk_purchase_returns_item`
        FOREIGN KEY (`item_id`)
        REFERENCES `items` (`id`)
        ON UPDATE CASCADE ON DELETE RESTRICT

) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;