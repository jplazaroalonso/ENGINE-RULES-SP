# Dashboard Screen - Main Landing Page

## Overview
The Dashboard serves as the main landing page for business users, providing an overview of system status, key metrics, and quick access to frequently used features.

## Screen Layout

### Header Section
```vue
<template>
  <div class="dashboard-header">
    <div class="welcome-section">
      <h1>Welcome back, {{ userStore.currentUser?.name }}</h1>
      <p class="subtitle">{{ getCurrentTimeGreeting() }} - Here's what's happening with your rules today</p>
    </div>
    
    <div class="quick-actions">
      <q-btn 
        color="primary" 
        icon="add" 
        label="Create Rule" 
        @click="navigateToCreateRule"
        class="q-mr-sm"
      />
      <q-btn 
        color="secondary" 
        icon="campaign" 
        label="New Campaign" 
        @click="navigateToCreateCampaign"
        outline
      />
    </div>
  </div>
</template>
```

### Key Metrics Cards
```vue
<template>
  <div class="metrics-grid">
    <MetricCard
      title="Active Rules"
      :value="metrics.activeRules"
      icon="rule"
      color="primary"
      :trend="metrics.activeRulesTrend"
      @click="navigateToRules({ status: 'ACTIVE' })"
    />
    
    <MetricCard
      title="Rules Pending Approval"
      :value="metrics.pendingApproval"
      icon="pending_actions"
      color="warning"
      :trend="metrics.pendingApprovalTrend"
      @click="navigateToApprovalQueue"
    />
    
    <MetricCard
      title="Campaign Performance"
      :value="formatPercentage(metrics.campaignPerformance)"
      icon="trending_up"
      color="positive"
      :trend="metrics.campaignPerformanceTrend"
      @click="navigateToAnalytics"
    />
    
    <MetricCard
      title="Revenue Impact"
      :value="formatCurrency(metrics.revenueImpact)"
      icon="monetization_on"
      color="accent"
      :trend="metrics.revenueImpactTrend"
      @click="navigateToRevenueReport"
    />
  </div>
</template>
```

### Recent Activity Section
```vue
<template>
  <div class="recent-activity">
    <div class="section-header">
      <h3>Recent Activity</h3>
      <q-btn flat label="View All" @click="navigateToActivity" />
    </div>
    
    <q-timeline color="primary">
      <q-timeline-entry
        v-for="activity in recentActivities"
        :key="activity.id"
        :title="activity.title"
        :subtitle="formatRelativeTime(activity.timestamp)"
        :icon="getActivityIcon(activity.type)"
        :body="activity.description"
      >
        <div class="activity-metadata">
          <q-chip 
            :color="getActivityColor(activity.type)" 
            text-color="white" 
            size="sm"
          >
            {{ activity.type }}
          </q-chip>
          <span class="activity-user">by {{ activity.user }}</span>
        </div>
      </q-timeline-entry>
    </q-timeline>
  </div>
</template>
```

### Quick Insights Charts
```vue
<template>
  <div class="insights-section">
    <div class="row q-gutter-md">
      <div class="col-md-6 col-12">
        <q-card class="chart-card">
          <q-card-section>
            <h4>Rule Execution Performance (Last 7 Days)</h4>
            <LineChart
              :data="performanceChartData"
              :options="performanceChartOptions"
              height="300"
            />
          </q-card-section>
        </q-card>
      </div>
      
      <div class="col-md-6 col-12">
        <q-card class="chart-card">
          <q-card-section>
            <h4>Rule Types Distribution</h4>
            <PieChart
              :data="distributionChartData"
              :options="distributionChartOptions"
              height="300"
            />
          </q-card-section>
        </q-card>
      </div>
    </div>
  </div>
</template>
```

