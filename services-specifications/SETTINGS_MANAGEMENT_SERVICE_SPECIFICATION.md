# Settings Management Service Specification

## Executive Summary

The Settings Management Service provides centralized configuration management, system settings, user preferences, and feature flags for the entire rules engine ecosystem. It ensures consistent configuration across all services while providing flexibility for different environments and user customizations.

## Service Overview

### Purpose
- Centralized configuration management for all services
- User and organization-specific settings
- Feature flags and A/B testing support
- Environment-specific configurations
- Audit trail for all configuration changes

### Business Value
- **Centralized Control**: Single source of truth for all system configurations
- **Flexibility**: Easy configuration changes without code deployments
- **Feature Management**: Controlled feature rollouts and A/B testing
- **Compliance**: Audit trail for all configuration changes
- **Multi-tenancy**: Support for multiple organizations with isolated settings

## Technical Architecture

### Service Structure
```
settings-management-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── domain/
│   │   ├── settings/
│   │   │   ├── configuration.go
│   │   │   ├── feature_flag.go
│   │   │   ├── user_preference.go
│   │   │   ├── organization_setting.go
│   │   │   └── service.go
│   │   └── shared/
│   │       ├── errors.go
│   │       └── events.go
│   ├── application/
│   │   ├── commands/
│   │   │   ├── create_configuration.go
│   │   │   ├── update_configuration.go
│   │   │   ├── create_feature_flag.go
│   │   │   └── update_user_preference.go
│   │   └── queries/
│   │       ├── get_configuration.go
│   │       ├── get_feature_flags.go
│   │       └── get_user_preferences.go
│   ├── infrastructure/
│   │   ├── persistence/
│   │   │   └── postgres/
│   │   │       ├── configuration_repository.go
│   │   │       ├── feature_flag_repository.go
│   │   │       └── migrations/
│   │   ├── cache/
│   │   │   └── redis/
│   │   │       └── settings_cache.go
│   │   └── messaging/
│   │       └── nats/
│   │           └── event_publisher.go
│   └── interfaces/
│       ├── rest/
│       │   ├── handlers/
│       │   │   ├── configuration_handler.go
│       │   │   ├── feature_flag_handler.go
│       │   │   └── preference_handler.go
│       │   └── dto/
│       │       ├── configuration_dto.go
│       │       └── feature_flag_dto.go
│       └── grpc/
│           └── settings_service.go
├── api/
│   ├── openapi/
│   │   └── settings-api.v1.yaml
│   └── proto/
│       └── settings.proto
├── tests/
│   ├── unit/
│   ├── integration/
│   └── behavioral/
├── deployments/
│   └── k8s/
├── Dockerfile
├── go.mod
└── go.sum
```

## Domain Model

### Core Entities

#### Configuration Aggregate
```go
type Configuration struct {
    ID              ConfigurationID   `json:"id"`
    Key             string            `json:"key"`
    Value           interface{}       `json:"value"`
    Type            ConfigurationType `json:"type"`
    Category        string            `json:"category"`
    Description     string            `json:"description"`
    Environment     Environment       `json:"environment"`
    OrganizationID  *OrganizationID   `json:"organizationId,omitempty"`
    Service         *ServiceName      `json:"service,omitempty"`
    IsEncrypted     bool              `json:"isEncrypted"`
    IsSensitive     bool              `json:"isSensitive"`
    ValidationRules ValidationRules   `json:"validationRules"`
    DefaultValue    interface{}       `json:"defaultValue"`
    CreatedBy       UserID            `json:"createdBy"`
    CreatedAt       time.Time         `json:"createdAt"`
    UpdatedAt       time.Time         `json:"updatedAt"`
    Version         int               `json:"version"`
    Events          []DomainEvent     `json:"-"`
}

type ConfigurationType string

const (
    ConfigurationTypeString  ConfigurationType = "STRING"
    ConfigurationTypeNumber  ConfigurationType = "NUMBER"
    ConfigurationTypeBoolean ConfigurationType = "BOOLEAN"
    ConfigurationTypeJSON    ConfigurationType = "JSON"
    ConfigurationTypeArray   ConfigurationType = "ARRAY"
    ConfigurationTypeObject  ConfigurationType = "OBJECT"
)

type Environment string

const (
    EnvironmentDevelopment Environment = "DEVELOPMENT"
    EnvironmentStaging     Environment = "STAGING"
    EnvironmentProduction  Environment = "PRODUCTION"
    EnvironmentTesting     Environment = "TESTING"
)
```

