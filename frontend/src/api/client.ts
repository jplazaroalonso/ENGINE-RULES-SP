import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notifications'

// Create axios instances for each service
export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'https://rules-management.local.dev',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const evaluationApiClient = axios.create({
  baseURL: import.meta.env.VITE_EVALUATION_API_URL || 'https://rules-evaluation.local.dev',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const calculatorApiClient = axios.create({
  baseURL: import.meta.env.VITE_CALCULATOR_API_URL || 'https://rules-calculator.local.dev',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor for all clients
const setupRequestInterceptor = (client: typeof apiClient) => {
  client.interceptors.request.use(
    (config) => {
      // Add auth token if available
      const authStore = useAuthStore()
      if (authStore.token) {
        config.headers.Authorization = `Bearer ${authStore.token}`
      }
      
      // Add request timestamp for debugging
      config.metadata = { startTime: new Date() }
      
      return config
    },
    (error) => {
      return Promise.reject(error)
    }
  )
}

// Setup request interceptors for all clients
setupRequestInterceptor(apiClient)
setupRequestInterceptor(evaluationApiClient)
setupRequestInterceptor(calculatorApiClient)

// Response interceptor for all clients
const setupResponseInterceptor = (client: typeof apiClient) => {
  client.interceptors.response.use(
    (response) => {
      // Log response time for debugging
      if (response.config.metadata?.startTime) {
        const endTime = new Date()
        const duration = endTime.getTime() - response.config.metadata.startTime.getTime()
        console.log(`API Request to ${response.config.url} took ${duration}ms`)
      }
      
      return response
    },
    async (error) => {
      const notificationStore = useNotificationStore()
      const authStore = useAuthStore()
      
      // Handle different error types
      if (error.response) {
        const { status, data } = error.response
        
        switch (status) {
          case 401:
            // Unauthorized - redirect to login
            await authStore.logout()
            window.location.href = '/login'
            break
            
          case 403:
            // Forbidden
            notificationStore.showError('Access Denied', 'You do not have permission to perform this action')
            break
            
          case 404:
            // Not found
            notificationStore.showError('Not Found', 'The requested resource was not found')
            break
            
          case 422:
            // Validation error
            const message = data?.message || 'Validation failed'
            notificationStore.showError('Validation Error', message)
            break
            
          case 429:
            // Rate limited
            notificationStore.showWarning('Rate Limited', 'Too many requests. Please try again later.')
            break
            
          case 500:
            // Server error
            notificationStore.showError('Server Error', 'An unexpected error occurred. Please try again.')
            break
            
          default:
            // Other errors
            const errorMessage = data?.message || `Request failed with status ${status}`
            notificationStore.showError('Request Failed', errorMessage)
        }
      } else if (error.request) {
        // Network error
        notificationStore.showError('Network Error', 'Unable to connect to the server. Please check your connection.')
      } else {
        // Other error
        notificationStore.showError('Error', error.message || 'An unexpected error occurred')
      }
      
      return Promise.reject(error)
    }
  )
}

// Setup response interceptors for all clients
setupResponseInterceptor(apiClient)
setupResponseInterceptor(evaluationApiClient)
setupResponseInterceptor(calculatorApiClient)

// Helper function to handle API responses
export const handleApiResponse = <T>(response: any): T => {
  if (response.data) {
    return response.data
  }
  throw new Error('Invalid API response format')
}

// Helper function to handle API errors
export const handleApiError = (error: any): never => {
  if (error.response?.data?.message) {
    throw new Error(error.response.data.message)
  }
  if (error.message) {
    throw new Error(error.message)
  }
  throw new Error('An unexpected error occurred')
}
