<template>
  <q-page class="dashboard-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="page-header-content">
        <div class="page-title-section">
          <h1 class="page-title">Dashboard</h1>
          <p class="page-subtitle">
            Welcome back, {{ user?.name }}! Here's what's happening with your rules engine.
          </p>
        </div>
        
        <div class="page-actions">
          <q-btn
            color="primary"
            icon="add"
            label="Create Rule"
            @click="$router.push('/rules/create')"
          />
          <q-btn
            color="secondary"
            icon="campaign"
            label="New Campaign"
            @click="$router.push('/campaigns/create')"
          />
        </div>
      </div>
    </div>

    <!-- Dashboard Content -->
    <div class="dashboard-content">
      <!-- Key Metrics Row -->
      <div class="metrics-row">
        <div class="row q-gutter-md">
          <div class="col-12 col-sm-6 col-md-3">
            <MetricCard
              title="Total Rules"
              :value="metrics.totalRules"
              icon="rule"
              color="primary"
              :trend="metrics.rulesTrend"
              @click="$router.push('/rules')"
            />
          </div>
          
          <div class="col-12 col-sm-6 col-md-3">
            <MetricCard
              title="Active Rules"
              :value="metrics.activeRules"
              icon="play_circle"
              color="positive"
              :trend="metrics.activeTrend"
              @click="$router.push('/rules?status=ACTIVE')"
            />
          </div>
          
          <div class="col-12 col-sm-6 col-md-3">
            <MetricCard
              title="Pending Approval"
              :value="metrics.pendingApproval"
              icon="pending"
              color="warning"
              :trend="metrics.pendingTrend"
              @click="$router.push('/rules?status=UNDER_REVIEW')"
            />
          </div>
          
          <div class="col-12 col-sm-6 col-md-3">
            <MetricCard
              title="Active Campaigns"
              :value="metrics.activeCampaigns"
              icon="campaign"
              color="info"
              :trend="metrics.campaignsTrend"
              @click="$router.push('/campaigns')"
            />
          </div>
        </div>
      </div>

      <!-- Main Content Row -->
      <div class="row q-gutter-md q-mt-md">
        <!-- Left Column -->
        <div class="col-12 col-lg-8">
          <!-- Recent Rules -->
          <q-card class="dashboard-card">
            <q-card-section>
              <div class="card-header">
                <div class="card-title">
                  <q-icon name="rule" color="primary" size="20px" />
                  <span>Recent Rules</span>
                </div>
                <q-btn
                  flat
                  dense
                  label="View All"
                  color="primary"
                  @click="$router.push('/rules')"
                />
              </div>
            </q-card-section>
            
            <q-card-section class="q-pt-none">
              <div v-if="recentRules.length > 0">
                <div
                  v-for="rule in recentRules"
                  :key="rule.id"
                  class="rule-item"
                  @click="$router.push(`/rules/${rule.id}`)"
                >
                  <div class="rule-item-content">
                    <div class="rule-name">{{ rule.name }}</div>
                    <div class="rule-description">{{ rule.description }}</div>
                    <div class="rule-meta">
                      <q-chip
                        :color="getStatusColor(rule.status)"
                        text-color="white"
                        dense
                        size="sm"
                      >
                        {{ rule.status }}
                      </q-chip>
                      <span class="rule-date">{{ formatDate(rule.created_at) }}</span>
                    </div>
                  </div>
                  <q-icon name="chevron_right" color="grey-5" />
                </div>
              </div>
              
              <div v-else class="empty-state">
                <q-icon name="rule" size="48px" color="grey-4" />
                <p>No rules created yet</p>
                <q-btn
                  color="primary"
                  label="Create your first rule"
                  @click="$router.push('/rules/create')"
                />
              </div>
            </q-card-section>
          </q-card>

          <!-- Performance Chart -->
          <q-card class="dashboard-card q-mt-md">
            <q-card-section>
              <div class="card-header">
                <div class="card-title">
                  <q-icon name="trending_up" color="primary" size="20px" />
                  <span>Rule Performance</span>
                </div>
                <q-btn-toggle
                  v-model="chartPeriod"
                  :options="chartPeriodOptions"
                  dense
                  color="primary"
                  toggle-color="primary"
                />
              </div>
            </q-card-section>
            
            <q-card-section class="q-pt-none">
              <div class="chart-container">
                <PerformanceChart
                  :data="chartData"
                  :period="chartPeriod"
                />
              </div>
            </q-card-section>
          </q-card>
        </div>

        <!-- Right Column -->
        <div class="col-12 col-lg-4">
          <!-- Quick Actions -->
          <q-card class="dashboard-card">
            <q-card-section>
              <div class="card-title">
                <q-icon name="flash_on" color="primary" size="20px" />
                <span>Quick Actions</span>
              </div>
            </q-card-section>
            
            <q-card-section class="q-pt-none">
              <div class="quick-actions">
                <q-btn
                  v-for="action in quickActions"
                  :key="action.name"
                  :color="action.color"
                  :icon="action.icon"
                  :label="action.label"
                  class="quick-action-btn"
                  @click="handleQuickAction(action)"
                />
              </div>
            </q-card-section>
          </q-card>

          <!-- System Status -->
          <q-card class="dashboard-card q-mt-md">
            <q-card-section>
              <div class="card-title">
                <q-icon name="monitor" color="primary" size="20px" />
                <span>System Status</span>
              </div>
            </q-card-section>
            
            <q-card-section class="q-pt-none">
              <div class="status-list">
                <div
                  v-for="service in systemServices"
                  :key="service.name"
                  class="status-item"
                >
                  <div class="status-info">
                    <span class="status-name">{{ service.name }}</span>
                    <span class="status-description">{{ service.description }}</span>
                  </div>
                  <q-icon
                    :name="service.status === 'healthy' ? 'check_circle' : 'error'"
                    :color="service.status === 'healthy' ? 'positive' : 'negative'"
                    size="20px"
                  />
                </div>
              </div>
            </q-card-section>
          </q-card>

          <!-- Recent Activity -->
          <q-card class="dashboard-card q-mt-md">
            <q-card-section>
              <div class="card-header">
                <div class="card-title">
                  <q-icon name="history" color="primary" size="20px" />
                  <span>Recent Activity</span>
                </div>
                <q-btn
                  flat
                  dense
                  label="View All"
                  color="primary"
                  @click="$router.push('/activity')"
                />
              </div>
            </q-card-section>
            
            <q-card-section class="q-pt-none">
              <div class="activity-list">
                <div
                  v-for="activity in recentActivity"
                  :key="activity.id"
                  class="activity-item"
                >
                  <q-avatar size="32px" class="activity-avatar">
                    <q-icon :name="activity.icon" :color="activity.color" />
                  </q-avatar>
                  
                  <div class="activity-content">
                    <div class="activity-text">{{ activity.text }}</div>
                    <div class="activity-time">{{ formatTime(activity.timestamp) }}</div>
                  </div>
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRulesStore } from '@/stores/rules'
import { formatDistanceToNow } from 'date-fns'
import MetricCard from '@/components/common/MetricCard.vue'
import PerformanceChart from '@/components/charts/PerformanceChart.vue'
import type { Rule } from '@/types'

