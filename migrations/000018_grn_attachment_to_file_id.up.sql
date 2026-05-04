ALTER TABLE `goods_receipt_notes`
    DROP COLUMN `attachment_path`,
    ADD COLUMN `file_id` BIGINT UNSIGNED NULL AFTER `remarks`,
    ADD INDEX `idx_grn_file_id` (`file_id`),
    ADD CONSTRAINT `fk_grn_file`
        FOREIGN KEY (`file_id`)
        REFERENCES `files` (`id`)
        ON UPDATE CASCADE ON DELETE SET NULL;
