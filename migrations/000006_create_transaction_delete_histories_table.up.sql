CREATE TABLE IF NOT EXISTS `transaction_delete_histories` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `transaction_type` VARCHAR(100) NOT NULL,
    `invoice_ref_no` VARCHAR(100) NOT NULL,
    `amount` DECIMAL(15, 2) DEFAULT 0.00,
    `deleted_by_id` BIGINT UNSIGNED,
    `reason_remarks` TEXT,
    `deleted_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX (`transaction_type`),
    INDEX (`deleted_by_id`),
    CONSTRAINT `fk_tdh_deleted_by` FOREIGN KEY (`deleted_by_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
