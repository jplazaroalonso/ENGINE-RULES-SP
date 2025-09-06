# Behaviour Tests (BDD) - Rule Creation and Management

## TC-BDD-01: End-to-End Rule Creation Journey

**Scenario**: Business analyst creates a complete promotional rule  
**Given** I am a business analyst with rule creation permissions  
**And** I need to create a "Summer Electronics Discount" promotion  
**And** the rule template library contains discount templates  
**When** I navigate to the rule creation interface  
**And** I select the "Percentage Discount" template  
**And** I configure the rule with electronics category and 20% discount  
**And** I validate the rule syntax  
**And** I test the rule with sample transaction data  
**And** I save the rule as a draft  
**Then** the rule should be created successfully  
**And** the rule should appear in my rules list with DRAFT status  
**And** the rule should be ready for submission to approval workflow

---

## TC-BDD-02: Template-Based Rule Creation

**Scenario**: User creates rule using existing template  
**Given** I am a business user with rule creation access  
**And** I want to create a loyalty points rule quickly  
**And** a "Loyalty Points" template exists with standard configuration  
**When** I browse the template library  
**And** I preview the "Loyalty Points" template details  
**And** I select and apply the template  
**And** I customize the points calculation formula  
**And** I set customer tier conditions  
**Then** the template should populate the rule editor correctly  
**And** my customizations should be preserved  
**And** the rule should validate successfully with custom content  
**And** I should be able to save the customized rule

---

## TC-BDD-03: Real-Time DSL Validation

**Scenario**: User receives immediate feedback on rule syntax errors  
**Given** I am creating a new discount rule  
**And** I am entering DSL conditions manually  
**And** the real-time validation system is active  
**When** I type invalid DSL syntax: "category == electronics AND amount >"  
**Then** the system should immediately highlight the syntax error  
**And** specific error messages should appear with correction suggestions  
**And** the error should be positioned at the exact location  
**When** I correct the syntax to: "category = 'electronics' AND amount > 100"  
**Then** the validation error should disappear  
**And** the rule should show as syntactically valid  
**And** I should be able to proceed to testing

---

## TC-BDD-04: Rule Testing with Multiple Scenarios

**Scenario**: User tests rule with various transaction scenarios  
**Given** I have created a "Volume Discount" rule  
**And** the rule offers 10% discount for purchases over $200  
**And** I want to verify the rule works correctly  
**When** I test with a $150 purchase (below threshold)  
**Then** no discount should be applied  
**When** I test with a $250 purchase (above threshold)  
**Then** a $25 discount should be calculated and displayed  
**When** I test with a $200 purchase (exactly at threshold)  
**Then** no discount should be applied (exclusive threshold)  
**And** all test results should be displayed with step-by-step execution details

---

## TC-BDD-05: Rule Cloning for Similar Promotions

**Scenario**: Business user clones existing rule to create similar promotion  
**Given** I have an approved "Holiday Discount" rule offering 15% off  
**And** I need to create a similar "New Year Discount" rule with 25% off  
**And** I want to avoid recreating the entire rule structure  
**When** I select the "Holiday Discount" rule from my rules list  
**And** I choose the "Clone Rule" option  
**And** I modify the cloned rule name to "New Year Discount"  
**And** I update the discount percentage from 15% to 25%  
**And** I save the modified clone  
**Then** a new independent rule should be created  
**And** the original "Holiday Discount" rule should remain unchanged  
**And** the clone relationship should be tracked in the rule metadata  
**And** both rules should exist independently in the system

---

## TC-BDD-06: Draft Rule Management Workflow

**Scenario**: User saves incomplete rule as draft and continues later  
**Given** I am creating a complex promotional rule  
**And** I need to leave the system before completing the rule  
**And** I have filled in the rule name and category  
**When** I save the incomplete rule as a draft  
**And** I navigate away from the rule editor  
**And** I return to the system later  
**And** I open my saved draft  
**Then** all previously entered information should be preserved  
**And** I should be able to continue editing from where I left off  
**And** I should be able to complete the rule definition  
**And** I should be able to change the status from DRAFT to UNDER_REVIEW

---

## TC-BDD-07: Rule History and Change Tracking

**Scenario**: User reviews complete history of rule changes  
**Given** I have a rule that has gone through multiple revisions  
**And** the rule has changed status several times  
**And** multiple users have made modifications  
**When** I access the rule history view  
**Then** I should see all changes in chronological order  
**And** each change should show the user who made it  
**And** status transitions should be clearly marked with timestamps  
**And** I should be able to compare different versions of the rule  
**And** I should be able to export the history for audit purposes

---

## TC-BDD-08: Business Rule Conflict Detection

