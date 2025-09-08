import { config } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { vi } from 'vitest'

// Mock Quasar components
vi.mock('quasar', () => ({
  Quasar: {
    install: vi.fn()
  },
  QBtn: {
    name: 'QBtn',
    template: '<button><slot /></button>'
  },
  QCard: {
    name: 'QCard',
    template: '<div class="q-card"><slot /></div>'
  },
  QCardSection: {
    name: 'QCardSection',
    template: '<div class="q-card__section"><slot /></div>'
  },
  QInput: {
    name: 'QInput',
    template: '<input class="q-input" />'
  },
  QSelect: {
    name: 'QSelect',
    template: '<select class="q-select" />'
  },
  QTable: {
    name: 'QTable',
    template: '<table class="q-table" />'
  },
  QChip: {
    name: 'QChip',
    template: '<span class="q-chip"><slot /></span>'
  },
  QIcon: {
    name: 'QIcon',
    template: '<i class="q-icon" />'
  },
  QSpinner: {
    name: 'QSpinner',
    template: '<div class="q-spinner" />'
  },
  QDialog: {
    name: 'QDialog',
    template: '<div class="q-dialog" v-if="modelValue"><slot /></div>'
  },
  QMenu: {
    name: 'QMenu',
    template: '<div class="q-menu" v-if="modelValue"><slot /></div>'
  },
  QList: {
    name: 'QList',
    template: '<ul class="q-list"><slot /></ul>'
  },
  QItem: {
    name: 'QItem',
    template: '<li class="q-item"><slot /></li>'
  },
  QItemSection: {
    name: 'QItemSection',
    template: '<div class="q-item__section"><slot /></div>'
  },
  QSeparator: {
    name: 'QSeparator',
    template: '<hr class="q-separator" />'
  },
  QTooltip: {
    name: 'QTooltip',
    template: '<div class="q-tooltip"><slot /></div>'
  },
  QCheckbox: {
    name: 'QCheckbox',
    template: '<input type="checkbox" class="q-checkbox" />'
  },
  QBtnGroup: {
    name: 'QBtnGroup',
    template: '<div class="q-btn-group"><slot /></div>'
  },
  QLayout: {
    name: 'QLayout',
    template: '<div class="q-layout"><slot /></div>'
  },
  QHeader: {
    name: 'QHeader',
    template: '<header class="q-header"><slot /></header>'
  },
  QDrawer: {
    name: 'QDrawer',
    template: '<aside class="q-drawer"><slot /></aside>'
  },
  QPageContainer: {
    name: 'QPageContainer',
    template: '<main class="q-page-container"><slot /></main>'
  },
  QFooter: {
    name: 'QFooter',
    template: '<footer class="q-footer"><slot /></footer>'
  },
  QPage: {
    name: 'QPage',
    template: '<div class="q-page"><slot /></div>'
  },
  Notify: {
    create: vi.fn()
  },
  Dialog: {
    create: vi.fn()
  }
}))

// Mock Vue Router
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn(),
    back: vi.fn(),
    replace: vi.fn()
  }),
  useRoute: () => ({
    params: {},
    query: {},
    path: '/',
    name: 'Test'
  })
}))

// Mock Axios
vi.mock('axios', () => ({
  default: {
    create: vi.fn(() => ({
      get: vi.fn(),
      post: vi.fn(),
      put: vi.fn(),
      delete: vi.fn(),
      interceptors: {
        request: { use: vi.fn() },
        response: { use: vi.fn() }
      }
    }))
  }
}))

// Global test configuration
config.global.plugins = [createPinia()]

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn()
  }))
})

// Mock ResizeObserver
global.ResizeObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn()
}))

// Mock IntersectionObserver
global.IntersectionObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn()
}))
