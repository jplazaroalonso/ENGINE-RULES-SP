<template>
  <div class="main-navigation">
    <!-- User profile section -->
    <div class="nav-user-section">
      <div class="nav-user-info">
        <q-avatar size="48px" class="nav-user-avatar">
          <img v-if="user?.avatar" :src="user.avatar" :alt="user.name" />
          <q-icon v-else name="person" size="24px" />
        </q-avatar>
        <div class="nav-user-details">
          <div class="nav-user-name">{{ user?.name }}</div>
          <div class="nav-user-role">{{ formatUserRole(user?.role) }}</div>
        </div>
      </div>
    </div>

    <q-separator class="nav-separator" />

    <!-- Navigation menu -->
    <q-list class="nav-menu">
      <q-item
        v-for="item in navigationItems"
        :key="item.name"
        clickable
        v-ripple
        :active="currentRoute === item.name"
        class="nav-item"
        @click="handleNavigation(item.path)"
      >
        <q-item-section avatar>
          <q-icon 
            :name="item.icon" 
            :color="currentRoute === item.name ? 'primary' : 'grey-6'"
            size="20px"
          />
        </q-item-section>

        <q-item-section>
          <q-item-label class="nav-item-label">{{ item.label }}</q-item-label>
        </q-item-section>

        <!-- Badge for notifications or counts -->
        <q-item-section side v-if="item.badge">
          <q-badge 
            :color="item.badgeColor || 'primary'" 
            :label="item.badge"
            rounded
          />
        </q-item-section>
      </q-item>
    </q-list>

    <q-separator class="nav-separator" />

    <!-- Quick actions -->
    <div class="nav-quick-actions">
      <div class="nav-section-title">Quick Actions</div>
      <q-list class="nav-menu">
        <q-item
          v-for="action in quickActions"
          :key="action.name"
          clickable
          v-ripple
          class="nav-item nav-item-small"
          @click="handleNavigation(action.path)"
        >
          <q-item-section avatar>
            <q-icon 
              :name="action.icon" 
              color="grey-6"
              size="18px"
            />
          </q-item-section>

          <q-item-section>
            <q-item-label class="nav-item-label-small">{{ action.label }}</q-item-label>
          </q-item-section>
        </q-item>
      </q-list>
    </div>

    <q-space />

    <!-- Footer section -->
    <div class="nav-footer">
      <q-list class="nav-menu">
        <q-item
          clickable
          v-ripple
          class="nav-item nav-item-small"
          @click="handleNavigation('/help')"
        >
          <q-item-section avatar>
            <q-icon name="help" color="grey-6" size="18px" />
          </q-item-section>

          <q-item-section>
            <q-item-label class="nav-item-label-small">Help & Support</q-item-label>
          </q-item-section>
        </q-item>

        <q-item
          clickable
          v-ripple
          class="nav-item nav-item-small"
          @click="toggleTheme"
        >
          <q-item-section avatar>
            <q-icon 
              :name="isDarkMode ? 'light_mode' : 'dark_mode'" 
              color="grey-6" 
              size="18px" 
            />
          </q-item-section>

          <q-item-section>
            <q-item-label class="nav-item-label-small">
              {{ isDarkMode ? 'Light Mode' : 'Dark Mode' }}
            </q-item-label>
          </q-item-section>
        </q-item>
      </q-list>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUIStore } from '@/stores/ui'
import { useRulesStore } from '@/stores/rules'
import type { User, UserRole } from '@/types'

interface Props {
  user: User | null
  currentRoute: string | symbol | null | undefined
}

