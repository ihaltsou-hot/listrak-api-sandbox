-- Up Migration: Create subscriptions table
CREATE TABLE subscriptions
(
    id                    SERIAL PRIMARY KEY,
    contact_id            INT                                 NOT NULL,
    phone_list            INT                                 NOT NULL,
    subscribed            BOOLEAN                             NOT NULL,
    pending_double_opt_in BOOLEAN                             NOT NULL,
    subscribe_date        TIMESTAMP,
    created_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_contact FOREIGN KEY (contact_id) REFERENCES contacts (id)
);

-- Create an index for searching by contact_id
CREATE INDEX idx_contact_id ON subscriptions (contact_id);