#### Feature Flag Aggregate
```go
type FeatureFlag struct {
    ID              FeatureFlagID     `json:"id"`
    Name            string            `json:"name"`
    Description     string            `json:"description"`
    Key             string            `json:"key"`
    IsEnabled       bool              `json:"isEnabled"`
    RolloutStrategy RolloutStrategy   `json:"rolloutStrategy"`
    TargetAudience  TargetAudience    `json:"targetAudience"`
    Variants        []Variant         `json:"variants"`
    Rules           []TargetingRule   `json:"rules"`
    Environment     Environment       `json:"environment"`
    OrganizationID  *OrganizationID   `json:"organizationId,omitempty"`
    Service         *ServiceName      `json:"service,omitempty"`
    CreatedBy       UserID            `json:"createdBy"`
    CreatedAt       time.Time         `json:"createdAt"`
    UpdatedAt       time.Time         `json:"updatedAt"`
    Version         int               `json:"version"`
    Events          []DomainEvent     `json:"-"`
}

type RolloutStrategy string

const (
    RolloutStrategyAll       RolloutStrategy = "ALL"
    RolloutStrategyPercentage RolloutStrategy = "PERCENTAGE"
    RolloutStrategyUserList  RolloutStrategy = "USER_LIST"
    RolloutStrategyRules     RolloutStrategy = "RULES"
)

type Variant struct {
    Key         string      `json:"key"`
    Name        string      `json:"name"`
    Value       interface{} `json:"value"`
    Weight      int         `json:"weight"`
    Description string      `json:"description"`
}

type TargetingRule struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Conditions  []TargetingCondition   `json:"conditions"`
    Variant     string                 `json:"variant"`
    IsEnabled   bool                   `json:"isEnabled"`
}

type TargetingCondition struct {
    Attribute string      `json:"attribute"`
    Operator  string      `json:"operator"`
    Value     interface{} `json:"value"`
}
```

#### User Preference Aggregate
```go
type UserPreference struct {
    ID             UserPreferenceID   `json:"id"`
    UserID         UserID             `json:"userId"`
    OrganizationID *OrganizationID    `json:"organizationId,omitempty"`
    Category       string             `json:"category"`
    Key            string             `json:"key"`
    Value          interface{}        `json:"value"`
    Type           PreferenceType     `json:"type"`
    IsDefault      bool               `json:"isDefault"`
    CreatedAt      time.Time          `json:"createdAt"`
    UpdatedAt      time.Time          `json:"updatedAt"`
    Version        int                `json:"version"`
    Events         []DomainEvent      `json:"-"`
}

type PreferenceType string

const (
    PreferenceTypeString  PreferenceType = "STRING"
    PreferenceTypeNumber  PreferenceType = "NUMBER"
    PreferenceTypeBoolean PreferenceType = "BOOLEAN"
    PreferenceTypeJSON    PreferenceType = "JSON"
    PreferenceTypeArray   PreferenceType = "ARRAY"
)
```

#### Organization Setting Aggregate
```go
type OrganizationSetting struct {
    ID             OrganizationSettingID `json:"id"`
    OrganizationID OrganizationID        `json:"organizationId"`
    Category       string                `json:"category"`
    Key            string                `json:"key"`
    Value          interface{}           `json:"value"`
    Type           SettingType           `json:"type"`
    IsInherited    bool                  `json:"isInherited"`
    ParentID       *OrganizationID       `json:"parentId,omitempty"`
    CreatedBy      UserID                `json:"createdBy"`
    CreatedAt      time.Time             `json:"createdAt"`
    UpdatedAt      time.Time             `json:"updatedAt"`
    Version        int                   `json:"version"`
    Events         []DomainEvent         `json:"-"`
}

type SettingType string

const (
    SettingTypeString  SettingType = "STRING"
    SettingTypeNumber  SettingType = "NUMBER"
    SettingTypeBoolean SettingType = "BOOLEAN"
    SettingTypeJSON    SettingType = "JSON"
    SettingTypeArray   SettingType = "ARRAY"
)
```

## API Specification

### REST API Endpoints

#### Configuration Management
```
GET    /api/v1/configurations                 # List configurations
POST   /api/v1/configurations                 # Create configuration
GET    /api/v1/configurations/:key            # Get configuration by key
PUT    /api/v1/configurations/:key            # Update configuration
DELETE /api/v1/configurations/:key            # Delete configuration
GET    /api/v1/configurations/service/:service # Get service configurations
GET    /api/v1/configurations/env/:env        # Get environment configurations
```

