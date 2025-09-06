# Rules List Screen - Rule Management Overview

## Overview
The Rules List screen provides a comprehensive interface for viewing, filtering, and managing all business rules. It serves as the primary entry point for rule management operations.

## Screen Layout

### Header Section with Filters
```vue
<template>
  <div class="rules-list-header">
    <div class="page-title-section">
      <h1>Rules Management</h1>
      <p class="subtitle">Manage and monitor your business rules</p>
    </div>
    
    <div class="header-actions">
      <q-btn 
        color="primary" 
        icon="add" 
        label="Create New Rule" 
        @click="navigateToCreateRule"
      />
      <q-btn 
        flat 
        icon="download" 
        label="Export" 
        @click="showExportDialog = true"
      />
      <q-btn 
        flat 
        icon="import_export" 
        label="Import" 
        @click="showImportDialog = true"
      />
    </div>
  </div>

  <!-- Advanced Filters -->
  <div class="filters-section">
    <q-expansion-item 
      icon="filter_list" 
      label="Filters" 
      v-model="filtersExpanded"
      class="filters-expansion"
    >
      <q-card class="filters-card">
        <q-card-section>
          <div class="row q-gutter-md">
            <div class="col-md-2 col-sm-6 col-12">
              <q-select
                v-model="filters.status"
                label="Status"
                :options="statusOptions"
                clearable
                multiple
                use-chips
                outlined
                dense
              />
            </div>
            
            <div class="col-md-2 col-sm-6 col-12">
              <q-select
                v-model="filters.priority"
                label="Priority"
                :options="priorityOptions"
                clearable
                multiple
                use-chips
                outlined
                dense
              />
            </div>
            
            <div class="col-md-2 col-sm-6 col-12">
              <q-select
                v-model="filters.category"
                label="Category"
                :options="categoryOptions"
                clearable
                multiple
                use-chips
                outlined
                dense
              />
            </div>
            
            <div class="col-md-2 col-sm-6 col-12">
              <q-input
                v-model="filters.createdBy"
                label="Created By"
                clearable
                outlined
                dense
              />
            </div>
            
            <div class="col-md-2 col-sm-6 col-12">
              <q-input
                v-model="filters.dateRange"
                label="Created Date"
                mask="##/##/#### - ##/##/####"
                outlined
                dense
              >
                <template v-slot:append>
                  <q-icon name="event" class="cursor-pointer">
                    <q-popup-proxy>
                      <q-date 
                        v-model="filters.dateRange" 
                        range 
                        @update:model-value="applyFilters"
                      />
                    </q-popup-proxy>
                  </q-icon>
                </template>
              </q-input>
            </div>
            
            <div class="col-md-2 col-sm-6 col-12">
              <div class="filter-actions">
                <q-btn 
                  color="primary" 
                  label="Apply" 
                  @click="applyFilters"
                  dense
                />
                <q-btn 
                  flat 
                  label="Clear" 
                  @click="clearFilters"
                  dense
                />
              </div>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </q-expansion-item>
  </div>
</template>
```

### Search and Quick Filters
```vue
<template>
  <div class="search-section">
    <div class="row q-gutter-md items-center">
      <div class="col-md-4 col-sm-6 col-12">
        <q-input
          v-model="searchQuery"
          label="Search rules..."
          outlined
          dense
          clearable
          debounce="300"
          @update:model-value="onSearch"
        >
          <template v-slot:prepend>
            <q-icon name="search" />
          </template>
        </q-input>
      </div>
      
      <div class="col-auto">
        <q-chip 
          v-for="quickFilter in quickFilters"
          :key="quickFilter.key"
          :color="activeQuickFilter === quickFilter.key ? 'primary' : 'grey-3'"
          :text-color="activeQuickFilter === quickFilter.key ? 'white' : 'grey-8'"
          clickable
          @click="applyQuickFilter(quickFilter.key)"
          class="q-mr-sm"
        >
          {{ quickFilter.label }} ({{ quickFilter.count }})
        </q-chip>
      </div>
      
      <div class="col-auto">
        <q-btn-dropdown
          color="grey-7"
          icon="sort"
          label="Sort"
          outline
          dense
        >
          <q-list>
            <q-item
              v-for="sortOption in sortOptions"
              :key="sortOption.value"
              clickable
              @click="applySorting(sortOption.value)"
              :class="{ 'bg-grey-2': currentSort === sortOption.value }"
            >
              <q-item-section>
                <q-item-label>{{ sortOption.label }}</q-item-label>
              </q-item-section>
              <q-item-section side v-if="currentSort === sortOption.value">
                <q-icon name="check" color="primary" />
              </q-item-section>
            </q-item>
          </q-list>
        </q-btn-dropdown>
      </div>
    </div>
  </div>
</template>
```

