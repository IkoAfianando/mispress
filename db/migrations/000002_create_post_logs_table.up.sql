CREATE TABLE IF NOT EXISTS post_logs (
                                         id SERIAL PRIMARY KEY,
                                         post_id INT NOT NULL,
                                         operation VARCHAR(20) NOT NULL,
                                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);