### Alerts and Notifications
```vue
<template>
  <div class="alerts-section" v-if="alerts.length > 0">
    <div class="section-header">
      <h3>System Alerts</h3>
      <q-btn flat label="Dismiss All" @click="dismissAllAlerts" />
    </div>
    
    <q-banner
      v-for="alert in alerts"
      :key="alert.id"
      :class="`alert-${alert.severity}`"
      :icon="getAlertIcon(alert.severity)"
      dense
    >
      <template v-slot:avatar>
        <q-icon :name="getAlertIcon(alert.severity)" />
      </template>
      
      <div class="alert-content">
        <div class="alert-title">{{ alert.title }}</div>
        <div class="alert-message">{{ alert.message }}</div>
        <div class="alert-timestamp">{{ formatTimestamp(alert.timestamp) }}</div>
      </div>
      
      <template v-slot:action>
        <q-btn 
          flat 
          dense 
          label="View" 
          @click="viewAlert(alert)" 
          v-if="alert.actionUrl"
        />
        <q-btn 
          flat 
          dense 
          icon="close" 
          @click="dismissAlert(alert.id)" 
        />
      </template>
    </q-banner>
  </div>
</template>
```

## Component Data Structure

### Dashboard Data Interface
```typescript
interface DashboardData {
  metrics: {
    activeRules: number
    pendingApproval: number
    campaignPerformance: number
    revenueImpact: number
    activeRulesTrend: TrendData
    pendingApprovalTrend: TrendData
    campaignPerformanceTrend: TrendData
    revenueImpactTrend: TrendData
  }
  
  recentActivities: Activity[]
  alerts: SystemAlert[]
  performanceChartData: ChartData
  distributionChartData: ChartData
}

interface TrendData {
  direction: 'up' | 'down' | 'stable'
  percentage: number
  period: string
}

interface Activity {
  id: string
  type: 'rule_created' | 'rule_approved' | 'campaign_launched' | 'rule_executed'
  title: string
  description: string
  user: string
  timestamp: string
}

interface SystemAlert {
  id: string
  severity: 'info' | 'warning' | 'error' | 'critical'
  title: string
  message: string
  timestamp: string
  actionUrl?: string
  dismissed: boolean
}
```

## Screen Interactions

### Navigation Actions
- **Create Rule**: Navigate to `/rules/create`
- **New Campaign**: Navigate to `/campaigns/create`
- **View Active Rules**: Navigate to `/rules?status=ACTIVE`
- **Approval Queue**: Navigate to `/approvals`
- **Analytics**: Navigate to `/analytics`

### Real-time Updates
- Metrics update every 30 seconds via WebSocket
- Recent activities stream in real-time
- System alerts appear immediately when triggered

### Performance Considerations
- Lazy load chart data on component mount
- Cache dashboard data for 5 minutes
- Use skeleton loaders while data loads
- Optimize chart rendering for large datasets

## Responsive Design
- **Desktop (>1024px)**: 4-column metric cards, side-by-side charts
- **Tablet (768-1024px)**: 2-column metric cards, stacked charts
- **Mobile (<768px)**: 1-column layout, simplified chart views

## Accessibility Features
- Keyboard navigation for all interactive elements
- ARIA labels for screen readers
- High contrast mode support
- Reduced motion respect for chart animations

## State Management
```typescript
// Dashboard Store
export const useDashboardStore = defineStore('dashboard', () => {
  const metrics = ref<DashboardMetrics>()
  const activities = ref<Activity[]>([])
  const alerts = ref<SystemAlert[]>([])
  const loading = ref(false)
  
  const fetchDashboardData = async () => {
    loading.value = true
    try {
      const [metricsData, activitiesData, alertsData] = await Promise.all([
        dashboardApi.getMetrics(),
        dashboardApi.getRecentActivities(),
        dashboardApi.getAlerts()
      ])
      
      metrics.value = metricsData
      activities.value = activitiesData
      alerts.value = alertsData
    } finally {
      loading.value = false
    }
  }
  
  return {
    metrics: readonly(metrics),
    activities: readonly(activities),
    alerts: readonly(alerts),
    loading: readonly(loading),
    fetchDashboardData
  }
})
```

## Testing Strategy
- Unit tests for all computed properties and methods
- Integration tests for API data fetching
- E2E tests for navigation flows
- Accessibility testing with axe-core
- Performance testing for chart rendering
