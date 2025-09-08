-- Analytics Dashboard Service Database Setup
-- This script creates the necessary database objects for the analytics dashboard service
-- It validates existing objects to avoid unnecessary deletions

-- Create database if it doesn't exist
SELECT 'CREATE DATABASE analytics_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'analytics_db')\gexec

-- Connect to the analytics database
\c analytics_db;

-- Create schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS analytics;

-- Set search path to include analytics schema
SET search_path TO analytics, public;

-- Create dashboards table
CREATE TABLE IF NOT EXISTS dashboards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    layout JSONB NOT NULL DEFAULT '{}',
    widgets JSONB NOT NULL DEFAULT '[]',
    filters JSONB NOT NULL DEFAULT '{}',
    refresh_interval INTEGER NOT NULL DEFAULT 300,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    owner_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1,
    
    CONSTRAINT dashboards_name_owner_unique UNIQUE (name, owner_id),
    CONSTRAINT dashboards_refresh_interval_check CHECK (refresh_interval >= 30 AND refresh_interval <= 3600)
);

-- Create indexes for dashboards if they don't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_dashboards_owner_id') THEN
        CREATE INDEX idx_dashboards_owner_id ON dashboards(owner_id);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_dashboards_is_public') THEN
        CREATE INDEX idx_dashboards_is_public ON dashboards(is_public);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_dashboards_created_at') THEN
        CREATE INDEX idx_dashboards_created_at ON dashboards(created_at);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_dashboards_updated_at') THEN
        CREATE INDEX idx_dashboards_updated_at ON dashboards(updated_at);
    END IF;
END $$;

-- Create reports table
CREATE TABLE IF NOT EXISTS reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    template JSONB NOT NULL,
    parameters JSONB NOT NULL DEFAULT '{}',
    schedule JSONB,
    output_format VARCHAR(20) NOT NULL,
    recipients JSONB NOT NULL DEFAULT '[]',
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    last_generated TIMESTAMP WITH TIME ZONE,
    next_run TIMESTAMP WITH TIME ZONE,
    owner_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1,
    
    CONSTRAINT reports_name_owner_unique UNIQUE (name, owner_id),
    CONSTRAINT reports_type_check CHECK (type IN ('PERFORMANCE', 'COMPLIANCE', 'BUSINESS', 'CUSTOM')),
    CONSTRAINT reports_output_format_check CHECK (output_format IN ('PDF', 'EXCEL', 'CSV', 'JSON', 'HTML')),
    CONSTRAINT reports_status_check CHECK (status IN ('ACTIVE', 'INACTIVE', 'GENERATING', 'ERROR'))
);

-- Create indexes for reports if they don't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_reports_owner_id') THEN
        CREATE INDEX idx_reports_owner_id ON reports(owner_id);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_reports_type') THEN
        CREATE INDEX idx_reports_type ON reports(type);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_reports_status') THEN
        CREATE INDEX idx_reports_status ON reports(status);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_reports_next_run') THEN
        CREATE INDEX idx_reports_next_run ON reports(next_run);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_reports_created_at') THEN
        CREATE INDEX idx_reports_created_at ON reports(created_at);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_reports_updated_at') THEN
        CREATE INDEX idx_reports_updated_at ON reports(updated_at);
    END IF;
END $$;

-- Create metrics table
CREATE TABLE IF NOT EXISTS metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    type VARCHAR(20) NOT NULL,
    category VARCHAR(20) NOT NULL,
    unit VARCHAR(50),
    aggregation VARCHAR(20) NOT NULL,
    data_source JSONB NOT NULL,
    dimensions JSONB NOT NULL DEFAULT '[]',
    filters JSONB NOT NULL DEFAULT '{}',
    calculation JSONB,
    is_calculated BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1,
    
    CONSTRAINT metrics_type_check CHECK (type IN ('COUNTER', 'GAUGE', 'HISTOGRAM', 'SUMMARY')),
    CONSTRAINT metrics_category_check CHECK (category IN ('PERFORMANCE', 'BUSINESS', 'SYSTEM', 'USER')),
    CONSTRAINT metrics_aggregation_check CHECK (aggregation IN ('SUM', 'AVG', 'MIN', 'MAX', 'COUNT', 'DISTINCT'))
);

