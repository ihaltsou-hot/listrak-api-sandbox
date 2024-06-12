-- Up Migration: Create contacts table
CREATE TABLE contacts
(
    id           SERIAL PRIMARY KEY,
    short_code   INT                                 NOT NULL,
    phone_number VARCHAR(255)                        NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE (short_code, phone_number)
);
