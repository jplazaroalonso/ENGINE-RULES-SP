-- Migration: Create campaign events table
-- Version: 003
-- Description: Creates the campaign_events table for storing campaign event tracking data

CREATE TABLE IF NOT EXISTS campaign_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL,
    customer_id VARCHAR(255),
    event_data JSONB NOT NULL DEFAULT '{}',
    occurred_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    revenue_amount DECIMAL(15,2),
    revenue_currency VARCHAR(3),
    cost_amount DECIMAL(15,2),
    cost_currency VARCHAR(3),
    
    -- Constraints
    CONSTRAINT campaign_events_type_check CHECK (
        event_type IN ('IMPRESSION', 'CLICK', 'CONVERSION', 'BOUNCE', 'UNSUBSCRIBE')
    ),
    CONSTRAINT campaign_events_currency_check CHECK (
        (revenue_currency IS NULL OR LENGTH(revenue_currency) = 3) AND
        (cost_currency IS NULL OR LENGTH(cost_currency) = 3)
    ),
    CONSTRAINT campaign_events_amounts_check CHECK (
        (revenue_amount IS NULL OR revenue_amount >= 0) AND
        (cost_amount IS NULL OR cost_amount >= 0)
    )
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_campaign_events_campaign_id ON campaign_events(campaign_id);
CREATE INDEX IF NOT EXISTS idx_campaign_events_type ON campaign_events(event_type);
CREATE INDEX IF NOT EXISTS idx_campaign_events_occurred_at ON campaign_events(occurred_at);
CREATE INDEX IF NOT EXISTS idx_campaign_events_customer_id ON campaign_events(customer_id);
CREATE INDEX IF NOT EXISTS idx_campaign_events_campaign_type ON campaign_events(campaign_id, event_type);

-- Composite index for common queries
CREATE INDEX IF NOT EXISTS idx_campaign_events_campaign_occurred ON campaign_events(campaign_id, occurred_at);

-- GIN index for JSONB field
CREATE INDEX IF NOT EXISTS idx_campaign_events_data ON campaign_events USING GIN(event_data);

-- Comments for documentation
COMMENT ON TABLE campaign_events IS 'Stores individual campaign events for tracking and analytics';
COMMENT ON COLUMN campaign_events.id IS 'Unique event identifier';
COMMENT ON COLUMN campaign_events.campaign_id IS 'Reference to the campaign';
COMMENT ON COLUMN campaign_events.event_type IS 'Type of event (impression, click, conversion, etc.)';
COMMENT ON COLUMN campaign_events.customer_id IS 'Customer who triggered the event (optional)';
COMMENT ON COLUMN campaign_events.event_data IS 'Additional event-specific data in JSON format';
COMMENT ON COLUMN campaign_events.occurred_at IS 'When the event occurred';
COMMENT ON COLUMN campaign_events.revenue_amount IS 'Revenue associated with this event (optional)';
COMMENT ON COLUMN campaign_events.revenue_currency IS 'Revenue currency (ISO 4217)';
COMMENT ON COLUMN campaign_events.cost_amount IS 'Cost associated with this event (optional)';
COMMENT ON COLUMN campaign_events.cost_currency IS 'Cost currency (ISO 4217)';
