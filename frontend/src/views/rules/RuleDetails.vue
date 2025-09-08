<template>
  <q-page class="rule-details-page">
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
            <h1 class="page-title">{{ rule?.name || 'Loading...' }}</h1>
            <p class="page-subtitle">
              {{ rule?.description || 'Rule details and configuration' }}
            </p>
          </div>
        </div>
        
        <div class="page-actions" v-if="rule">
          <q-btn
            flat
            icon="edit"
            label="Edit"
            @click="editRule"
          />
          <q-btn
            :color="rule.status === 'ACTIVE' ? 'orange' : 'positive'"
            :icon="rule.status === 'ACTIVE' ? 'pause' : 'play_arrow'"
            :label="rule.status === 'ACTIVE' ? 'Deactivate' : 'Activate'"
            @click="toggleRuleStatus"
          />
          <q-btn
            color="primary"
            icon="play_circle"
            label="Test Rule"
            @click="showTestDialog = true"
          />
        </div>
      </div>
    </div>

    <div class="details-content" v-if="rule">
      <div class="details-main">
        <!-- Rule Information -->
        <q-card class="rule-info-card">
          <q-card-section>
            <div class="rule-header">
              <div class="rule-title-section">
                <h2 class="rule-title">{{ rule.name }}</h2>
                <div class="rule-badges">
                  <q-chip
                    :color="getStatusColor(rule.status)"
                    text-color="white"
                    :icon="getStatusIcon(rule.status)"
                  >
                    {{ rule.status }}
                  </q-chip>
                  <q-chip
                    :color="getPriorityColor(rule.priority)"
                    text-color="white"
                    :icon="getPriorityIcon(rule.priority)"
                  >
                    {{ rule.priority }}
                  </q-chip>
                  <q-chip
                    color="info"
                    text-color="white"
                  >
                    {{ rule.category }}
                  </q-chip>
                </div>
              </div>
              
              <div class="rule-meta">
                <div class="meta-item">
                  <span class="meta-label">Version:</span>
                  <span class="meta-value">{{ rule.version }}</span>
                </div>
                <div class="meta-item">
                  <span class="meta-label">Created:</span>
                  <span class="meta-value">{{ formatDate(rule.created_at) }}</span>
                </div>
                <div class="meta-item">
                  <span class="meta-label">Updated:</span>
                  <span class="meta-value">{{ formatDate(rule.updated_at) }}</span>
                </div>
                <div class="meta-item">
                  <span class="meta-label">Created By:</span>
                  <span class="meta-value">{{ rule.created_by }}</span>
                </div>
              </div>
            </div>
            
            <div class="rule-description">
              <h3 class="section-title">Description</h3>
              <p class="description-text">{{ rule.description }}</p>
            </div>
            
            <div class="rule-tags" v-if="rule.tags && rule.tags.length > 0">
              <h3 class="section-title">Tags</h3>
              <div class="tags-list">
                <q-chip
                  v-for="tag in rule.tags"
                  :key="tag"
                  color="grey-6"
                  text-color="white"
                  dense
                >
                  {{ tag }}
                </q-chip>
              </div>
            </div>
          </q-card-section>
        </q-card>

        <!-- DSL Content -->
        <q-card class="dsl-card">
          <q-card-section>
            <div class="dsl-header">
              <h3 class="section-title">Rule Definition (DSL)</h3>
              <div class="dsl-actions">
                <q-btn
                  flat
                  dense
                  icon="content_copy"
                  @click="copyDSL"
                >
                  <q-tooltip>Copy DSL</q-tooltip>
                </q-btn>
                <q-btn
                  flat
                  dense
                  icon="play_circle"
                  @click="showTestDialog = true"
                >
                  <q-tooltip>Test Rule</q-tooltip>
                </q-btn>
                <q-btn
                  flat
                  dense
                  icon="check_circle"
                  @click="validateRule"
                >
                  <q-tooltip>Validate Rule</q-tooltip>
                </q-btn>
              </div>
            </div>
            
            <div class="dsl-content">
              <pre class="dsl-code">{{ rule.dsl_content }}</pre>
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

      <div class="details-sidebar">
        <!-- Rule Statistics -->
        <q-card class="stats-card">
          <q-card-section>
            <h3 class="section-title">Statistics</h3>
            
            <div class="stats-content">
              <div class="stat-item">
                <div class="stat-value">{{ rule.evaluation_count || 0 }}</div>
                <div class="stat-label">Evaluations</div>
              </div>
              
              <div class="stat-item">
                <div class="stat-value">{{ rule.success_rate || 0 }}%</div>
                <div class="stat-label">Success Rate</div>
              </div>
              
              <div class="stat-item">
                <div class="stat-value">{{ rule.last_evaluated || 'Never' }}</div>
                <div class="stat-label">Last Evaluated</div>
              </div>
            </div>
          </q-card-section>
        </q-card>

        <!-- Configuration -->
        <q-card class="config-card">
          <q-card-section>
            <h3 class="section-title">Configuration</h3>
            
            <div class="config-content">
              <div class="config-item">
                <span class="config-label">Effective Date:</span>
                <span class="config-value">{{ rule.effective_date || 'Not set' }}</span>
              </div>
              
              <div class="config-item">
                <span class="config-label">Expiration Date:</span>
                <span class="config-value">{{ rule.expiration_date || 'Not set' }}</span>
              </div>
              
              <div class="config-item">
                <span class="config-label">Requires Approval:</span>
                <q-icon
                  :name="rule.requires_approval ? 'check_circle' : 'cancel'"
                  :color="rule.requires_approval ? 'positive' : 'negative'"
                  size="sm"
                />
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
                icon="edit"
                label="Edit Rule"
                @click="editRule"
                class="action-button"
              />
              
              <q-btn
                color="secondary"
                icon="content_copy"
                label="Duplicate"
                @click="duplicateRule"
                class="action-button"
              />
              
              <q-btn
                color="info"
                icon="play_circle"
                label="Test Rule"
                @click="showTestDialog = true"
                class="action-button"
              />
              
              <q-btn
                color="warning"
                icon="history"
                label="View History"
                @click="viewHistory"
                class="action-button"
              />
              
              <q-btn
                color="negative"
                icon="delete"
                label="Delete Rule"
                @click="deleteRule"
                class="action-button"
              />
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Loading State -->
    <div v-else-if="loading" class="loading-state">
      <q-spinner-dots size="50px" color="primary" />
      <p>Loading rule details...</p>
    </div>

    <!-- Error State -->
    <div v-else class="error-state">
      <q-icon name="error" size="50px" color="negative" />
      <p>Failed to load rule details</p>
      <q-btn color="primary" label="Retry" @click="loadRule" />
    </div>

    <!-- Test Dialog -->
    <q-dialog v-model="showTestDialog">
      <q-card style="min-width: 500px">
        <q-card-section>
          <div class="text-h6">Test Rule</div>
        </q-card-section>
        
        <q-card-section>
          <div class="test-form">
            <q-input
              v-model="testContext.quantity"
              label="Quantity"
              type="number"
              outlined
              dense
              class="test-input"
            />
            
            <q-input
              v-model="testContext.price"
              label="Price"
              type="number"
              outlined
              dense
              class="test-input"
            />
            
            <q-input
              v-model="testContext.customer_age"
              label="Customer Age"
              type="number"
              outlined
              dense
              class="test-input"
            />
            
            <q-input
              v-model="testContext.product_family"
              label="Product Family"
              outlined
              dense
              class="test-input"
            />
          </div>
        </q-card-section>
        
        <q-card-actions align="right">
          <q-btn flat label="Cancel" v-close-popup />
          <q-btn 
            color="primary" 
            label="Run Test" 
            @click="runTest"
            :loading="testing"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useQuasar } from 'quasar'
