Feature: Settings Management
  As a system administrator
  I want to manage configurations, feature flags, user preferences, and organization settings
  So that I can control the behavior of the Rules Engine platform

  Background:
    Given the settings management service is running
    And the database is clean
    And the event bus is ready

  Scenario: Create a new configuration
    Given I have a valid configuration request
    When I create a configuration with key "database.host"
    Then the configuration should be created successfully
    And the configuration should have the correct values
    And a configuration created event should be published

  Scenario: Create a configuration with duplicate key
    Given a configuration with key "database.host" already exists
    When I try to create a configuration with key "database.host"
    Then the creation should fail with error "configuration with key 'database.host' already exists"

  Scenario: Update an existing configuration
    Given a configuration with key "database.host" exists
    When I update the configuration value
    Then the configuration should be updated successfully
    And the updated_at timestamp should be changed
    And a configuration updated event should be published

  Scenario: Delete a configuration
    Given a configuration with key "database.host" exists
    When I delete the configuration
    Then the configuration should be deleted successfully
    And a configuration deleted event should be published

  Scenario: List configurations with pagination
    Given I have 25 configurations in the system
    When I request the first page with limit 10
    Then I should receive 10 configurations
    And the pagination should show page 1 of 3
    And the total count should be 25

  Scenario: Filter configurations by category
    Given I have configurations in categories "database", "cache", and "api"
    When I filter configurations by category "database"
    Then I should receive only database configurations
    And all returned configurations should have category "database"

  Scenario: Filter configurations by environment
    Given I have configurations for environments "development", "staging", and "production"
    When I filter configurations by environment "development"
    Then I should receive only development configurations
    And all returned configurations should have environment "development"

  Scenario: Create a new feature flag
    Given I have a valid feature flag request
    When I create a feature flag with key "new-feature"
    Then the feature flag should be created successfully
    And the feature flag should have the correct values
    And a feature flag created event should be published

  Scenario: Update feature flag rollout percentage
    Given a feature flag with key "new-feature" exists
    When I update the rollout percentage to 75
    Then the feature flag should be updated successfully
    And the rollout percentage should be 75
    And a feature flag updated event should be published

  Scenario: Activate a feature flag
    Given a feature flag with key "new-feature" exists and is inactive
    When I update the status to active
    Then the feature flag should be active
    And a feature flag updated event should be published

  Scenario: Deactivate a feature flag
    Given a feature flag with key "new-feature" exists and is active
    When I update the status to inactive
    Then the feature flag should be inactive
    And a feature flag updated event should be published

  Scenario: Create a user preference
    Given I have a valid user preference request
    When I create a user preference for user "user123"
    Then the user preference should be created successfully
    And the user preference should have the correct values
    And a user preference created event should be published

  Scenario: Update user preference
    Given a user preference for user "user123" exists
    When I update the preference value
    Then the user preference should be updated successfully
    And the updated_at timestamp should be changed
    And a user preference updated event should be published

  Scenario: Delete user preference
    Given a user preference for user "user123" exists
    When I delete the user preference
    Then the user preference should be deleted successfully
    And a user preference deleted event should be published

  Scenario: List user preferences by user
    Given I have preferences for users "user123" and "user456"
    When I filter preferences by user "user123"
    Then I should receive only preferences for user "user123"
    And all returned preferences should have user_id "user123"

  Scenario: Create an organization setting
    Given I have a valid organization setting request
    When I create an organization setting for organization "org123"
    Then the organization setting should be created successfully
    And the organization setting should have the correct values
    And an organization setting created event should be published

  Scenario: Update organization setting
    Given an organization setting for organization "org123" exists
    When I update the setting value
    Then the organization setting should be updated successfully
    And the updated_at timestamp should be changed
    And an organization setting updated event should be published

  Scenario: Delete organization setting
    Given an organization setting for organization "org123" exists
    When I delete the organization setting
    Then the organization setting should be deleted successfully
    And an organization setting deleted event should be published

  Scenario: List organization settings by organization
    Given I have settings for organizations "org123" and "org456"
    When I filter settings by organization "org123"
    Then I should receive only settings for organization "org123"
    And all returned settings should have organization_id "org123"

  Scenario: Encrypt sensitive configuration
    Given I have a configuration with sensitive data
    When I create the configuration with encryption enabled
    Then the configuration should be created successfully
    And the is_encrypted flag should be true
    And the value should be encrypted

  Scenario: Cache configuration for performance
    Given a configuration with key "database.host" exists
    When I retrieve the configuration multiple times
    Then the first request should hit the database
    And subsequent requests should hit the cache
    And the response should be consistent

  Scenario: Handle concurrent configuration updates
    Given a configuration with key "database.host" exists
    When two users update the configuration simultaneously
    Then both updates should be processed
    And the final state should be consistent
    And both users should receive updated events

  Scenario: Validate configuration constraints
    Given I have an invalid configuration request
    When I try to create the configuration
    Then the creation should fail
    And I should receive validation error messages
    And no configuration should be created

  Scenario: Handle database connection failure
    Given the database is unavailable
    When I try to create a configuration
    Then the operation should fail gracefully
    And I should receive an appropriate error message
    And no partial state should be created

  Scenario: Handle event bus failure
    Given the event bus is unavailable
    When I create a configuration
    Then the configuration should be created successfully
    But the event publishing should fail
    And I should receive a warning about the event failure

  Scenario: Audit configuration changes
    Given a configuration with key "database.host" exists
    When I update the configuration
    Then the change should be logged
    And the audit log should contain the old and new values
    And the audit log should contain the user who made the change

  Scenario: Backup and restore configurations
    Given I have multiple configurations in the system
    When I backup the configurations
    Then the backup should contain all configurations
    And the backup should be in a consistent state
    When I restore from the backup
    Then all configurations should be restored correctly
    And the restored configurations should be identical to the original
