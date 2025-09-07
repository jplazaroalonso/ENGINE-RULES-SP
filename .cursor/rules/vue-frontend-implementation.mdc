---
description: Vue Frontend Implementation Rule
globs:
alwaysApply: false
---

# Vue Frontend Implementation Rule

## Overview

This rule provides comprehensive guidance for implementing enterprise-grade Vue 3 frontend applications with TypeScript, component libraries, and comprehensive testing. Based on successful Rules Engine frontend implementation achieving production-ready user experience.

## 1. Feedback Section

### Identified Issues in Generic Frontend Development:
- **Inconsistent Component Design**: Duplicated components across different screens
- **Poor State Management**: Direct API calls in components without proper caching
- **Missing Accessibility**: No keyboard navigation, screen reader support
- **Weak Type Safety**: JavaScript instead of TypeScript, no interface definitions
- **Limited Testing**: Missing behavioral tests and component integration tests
- **No Design System**: Inconsistent styling and component behavior

### Recommendations:
- Create shared component library with consistent design patterns
- Implement centralized state management with Pinia and proper caching
- Add comprehensive accessibility features following WCAG guidelines
- Use TypeScript throughout with strict typing and interface definitions
- Implement comprehensive testing pyramid with behavioral and visual tests
- Design component system with clear patterns and documentation

## 2. Role and Context Definition

### Target Role: Frontend Developer (Vue/TypeScript)
### Background Context:
- **Framework**: Vue 3 with Composition API and script setup syntax
- **Type System**: TypeScript with strict mode and interface definitions
- **UI Library**: Quasar Framework with Material Design components
- **State Management**: Pinia with reactive stores and computed properties
- **Testing**: Vitest (unit), Cypress (E2E), Storybook (component documentation)
- **Build System**: Vite with optimized production builds

## 3. Objective and Goals

### Primary Objective:
Create production-ready Vue 3 frontend application with detailed screen specifications, reusable component library, comprehensive testing, and enterprise-grade user experience.

### Success Criteria:
- **Component Library**: Reusable components with consistent design patterns
- **Screen Implementation**: Complete screen specifications with user workflows
- **Type Safety**: Full TypeScript coverage with strict typing
- **Accessibility**: WCAG 2.1 AA compliance with keyboard navigation
- **Performance**: <3s initial load, 60fps interactions, optimized bundle size
- **Testing**: >80% coverage with unit, integration, and E2E tests

## 4. Key Terms and Definitions

### Technical Terminology:
- **Composition API**: Vue 3's reactive composition system for reusable logic
- **Pinia Store**: Centralized state management with type-safe reactive stores
- **Component Library**: Shared collection of reusable UI components
- **Design Tokens**: Standardized design values (colors, spacing, typography)
- **Props Interface**: TypeScript interface defining component input properties
- **Composable**: Reusable Vue composition function encapsulating logic
- **Emit Events**: Type-safe component event system for parent communication

## 5. Task Decomposition (Chain-of-Thought)

### Step 1: Project Architecture Setup
- **Input**: Application requirements, screen specifications, design system
- **Process**: Setup Vue 3 project with TypeScript, Vite, and component structure
- **Output**: Complete project skeleton with shared component library foundation
- **Human Validation Point**: Review project structure follows Vue 3 best practices

### Step 2: Shared Component Library Development
- **Input**: Design tokens, component specifications, accessibility requirements
- **Process**: Create reusable components with consistent design patterns
- **Output**: Component library with Storybook documentation and testing
- **Human Validation Point**: Verify components meet design system and accessibility standards

### Step 3: Screen Implementation
- **Input**: Screen specifications, user workflows, API contracts
- **Process**: Implement complete screens using shared components and state management
- **Output**: Functional screens with proper navigation and user interactions
- **Human Validation Point**: Confirm screens match specifications and handle edge cases

### Step 4: State Management Implementation
- **Input**: Data requirements, API endpoints, caching strategies
- **Process**: Create Pinia stores with reactive data and computed properties
- **Output**: Centralized state management with proper error handling
- **Human Validation Point**: Validate stores handle all data scenarios correctly

### Step 5: Testing Implementation
- **Input**: Component behavior, user workflows, accessibility requirements
- **Process**: Create comprehensive test suite covering all interaction patterns
- **Output**: Complete test coverage with unit, integration, and E2E tests
- **Human Validation Point**: Confirm tests cover all user scenarios and edge cases

### Step 6: Performance Optimization and Deployment
- **Input**: Performance requirements, bundle size targets, deployment specs
- **Process**: Optimize application performance and create deployment artifacts
- **Output**: Production-ready application with monitoring and deployment setup
- **Human Validation Point**: Verify performance metrics meet requirements

## 6. Context and Constraints

### Technical Context:
- **Vue Version**: 3.4+ with Composition API and script setup
- **TypeScript**: 5.0+ with strict mode and interface definitions
- **Build Tool**: Vite 5.0+ with optimized production builds
- **UI Framework**: Quasar 2.14+ with Material Design components
- **Testing**: Vitest for unit tests, Cypress for E2E testing

### Business Context:
- **User Experience**: Intuitive interfaces for business rule management
- **Performance**: <3s initial load time, 60fps smooth interactions
- **Accessibility**: WCAG 2.1 AA compliance for inclusive usage
- **Responsive**: Mobile-first design with progressive enhancement
- **Browser Support**: Modern browsers (Chrome 90+, Firefox 88+, Safari 14+)

