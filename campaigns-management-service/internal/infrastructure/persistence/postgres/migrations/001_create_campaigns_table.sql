-- Migration: Create campaigns table
-- Version: 001
-- Description: Creates the campaigns table with all necessary fields and constraints

CREATE TABLE IF NOT EXISTS campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'DRAFT',
    campaign_type VARCHAR(50) NOT NULL,
    targeting_rules JSONB NOT NULL DEFAULT '[]',
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE,
    budget_amount DECIMAL(15,2),
    budget_currency VARCHAR(3),
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    settings JSONB NOT NULL DEFAULT '{}',
    version INTEGER NOT NULL DEFAULT 1,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Constraints
    CONSTRAINT campaigns_status_check CHECK (status IN ('DRAFT', 'ACTIVE', 'PAUSED', 'COMPLETED', 'CANCELLED')),
    CONSTRAINT campaigns_type_check CHECK (campaign_type IN ('PROMOTION', 'LOYALTY', 'COUPON', 'SEGMENTATION', 'RETARGETING')),
    CONSTRAINT campaigns_dates_check CHECK (end_date IS NULL OR end_date > start_date),
    CONSTRAINT campaigns_budget_check CHECK (budget_amount IS NULL OR budget_amount > 0),
    CONSTRAINT campaigns_currency_check CHECK (budget_currency IS NULL OR LENGTH(budget_currency) = 3)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_campaigns_status ON campaigns(status);
CREATE INDEX IF NOT EXISTS idx_campaigns_type ON campaigns(campaign_type);
CREATE INDEX IF NOT EXISTS idx_campaigns_dates ON campaigns(start_date, end_date);
CREATE INDEX IF NOT EXISTS idx_campaigns_created_by ON campaigns(created_by);
CREATE INDEX IF NOT EXISTS idx_campaigns_created_at ON campaigns(created_at);
CREATE INDEX IF NOT EXISTS idx_campaigns_updated_at ON campaigns(updated_at);
CREATE INDEX IF NOT EXISTS idx_campaigns_deleted_at ON campaigns(deleted_at);

-- GIN index for JSONB fields
CREATE INDEX IF NOT EXISTS idx_campaigns_targeting_rules ON campaigns USING GIN(targeting_rules);
CREATE INDEX IF NOT EXISTS idx_campaigns_settings ON campaigns USING GIN(settings);

-- Comments for documentation
COMMENT ON TABLE campaigns IS 'Stores campaign information and configuration';
COMMENT ON COLUMN campaigns.id IS 'Unique campaign identifier';
COMMENT ON COLUMN campaigns.name IS 'Human-readable campaign name (must be unique)';
COMMENT ON COLUMN campaigns.description IS 'Detailed campaign description';
COMMENT ON COLUMN campaigns.status IS 'Campaign lifecycle status';
COMMENT ON COLUMN campaigns.campaign_type IS 'Type of campaign (promotion, loyalty, etc.)';
COMMENT ON COLUMN campaigns.targeting_rules IS 'JSON array of rule IDs used for targeting';
COMMENT ON COLUMN campaigns.start_date IS 'Campaign start date and time';
COMMENT ON COLUMN campaigns.end_date IS 'Campaign end date and time (optional)';
COMMENT ON COLUMN campaigns.budget_amount IS 'Campaign budget amount';
COMMENT ON COLUMN campaigns.budget_currency IS 'Campaign budget currency (ISO 4217)';
COMMENT ON COLUMN campaigns.created_by IS 'User who created the campaign';
COMMENT ON COLUMN campaigns.settings IS 'Campaign configuration and settings';
COMMENT ON COLUMN campaigns.version IS 'Campaign version for optimistic locking';
COMMENT ON COLUMN campaigns.deleted_at IS 'Soft delete timestamp';
