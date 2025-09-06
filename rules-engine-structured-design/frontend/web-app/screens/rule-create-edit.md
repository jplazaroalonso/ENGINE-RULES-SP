# Rule Create/Edit Screen - Rule Form Management

## Overview
The Rule Create/Edit screen provides a comprehensive form interface for creating new business rules or editing existing ones. It includes DSL editing, validation, testing, and template integration.

## Screen Layout

### Header Section with Context
```vue
<template>
  <div class="rule-form-header">
    <div class="breadcrumb-section">
      <q-breadcrumbs>
        <q-breadcrumbs-el label="Rules" to="/rules" />
        <q-breadcrumbs-el v-if="isEditMode" :label="rule?.name || 'Edit Rule'" />
        <q-breadcrumbs-el v-else label="Create New Rule" />
      </q-breadcrumbs>
    </div>
    
    <div class="form-header-info">
      <h1>{{ isEditMode ? 'Edit Rule' : 'Create New Rule' }}</h1>
      <p class="subtitle">
        {{ isEditMode 
          ? `Modify rule configuration and settings` 
          : 'Define a new business rule with DSL configuration' }}
      </p>
    </div>
    
    <div class="form-actions">
      <q-btn
        flat
        label="Cancel"
        @click="handleCancel"
        class="q-mr-sm"
      />
      
      <q-btn
        color="secondary"
        icon="play_arrow"
        label="Test Rule"
        @click="showTestDialog = true"
        :disable="!form.dslContent || validationErrors.length > 0"
        outline
        class="q-mr-sm"
      />
      
      <q-btn
        color="warning"
        icon="save"
        label="Save as Draft"
        @click="saveDraft"
        :loading="savingDraft"
        outline
        class="q-mr-sm"
        v-if="!isEditMode || rule?.status === 'DRAFT'"
      />
      
      <q-btn
        color="primary"
        icon="check"
        :label="isEditMode ? 'Update Rule' : 'Create Rule'"
        @click="handleSubmit"
        :loading="saving"
        :disable="!isFormValid"
      />
    </div>
  </div>
</template>
```

