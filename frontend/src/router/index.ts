import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// Route components
const Dashboard = () => import('@/views/Dashboard.vue')
const RulesList = () => import('@/views/rules/RulesList.vue')
const RuleEditor = () => import('@/views/rules/RuleEditor.vue')
const RuleDetail = () => import('@/views/rules/RuleDetails.vue')
const CampaignsList = () => import('@/views/campaigns/CampaignsList.vue')
const CampaignEditor = () => import('@/views/campaigns/CampaignEditor.vue')
const Analytics = () => import('@/views/Analytics.vue')
const Customers = () => import('@/views/Customers.vue')
const Settings = () => import('@/views/Settings.vue')
const Login = () => import('@/views/auth/Login.vue')

const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { 
      requiresAuth: false,
      hideForAuth: true 
    }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard,
    meta: { 
      requiresAuth: false,
      title: 'Dashboard',
      icon: 'dashboard'
    }
  },
  {
    path: '/rules',
    name: 'RulesList',
    component: RulesList,
    meta: { 
      requiresAuth: false,
      title: 'Rules Management',
      icon: 'rule'
    }
  },
  {
    path: '/rules/create',
    name: 'RuleCreate',
    component: RuleEditor,
    meta: { 
      requiresAuth: false,
      title: 'Create Rule',
      icon: 'add_circle'
    }
  },
  {
    path: '/rules/:id/edit',
    name: 'RuleEdit',
    component: RuleEditor,
    meta: { 
      requiresAuth: false,
      title: 'Edit Rule',
      icon: 'edit'
    }
  },
  {
    path: '/rules/:id',
    name: 'RuleDetail',
    component: RuleDetail,
    meta: { 
      requiresAuth: false,
      title: 'Rule Details',
      icon: 'visibility'
    }
  },
  {
    path: '/campaigns',
    name: 'CampaignsList',
    component: CampaignsList,
    meta: { 
      requiresAuth: false,
      title: 'Campaigns',
      icon: 'campaign'
    }
  },
  {
    path: '/campaigns/create',
    name: 'CampaignCreate',
    component: CampaignEditor,
    meta: { 
      requiresAuth: false,
      title: 'Create Campaign',
      icon: 'add_circle'
    }
  },
  {
    path: '/campaigns/:id/edit',
    name: 'CampaignEdit',
    component: CampaignEditor,
    meta: { 
      requiresAuth: false,
      title: 'Edit Campaign',
      icon: 'edit'
    }
  },
  {
    path: '/analytics',
    name: 'Analytics',
    component: Analytics,
    meta: { 
      requiresAuth: false,
      title: 'Analytics',
      icon: 'analytics'
    }
  },
  {
    path: '/customers',
    name: 'Customers',
    component: Customers,
    meta: { 
      requiresAuth: false,
      title: 'Customers',
      icon: 'people'
    }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: Settings,
    meta: { 
      requiresAuth: false,
      title: 'Settings',
      icon: 'settings'
    }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    redirect: '/dashboard'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// Navigation guards
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // Check if route requires authentication
  if (to.meta.requiresAuth) {
    if (!authStore.isAuthenticated) {
      next('/login')
      return
    }
  }
  
  // Redirect authenticated users away from login page
  if (to.meta.hideForAuth && authStore.isAuthenticated) {
    next('/dashboard')
    return
  }
  
  // Set page title
  if (to.meta.title) {
    document.title = `${to.meta.title} - Rules Engine`
  }
  
  next()
})

export default router
