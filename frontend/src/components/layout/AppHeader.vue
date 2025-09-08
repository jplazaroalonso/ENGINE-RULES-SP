<template>
  <q-toolbar class="app-header-toolbar">
    <!-- Mobile menu button -->
    <q-btn
      flat
      dense
      round
      icon="menu"
      class="q-mr-sm"
      @click="$emit('toggle-drawer')"
    />

    <!-- Logo and title -->
    <div class="app-header-brand">
      <img 
        src="/favicon.svg" 
        alt="Carrefour" 
        class="app-header-logo"
      />
      <div class="app-header-title">
        <span class="app-header-title-main">Rules Engine</span>
        <span class="app-header-title-sub">Carrefour</span>
      </div>
    </div>

    <q-space />

    <!-- Search bar -->
    <div class="app-header-search">
      <q-input
        v-model="searchQuery"
        placeholder="Search rules, campaigns..."
        dense
        outlined
        class="search-input"
        @keyup.enter="handleSearch"
      >
        <template v-slot:prepend>
          <q-icon name="search" />
        </template>
        <template v-slot:append>
          <q-icon 
            v-if="searchQuery" 
            name="clear" 
            class="cursor-pointer" 
            @click="searchQuery = ''"
          />
        </template>
      </q-input>
    </div>

    <q-space />

    <!-- Header actions -->
    <div class="app-header-actions">
      <!-- Notifications -->
      <q-btn
        flat
        round
        dense
        icon="notifications"
        class="action-button"
        :class="{ 'has-notifications': unreadCount > 0 }"
      >
        <q-badge 
          v-if="unreadCount > 0" 
          color="negative" 
          floating 
          rounded
        >
          {{ unreadCount > 99 ? '99+' : unreadCount }}
        </q-badge>
        
        <q-menu
          fit
          anchor="bottom right"
          self="top right"
          class="notifications-menu"
        >
          <q-list style="min-width: 300px">
            <q-item-label header>
              <div class="row items-center justify-between">
                <span>Notifications</span>
                <q-btn
                  flat
                  dense
                  size="sm"
                  label="Mark all read"
                  @click="markAllAsRead"
                />
              </div>
            </q-item-label>
            
            <q-separator />
            
            <template v-if="notifications.length > 0">
              <q-item
                v-for="notification in notifications.slice(0, 5)"
                :key="notification.id"
                clickable
                v-close-popup
                :class="{ 'unread': !notification.read }"
                @click="handleNotificationClick(notification)"
              >
                <q-item-section avatar>
                  <q-icon
                    :name="getNotificationIcon(notification.type)"
                    :color="getNotificationColor(notification.type)"
                  />
                </q-item-section>
                
                <q-item-section>
                  <q-item-label>{{ notification.title }}</q-item-label>
                  <q-item-label caption lines="2">
                    {{ notification.message }}
                  </q-item-label>
                  <q-item-label caption>
                    {{ formatTime(notification.createdAt) }}
                  </q-item-label>
                </q-item-section>
                
                <q-item-section side>
                  <q-btn
                    flat
                    round
                    dense
                    icon="close"
                    size="sm"
                    @click.stop="deleteNotification(notification.id)"
                  />
                </q-item-section>
              </q-item>
            </template>
            
            <q-item v-else>
              <q-item-section class="text-center text-grey-6">
                No notifications
              </q-item-section>
            </q-item>
            
            <q-separator />
            
            <q-item clickable v-close-popup @click="$router.push('/notifications')">
              <q-item-section class="text-center">
                View all notifications
              </q-item-section>
            </q-item>
          </q-list>
        </q-menu>
      </q-btn>

      <!-- User menu -->
      <q-btn
        flat
        round
        dense
        class="user-button"
      >
        <q-avatar size="32px" class="user-avatar">
          <img v-if="user?.avatar" :src="user.avatar" :alt="user.name" />
          <q-icon v-else name="person" />
        </q-avatar>
        
        <q-menu
          fit
          anchor="bottom right"
          self="top right"
          class="user-menu"
        >
          <q-list style="min-width: 200px">
            <q-item>
              <q-item-section avatar>
                <q-avatar>
                  <img v-if="user?.avatar" :src="user.avatar" :alt="user.name" />
                  <q-icon v-else name="person" />
                </q-avatar>
              </q-item-section>
              
              <q-item-section>
                <q-item-label>{{ user?.name }}</q-item-label>
                <q-item-label caption>{{ user?.email }}</q-item-label>
              </q-item-section>
            </q-item>
            
            <q-separator />
            
            <q-item clickable v-close-popup @click="$router.push('/profile')">
              <q-item-section avatar>
                <q-icon name="person" />
              </q-item-section>
              <q-item-section>Profile</q-item-section>
            </q-item>
            
            <q-item clickable v-close-popup @click="$router.push('/settings')">
              <q-item-section avatar>
                <q-icon name="settings" />
              </q-item-section>
              <q-item-section>Settings</q-item-section>
            </q-item>
            
            <q-separator />
            
            <q-item clickable v-close-popup @click="handleLogout">
              <q-item-section avatar>
                <q-icon name="logout" />
              </q-item-section>
              <q-item-section>Logout</q-item-section>
            </q-item>
          </q-list>
        </q-menu>
      </q-btn>
    </div>
  </q-toolbar>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useNotificationStore } from '@/stores/notifications'
