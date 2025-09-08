<template>
  <q-page class="rules-list-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="page-header-content">
        <div class="page-title-section">
          <h1 class="page-title">Rules Management</h1>
          <p class="page-subtitle">
            Create, manage, and monitor your business rules
          </p>
        </div>
        
        <div class="page-actions">
          <q-btn
            color="primary"
            icon="add"
            label="Create Rule"
            @click="$router.push('/rules/create')"
          />
          <q-btn
            color="secondary"
            icon="upload"
            label="Import Rules"
            @click="showImportDialog = true"
          />
        </div>
      </div>
    </div>

    <!-- Filters and Search -->
    <q-card class="filters-card">
      <q-card-section>
        <div class="filters-row">
          <!-- Search -->
          <div class="search-section">
            <q-input
              v-model="searchQuery"
              placeholder="Search rules..."
              outlined
              dense
              class="search-input"
              @input="handleSearch"
            >
              <template v-slot:prepend>
                <q-icon name="search" />
              </template>
              <template v-slot:append v-if="searchQuery">
                <q-icon 
                  name="clear" 
                  class="cursor-pointer" 
                  @click="clearSearch"
                />
              </template>
            </q-input>
          </div>

          <!-- Filters -->
          <div class="filters-section">
            <q-select
              v-model="statusFilter"
              :options="statusOptions"
              label="Status"
              outlined
              dense
              multiple
              emit-value
              map-options
              class="filter-select"
              @update:model-value="handleFilterChange"
            />
            
            <q-select
              v-model="priorityFilter"
              :options="priorityOptions"
              label="Priority"
              outlined
              dense
              multiple
              emit-value
              map-options
              class="filter-select"
              @update:model-value="handleFilterChange"
            />
            
            <q-select
              v-model="categoryFilter"
              :options="categoryOptions"
              label="Category"
              outlined
              dense
              emit-value
              map-options
              class="filter-select"
              @update:model-value="handleFilterChange"
            />
          </div>

          <!-- Actions -->
          <div class="actions-section">
            <q-btn
              flat
              dense
              icon="refresh"
              @click="refreshRules"
              :loading="loading"
            >
              <q-tooltip>Refresh</q-tooltip>
            </q-btn>
            
            <q-btn
              flat
              dense
              icon="download"
              @click="exportRules"
            >
              <q-tooltip>Export</q-tooltip>
            </q-btn>
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Rules Table -->
    <q-card class="rules-table-card">
      <q-card-section>
        <q-table
          :rows="filteredRules"
          :columns="columns"
          :loading="loading"
          :pagination="pagination"
          row-key="id"
          selection="multiple"
          v-model:selected="selectedRules"
          class="rules-table"
          @request="onRequest"
        >
          <!-- Status Column -->
          <template v-slot:body-cell-status="props">
            <q-td :props="props">
              <q-chip
                :color="getStatusColor(props.value)"
                text-color="white"
                dense
                :icon="getStatusIcon(props.value)"
              >
                {{ props.value }}
              </q-chip>
            </q-td>
          </template>

          <!-- Priority Column -->
          <template v-slot:body-cell-priority="props">
            <q-td :props="props">
              <div class="priority-cell">
                <q-icon
                  :name="getPriorityIcon(props.value)"
                  :color="getPriorityColor(props.value)"
                  size="sm"
                />
                <span class="priority-text">{{ props.value }}</span>
              </div>
            </q-td>
          </template>

          <!-- Created At Column -->
          <template v-slot:body-cell-created_at="props">
            <q-td :props="props">
              <div class="date-cell">
                <div class="date-text">{{ formatDate(props.value) }}</div>
                <div class="time-text">{{ formatTime(props.value) }}</div>
              </div>
            </q-td>
          </template>

          <!-- Actions Column -->
          <template v-slot:body-cell-actions="props">
            <q-td :props="props">
              <q-btn-group flat>
                <q-btn
                  flat
                  dense
                  icon="visibility"
                  color="primary"
                  @click="viewRule(props.row)"
                >
                  <q-tooltip>View Details</q-tooltip>
                </q-btn>
                
                <q-btn
                  flat
                  dense
                  icon="edit"
                  color="secondary"
                  @click="editRule(props.row)"
                >
                  <q-tooltip>Edit Rule</q-tooltip>
                </q-btn>
                
                <q-btn
                  flat
                  dense
                  :icon="props.row.status === 'ACTIVE' ? 'pause' : 'play_arrow'"
                  :color="props.row.status === 'ACTIVE' ? 'orange' : 'positive'"
                  @click="toggleRuleStatus(props.row)"
                >
                  <q-tooltip>
                    {{ props.row.status === 'ACTIVE' ? 'Deactivate' : 'Activate' }}
                  </q-tooltip>
                </q-btn>
                
                <q-btn
                  flat
                  dense
                  icon="more_vert"
                  color="grey-7"
                >
                  <q-menu>
                    <q-list dense>
                      <q-item clickable @click="duplicateRule(props.row)">
                        <q-item-section avatar>
                          <q-icon name="content_copy" />
                        </q-item-section>
                        <q-item-section>Duplicate</q-item-section>
                      </q-item>
                      
                      <q-item clickable @click="testRule(props.row)">
                        <q-item-section avatar>
                          <q-icon name="play_circle" />
                        </q-item-section>
                        <q-item-section>Test Rule</q-item-section>
                      </q-item>
                      
                      <q-item clickable @click="viewHistory(props.row)">
                        <q-item-section avatar>
                          <q-icon name="history" />
                        </q-item-section>
                        <q-item-section>View History</q-item-section>
                      </q-item>
                      
                      <q-separator />
                      
                      <q-item 
                        clickable 
                        @click="deleteRule(props.row)"
                        class="text-negative"
                      >
                        <q-item-section avatar>
                          <q-icon name="delete" />
                        </q-item-section>
                        <q-item-section>Delete</q-item-section>
                      </q-item>
                    </q-list>
                  </q-menu>
                </q-btn>
              </q-btn-group>
            </q-td>
          </template>

          <!-- Bulk Actions -->
          <template v-slot:top-selection="scope">
            <q-btn-group class="q-mr-md">
              <q-btn
                color="positive"
                icon="play_arrow"
                label="Activate Selected"
                @click="bulkActivate"
                :disable="selectedRules.length === 0"
              />
              <q-btn
                color="orange"
                icon="pause"
                label="Deactivate Selected"
                @click="bulkDeactivate"
                :disable="selectedRules.length === 0"
              />
              <q-btn
                color="negative"
                icon="delete"
                label="Delete Selected"
                @click="bulkDelete"
                :disable="selectedRules.length === 0"
              />
            </q-btn-group>
            
            <div class="text-subtitle2">
              {{ scope.selectedRows.length }} rule(s) selected
            </div>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- Import Dialog -->
    <q-dialog v-model="showImportDialog">
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">Import Rules</div>
        </q-card-section>
        
        <q-card-section>
          <q-file
            v-model="importFile"
            label="Select rules file"
            accept=".json,.csv,.xlsx"
            outlined
            dense
          >
            <template v-slot:prepend>
              <q-icon name="attach_file" />
            </template>
          </q-file>
        </q-card-section>
        
        <q-card-actions align="right">
          <q-btn flat label="Cancel" v-close-popup />
          <q-btn 
            color="primary" 
            label="Import" 
            @click="handleImport"
            :disable="!importFile"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { useRulesStore } from '@/stores/rules'
