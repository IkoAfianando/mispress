CREATE TABLE IF NOT EXISTS post_logs
(
    post_log_id
               varchar(50)
        PRIMARY
            KEY,
    post_id
               varchar(50) not null,
    operation
               VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);