### Rule Form with Sections
```vue
<template>
  <div class="rule-form-content">
    <q-form @submit="handleSubmit" class="rule-form">
      
      <!-- Basic Information Section -->
      <q-card class="form-section">
        <q-card-section>
          <h3>Basic Information</h3>
          
          <div class="row q-gutter-md">
            <div class="col-md-6 col-12">
              <q-input
                v-model="form.name"
                label="Rule Name *"
                outlined
                dense
                :rules="[rules.required, rules.minLength(3), rules.maxLength(100)]"
                :error="!!fieldErrors.name"
                :error-message="fieldErrors.name"
                @blur="validateField('name')"
              >
                <template v-slot:hint>
                  Provide a clear, descriptive name for this rule
                </template>
              </q-input>
            </div>
            
            <div class="col-md-3 col-sm-6 col-12">
              <q-select
                v-model="form.priority"
                label="Priority *"
                :options="priorityOptions"
                outlined
                dense
                emit-value
                map-options
                :rules="[rules.required]"
              >
                <template v-slot:hint>
                  Set execution priority for conflict resolution
                </template>
              </q-select>
            </div>
            
            <div class="col-md-3 col-sm-6 col-12">
              <q-select
                v-model="form.category"
                label="Category"
                :options="categoryOptions"
                outlined
                dense
                use-input
                hide-selected
                fill-input
                input-debounce="300"
                @filter="filterCategories"
                @new-value="createCategory"
              >
                <template v-slot:no-option>
                  <q-item>
                    <q-item-section class="text-grey">
                      Type to create new category
                    </q-item-section>
                  </q-item>
                </template>
                <template v-slot:hint>
                  Organize rules by business domain
                </template>
              </q-select>
            </div>
          </div>
          
          <div class="row q-gutter-md q-mt-sm">
            <div class="col-12">
              <q-input
                v-model="form.description"
                label="Description"
                type="textarea"
                rows="3"
                outlined
                dense
                :rules="[rules.maxLength(500)]"
              >
                <template v-slot:hint>
                  Describe what this rule does and when it applies
                </template>
              </q-input>
            </div>
          </div>
          
          <div class="row q-gutter-md q-mt-sm">
            <div class="col-12">
              <q-input
                v-model="tagsInput"
                label="Tags"
                outlined
                dense
                @keyup.enter="addTag"
                @keyup.comma="addTag"
              >
                <template v-slot:hint>
                  Press Enter or comma to add tags
                </template>
              </q-input>
              
              <div class="tags-display q-mt-sm" v-if="form.tags.length > 0">
                <q-chip
                  v-for="tag in form.tags"
                  :key="tag"
                  removable
                  @remove="removeTag(tag)"
                  color="primary"
                  text-color="white"
                  size="sm"
                >
                  {{ tag }}
                </q-chip>
              </div>
            </div>
          </div>
        </q-card-section>
      </q-card>
      
      <!-- Template Selection Section -->
      <q-card class="form-section q-mt-md" v-if="!isEditMode">
        <q-card-section>
          <div class="template-header">
            <h3>Template Selection (Optional)</h3>
            <q-toggle
              v-model="useTemplate"
              label="Use template"
              color="primary"
            />
          </div>
          
          <div v-if="useTemplate" class="template-selection">
            <div class="row q-gutter-md">
              <div class="col-md-6 col-12">
                <q-select
                  v-model="selectedTemplate"
                  label="Select Template"
                  :options="templateOptions"
                  outlined
                  dense
                  option-label="name"
                  option-value="id"
                  @update:model-value="onTemplateSelected"
                >
                  <template v-slot:option="scope">
                    <q-item v-bind="scope.itemProps">
                      <q-item-section>
                        <q-item-label>{{ scope.opt.name }}</q-item-label>
                        <q-item-label caption>{{ scope.opt.description }}</q-item-label>
                      </q-item-section>
                      <q-item-section side>
                        <q-chip color="grey-3" text-color="grey-8" size="sm">
                          {{ scope.opt.category }}
                        </q-chip>
                      </q-item-section>
                    </q-item>
                  </template>
                </q-select>
              </div>
            </div>
            
            <!-- Template Parameters -->
            <div v-if="selectedTemplate && templateParameters.length > 0" class="template-parameters q-mt-md">
              <h4>Template Parameters</h4>
              
              <div class="row q-gutter-md">
                <div
                  v-for="param in templateParameters"
                  :key="param.id"
                  class="col-md-4 col-sm-6 col-12"
                >
                  <q-input
                    v-if="param.type === 'string' || param.type === 'number'"
                    v-model="templateParameterValues[param.name]"
                    :label="param.name + (param.required ? ' *' : '')"
                    :type="param.type === 'number' ? 'number' : 'text'"
                    outlined
                    dense
                    :rules="param.required ? [rules.required] : []"
                  >
                    <template v-slot:hint>
                      {{ param.description }}
                    </template>
                  </q-input>
                  
                  <q-select
                    v-else-if="param.type === 'select'"
                    v-model="templateParameterValues[param.name]"
                    :label="param.name + (param.required ? ' *' : '')"
                    :options="param.options"
                    outlined
                    dense
                    :rules="param.required ? [rules.required] : []"
                  >
                    <template v-slot:hint>
                      {{ param.description }}
                    </template>
                  </q-select>
                  
                  <q-toggle
                    v-else-if="param.type === 'boolean'"
                    v-model="templateParameterValues[param.name]"
                    :label="param.name"
                    color="primary"
                  />
                </div>
              </div>
              
              <div class="template-actions q-mt-md">
                <q-btn
                  color="primary"
                  label="Generate DSL from Template"
                  @click="generateDSLFromTemplate"
                  :loading="generatingDSL"
                  outline
                />
              </div>
            </div>
          </div>
        </q-card-section>
      </q-card>
      
      <!-- DSL Configuration Section -->
      <q-card class="form-section q-mt-md">
        <q-card-section>
          <div class="dsl-header">
            <h3>DSL Configuration *</h3>
            <div class="dsl-actions">
              <q-btn
                flat
                icon="help"
                label="DSL Help"
                @click="showDSLHelp = true"
                size="sm"
              />
              <q-btn
                flat
                icon="code"
                label="Validate"
                @click="validateDSL"
                :loading="validating"
                size="sm"
              />
              <q-btn
                flat
                icon="format_align_left"
                label="Format"
                @click="formatDSL"
                size="sm"
              />
            </div>
          </div>
          
          <!-- DSL Editor -->
          <div class="dsl-editor-container">
            <DSLEditor
              v-model="form.dslContent"
              :validation-errors="validationErrors"
              :validation-warnings="validationWarnings"
              @validate="onDSLValidation"
              @change="onDSLChange"
              class="dsl-editor"
              :rules="[rules.required]"
            />
            
            <!-- Validation Results -->
            <div v-if="validationResults" class="validation-results q-mt-sm">
              <q-banner
                v-if="validationErrors.length > 0"
                class="bg-negative text-white"
                icon="error"
                dense
              >
                <div class="validation-errors">
                  <div class="error-title">Validation Errors:</div>
                  <ul>
                    <li v-for="error in validationErrors" :key="error.line">
                      Line {{ error.line }}: {{ error.message }}
                    </li>
                  </ul>
                </div>
              </q-banner>
              
              <q-banner
                v-if="validationWarnings.length > 0"
                class="bg-warning text-white q-mt-sm"
                icon="warning"
                dense
              >
                <div class="validation-warnings">
                  <div class="warning-title">Warnings:</div>
                  <ul>
                    <li v-for="warning in validationWarnings" :key="warning.line">
                      Line {{ warning.line }}: {{ warning.message }}
                    </li>
                  </ul>
                </div>
              </q-banner>
              
              <q-banner
                v-if="validationErrors.length === 0 && validationWarnings.length === 0"
                class="bg-positive text-white q-mt-sm"
                icon="check_circle"
                dense
              >
                DSL syntax is valid
              </q-banner>
            </div>
          </div>
        </q-card-section>
      </q-card>
      
      <!-- Advanced Settings Section -->
      <q-card class="form-section q-mt-md">
        <q-card-section>
          <q-expansion-item
            icon="settings"
            label="Advanced Settings"
            v-model="advancedSettingsExpanded"
          >
            <div class="advanced-settings-content q-pa-md">
              <div class="row q-gutter-md">
                <div class="col-md-4 col-sm-6 col-12">
                  <q-input
                    v-model.number="form.executionTimeout"
                    label="Execution Timeout (ms)"
                    type="number"
                    outlined
                    dense
                    :min="100"
                    :max="30000"
                  >
                    <template v-slot:hint>
                      Maximum execution time in milliseconds
                    </template>
                  </q-input>
                </div>
                
                <div class="col-md-4 col-sm-6 col-12">
                  <q-input
                    v-model.number="form.maxRetries"
                    label="Max Retries"
                    type="number"
                    outlined
                    dense
                    :min="0"
                    :max="5"
                  >
                    <template v-slot:hint>
                      Number of retry attempts on failure
                    </template>
                  </q-input>
                </div>
                
                <div class="col-md-4 col-sm-6 col-12">
                  <q-toggle
                    v-model="form.enableCaching"
                    label="Enable Result Caching"
                    color="primary"
                  />
                </div>
              </div>
              
              <div class="row q-gutter-md q-mt-md">
                <div class="col-md-6 col-12">
                  <q-input
                    v-model="form.effectiveDate"
                    label="Effective Date"
                    type="date"
                    outlined
                    dense
                  >
                    <template v-slot:hint>
                      When this rule becomes active
                    </template>
                  </q-input>
                </div>
                
                <div class="col-md-6 col-12">
                  <q-input
                    v-model="form.expiryDate"
                    label="Expiry Date (Optional)"
                    type="date"
                    outlined
                    dense
                  >
                    <template v-slot:hint>
                      When this rule expires (leave empty for no expiry)
                    </template>
                  </q-input>
                </div>
              </div>
              
              <div class="row q-gutter-md q-mt-md">
                <div class="col-12">
                  <q-input
                    v-model="form.notes"
                    label="Implementation Notes"
                    type="textarea"
                    rows="3"
                    outlined
                    dense
                  >
                    <template v-slot:hint>
                      Internal notes for developers and reviewers
                    </template>
                  </q-input>
                </div>
              </div>
            </div>
          </q-expansion-item>
        </q-card-section>
      </q-card>
      
      <!-- Dependencies Section -->
      <q-card class="form-section q-mt-md" v-if="isEditMode && ruleDependencies.length > 0">
        <q-card-section>
          <h3>Dependencies</h3>
          
          <q-list>
            <q-item
              v-for="dependency in ruleDependencies"
              :key="dependency.id"
            >
              <q-item-section avatar>
                <q-icon :name="getDependencyIcon(dependency.type)" />
              </q-item-section>
              <q-item-section>
                <q-item-label>{{ dependency.name }}</q-item-label>
                <q-item-label caption>{{ dependency.description }}</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-chip
                  :color="dependency.status === 'ACTIVE' ? 'positive' : 'warning'"
                  text-color="white"
                  size="sm"
                >
                  {{ dependency.status }}
                </q-chip>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
      </q-card>
    </q-form>
    
    <!-- Test Dialog -->
    <q-dialog v-model="showTestDialog" persistent>
      <q-card style="min-width: 800px">
        <q-card-section>
          <div class="text-h6">Test Rule</div>
          <p class="text-subtitle2">Test your rule with sample data</p>
        </q-card-section>
        
        <q-card-section>
          <RuleTestInterface
            :dsl-content="form.dslContent"
            @test-complete="onTestComplete"
            @close="showTestDialog = false"
          />
        </q-card-section>
      </q-card>
    </q-dialog>
    
    <!-- DSL Help Dialog -->
    <q-dialog v-model="showDSLHelp">
      <q-card style="max-width: 1000px; width: 90vw">
        <q-card-section>
          <div class="text-h6">DSL Reference Guide</div>
        </q-card-section>
        
        <q-card-section>
          <DSLHelpContent />
        </q-card-section>
        
        <q-card-actions align="right">
          <q-btn flat label="Close" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>
```

