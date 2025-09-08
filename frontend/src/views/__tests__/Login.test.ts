import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createTestWrapper, waitForNextTick } from '../../test/utils/test-utils'
import Login from '../auth/Login.vue'
import { useAuthStore } from '../../stores/auth'

// Mock the auth store
vi.mock('../../stores/auth', () => ({
  useAuthStore: vi.fn()
}))

describe('Login', () => {
  let mockAuthStore: any

  beforeEach(() => {
    mockAuthStore = {
      login: vi.fn(),
      loading: false,
      error: null,
      clearError: vi.fn()
    }
    
    vi.mocked(useAuthStore).mockReturnValue(mockAuthStore)
  })

  it('should render login form', () => {
    const wrapper = createTestWrapper(Login)

    expect(wrapper.find('input[type="text"]').exists()).toBe(true)
    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('should have correct form labels', () => {
    const wrapper = createTestWrapper(Login)

    expect(wrapper.text()).toContain('Username')
    expect(wrapper.text()).toContain('Password')
  })

  it('should show login button', () => {
    const wrapper = createTestWrapper(Login)

    const loginButton = wrapper.find('button[type="submit"]')
    expect(loginButton.exists()).toBe(true)
    expect(loginButton.text()).toContain('Login')
  })

  it('should handle form submission', async () => {
    const wrapper = createTestWrapper(Login)

    // Fill form
    await wrapper.find('input[type="text"]').setValue('testuser')
    await wrapper.find('input[type="password"]').setValue('password123')

    // Submit form
    await wrapper.find('form').trigger('submit')

    expect(mockAuthStore.login).toHaveBeenCalledWith({
      username: 'testuser',
      password: 'password123'
    })
  })

  it('should show loading state', () => {
    mockAuthStore.loading = true
    
    const wrapper = createTestWrapper(Login)

    const loginButton = wrapper.find('button[type="submit"]')
    expect(loginButton.attributes('disabled')).toBeDefined()
  })

  it('should show error message', () => {
    mockAuthStore.error = 'Invalid credentials'
    
    const wrapper = createTestWrapper(Login)

    expect(wrapper.text()).toContain('Invalid credentials')
  })

  it('should clear error when form is modified', async () => {
    mockAuthStore.error = 'Some error'
    
    const wrapper = createTestWrapper(Login)

    // Modify form
    await wrapper.find('input[type="text"]').setValue('newuser')

    expect(mockAuthStore.clearError).toHaveBeenCalled()
  })

  it('should validate required fields', async () => {
    const wrapper = createTestWrapper(Login)

    // Submit empty form
    await wrapper.find('form').trigger('submit')

    expect(mockAuthStore.login).not.toHaveBeenCalled()
  })

  it('should handle login success', async () => {
    mockAuthStore.login.mockResolvedValueOnce({})
    
    const wrapper = createTestWrapper(Login)

    // Fill and submit form
    await wrapper.find('input[type="text"]').setValue('testuser')
    await wrapper.find('input[type="password"]').setValue('password123')
    await wrapper.find('form').trigger('submit')

    expect(mockAuthStore.login).toHaveBeenCalledWith({
      username: 'testuser',
      password: 'password123'
    })
  })

  it('should handle login failure', async () => {
    const error = new Error('Login failed')
    mockAuthStore.login.mockRejectedValueOnce(error)
    
    const wrapper = createTestWrapper(Login)

    // Fill and submit form
    await wrapper.find('input[type="text"]').setValue('testuser')
    await wrapper.find('input[type="password"]').setValue('wrongpassword')
    await wrapper.find('form').trigger('submit')

    await waitForNextTick()

    expect(mockAuthStore.login).toHaveBeenCalledWith({
      username: 'testuser',
      password: 'wrongpassword'
    })
  })

  it('should have proper form structure', () => {
    const wrapper = createTestWrapper(Login)

    const form = wrapper.find('form')
    expect(form.exists()).toBe(true)

    const inputs = wrapper.findAll('input')
    expect(inputs).toHaveLength(2)

    const submitButton = wrapper.find('button[type="submit"]')
    expect(submitButton.exists()).toBe(true)
  })

  it('should have correct input types', () => {
    const wrapper = createTestWrapper(Login)

    const usernameInput = wrapper.find('input[type="text"]')
    const passwordInput = wrapper.find('input[type="password"]')

    expect(usernameInput.exists()).toBe(true)
    expect(passwordInput.exists()).toBe(true)
  })
})
