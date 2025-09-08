-- Migration: Create campaign metrics table
-- Version: 002
-- Description: Creates the campaign_metrics table for storing campaign performance data

CREATE TABLE IF NOT EXISTS campaign_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    impressions BIGINT NOT NULL DEFAULT 0,
    clicks BIGINT NOT NULL DEFAULT 0,
    conversions BIGINT NOT NULL DEFAULT 0,
    revenue_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    revenue_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    cost_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    cost_currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
    ctr DECIMAL(10,4) NOT NULL DEFAULT 0,
    conversion_rate DECIMAL(10,4) NOT NULL DEFAULT 0,
    cost_per_click DECIMAL(15,2) NOT NULL DEFAULT 0,
    cost_per_conversion DECIMAL(15,2) NOT NULL DEFAULT 0,
    roas DECIMAL(10,4) NOT NULL DEFAULT 0,
    roi DECIMAL(10,4) NOT NULL DEFAULT 0,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT campaign_metrics_positive CHECK (
        impressions >= 0 AND clicks >= 0 AND conversions >= 0 AND
        revenue_amount >= 0 AND cost_amount >= 0 AND
        ctr >= 0 AND conversion_rate >= 0 AND
        cost_per_click >= 0 AND cost_per_conversion >= 0
    ),
    CONSTRAINT campaign_metrics_currency_check CHECK (
        LENGTH(revenue_currency) = 3 AND LENGTH(cost_currency) = 3
    ),
    CONSTRAINT campaign_metrics_campaign_unique UNIQUE (campaign_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_campaign_metrics_campaign_id ON campaign_metrics(campaign_id);
CREATE INDEX IF NOT EXISTS idx_campaign_metrics_last_updated ON campaign_metrics(last_updated);
CREATE INDEX IF NOT EXISTS idx_campaign_metrics_impressions ON campaign_metrics(impressions);
CREATE INDEX IF NOT EXISTS idx_campaign_metrics_clicks ON campaign_metrics(clicks);
CREATE INDEX IF NOT EXISTS idx_campaign_metrics_conversions ON campaign_metrics(conversions);
CREATE INDEX IF NOT EXISTS idx_campaign_metrics_roi ON campaign_metrics(roi);

-- Comments for documentation
COMMENT ON TABLE campaign_metrics IS 'Stores campaign performance metrics and analytics';
COMMENT ON COLUMN campaign_metrics.id IS 'Unique metrics identifier';
COMMENT ON COLUMN campaign_metrics.campaign_id IS 'Reference to the campaign';
COMMENT ON COLUMN campaign_metrics.impressions IS 'Total number of impressions';
COMMENT ON COLUMN campaign_metrics.clicks IS 'Total number of clicks';
COMMENT ON COLUMN campaign_metrics.conversions IS 'Total number of conversions';
COMMENT ON COLUMN campaign_metrics.revenue_amount IS 'Total revenue generated';
COMMENT ON COLUMN campaign_metrics.revenue_currency IS 'Revenue currency (ISO 4217)';
COMMENT ON COLUMN campaign_metrics.cost_amount IS 'Total cost spent';
COMMENT ON COLUMN campaign_metrics.cost_currency IS 'Cost currency (ISO 4217)';
COMMENT ON COLUMN campaign_metrics.ctr IS 'Click-through rate (percentage)';
COMMENT ON COLUMN campaign_metrics.conversion_rate IS 'Conversion rate (percentage)';
COMMENT ON COLUMN campaign_metrics.cost_per_click IS 'Cost per click (CPC)';
COMMENT ON COLUMN campaign_metrics.cost_per_conversion IS 'Cost per acquisition (CPA)';
COMMENT ON COLUMN campaign_metrics.roas IS 'Return on ad spend';
COMMENT ON COLUMN campaign_metrics.roi IS 'Return on investment (percentage)';
COMMENT ON COLUMN campaign_metrics.last_updated IS 'Last time metrics were updated';
