# Rule Detail Screen - Individual Rule Management

## Overview
The Rule Detail screen provides a comprehensive view of a specific business rule, including its configuration, execution history, performance metrics, and management options.

## Screen Layout

### Header Section with Rule Information
```vue
<template>
  <div class="rule-detail-header">
    <div class="breadcrumb-section">
      <q-breadcrumbs>
        <q-breadcrumbs-el label="Rules" to="/rules" />
        <q-breadcrumbs-el :label="rule?.name || 'Rule Detail'" />
      </q-breadcrumbs>
    </div>
    
    <div class="rule-header-info">
      <div class="rule-title-section">
        <h1 class="rule-name">{{ rule?.name }}</h1>
        <div class="rule-metadata">
          <q-chip
            :color="getStatusColor(rule?.status)"
            text-color="white"
            :icon="getStatusIcon(rule?.status)"
            class="q-mr-sm"
          >
            {{ rule?.status }}
          </q-chip>
          
          <q-chip
            :color="getPriorityColor(rule?.priority)"
            text-color="white"
            class="q-mr-sm"
          >
            {{ rule?.priority }} Priority
          </q-chip>
          
          <q-chip
            color="grey-3"
            text-color="grey-8"
            icon="schedule"
            class="q-mr-sm"
          >
            Version {{ rule?.version }}
          </q-chip>
          
          <q-chip
            color="grey-3"
            text-color="grey-8"
            icon="person"
          >
            Created by {{ rule?.createdBy }}
          </q-chip>
        </div>
      </div>
      
      <div class="rule-actions">
        <q-btn
          color="primary"
          icon="edit"
          label="Edit Rule"
          @click="navigateToEdit"
          v-if="canEdit"
        />
        
        <q-btn
          color="positive"
          icon="play_arrow"
          label="Test Rule"
          @click="showTestDialog = true"
          outline
          class="q-ml-sm"
        />
        
        <q-btn-dropdown
          color="grey-7"
          icon="more_vert"
          outline
          class="q-ml-sm"
        >
          <q-list>
            <q-item clickable @click="duplicateRule" v-if="canDuplicate">
              <q-item-section avatar>
                <q-icon name="content_copy" />
              </q-item-section>
              <q-item-section>Duplicate Rule</q-item-section>
            </q-item>
            
            <q-item clickable @click="exportRule">
              <q-item-section avatar>
                <q-icon name="download" />
              </q-item-section>
              <q-item-section>Export Rule</q-item-section>
            </q-item>
            
            <q-item clickable @click="shareRule">
              <q-item-section avatar>
                <q-icon name="share" />
              </q-item-section>
              <q-item-section>Share Rule</q-item-section>
            </q-item>
            
            <q-separator />
            
            <q-item clickable @click="submitForApproval" v-if="canSubmitForApproval">
              <q-item-section avatar>
                <q-icon name="approval" color="warning" />
              </q-item-section>
              <q-item-section>Submit for Approval</q-item-section>
            </q-item>
            
            <q-item clickable @click="approveRule" v-if="canApprove">
              <q-item-section avatar>
                <q-icon name="check_circle" color="positive" />
              </q-item-section>
              <q-item-section>Approve Rule</q-item-section>
            </q-item>
            
            <q-item clickable @click="rejectRule" v-if="canReject">
              <q-item-section avatar>
                <q-icon name="cancel" color="negative" />
              </q-item-section>
              <q-item-section>Reject Rule</q-item-section>
            </q-item>
            
            <q-separator />
            
            <q-item clickable @click="activateRule" v-if="canActivate">
              <q-item-section avatar>
                <q-icon name="play_arrow" color="primary" />
              </q-item-section>
              <q-item-section>Activate</q-item-section>
            </q-item>
            
            <q-item clickable @click="deactivateRule" v-if="canDeactivate">
              <q-item-section avatar>
                <q-icon name="pause" color="warning" />
              </q-item-section>
              <q-item-section>Deactivate</q-item-section>
            </q-item>
            
            <q-separator />
            
            <q-item clickable @click="deleteRule" v-if="canDelete">
              <q-item-section avatar>
                <q-icon name="delete" color="negative" />
              </q-item-section>
              <q-item-section>Delete Rule</q-item-section>
            </q-item>
          </q-list>
        </q-btn-dropdown>
      </div>
    </div>
  </div>
</template>
```