import { useRulesStore } from '@/stores/rules'
import { useNotificationStore } from '@/stores/notifications'
import { formatDistanceToNow } from 'date-fns'
import type { Rule } from '@/types'

const router = useRouter()
const route = useRoute()
const $q = useQuasar()
const rulesStore = useRulesStore()
const notificationStore = useNotificationStore()

// Reactive state
const loading = ref(false)
const testing = ref(false)
const showTestDialog = ref(false)
const testResult = ref<any>(null)

const testContext = ref({
  quantity: 3,
  price: 100,
  customer_age: 30,
  product_family: 'electronics'
})

// Computed properties
const ruleId = computed(() => route.params.id as string)
const rule = computed(() => rulesStore.currentRule)

// Methods
const goBack = () => {
  router.back()
}

const editRule = () => {
  router.push(`/rules/${ruleId.value}/edit`)
}

const toggleRuleStatus = async () => {
  if (!rule.value) return
  
  try {
    if (rule.value.status === 'ACTIVE') {
      await rulesStore.deactivateRule(ruleId.value)
    } else {
      await rulesStore.activateRule(ruleId.value)
    }
  } catch (error) {
    // Error handling is done in the store
  }
}

const copyDSL = () => {
  if (rule.value?.dsl_content) {
    navigator.clipboard.writeText(rule.value.dsl_content)
    notificationStore.showSuccess('Success', 'DSL content copied to clipboard')
  }
}

const validateRule = async () => {
  if (!rule.value) return
  
  try {
    await rulesStore.validateRule({
      dsl_content: rule.value.dsl_content,
      context: {},
      rule_category: rule.value.category
    })
    
    notificationStore.showSuccess('Success', 'Rule validation passed')
  } catch (error) {
    notificationStore.showError('Error', 'Rule validation failed')
  }
}

const runTest = async () => {
  if (!rule.value) return
  
  testing.value = true
  try {
    const result = await rulesStore.validateRule({
      dsl_content: rule.value.dsl_content,
      context: testContext.value,
      rule_category: rule.value.category
    })
    
    testResult.value = {
      success: true,
      message: 'Rule test completed successfully',
      result: result
    }
    
    showTestDialog.value = false
  } catch (error) {
    testResult.value = {
      success: false,
      message: 'Rule test failed',
      result: null
    }
  } finally {
    testing.value = false
  }
}