### Negative Constraints:
- **Do NOT** use Options API - use Composition API exclusively
- **Do NOT** create components without TypeScript interfaces
- **Do NOT** skip accessibility features and keyboard navigation
- **Do NOT** implement custom UI components without design system
- **Do NOT** ignore performance optimization and bundle analysis

## 7. Examples and Illustrations (Few-Shot)

### Example 1: Shared Component Library Structure

#### Project Structure:
```
shared-components/
├── src/
│   ├── components/
│   │   ├── core/                    # Universal components
│   │   │   ├── BaseButton/
│   │   │   │   ├── BaseButton.vue
│   │   │   │   ├── BaseButton.stories.ts
│   │   │   │   ├── BaseButton.test.ts
│   │   │   │   └── types.ts
│   │   │   ├── BaseInput/
│   │   │   │   ├── BaseInput.vue
│   │   │   │   ├── BaseInput.stories.ts
│   │   │   │   ├── BaseInput.test.ts
│   │   │   │   └── types.ts
│   │   │   └── BaseSelect/
│   │   ├── entities/                # Business-specific components
│   │   │   ├── RuleCard/
│   │   │   │   ├── RuleCard.vue
│   │   │   │   ├── RuleCard.stories.ts
│   │   │   │   ├── RuleCard.test.ts
│   │   │   │   └── types.ts
│   │   │   ├── CampaignCard/
│   │   │   └── CustomerCard/
│   │   ├── forms/                   # Form components
│   │   │   ├── DynamicForm/
│   │   │   ├── DSLEditor/
│   │   │   └── ValidationSummary/
│   │   ├── data/                    # Data display components
│   │   │   ├── DataTable/
│   │   │   ├── DataChart/
│   │   │   └── DataExporter/
│   │   └── layout/                  # Layout components
│   │       ├── AppLayout/
│   │       ├── PageHeader/
│   │       └── NavigationMenu/
│   ├── composables/                 # Reusable logic
│   │   ├── useApi.ts
│   │   ├── useForm.ts
│   │   ├── useTable.ts
│   │   ├── usePagination.ts
│   │   └── useNotification.ts
│   ├── stores/                      # Pinia stores
│   │   ├── auth.ts
│   │   ├── rules.ts
│   │   ├── campaigns.ts
│   │   └── ui.ts
│   ├── types/                       # TypeScript definitions
│   │   ├── api.ts
│   │   ├── entities.ts
│   │   ├── forms.ts
│   │   └── ui.ts
│   ├── utils/                       # Utility functions
│   │   ├── formatters.ts
│   │   ├── validators.ts
│   │   └── constants.ts
│   └── styles/                      # Design system
│       ├── tokens.scss
│       ├── components.scss
│       └── utilities.scss
├── stories/                         # Storybook configuration
├── tests/                          # Test utilities
└── docs/                           # Component documentation
```

#### Base Button Component Example:
```vue
<!-- src/components/core/BaseButton/BaseButton.vue -->
<template>
  <q-btn
    :class="buttonClasses"
    :color="computedColor"
    :size="size"
    :flat="variant === 'text'"
    :outlined="variant === 'outlined'"
    :unelevated="variant === 'filled'"
    :disable="disabled || loading"
    :loading="loading"
    :icon="icon"
    :icon-right="iconRight"
    :to="to"
    :href="href"
    :type="type"
    @click="handleClick"
    v-bind="$attrs"
  >
    <slot>{{ label }}</slot>
  </q-btn>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { ButtonProps, ButtonEmits } from './types'

// Define props with default values
const props = withDefaults(defineProps<ButtonProps>(), {
  variant: 'filled',
  size: 'md',
  color: 'primary',
  type: 'button',
  disabled: false,
  loading: false,
  fullWidth: false
})

// Define emits
const emit = defineEmits<ButtonEmits>()

// Computed properties
const buttonClasses = computed(() => ({
  'base-button': true,
  'base-button--full-width': props.fullWidth,
  [`base-button--${props.variant}`]: true,
  [`base-button--${props.size}`]: true
}))

const computedColor = computed(() => {
  // Map custom colors to Quasar colors if needed
  const colorMap: Record<string, string> = {
    danger: 'negative',
    success: 'positive',
    warning: 'orange'
  }
  return colorMap[props.color] || props.color
})

// Event handlers
const handleClick = (event: MouseEvent) => {
  if (!props.disabled && !props.loading) {
    emit('click', event)
  }
}
</script>

<style lang="scss" scoped>
.base-button {
  &--full-width {
    width: 100%;
  }
  
  // Custom variant styles
  &--filled {
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }
  
  &--outlined {
    border-width: 2px;
  }
  
  &--text {
    text-transform: none;
  }
  
  // Size variants
  &--xs {
    font-size: 0.75rem;
    padding: 4px 8px;
  }
  
  &--sm {
    font-size: 0.875rem;
    padding: 6px 12px;
  }
  
  &--md {
    font-size: 1rem;
    padding: 8px 16px;
  }
  
  &--lg {
    font-size: 1.125rem;
    padding: 12px 24px;
  }
  
  &--xl {
    font-size: 1.25rem;
    padding: 16px 32px;
  }
}
</style>
```

