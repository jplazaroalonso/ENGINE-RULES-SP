import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import { authApi } from '@/api/auth'
import type { User, UserRole, LoginRequest, UpdateUserRequest, ChangePasswordRequest } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('auth_token'))
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Getters
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const userRole = computed(() => user.value?.role || 'VIEWER')
  const isAdmin = computed(() => userRole.value === 'ADMIN')
  const isManager = computed(() => ['ADMIN', 'MANAGER'].includes(userRole.value))
  const canEdit = computed(() => ['ADMIN', 'MANAGER', 'USER'].includes(userRole.value))
  const canView = computed(() => true) // All authenticated users can view

  // Actions
  const login = async (email: string, password: string) => {
    loading.value = true
    error.value = null

    try {
      const response = await authApi.login({ email, password })
      
      token.value = response.token
      user.value = response.user
      
      // Store token in localStorage
      localStorage.setItem('auth_token', response.token)
      
      return response
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Login failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  const logout = async () => {
    loading.value = true

    try {
      // Call logout API
      await authApi.logout()
    } catch (err) {
      console.error('Logout API error:', err)
    } finally {
      // Always clear local state
      token.value = null
      user.value = null
      localStorage.removeItem('auth_token')
      loading.value = false
    }
  }

  const initialize = async () => {
    if (token.value) {
      try {
        const userData = await authApi.getCurrentUser()
        user.value = userData
      } catch (err) {
        console.error('Failed to initialize auth:', err)
        // Clear invalid token
        await logout()
      }
    }
  }

  const updateProfile = async (updates: UpdateUserRequest) => {
    if (!user.value) throw new Error('No user logged in')

    loading.value = true
    error.value = null

    try {
      const updatedUser = await authApi.updateProfile(user.value.id, updates)
      user.value = updatedUser
      return updatedUser
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Update failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  const changePassword = async (currentPassword: string, newPassword: string) => {
    if (!user.value) throw new Error('No user logged in')

    loading.value = true
    error.value = null

    try {
      await authApi.changePassword(user.value.id, { currentPassword, newPassword })
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Password change failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  const clearError = () => {
    error.value = null
  }


  return {
    // State
    user: readonly(user),
    token: readonly(token),
    loading: readonly(loading),
    error: readonly(error),
    
    // Getters
    isAuthenticated,
    userRole,
    isAdmin,
    isManager,
    canEdit,
    canView,
    
    // Actions
    login,
    logout,
    initialize,
    updateProfile,
    changePassword,
    clearError
  }
})