### Tab Navigation for Rule Information
```vue
<template>
  <div class="rule-detail-tabs">
    <q-tabs
      v-model="activeTab"
      dense
      class="text-grey"
      active-color="primary"
      indicator-color="primary"
      align="justify"
      narrow-indicator
    >
      <q-tab name="overview" label="Overview" icon="info" />
      <q-tab name="configuration" label="Configuration" icon="settings" />
      <q-tab name="execution" label="Execution History" icon="history" />
      <q-tab name="performance" label="Performance" icon="analytics" />
      <q-tab name="versions" label="Versions" icon="source" />
      <q-tab name="approvals" label="Approvals" icon="approval" />
    </q-tabs>

    <q-separator />

    <q-tab-panels v-model="activeTab" animated class="rule-detail-panels">
      
      <!-- Overview Tab -->
      <q-tab-panel name="overview" class="q-pa-md">
        <div class="overview-content">
          <div class="row q-gutter-md">
            <!-- Rule Basic Information -->
            <div class="col-md-8 col-12">
              <q-card class="rule-info-card">
                <q-card-section>
                  <h3>Rule Information</h3>
                  
                  <div class="info-grid">
                    <div class="info-item">
                      <label>Description:</label>
                      <p>{{ rule?.description || 'No description provided' }}</p>
                    </div>
                    
                    <div class="info-item">
                      <label>Category:</label>
                      <p>{{ rule?.category || 'Uncategorized' }}</p>
                    </div>
                    
                    <div class="info-item">
                      <label>Tags:</label>
                      <div class="tags-display">
                        <q-chip
                          v-for="tag in rule?.tags"
                          :key="tag"
                          color="grey-3"
                          text-color="grey-8"
                          size="sm"
                        >
                          {{ tag }}
                        </q-chip>
                        <span v-if="!rule?.tags?.length" class="text-grey-6">No tags</span>
                      </div>
                    </div>
                    
                    <div class="info-item">
                      <label>Created:</label>
                      <p>{{ formatDateTime(rule?.createdAt) }} by {{ rule?.createdBy }}</p>
                    </div>
                    
                    <div class="info-item">
                      <label>Last Modified:</label>
                      <p>{{ formatDateTime(rule?.updatedAt) }}</p>
                    </div>
                    
                    <div class="info-item" v-if="rule?.approvedAt">
                      <label>Approved:</label>
                      <p>{{ formatDateTime(rule?.approvedAt) }} by {{ rule?.approvedBy }}</p>
                    </div>
                  </div>
                </q-card-section>
              </q-card>
            </div>
            
            <!-- Quick Stats -->
            <div class="col-md-4 col-12">
              <q-card class="stats-card">
                <q-card-section>
                  <h3>Quick Stats</h3>
                  
                  <div class="stats-grid">
                    <div class="stat-item">
                      <div class="stat-value">{{ ruleStats.executionCount }}</div>
                      <div class="stat-label">Total Executions</div>
                    </div>
                    
                    <div class="stat-item">
                      <div class="stat-value">{{ formatPercentage(ruleStats.successRate) }}</div>
                      <div class="stat-label">Success Rate</div>
                    </div>
                    
                    <div class="stat-item">
                      <div class="stat-value">{{ ruleStats.avgExecutionTime }}ms</div>
                      <div class="stat-label">Avg Execution Time</div>
                    </div>
                    
                    <div class="stat-item">
                      <div class="stat-value">{{ formatCurrency(ruleStats.impactValue) }}</div>
                      <div class="stat-label">Business Impact</div>
                    </div>
                  </div>
                </q-card-section>
              </q-card>
              
              <!-- Related Rules -->
              <q-card class="related-rules-card q-mt-md" v-if="relatedRules.length > 0">
                <q-card-section>
                  <h3>Related Rules</h3>
                  
                  <q-list>
                    <q-item
                      v-for="relatedRule in relatedRules"
                      :key="relatedRule.id"
                      clickable
                      @click="navigateToRule(relatedRule.id)"
                    >
                      <q-item-section>
                        <q-item-label>{{ relatedRule.name }}</q-item-label>
                        <q-item-label caption>{{ relatedRule.category }}</q-item-label>
                      </q-item-section>
                      <q-item-section side>
                        <q-chip
                          :color="getStatusColor(relatedRule.status)"
                          text-color="white"
                          size="sm"
                        >
                          {{ relatedRule.status }}
                        </q-chip>
                      </q-item-section>
                    </q-item>
                  </q-list>
                </q-card-section>
              </q-card>
            </div>
          </div>
        </div>
      </q-tab-panel>

      <!-- Configuration Tab -->
      <q-tab-panel name="configuration" class="q-pa-md">
        <div class="configuration-content">
          <q-card class="dsl-viewer-card">
            <q-card-section>
              <div class="dsl-header">
                <h3>Rule Configuration (DSL)</h3>
                <div class="dsl-actions">
                  <q-btn
                    flat
                    icon="content_copy"
                    label="Copy"
                    @click="copyDSLToClipboard"
                    size="sm"
                  />
                  <q-btn
                    flat
                    icon="download"
                    label="Download"
                    @click="downloadDSL"
                    size="sm"
                  />
                </div>
              </div>
              
              <div class="dsl-content">
                <pre class="dsl-code"><code>{{ rule?.dslContent }}</code></pre>
              </div>
            </q-card-section>
          </q-card>
          
          <!-- Rule Dependencies -->
          <q-card class="dependencies-card q-mt-md" v-if="ruleDependencies.length > 0">
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
                    <q-item-label caption>{{ dependency.type }} - {{ dependency.description }}</q-item-label>
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
          
          <!-- Template Information -->
          <q-card class="template-card q-mt-md" v-if="rule?.template">
            <q-card-section>
              <h3>Template Information</h3>
              
              <div class="template-info">
                <div class="template-basic">
                  <p><strong>Template:</strong> {{ rule.template.name }}</p>
                  <p><strong>Category:</strong> {{ rule.template.category }}</p>
                  <p><strong>Description:</strong> {{ rule.template.description }}</p>
                </div>
                
                <div class="template-parameters" v-if="rule.template.parameters.length > 0">
                  <h4>Parameters Used:</h4>
                  <q-list>
                    <q-item
                      v-for="param in rule.template.parameters"
                      :key="param.id"
                      dense
                    >
                      <q-item-section>
                        <q-item-label>{{ param.name }}</q-item-label>
                        <q-item-label caption>{{ param.description }}</q-item-label>
                      </q-item-section>
                      <q-item-section side>
                        <q-chip color="grey-3" text-color="grey-8" size="sm">
                          {{ param.type }}
                        </q-chip>
                      </q-item-section>
                    </q-item>
                  </q-list>
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </q-tab-panel>

      <!-- Execution History Tab -->
      <q-tab-panel name="execution" class="q-pa-md">
        <div class="execution-content">
          <div class="execution-filters">
            <div class="row q-gutter-md items-center">
              <div class="col-auto">
                <q-select
                  v-model="executionFilters.timeRange"
                  label="Time Range"
                  :options="timeRangeOptions"
                  outlined
                  dense
                  style="width: 200px"
                />
              </div>
              
              <div class="col-auto">
                <q-select
                  v-model="executionFilters.status"
                  label="Execution Status"
                  :options="executionStatusOptions"
                  outlined
                  dense
                  style="width: 180px"
                />
              </div>
              
              <div class="col-auto">
                <q-btn
                  color="primary"
                  label="Refresh"
                  icon="refresh"
                  @click="fetchExecutionHistory"
                  outline
                  dense
                />
              </div>
            </div>
          </div>
          
          <q-table
            :rows="executionHistory"
            :columns="executionColumns"
            :loading="executionLoading"
            :pagination="executionPagination"
            @request="onExecutionTableRequest"
            row-key="id"
            class="execution-table q-mt-md"
          >
            <template v-slot:body-cell-status="props">
              <q-td :props="props">
                <q-chip
                  :color="getExecutionStatusColor(props.value)"
                  text-color="white"
                  size="sm"
                  :icon="getExecutionStatusIcon(props.value)"
                >
                  {{ props.value }}
                </q-chip>
              </q-td>
            </template>
            
            <template v-slot:body-cell-executionTime="props">
              <q-td :props="props">
                <span :class="{
                  'text-positive': props.value < 100,
                  'text-warning': props.value >= 100 && props.value < 500,
                  'text-negative': props.value >= 500
                }">
                  {{ props.value }}ms
                </span>
              </q-td>
            </template>
            
            <template v-slot:body-cell-actions="props">
              <q-td :props="props">
                <q-btn
                  flat
                  dense
                  icon="visibility"
                  @click="viewExecutionDetails(props.row)"
                  size="sm"
                >
                  <q-tooltip>View Details</q-tooltip>
                </q-btn>
                
                <q-btn
                  flat
                  dense
                  icon="replay"
                  @click="replayExecution(props.row)"
                  size="sm"
                  v-if="canReplay(props.row)"
                >
                  <q-tooltip>Replay Execution</q-tooltip>
                </q-btn>
              </q-td>
            </template>
          </q-table>
        </div>
      </q-tab-panel>

      <!-- Performance Tab -->
      <q-tab-panel name="performance" class="q-pa-md">
        <div class="performance-content">
          <!-- Performance Metrics -->
          <div class="row q-gutter-md">
            <div class="col-md-6 col-12">
              <q-card class="performance-chart-card">
                <q-card-section>
                  <h4>Execution Time Trend</h4>
                  <LineChart
                    :data="performanceChartData.executionTime"
                    :options="executionTimeChartOptions"
                    height="300"
                  />
                </q-card-section>
              </q-card>
            </div>
            
            <div class="col-md-6 col-12">
              <q-card class="performance-chart-card">
                <q-card-section>
                  <h4>Success Rate Trend</h4>
                  <LineChart
                    :data="performanceChartData.successRate"
                    :options="successRateChartOptions"
                    height="300"
                  />
                </q-card-section>
              </q-card>
            </div>
          </div>
          
          <!-- Performance Metrics Table -->
          <q-card class="performance-metrics-card q-mt-md">
            <q-card-section>
              <h4>Performance Metrics Summary</h4>
              
              <q-table
                :rows="performanceMetrics"
                :columns="performanceColumns"
                row-key="period"
                flat
                bordered
              />
            </q-card-section>
          </q-card>
        </div>
      </q-tab-panel>

      <!-- Versions Tab -->
      <q-tab-panel name="versions" class="q-pa-md">
        <div class="versions-content">
          <q-timeline color="primary">
            <q-timeline-entry
              v-for="version in ruleVersions"
              :key="version.id"
              :title="`Version ${version.version}`"
              :subtitle="formatDateTime(version.createdAt)"
              :icon="getVersionIcon(version)"
            >
              <div class="version-details">
                <p><strong>Changes:</strong> {{ version.changeDescription }}</p>
                <p><strong>Modified by:</strong> {{ version.createdBy }}</p>
                
                <div class="version-actions q-mt-sm">
                  <q-btn
                    flat
                    dense
                    icon="visibility"
                    label="View"
                    @click="viewVersion(version)"
                    size="sm"
                  />
                  <q-btn
                    flat
                    dense
                    icon="compare_arrows"
                    label="Compare"
                    @click="compareVersion(version)"
                    size="sm"
                  />
                  <q-btn
                    flat
                    dense
                    icon="restore"
                    label="Restore"
                    @click="restoreVersion(version)"
                    size="sm"
                    v-if="canRestoreVersion(version)"
                  />
                </div>
              </div>
            </q-timeline-entry>
          </q-timeline>
        </div>
      </q-tab-panel>

      <!-- Approvals Tab -->
      <q-tab-panel name="approvals" class="q-pa-md">
        <div class="approvals-content">
          <q-card class="approval-workflow-card">
            <q-card-section>
              <h4>Approval Workflow</h4>
              
              <q-stepper
                v-model="approvalStep"
                color="primary"
                animated
                flat
                bordered
              >
                <q-step
                  :name="1"
                  title="Submission"
                  icon="edit"
                  :done="approvalStep > 1"
                >
                  <div class="approval-step-content">
                    <p>Rule submitted for approval on {{ formatDateTime(rule?.submittedAt) }}</p>
                    <p>Submitted by: {{ rule?.submittedBy }}</p>
                  </div>
                </q-step>

                <q-step
                  :name="2"
                  title="Review"
                  icon="visibility"
                  :done="approvalStep > 2"
                >
                  <div class="approval-step-content">
                    <p v-if="rule?.reviewedAt">
                      Reviewed on {{ formatDateTime(rule?.reviewedAt) }} by {{ rule?.reviewedBy }}
                    </p>
                    <p v-else class="text-grey-6">
                      Pending review...
                    </p>
                  </div>
                </q-step>

                <q-step
                  :name="3"
                  title="Approval"
                  icon="check_circle"
                  :done="approvalStep > 3"
                >
                  <div class="approval-step-content">
                    <p v-if="rule?.approvedAt">
                      Approved on {{ formatDateTime(rule?.approvedAt) }} by {{ rule?.approvedBy }}
                    </p>
                    <p v-else-if="rule?.rejectedAt">
                      Rejected on {{ formatDateTime(rule?.rejectedAt) }} by {{ rule?.rejectedBy }}
                      <br>
                      Reason: {{ rule?.rejectionReason }}
                    </p>
                    <p v-else class="text-grey-6">
                      Pending approval...
                    </p>
                  </div>
                </q-step>
              </q-stepper>
            </q-card-section>
          </q-card>
          
          <!-- Approval Actions -->
          <q-card class="approval-actions-card q-mt-md" v-if="showApprovalActions">
            <q-card-section>
              <h4>Approval Actions</h4>
              
              <div class="approval-form">
                <q-input
                  v-model="approvalComment"
                  label="Comments"
                  type="textarea"
                  rows="3"
                  outlined
                  class="q-mb-md"
                />
                
                <div class="approval-buttons">
                  <q-btn
                    color="positive"
                    icon="check_circle"
                    label="Approve"
                    @click="approveRule"
                    v-if="canApprove"
                  />
                  
                  <q-btn
                    color="negative"
                    icon="cancel"
                    label="Reject"
                    @click="rejectRule"
                    class="q-ml-sm"
                    v-if="canReject"
                  />
                  
                  <q-btn
                    color="warning"
                    icon="edit"
                    label="Request Changes"
                    @click="requestChanges"
                    outline
                    class="q-ml-sm"
                    v-if="canRequestChanges"
                  />
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </q-tab-panel>
    </q-tab-panels>
  </div>
</template>
```

