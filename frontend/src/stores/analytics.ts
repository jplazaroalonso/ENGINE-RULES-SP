import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import { apiClient, handleApiResponse, handleApiError } from '@/api/client'
import type { ApiResponse } from '@/types'

export interface AnalyticsMetrics {
  rulesMetrics: {
    totalRules: number
    activeRules: number
    executionsToday: number
    averageExecutionTime: number
    errorRate: number
  }
  systemHealth: {
    serviceStatus: Array<{
      service: string
      status: 'healthy' | 'degraded' | 'down'
      responseTime: number
      uptime: number
    }>
    totalUptime: number
  }
  performanceMetrics: {
    totalRequests: number
    averageResponseTime: number
    errorRate: number
    throughput: number
  }
}

export const useAnalyticsStore = defineStore('analytics', () => {
  // State
  const metrics = ref<AnalyticsMetrics | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const lastUpdated = ref<string | null>(null)

  // Getters
  const hasMetrics = computed(() => metrics.value !== null)
  const isHealthy = computed(() => {
    if (!metrics.value) return false
    return metrics.value.systemHealth.serviceStatus.every(
      service => service.status === 'healthy'
    )
  })

  // Actions
  const fetchMetrics = async () => {
    loading.value = true
    error.value = null

    try {
      // TODO: Replace with actual API call when backend implements analytics endpoints
      // For now, return mock data
      const mockMetrics: AnalyticsMetrics = {
        rulesMetrics: {
          totalRules: 15,
          activeRules: 12,
          executionsToday: 1247,
          averageExecutionTime: 45,
          errorRate: 0.02
        },
        systemHealth: {
          serviceStatus: [
            {
              service: 'rules-management',
              status: 'healthy',
              responseTime: 23,
              uptime: 99.9
            },
            {
              service: 'rules-evaluation',
              status: 'healthy',
              responseTime: 18,
              uptime: 99.8
            },
            {
              service: 'rules-calculator',
              status: 'healthy',
              responseTime: 31,
              uptime: 99.7
            }
          ],
          totalUptime: 99.8
        },
        performanceMetrics: {
          totalRequests: 15678,
          averageResponseTime: 24,
          errorRate: 0.015,
          throughput: 234
        }
      }

      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 500))

      metrics.value = mockMetrics
      lastUpdated.value = new Date().toISOString()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch analytics'
      throw err
    } finally {
      loading.value = false
    }
  }

  const refreshMetrics = async () => {
    await fetchMetrics()
  }

  const clearError = () => {
    error.value = null
  }

  return {
    // State
    metrics: readonly(metrics),
    loading: readonly(loading),
    error: readonly(error),
    lastUpdated: readonly(lastUpdated),
    
    // Getters
    hasMetrics,
    isHealthy,
    
    // Actions
    fetchMetrics,
    refreshMetrics,
    clearError
  }
})
