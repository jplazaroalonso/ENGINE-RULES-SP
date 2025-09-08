<template>
  <q-page class="rule-editor-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="page-header-content">
        <div class="page-title-section">
          <q-btn
            flat
            icon="arrow_back"
            @click="goBack"
            class="back-button"
          />
          <div>
            <h1 class="page-title">
              {{ isEditing ? 'Edit Rule' : 'Create New Rule' }}
            </h1>
            <p class="page-subtitle">
              {{ isEditing ? 'Modify rule configuration and settings' : 'Define a new business rule' }}
            </p>
          </div>
        </div>
        
        <div class="page-actions">
          <q-btn
            flat
            label="Cancel"
            @click="goBack"
            :disable="saving"
          />
          <q-btn
            color="secondary"
            label="Save Draft"
            @click="saveDraft"
            :loading="saving"
            :disable="!isFormValid"
          />
          <q-btn
            color="primary"
            :label="isEditing ? 'Update Rule' : 'Create Rule'"
            @click="saveRule"
            :loading="saving"
            :disable="!isFormValid"
          />
        </div>
      </div>
    </div>

    <div class="editor-content">
      <div class="editor-main">
        <!-- Rule Form -->
        <q-card class="rule-form-card">
          <q-card-section>
            <div class="form-section">
              <h3 class="section-title">Basic Information</h3>
              
              <div class="form-row">
                <q-input
                  v-model="ruleForm.name"
                  label="Rule Name"
                  outlined
                  dense
                  :rules="[val => !!val || 'Rule name is required']"
                  class="form-field"
                />
                
                <q-select
                  v-model="ruleForm.category"
                  :options="categoryOptions"
                  label="Category"
                  outlined
                  dense
                  emit-value
                  map-options
                  :rules="[val => !!val || 'Category is required']"
                  class="form-field"
                />
              </div>
              
              <div class="form-row">
                <q-input
                  v-model="ruleForm.description"
                  label="Description"
                  outlined
                  dense
                  type="textarea"
                  rows="3"
                  class="form-field full-width"
                />
              </div>
              
              <div class="form-row">
                <q-select
                  v-model="ruleForm.priority"
                  :options="priorityOptions"
                  label="Priority"
                  outlined
                  dense
                  emit-value
                  map-options
                  class="form-field"
                />
                
                <q-input
                  v-model="ruleForm.version"
                  label="Version"
                  outlined
                  dense
                  readonly
                  class="form-field"
                />
              </div>
            </div>

            <q-separator class="section-separator" />

            <div class="form-section">
              <h3 class="section-title">Rule Definition</h3>
              
              <div class="dsl-editor-section">
                <div class="dsl-editor-header">
                  <span class="dsl-label">DSL Content</span>
                  <div class="dsl-actions">
                    <q-btn
                      flat
                      dense
                      icon="help"
                      @click="showDSLHelp = true"
                    >
                      <q-tooltip>DSL Help</q-tooltip>
                    </q-btn>
                    <q-btn
                      flat
                      dense
                      icon="play_circle"
                      @click="testRule"
                      :disable="!ruleForm.dsl_content"
                    >
                      <q-tooltip>Test Rule</q-tooltip>
                    </q-btn>
                    <q-btn
                      flat
                      dense
                      icon="check_circle"
                      @click="validateRule"
                      :disable="!ruleForm.dsl_content"
                    >
                      <q-tooltip>Validate Rule</q-tooltip>
                    </q-btn>
                  </div>
                </div>
                
                <q-input
                  v-model="ruleForm.dsl_content"
                  outlined
                  dense
                  type="textarea"
                  rows="10"
                  placeholder="Enter your rule definition using the DSL syntax..."
                  :rules="[val => !!val || 'DSL content is required']"
                  class="dsl-editor"
                />
              </div>
              
              <div class="form-row">
                <q-input
                  v-model="ruleForm.tags"
                  label="Tags (comma-separated)"
                  outlined
                  dense
                  class="form-field full-width"
                />
              </div>
            </div>

            <q-separator class="section-separator" />

            <div class="form-section">
              <h3 class="section-title">Configuration</h3>
              
              <div class="form-row">
                <q-checkbox
                  v-model="ruleForm.is_active"
                  label="Active"
                  class="form-field"
                />
                
                <q-checkbox
                  v-model="ruleForm.requires_approval"
                  label="Requires Approval"
                  class="form-field"
                />
              </div>
              
              <div class="form-row">
                <q-input
                  v-model="ruleForm.effective_date"
                  label="Effective Date"
                  outlined
                  dense
                  type="date"
                  class="form-field"
                />
                
                <q-input
                  v-model="ruleForm.expiration_date"
                  label="Expiration Date"
                  outlined
                  dense
                  type="date"
                  class="form-field"
                />
              </div>
            </div>
          </q-card-section>
        </q-card>

        <!-- Test Results -->
        <q-card v-if="testResult" class="test-results-card">
          <q-card-section>
            <div class="test-results-header">
              <h3 class="section-title">Test Results</h3>
              <q-btn
                flat
                dense
                icon="close"
                @click="testResult = null"
              />
            </div>
            
            <div class="test-result-content">
              <div class="test-result-item">
                <span class="test-label">Status:</span>
                <q-chip
                  :color="testResult.success ? 'positive' : 'negative'"
                  text-color="white"
                  :icon="testResult.success ? 'check_circle' : 'error'"
                >
                  {{ testResult.success ? 'Success' : 'Failed' }}
                </q-chip>
              </div>
              
              <div v-if="testResult.message" class="test-result-item">
                <span class="test-label">Message:</span>
                <span class="test-value">{{ testResult.message }}</span>
              </div>
              
              <div v-if="testResult.result" class="test-result-item">
                <span class="test-label">Result:</span>
                <pre class="test-result-json">{{ JSON.stringify(testResult.result, null, 2) }}</pre>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <div class="editor-sidebar">
        <!-- Rule Preview -->
        <q-card class="preview-card">
          <q-card-section>
            <h3 class="section-title">Rule Preview</h3>
            
            <div class="preview-content">
              <div class="preview-item">
                <span class="preview-label">Name:</span>
                <span class="preview-value">{{ ruleForm.name || 'Untitled Rule' }}</span>
              </div>
              
              <div class="preview-item">
                <span class="preview-label">Category:</span>
                <span class="preview-value">{{ ruleForm.category || 'Not specified' }}</span>
              </div>
              
              <div class="preview-item">
                <span class="preview-label">Priority:</span>
                <span class="preview-value">{{ ruleForm.priority || 'MEDIUM' }}</span>
              </div>
              
              <div class="preview-item">
                <span class="preview-label">Status:</span>
                <q-chip
                  :color="ruleForm.is_active ? 'positive' : 'grey'"
                  text-color="white"
                  dense
                >
                  {{ ruleForm.is_active ? 'Active' : 'Inactive' }}
                </q-chip>
              </div>
            </div>
          </q-card-section>
        </q-card>

        <!-- Validation Status -->
        <q-card class="validation-card">
          <q-card-section>
            <h3 class="section-title">Validation</h3>
            
            <div class="validation-content">
              <div class="validation-item">
                <q-icon
                  :name="isFormValid ? 'check_circle' : 'error'"
                  :color="isFormValid ? 'positive' : 'negative'"
                  size="sm"
                />
                <span class="validation-text">
                  {{ isFormValid ? 'Form is valid' : 'Form has errors' }}
                </span>
              </div>
              
              <div class="validation-item">
                <q-icon
                  :name="ruleForm.dsl_content ? 'check_circle' : 'error'"
                  :color="ruleForm.dsl_content ? 'positive' : 'negative'"
                  size="sm"
                />
                <span class="validation-text">
                  {{ ruleForm.dsl_content ? 'DSL content provided' : 'DSL content required' }}
                </span>
              </div>
            </div>
          </q-card-section>
        </q-card>

        <!-- Quick Actions -->
        <q-card class="actions-card">
          <q-card-section>
            <h3 class="section-title">Quick Actions</h3>
            
            <div class="quick-actions">
              <q-btn
                color="primary"
                icon="save"
                label="Save"
                @click="saveRule"
                :loading="saving"
                :disable="!isFormValid"
                class="action-button"
              />
              
              <q-btn
                color="secondary"
                icon="play_circle"
                label="Test"
                @click="testRule"
                :disable="!ruleForm.dsl_content"
                class="action-button"
              />
              
              <q-btn
                color="info"
                icon="visibility"
                label="Preview"
                @click="previewRule"
                :disable="!ruleForm.dsl_content"
                class="action-button"
              />
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- DSL Help Dialog -->
    <DSLHelp v-model="showDSLHelp" />
  </q-page>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useQuasar } from 'quasar'
