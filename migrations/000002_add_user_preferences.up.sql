-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_preferences (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    deleted_at DATETIME(3) NULL,
    user_id VARCHAR(255) NOT NULL,
    report_frequency VARCHAR(50) NOT NULL DEFAULT 'none',
    last_report_sent DATETIME(3) NULL,
    UNIQUE INDEX idx_user_id (user_id),
    INDEX idx_user_preferences_deleted_at (deleted_at)
);
-- +goose StatementEnd
