Feature: Campaign Management
  As a marketing manager
  I want to manage marketing campaigns
  So that I can create, activate, and track promotional activities

  Background:
    Given I am authenticated as a marketing manager
    And the campaigns management service is running
    And the rules management service is available

  Scenario: Successfully create a new campaign
    Given I have a valid campaign definition
    When I submit a create campaign request with:
      | name        | Summer Sale 2024 |
      | description | Summer promotion campaign with 20% discount |
      | type        | PROMOTION |
      | targetingRules | rule-123, rule-456 |
      | startDate   | 2024-06-01T00:00:00Z |
      | endDate     | 2024-08-31T23:59:59Z |
      | budget      | 10000.00 EUR |
      | createdBy   | john.doe |
    Then the campaign should be created successfully
    And the campaign status should be "DRAFT"
    And the campaign should have version 1
    And a "CampaignCreated" event should be published

  Scenario: Prevent duplicate campaign names
    Given a campaign with name "Summer Sale 2024" already exists
    When I try to create a campaign with name "Summer Sale 2024"
    Then I should receive a "BUSINESS_ERROR" response
    And the error message should contain "campaign name already exists"

  Scenario: Activate a draft campaign
    Given I have a campaign in "DRAFT" status
    And the campaign start date is in the past
    When I activate the campaign
    Then the campaign status should change to "ACTIVE"
    And a "CampaignActivated" event should be published
    And the campaign should be trackable for metrics

  Scenario: Pause an active campaign
    Given I have a campaign in "ACTIVE" status
    When I pause the campaign
    Then the campaign status should change to "PAUSED"
    And a "CampaignPaused" event should be published
    And the campaign should stop generating new impressions

  Scenario: Complete a campaign
    Given I have a campaign in "ACTIVE" status
    When I complete the campaign
    Then the campaign status should change to "COMPLETED"
    And a "CampaignCompleted" event should be published
    And the campaign should be archived

  Scenario: Cancel a campaign with reason
    Given I have a campaign in "ACTIVE" status
    When I cancel the campaign with reason "Budget exceeded"
    Then the campaign status should change to "CANCELLED"
    And a "CampaignCancelled" event should be published
    And the cancellation reason should be recorded

  Scenario: Update campaign targeting rules
    Given I have a campaign in "DRAFT" status
    When I update the targeting rules to "rule-789, rule-101"
    Then the campaign targeting rules should be updated
    And the campaign version should be incremented
    And a "CampaignUpdated" event should be published

  Scenario: Update campaign budget
    Given I have a campaign in "DRAFT" status
    When I update the budget to "15000.00 USD"
    Then the campaign budget should be updated
    And the campaign version should be incremented
    And a "CampaignUpdated" event should be published

  Scenario: Track campaign events
    Given I have an active campaign
    When a customer interacts with the campaign
    And the interaction type is "impression"
    Then the campaign metrics should be updated
    And the impression count should increase by 1
    And the last updated timestamp should be set

  Scenario: Validate campaign before activation
    Given I have a campaign in "DRAFT" status
    When I validate the campaign
    Then the validation should check targeting rules
    And the validation should check budget constraints
    And the validation should check date constraints
    And the validation result should be returned

  Scenario: List campaigns with filtering
    Given I have multiple campaigns with different statuses
    When I request campaigns with status filter "ACTIVE"
    Then only active campaigns should be returned
    And the response should include pagination information

  Scenario: Search campaigns by name
    Given I have campaigns with names "Summer Sale", "Winter Sale", "Spring Promotion"
    When I search for campaigns with query "Sale"
    Then campaigns containing "Sale" in the name should be returned
    And the search should be case-insensitive

  Scenario: Get campaign metrics
    Given I have an active campaign with tracked events
    When I request campaign metrics
    Then the response should include impressions, clicks, conversions
    And the response should include revenue and cost data
    And the response should include performance indicators (CTR, conversion rate, ROAS, ROI)

  Scenario: Handle campaign budget exceeded
    Given I have an active campaign with budget limit
    When the campaign cost exceeds the budget
    Then the campaign should be marked as budget exceeded
    And a "CampaignBudgetExceeded" event should be published
    And the campaign should be automatically paused

  Scenario: Handle campaign end date reached
    Given I have an active campaign with end date
    When the current date reaches the campaign end date
    Then the campaign should be automatically completed
    And a "CampaignCompleted" event should be published
    And the campaign should stop generating new impressions

  Scenario: Validate targeting rules integration
    Given I have a campaign with targeting rules
    When I validate the campaign
    Then the system should verify targeting rules exist in Rules Engine
    And the system should verify targeting rules are active
    And validation should fail if any targeting rule is invalid

  Scenario: Handle campaign deletion
    Given I have a campaign in "DRAFT" status
    When I delete the campaign
    Then the campaign should be soft deleted
    And a "CampaignDeleted" event should be published
    And the campaign should not appear in active listings

  Scenario: Prevent deletion of active campaigns
    Given I have a campaign in "ACTIVE" status
    When I try to delete the campaign
    Then I should receive a "BUSINESS_ERROR" response
    And the error message should contain "cannot delete active campaign"
    And the campaign should remain unchanged

  Scenario: Campaign version control
    Given I have a campaign with version 1
    When I update the campaign
    Then the campaign version should be incremented to 2
    And the updated timestamp should be set
    And a "CampaignUpdated" event should be published

  Scenario: Campaign settings validation
    Given I have a campaign with invalid settings
    When I try to create the campaign
    Then I should receive a "VALIDATION_ERROR" response
    And the error message should contain "invalid campaign settings"
    And the campaign should not be created

  Scenario: Campaign A/B testing
    Given I have a campaign with A/B testing enabled
    When the campaign is activated
    Then the system should create multiple variants
    And the system should distribute traffic according to the split
    And the system should track metrics for each variant

  Scenario: Campaign personalization
    Given I have a campaign with personalization enabled
    When a customer interacts with the campaign
    Then the system should apply personalization rules
    And the system should deliver personalized content
    And the system should track personalization effectiveness