import { useNotificationStore } from '@/stores/notifications'
import { formatDistanceToNow } from 'date-fns'
import type { Rule, RuleStatus, Priority } from '@/types'

const router = useRouter()
const $q = useQuasar()
const rulesStore = useRulesStore()
const notificationStore = useNotificationStore()

// Reactive state
const searchQuery = ref('')
const statusFilter = ref<RuleStatus[]>([])
const priorityFilter = ref<Priority[]>([])
const categoryFilter = ref('')
const selectedRules = ref<Rule[]>([])
const showImportDialog = ref(false)
const importFile = ref<File | null>(null)

// Computed properties
const loading = computed(() => rulesStore.loading)
const filteredRules = computed(() => rulesStore.filteredRules)
const pagination = computed(() => rulesStore.pagination)

const statusOptions = [
  { label: 'Draft', value: 'DRAFT' },
  { label: 'Under Review', value: 'UNDER_REVIEW' },
  { label: 'Approved', value: 'APPROVED' },
  { label: 'Active', value: 'ACTIVE' },
  { label: 'Inactive', value: 'INACTIVE' },
  { label: 'Deprecated', value: 'DEPRECATED' }
]

const priorityOptions = [
  { label: 'Low', value: 'LOW' },
  { label: 'Medium', value: 'MEDIUM' },
  { label: 'High', value: 'HIGH' },
  { label: 'Critical', value: 'CRITICAL' }
]

