import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'

export const useUIStore = defineStore('ui', () => {
  // State
  const globalLoading = ref(false)
  const sidebarOpen = ref(true)
  const theme = ref<'light' | 'dark'>('light')
  const language = ref('en')
  const pageTitle = ref('Rules Engine')
  const breadcrumbs = ref<Array<{ label: string; to?: string }>>([])

  // Getters
  const isDarkMode = computed(() => theme.value === 'dark')
  const sidebarWidth = computed(() => sidebarOpen.value ? 280 : 0)

  // Actions
  const setGlobalLoading = (loading: boolean) => {
    globalLoading.value = loading
  }

  const toggleSidebar = () => {
    sidebarOpen.value = !sidebarOpen.value
  }

  const setSidebarOpen = (open: boolean) => {
    sidebarOpen.value = open
  }

  const toggleTheme = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
    // Store theme preference
    localStorage.setItem('theme', theme.value)
  }

  const setTheme = (newTheme: 'light' | 'dark') => {
    theme.value = newTheme
    localStorage.setItem('theme', newTheme)
  }

  const setLanguage = (lang: string) => {
    language.value = lang
    localStorage.setItem('language', lang)
  }

  const setPageTitle = (title: string) => {
    pageTitle.value = title
    document.title = `${title} - Rules Engine`
  }

  const setBreadcrumbs = (crumbs: Array<{ label: string; to?: string }>) => {
    breadcrumbs.value = crumbs
  }

  const addBreadcrumb = (crumb: { label: string; to?: string }) => {
    breadcrumbs.value.push(crumb)
  }

  const clearBreadcrumbs = () => {
    breadcrumbs.value = []
  }

  // Initialize from localStorage
  const initialize = () => {
    const savedTheme = localStorage.getItem('theme') as 'light' | 'dark'
    const savedLanguage = localStorage.getItem('language')
    const savedSidebarOpen = localStorage.getItem('sidebarOpen')

    if (savedTheme) {
      theme.value = savedTheme
    }
    if (savedLanguage) {
      language.value = savedLanguage
    }
    if (savedSidebarOpen !== null) {
      sidebarOpen.value = savedSidebarOpen === 'true'
    }
  }

  return {
    // State
    globalLoading: readonly(globalLoading),
    sidebarOpen: readonly(sidebarOpen),
    theme: readonly(theme),
    language: readonly(language),
    pageTitle: readonly(pageTitle),
    breadcrumbs: readonly(breadcrumbs),
    
    // Getters
    isDarkMode,
    sidebarWidth,
    
    // Actions
    setGlobalLoading,
    toggleSidebar,
    setSidebarOpen,
    toggleTheme,
    setTheme,
    setLanguage,
    setPageTitle,
    setBreadcrumbs,
    addBreadcrumb,
    clearBreadcrumbs,
    initialize
  }
})