-- Create indexes for metrics if they don't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_metrics_type') THEN
        CREATE INDEX idx_metrics_type ON metrics(type);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_metrics_category') THEN
        CREATE INDEX idx_metrics_category ON metrics(category);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_metrics_is_calculated') THEN
        CREATE INDEX idx_metrics_is_calculated ON metrics(is_calculated);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_metrics_created_at') THEN
        CREATE INDEX idx_metrics_created_at ON metrics(created_at);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_metrics_updated_at') THEN
        CREATE INDEX idx_metrics_updated_at ON metrics(updated_at);
    END IF;
END $$;

-- Create metric_data table
CREATE TABLE IF NOT EXISTS metric_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    metric_id UUID NOT NULL REFERENCES metrics(id) ON DELETE CASCADE,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    value DECIMAL(20,6) NOT NULL,
    dimensions JSONB NOT NULL DEFAULT '{}',
    labels JSONB NOT NULL DEFAULT '{}',
    
    CONSTRAINT metric_data_metric_timestamp_unique UNIQUE (metric_id, timestamp)
);

-- Create indexes for metric_data if they don't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_metric_data_metric_id') THEN
        CREATE INDEX idx_metric_data_metric_id ON metric_data(metric_id);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_metric_data_timestamp') THEN
        CREATE INDEX idx_metric_data_timestamp ON metric_data(timestamp);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_metric_data_metric_timestamp') THEN
        CREATE INDEX idx_metric_data_metric_timestamp ON metric_data(metric_id, timestamp);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_metric_data_value') THEN
        CREATE INDEX idx_metric_data_value ON metric_data(value);
    END IF;
END $$;

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

-- Create indexes for domain_events if they don't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_domain_events_aggregate_id') THEN
        CREATE INDEX idx_domain_events_aggregate_id ON domain_events(aggregate_id);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_domain_events_event_type') THEN
        CREATE INDEX idx_domain_events_event_type ON domain_events(event_type);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_domain_events_timestamp') THEN
        CREATE INDEX idx_domain_events_timestamp ON domain_events(timestamp);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_domain_events_aggregate_type') THEN
        CREATE INDEX idx_domain_events_aggregate_type ON domain_events(aggregate_type);
    END IF;
END $$;

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

-- Create indexes for audit_logs if they don't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_audit_logs_entity_id') THEN
        CREATE INDEX idx_audit_logs_entity_id ON audit_logs(entity_id);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_audit_logs_entity_type') THEN
        CREATE INDEX idx_audit_logs_entity_type ON audit_logs(entity_type);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_audit_logs_user_id') THEN
        CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_audit_logs_timestamp') THEN
        CREATE INDEX idx_audit_logs_timestamp ON audit_logs(timestamp);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_audit_logs_action') THEN
        CREATE INDEX idx_audit_logs_action ON audit_logs(action);
    END IF;
END $$;

-- Create migrations tracking table
CREATE TABLE IF NOT EXISTS schema_migrations (
    id SERIAL PRIMARY KEY,
    version VARCHAR(255) NOT NULL UNIQUE,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert initial migration records if they don't exist
INSERT INTO schema_migrations (version) VALUES 
    ('001_create_dashboards_table.sql'),
    ('002_create_reports_table.sql'),
    ('003_create_metrics_table.sql'),
    ('004_create_metric_data_table.sql'),
    ('005_create_events_table.sql')
ON CONFLICT (version) DO NOTHING;

-- Create a user for the analytics service if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'analytics_user') THEN
        CREATE ROLE analytics_user WITH LOGIN PASSWORD 'analytics_password';
    END IF;
END $$;

-- Grant permissions to analytics_user
GRANT USAGE ON SCHEMA analytics TO analytics_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA analytics TO analytics_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA analytics TO analytics_user;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA analytics TO analytics_user;

-- Set default privileges for future objects
ALTER DEFAULT PRIVILEGES IN SCHEMA analytics GRANT ALL ON TABLES TO analytics_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA analytics GRANT ALL ON SEQUENCES TO analytics_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA analytics GRANT ALL ON FUNCTIONS TO analytics_user;

-- Display completion message
SELECT 'Analytics Dashboard Service database setup completed successfully!' as status;