const categoryOptions = [
  { label: 'All Categories', value: '' },
  { label: 'Promotions', value: 'PROMOTIONS' },
  { label: 'Coupons', value: 'COUPONS' },
  { label: 'Loyalty', value: 'LOYALTY' },
  { label: 'Taxes', value: 'TAXES' },
  { label: 'Payments', value: 'PAYMENTS' }
]

const columns = [
  {
    name: 'name',
    label: 'Rule Name',
    field: 'name',
    align: 'left',
    sortable: true
  },
  {
    name: 'category',
    label: 'Category',
    field: 'category',
    align: 'left',
    sortable: true
  },
  {
    name: 'status',
    label: 'Status',
    field: 'status',
    align: 'center',
    sortable: true
  },
  {
    name: 'priority',
    label: 'Priority',
    field: 'priority',
    align: 'center',
    sortable: true
  },
  {
    name: 'created_at',
    label: 'Created',
    field: 'created_at',
    align: 'center',
    sortable: true
  },
  {
    name: 'created_by',
    label: 'Created By',
    field: 'created_by',
    align: 'center',
    sortable: true
  },
  {
    name: 'actions',
    label: 'Actions',
    field: 'actions',
    align: 'center'
  }
]

// Methods
const handleSearch = () => {
  rulesStore.updateSearch(searchQuery.value)
}

const clearSearch = () => {
  searchQuery.value = ''
  rulesStore.updateSearch('')
}

const handleFilterChange = () => {
  rulesStore.updateFilters({
    status: statusFilter.value,
    priority: priorityFilter.value,
    category: categoryFilter.value
  })
}

const refreshRules = () => {
  rulesStore.fetchRules()
}

const exportRules = () => {
  // Implement export functionality
  notificationStore.showInfo('Export', 'Export functionality will be implemented')
}

const onRequest = (props: any) => {
  rulesStore.fetchRules({
    page: props.pagination.page,
    limit: props.pagination.rowsPerPage,
    sortBy: props.pagination.sortBy,
    sortOrder: props.pagination.descending ? 'desc' : 'asc'
  })
}

const viewRule = (rule: Rule) => {
  router.push(`/rules/${rule.id}`)
}

const editRule = (rule: Rule) => {
  router.push(`/rules/${rule.id}/edit`)
}

const toggleRuleStatus = async (rule: Rule) => {
  try {
    if (rule.status === 'ACTIVE') {
      await rulesStore.deactivateRule(rule.id)
    } else {
      await rulesStore.activateRule(rule.id)
    }
  } catch (error) {
    // Error handling is done in the store
  }
}

const duplicateRule = async (rule: Rule) => {
  try {
    await rulesStore.duplicateRule(rule.id)
  } catch (error) {
    // Error handling is done in the store
  }
}