import { useRulesStore } from '@/stores/rules'
import { useNotificationStore } from '@/stores/notifications'
import DSLHelp from '@/components/DSLHelp.vue'
import type { Rule, RuleCreateRequest, RuleUpdateRequest } from '@/types'

const router = useRouter()
const route = useRoute()
const $q = useQuasar()
const rulesStore = useRulesStore()
const notificationStore = useNotificationStore()

// Reactive state
const saving = ref(false)
const showDSLHelp = ref(false)
const testResult = ref<any>(null)

const ruleForm = ref({
  name: '',
  description: '',
  dsl_content: '',
  category: '',
  priority: 'MEDIUM',
  version: '1.0.0',
  tags: '',
  is_active: true,
  requires_approval: false,
  effective_date: '',
  expiration_date: ''
})

// Computed properties
const isEditing = computed(() => !!route.params.id)
const ruleId = computed(() => route.params.id as string)

const isFormValid = computed(() => {
  return ruleForm.value.name && 
         ruleForm.value.category && 
         ruleForm.value.dsl_content
})

const categoryOptions = [
  { label: 'Promotions', value: 'PROMOTIONS' },
  { label: 'Coupons', value: 'COUPONS' },
  { label: 'Loyalty', value: 'LOYALTY' },
  { label: 'Taxes', value: 'TAXES' },
  { label: 'Payments', value: 'PAYMENTS' }
]

