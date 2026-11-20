ALTER TABLE users
    DROP COLUMN IF EXISTS email_confirmation_expires_at,
    DROP COLUMN IF EXISTS email_confirmation_token;