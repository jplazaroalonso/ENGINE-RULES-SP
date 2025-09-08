import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '../auth'
import { mockApiResponse, mockApiError, mockUser } from '../../test/utils/test-utils'

// Mock the API client
vi.mock('@/api/client', () => ({
  apiClient: {
    post: vi.fn()
  }
}))

describe('Auth Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should initialize with default state', () => {
    const store = useAuthStore()
    
    expect(store.user).toBeNull()
    expect(store.isAuthenticated).toBe(false)
    expect(store.loading).toBe(false)
    expect(store.error).toBeNull()
  })

  it('should login successfully', async () => {
    const store = useAuthStore()
    const loginResponse = {
      user: mockUser,
      token: 'mock-jwt-token'
    }

    // Mock successful API response
    const { apiClient } = await import('@/api/client')
    vi.mocked(apiClient.post).mockResolvedValueOnce(mockApiResponse(loginResponse))

    await store.login({ username: 'testuser', password: 'password' })

    expect(store.user).toEqual(mockUser)
    expect(store.isAuthenticated).toBe(true)
    expect(store.error).toBeNull()
    expect(apiClient.post).toHaveBeenCalledWith('/auth/login', {
      username: 'testuser',
      password: 'password'
    })
  })

  it('should handle login failure', async () => {
    const store = useAuthStore()
    const error = mockApiError('Invalid credentials', 401)

    // Mock failed API response
    const { apiClient } = await import('@/api/client')
    vi.mocked(apiClient.post).mockRejectedValueOnce(error)

    await store.login({ username: 'testuser', password: 'wrongpassword' })

    expect(store.user).toBeNull()
    expect(store.isAuthenticated).toBe(false)
    expect(store.error).toBe('Invalid credentials')
  })

  it('should logout successfully', () => {
    const store = useAuthStore()
    
    // Set initial state
    store.user = mockUser
    store.isAuthenticated = true

    store.logout()

    expect(store.user).toBeNull()
    expect(store.isAuthenticated).toBe(false)
    expect(store.error).toBeNull()
  })

  it('should clear error', () => {
    const store = useAuthStore()
    store.error = 'Some error'

    store.clearError()

    expect(store.error).toBeNull()
  })

  it('should set loading state during login', async () => {
    const store = useAuthStore()
    const loginResponse = {
      user: mockUser,
      token: 'mock-jwt-token'
    }

    // Mock API response with delay
    const { apiClient } = await import('@/api/client')
    vi.mocked(apiClient.post).mockImplementationOnce(
      () => new Promise(resolve => 
        setTimeout(() => resolve(mockApiResponse(loginResponse)), 100)
      )
    )

    const loginPromise = store.login({ username: 'testuser', password: 'password' })
    
    // Check loading state during request
    expect(store.loading).toBe(true)
    
    await loginPromise
    
    // Check loading state after request
    expect(store.loading).toBe(false)
  })
})
