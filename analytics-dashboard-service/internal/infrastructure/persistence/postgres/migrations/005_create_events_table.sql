-- Create domain_events table for event sourcing
CREATE TABLE IF NOT EXISTS domain_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(100) NOT NULL,
    aggregate_id VARCHAR(100) NOT NULL,
    aggregate_type VARCHAR(50) NOT NULL,
    version INTEGER NOT NULL,
    data JSONB NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}',
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT domain_events_aggregate_version_unique UNIQUE (aggregate_id, version)
);

-- Create indexes for domain_events
CREATE INDEX IF NOT EXISTS idx_domain_events_aggregate_id ON domain_events(aggregate_id);
CREATE INDEX IF NOT EXISTS idx_domain_events_event_type ON domain_events(event_type);
CREATE INDEX IF NOT EXISTS idx_domain_events_timestamp ON domain_events(timestamp);
CREATE INDEX IF NOT EXISTS idx_domain_events_aggregate_type ON domain_events(aggregate_type);

-- Create audit_logs table for tracking changes
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    user_id UUID,
    changes JSONB,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ip_address INET,
    user_agent TEXT
);

-- Create indexes for audit_logs
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity_id ON audit_logs(entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity_type ON audit_logs(entity_type);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp ON audit_logs(timestamp);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