### Rules Data Table
```vue
<template>
  <div class="rules-table-section">
    <q-table
      :rows="filteredRules"
      :columns="tableColumns"
      :loading="loading"
      :pagination="tablePagination"
      @request="onTableRequest"
      row-key="id"
      selection="multiple"
      v-model:selected="selectedRules"
      class="rules-table"
    >
      <!-- Table Header -->
      <template v-slot:top>
        <div class="table-header-actions" v-if="selectedRules.length > 0">
          <span class="selected-count">{{ selectedRules.length }} rule(s) selected</span>
          <q-btn
            color="primary"
            icon="approval"
            label="Bulk Approve"
            @click="showBulkApprovalDialog = true"
            :disable="!canBulkApprove"
            outline
            dense
          />
          <q-btn
            color="negative"
            icon="delete"
            label="Bulk Delete"
            @click="showBulkDeleteDialog = true"
            outline
            dense
          />
          <q-btn
            color="grey"
            icon="more_vert"
            label="More Actions"
            outline
            dense
          >
            <q-menu>
              <q-list>
                <q-item clickable @click="bulkActivate">
                  <q-item-section avatar>
                    <q-icon name="play_arrow" />
                  </q-item-section>
                  <q-item-section>Activate</q-item-section>
                </q-item>
                <q-item clickable @click="bulkDeactivate">
                  <q-item-section avatar>
                    <q-icon name="pause" />
                  </q-item-section>
                  <q-item-section>Deactivate</q-item-section>
                </q-item>
                <q-item clickable @click="bulkExport">
                  <q-item-section avatar>
                    <q-icon name="download" />
                  </q-item-section>
                  <q-item-section>Export Selected</q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </div>
      </template>

      <!-- Custom Column Renderers -->
      <template v-slot:body-cell-name="props">
        <q-td :props="props">
          <div class="rule-name-cell">
            <div class="rule-name">
              <q-btn
                flat
                dense
                :label="props.value"
                @click="navigateToRuleDetail(props.row.id)"
                class="rule-name-link"
              />
            </div>
            <div class="rule-description" v-if="props.row.description">
              {{ truncateText(props.row.description, 60) }}
            </div>
          </div>
        </q-td>
      </template>

      <template v-slot:body-cell-status="props">
        <q-td :props="props">
          <q-chip
            :color="getStatusColor(props.value)"
            text-color="white"
            size="sm"
            :icon="getStatusIcon(props.value)"
          >
            {{ props.value }}
          </q-chip>
        </q-td>
      </template>

      <template v-slot:body-cell-priority="props">
        <q-td :props="props">
          <q-chip
            :color="getPriorityColor(props.value)"
            text-color="white"
            size="sm"
          >
            {{ props.value }}
          </q-chip>
        </q-td>
      </template>

      <template v-slot:body-cell-tags="props">
        <q-td :props="props">
          <div class="tags-cell">
            <q-chip
              v-for="tag in props.value.slice(0, 2)"
              :key="tag"
              color="grey-3"
              text-color="grey-8"
              size="sm"
              dense
            >
              {{ tag }}
            </q-chip>
            <q-chip
              v-if="props.value.length > 2"
              color="grey-3"
              text-color="grey-8"
              size="sm"
              dense
            >
              +{{ props.value.length - 2 }}
            </q-chip>
          </div>
        </q-td>
      </template>

      <template v-slot:body-cell-actions="props">
        <q-td :props="props">
          <div class="action-buttons">
            <q-btn
              flat
              dense
              icon="visibility"
              @click="navigateToRuleDetail(props.row.id)"
              size="sm"
            >
              <q-tooltip>View Details</q-tooltip>
            </q-btn>
            
            <q-btn
              flat
              dense
              icon="edit"
              @click="navigateToRuleEdit(props.row.id)"
              size="sm"
              v-if="canEditRule(props.row)"
            >
              <q-tooltip>Edit Rule</q-tooltip>
            </q-btn>
            
            <q-btn
              flat
              dense
              icon="content_copy"
              @click="duplicateRule(props.row)"
              size="sm"
            >
              <q-tooltip>Duplicate Rule</q-tooltip>
            </q-btn>
            
            <q-btn
              flat
              dense
              icon="play_arrow"
              @click="testRule(props.row)"
              size="sm"
              color="positive"
            >
              <q-tooltip>Test Rule</q-tooltip>
            </q-btn>
            
            <q-btn
              flat
              dense
              icon="more_vert"
              size="sm"
            >
              <q-menu>
                <q-list>
                  <q-item clickable @click="approveRule(props.row)" v-if="canApprove(props.row)">
                    <q-item-section avatar>
                      <q-icon name="check_circle" color="positive" />
                    </q-item-section>
                    <q-item-section>Approve</q-item-section>
                  </q-item>
                  
                  <q-item clickable @click="activateRule(props.row)" v-if="canActivate(props.row)">
                    <q-item-section avatar>
                      <q-icon name="play_arrow" color="primary" />
                    </q-item-section>
                    <q-item-section>Activate</q-item-section>
                  </q-item>
                  
                  <q-item clickable @click="deactivateRule(props.row)" v-if="canDeactivate(props.row)">
                    <q-item-section avatar>
                      <q-icon name="pause" color="warning" />
                    </q-item-section>
                    <q-item-section>Deactivate</q-item-section>
                  </q-item>
                  
                  <q-separator />
                  
                  <q-item clickable @click="viewRuleHistory(props.row)">
                    <q-item-section avatar>
                      <q-icon name="history" />
                    </q-item-section>
                    <q-item-section>View History</q-item-section>
                  </q-item>
                  
                  <q-item clickable @click="exportRule(props.row)">
                    <q-item-section avatar>
                      <q-icon name="download" />
                    </q-item-section>
                    <q-item-section>Export</q-item-section>
                  </q-item>
                  
                  <q-separator />
                  
                  <q-item clickable @click="deleteRule(props.row)" v-if="canDelete(props.row)">
                    <q-item-section avatar>
                      <q-icon name="delete" color="negative" />
                    </q-item-section>
                    <q-item-section>Delete</q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-btn>
          </div>
        </q-td>
      </template>

      <!-- Empty State -->
      <template v-slot:no-data>
        <div class="empty-state">
          <q-icon name="rule" size="4rem" color="grey-4" />
          <div class="empty-title">No rules found</div>
          <div class="empty-subtitle">
            {{ searchQuery ? 'Try adjusting your search criteria' : 'Create your first rule to get started' }}
          </div>
          <q-btn
            color="primary"
            icon="add"
            label="Create New Rule"
            @click="navigateToCreateRule"
            class="q-mt-md"
            v-if="!searchQuery"
          />
        </div>
      </template>
    </q-table>
  </div>
</template>
```