**Scenario**: System detects and warns about conflicting rules  
**Given** I have an active rule offering "10% discount on electronics"  
**And** I am creating a new rule offering "15% discount on electronics"  
**And** both rules have overlapping conditions  
**When** I validate the new rule  
**Then** the system should detect the potential conflict  
**And** a warning should be displayed explaining the conflict  
**And** suggestions for resolution should be provided  
**And** I should be able to proceed with awareness of the conflict  
**Or** I should be able to modify the rule to avoid the conflict

---

## TC-BDD-09: Template Customization and Reuse

**Scenario**: User customizes template for specific business needs  
**Given** I frequently create similar discount rules  
**And** the standard discount template doesn't match my needs exactly  
**And** I want to create a customized template version  
**When** I select the standard discount template  
**And** I modify the default variables and conditions  
**And** I save my customizations  
**Then** my customized version should be available for future use  
**And** the original template should remain unchanged  
**And** I should be able to create rules from my customized template  
**And** I should be able to reset to the original template if needed

---

## TC-BDD-10: Rule Status Workflow Management

**Scenario**: Rule progresses through complete approval workflow  
**Given** I have created and tested a promotional rule  
**And** the rule is ready for production use  
**And** my organization requires approval for rule activation  
**When** I submit the rule for approval (UNDER_REVIEW status)  
**Then** the rule should be locked from editing  
**And** an approver should be notified of the pending review  
**When** the approver reviews and approves the rule  
**Then** the rule status should change to APPROVED  
**And** I should be notified of the approval  
**When** I activate the approved rule  
**Then** the rule status should change to ACTIVE  
**And** the rule should be available for transaction evaluation

---

## Cross-Cutting Behaviour Scenarios

### Performance Behaviour

**Scenario**: System maintains performance under normal load  
**Given** 10 business users are creating rules simultaneously  
**And** the system is processing rule validations and tests  
**When** I perform any rule creation operation  
**Then** the response time should remain within acceptable limits  
**And** the system should remain responsive  
**And** no operations should fail due to load

### Error Recovery Behaviour

**Scenario**: System gracefully handles connection interruptions  
**Given** I am in the middle of creating a complex rule  
**And** I have unsaved changes in the rule editor  
**When** my network connection is interrupted  
**And** the connection is restored after 30 seconds  
**Then** my unsaved changes should be preserved  
**And** I should be able to continue editing  
**And** the system should attempt to save my work automatically

### Security Behaviour

**Scenario**: System enforces proper access controls  
**Given** I am a business user with limited permissions  
**And** I should not have access to administrative features  
**When** I attempt to access restricted functionality  
**Then** the system should deny access appropriately  
**And** I should receive a clear message about permission requirements  
**And** my legitimate actions should continue to work normally

### Integration Behaviour

**Scenario**: Rule creation integrates with approval workflow  
**Given** I submit a rule for approval  
**And** the approval workflow system is integrated  
**When** the rule moves through approval states  
**Then** notifications should be sent to appropriate stakeholders  
**And** audit events should be properly logged  
**And** the rule status should synchronize correctly across systems

---

## Scenario Outline: Multiple Rule Types Creation

**Scenario Outline**: Creating different types of business rules  
**Given** I am a business analyst creating a `<rule_type>` rule  
**And** I want to apply `<action>` when `<condition>` is met  
**When** I create the rule using appropriate template or custom DSL  
**And** I test the rule with relevant transaction data  
**Then** the rule should correctly apply `<expected_result>`

**Examples**:
| rule_type | condition | action | expected_result |
|-----------|-----------|---------|----------------|
| Discount | purchase > $100 | 10% discount | $10 off $100 purchase |
| Loyalty | customer = 'gold' | 2x points | Double points awarded |
| Coupon | category = 'books' | Generate coupon | New coupon created |
| Shipping | order_total > $50 | Free shipping | $0 shipping cost |
| Bundle | buy 2 items | Get 1 free | Third item free |

---

## Background Scenarios

### Common Setup
**Given** the Rules Engine system is running  
**And** the database contains sample templates and data  
**And** I am authenticated as a business user  
**And** I have appropriate permissions for rule creation  
**And** the validation and testing services are available

### Common Teardown
**After** each scenario completes  
**Then** test data should be cleaned up appropriately  
**And** any temporary rules should be removed  
**And** system state should be restored for next test

---

## Tags and Organization

### Feature Tags
- @rule-creation
- @template-management
- @validation
- @testing
- @workflow
- @audit

### Priority Tags
- @critical (TC-BDD-01, TC-BDD-03, TC-BDD-04, TC-BDD-10)
- @high (TC-BDD-02, TC-BDD-05, TC-BDD-08)
- @medium (TC-BDD-06, TC-BDD-07, TC-BDD-09)

### Component Tags
- @dsl-engine
- @template-system
- @validation-service
- @testing-framework
- @workflow-engine
- @audit-system
