ALTER TABLE users
    ADD COLUMN email_confirmation_token      VARCHAR(255),
    ADD COLUMN email_confirmation_expires_at TIMESTAMP;