const authStore = useAuthStore()
const rulesStore = useRulesStore()

// Reactive state
const chartPeriod = ref('7d')
const chartPeriodOptions = [
  { label: '7D', value: '7d' },
  { label: '30D', value: '30d' },
  { label: '90D', value: '90d' }
]

// Computed properties
const user = computed(() => authStore.user)
const recentRules = computed(() => rulesStore.rules.slice(0, 5))
const metrics = computed(() => ({
  totalRules: rulesStore.rulesStats.total,
  activeRules: rulesStore.rulesStats.active,
  pendingApproval: rulesStore.rulesStats.underReview,
  activeCampaigns: 12, // Mock data
  rulesTrend: '+12%',
  activeTrend: '+8%',
  pendingTrend: '-5%',
  campaignsTrend: '+15%'
}))

const chartData = computed(() => ({
  labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
  datasets: [{
    label: 'Rule Executions',
    data: [120, 190, 300, 500, 200, 300, 450],
    borderColor: '#004B87',
    backgroundColor: 'rgba(0, 75, 135, 0.1)',
    tension: 0.4
  }]
}))

const quickActions = computed(() => [
  {
    name: 'create-rule',
    label: 'Create Rule',
    icon: 'add_circle',
    color: 'primary'
  },
  {
    name: 'test-rule',
    label: 'Test Rule',
    icon: 'play_circle',
    color: 'secondary'
  },
  {
    name: 'view-analytics',
    label: 'View Analytics',
    icon: 'analytics',
    color: 'info'
  },
  {
    name: 'export-data',
    label: 'Export Data',
    icon: 'download',
    color: 'positive'
  }
])

const systemServices = computed(() => [
  {
    name: 'Rules Engine',
    description: 'Core rule processing',
    status: 'healthy'
  },
  {
    name: 'Database',
    description: 'PostgreSQL',
    status: 'healthy'
  },
  {
    name: 'Message Queue',
    description: 'NATS JetStream',
    status: 'healthy'
  },
  {
    name: 'API Gateway',
    description: 'Traefik Ingress',
    status: 'healthy'
  }
])

