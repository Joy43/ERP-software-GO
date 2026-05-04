CREATE TABLE IF NOT EXISTS `wallets` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(100) NOT NULL UNIQUE,
    `commission_percent` DECIMAL(5, 2) DEFAULT 0.00,
    `bank_account` VARCHAR(50),
    `reference_code` VARCHAR(50),
    `is_rounding` BOOLEAN DEFAULT FALSE,
    `is_coupon` BOOLEAN DEFAULT FALSE,
    `is_wallet_charge` BOOLEAN DEFAULT FALSE,
    `is_active` BOOLEAN DEFAULT TRUE,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Seed initial data
INSERT IGNORE INTO `wallets` (`name`, `commission_percent`, `bank_account`, `reference_code`) VALUES
('bKash', 12.9, '2073088790001', 'N/A'),
('Eastern Bank PLC', 1.2, '1111070004703', 'N/A'),
('Dutch-Bangla Bank PLC', 1.2, '2581100236812', 'N/A'),
('Pubali Bank PLC', 1, '4648901010679', 'N/A'),
('BRAC Bank PLC', 1.3, '2073088790001', 'N/A'),
('City Bank PLC', 1.8, '1264551582001', 'N/A'),
('Nagad', 12.5, '4648901010679', 'N/A'),
('PBL Mobile Banking', 1, '4648901010679', 'N/A');
