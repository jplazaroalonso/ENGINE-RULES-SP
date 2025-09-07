Feature: Rule Calculation
  As a system
  I want to calculate the result of multiple business rules
  So that I can determine a final value based on a given context

  Background:
    Given the rules calculator service is running
    And the rules evaluation service is available and configured

  Scenario: Successfully calculate a set of rules
    When I send a calculation request with the following rules:
      | ruleId |
      | rule-A |
      | rule-B |
    And the context:
      """
      {
        "customer_tier": "gold",
        "purchase_amount": 250
      }
      """
    And the evaluation service will return:
      | ruleId | value |
      | rule-A | 25.0  |
      | rule-B | 10.5  |
    Then I should receive a successful calculation response
    And the total value should be 35.5
    And the breakdown should contain "rule-A" with value 25.0
    And the breakdown should contain "rule-B" with value 10.5

  Scenario: Calculate rules where one rule fails to evaluate
    When I send a calculation request with the following rules:
      | ruleId       |
      | rule-A       |
      | rule-FAILURE |
      | rule-C       |
    And the context:
      """
      {
        "customer_tier": "silver",
        "item_count": 5
      }
      """
    And the evaluation service will return:
      | ruleId       | value | success |
      | rule-A       | 5.0   | true    |
      | rule-FAILURE | 0.0   | false   |
      | rule-C       | 7.5   | true    |
    Then I should receive a successful calculation response
    And the total value should be 12.5
    And the breakdown should contain "rule-A" with value 5.0
    And the breakdown should contain "rule-C" with value 7.5
    And the breakdown should not contain "rule-FAILURE"

  Scenario: Send a calculation request with no rules
    When I send a calculation request with an empty list of rules
    And any context
    Then I should receive a successful calculation response
    And the total value should be 0.0
    And the breakdown should be empty

  Scenario: Send a calculation request with an invalid body
    When I send a calculation request with a missing context
    Then I should receive a "400 Bad Request" response
    And the response body should indicate a validation error
