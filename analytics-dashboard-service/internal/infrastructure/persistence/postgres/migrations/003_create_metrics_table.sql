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

-- Create indexes for metrics
CREATE INDEX IF NOT EXISTS idx_metrics_type ON metrics(type);
CREATE INDEX IF NOT EXISTS idx_metrics_category ON metrics(category);
CREATE INDEX IF NOT EXISTS idx_metrics_is_calculated ON metrics(is_calculated);
CREATE INDEX IF NOT EXISTS idx_metrics_created_at ON metrics(created_at);
CREATE INDEX IF NOT EXISTS idx_metrics_updated_at ON metrics(updated_at);
