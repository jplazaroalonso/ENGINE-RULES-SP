import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useRulesStore } from '../rules'
import { mockApiResponse, mockApiError, mockRule, generateRules } from '../../test/utils/test-utils'

// Mock the API client
vi.mock('@/api/rules', () => ({
  rulesApi: {
    getRules: vi.fn(),
    getRule: vi.fn(),
    createRule: vi.fn(),
    updateRule: vi.fn(),
    deleteRule: vi.fn(),
    activateRule: vi.fn(),
    deactivateRule: vi.fn(),
    validateRule: vi.fn(),
    submitForApproval: vi.fn()
  }
}))

describe('Rules Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should initialize with default state', () => {
    const store = useRulesStore()
    
    expect(store.rules).toEqual([])
    expect(store.currentRule).toBeNull()
    expect(store.loading).toBe(false)
    expect(store.saving).toBe(false)
    expect(store.error).toBeNull()
    expect(store.pagination).toEqual({
      page: 1,
      limit: 20,
      total: 0,
      totalPages: 0
    })
  })

  it('should fetch rules successfully', async () => {
    const store = useRulesStore()
    const mockRules = generateRules(3)
    const response = {
      data: mockRules,
      pagination: {
        page: 1,
        limit: 20,
        total: 3,
        totalPages: 1
      }
    }

    const { rulesApi } = await import('@/api/rules')
    vi.mocked(rulesApi.getRules).mockResolvedValueOnce(mockApiResponse(response))

    await store.fetchRules()

    expect(store.rules).toEqual(mockRules)
    expect(store.pagination).toEqual(response.pagination)
    expect(store.loading).toBe(false)
    expect(store.error).toBeNull()
  })

  it('should handle fetch rules error', async () => {
    const store = useRulesStore()
    const error = mockApiError('Failed to fetch rules', 500)

    const { rulesApi } = await import('@/api/rules')
    vi.mocked(rulesApi.getRules).mockRejectedValueOnce(error)

    await expect(store.fetchRules()).rejects.toThrow('Failed to fetch rules')
    expect(store.error).toBe('Failed to fetch rules')
    expect(store.loading).toBe(false)
  })

  it('should fetch single rule successfully', async () => {
    const store = useRulesStore()
    const response = { data: mockRule }

    const { rulesApi } = await import('@/api/rules')
    vi.mocked(rulesApi.getRule).mockResolvedValueOnce(mockApiResponse(response))

    const result = await store.fetchRule('1')

    expect(result).toEqual(mockRule)
    expect(store.currentRule).toEqual(mockRule)
    expect(store.loading).toBe(false)
  })

  it('should create rule successfully', async () => {
    const store = useRulesStore()
    const createRequest = {
      name: 'New Rule',
      description: 'A new rule',
      dsl_content: 'IF condition THEN action',
      category: 'PROMOTIONS',
      priority: 'MEDIUM',
      tags: ['new']
    }
    const response = { data: { ...mockRule, ...createRequest } }

    const { rulesApi } = await import('@/api/rules')
    vi.mocked(rulesApi.createRule).mockResolvedValueOnce(mockApiResponse(response))

    const result = await store.createRule(createRequest)

    expect(result).toEqual(response.data)
    expect(store.rules).toContain(response.data)
    expect(store.currentRule).toEqual(response.data)
    expect(store.saving).toBe(false)
  })

  it('should update rule successfully', async () => {
    const store = useRulesStore()
    const updateRequest = {
      name: 'Updated Rule',
      description: 'An updated rule'
    }
    const updatedRule = { ...mockRule, ...updateRequest }
    const response = { data: updatedRule }

    // Add rule to store first
    store.rules = [mockRule]

    const { rulesApi } = await import('@/api/rules')
    vi.mocked(rulesApi.updateRule).mockResolvedValueOnce(mockApiResponse(response))

    const result = await store.updateRule('1', updateRequest)

    expect(result).toEqual(updatedRule)
    expect(store.rules[0]).toEqual(updatedRule)
    expect(store.saving).toBe(false)
  })

  it('should delete rule successfully', async () => {
    const store = useRulesStore()
    
    // Add rule to store first
    store.rules = [mockRule]

    const { rulesApi } = await import('@/api/rules')
    vi.mocked(rulesApi.deleteRule).mockResolvedValueOnce(mockApiResponse({}))

    await store.deleteRule('1')

    expect(store.rules).toEqual([])
    expect(store.saving).toBe(false)
  })

  it('should activate rule successfully', async () => {
    const store = useRulesStore()
    const inactiveRule = { ...mockRule, status: 'INACTIVE' as const }
    
    // Add rule to store first
    store.rules = [inactiveRule]

    const { rulesApi } = await import('@/api/rules')
    vi.mocked(rulesApi.activateRule).mockResolvedValueOnce(mockApiResponse({}))

    await store.activateRule('1')

    expect(store.rules[0].status).toBe('ACTIVE')
    expect(store.saving).toBe(false)
  })

  it('should deactivate rule successfully', async () => {
    const store = useRulesStore()
    
    // Add rule to store first
    store.rules = [mockRule]

    const { rulesApi } = await import('@/api/rules')
    vi.mocked(rulesApi.deactivateRule).mockResolvedValueOnce(mockApiResponse({}))

    await store.deactivateRule('1')

    expect(store.rules[0].status).toBe('INACTIVE')
    expect(store.saving).toBe(false)
  })

  it('should validate rule successfully', async () => {
    const store = useRulesStore()
    const validationRequest = {
      dsl_content: 'IF quantity >= 3 THEN discount = price * 0.1',
      context: { quantity: 5, price: 100 },
      rule_category: 'PROMOTIONS'
    }
    const response = { data: { valid: true, result: { discount: 10 } } }

    const { rulesApi } = await import('@/api/rules')
    vi.mocked(rulesApi.validateRule).mockResolvedValueOnce(mockApiResponse(response))

    const result = await store.validateRule(validationRequest)

    expect(result).toEqual(response.data)
  })

  it('should duplicate rule successfully', async () => {
    const store = useRulesStore()
    const duplicatedRule = { ...mockRule, id: '2', name: 'Test Rule (Copy)' }
    const response = { data: duplicatedRule }

    // Add original rule to store first
    store.rules = [mockRule]

    const { rulesApi } = await import('@/api/rules')
    vi.mocked(rulesApi.getRule).mockResolvedValueOnce(mockApiResponse({ data: mockRule }))
    vi.mocked(rulesApi.createRule).mockResolvedValueOnce(mockApiResponse(response))

    const result = await store.duplicateRule('1')

    expect(result).toEqual(duplicatedRule)
    expect(store.rules).toHaveLength(2)
    expect(store.rules[0]).toEqual(duplicatedRule) // New rule should be at the beginning
  })

  it('should filter rules correctly', () => {
    const store = useRulesStore()
    const rules = [
      { ...mockRule, id: '1', status: 'ACTIVE', category: 'PROMOTIONS' },
      { ...mockRule, id: '2', status: 'DRAFT', category: 'COUPONS' },
      { ...mockRule, id: '3', status: 'ACTIVE', category: 'PROMOTIONS' }
    ]
    
    store.rules = rules

    // Test status filter
    store.updateFilters({ status: ['ACTIVE'] })
    expect(store.filteredRules).toHaveLength(2)
    expect(store.filteredRules.every(rule => rule.status === 'ACTIVE')).toBe(true)

    // Test category filter
    store.updateFilters({ status: [], category: 'PROMOTIONS' })
    expect(store.filteredRules).toHaveLength(2)
    expect(store.filteredRules.every(rule => rule.category === 'PROMOTIONS')).toBe(true)

    // Test search filter
    store.updateSearch('Test Rule 1')
    expect(store.filteredRules).toHaveLength(1)
    expect(store.filteredRules[0].id).toBe('1')
  })

  it('should calculate rules stats correctly', () => {
    const store = useRulesStore()
    const rules = [
      { ...mockRule, id: '1', status: 'ACTIVE' },
      { ...mockRule, id: '2', status: 'DRAFT' },
      { ...mockRule, id: '3', status: 'ACTIVE' },
      { ...mockRule, id: '4', status: 'UNDER_REVIEW' },
      { ...mockRule, id: '5', status: 'DEPRECATED' }
    ]
    
    store.rules = rules

    const stats = store.rulesStats
    expect(stats.total).toBe(5)
    expect(stats.active).toBe(2)
    expect(stats.draft).toBe(1)
    expect(stats.underReview).toBe(1)
    expect(stats.deprecated).toBe(1)
  })

  it('should perform bulk operations successfully', async () => {
    const store = useRulesStore()
    const rules = [
      { ...mockRule, id: '1', status: 'INACTIVE' },
      { ...mockRule, id: '2', status: 'INACTIVE' },
      { ...mockRule, id: '3', status: 'ACTIVE' }
    ]
    
    store.rules = rules

    const { rulesApi } = await import('@/api/rules')
    vi.mocked(rulesApi.activateRule).mockResolvedValue(mockApiResponse({}))

    await store.bulkActivate(['1', '2'])

    expect(store.rules[0].status).toBe('ACTIVE')
    expect(store.rules[1].status).toBe('ACTIVE')
    expect(store.rules[2].status).toBe('ACTIVE') // Unchanged
  })
})
