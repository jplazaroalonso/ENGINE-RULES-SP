<template>
  <!-- Login page without layout -->
  <div v-if="isLoginPage" class="login-wrapper">
    <router-view />
  </div>
  
  <!-- Main app with layout -->
  <q-layout v-else view="lHh Lpr lFf" class="app-layout">
    <!-- Header -->
    <q-header elevated class="app-header">
      <AppHeader 
        :user="currentUser"
        :notifications="notifications"
        @logout="handleLogout"
        @toggle-drawer="toggleLeftDrawer"
      />
    </q-header>

    <!-- Navigation Drawer -->
    <q-drawer
      v-model="leftDrawerOpen"
      show-if-above
      bordered
      class="app-drawer"
    >
      <MainNavigation 
        :user="currentUser"
        :current-route="$route.name"
        @navigate="handleNavigation"
      />
    </q-drawer>

    <!-- Main Content -->
    <q-page-container class="app-content">
      <router-view v-slot="{ Component, route }">
        <transition name="page" mode="out-in">
          <component :is="Component" :key="route.path" />
        </transition>
      </router-view>
    </q-page-container>

    <!-- Footer -->
    <q-footer elevated class="app-footer">
      <AppFooter />
    </q-footer>

    <!-- Global Loading -->
    <q-inner-loading :showing="globalLoading">
      <q-spinner-dots size="50px" color="primary" />
    </q-inner-loading>
  </q-layout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notifications'
import { useUIStore } from '@/stores/ui'

// Components
import AppHeader from '@/components/layout/AppHeader.vue'
import MainNavigation from '@/components/layout/MainNavigation.vue'
import AppFooter from '@/components/layout/AppFooter.vue'

const router = useRouter()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()
const uiStore = useUIStore()

// Reactive state
const leftDrawerOpen = ref(false)

// Computed properties
const currentUser = computed(() => authStore.user)
const notifications = computed(() => notificationStore.unreadNotifications)
const globalLoading = computed(() => uiStore.globalLoading)
const isLoginPage = computed(() => router.currentRoute.value.name === 'Login')

// Methods
const toggleLeftDrawer = () => {
  leftDrawerOpen.value = !leftDrawerOpen.value
}

const handleNavigation = (route: string) => {
  router.push(route)
  leftDrawerOpen.value = false
}

const handleLogout = async () => {
  try {
    await authStore.logout()
    await router.push('/login')
  } catch (error) {
    console.error('Logout failed:', error)
  }
}

// Lifecycle
onMounted(async () => {
  // Initialize stores
  await authStore.initialize()
  await notificationStore.fetchNotifications()
})
</script>

<style lang="scss" scoped>
.login-wrapper {
  min-height: 100vh;
  width: 100%;
}

.app-layout {
  background-color: var(--carrefour-gray-50);
}

.app-header {
  background: linear-gradient(135deg, var(--carrefour-blue) 0%, var(--carrefour-dark-blue) 100%);
  box-shadow: 0 2px 8px rgba(0, 75, 135, 0.15);
}

.app-drawer {
  background-color: white;
  border-right: 1px solid var(--carrefour-gray-200);
}

.app-content {
  background-color: var(--carrefour-gray-50);
  min-height: calc(100vh - 64px - 48px); // Subtract header and footer height
}

.app-footer {
  background-color: var(--carrefour-gray-800);
  color: white;
}

// Page transitions
.page-enter-active,
.page-leave-active {
  transition: all 0.3s ease;
}

.page-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.page-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}
</style>