#### Button Types Definition:
```typescript
// src/components/core/BaseButton/types.ts
import type { RouteLocationRaw } from 'vue-router'

export interface ButtonProps {
  /** Button text content */
  label?: string
  
  /** Button visual variant */
  variant?: 'filled' | 'outlined' | 'text'
  
  /** Button size */
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl'
  
  /** Button color theme */
  color?: 'primary' | 'secondary' | 'accent' | 'positive' | 'negative' | 'warning' | 'info' | 'dark' | 'light'
  
  /** Left icon name */
  icon?: string
  
  /** Right icon name */
  iconRight?: string
  
  /** Button HTML type */
  type?: 'button' | 'submit' | 'reset'
  
  /** Disabled state */
  disabled?: boolean
  
  /** Loading state with spinner */
  loading?: boolean
  
  /** Full width button */
  fullWidth?: boolean
  
  /** Vue Router route object */
  to?: RouteLocationRaw
  
  /** External href */
  href?: string
  
  /** Target for href */
  target?: string
}

export interface ButtonEmits {
  /** Emitted when button is clicked */
  click: [event: MouseEvent]
}

export interface ButtonSlots {
  /** Default slot for button content */
  default(): any
}
```

#### Button Storybook Stories:
```typescript
// src/components/core/BaseButton/BaseButton.stories.ts
import type { Meta, StoryObj } from '@storybook/vue3'
import BaseButton from './BaseButton.vue'

const meta: Meta<typeof BaseButton> = {
  title: 'Core/BaseButton',
  component: BaseButton,
  parameters: {
    docs: {
      description: {
        component: 'A versatile button component with multiple variants, sizes, and states.'
      }
    }
  },
  argTypes: {
    variant: {
      control: 'select',
      options: ['filled', 'outlined', 'text']
    },
    size: {
      control: 'select',
      options: ['xs', 'sm', 'md', 'lg', 'xl']
    },
    color: {
      control: 'select',
      options: ['primary', 'secondary', 'accent', 'positive', 'negative', 'warning', 'info']
    },
    onClick: { action: 'clicked' }
  }
}

export default meta
type Story = StoryObj<typeof BaseButton>

export const Default: Story = {
  args: {
    label: 'Button'
  }
}

export const Variants: Story = {
  render: () => ({
    components: { BaseButton },
    template: `
      <div class="flex gap-4">
        <BaseButton label="Filled" variant="filled" />
        <BaseButton label="Outlined" variant="outlined" />
        <BaseButton label="Text" variant="text" />
      </div>
    `
  })
}

export const Sizes: Story = {
  render: () => ({
    components: { BaseButton },
    template: `
      <div class="flex items-center gap-4">
        <BaseButton label="Extra Small" size="xs" />
        <BaseButton label="Small" size="sm" />
        <BaseButton label="Medium" size="md" />
        <BaseButton label="Large" size="lg" />
        <BaseButton label="Extra Large" size="xl" />
      </div>
    `
  })
}

export const Colors: Story = {
  render: () => ({
    components: { BaseButton },
    template: `
      <div class="flex flex-wrap gap-2">
        <BaseButton label="Primary" color="primary" />
        <BaseButton label="Secondary" color="secondary" />
        <BaseButton label="Accent" color="accent" />
        <BaseButton label="Positive" color="positive" />
        <BaseButton label="Negative" color="negative" />
        <BaseButton label="Warning" color="warning" />
        <BaseButton label="Info" color="info" />
      </div>
    `
  })
}

export const States: Story = {
  render: () => ({
    components: { BaseButton },
    template: `
      <div class="flex gap-4">
        <BaseButton label="Normal" />
        <BaseButton label="Loading" :loading="true" />
        <BaseButton label="Disabled" :disabled="true" />
      </div>
    `
  })
}

export const WithIcons: Story = {
  render: () => ({
    components: { BaseButton },
    template: `
      <div class="flex gap-4">
        <BaseButton label="Save" icon="save" />
        <BaseButton label="Delete" icon="delete" color="negative" />
        <BaseButton label="Next" icon-right="arrow_forward" />
        <BaseButton icon="search" />
      </div>
    `
  })
}

export const Interactive: Story = {
  args: {
    label: 'Click me!'
  },
  play: async ({ args, canvasElement }) => {
    // Add interaction tests here
  }
}
```

### Example 2: Rule Card Component Implementation

