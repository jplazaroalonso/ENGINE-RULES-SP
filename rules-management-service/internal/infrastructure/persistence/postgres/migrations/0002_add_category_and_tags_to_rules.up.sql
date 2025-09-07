-- 0002_add_category_and_tags_to_rules.up.sql
ALTER TABLE rules
ADD COLUMN category VARCHAR(100),
ADD COLUMN tags TEXT[];

-- Add indexes for the new columns
CREATE INDEX IF NOT EXISTS idx_rules_category ON rules(category);
CREATE INDEX IF NOT EXISTS idx_rules_tags ON rules USING GIN(tags);