const duplicateRule = async () => {
  if (!rule.value) return
  
  try {
    await rulesStore.duplicateRule(ruleId.value)
  } catch (error) {
    // Error handling is done in the store
  }
}

const viewHistory = () => {
  // Implement history view
  notificationStore.showInfo('Info', 'History view will be implemented')
}

const deleteRule = () => {
  if (!rule.value) return
  
  $q.dialog({
    title: 'Confirm Delete',
    message: `Are you sure you want to delete the rule "${rule.value.name}"?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await rulesStore.deleteRule(ruleId.value)
      router.push('/rules')
    } catch (error) {
      // Error handling is done in the store
    }
  })
}

const loadRule = async () => {
  loading.value = true
  try {
    await rulesStore.fetchRule(ruleId.value)
  } catch (error) {
    notificationStore.showError('Error', 'Failed to load rule')
  } finally {
    loading.value = false
  }
}

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    'ACTIVE': 'positive',
    'INACTIVE': 'grey',
    'DRAFT': 'info',
    'UNDER_REVIEW': 'warning',
    'APPROVED': 'positive',
    'DEPRECATED': 'negative'
  }
  return colors[status] || 'grey'
}

const getStatusIcon = (status: string) => {
  const icons: Record<string, string> = {
    'ACTIVE': 'play_circle',
    'INACTIVE': 'pause_circle',
    'DRAFT': 'edit',
    'UNDER_REVIEW': 'pending',
    'APPROVED': 'check_circle',
    'DEPRECATED': 'block'
  }
  return icons[status] || 'help'
}

const getPriorityColor = (priority: string) => {
  const colors: Record<string, string> = {
    'LOW': 'grey',
    'MEDIUM': 'info',
    'HIGH': 'warning',
    'CRITICAL': 'negative'
  }
  return colors[priority] || 'grey'
}

const getPriorityIcon = (priority: string) => {
  const icons: Record<string, string> = {
    'LOW': 'keyboard_arrow_down',
    'MEDIUM': 'remove',
    'HIGH': 'keyboard_arrow_up',
    'CRITICAL': 'priority_high'
  }
  return icons[priority] || 'remove'
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

// Lifecycle
onMounted(() => {
  loadRule()
})
</script>

<style lang="scss" scoped>
.rule-details-page {
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

.details-content {
  display: flex;
  gap: 24px;
  padding: 0 24px 24px;
}

.details-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.details-sidebar {
  width: 320px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.rule-info-card,
.dsl-card,
.test-results-card,
.stats-card,
.config-card,
.actions-card {
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.rule-header {
  margin-bottom: 24px;
}

.rule-title-section {
  margin-bottom: 16px;
}

.rule-title {
  font-size: 24px;
  font-weight: 600;
  color: var(--carrefour-dark-blue);
  margin: 0 0 12px 0;
}

.rule-badges {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.rule-meta {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
}

.meta-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.meta-label {
  font-size: 12px;
  font-weight: 500;
  color: var(--carrefour-gray-600);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.meta-value {
  font-size: 14px;
  color: var(--carrefour-gray-800);
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--carrefour-dark-blue);
  margin: 0 0 16px 0;
}

.rule-description {
  margin-bottom: 24px;
}

.description-text {
  font-size: 16px;
  line-height: 1.6;
  color: var(--carrefour-gray-700);
  margin: 0;
}

.rule-tags {
  margin-bottom: 24px;
}

.tags-list {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.dsl-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.dsl-actions {
  display: flex;
  gap: 4px;
}

.dsl-content {
  background: var(--carrefour-gray-100);
  border-radius: 8px;
  padding: 16px;
  overflow-x: auto;
}

.dsl-code {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.5;
  color: var(--carrefour-gray-800);
  margin: 0;
  white-space: pre-wrap;
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

.stats-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(100px, 1fr));
  gap: 16px;
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: var(--carrefour-dark-blue);
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: var(--carrefour-gray-600);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.config-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.config-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.config-label {
  font-weight: 500;
  color: var(--carrefour-gray-800);
}

.config-value {
  color: var(--carrefour-gray-600);
}

.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.action-button {
  width: 100%;
}

.loading-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  gap: 16px;
  
  p {
    font-size: 16px;
    color: var(--carrefour-gray-600);
    margin: 0;
  }
}

.test-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.test-input {
  width: 100%;
}

// Responsive design
@media (max-width: 1024px) {
  .details-content {
    flex-direction: column;
  }
  
  .details-sidebar {
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
  
  .details-content {
    padding: 0 16px 16px;
  }
  
  .rule-meta {
    grid-template-columns: 1fr;
  }
  
  .stats-content {
    grid-template-columns: repeat(3, 1fr);
  }
}
</style>