## Component Data Structure

### Rule Detail Data Interface
```typescript
interface RuleDetailData {
  rule: Rule | null
  ruleStats: RuleStatistics
  relatedRules: Rule[]
  ruleDependencies: RuleDependency[]
  executionHistory: RuleExecution[]
  ruleVersions: RuleVersion[]
  performanceMetrics: PerformanceMetric[]
  activeTab: string
  loading: boolean
  
  // Test Dialog
  showTestDialog: boolean
  testResults: RuleTestResult | null
  
  // Execution History
  executionFilters: ExecutionFilters
  executionLoading: boolean
  executionPagination: TablePagination
  
  // Performance Charts
  performanceChartData: {
    executionTime: ChartData
    successRate: ChartData
  }
  
  // Approval
  approvalStep: number
  showApprovalActions: boolean
  approvalComment: string
}

interface RuleStatistics {
  executionCount: number
  successRate: number
  avgExecutionTime: number
  impactValue: number
}

interface RuleDependency {
  id: string
  name: string
  type: 'SERVICE' | 'DATA_SOURCE' | 'EXTERNAL_API'
  status: 'ACTIVE' | 'INACTIVE'
  description: string
}

interface RuleExecution {
  id: string
  timestamp: string
  status: 'SUCCESS' | 'FAILURE' | 'TIMEOUT'
  executionTime: number
  inputData: Record<string, any>
  outputData: Record<string, any>
  errorMessage?: string
}

interface RuleVersion {
  id: string
  version: number
  dslContent: string
  changeDescription: string
  createdAt: string
  createdBy: string
  isActive: boolean
}
```