const priorityOptions = [
  { label: 'Low', value: 'LOW' },
  { label: 'Medium', value: 'MEDIUM' },
  { label: 'High', value: 'HIGH' },
  { label: 'Critical', value: 'CRITICAL' }
]

// Methods
const goBack = () => {
  router.back()
}

const saveDraft = async () => {
  saving.value = true
  try {
    const ruleData: RuleCreateRequest = {
      ...ruleForm.value,
      tags: ruleForm.value.tags.split(',').map(tag => tag.trim()).filter(tag => tag),
      status: 'DRAFT'
    }
    
    if (isEditing.value) {
      await rulesStore.updateRule(ruleId.value, ruleData as RuleUpdateRequest)
    } else {
      await rulesStore.createRule(ruleData)
    }
    
    notificationStore.showSuccess('Success', 'Rule saved as draft')
  } catch (error) {
    // Error handling is done in the store
  } finally {
    saving.value = false
  }
}

const saveRule = async () => {
  saving.value = true
  try {
    const ruleData: RuleCreateRequest = {
      ...ruleForm.value,
      tags: ruleForm.value.tags.split(',').map(tag => tag.trim()).filter(tag => tag),
      status: ruleForm.value.requires_approval ? 'UNDER_REVIEW' : 'ACTIVE'
    }
    
    if (isEditing.value) {
      await rulesStore.updateRule(ruleId.value, ruleData as RuleUpdateRequest)
    } else {
      await rulesStore.createRule(ruleData)
    }
    
    notificationStore.showSuccess('Success', isEditing.value ? 'Rule updated successfully' : 'Rule created successfully')
    router.push('/rules')
  } catch (error) {
    // Error handling is done in the store
  } finally {
    saving.value = false
  }
}

const testRule = async () => {
  if (!ruleForm.value.dsl_content) {
    notificationStore.showError('Error', 'Please provide DSL content to test')
    return
  }
  
  try {
    const testContext = {
      quantity: 3,
      price: 100,
      customer_age: 30,
      product_family: 'electronics'
    }
    
    const result = await rulesStore.validateRule({
      dsl_content: ruleForm.value.dsl_content,
      context: testContext,
      rule_category: ruleForm.value.category
    })
    
    testResult.value = {
      success: true,
      message: 'Rule test completed successfully',
      result: result
    }
  } catch (error) {
    testResult.value = {
      success: false,
      message: 'Rule test failed',
      result: null
    }
  }
}

const validateRule = async () => {
  if (!ruleForm.value.dsl_content) {
    notificationStore.showError('Error', 'Please provide DSL content to validate')
    return
  }
  
  try {
    await rulesStore.validateRule({
      dsl_content: ruleForm.value.dsl_content,
      context: {},
      rule_category: ruleForm.value.category
    })
    
    notificationStore.showSuccess('Success', 'Rule validation passed')
  } catch (error) {
    notificationStore.showError('Error', 'Rule validation failed')
  }
}

const previewRule = () => {
  // Implement rule preview functionality
  notificationStore.showInfo('Preview', 'Rule preview functionality will be implemented')
}