#### Feature Flag Management
```
GET    /api/v1/feature-flags                  # List feature flags
POST   /api/v1/feature-flags                  # Create feature flag
GET    /api/v1/feature-flags/:key             # Get feature flag
PUT    /api/v1/feature-flags/:key             # Update feature flag
DELETE /api/v1/feature-flags/:key             # Delete feature flag
POST   /api/v1/feature-flags/:key/evaluate    # Evaluate feature flag
GET    /api/v1/feature-flags/service/:service # Get service feature flags
```

#### User Preferences
```
GET    /api/v1/users/:userId/preferences      # Get user preferences
POST   /api/v1/users/:userId/preferences      # Create user preference
PUT    /api/v1/users/:userId/preferences/:key # Update user preference
DELETE /api/v1/users/:userId/preferences/:key # Delete user preference
GET    /api/v1/users/:userId/preferences/category/:category # Get category preferences
```

#### Organization Settings
```
GET    /api/v1/organizations/:orgId/settings  # Get organization settings
POST   /api/v1/organizations/:orgId/settings  # Create organization setting
PUT    /api/v1/organizations/:orgId/settings/:key # Update organization setting
DELETE /api/v1/organizations/:orgId/settings/:key # Delete organization setting
GET    /api/v1/organizations/:orgId/settings/category/:category # Get category settings
```

#### Bulk Operations
```
POST   /api/v1/configurations/bulk            # Bulk update configurations
POST   /api/v1/feature-flags/bulk             # Bulk update feature flags
POST   /api/v1/users/:userId/preferences/bulk # Bulk update user preferences
```

#### Import/Export
```
GET    /api/v1/configurations/export          # Export configurations
POST   /api/v1/configurations/import          # Import configurations
GET    /api/v1/feature-flags/export           # Export feature flags
POST   /api/v1/feature-flags/import           # Import feature flags
```

## Business Rules and Invariants

### Configuration Rules
1. **Configuration Keys**: 
   - Must be unique per environment and service
   - Must follow naming convention (e.g., `service.category.key`)
   - Cannot be modified once created (versioned instead)

2. **Value Validation**:
   - Values must match declared type
   - Must pass validation rules if defined
   - Sensitive values must be encrypted

3. **Environment Isolation**:
   - Production configurations cannot be modified directly
   - Changes must go through staging environment first
   - Rollback capability for all changes

### Feature Flag Rules
1. **Flag Evaluation**:
   - Must be deterministic for same user/context
   - Evaluation must be fast (< 10ms)
   - Must handle missing flags gracefully

2. **Rollout Strategy**:
   - Percentage rollout must be 0-100
   - User list must contain valid user IDs
   - Rules must be valid targeting conditions

3. **Flag Lifecycle**:
   - Flags must have expiration date
   - Deprecated flags must be marked
   - Cleanup of unused flags

### User Preference Rules
1. **Preference Hierarchy**:
   - User preferences override organization defaults
   - Organization preferences override system defaults
   - System defaults are fallback

2. **Preference Validation**:
   - Values must match preference type
   - Must be within allowed ranges/options
   - Cannot exceed size limits

## Database Schema

### Configurations Table
```sql
CREATE TABLE configurations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(255) NOT NULL,
    value JSONB NOT NULL,
    type VARCHAR(20) NOT NULL,
    category VARCHAR(100) NOT NULL,
    description TEXT,
    environment VARCHAR(20) NOT NULL,
    organization_id UUID,
    service VARCHAR(100),
    is_encrypted BOOLEAN NOT NULL DEFAULT FALSE,
    is_sensitive BOOLEAN NOT NULL DEFAULT FALSE,
    validation_rules JSONB,
    default_value JSONB,
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1,
    
    CONSTRAINT configurations_key_env_unique UNIQUE (key, environment, organization_id, service),
    CONSTRAINT configurations_type_check CHECK (type IN ('STRING', 'NUMBER', 'BOOLEAN', 'JSON', 'ARRAY', 'OBJECT')),
    CONSTRAINT configurations_environment_check CHECK (environment IN ('DEVELOPMENT', 'STAGING', 'PRODUCTION', 'TESTING'))
);

CREATE INDEX idx_configurations_key ON configurations(key);
CREATE INDEX idx_configurations_environment ON configurations(environment);
CREATE INDEX idx_configurations_organization_id ON configurations(organization_id);
CREATE INDEX idx_configurations_service ON configurations(service);
CREATE INDEX idx_configurations_category ON configurations(category);
```

