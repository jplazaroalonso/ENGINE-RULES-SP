import { mount, VueWrapper } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import { vi } from 'vitest'
import type { App } from 'vue'

// Mock router
const createMockRouter = () => {
  return createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/', component: { template: '<div>Home</div>' } },
      { path: '/login', component: { template: '<div>Login</div>' } },
      { path: '/dashboard', component: { template: '<div>Dashboard</div>' } },
      { path: '/rules', component: { template: '<div>Rules</div>' } }
    ]
  })
}

// Mock stores
export const createMockStores = () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  return pinia
}

// Test wrapper factory
export const createTestWrapper = (
  component: any,
  options: any = {}
): VueWrapper<any> => {
  const pinia = createMockStores()
  const router = createMockRouter()
  
  return mount(component, {
    global: {
      plugins: [pinia, router],
      stubs: {
        'router-link': true,
        'router-view': true
      },
      mocks: {
        $router: router,
        $route: router.currentRoute.value
      }
    },
    ...options
  })
}

// Mock API responses
export const mockApiResponse = <T>(data: T, status = 200) => ({
  data,
  status,
  statusText: 'OK',
  headers: {},
  config: {}
})

// Mock error response
export const mockApiError = (message = 'API Error', status = 500) => {
  const error = new Error(message)
  ;(error as any).response = {
    data: { message },
    status,
    statusText: 'Internal Server Error',
    headers: {},
    config: {}
  }
  return error
}

// Mock rule data
export const mockRule = {
  id: '1',
  name: 'Test Rule',
  description: 'A test rule for unit testing',
  dsl_content: 'IF quantity >= 3 THEN discount = price * 0.1',
  category: 'PROMOTIONS',
  priority: 'MEDIUM',
  status: 'ACTIVE',
  version: '1.0.0',
  tags: ['test', 'promotion'],
  created_by: 'test-user',
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
  effective_date: null,
  expiration_date: null,
  requires_approval: false,
  evaluation_count: 0,
  success_rate: 0,
  last_evaluated: null
}

// Mock user data
export const mockUser = {
  id: '1',
  username: 'testuser',
  email: 'test@example.com',
  name: 'Test User',
  role: 'admin',
  permissions: ['read', 'write', 'admin']
}

// Mock notification data
export const mockNotification = {
  id: '1',
  type: 'success',
  title: 'Success',
  message: 'Operation completed successfully',
  timestamp: new Date().toISOString(),
  read: false
}

// Wait for next tick
export const waitForNextTick = () => new Promise(resolve => setTimeout(resolve, 0))

// Mock localStorage
export const mockLocalStorage = () => {
  const store: Record<string, string> = {}
  
  return {
    getItem: vi.fn((key: string) => store[key] || null),
    setItem: vi.fn((key: string, value: string) => {
      store[key] = value
    }),
    removeItem: vi.fn((key: string) => {
      delete store[key]
    }),
    clear: vi.fn(() => {
      Object.keys(store).forEach(key => delete store[key])
    })
  }
}

// Mock sessionStorage
export const mockSessionStorage = () => {
  const store: Record<string, string> = {}
  
  return {
    getItem: vi.fn((key: string) => store[key] || null),
    setItem: vi.fn((key: string, value: string) => {
      store[key] = value
    }),
    removeItem: vi.fn((key: string) => {
      delete store[key]
    }),
    clear: vi.fn(() => {
      Object.keys(store).forEach(key => delete store[key])
    })
  }
}

// Mock fetch
export const mockFetch = (response: any, ok = true) => {
  global.fetch = vi.fn(() =>
    Promise.resolve({
      ok,
      status: ok ? 200 : 400,
      json: () => Promise.resolve(response),
      text: () => Promise.resolve(JSON.stringify(response))
    })
  ) as any
}

// Mock console methods
export const mockConsole = () => {
  const originalConsole = { ...console }
  
  beforeEach(() => {
    console.log = vi.fn()
    console.warn = vi.fn()
    console.error = vi.fn()
  })
  
  afterEach(() => {
    Object.assign(console, originalConsole)
  })
}

// Test data generators
export const generateRules = (count: number) => {
  return Array.from({ length: count }, (_, index) => ({
    ...mockRule,
    id: `${index + 1}`,
    name: `Test Rule ${index + 1}`,
    description: `Test rule ${index + 1} description`
  }))
}

export const generateUsers = (count: number) => {
  return Array.from({ length: count }, (_, index) => ({
    ...mockUser,
    id: `${index + 1}`,
    username: `user${index + 1}`,
    email: `user${index + 1}@example.com`,
    name: `User ${index + 1}`
  }))
}