## Component Data Structure

### Rules List Data Interface
```typescript
interface RulesListData {
  rules: Rule[]
  filteredRules: Rule[]
  selectedRules: Rule[]
  loading: boolean
  searchQuery: string
  activeQuickFilter: string | null
  currentSort: string
  filtersExpanded: boolean
  
  filters: {
    status: RuleStatus[]
    priority: Priority[]
    category: string[]
    createdBy: string
    dateRange: string
  }
  
  tablePagination: {
    page: number
    rowsPerPage: number
    rowsNumber: number
    sortBy: string
    descending: boolean
  }
}

interface QuickFilter {
  key: string
  label: string
  count: number
  filter: (rules: Rule[]) => Rule[]
}

interface SortOption {
  label: string
  value: string
  field: keyof Rule
  direction: 'asc' | 'desc'
}
```

### Table Columns Configuration
```typescript
const tableColumns = [
  {
    name: 'name',
    label: 'Rule Name',
    field: 'name',
    align: 'left',
    sortable: true,
    style: 'width: 300px'
  },
  {
    name: 'status',
    label: 'Status',
    field: 'status',
    align: 'center',
    sortable: true,
    style: 'width: 120px'
  },
  {
    name: 'priority',
    label: 'Priority',
    field: 'priority',
    align: 'center',
    sortable: true,
    style: 'width: 100px'
  },
  {
    name: 'category',
    label: 'Category',
    field: 'category',
    align: 'left',
    sortable: true,
    style: 'width: 150px'
  },
  {
    name: 'createdBy',
    label: 'Created By',
    field: 'createdBy',
    align: 'left',
    sortable: true,
    style: 'width: 140px'
  },
  {
    name: 'createdAt',
    label: 'Created Date',
    field: 'createdAt',
    format: (val: string) => formatDate(val),
    align: 'center',
    sortable: true,
    style: 'width: 120px'
  },
  {
    name: 'tags',
    label: 'Tags',
    field: 'tags',
    align: 'left',
    style: 'width: 200px'
  },
  {
    name: 'actions',
    label: 'Actions',
    field: 'actions',
    align: 'center',
    style: 'width: 120px'
  }
]
```