const loadRule = async () => {
  if (isEditing.value) {
    try {
      const rule = await rulesStore.fetchRule(ruleId.value)
      if (rule) {
        ruleForm.value = {
          name: rule.name,
          description: rule.description,
          dsl_content: rule.dsl_content,
          category: rule.category,
          priority: rule.priority,
          version: rule.version,
          tags: rule.tags?.join(', ') || '',
          is_active: rule.status === 'ACTIVE',
          requires_approval: false,
          effective_date: rule.effective_date || '',
          expiration_date: rule.expiration_date || ''
        }
      }
    } catch (error) {
      notificationStore.showError('Error', 'Failed to load rule')
      router.push('/rules')
    }
  }
}

// Lifecycle
onMounted(() => {
  loadRule()
})

// Watch for form changes to update preview
watch(ruleForm, () => {
  // Update preview in real-time
}, { deep: true })
</script>

<style lang="scss" scoped>
.rule-editor-page {
  background-color: var(--carrefour-gray-50);
  min-height: 100vh;
}

.page-header {
  background: white;
  border-bottom: 1px solid var(--carrefour-gray-200);
  padding: 24px;
  margin-bottom: 24px;
}

.page-header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 24px;
}

.page-title-section {
  display: flex;
  align-items: center;
  gap: 16px;
  flex: 1;
}

.back-button {
  color: var(--carrefour-gray-600);
}

.page-title {
  font-size: 28px;
  font-weight: 600;
  color: var(--carrefour-dark-blue);
  margin: 0 0 8px 0;
}

.page-subtitle {
  font-size: 16px;
  color: var(--carrefour-gray-600);
  margin: 0;
}

.page-actions {
  display: flex;
  gap: 12px;
  flex-shrink: 0;
}

.editor-content {
  display: flex;
  gap: 24px;
  padding: 0 24px 24px;
}

.editor-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.editor-sidebar {
  width: 320px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.rule-form-card,
.test-results-card,
.preview-card,
.validation-card,
.actions-card {
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.form-section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--carrefour-dark-blue);
  margin: 0 0 16px 0;
}

.section-separator {
  margin: 24px 0;
}

.form-row {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.form-field {
  flex: 1;
}

.form-field.full-width {
  width: 100%;
}

.dsl-editor-section {
  margin-bottom: 16px;
}

.dsl-editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.dsl-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--carrefour-gray-800);
}

.dsl-actions {
  display: flex;
  gap: 4px;
}

.dsl-editor {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
}

.test-results-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.test-result-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.test-result-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.test-label {
  font-weight: 500;
  color: var(--carrefour-gray-800);
  min-width: 80px;
}

.test-value {
  color: var(--carrefour-gray-600);
}

.test-result-json {
  background: var(--carrefour-gray-100);
  padding: 12px;
  border-radius: 8px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
  overflow-x: auto;
  margin: 0;
}

.preview-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.preview-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.preview-label {
  font-weight: 500;
  color: var(--carrefour-gray-800);
}

.preview-value {
  color: var(--carrefour-gray-600);
  text-align: right;
}

.validation-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.validation-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.validation-text {
  font-size: 14px;
  color: var(--carrefour-gray-800);
}

.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.action-button {
  width: 100%;
}

.dsl-help-content {
  h4 {
    margin: 16px 0 8px 0;
    color: var(--carrefour-dark-blue);
  }
  
  .dsl-example {
    background: var(--carrefour-gray-100);
    padding: 12px;
    border-radius: 8px;
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 14px;
    overflow-x: auto;
    margin: 8px 0;
  }
  
  ul {
    margin: 8px 0;
    padding-left: 20px;
  }
  
  li {
    margin: 4px 0;
  }
  
  code {
    background: var(--carrefour-gray-100);
    padding: 2px 6px;
    border-radius: 4px;
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 12px;
  }
}

// Responsive design
@media (max-width: 1024px) {
  .editor-content {
    flex-direction: column;
  }
  
  .editor-sidebar {
    width: 100%;
  }
}

@media (max-width: 768px) {
  .page-header {
    padding: 16px;
  }
  
  .page-header-content {
    flex-direction: column;
    gap: 16px;
  }
  
  .page-title-section {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .page-actions {
    width: 100%;
    justify-content: stretch;
    
    .q-btn {
      flex: 1;
    }
  }
  
  .editor-content {
    padding: 0 16px 16px;
  }
  
  .form-row {
    flex-direction: column;
    gap: 12px;
  }
  
  .dsl-editor-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .dsl-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>
