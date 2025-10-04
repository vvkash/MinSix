-- Transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    tx_hash VARCHAR(66) UNIQUE NOT NULL,
    block_number BIGINT NOT NULL,
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42),
    value NUMERIC(78, 0) NOT NULL,
    gas_price NUMERIC(78, 0) NOT NULL,
    gas_used BIGINT NOT NULL,
    input_data TEXT,
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_from_address ON transactions(from_address);
CREATE INDEX IF NOT EXISTS idx_to_address ON transactions(to_address);
CREATE INDEX IF NOT EXISTS idx_block_number ON transactions(block_number);
CREATE INDEX IF NOT EXISTS idx_timestamp ON transactions(timestamp);

-- Flagged transactions table
CREATE TABLE IF NOT EXISTS flagged_transactions (
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER REFERENCES transactions(id) ON DELETE CASCADE,
    tx_hash VARCHAR(66) NOT NULL,
    risk_score INTEGER NOT NULL CHECK (risk_score >= 0 AND risk_score <= 100),
    reasons TEXT[] NOT NULL,
    flagged_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'reviewed', 'false_positive', 'confirmed'))
);

CREATE INDEX IF NOT EXISTS idx_tx_hash ON flagged_transactions(tx_hash);
CREATE INDEX IF NOT EXISTS idx_risk_score ON flagged_transactions(risk_score);
CREATE INDEX IF NOT EXISTS idx_flagged_at ON flagged_transactions(flagged_at);
CREATE INDEX IF NOT EXISTS idx_status ON flagged_transactions(status);

-- Blacklisted addresses
CREATE TABLE IF NOT EXISTS blacklisted_addresses (
    id SERIAL PRIMARY KEY,
    address VARCHAR(42) UNIQUE NOT NULL,
    reason TEXT NOT NULL,
    source VARCHAR(100),
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_blacklist_address ON blacklisted_addresses(address);

-- Monitored wallets
CREATE TABLE IF NOT EXISTS monitored_wallets (
    id SERIAL PRIMARY KEY,
    address VARCHAR(42) UNIQUE NOT NULL,
    label VARCHAR(255),
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_checked TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_wallet_address ON monitored_wallets(address);

-- Platform statistics
CREATE TABLE IF NOT EXISTS statistics (
    id SERIAL PRIMARY KEY,
    metric_name VARCHAR(100) UNIQUE NOT NULL,
    metric_value NUMERIC NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert some known blacklisted addresses (examples from public sources)
INSERT INTO blacklisted_addresses (address, reason, source) VALUES
('0x0000000000000000000000000000000000000000', 'Null address', 'System'),
('0x000000000000000000000000000000000000dead', 'Burn address', 'System')
ON CONFLICT (address) DO NOTHING;

-- Initialize statistics
INSERT INTO statistics (metric_name, metric_value) VALUES
('total_transactions', 0),
('total_flagged', 0),
('false_positives', 0),
('confirmed_fraud', 0)
ON CONFLICT (metric_name) DO NOTHING;
