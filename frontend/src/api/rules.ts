import { apiClient, evaluationApiClient, calculatorApiClient, handleApiResponse, handleApiError } from './client'
import type { 
  Rule, 
  CreateRuleRequest, 
  UpdateRuleRequest, 
  RuleListParams,
  ApiResponse,
  EvaluationRequest,
  EvaluationResponse,
  CalculationRequest,
  CalculationResponse
} from '@/types'

export const rulesApi = {
  // Get all rules with pagination and filtering
  async getRules(params?: RuleListParams): Promise<ApiResponse<Rule[]>> {
    try {
      const queryParams = new URLSearchParams()
      
      if (params) {
        if (params.page) queryParams.append('page', params.page.toString())
        if (params.limit) queryParams.append('limit', params.limit.toString())
        if (params.sortBy) queryParams.append('sort_by', params.sortBy)
        if (params.sortOrder) queryParams.append('sort_order', params.sortOrder)
        if (params.status) queryParams.append('status', params.status)
        if (params.category) queryParams.append('category', params.category)
        if (params.search) queryParams.append('search', params.search)
      }
      
      const url = `/api/v1/rules${queryParams.toString() ? '?' + queryParams.toString() : ''}`
      const response = await apiClient.get(url)
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Get a single rule by ID
  async getRule(id: string): Promise<ApiResponse<Rule>> {
    try {
      const response = await apiClient.get(`/api/v1/rules/${id}`)
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Create a new rule
  async createRule(rule: CreateRuleRequest): Promise<ApiResponse<Rule>> {
    try {
      const response = await apiClient.post('/api/v1/rules', rule)
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Update an existing rule
  async updateRule(id: string, rule: UpdateRuleRequest): Promise<ApiResponse<Rule>> {
    try {
      const response = await apiClient.put(`/api/v1/api/v1/rules/${id}`, rule)
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Delete a rule
  async deleteRule(id: string): Promise<void> {
    try {
      await apiClient.delete(`/api/v1/api/v1/rules/${id}`)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Submit rule for approval
  async submitForApproval(id: string): Promise<ApiResponse<Rule>> {
    try {
      const response = await apiClient.post(`/api/v1/api/v1/rules/${id}/submit-approval`)
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Approve a rule
  async approveRule(id: string): Promise<ApiResponse<Rule>> {
    try {
      const response = await apiClient.post(`/api/v1/api/v1/rules/${id}/approve`)
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Activate a rule
  async activateRule(id: string): Promise<ApiResponse<Rule>> {
    try {
      const response = await apiClient.post(`/api/v1/rules/${id}/activate`)
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Deactivate a rule
  async deactivateRule(id: string): Promise<ApiResponse<Rule>> {
    try {
      const response = await apiClient.post(`/api/v1/rules/${id}/deactivate`)
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Get rule metrics
  async getRuleMetrics(id: string): Promise<ApiResponse<Record<string, any>>> {
    try {
      const response = await apiClient.get(`/api/v1/rules/${id}/metrics`)
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Validate rule DSL
  async validateRule(request: { dslContent: string; testData?: Record<string, unknown> }): Promise<ApiResponse<unknown>> {
    try {
      const response = await apiClient.post('/api/v1/rules/validate', request)
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Test rule with sample data
  async testRule(id: string, testData: Record<string, unknown>): Promise<ApiResponse<unknown>> {
    try {
      const response = await apiClient.post(`/api/v1/rules/${id}/test`, { testData })
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Duplicate a rule
  async duplicateRule(id: string, newName?: string): Promise<ApiResponse<Rule>> {
    try {
      const response = await apiClient.post(`/api/v1/rules/${id}/duplicate`, { newName })
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Get rule history
  async getRuleHistory(id: string): Promise<ApiResponse<unknown[]>> {
    try {
      const response = await apiClient.get(`/api/v1/rules/${id}/history`)
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Bulk operations
  async bulkActivate(ruleIds: string[]): Promise<ApiResponse<Rule[]>> {
    try {
      const response = await apiClient.post('/api/v1/rules/bulk/activate', { ruleIds })
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  async bulkDeactivate(ruleIds: string[]): Promise<ApiResponse<Rule[]>> {
    try {
      const response = await apiClient.post('/api/v1/rules/bulk/deactivate', { ruleIds })
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  async bulkDelete(ruleIds: string[]): Promise<void> {
    try {
      await apiClient.post('/api/v1/rules/bulk/delete', { ruleIds })
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Export rules
  async exportRules(format: 'json' | 'csv' | 'xlsx', filters?: Record<string, unknown>): Promise<Blob> {
    try {
      const response = await apiClient.post('/api/v1/rules/export', { format, filters }, {
        responseType: 'blob'
      })
      return response.data
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Import rules
  async importRules(file: File): Promise<ApiResponse<{ imported: number; errors: unknown[] }>> {
    try {
      const formData = new FormData()
      formData.append('file', file)
      
      const response = await apiClient.post('/api/v1/rules/import', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  }
}

// Evaluation API
export const evaluationApi = {
  // Evaluate a rule
  async evaluateRule(request: EvaluationRequest): Promise<ApiResponse<EvaluationResponse>> {
    try {
      const response = await evaluationApiClient.post('/api/v1/evaluate', request)
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  }
}

// Calculator API
export const calculatorApi = {
  // Calculate rule results
  async calculateRules(request: CalculationRequest): Promise<ApiResponse<CalculationResponse>> {
    try {
      const response = await calculatorApiClient.post('/api/v1/calculate', request)
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  }
}