```vue
<!-- src/components/entities/RuleCard/RuleCard.vue -->
<template>
  <q-card 
    class="rule-card" 
    :class="cardClasses"
    @click="handleCardClick"
    :tabindex="interactive ? 0 : -1"
    @keydown.enter="handleCardClick"
    @keydown.space.prevent="handleCardClick"
  >
    <q-card-section>
      <!-- Header with title and actions -->
      <div class="rule-card__header">
        <div class="rule-card__title-section">
          <h3 class="rule-card__title">{{ rule.name }}</h3>
          <div class="rule-card__metadata">
            <StatusChip 
              :status="rule.status" 
              :size="compact ? 'sm' : 'md'"
              class="rule-card__status"
            />
            <PriorityChip 
              :priority="rule.priority"
              :size="compact ? 'sm' : 'md'" 
              class="rule-card__priority"
            />
            <VersionBadge 
              :version="rule.version"
              class="rule-card__version"
            />
          </div>
        </div>
        
        <div class="rule-card__actions" v-if="showActions">
          <BaseButton
            icon="visibility"
            variant="text"
            size="sm"
            @click.stop="emit('view', rule)"
            :aria-label="`View ${rule.name} details`"
          />
          
          <BaseButton
            icon="edit"
            variant="text"
            size="sm"
            @click.stop="emit('edit', rule)"
            v-if="canEdit"
            :aria-label="`Edit ${rule.name}`"
          />
          
          <q-btn-dropdown
            icon="more_vert"
            flat
            dense
            auto-close
            @click.stop
          >
            <q-list>
              <q-item clickable @click="emit('duplicate', rule)">
                <q-item-section avatar>
                  <q-icon name="content_copy" />
                </q-item-section>
                <q-item-section>Duplicate</q-item-section>
              </q-item>
              
              <q-item clickable @click="emit('test', rule)">
                <q-item-section avatar>
                  <q-icon name="play_arrow" />
                </q-item-section>
                <q-item-section>Test Rule</q-item-section>
              </q-item>
              
              <q-separator />
              
              <q-item 
                clickable 
                @click="emit('delete', rule)"
                v-if="canDelete"
                class="text-negative"
              >
                <q-item-section avatar>
                  <q-icon name="delete" />
                </q-item-section>
                <q-item-section>Delete</q-item-section>
              </q-item>
            </q-list>
          </q-btn-dropdown>
        </div>
      </div>
      
      <!-- Description -->
      <div class="rule-card__description" v-if="rule.description && !compact">
        <p>{{ truncatedDescription }}</p>
        <BaseButton
          v-if="isDescriptionTruncated"
          variant="text"
          size="sm"
          :label="showFullDescription ? 'Show less' : 'Show more'"
          @click.stop="showFullDescription = !showFullDescription"
        />
      </div>
      
      <!-- Tags -->
      <div class="rule-card__tags" v-if="rule.tags?.length">
        <TagChip
          v-for="tag in displayedTags"
          :key="tag"
          :tag="tag"
          size="sm"
        />
        <TagChip
          v-if="hasMoreTags"
          :tag="`+${rule.tags.length - maxTags}`"
          variant="more"
          size="sm"
        />
      </div>
      
      <!-- Footer with metadata -->
      <div class="rule-card__footer">
        <div class="rule-card__author">
          <UserAvatar 
            :user="rule.createdBy"
            size="sm"
            show-name
          />
          <TimeAgo 
            :timestamp="rule.createdAt"
            class="rule-card__timestamp"
          />
        </div>
        
        <!-- Performance metrics -->
        <div class="rule-card__metrics" v-if="metrics && showMetrics">
          <PerformanceIndicator
            :success-rate="metrics.successRate"
            :execution-count="metrics.executionCount"
            :avg-execution-time="metrics.avgExecutionTime"
            compact
          />
        </div>
      </div>
    </q-card-section>
    
    <!-- Loading overlay -->
    <q-inner-loading :showing="loading">
      <q-spinner-dots size="40px" color="primary" />
    </q-inner-loading>
  </q-card>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { truncateText } from '@/utils/formatters'
import type { RuleCardProps, RuleCardEmits } from './types'

// Components
import StatusChip from '@/components/core/StatusChip/StatusChip.vue'
import PriorityChip from '@/components/core/PriorityChip/PriorityChip.vue'
import VersionBadge from '@/components/core/VersionBadge/VersionBadge.vue'
import TagChip from '@/components/core/TagChip/TagChip.vue'
import UserAvatar from '@/components/core/UserAvatar/UserAvatar.vue'
import TimeAgo from '@/components/core/TimeAgo/TimeAgo.vue'
import PerformanceIndicator from '@/components/core/PerformanceIndicator/PerformanceIndicator.vue'
import BaseButton from '@/components/core/BaseButton/BaseButton.vue'

// Props with defaults
const props = withDefaults(defineProps<RuleCardProps>(), {
  compact: false,
  interactive: true,
  showActions: true,
  showMetrics: true,
  canEdit: true,
  canDelete: true,
  loading: false,
  maxTags: 3,
  maxDescriptionLength: 120
})

// Emits
const emit = defineEmits<RuleCardEmits>()

// Local state
const showFullDescription = ref(false)

// Computed properties
const cardClasses = computed(() => ({
  'rule-card--compact': props.compact,
  'rule-card--interactive': props.interactive,
  'rule-card--loading': props.loading,
  [`rule-card--status-${props.rule.status.toLowerCase()}`]: true,
  [`rule-card--priority-${props.rule.priority.toLowerCase()}`]: true
}))

const truncatedDescription = computed(() => {
  if (!props.rule.description) return ''
  if (showFullDescription.value) return props.rule.description
  return truncateText(props.rule.description, props.maxDescriptionLength)
})

const isDescriptionTruncated = computed(() => {
  return props.rule.description && 
         props.rule.description.length > props.maxDescriptionLength
})

const displayedTags = computed(() => {
  if (!props.rule.tags) return []
  return props.rule.tags.slice(0, props.maxTags)
})

const hasMoreTags = computed(() => {
  return props.rule.tags && props.rule.tags.length > props.maxTags
})

// Event handlers
const handleCardClick = () => {
  if (props.interactive && !props.loading) {
    emit('click', props.rule)
  }
}
</script>

<style lang="scss" scoped>
.rule-card {
  transition: all 0.3s ease;
  border-radius: 12px;
  position: relative;
  
  &--interactive {
    cursor: pointer;
    
    &:hover {
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
      transform: translateY(-2px);
    }
    
    &:focus {
      outline: 2px solid var(--q-primary);
      outline-offset: 2px;
    }
  }
  
  &--compact {
    .rule-card__title {
      font-size: 1rem;
    }
    
    .rule-card__description {
      display: none;
    }
  }
  
  &--loading {
    pointer-events: none;
  }
  
  // Status-based styling
  &--status-active {
    border-left: 4px solid var(--q-positive);
  }
  
  &--status-inactive {
    border-left: 4px solid var(--q-warning);
  }
  
  &--status-draft {
    border-left: 4px solid var(--q-grey);
  }
  
  &--status-deprecated {
    border-left: 4px solid var(--q-negative);
    opacity: 0.8;
  }
  
  &__header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 16px;
  }
  
  &__title-section {
    flex: 1;
    min-width: 0; // Allow text truncation
  }
  
  &__title {
    margin: 0 0 8px 0;
    font-size: 1.1rem;
    font-weight: 600;
    line-height: 1.3;
    color: var(--q-dark);
    
    // Text truncation for long titles
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  
  &__metadata {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
    align-items: center;
  }
  
  &__actions {
    display: flex;
    gap: 4px;
    margin-left: 16px;
    flex-shrink: 0;
  }
  
  &__description {
    margin-bottom: 16px;
    
    p {
      margin: 0;
      color: var(--q-grey-7);
      line-height: 1.4;
    }
  }
  
  &__tags {
    display: flex;
    gap: 6px;
    flex-wrap: wrap;
    margin-bottom: 16px;
  }
  
  &__footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 16px;
  }
  
  &__author {
    display: flex;
    align-items: center;
    gap: 8px;
    flex: 1;
    min-width: 0;
  }
  
  &__timestamp {
    color: var(--q-grey-6);
    font-size: 0.875rem;
  }
  
  &__metrics {
    flex-shrink: 0;
  }
}

// Responsive design
@media (max-width: 600px) {
  .rule-card {
    &__header {
      flex-direction: column;
      align-items: stretch;
      gap: 12px;
    }
    
    &__actions {
      margin-left: 0;
      justify-content: flex-end;
    }
    
    &__footer {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;
    }
  }
}

// High contrast mode support
@media (prefers-contrast: high) {
  .rule-card {
    border: 2px solid var(--q-dark);
    
    &--status-active {
      border-left: 6px solid var(--q-positive);
    }
    
    &--status-inactive {
      border-left: 6px solid var(--q-warning);
    }
    
    &--status-draft {
      border-left: 6px solid var(--q-grey);
    }
    
    &--status-deprecated {
      border-left: 6px solid var(--q-negative);
    }
  }
}

// Reduced motion support
@media (prefers-reduced-motion: reduce) {
  .rule-card {
    transition: none;
    
    &--interactive:hover {
      transform: none;
    }
  }
}
</style>
```

