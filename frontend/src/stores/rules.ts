import { defineStore } from 'pinia'
import { ref, computed, watch, readonly } from 'vue'
import { rulesApi } from '@/api/rules'
import { useNotificationStore } from './notifications'
import type { 
  Rule, 
  RuleListParams, 
  RuleFilters,
  CreateRuleRequest,
  UpdateRuleRequest,
  RuleMetrics,
  PaginationInfo
} from '@/types'

export const useRulesStore = defineStore('rules', () => {
  // State
  const rules = ref<Rule[]>([])
  const currentRule = ref<Rule | null>(null)
  const loading = ref(false)
  const saving = ref(false)
  const error = ref<string | null>(null)
  
  // Pagination
  const pagination = ref<PaginationInfo>({
    page: 1,
    limit: 20,
    total: 0,
    totalPages: 0
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
        rule.created_by.toLowerCase().includes(filters.value.createdBy.toLowerCase())
      )
    }
    
    // Apply date range filter
    if (filters.value.dateRange) {
      const { start, end } = filters.value.dateRange
      result = result.filter(rule => {
        const createdDate = new Date(rule.created_at)
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
        limit: pagination.value.limit,
        sortBy: 'created_at',
        sortOrder: 'desc',
        filters: filters.value,
        search: searchQuery.value,
        ...params
      }
      
      const response = await rulesApi.getRules(requestParams)
      
      rules.value = response.data || []
      pagination.value = response.pagination || pagination.value
      
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
      dsl_content: originalRule.dsl_content,
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
  
  const validateRule = async (request: { dsl_content: string; context: Record<string, any>; rule_category: string }) => {
    try {
      const response = await rulesApi.validateRule({
        dsl_content: request.dsl_content,
        context: request.context,
        rule_category: request.rule_category
      })
      return response.data
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to validate rule'
      notificationStore.showError('Rule validation failed', message)
      throw err
    }
  }

  const bulkActivate = async (ruleIds: string[]) => {
    saving.value = true
    error.value = null
    
    try {
      await Promise.all(ruleIds.map(id => rulesApi.activateRule(id)))
      
      // Update local state
      rules.value = rules.value.map(rule => 
        ruleIds.includes(rule.id) ? { ...rule, status: 'ACTIVE' } : rule
      )
      
      notificationStore.showSuccess(`${ruleIds.length} rules activated successfully`)
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to activate rules'
      error.value = message
      notificationStore.showError('Failed to activate rules', message)
      throw err
    } finally {
      saving.value = false
    }
  }

  const bulkDeactivate = async (ruleIds: string[]) => {
    saving.value = true
    error.value = null
    
    try {
      await Promise.all(ruleIds.map(id => rulesApi.deactivateRule(id)))
      
      // Update local state
      rules.value = rules.value.map(rule => 
        ruleIds.includes(rule.id) ? { ...rule, status: 'INACTIVE' } : rule
      )
      
      notificationStore.showSuccess(`${ruleIds.length} rules deactivated successfully`)
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to deactivate rules'
      error.value = message
      notificationStore.showError('Failed to deactivate rules', message)
      throw err
    } finally {
      saving.value = false
    }
  }

  const bulkDelete = async (ruleIds: string[]) => {
    saving.value = true
    error.value = null
    
    try {
      await Promise.all(ruleIds.map(id => rulesApi.deleteRule(id)))
      
      // Remove from local state
      rules.value = rules.value.filter(rule => !ruleIds.includes(rule.id))
      
      // Remove metrics cache
      ruleIds.forEach(id => ruleMetrics.value.delete(id))
      
      notificationStore.showSuccess(`${ruleIds.length} rules deleted successfully`)
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to delete rules'
      error.value = message
      notificationStore.showError('Failed to delete rules', message)
      throw err
    } finally {
      saving.value = false
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
      limit: 20,
      total: 0,
      totalPages: 0
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
    activateRule,
    deactivateRule,
    fetchRuleMetrics,
    validateRule,
    bulkActivate,
    bulkDeactivate,
    bulkDelete,
    
    // Utility actions
    clearError,
    clearCurrentRule,
    updateFilters,
    updateSearch,
    reset
  }
})