## Screen Interactions

### Navigation Actions
- **Create Rule**: Navigate to `/rules/create`
- **View Rule**: Navigate to `/rules/:id`
- **Edit Rule**: Navigate to `/rules/:id/edit`
- **Rule History**: Navigate to `/rules/:id/history`

### Bulk Operations
- **Bulk Approve**: Approve multiple selected rules
- **Bulk Activate/Deactivate**: Change status of multiple rules
- **Bulk Delete**: Delete multiple rules with confirmation
- **Bulk Export**: Export selected rules to CSV/JSON

### Real-time Features
- Rule status updates via WebSocket
- Real-time rule execution metrics
- Live notifications for rule changes

## State Management
```typescript
// Rules List Store
export const useRulesListStore = defineStore('rulesList', () => {
  const rules = ref<Rule[]>([])
  const filters = ref<RulesFilters>({})
  const pagination = ref<TablePagination>({})
  const loading = ref(false)
  
  const filteredRules = computed(() => {
    let result = rules.value
    
    // Apply search filter
    if (filters.value.search) {
      result = result.filter(rule => 
        rule.name.toLowerCase().includes(filters.value.search!.toLowerCase()) ||
        rule.description.toLowerCase().includes(filters.value.search!.toLowerCase())
      )
    }
    
    // Apply status filter
    if (filters.value.status?.length) {
      result = result.filter(rule => filters.value.status!.includes(rule.status))
    }
    
    // Apply other filters...
    
    return result
  })
  
  const fetchRules = async (request: ListRulesRequest) => {
    loading.value = true
    try {
      const response = await rulesApi.list(request)
      rules.value = response.data
      pagination.value = response.pagination
    } finally {
      loading.value = false
    }
  }
  
  return {
    rules: readonly(rules),
    filteredRules,
    filters,
    pagination,
    loading: readonly(loading),
    fetchRules
  }
})
```

## Performance Optimizations
- Virtual scrolling for large rule lists
- Debounced search input (300ms)
- Lazy loading of rule details
- Optimistic updates for status changes
- Client-side pagination with server-side sorting

## Accessibility Features
- Keyboard navigation for table rows
- Screen reader support for status chips
- High contrast mode for priority indicators
- Focus management for bulk operations

## Testing Strategy
- Unit tests for filtering and sorting logic
- Integration tests for API interactions
- E2E tests for complete user workflows
- Performance tests for large datasets
- Accessibility tests with keyboard navigation