import { formatDistanceToNow } from 'date-fns'
import type { User, Notification } from '@/types'

interface Props {
  user: User | null
  notifications: Notification[]
}

interface Emits {
  (e: 'logout'): void
  (e: 'toggle-drawer'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const router = useRouter()
const notificationStore = useNotificationStore()

// Reactive state
const searchQuery = ref('')

// Computed properties
const unreadCount = computed(() => 
  props.notifications.filter(n => !n.read).length
)

// Methods
const handleSearch = () => {
  if (searchQuery.value.trim()) {
    router.push({
      path: '/rules',
      query: { search: searchQuery.value }
    })
  }
}

const handleNotificationClick = (notification: Notification) => {
  // Mark as read
  notificationStore.markAsRead(notification.id)
  
  // Navigate if action URL provided
  if (notification.actionUrl) {
    router.push(notification.actionUrl)
  }
}

const markAllAsRead = () => {
  notificationStore.markAllAsRead()
}

const deleteNotification = (id: string) => {
  notificationStore.deleteNotification(id)
}

const handleLogout = () => {
  emit('logout')
}

const getNotificationIcon = (type: string) => {
  const icons: Record<string, string> = {
    'INFO': 'info',
    'SUCCESS': 'check_circle',
    'WARNING': 'warning',
    'ERROR': 'error'
  }
  return icons[type] || 'info'
}

const getNotificationColor = (type: string) => {
  const colors: Record<string, string> = {
    'INFO': 'info',
    'SUCCESS': 'positive',
    'WARNING': 'warning',
    'ERROR': 'negative'
  }
  return colors[type] || 'info'
}

const formatTime = (dateString: string) => {
  return formatDistanceToNow(new Date(dateString), { addSuffix: true })
}
</script>

<style lang="scss" scoped>
.app-header-toolbar {
  background: linear-gradient(135deg, var(--carrefour-blue) 0%, var(--carrefour-dark-blue) 100%);
  color: white;
  min-height: 64px;
  padding: 0 16px;
}

.app-header-brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.app-header-logo {
  width: 32px;
  height: 32px;
  filter: brightness(0) invert(1); // Make logo white
}

.app-header-title {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
}

.app-header-title-main {
  font-size: 18px;
  font-weight: 600;
  color: white;
}

.app-header-title-sub {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.8);
  font-weight: 400;
}

.app-header-search {
  max-width: 400px;
  width: 100%;
  margin: 0 24px;
}

.search-input {
  :deep(.q-field__control) {
    background-color: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    color: white;
  }
  
  :deep(.q-field__native) {
    color: white;
    
    &::placeholder {
      color: rgba(255, 255, 255, 0.7);
    }
  }
  
  :deep(.q-field__prepend) {
    color: rgba(255, 255, 255, 0.7);
  }
  
  :deep(.q-field__append) {
    color: rgba(255, 255, 255, 0.7);
  }
}

.app-header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.action-button {
  color: white;
  
  &:hover {
    background-color: rgba(255, 255, 255, 0.1);
  }
  
  &.has-notifications {
    color: var(--carrefour-yellow);
  }
}

.user-button {
  color: white;
  
  &:hover {
    background-color: rgba(255, 255, 255, 0.1);
  }
}

.user-avatar {
  border: 2px solid rgba(255, 255, 255, 0.3);
}

.notifications-menu {
  :deep(.q-list) {
    max-height: 400px;
    overflow-y: auto;
  }
}

.user-menu {
  :deep(.q-list) {
    padding: 8px 0;
  }
}

// Notification item styles
:deep(.q-item.unread) {
  background-color: var(--carrefour-light-blue);
  border-left: 3px solid var(--carrefour-blue);
}

// Responsive design
@media (max-width: 768px) {
  .app-header-search {
    display: none;
  }
  
  .app-header-title-sub {
    display: none;
  }
}

@media (max-width: 480px) {
  .app-header-brand {
    gap: 8px;
  }
  
  .app-header-title-main {
    font-size: 16px;
  }
}
</style>
