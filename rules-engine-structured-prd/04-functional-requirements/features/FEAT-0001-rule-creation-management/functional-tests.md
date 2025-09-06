# Functional Tests - Rule Creation and Management

## TC-FUNC-01: Complete Rule Creation Flow

**Test ID**: TC-FUNC-01  
**Priority**: High  
**Prerequisites**: User has rule creation permissions, system is accessible, template library is populated  

**Test Steps**:
1. Navigate to rule creation interface from main dashboard
2. Select "Percentage Discount" template from promotions category
3. Fill in rule name: "Summer Sale - 20% Electronics"
4. Define DSL conditions: `category = 'electronics' AND purchase_amount > 100`
5. Set discount action: `discount_percentage = 20`
6. Validate rule syntax using validation button
7. Test rule with sample data: Electronics item, $150 purchase
8. Review test results showing 20% discount applied
9. Save rule with DRAFT status
10. Verify rule appears in rules list with correct status

**Test Data**:
- Template: Percentage Discount template
- Rule Name: "Summer Sale - 20% Electronics"
- DSL: `category = 'electronics' AND purchase_amount > 100`
- Test Transaction: {category: 'electronics', amount: 150.00, customer_id: 'CUST001'}

**Expected Results**:
- Rule is created successfully with DRAFT status
- Validation passes without errors
- Test shows $30.00 discount (20% of $150)
- Rule appears in user's rules list
- All operations complete within 5 seconds

---

## TC-FUNC-02: Template Selection and Application

**Test ID**: TC-FUNC-02  
**Priority**: High  
**Prerequisites**: Multiple templates available, user has permissions  

**Test Steps**:
1. Access rule creation interface
2. Browse available templates by category
3. Preview "Loyalty Points" template details
4. Select "Loyalty Points" template
5. Verify template content populates editor
6. Customize points calculation: `points = purchase_amount * 0.1`
7. Modify conditions: `customer_tier = 'gold'`
8. Validate customized rule
9. Save as "Gold Customer Points"

**Test Data**:
- Template: Loyalty Points template
- Customer Tier: 'gold'
- Points Formula: purchase_amount * 0.1
- Test Transaction: {amount: 200.00, customer_tier: 'gold'}

**Expected Results**:
- Template applies successfully
- Customization preserves template structure
- Validation passes with custom content
- Test shows 20 points awarded (200 * 0.1)
- Rule saves with all customizations intact

---

## TC-FUNC-03: DSL Validation and Error Handling

**Test ID**: TC-FUNC-03  
**Priority**: High  
**Prerequisites**: Rule creation interface available, validation engine active  

**Test Steps**:
1. Create new rule without template
2. Enter invalid DSL syntax: `category == electronics AND amount >`
3. Trigger validation
4. Verify specific error messages displayed
5. Correct syntax: `category = 'electronics' AND amount > 100`
6. Re-validate
7. Verify validation passes
8. Test with business logic conflict: duplicate existing rule
9. Verify conflict warning displayed

**Test Data**:
- Invalid DSL: `category == electronics AND amount >`
- Valid DSL: `category = 'electronics' AND amount > 100`
- Existing Rule: Same conditions as test rule

**Expected Results**:
- Invalid syntax shows specific error messages
- Error highlights exact location and suggests correction
- Valid syntax passes validation
- Conflict detection identifies potential duplicate rule
- All validation responses within 2 seconds

---

## TC-FUNC-04: Rule Testing with Sample Data

**Test ID**: TC-FUNC-04  
**Priority**: High  
**Prerequisites**: Rule created and validated, test data available  

**Test Steps**:
1. Open rule testing interface
2. Use predefined test data set: "Standard Purchase"
3. Execute test
4. Review step-by-step execution results
5. Modify test data: increase purchase amount
6. Re-run test
7. Compare results between test runs
8. Save test scenario for future use
9. Test with edge case data: minimum qualifying amount

**Test Data**:
- Standard Test: {category: 'electronics', amount: 150.00}
- Modified Test: {category: 'electronics', amount: 300.00}
- Edge Case: {category: 'electronics', amount: 100.01}

**Expected Results**:
- Test executes successfully with detailed results
- Step-by-step execution shows rule logic flow
- Results change appropriately with different data
- Edge case properly handled at boundary
- Test scenarios save for reuse

---

## TC-FUNC-05: Rule Cloning Functionality

**Test ID**: TC-FUNC-05  
**Priority**: Medium  
**Prerequisites**: Existing approved rule available, cloning permissions  

**Test Steps**:
1. Navigate to existing rules list
2. Select rule "Holiday Discount - 15%"
3. Choose "Clone Rule" option
4. Verify cloned rule opens in editor
5. Modify cloned rule name: "New Year Discount - 25%"
6. Update discount percentage from 15% to 25%
7. Validate modified clone
8. Save clone as draft
9. Verify original rule unchanged
10. Verify clone relationship tracked

**Test Data**:
- Original Rule: "Holiday Discount - 15%"
- Cloned Rule: "New Year Discount - 25%"
- Modification: discount_percentage = 25

