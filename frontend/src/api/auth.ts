import { apiClient, handleApiResponse, handleApiError } from './client'
import type { User, LoginRequest, LoginResponse, UpdateUserRequest, ChangePasswordRequest } from '@/types'

export const authApi = {
  // Login user
  async login(request: LoginRequest): Promise<LoginResponse> {
    try {
      const response = await apiClient.post('/api/v1/api/v1/auth/login', request)
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Logout user
  async logout(): Promise<void> {
    try {
      await apiClient.post('/api/v1/auth/logout')
    } catch (error) {
      // Don't throw error on logout failure
      console.warn('Logout request failed:', error)
    }
  },

  // Get current user
  async getCurrentUser(): Promise<User> {
    try {
      const response = await apiClient.get('/api/v1/auth/me')
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Update user profile
  async updateProfile(userId: string, updates: UpdateUserRequest): Promise<User> {
    try {
      const response = await apiClient.put(`/api/v1/auth/users/${userId}`, updates)
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Change password
  async changePassword(userId: string, request: ChangePasswordRequest): Promise<void> {
    try {
      await apiClient.post(`/api/v1/auth/users/${userId}/change-password`, request)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Refresh token
  async refreshToken(): Promise<LoginResponse> {
    try {
      const response = await apiClient.post('/api/v1/auth/refresh')
      return handleApiResponse(response)
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Verify email
  async verifyEmail(token: string): Promise<void> {
    try {
      await apiClient.post('/api/v1/auth/verify-email', { token })
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Request password reset
  async requestPasswordReset(email: string): Promise<void> {
    try {
      await apiClient.post('/api/v1/auth/forgot-password', { email })
    } catch (error) {
      throw handleApiError(error)
    }
  },

  // Reset password
  async resetPassword(token: string, newPassword: string): Promise<void> {
    try {
      await apiClient.post('/api/v1/auth/reset-password', { token, newPassword })
    } catch (error) {
      throw handleApiError(error)
    }
  }
}
