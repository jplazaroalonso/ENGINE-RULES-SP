<template>
  <q-page class="analytics-page">
    <div class="page-header">
      <div class="page-header-content">
        <div class="page-title-section">
          <h1 class="page-title">Analytics Dashboard</h1>
          <p class="page-subtitle">
            Monitor system performance and rule execution metrics
          </p>
        </div>
        
        <div class="page-actions">
          <q-btn
            color="primary"
            icon="refresh"
            label="Refresh"
            @click="refreshMetrics"
            :loading="loading"
          />
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading && !hasMetrics" class="loading-container">
      <q-spinner-dots size="50px" color="primary" />
      <p>Loading analytics data...</p>
    </div>

    <!-- Error State -->
    <q-banner v-if="error" class="bg-negative text-white">
      <template v-slot:avatar>
        <q-icon name="error" />
      </template>
      {{ error }}
      <template v-slot:action>
        <q-btn flat label="Retry" @click="fetchMetrics" />
        <q-btn flat label="Dismiss" @click="clearError" />
      </template>
    </q-banner>

    <!-- Analytics Content -->
    <div v-if="hasMetrics" class="analytics-content">
      <!-- System Health Overview -->
      <div class="metrics-grid">
        <q-card class="metric-card">
          <q-card-section>
            <div class="metric-header">
              <q-icon name="health_and_safety" size="24px" color="positive" />
              <span class="metric-title">System Health</span>
            </div>
            <div class="metric-value">
              {{ metrics?.systemHealth.totalUptime.toFixed(1) }}%
            </div>
            <div class="metric-subtitle">Overall Uptime</div>
          </q-card-section>
        </q-card>

        <q-card class="metric-card">
          <q-card-section>
            <div class="metric-header">
              <q-icon name="rule" size="24px" color="primary" />
              <span class="metric-title">Active Rules</span>
            </div>
            <div class="metric-value">
              {{ metrics?.rulesMetrics.activeRules }}/{{ metrics?.rulesMetrics.totalRules }}
            </div>
            <div class="metric-subtitle">Rules Active</div>
          </q-card-section>
        </q-card>

        <q-card class="metric-card">
          <q-card-section>
            <div class="metric-header">
              <q-icon name="play_arrow" size="24px" color="info" />
              <span class="metric-title">Executions Today</span>
            </div>
            <div class="metric-value">
              {{ metrics?.rulesMetrics.executionsToday.toLocaleString() }}
            </div>
            <div class="metric-subtitle">Rule Executions</div>
          </q-card-section>
        </q-card>

        <q-card class="metric-card">
          <q-card-section>
            <div class="metric-header">
              <q-icon name="speed" size="24px" color="warning" />
              <span class="metric-title">Avg Response Time</span>
            </div>
            <div class="metric-value">
              {{ metrics?.performanceMetrics.averageResponseTime }}ms
            </div>
            <div class="metric-subtitle">System Performance</div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Service Status -->
      <q-card class="service-status-card">
        <q-card-section>
          <div class="card-header">
            <h3>Service Status</h3>
            <q-chip 
              :color="isHealthy ? 'positive' : 'negative'"
              :icon="isHealthy ? 'check_circle' : 'error'"
              :label="isHealthy ? 'All Systems Operational' : 'Issues Detected'"
            />
          </div>
          
          <div class="service-list">
            <div 
              v-for="service in metrics?.systemHealth.serviceStatus" 
              :key="service.service"
              class="service-item"
            >
              <div class="service-info">
                <q-icon 
                  :name="service.status === 'healthy' ? 'check_circle' : 'error'"
                  :color="service.status === 'healthy' ? 'positive' : 'negative'"
                />
                <span class="service-name">{{ service.service }}</span>
              </div>
              <div class="service-metrics">
                <span class="response-time">{{ service.responseTime }}ms</span>
                <span class="uptime">{{ service.uptime }}% uptime</span>
              </div>
            </div>
          </div>
        </q-card-section>
      </q-card>

      <!-- Rules Performance -->
      <q-card class="rules-performance-card">
        <q-card-section>
          <div class="card-header">
            <h3>Rules Performance</h3>
            <span class="last-updated">
              Last updated: {{ lastUpdated ? new Date(lastUpdated).toLocaleTimeString() : 'Never' }}
            </span>
          </div>
          
          <div class="performance-grid">
            <div class="performance-item">
              <div class="performance-label">Total Rules</div>
              <div class="performance-value">{{ metrics?.rulesMetrics.totalRules }}</div>
            </div>
            <div class="performance-item">
              <div class="performance-label">Active Rules</div>
              <div class="performance-value">{{ metrics?.rulesMetrics.activeRules }}</div>
            </div>
            <div class="performance-item">
              <div class="performance-label">Executions Today</div>
              <div class="performance-value">{{ metrics?.rulesMetrics.executionsToday.toLocaleString() }}</div>
            </div>
            <div class="performance-item">
              <div class="performance-label">Avg Execution Time</div>
              <div class="performance-value">{{ metrics?.rulesMetrics.averageExecutionTime }}ms</div>
            </div>
            <div class="performance-item">
              <div class="performance-label">Error Rate</div>
              <div class="performance-value">{{ (metrics?.rulesMetrics.errorRate * 100).toFixed(2) }}%</div>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAnalyticsStore } from '@/stores/analytics'