## Component Data Structure

### Form Data Interface
```typescript
interface RuleFormData {
  name: string
  description: string
  category: string
  priority: Priority
  dslContent: string
  tags: string[]
  executionTimeout: number
  maxRetries: number
  enableCaching: boolean
  effectiveDate: string
  expiryDate?: string
  notes: string
}

interface FormState {
  form: RuleFormData
  isEditMode: boolean
  loading: boolean
  saving: boolean
  savingDraft: boolean
  
  // Template
  useTemplate: boolean
  selectedTemplate: RuleTemplate | null
  templateParameters: TemplateParameter[]
  templateParameterValues: Record<string, any>
  generatingDSL: boolean
  
  // Validation
  validationResults: ValidationResult | null
  validationErrors: ValidationError[]
  validationWarnings: ValidationWarning[]
  validating: boolean
  fieldErrors: Record<string, string>
  
  // UI State
  tagsInput: string
  advancedSettingsExpanded: boolean
  showTestDialog: boolean
  showDSLHelp: boolean
}

interface ValidationError {
  line: number
  column: number
  message: string
  severity: 'error' | 'warning'
}

interface TemplateParameter {
  id: string
  name: string
  type: 'string' | 'number' | 'boolean' | 'select'
  description: string
  required: boolean
  defaultValue?: any
  options?: SelectOption[]
}
```