### Feature Flags Table
```sql
CREATE TABLE feature_flags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    key VARCHAR(255) NOT NULL UNIQUE,
    is_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    rollout_strategy VARCHAR(20) NOT NULL,
    target_audience JSONB NOT NULL DEFAULT '{}',
    variants JSONB NOT NULL DEFAULT '[]',
    rules JSONB NOT NULL DEFAULT '[]',
    environment VARCHAR(20) NOT NULL,
    organization_id UUID,
    service VARCHAR(100),
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1,
    
    CONSTRAINT feature_flags_rollout_strategy_check CHECK (rollout_strategy IN ('ALL', 'PERCENTAGE', 'USER_LIST', 'RULES')),
    CONSTRAINT feature_flags_environment_check CHECK (environment IN ('DEVELOPMENT', 'STAGING', 'PRODUCTION', 'TESTING'))
);

CREATE INDEX idx_feature_flags_key ON feature_flags(key);
CREATE INDEX idx_feature_flags_environment ON feature_flags(environment);
CREATE INDEX idx_feature_flags_organization_id ON feature_flags(organization_id);
CREATE INDEX idx_feature_flags_service ON feature_flags(service);
CREATE INDEX idx_feature_flags_is_enabled ON feature_flags(is_enabled);
```

### User Preferences Table
```sql
CREATE TABLE user_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    organization_id UUID,
    category VARCHAR(100) NOT NULL,
    key VARCHAR(255) NOT NULL,
    value JSONB NOT NULL,
    type VARCHAR(20) NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1,
    
    CONSTRAINT user_preferences_user_key_unique UNIQUE (user_id, organization_id, category, key),
    CONSTRAINT user_preferences_type_check CHECK (type IN ('STRING', 'NUMBER', 'BOOLEAN', 'JSON', 'ARRAY'))
);

CREATE INDEX idx_user_preferences_user_id ON user_preferences(user_id);
CREATE INDEX idx_user_preferences_organization_id ON user_preferences(organization_id);
CREATE INDEX idx_user_preferences_category ON user_preferences(category);
CREATE INDEX idx_user_preferences_key ON user_preferences(key);
```

### Organization Settings Table
```sql
CREATE TABLE organization_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL,
    category VARCHAR(100) NOT NULL,
    key VARCHAR(255) NOT NULL,
    value JSONB NOT NULL,
    type VARCHAR(20) NOT NULL,
    is_inherited BOOLEAN NOT NULL DEFAULT FALSE,
    parent_id UUID,
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1,
    
    CONSTRAINT organization_settings_org_key_unique UNIQUE (organization_id, category, key),
    CONSTRAINT organization_settings_type_check CHECK (type IN ('STRING', 'NUMBER', 'BOOLEAN', 'JSON', 'ARRAY'))
);

CREATE INDEX idx_organization_settings_organization_id ON organization_settings(organization_id);
CREATE INDEX idx_organization_settings_category ON organization_settings(category);
CREATE INDEX idx_organization_settings_key ON organization_settings(key);
CREATE INDEX idx_organization_settings_parent_id ON organization_settings(parent_id);
```

## Performance Requirements

### Response Time Targets
- **Get Configuration**: < 50ms (cached)
- **Get Feature Flag**: < 10ms (cached)
- **Get User Preferences**: < 100ms
- **Bulk Operations**: < 1 second for 100 items
- **Configuration Export**: < 5 seconds

### Scalability Requirements
- **Configuration Volume**: Support 100,000+ configurations
- **Feature Flags**: Support 10,000+ active feature flags
- **User Preferences**: Support 1,000,000+ user preferences
- **Concurrent Requests**: Handle 10,000+ concurrent requests
- **Cache Hit Rate**: 95%+ for frequently accessed settings

## Security Requirements

### Access Control
- **Role-Based Access**: Different permissions for different setting types
- **Environment Isolation**: Production settings require special permissions
- **Audit Logging**: Complete audit trail for all changes
- **Encryption**: Sensitive configurations encrypted at rest

### Data Protection
- **Sensitive Data**: Special handling for sensitive configurations
- **Access Logging**: Log all configuration access
- **Change Approval**: Critical changes require approval workflow
- **Backup**: Regular backups of all configurations

## Implementation Timeline

### Phase 1: Core Settings (2 weeks)
- Basic configuration management
- User preferences
- Database schema and migrations
- Unit tests

### Phase 2: Feature Flags (2 weeks)
- Feature flag management
- Evaluation engine
- Targeting rules
- Integration tests

### Phase 3: Advanced Features (2 weeks)
- Organization settings
- Bulk operations
- Import/export functionality
- Caching implementation

### Phase 4: Production Readiness (2 weeks)
- Security implementation
- Performance optimization
- Monitoring and alerting
- Deployment configuration

**Total Estimated Effort**: 8 weeks
**Team Size**: 2-3 developers
**Dependencies**: None (foundational service)
