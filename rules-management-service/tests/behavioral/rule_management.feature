Feature: Rule Management API

  As a system administrator or business user,
  I want to create, manage, and validate business rules through the API,
  So that I can control the system's behavior dynamically.

  Background:
    Given the rules-management-service is running

  Scenario: Successfully create a new rule
    When I send a "POST" request to "/v1/rules" with the following body:
    """
    {
      "name": "Summer Promotion 2024",
      "description": "15% discount for all summer products.",
      "dsl_content": "IF product.category == 'summer' THEN discount.percentage = 15",
      "priority": "HIGH",
      "category": "PROMOTIONS",
      "tags": ["summer", "sale"]
    }
    """
    Then the response status code should be 201
    And the response body should be a JSON object with a "rule_id"

  Scenario: Attempt to create a rule with a duplicate name
    Given a rule with the name "Summer Promotion 2024" already exists
    When I send a "POST" request to "/v1/rules" with the following body:
    """
    {
      "name": "Summer Promotion 2024",
      "description": "Another promotion with the same name.",
      "dsl_content": "IF product.category == 'summer' THEN discount.percentage = 20",
      "priority": "MEDIUM",
      "category": "PROMOTIONS"
    }
    """
    Then the response status code should be 409
    And the response body should contain the error message "rule name already exists"

  Scenario: Validate a syntactically correct rule DSL
    When I send a "POST" request to "/v1/rules/validate" with the following body:
    """
    {
      "dsl_content": "IF customer.tier == 'GOLD' THEN loyalty.points = 100"
    }
    """
    Then the response status code should be 200
    And the response body should contain a JSON object where "is_valid" is true

  Scenario: Validate a syntactically incorrect rule DSL
    When I send a "POST" request to "/v1/rules/validate" with the following body:
    """
    {
      "dsl_content": "IF customer.tier = 'GOLD' APPLY loyalty.points = 100"
    }
    """
    Then the response status code should be 200
    And the response body should contain a JSON object where "is_valid" is false
    And the "errors" array should contain "DSL is missing 'THEN' keyword"