const recentActivity = computed(() => [
  {
    id: '1',
    text: 'New rule "Senior Discount" was created',
    icon: 'add_circle',
    color: 'primary',
    timestamp: new Date(Date.now() - 1000 * 60 * 5).toISOString()
  },
  {
    id: '2',
    text: 'Campaign "Black Friday" was activated',
    icon: 'campaign',
    color: 'positive',
    timestamp: new Date(Date.now() - 1000 * 60 * 15).toISOString()
  },
  {
    id: '3',
    text: 'Rule "VAT Calculation" was updated',
    icon: 'edit',
    color: 'info',
    timestamp: new Date(Date.now() - 1000 * 60 * 30).toISOString()
  },
  {
    id: '4',
    text: 'System backup completed',
    icon: 'backup',
    color: 'positive',
    timestamp: new Date(Date.now() - 1000 * 60 * 60).toISOString()
  }
])

// Methods
const handleQuickAction = (action: any) => {
  switch (action.name) {
    case 'create-rule':
      $router.push('/rules/create')
      break
    case 'test-rule':
      $router.push('/rules?action=test')
      break
    case 'view-analytics':
      $router.push('/analytics')
      break
    case 'export-data':
      // Handle export
      break
  }
}

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    'ACTIVE': 'positive',
    'DRAFT': 'info',
    'UNDER_REVIEW': 'warning',
    'INACTIVE': 'grey',
    'DEPRECATED': 'negative'
  }
  return colors[status] || 'grey'
}

const formatDate = (dateString: string) => {
  return formatDistanceToNow(new Date(dateString), { addSuffix: true })
}

const formatTime = (dateString: string) => {
  return formatDistanceToNow(new Date(dateString), { addSuffix: true })
}

// Lifecycle
onMounted(() => {
  rulesStore.fetchRules()
})
</script>

<style lang="scss" scoped>
.dashboard-page {
  background-color: var(--carrefour-gray-50);
  min-height: 100vh;
}

.page-header {
  background: white;
  border-bottom: 1px solid var(--carrefour-gray-200);
  padding: 24px;
  margin-bottom: 24px;
}

.page-header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 24px;
}

.page-title-section {
  flex: 1;
}

.page-title {
  font-size: 28px;
  font-weight: 600;
  color: var(--carrefour-dark-blue);
  margin: 0 0 8px 0;
}

.page-subtitle {
  font-size: 16px;
  color: var(--carrefour-gray-600);
  margin: 0;
}

.page-actions {
  display: flex;
  gap: 12px;
  flex-shrink: 0;
}

.dashboard-content {
  padding: 0 24px 24px;
}

.metrics-row {
  margin-bottom: 24px;
}

.dashboard-card {
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 1px solid var(--carrefour-gray-200);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
  color: var(--carrefour-dark-blue);
}

.rule-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border: 1px solid var(--carrefour-gray-200);
  border-radius: 8px;
  margin-bottom: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    border-color: var(--carrefour-blue);
    background-color: var(--carrefour-light-blue);
  }
  
  &:last-child {
    margin-bottom: 0;
  }
}

.rule-item-content {
  flex: 1;
  min-width: 0;
}

.rule-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--carrefour-dark-blue);
  margin-bottom: 4px;
}

.rule-description {
  font-size: 14px;
  color: var(--carrefour-gray-600);
  margin-bottom: 8px;
  line-height: 1.4;
}

.rule-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.rule-date {
  font-size: 12px;
  color: var(--carrefour-gray-500);
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: var(--carrefour-gray-500);
  
  p {
    margin: 16px 0;
    font-size: 16px;
  }
}

.chart-container {
  height: 300px;
  position: relative;
}

.quick-actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.quick-action-btn {
  height: 48px;
  font-weight: 500;
}

.status-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.status-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  background-color: var(--carrefour-gray-50);
  border-radius: 8px;
}

.status-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.status-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--carrefour-dark-blue);
}

.status-description {
  font-size: 12px;
  color: var(--carrefour-gray-600);
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.activity-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.activity-avatar {
  background-color: var(--carrefour-gray-100);
  flex-shrink: 0;
}

.activity-content {
  flex: 1;
  min-width: 0;
}

.activity-text {
  font-size: 14px;
  color: var(--carrefour-gray-800);
  line-height: 1.4;
  margin-bottom: 4px;
}

.activity-time {
  font-size: 12px;
  color: var(--carrefour-gray-500);
}

// Responsive design
@media (max-width: 768px) {
  .page-header {
    padding: 16px;
  }
  
  .page-header-content {
    flex-direction: column;
    gap: 16px;
  }
  
  .page-actions {
    width: 100%;
    justify-content: stretch;
    
    .q-btn {
      flex: 1;
    }
  }
  
  .dashboard-content {
    padding: 0 16px 16px;
  }
  
  .quick-actions {
    grid-template-columns: 1fr;
  }
  
  .rule-item {
    padding: 12px;
  }
  
  .rule-name {
    font-size: 14px;
  }
  
  .rule-description {
    font-size: 13px;
  }
}
</style>
