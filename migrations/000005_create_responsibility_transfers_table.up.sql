CREATE TABLE IF NOT EXISTS `responsibility_transfers` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `from_user_id` BIGINT UNSIGNED NOT NULL,
    `to_user_id` BIGINT UNSIGNED NOT NULL,
    `from_date` DATE NOT NULL,
    `to_date` DATE NOT NULL,
    `remarks` TEXT,
    `status` ENUM('Pending', 'Approved', 'Rejected') DEFAULT 'Pending',
    `approved_by_id` BIGINT UNSIGNED,
    `approved_at` TIMESTAMP NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX (`from_user_id`),
    INDEX (`to_user_id`),
    INDEX (`approved_by_id`),
    CONSTRAINT `fk_rt_from_user` FOREIGN KEY (`from_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_rt_to_user` FOREIGN KEY (`to_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_rt_approved_by` FOREIGN KEY (`approved_by_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
