# UI/UX Generation Instructions for Rules Engine

**Version:** 1.0.0  
**Last Updated:** 2024-12-19  
**Document Type:** UI/UX Requirements Specification  
**Target Audience:** UI/UX Designers, Frontend Developers, Product Managers  
**Status:** Complete

## Table of Contents

1. [Overview](#overview)
2. [Generation Guidelines](#generation-guidelines)
3. [Business Rule Management Interfaces](#1-business-rule-management-interfaces)
4. [CRUD Interfaces for Database Entities](#2-crud-interfaces-for-database-entities)
5. [Rule Creation Interfaces (Type-Specific)](#3-rule-creation-interfaces-type-specific)
6. [Rule Approval Workflow Interfaces](#4-rule-approval-workflow-interfaces)
7. [Rule Evaluation and Calculation Interfaces](#5-rule-evaluation-and-calculation-interfaces)
8. [Advanced Interface Types](#6-advanced-interface-types)
9. [Mobile-Specific Interfaces](#7-mobile-specific-interfaces)
10. [Integration and API Interfaces](#8-integration-and-api-interfaces)
11. [Cross-Interface Standards](#cross-interface-standards)
12. [Implementation Guidelines](#implementation-guidelines)

---

## Overview

This document provides a comprehensive set of instructions for generating UI/UX specifications for the Rules Engine system. Each instruction follows the established Generate UI/UX rule (@ui-ux.mdc) pattern and is designed to create complete, accessible, and user-centered interface specifications.

### Purpose
- Standardize UI/UX specification generation across all Rules Engine interfaces
- Ensure consistency with DDD domain models and business requirements
- Provide clear, actionable instructions for creating user-centered designs
- Maintain alignment with accessibility and performance standards

### Scope
This document covers all user-facing interfaces for the Rules Engine system, including:
- Management and administrative interfaces
- Rule creation and editing workflows
- Approval and governance interfaces
- Monitoring and analytics dashboards
- Mobile and responsive interfaces

---

## Generation Guidelines

### Mandatory Output Elements
Each generated UI/UX specification must include:

1. **User Flow Diagram** (Mermaid syntax)
   - Main user journey path
   - Alternative and error flows
   - Decision points and branching logic

2. **UI Components and Wireframes** (Text-based)
   - Screen layout descriptions
   - Component specifications
   - Responsive behavior definitions

3. **Interaction and States**
   - UI state definitions (loading, error, success, etc.)
   - User action mappings to domain events
   - Accessibility considerations

### Output Structure
```
07-ui-ux/UI-XXX-interface-name/
├── README.md           # Main interface specification
├── flows.md           # Detailed user flows (if complex)
├── components.md      # Component library additions
└── interactions.md    # Detailed interaction specifications
```

### Domain Event Integration
All user actions must map to appropriate domain events:

```yaml
Rule Management Events:
  - RULE_CREATED
  - RULE_UPDATED
  - RULE_ACTIVATED
  - RULE_DEACTIVATED
  - RULE_DELETED

Approval Events:
  - APPROVAL_REQUESTED
  - RULE_APPROVED
  - RULE_REJECTED
  - APPROVAL_DELEGATED

Evaluation Events:
  - RULE_EVALUATED
  - CONFLICT_DETECTED
  - CONFLICT_RESOLVED
  - CALCULATION_COMPLETED
```

---

## 1. Business Rule Management Interfaces

### 1.1 Primary Rule Management Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the business rule management interface. The interface should allow users to view a list of all existing rules, filter them by status (DRAFT, REVIEW, APPROVED, ACTIVE, INACTIVE, ARCHIVED), priority (LOW, MEDIUM, HIGH, CRITICAL), and category (Promotions, Loyalty, Taxes, Coupons). Include search functionality, bulk actions (activate, deactivate, delete), and individual rule actions (edit, clone, view details, view history). The interface should support pagination for large rule sets and provide quick access to create new rules.
```

**Expected Output:** `07-ui-ux/UI-001-rule-management-interface/`

**Key Features:**
- Rule listing with advanced filtering
- Bulk operations for rule management
- Quick action buttons for common tasks
- Search and pagination capabilities
- Status-based visual indicators

### 1.2 Rule Details Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the rule details interface. The interface should display comprehensive rule information including: rule metadata (name, description, status, priority), DSL content with syntax highlighting, rule performance metrics, version history, approval status, and usage statistics. Include edit capabilities, clone functionality, status change actions, and related rules display. The interface should support rule comparison with previous versions and provide detailed audit trail information.
```

**Expected Output:** `07-ui-ux/UI-002-rule-details-interface/`

**Key Features:**
- Comprehensive rule information display
- Syntax-highlighted DSL editor
- Performance metrics visualization
- Version history and comparison
- Related rules and dependencies

---

## 2. CRUD Interfaces for Database Entities

### 2.1 Rule Template Management Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the rule template management interface. The interface should provide CRUD operations for rule templates including: create new template with name, description, category, DSL template content, and parameter definitions; read/view template details with preview functionality; update existing templates with version control; delete templates with safety checks. Include template categorization (Promotions, Loyalty, Taxes, Coupons), template usage statistics, and template testing capabilities. The interface should support template import/export and template sharing between users.
```

**Expected Output:** `07-ui-ux/UI-003-template-management-interface/`

### 2.2 Customer Segment Management Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the customer segment management interface. The interface should provide CRUD operations for customer segments including: create new segments with criteria definitions; view segment details with member count and growth trends; update segment criteria with impact preview; delete segments with dependency checks. Include segment visualization with charts and metrics, segment testing with sample customers, and segment performance tracking. The interface should support segment hierarchy and nested segment relationships.
```

**Expected Output:** `07-ui-ux/UI-004-customer-segment-interface/`

### 2.3 Product Category Management Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the product category management interface. The interface should provide CRUD operations for product categories including: create hierarchical category structures; view category details with product counts and sales metrics; update category properties and hierarchy; delete categories with product reassignment workflow. Include category tree visualization, category performance analytics, and bulk category operations. The interface should support drag-and-drop hierarchy management and category merchandising features.
```

**Expected Output:** `07-ui-ux/UI-005-product-category-interface/`

### 2.4 Audit Trail Management Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the audit trail management interface. The interface should provide read-only access to audit logs including: view comprehensive audit trails with timestamps, user attribution, and action details; search and filter audit logs by date range, user, action type, and entity; export audit reports for compliance; view audit timeline with visual representation. Include advanced filtering, audit report generation, and compliance dashboard views. The interface should support real-time audit monitoring and alerting for critical events.
```

**Expected Output:** `07-ui-ux/UI-006-audit-trail-interface/`

---

## 3. Rule Creation Interfaces (Type-Specific)

### 3.1 Promotional Rule Creation Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the promotional rule creation interface. The interface should be specifically designed for creating promotional rules including: discount percentage/amount configuration, customer eligibility criteria (segments, tiers, behaviors), product/category targeting, time-based restrictions (date ranges, days of week, time of day), usage limits (per customer, total usage, budget caps), and promotion stacking rules. Include real-time promotion preview, impact estimation calculator, and promotion conflict detection. The interface should provide guided wizard for complex promotions and template-based quick setup for standard promotions.
```

**Expected Output:** `07-ui-ux/UI-007-promotional-rule-creation/`

**Specialized Features:**
- Discount configuration tools
- Customer eligibility builder
- Time-based restriction calendar
- Usage limit controls
- Impact estimation calculator

### 3.2 Loyalty Rule Creation Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the loyalty rule creation interface. The interface should be specifically designed for creating loyalty rules including: points earning configuration (base rate, multipliers, bonus events), tier progression rules, points redemption rules, expiration policies, and special loyalty events. Include loyalty simulation tools, points calculation preview, tier impact analysis, and customer journey visualization. The interface should support complex loyalty scenarios like tiered earning rates, seasonal bonuses, and partnership point transfers.
```

**Expected Output:** `07-ui-ux/UI-008-loyalty-rule-creation/`

**Specialized Features:**
- Points earning rate calculator
- Tier progression builder
- Redemption rule configurator
- Loyalty event scheduler
- Customer journey mapping

### 3.3 Tax Rule Creation Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the tax rule creation interface. The interface should be specifically designed for creating tax rules including: jurisdiction selection (country, state, city), tax rate configuration, tax exemption rules, product category tax mappings, and tax calculation methods (inclusive/exclusive). Include tax jurisdiction mapping tools, tax rate calculator, exemption verification, and compliance validation. The interface should support complex tax scenarios like multi-jurisdictional sales, tax holidays, and special tax categories.
```

**Expected Output:** `07-ui-ux/UI-009-tax-rule-creation/`

**Specialized Features:**
- Jurisdiction mapping interface
- Tax rate calculator
- Exemption rule builder
- Compliance validation tools
- Multi-jurisdiction support

### 3.4 Fee Rule Creation Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the fee rule creation interface. The interface should be specifically designed for creating fee rules including: fee type selection (shipping, processing, service), fee calculation methods (flat rate, percentage, tiered), fee waiver conditions, fee bundling rules, and fee refund policies. Include fee calculator, fee impact analysis, fee optimization suggestions, and fee comparison tools. The interface should support complex fee structures like distance-based shipping, volume discounts, and conditional fee waivers.
```

**Expected Output:** `07-ui-ux/UI-010-fee-rule-creation/`

**Specialized Features:**
- Fee calculation builder
- Waiver condition configurator
- Fee bundling interface
- Impact analysis tools
- Optimization recommendations

### 3.5 Coupon Rule Creation Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the coupon rule creation interface. The interface should be specifically designed for creating coupon rules including: coupon code generation, discount configuration, usage restrictions, distribution channels, expiration settings, and stacking policies. Include coupon performance predictions, fraud prevention settings, coupon batch management, and redemption tracking. The interface should support personalized coupons, dynamic coupon values, and A/B testing of coupon strategies.
```

**Expected Output:** `07-ui-ux/UI-011-coupon-rule-creation/`

**Specialized Features:**
- Coupon code generator
- Distribution channel manager
- Usage restriction builder
- Fraud prevention controls
- A/B testing interface

---

## 4. Rule Approval Workflow Interfaces

### 4.1 Main Approval Workflow Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the rule approval workflow interface. The interface should support multi-level approval processes including: approval request submission with impact analysis, reviewer assignment and notification, approval dashboard with pending/approved/rejected rules, approval history and comments, escalation workflows, and bulk approval capabilities. Include impact assessment visualization, risk analysis reports, compliance checklist validation, and approval workflow configuration. The interface should support different approval paths based on rule type, impact level, and organizational hierarchy.
```

**Expected Output:** `07-ui-ux/UI-012-approval-workflow-interface/`

**Key Features:**
- Multi-level approval dashboard
- Impact assessment visualization
- Workflow configuration tools
- Bulk approval capabilities
- Escalation management

### 4.2 Approval Request Submission Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the approval request submission interface. The interface should allow rule creators to submit rules for approval including: rule summary and impact statement, approval justification, supporting documentation upload, reviewer selection, and priority assignment. Include automated impact analysis, compliance pre-check, stakeholder notification settings, and submission validation. The interface should guide users through the submission process and provide submission status tracking.
```

**Expected Output:** `07-ui-ux/UI-013-approval-submission-interface/`

### 4.3 Approval Review Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the approval review interface. The interface should enable approvers to review and approve rule changes including: rule comparison view (current vs proposed), impact analysis dashboard, compliance verification checklist, risk assessment summary, and approval decision workflow. Include rule testing capabilities, stakeholder feedback collection, approval comments and conditions, and delegation options. The interface should support batch approval for similar rules and approval workflow customization.
```

**Expected Output:** `07-ui-ux/UI-014-approval-review-interface/`

---

## 5. Rule Evaluation and Calculation Interfaces

### 5.1 Rule Testing Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the rule testing interface. The interface should allow users to test rules with sample data including: test scenario builder, sample data input forms, rule execution simulation, result visualization, and performance metrics display. Include predefined test datasets, test result comparison, test scenario saving/sharing, and automated test case generation. The interface should support A/B testing of rule variations and provide detailed execution traces for debugging.
```

**Expected Output:** `07-ui-ux/UI-015-rule-testing-interface/`

**Key Features:**
- Test scenario builder
- Sample data input forms
- Result visualization tools
- Performance metrics display
- Execution trace analyzer

### 5.2 Rule Performance Monitoring Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the rule performance monitoring interface. The interface should provide real-time monitoring of rule evaluation performance including: evaluation response time metrics, throughput statistics, error rate monitoring, conflict resolution tracking, and resource utilization charts. Include performance alerts, trend analysis, capacity planning tools, and performance optimization recommendations. The interface should support performance comparison between rule versions and performance impact analysis of rule changes.
```

**Expected Output:** `07-ui-ux/UI-016-performance-monitoring-interface/`

### 5.3 Rule Calculation Results Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the rule calculation results interface. The interface should display detailed calculation results including: applied rules breakdown, calculation step-by-step process, conflict resolution decisions, final benefit/cost calculations, and calculation audit trail. Include result export capabilities, calculation comparison tools, result visualization charts, and calculation verification features. The interface should support drill-down analysis of complex calculations and provide calculation explanation for business users.
```

**Expected Output:** `07-ui-ux/UI-017-calculation-results-interface/`

### 5.4 Rule Impact Analysis Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the rule impact analysis interface. The interface should provide comprehensive impact analysis including: business impact projections, financial impact calculations, customer impact assessment, system performance impact, and competitive impact analysis. Include impact simulation tools, scenario modeling, impact comparison charts, and impact report generation. The interface should support what-if analysis, sensitivity analysis, and impact trend forecasting.
```

**Expected Output:** `07-ui-ux/UI-018-impact-analysis-interface/`

---

## 6. Advanced Interface Types

### 6.1 Rule Conflict Resolution Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the rule conflict resolution interface. The interface should help users identify and resolve rule conflicts including: conflict detection dashboard, conflict visualization tools, resolution strategy options, conflict impact assessment, and resolution workflow. Include automated conflict detection, manual conflict review, resolution testing, and conflict prevention recommendations. The interface should support complex conflict scenarios and provide guided resolution workflows.
```

**Expected Output:** `07-ui-ux/UI-019-conflict-resolution-interface/`

**Key Features:**
- Conflict detection dashboard
- Visual conflict mapper
- Resolution strategy selector
- Impact assessment tools
- Guided resolution workflow

### 6.2 Rule Analytics and Reporting Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the rule analytics and reporting interface. The interface should provide comprehensive analytics including: rule usage statistics, business impact metrics, performance analytics, user adoption tracking, and trend analysis. Include customizable dashboards, automated report generation, data export capabilities, and interactive visualizations. The interface should support drill-down analysis, comparative reporting, and predictive analytics.
```

**Expected Output:** `07-ui-ux/UI-020-analytics-reporting-interface/`

### 6.3 Rule Version Control Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the rule version control interface. The interface should manage rule versions including: version history display, version comparison tools, rollback capabilities, branching and merging workflows, and version approval tracking. Include visual diff displays, version tagging, release management, and version impact analysis. The interface should support complex versioning scenarios and provide version governance features.
```

**Expected Output:** `07-ui-ux/UI-021-version-control-interface/`

---

## 7. Mobile-Specific Interfaces

### 7.1 Mobile Rule Management Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the mobile rule management interface. The interface should provide mobile-optimized rule management including: touch-friendly rule listing, swipe actions for common operations, mobile-optimized search and filters, quick rule status changes, and mobile notifications. Include offline capabilities, mobile-specific navigation patterns, touch gestures, and mobile performance optimization. The interface should prioritize essential functions for mobile use cases.
```

**Expected Output:** `07-ui-ux/UI-022-mobile-rule-management/`

**Mobile Features:**
- Touch-optimized interactions
- Swipe gestures for actions
- Offline functionality
- Mobile notifications
- Simplified navigation

### 7.2 Mobile Rule Approval Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the mobile rule approval interface. The interface should enable mobile approval workflows including: approval notifications, mobile-friendly rule review, quick approval actions, mobile signatures, and approval delegation. Include push notifications, biometric authentication, offline approval capabilities, and mobile-optimized document viewing. The interface should support urgent approval scenarios and provide mobile approval audit trails.
```

**Expected Output:** `07-ui-ux/UI-023-mobile-approval-interface/`

---

## 8. Integration and API Interfaces

### 8.1 API Management Interface

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the API management interface. The interface should provide API management capabilities including: API endpoint documentation, API key management, rate limiting configuration, API usage monitoring, and API performance analytics. Include API testing tools, developer portal features, API versioning management, and API security monitoring. The interface should support API lifecycle management and developer onboarding.
```

**Expected Output:** `07-ui-ux/UI-024-api-management-interface/`

### 8.2 System Integration Dashboard

**Instruction:**
```
Following the Generate UI/UX rule (@ui-ux.mdc), create the specifications for the system integration dashboard. The interface should monitor system integrations including: integration health monitoring, data synchronization status, error tracking and resolution, integration performance metrics, and connectivity status. Include integration testing tools, data mapping visualization, error alerting, and integration documentation. The interface should support real-time monitoring and integration troubleshooting.
```

**Expected Output:** `07-ui-ux/UI-025-integration-dashboard/`

---

## Cross-Interface Standards

### Design System Requirements

#### Color Palette
```yaml
Primary Colors:
  - Primary Blue: #0066CC (Actions, Links)
  - Success Green: #28A745 (Confirmations, Success)
  - Warning Orange: #FFC107 (Warnings, Attention)
  - Error Red: #DC3545 (Errors, Deletion)
  - Info Blue: #17A2B8 (Information, Neutral)

Neutral Colors:
  - Background: #F8F9FA
  - Surface: #FFFFFF
  - Border: #DEE2E6
  - Text Primary: #212529
  - Text Secondary: #6C757D
```

#### Typography
```yaml
Font Family: System font stack (San Francisco, Roboto, etc.)
Font Sizes:
  - Display: 32px (Major headings)
  - Heading 1: 24px (Page titles)
  - Heading 2: 20px (Section headers)
  - Heading 3: 16px (Subsection headers)
  - Body: 14px (Regular text)
  - Caption: 12px (Labels, captions)
```

#### Spacing System
```yaml
Base Unit: 8px
Spacing Scale:
  - xs: 4px
  - sm: 8px
  - md: 16px
  - lg: 24px
  - xl: 32px
  - xxl: 48px
```

### Component Standards

#### Common Components
- **Buttons**: Primary, Secondary, Tertiary, Icon buttons
- **Form Controls**: Input fields, Select dropdowns, Checkboxes, Radio buttons
- **Data Display**: Tables, Cards, Lists, Tags, Badges
- **Navigation**: Breadcrumbs, Tabs, Pagination, Sidebar
- **Feedback**: Alerts, Toasts, Progress indicators, Loading states
- **Overlays**: Modals, Popovers, Tooltips, Drawers

#### Component States
```yaml
Interactive States:
  - Default: Base appearance
  - Hover: Interactive feedback
  - Active: Currently pressed/selected
  - Focus: Keyboard focus indicator
  - Disabled: Non-interactive state
  - Loading: Processing state
  - Error: Error state with validation
```

### Accessibility Standards

#### WCAG 2.1 AA Compliance
- **Color Contrast**: Minimum 4.5:1 for normal text, 3:1 for large text
- **Focus Management**: Visible focus indicators, logical tab order
- **Keyboard Navigation**: Full keyboard accessibility
- **Screen Reader Support**: Proper ARIA labels and descriptions
- **Text Alternatives**: Alt text for images, descriptive link text

#### Assistive Technology Support
- **Screen Readers**: NVDA, JAWS, VoiceOver compatibility
- **Voice Control**: Voice command support where applicable
- **Magnification**: High DPI and zoom support up to 200%
- **Motor Impairments**: Large touch targets (minimum 44px)

### Performance Standards

#### Loading Performance
```yaml
Performance Targets:
  - First Contentful Paint: < 1.5s
  - Largest Contentful Paint: < 2.5s
  - Cumulative Layout Shift: < 0.1
  - First Input Delay: < 100ms
  
Loading Strategies:
  - Skeleton screens for content loading
  - Progressive loading for large datasets
  - Lazy loading for images and non-critical content
  - Optimistic UI updates where appropriate
```

#### Responsive Breakpoints
```yaml
Breakpoints:
  - Mobile: < 768px
  - Tablet: 768px - 1024px
  - Desktop: > 1024px
  
Design Approach:
  - Mobile-first responsive design
  - Progressive enhancement
  - Touch-friendly interactions
  - Flexible grid systems
```

---

## Implementation Guidelines

### Development Workflow

#### Phase 1: Design System Foundation
1. Establish design tokens and component library
2. Create base components with all states
3. Implement accessibility features
4. Set up responsive grid system

#### Phase 2: Core Interface Development
1. Implement rule management interfaces
2. Develop CRUD interfaces for entities
3. Create rule creation workflows
4. Build approval workflow interfaces

#### Phase 3: Advanced Features
1. Implement evaluation and calculation interfaces
2. Develop analytics and reporting
3. Create mobile-specific interfaces
4. Build integration dashboards

#### Phase 4: Optimization and Enhancement
1. Performance optimization
2. Advanced accessibility features
3. User experience refinements
4. Integration testing and validation

### Quality Assurance

#### Testing Requirements
- **Unit Testing**: Component functionality and state management
- **Integration Testing**: Interface workflows and data flow
- **Accessibility Testing**: WCAG compliance and assistive technology
- **Performance Testing**: Load times and responsiveness
- **Usability Testing**: User experience validation

#### Review Process
1. **Design Review**: Visual design and interaction patterns
2. **Technical Review**: Implementation feasibility and standards
3. **Accessibility Review**: WCAG compliance verification
4. **User Experience Review**: Usability and user journey validation

### Documentation Standards

#### Interface Documentation
Each interface specification must include:
- **Purpose and Context**: Why the interface exists and how it fits
- **User Stories**: Who uses it and what they need to accomplish
- **User Flow Diagrams**: Visual representation of user journeys
- **Component Specifications**: Detailed component descriptions
- **Interaction Patterns**: How users interact with the interface
- **Accessibility Notes**: Specific accessibility considerations
- **Performance Requirements**: Loading and response time expectations

#### Design System Documentation
- **Component Library**: Detailed component specifications
- **Design Tokens**: Color, typography, spacing definitions
- **Usage Guidelines**: When and how to use components
- **Accessibility Guidelines**: Component-specific accessibility notes
- **Code Examples**: Implementation examples and best practices

---

## Conclusion

This comprehensive set of UI/UX generation instructions provides the foundation for creating consistent, accessible, and user-centered interfaces for the Rules Engine system. Each instruction is designed to generate complete specifications that align with domain requirements, accessibility standards, and performance expectations.

### Key Benefits
- **Consistency**: Standardized approach across all interfaces
- **Accessibility**: Built-in WCAG 2.1 AA compliance
- **Performance**: Optimized for speed and responsiveness
- **Maintainability**: Clear documentation and design system
- **Scalability**: Flexible architecture for future enhancements

### Next Steps
1. Begin with core rule management interfaces
2. Establish design system foundation
3. Implement CRUD interfaces for data entities
4. Develop specialized rule creation workflows
5. Build approval and governance interfaces
6. Create evaluation and monitoring dashboards

By following these instructions and guidelines, the development team can create a comprehensive, user-friendly Rules Engine interface that meets all business requirements while maintaining the highest standards of usability and accessibility.