### Example 3: Pinia Store Implementation

```typescript
// src/stores/rules.ts
import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import { rulesApi } from '@/api/rules'
import { useNotificationStore } from './notifications'
import type { 
  Rule, 
  RuleListParams, 
  RuleFilters,
  CreateRuleRequest,
  UpdateRuleRequest,
  RuleMetrics 
} from '@/types/rules'

export const useRulesStore = defineStore('rules', () => {
  // State
  const rules = ref<Rule[]>([])
  const currentRule = ref<Rule | null>(null)
  const loading = ref(false)
  const saving = ref(false)
  const error = ref<string | null>(null)
  
  // Pagination
  const pagination = ref({
    page: 1,
    rowsPerPage: 20,
    rowsNumber: 0,
    sortBy: 'createdAt',
    descending: true
  })
  
  // Filters
  const filters = ref<RuleFilters>({
    status: [],
    priority: [],
    category: '',
    createdBy: '',
    dateRange: null,
    tags: []
  })
  
  // Search
  const searchQuery = ref('')
  
  // Cache for rule metrics
  const ruleMetrics = ref<Map<string, RuleMetrics>>(new Map())
  
  // Notification store
  const notificationStore = useNotificationStore()
  
  // Getters
  const filteredRules = computed(() => {
    let result = rules.value
    
    // Apply search filter
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      result = result.filter(rule => 
        rule.name.toLowerCase().includes(query) ||
        rule.description.toLowerCase().includes(query) ||
        rule.tags?.some(tag => tag.toLowerCase().includes(query))
      )
    }
    
    // Apply status filter
    if (filters.value.status.length > 0) {
      result = result.filter(rule => filters.value.status.includes(rule.status))
    }
    
    // Apply priority filter
    if (filters.value.priority.length > 0) {
      result = result.filter(rule => filters.value.priority.includes(rule.priority))
    }
    
    // Apply category filter
    if (filters.value.category) {
      result = result.filter(rule => rule.category === filters.value.category)
    }
    
    // Apply created by filter
    if (filters.value.createdBy) {
      result = result.filter(rule => 
        rule.createdBy.toLowerCase().includes(filters.value.createdBy.toLowerCase())
      )
    }
    
    // Apply date range filter
    if (filters.value.dateRange) {
      const { start, end } = filters.value.dateRange
      result = result.filter(rule => {
        const createdDate = new Date(rule.createdAt)
        return createdDate >= start && createdDate <= end
      })
    }
    
    // Apply tags filter
    if (filters.value.tags.length > 0) {
      result = result.filter(rule => 
        rule.tags?.some(tag => filters.value.tags.includes(tag))
      )
    }
    
    return result
  })
  
  const rulesByStatus = computed(() => {
    return (status: Rule['status']) => 
      rules.value.filter(rule => rule.status === status)
  })
  
  const rulesByPriority = computed(() => {
    return (priority: Rule['priority']) =>
      rules.value.filter(rule => rule.priority === priority)
  })
  
  const rulesStats = computed(() => {
    const stats = {
      total: rules.value.length,
      active: 0,
      draft: 0,
      underReview: 0,
      deprecated: 0
    }
    
    rules.value.forEach(rule => {
      switch (rule.status) {
        case 'ACTIVE':
          stats.active++
          break
        case 'DRAFT':
          stats.draft++
          break
        case 'UNDER_REVIEW':
          stats.underReview++
          break
        case 'DEPRECATED':
          stats.deprecated++
          break
      }
    })
    
    return stats
  })
  
  // Actions
  const fetchRules = async (params?: Partial<RuleListParams>) => {
    loading.value = true
    error.value = null
    
    try {
      const requestParams: RuleListParams = {
        page: pagination.value.page,
        limit: pagination.value.rowsPerPage,
        sortBy: pagination.value.sortBy,
        sortOrder: pagination.value.descending ? 'desc' : 'asc',
        filters: filters.value,
        search: searchQuery.value,
        ...params
      }
      
      const response = await rulesApi.getRules(requestParams)
      
      rules.value = response.data
      pagination.value.rowsNumber = response.total
      pagination.value.page = response.page
      
      return response
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to fetch rules'
      error.value = message
      notificationStore.showError('Failed to load rules', message)
      throw err
    } finally {
      loading.value = false
    }
  }
  
  const fetchRule = async (id: string, useCache = true) => {
    // Check cache first
    if (useCache) {
      const cachedRule = rules.value.find(rule => rule.id === id)
      if (cachedRule) {
        currentRule.value = cachedRule
        return cachedRule
      }
    }
    
    loading.value = true
    error.value = null
    
    try {
      const response = await rulesApi.getRule(id)
      currentRule.value = response.data
      
      // Update cache
      const index = rules.value.findIndex(rule => rule.id === id)
      if (index !== -1) {
        rules.value[index] = response.data
      }
      
      return response.data
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to fetch rule'
      error.value = message
      notificationStore.showError('Failed to load rule', message)
      throw err
    } finally {
      loading.value = false
    }
  }
  
  const createRule = async (request: CreateRuleRequest) => {
    saving.value = true
    error.value = null
    
    try {
      const response = await rulesApi.createRule(request)
      
      // Add to local state
      rules.value.unshift(response.data)
      currentRule.value = response.data
      
      notificationStore.showSuccess('Rule created successfully')
      return response.data
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to create rule'
      error.value = message
      notificationStore.showError('Failed to create rule', message)
      throw err
    } finally {
      saving.value = false
    }
  }
  
  const updateRule = async (id: string, request: UpdateRuleRequest) => {
    saving.value = true
    error.value = null
    
    try {
      const response = await rulesApi.updateRule(id, request)
      
      // Update local state
      const index = rules.value.findIndex(rule => rule.id === id)
      if (index !== -1) {
        rules.value[index] = response.data
      }
      
      if (currentRule.value?.id === id) {
        currentRule.value = response.data
      }
      
      notificationStore.showSuccess('Rule updated successfully')
      return response.data
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to update rule'
      error.value = message
      notificationStore.showError('Failed to update rule', message)
      throw err
    } finally {
      saving.value = false
    }
  }
  
  const deleteRule = async (id: string) => {
    saving.value = true
    error.value = null
    
    try {
      await rulesApi.deleteRule(id)
      
      // Remove from local state
      rules.value = rules.value.filter(rule => rule.id !== id)
      
      if (currentRule.value?.id === id) {
        currentRule.value = null
      }
      
      // Remove metrics cache
      ruleMetrics.value.delete(id)
      
      notificationStore.showSuccess('Rule deleted successfully')
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to delete rule'
      error.value = message
      notificationStore.showError('Failed to delete rule', message)
      throw err
    } finally {
      saving.value = false
    }
  }
  
  const duplicateRule = async (id: string, newName?: string) => {
    const originalRule = await fetchRule(id)
    
    const duplicateRequest: CreateRuleRequest = {
      name: newName || `${originalRule.name} (Copy)`,
      description: originalRule.description,
      dslContent: originalRule.dslContent,
      priority: originalRule.priority,
      category: originalRule.category,
      tags: [...(originalRule.tags || [])]
    }
    
    return createRule(duplicateRequest)
  }
  
  const submitForApproval = async (id: string) => {
    saving.value = true
    error.value = null
    
    try {
      const response = await rulesApi.submitForApproval(id)
      
      // Update local state
      const index = rules.value.findIndex(rule => rule.id === id)
      if (index !== -1) {
        rules.value[index] = { ...rules.value[index], status: 'UNDER_REVIEW' }
      }
      
      notificationStore.showSuccess('Rule submitted for approval')
      return response.data
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to submit rule for approval'
      error.value = message
      notificationStore.showError('Failed to submit for approval', message)
      throw err
    } finally {
      saving.value = false
    }
  }
  
  const approveRule = async (id: string) => {
    saving.value = true
    error.value = null
    
    try {
      const response = await rulesApi.approveRule(id)
      
      // Update local state
      const index = rules.value.findIndex(rule => rule.id === id)
      if (index !== -1) {
        rules.value[index] = { ...rules.value[index], status: 'APPROVED' }
      }
      
      notificationStore.showSuccess('Rule approved successfully')
      return response.data
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to approve rule'
      error.value = message
      notificationStore.showError('Failed to approve rule', message)
      throw err
    } finally {
      saving.value = false
    }
  }
  
  const activateRule = async (id: string) => {
    saving.value = true
    error.value = null
    
    try {
      const response = await rulesApi.activateRule(id)
      
      // Update local state
      const index = rules.value.findIndex(rule => rule.id === id)
      if (index !== -1) {
        rules.value[index] = { ...rules.value[index], status: 'ACTIVE' }
      }
      
      notificationStore.showSuccess('Rule activated successfully')
      return response.data
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to activate rule'
      error.value = message
      notificationStore.showError('Failed to activate rule', message)
      throw err
    } finally {
      saving.value = false
    }
  }
  
  const deactivateRule = async (id: string) => {
    saving.value = true
    error.value = null
    
    try {
      const response = await rulesApi.deactivateRule(id)
      
      // Update local state
      const index = rules.value.findIndex(rule => rule.id === id)
      if (index !== -1) {
        rules.value[index] = { ...rules.value[index], status: 'INACTIVE' }
      }
      
      notificationStore.showSuccess('Rule deactivated successfully')
      return response.data
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to deactivate rule'
      error.value = message
      notificationStore.showError('Failed to deactivate rule', message)
      throw err
    } finally {
      saving.value = false
    }
  }
  
  const fetchRuleMetrics = async (id: string) => {
    try {
      const response = await rulesApi.getRuleMetrics(id)
      ruleMetrics.value.set(id, response.data)
      return response.data
    } catch (err) {
      console.error('Failed to fetch rule metrics:', err)
      return null
    }
  }
  
  const validateRule = async (dslContent: string, testData?: Record<string, any>) => {
    try {
      const response = await rulesApi.validateRule({
        dslContent,
        testData
      })
      return response.data
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to validate rule'
      notificationStore.showError('Rule validation failed', message)
      throw err
    }
  }
  
  // Utility actions
  const clearError = () => {
    error.value = null
  }
  
  const clearCurrentRule = () => {
    currentRule.value = null
  }
  
  const updateFilters = (newFilters: Partial<RuleFilters>) => {
    filters.value = { ...filters.value, ...newFilters }
    pagination.value.page = 1 // Reset to first page when filters change
  }
  
  const updateSearch = (query: string) => {
    searchQuery.value = query
    pagination.value.page = 1 // Reset to first page when search changes
  }
  
  const reset = () => {
    rules.value = []
    currentRule.value = null
    loading.value = false
    saving.value = false
    error.value = null
    filters.value = {
      status: [],
      priority: [],
      category: '',
      createdBy: '',
      dateRange: null,
      tags: []
    }
    searchQuery.value = ''
    pagination.value = {
      page: 1,
      rowsPerPage: 20,
      rowsNumber: 0,
      sortBy: 'createdAt',
      descending: true
    }
    ruleMetrics.value.clear()
  }
  
  // Watch for filter changes to automatically refetch
  watch([filters, searchQuery], () => {
    fetchRules()
  }, { deep: true })
  
  return {
    // State
    rules: readonly(rules),
    currentRule: readonly(currentRule),
    loading: readonly(loading),
    saving: readonly(saving),
    error: readonly(error),
    pagination: readonly(pagination),
    filters: readonly(filters),
    searchQuery: readonly(searchQuery),
    ruleMetrics: readonly(ruleMetrics),
    
    // Getters
    filteredRules,
    rulesByStatus,
    rulesByPriority,
    rulesStats,
    
    // Actions
    fetchRules,
    fetchRule,
    createRule,
    updateRule,
    deleteRule,
    duplicateRule,
    submitForApproval,
    approveRule,
    activateRule,
    deactivateRule,
    fetchRuleMetrics,
    validateRule,
    
    // Utility actions
    clearError,
    clearCurrentRule,
    updateFilters,
    updateSearch,
    reset
  }
})
```