const analyticsStore = useAnalyticsStore()

const {
  metrics,
  loading,
  error,
  lastUpdated,
  hasMetrics,
  isHealthy,
  fetchMetrics,
  refreshMetrics,
  clearError
} = analyticsStore

onMounted(() => {
  fetchMetrics()
})
</script>

<style lang="scss" scoped>
.analytics-page {
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
  
  .page-header-content {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
  }
  
  .page-title-section {
    flex: 1;
  }
  
  .page-actions {
    display: flex;
    gap: 12px;
  }
  
  h1 {
    margin: 0 0 8px 0;
    color: var(--carrefour-dark-blue);
  }
  
  p {
    margin: 0;
    color: var(--carrefour-gray-600);
  }
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 24px;
  
  p {
    margin-top: 16px;
    color: var(--carrefour-gray-600);
  }
}

.analytics-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 16px;
}

.metric-card {
  .metric-header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 12px;
  }
  
  .metric-title {
    font-weight: 500;
    color: var(--carrefour-gray-700);
  }
  
  .metric-value {
    font-size: 2rem;
    font-weight: bold;
    color: var(--carrefour-blue);
    margin-bottom: 4px;
  }
  
  .metric-subtitle {
    font-size: 0.875rem;
    color: var(--carrefour-gray-600);
  }
}

.service-status-card,
.rules-performance-card {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    
    h3 {
      margin: 0;
      color: var(--carrefour-dark-blue);
    }
    
    .last-updated {
      font-size: 0.875rem;
      color: var(--carrefour-gray-600);
    }
  }
}

.service-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.service-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background-color: var(--carrefour-gray-50);
  border-radius: 8px;
  
  .service-info {
    display: flex;
    align-items: center;
    gap: 12px;
    
    .service-name {
      font-weight: 500;
      color: var(--carrefour-gray-800);
    }
  }
  
  .service-metrics {
    display: flex;
    gap: 16px;
    font-size: 0.875rem;
    color: var(--carrefour-gray-600);
  }
}

.performance-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.performance-item {
  text-align: center;
  padding: 16px;
  background-color: var(--carrefour-gray-50);
  border-radius: 8px;
  
  .performance-label {
    font-size: 0.875rem;
    color: var(--carrefour-gray-600);
    margin-bottom: 8px;
  }
  
  .performance-value {
    font-size: 1.5rem;
    font-weight: bold;
    color: var(--carrefour-blue);
  }
}

@media (max-width: 768px) {
  .analytics-page {
    padding: 16px;
  }
  
  .page-header-content {
    flex-direction: column;
    gap: 16px;
  }
  
  .metrics-grid {
    grid-template-columns: 1fr;
  }
  
  .performance-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .service-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}
</style>