### Form Validation Rules
```typescript
const validationRules = {
  required: (value: any) => !!value || 'This field is required',
  minLength: (min: number) => (value: string) => 
    value.length >= min || `Must be at least ${min} characters`,
  maxLength: (max: number) => (value: string) => 
    value.length <= max || `Must be at most ${max} characters`,
  uniqueName: async (value: string) => {
    if (!value) return true
    const exists = await rulesApi.checkNameExists(value, rule?.id)
    return !exists || 'Rule name already exists'
  }
}
```

## Screen Interactions

### Form Actions
- **Save Draft**: Save rule with DRAFT status
- **Create/Update Rule**: Submit form with validation
- **Test Rule**: Open test dialog with current DSL
- **Cancel**: Confirm and navigate back

### Template Integration
- **Select Template**: Load template DSL and parameters
- **Generate DSL**: Apply parameters to template DSL
- **Parameter Validation**: Validate required parameters

### DSL Editor Features
- **Syntax Highlighting**: Real-time syntax coloring
- **Auto-completion**: Context-aware suggestions
- **Validation**: Real-time error checking
- **Formatting**: Auto-format DSL code

## State Management
```typescript
// Rule Form Store
export const useRuleFormStore = defineStore('ruleForm', () => {
  const form = ref<RuleFormData>(getDefaultFormData())
  const validationErrors = ref<ValidationError[]>([])
  const loading = ref(false)
  
  const isFormValid = computed(() => {
    return form.value.name.trim() !== '' &&
           form.value.dslContent.trim() !== '' &&
           validationErrors.value.length === 0
  })
  
  const validateDSL = async (dslContent: string) => {
    try {
      const result = await rulesApi.validateDSL(dslContent)
      validationErrors.value = result.errors
      return result
    } catch (error) {
      validationErrors.value = [{ 
        line: 0, 
        column: 0, 
        message: 'Validation service unavailable', 
        severity: 'error' 
      }]
      throw error
    }
  }
  
  const saveRule = async (isDraft = false) => {
    loading.value = true
    try {
      const ruleData = {
        ...form.value,
        status: isDraft ? 'DRAFT' : 'UNDER_REVIEW'
      }
      
      if (isEditMode.value) {
        return await rulesApi.update(route.params.id, ruleData)
      } else {
        return await rulesApi.create(ruleData)
      }
    } finally {
      loading.value = false
    }
  }
  
  return {
    form,
    validationErrors: readonly(validationErrors),
    loading: readonly(loading),
    isFormValid,
    validateDSL,
    saveRule
  }
})
```

## Performance Optimizations
- Debounced DSL validation (500ms)
- Lazy load template options
- Virtual scrolling for large category lists
- Optimistic form updates
- Client-side validation before API calls

## Accessibility Features
- Keyboard navigation for all form fields
- Screen reader support for validation messages
- High contrast mode for DSL editor
- Focus management for modal dialogs
- ARIA labels for complex form sections

## Testing Strategy
- Unit tests for form validation logic
- Integration tests for template functionality
- E2E tests for complete form submission
- DSL editor functionality testing
- Accessibility tests with keyboard navigation