## Screen Interactions

### Rule Management Actions
- **Edit Rule**: Navigate to edit form
- **Test Rule**: Open test dialog with sample data
- **Duplicate Rule**: Create copy with incremented name
- **Export Rule**: Download as JSON/YAML
- **Share Rule**: Generate shareable link

### Approval Workflow Actions
- **Submit for Approval**: Change status to UNDER_REVIEW
- **Approve Rule**: Approve and activate rule
- **Reject Rule**: Reject with reason
- **Request Changes**: Send back for modifications

### Version Management
- **View Version**: Show version details in modal
- **Compare Versions**: Side-by-side diff view
- **Restore Version**: Revert to previous version

## State Management
```typescript
// Rule Detail Store
export const useRuleDetailStore = defineStore('ruleDetail', () => {
  const rule = ref<Rule | null>(null)
  const ruleStats = ref<RuleStatistics | null>(null)
  const executionHistory = ref<RuleExecution[]>([])
  const loading = ref(false)
  
  const fetchRuleDetail = async (id: string) => {
    loading.value = true
    try {
      const [ruleData, statsData, historyData] = await Promise.all([
        rulesApi.get(id),
        rulesApi.getStatistics(id),
        rulesApi.getExecutionHistory(id)
      ])
      
      rule.value = ruleData
      ruleStats.value = statsData
      executionHistory.value = historyData
    } finally {
      loading.value = false
    }
  }
  
  return {
    rule: readonly(rule),
    ruleStats: readonly(ruleStats),
    executionHistory: readonly(executionHistory),
    loading: readonly(loading),
    fetchRuleDetail
  }
})
```

## Performance Optimizations
- Lazy load tab content on first access
- Cache rule data for 5 minutes
- Virtual scrolling for execution history
- Debounced chart updates
- Optimistic updates for status changes

## Accessibility Features
- Keyboard navigation for tabs
- Screen reader support for charts
- High contrast mode for status indicators
- Focus management for modal dialogs

## Testing Strategy
- Unit tests for all tab components
- Integration tests for API interactions
- E2E tests for complete rule workflows
- Performance tests for large execution history
- Accessibility tests with keyboard navigation