interface Emits {
  (e: 'navigate', path: string): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const router = useRouter()
const uiStore = useUIStore()
const rulesStore = useRulesStore()

// Computed properties
const isDarkMode = computed(() => uiStore.isDarkMode)

const navigationItems = computed(() => [
  {
    name: 'Dashboard',
    label: 'Dashboard',
    icon: 'dashboard',
    path: '/dashboard'
  },
  {
    name: 'RulesList',
    label: 'Rules Management',
    icon: 'rule',
    path: '/rules',
    badge: rulesStore.rulesStats.underReview > 0 ? rulesStore.rulesStats.underReview : undefined,
    badgeColor: 'warning'
  },
  {
    name: 'CampaignsList',
    label: 'Campaigns',
    icon: 'campaign',
    path: '/campaigns'
  },
  {
    name: 'Analytics',
    label: 'Analytics',
    icon: 'analytics',
    path: '/analytics'
  },
  {
    name: 'Customers',
    label: 'Customers',
    icon: 'people',
    path: '/customers'
  },
  {
    name: 'Settings',
    label: 'Settings',
    icon: 'settings',
    path: '/settings'
  }
])

const quickActions = computed(() => [
  {
    name: 'CreateRule',
    label: 'Create Rule',
    icon: 'add_circle',
    path: '/rules/create'
  },
  {
    name: 'CreateCampaign',
    label: 'Create Campaign',
    icon: 'add_circle',
    path: '/campaigns/create'
  },
  {
    name: 'ViewReports',
    label: 'View Reports',
    icon: 'assessment',
    path: '/analytics'
  }
])

// Methods
const handleNavigation = (path: string) => {
  emit('navigate', path)
}

const toggleTheme = () => {
  uiStore.toggleTheme()
}

const formatUserRole = (role?: UserRole) => {
  if (!role) return 'User'
  
  const roleLabels: Record<UserRole, string> = {
    'ADMIN': 'Administrator',
    'MANAGER': 'Manager',
    'USER': 'User',
    'VIEWER': 'Viewer'
  }
  
  return roleLabels[role]
}
</script>

<style lang="scss" scoped>
.main-navigation {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: white;
}

.nav-user-section {
  padding: 20px 16px;
  background: linear-gradient(135deg, var(--carrefour-light-blue) 0%, rgba(0, 75, 135, 0.05) 100%);
}

.nav-user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.nav-user-avatar {
  border: 2px solid var(--carrefour-blue);
}

.nav-user-details {
  flex: 1;
  min-width: 0;
}

.nav-user-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--carrefour-dark-blue);
  line-height: 1.2;
}

.nav-user-role {
  font-size: 12px;
  color: var(--carrefour-gray-600);
  margin-top: 2px;
}

.nav-separator {
  margin: 8px 0;
  background-color: var(--carrefour-gray-200);
}

.nav-menu {
  padding: 0;
}

.nav-item {
  margin: 2px 8px;
  border-radius: 8px;
  transition: all 0.2s ease;
  
  &:hover {
    background-color: var(--carrefour-light-blue);
  }
  
  &.q-item--active {
    background-color: var(--carrefour-blue);
    color: white;
    
    .nav-item-label {
      color: white;
      font-weight: 500;
    }
  }
}

.nav-item-small {
  margin: 1px 8px;
  padding: 8px 16px;
}

.nav-item-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--carrefour-gray-800);
  transition: color 0.2s ease;
}

.nav-item-label-small {
  font-size: 13px;
  color: var(--carrefour-gray-600);
  transition: color 0.2s ease;
}

.nav-quick-actions {
  padding: 16px 0 8px;
}

.nav-section-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--carrefour-gray-500);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 0 24px 8px;
}

.nav-footer {
  padding: 8px 0 16px;
  border-top: 1px solid var(--carrefour-gray-200);
  margin-top: auto;
}

// Active state styles
:deep(.q-item--active) {
  .q-icon {
    color: white !important;
  }
}

// Hover effects
.nav-item:hover {
  .nav-item-label,
  .nav-item-label-small {
    color: var(--carrefour-blue);
  }
  
  .q-icon {
    color: var(--carrefour-blue) !important;
  }
}

// Badge styles
:deep(.q-badge) {
  font-size: 10px;
  min-width: 18px;
  height: 18px;
  line-height: 18px;
}

// Responsive design
@media (max-width: 1024px) {
  .nav-user-section {
    padding: 16px 12px;
  }
  
  .nav-user-name {
    font-size: 14px;
  }
  
  .nav-user-role {
    font-size: 11px;
  }
}

@media (max-width: 768px) {
  .nav-user-details {
    display: none;
  }
  
  .nav-user-info {
    justify-content: center;
  }
  
  .nav-section-title {
    display: none;
  }
  
  .nav-item-label,
  .nav-item-label-small {
    display: none;
  }
  
  .nav-item {
    justify-content: center;
  }
}
</style>