### Example 4: Behavioral Test Implementation

```gherkin
# tests/e2e/rule-management.feature
Feature: Rule Management
  As a business user
  I want to manage business rules through the web interface
  So that I can control system behavior effectively

  Background:
    Given I am logged in as a business user
    And I am on the rules management page

  Scenario: View rules list with filtering
    Given there are 10 rules with different statuses
    When I apply a filter for "ACTIVE" status
    Then I should see only active rules in the list
    And the rule count should be updated accordingly

  Scenario: Create a new rule successfully
    When I click the "Create New Rule" button
    And I fill in the rule form with:
      | name        | Customer Loyalty Rule |
      | description | Apply discount for VIP customers |
      | dslContent  | customer.tier == "VIP" |
      | priority    | HIGH |
    And I click "Create Rule"
    Then the rule should be created successfully
    And I should see a success notification
    And the rule should appear in the rules list

  Scenario: Edit an existing rule
    Given there is a rule named "Test Rule" in draft status
    When I click the edit button for "Test Rule"
    And I update the description to "Updated description"
    And I click "Update Rule"
    Then the rule should be updated successfully
    And I should see the updated description in the rules list

  Scenario: Delete a rule with confirmation
    Given there is a rule named "Delete Me" in draft status
    When I click the delete button for "Delete Me"
    And I confirm the deletion in the dialog
    Then the rule should be deleted successfully
    And I should see a success notification
    And the rule should no longer appear in the rules list

  Scenario: Test rule validation
    When I create a rule with invalid DSL content
    Then I should see validation errors
    And the rule should not be created
    And I should remain on the create form

  Scenario: Search rules by name
    Given there are rules with names containing "discount"
    When I enter "discount" in the search box
    Then I should see only rules with "discount" in their names
    And the search results should be highlighted

  Scenario: Rule card keyboard navigation
    Given I am on the rules list page
    When I navigate using the keyboard
    And I press Tab to focus on a rule card
    And I press Enter to select the rule
    Then the rule details should be displayed
    And the navigation should be accessible

  Scenario: Responsive rule cards on mobile
    Given I am viewing the rules list on a mobile device
    When I scroll through the rule cards
    Then the cards should be properly formatted for mobile
    And all actions should be accessible via touch
```

