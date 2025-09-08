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

-- Create indexes for reports
CREATE INDEX IF NOT EXISTS idx_reports_owner_id ON reports(owner_id);
CREATE INDEX IF NOT EXISTS idx_reports_type ON reports(type);
CREATE INDEX IF NOT EXISTS idx_reports_status ON reports(status);
CREATE INDEX IF NOT EXISTS idx_reports_next_run ON reports(next_run);
CREATE INDEX IF NOT EXISTS idx_reports_created_at ON reports(created_at);
CREATE INDEX IF NOT EXISTS idx_reports_updated_at ON reports(updated_at);
