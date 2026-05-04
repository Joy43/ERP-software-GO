-- Create folders table
CREATE TABLE IF NOT EXISTS folders (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    parent_id BIGINT UNSIGNED,
    created_by BIGINT UNSIGNED NOT NULL,
    level INT DEFAULT 0,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (parent_id) REFERENCES folders(id) ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
    
    INDEX idx_parent_id (parent_id),
    INDEX idx_created_by (created_by),
    INDEX idx_is_deleted (is_deleted),
    INDEX idx_name (name),
    INDEX idx_parent_deleted (parent_id, is_deleted)
);

-- Create files table
CREATE TABLE IF NOT EXISTS files (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    folder_id BIGINT UNSIGNED,
    name VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    extension VARCHAR(50),
    mime_type VARCHAR(100),
    size BIGINT UNSIGNED NOT NULL,
    storage_path VARCHAR(500) NOT NULL,
    uploaded_by BIGINT UNSIGNED NOT NULL,
    file_type ENUM('image', 'document', 'video', 'audio', 'archive','zip', 'other') DEFAULT 'other',
    is_public BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE SET NULL,
    FOREIGN KEY (uploaded_by) REFERENCES users(id) ON DELETE CASCADE,
    
    INDEX idx_folder_id (folder_id),
    INDEX idx_uploaded_by (uploaded_by),
    INDEX idx_is_deleted (is_deleted),
    INDEX idx_name (name),
    INDEX idx_file_type (file_type),
    INDEX idx_folder_deleted (folder_id, is_deleted),
    INDEX idx_uploaded_created (uploaded_by, created_at)
);




-- Create storage_usage table for tracking user storage
CREATE TABLE IF NOT EXISTS storage_usage (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL UNIQUE,
    total_files BIGINT UNSIGNED DEFAULT 0,
    total_bytes BIGINT UNSIGNED DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    INDEX idx_user_id (user_id)
);

-- Create file_history table for audit trail
CREATE TABLE IF NOT EXISTS file_history (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    file_id BIGINT UNSIGNED NOT NULL,
    action ENUM('upload', 'move', 'rename', 'delete', 'restore', 'share') NOT NULL,
    performed_by BIGINT UNSIGNED NOT NULL,
    details JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE,
    FOREIGN KEY (performed_by) REFERENCES users(id) ON DELETE CASCADE,
    
    INDEX idx_file_id (file_id),
    INDEX idx_performed_by (performed_by),
    INDEX idx_action (action),
    INDEX idx_created_at (created_at)
);