## 8. Output Specifications

### Format Requirements:
- **Component Structure**: Vue 3 Composition API with TypeScript interfaces
- **Code Quality**: ESLint, Prettier, and Vue/TypeScript strict mode compliance
- **Accessibility**: WCAG 2.1 AA compliance with keyboard navigation and screen reader support
- **Documentation**: Storybook stories for all components with interaction examples

### Quality Criteria:
- **Performance**: <3s initial load time, 60fps interactions, Core Web Vitals compliance
- **Type Safety**: 100% TypeScript coverage with strict mode enabled
- **Testing**: >80% code coverage with unit, integration, and E2E tests
- **Bundle Size**: Optimized bundle with code splitting and lazy loading

## 9. Validation Checkpoints

### Pre-execution Validation:
- [ ] Screen specifications and user workflows documented
- [ ] Design system tokens and component patterns defined
- [ ] API contracts and data types documented
- [ ] Accessibility requirements and testing strategy planned

### Mid-execution Validation:
- [ ] Shared components follow design system consistently
- [ ] Screens implement specified user workflows correctly
- [ ] State management handles all data scenarios properly
- [ ] Components meet accessibility standards with keyboard navigation

### Post-execution Validation:
- [ ] All screens functional with proper error handling
- [ ] Performance metrics meet requirements under load
- [ ] Accessibility audit passes with 100% compliance
- [ ] Cross-browser testing completed successfully
- [ ] Production deployment successful with monitoring

