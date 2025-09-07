Feature: Rule Evaluation API

  As a client system,
  I want to submit rules for evaluation with a given context,
  So that I can get a deterministic result based on the rule's logic.

  Background:
    Given the rules-evaluation-service is running

  Scenario: Successfully evaluate a promotion rule
    When I send a "POST" request to "/v1/evaluate" with the following body:
    """
    {
      "rule_category": "PROMOTIONS",
      "dsl_content": "IF order.amount > 100 THEN discount.percentage = 10",
      "context": {
        "order_amount": 150.0
      }
    }
    """
    Then the response status code should be 200
    And the response body should contain a "result" object where "eligible" is true and "discount_percentage" is 10

  Scenario: Evaluate a promotion rule where the condition is not met
    When I send a "POST" request to "/v1/evaluate" with the following body:
    """
    {
      "rule_category": "PROMOTIONS",
      "dsl_content": "IF order.amount > 100 THEN discount.percentage = 10",
      "context": {
        "order_amount": 50.0
      }
    }
    """
    Then the response status code should be 200
    And the response body should contain a "result" object where "eligible" is false

  Scenario: Successfully evaluate a tax rule
    When I send a "POST" request to "/v1/evaluate" with the following body:
    """
    {
      "rule_category": "TAXES",
      "dsl_content": "IF customer.region == 'CA' THEN tax.percentage = 9.5",
      "context": {
        "customer_region": "CA"
      }
    }
    """
    Then the response status code should be 200
    And the response body should contain a "result" object where "taxable" is true and "tax_percentage" is 9.5

  Scenario: Evaluate a rule with an unsupported category
    When I send a "POST" request to "/v1/evaluate" with the following body:
    """
    {
      "rule_category": "LOYALTY",
      "dsl_content": "IF customer.tier == 'GOLD' THEN points.multiplier = 2",
      "context": {
        "customer_tier": "GOLD"
      }
    }
    """
    Then the response status code should be 500
    And the response body should contain the error message "no evaluation strategy found for category: LOYALTY"