const testRule = (rule: Rule) => {
  router.push(`/rules/${rule.id}?action=test`)
}

const viewHistory = (rule: Rule) => {
  router.push(`/rules/${rule.id}?tab=history`)
}

const deleteRule = (rule: Rule) => {
  $q.dialog({
    title: 'Confirm Delete',
    message: `Are you sure you want to delete the rule "${rule.name}"?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await rulesStore.deleteRule(rule.id)
    } catch (error) {
      // Error handling is done in the store
    }
  })
}

const bulkActivate = async () => {
  try {
    const ruleIds = selectedRules.value.map(rule => rule.id)
    await rulesStore.bulkActivate(ruleIds)
    selectedRules.value = []
  } catch (error) {
    // Error handling is done in the store
  }
}

const bulkDeactivate = async () => {
  try {
    const ruleIds = selectedRules.value.map(rule => rule.id)
    await rulesStore.bulkDeactivate(ruleIds)
    selectedRules.value = []
  } catch (error) {
    // Error handling is done in the store
  }
}

const bulkDelete = () => {
  $q.dialog({
    title: 'Confirm Bulk Delete',
    message: `Are you sure you want to delete ${selectedRules.value.length} selected rules?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      const ruleIds = selectedRules.value.map(rule => rule.id)
      await rulesStore.bulkDelete(ruleIds)
      selectedRules.value = []
    } catch (error) {
      // Error handling is done in the store
    }
  })
}

const handleImport = () => {
  if (importFile.value) {
    // Implement import functionality
    notificationStore.showInfo('Import', 'Import functionality will be implemented')
    showImportDialog.value = false
    importFile.value = null
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

const getPriorityIcon = (priority: string) => {
  const icons: Record<string, string> = {
    'LOW': 'keyboard_arrow_down',
    'MEDIUM': 'remove',
    'HIGH': 'keyboard_arrow_up',
    'CRITICAL': 'priority_high'
  }
  return icons[priority] || 'remove'
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

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

const formatTime = (dateString: string) => {
  return formatDistanceToNow(new Date(dateString), { addSuffix: true })
}

// Lifecycle
onMounted(() => {
  rulesStore.fetchRules()
})
</script>

<style lang="scss" scoped>
.rules-list-page {
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
  flex: 1;
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

.filters-card {
  margin: 0 24px 24px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.filters-row {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.search-section {
  flex: 1;
  min-width: 300px;
}

.filters-section {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.filter-select {
  min-width: 150px;
}

.actions-section {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.rules-table-card {
  margin: 0 24px 24px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.rules-table {
  :deep(.q-table__top) {
    padding: 16px;
  }
  
  :deep(.q-table__bottom) {
    padding: 16px;
  }
}

.priority-cell {
  display: flex;
  align-items: center;
  gap: 6px;
}

.priority-text {
  font-size: 14px;
  font-weight: 500;
}

.date-cell {
  text-align: center;
}

.date-text {
  font-size: 14px;
  font-weight: 500;
  color: var(--carrefour-gray-800);
}

.time-text {
  font-size: 12px;
  color: var(--carrefour-gray-600);
}

// Responsive design
@media (max-width: 768px) {
  .page-header {
    padding: 16px;
  }
  
  .page-header-content {
    flex-direction: column;
    gap: 16px;
  }
  
  .page-actions {
    width: 100%;
    justify-content: stretch;
    
    .q-btn {
      flex: 1;
    }
  }
  
  .filters-card,
  .rules-table-card {
    margin: 0 16px 16px;
  }
  
  .filters-row {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
  
  .search-section {
    min-width: auto;
  }
  
  .filters-section {
    justify-content: stretch;
    
    .filter-select {
      flex: 1;
      min-width: auto;
    }
  }
  
  .actions-section {
    justify-content: center;
  }
}
</style>