## Implementation Tasks Breakdown

### Phase 1: Project Setup and Architecture (Days 1-2)
1. **Project Initialization**
   - Create Vue 3 + TypeScript + Vite project
   - Configure ESLint, Prettier, and Git hooks
   - Setup Quasar Framework with custom theme
   - Configure build optimization and bundle analysis

2. **Architecture Foundation**
   - Create project structure following best practices
   - Setup dependency injection and configuration
   - Configure routing with type-safe navigation
   - Add error boundary and global error handling

### Phase 2: Shared Component Library (Days 3-6)
1. **Design System Implementation**
   - Define design tokens (colors, spacing, typography)
   - Create base components (Button, Input, Select)
   - Implement layout components (AppLayout, PageHeader)
   - Add feedback components (Loading, Error, Success)

2. **Component Documentation**
   - Setup Storybook with Vue 3 support
   - Create stories for all components
   - Add interaction testing with Storybook
   - Generate component documentation

### Phase 3: Business Components (Days 7-10)
1. **Entity Components**
   - Create RuleCard, CampaignCard, CustomerCard
   - Implement data display components (DataTable, Charts)
   - Add form components (DynamicForm, DSLEditor)
   - Create specialized UI components

2. **Component Integration**
   - Test component composition patterns
   - Implement component prop validation
   - Add component performance monitoring
   - Create component usage guidelines

### Phase 4: State Management (Days 11-12)
1. **Pinia Stores**
   - Create stores for all entities (rules, campaigns, etc.)
   - Implement caching and optimistic updates
   - Add error handling and retry logic
   - Create computed properties and getters

2. **API Integration**
   - Implement API client with interceptors
   - Add request/response type safety
   - Create authentication and authorization
   - Implement offline capability with service workers

### Phase 5: Screen Implementation (Days 13-16)
1. **Core Screens**
   - Implement Dashboard with real-time data
   - Create Rules List with advanced filtering
   - Build Rule Detail with tabbed interface
   - Add Rule Create/Edit forms with validation

2. **Advanced Features**
   - Implement Campaign Management screens
   - Add bulk operations and data export
   - Create admin dashboard with user management
   - Add settings and configuration screens

### Phase 6: Testing and Quality Assurance (Days 17-18)
1. **Comprehensive Testing**
   - Write unit tests for all components
   - Create integration tests for user workflows
   - Implement E2E tests with Cypress
   - Add visual regression testing

2. **Accessibility and Performance**
   - Conduct accessibility audit and remediation
   - Optimize performance with lazy loading
   - Implement progressive web app features
   - Add monitoring and analytics

### Phase 7: Production Readiness (Days 19-20)
1. **Optimization and Deployment**
   - Optimize bundle size and loading performance
   - Configure CDN and static asset optimization
   - Setup deployment pipeline with Docker
   - Add monitoring and error tracking

2. **Documentation and Training**
   - Create user guides and tutorials
   - Document component library usage
   - Create developer onboarding guide
   - Setup maintenance and update procedures

## Best Practices

### Component Design:
- Use Composition API with script setup syntax
- Implement proper TypeScript interfaces for all props
- Follow single responsibility principle
- Create reusable and composable components

### Performance Optimization:
- Implement lazy loading for routes and components
- Use Vue's built-in performance optimizations
- Optimize bundle size with code splitting
- Implement effective caching strategies

### Accessibility Implementation:
- Use semantic HTML elements appropriately
- Implement keyboard navigation for all interactions
- Add ARIA labels and descriptions
- Test with screen readers and accessibility tools

### Testing Strategy:
- Follow testing pyramid (unit > integration > E2E)
- Test user workflows and edge cases
- Implement visual regression testing
- Add performance and accessibility testing
