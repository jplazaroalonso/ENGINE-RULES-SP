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

-- Create indexes for dashboards
CREATE INDEX IF NOT EXISTS idx_dashboards_owner_id ON dashboards(owner_id);
CREATE INDEX IF NOT EXISTS idx_dashboards_is_public ON dashboards(is_public);
CREATE INDEX IF NOT EXISTS idx_dashboards_created_at ON dashboards(created_at);
CREATE INDEX IF NOT EXISTS idx_dashboards_updated_at ON dashboards(updated_at);
