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

-- Create indexes for metric_data
CREATE INDEX IF NOT EXISTS idx_metric_data_metric_id ON metric_data(metric_id);
CREATE INDEX IF NOT EXISTS idx_metric_data_timestamp ON metric_data(timestamp);
CREATE INDEX IF NOT EXISTS idx_metric_data_metric_timestamp ON metric_data(metric_id, timestamp);
CREATE INDEX IF NOT EXISTS idx_metric_data_value ON metric_data(value);

-- Note: Partitioned tables can be added later for better performance with large datasets
-- For now, using a regular table for simplicity