**Expected Results**:
- Clone creates independent copy of rule
- Modifications don't affect original rule
- Clone relationship visible in metadata
- Both rules exist independently
- Clone has DRAFT status initially

---

## TC-FUNC-06: Rule History and Audit Trail

**Test ID**: TC-FUNC-06  
**Priority**: Medium  
**Prerequisites**: Rule with multiple status changes, history tracking enabled  

**Test Steps**:
1. Select rule with change history
2. Open rule history view
3. Verify chronological order of changes
4. Review status transition timestamps
5. Check user attribution for each change
6. View detailed change descriptions
7. Compare current version with previous version
8. Export history for audit purposes

**Test Data**:
- Rule with history: Multiple status changes, edits, approvals
- Time period: Last 30 days of changes

**Expected Results**:
- Complete history displayed chronologically
- All status changes properly attributed
- Detailed change descriptions available
- Version comparison shows differences
- History export generates proper audit report

---

## TC-FUNC-07: Draft Rule Management

**Test ID**: TC-FUNC-07  
**Priority**: Medium  
**Prerequisites**: User can create and save draft rules  

**Test Steps**:
1. Create new rule with partial information
2. Save as draft without completing all fields
3. Navigate away from rule editor
4. Return to drafts list
5. Open saved draft
6. Verify all partial information preserved
7. Complete remaining fields
8. Update draft with complete information
9. Change status from DRAFT to UNDER_REVIEW

**Test Data**:
- Partial Rule: Name and category only
- Complete Rule: Full DSL and actions

**Expected Results**:
- Draft saves with incomplete information
- All data preserved when reopening
- Draft can be updated and completed
- Status transition works correctly

---

## TC-FUNC-08: Advanced Validation Scenarios

**Test ID**: TC-FUNC-08  
**Priority**: Medium  
**Prerequisites**: Complex rule scenarios, validation engine with business rules  

**Test Steps**:
1. Create rule with complex nested conditions
2. Validate complex DSL syntax
3. Test rule with multiple customer segments
4. Verify performance impact assessment
5. Check for circular dependency warnings
6. Validate against existing rule conflicts
7. Test with maximum complexity limits
8. Verify optimization suggestions

**Test Data**:
- Complex DSL: Multiple nested AND/OR conditions
- Customer Segments: Gold, Silver, Bronze
- Performance Threshold: <100ms execution time

**Expected Results**:
- Complex rules validate correctly
- Performance impact properly assessed
- Conflicts identified and reported
- Optimization suggestions provided when appropriate

---

## TC-FUNC-09: Template Management and Customization

**Test ID**: TC-FUNC-09  
**Priority**: Low  
**Prerequisites**: Template system active, customization features enabled  

**Test Steps**:
1. Select template for customization
2. Modify template variables
3. Save template customization
4. Create rule from customized template
5. Verify customizations applied
6. Reset template to original state
7. Verify reset removes customizations
8. Test with multiple template variations

**Test Data**:
- Base Template: Standard discount template
- Customizations: Modified variable names and defaults

**Expected Results**:
- Template customizations save properly
- Rules created from customized templates work correctly
- Reset functionality restores original template
- Multiple customizations managed independently

---

## TC-FUNC-10: Rule Status Workflow Integration

**Test ID**: TC-FUNC-10  
**Priority**: High  
**Prerequisites**: Approval workflow active, multiple user roles configured  

**Test Steps**:
1. Create rule and save as DRAFT
2. Submit rule for approval (UNDER_REVIEW)
3. Verify rule cannot be edited in UNDER_REVIEW status
4. Simulate approval process (APPROVED status)
5. Activate rule (ACTIVE status)
6. Test deactivation (INACTIVE status)
7. Test reactivation process
8. Verify status history maintained

**Test Data**:
- Rule workflow: DRAFT → UNDER_REVIEW → APPROVED → ACTIVE
- User roles: Creator, Approver, Administrator

**Expected Results**:
- Status transitions follow defined workflow
- Appropriate permissions enforced at each status
- Rule behavior changes correctly with status
- Complete status history maintained

---

## Performance Requirements

### Response Time Requirements
- Rule creation: <3 seconds
- Validation: <2 seconds
- Testing: <3 seconds
- Save operations: <1 second
- Template loading: <2 seconds

### Throughput Requirements
- Concurrent rule creation: 10 users
- Validation operations: 50 per minute
- Test executions: 100 per minute

### Resource Requirements
- Memory usage: <100MB per active session
- CPU usage: <50% during peak operations
- Database connections: <5 per user session

## Error Handling Requirements

### Validation Errors
- Syntax errors: Specific line and character position
- Semantic errors: Business context and suggestions
- Conflict warnings: Detailed conflict analysis
- Performance warnings: Impact assessment

### System Errors
- Network timeouts: Graceful degradation
- Database errors: Data preservation and recovery
- Service unavailability: Clear user communication
- Concurrent editing: Conflict resolution

### Recovery Procedures
- Auto-save: Every 30 seconds
- Session recovery: Resume after timeout
- Data backup: Automatic on save
- Error logging: Complete audit trail
