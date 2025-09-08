<template>
  <q-card 
    class="metric-card"
    :class="{ 'metric-card--clickable': clickable }"
    @click="handleClick"
  >
    <q-card-section class="metric-card-content">
      <!-- Icon and Value -->
      <div class="metric-header">
        <div class="metric-icon" :style="{ backgroundColor: iconBgColor }">
          <q-icon :name="icon" :color="iconColor" size="24px" />
        </div>
        
        <div class="metric-value-section">
          <div class="metric-value">{{ formattedValue }}</div>
          <div v-if="trend" class="metric-trend" :class="trendClass">
            <q-icon :name="trendIcon" size="14px" />
            <span>{{ trend }}</span>
          </div>
        </div>
      </div>
      
      <!-- Title -->
      <div class="metric-title">{{ title }}</div>
      
      <!-- Description -->
      <div v-if="description" class="metric-description">{{ description }}</div>
      
      <!-- Loading overlay -->
      <q-inner-loading :showing="loading">
        <q-spinner-dots size="40px" color="primary" />
      </q-inner-loading>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  title: string
  value: number | string
  icon: string
  color?: string
  trend?: string
  description?: string
  loading?: boolean
  clickable?: boolean
  format?: 'number' | 'currency' | 'percentage' | 'text'
  precision?: number
}

interface Emits {
  (e: 'click'): void
}

const props = withDefaults(defineProps<Props>(), {
  color: 'primary',
  loading: false,
  clickable: false,
  format: 'number',
  precision: 0
})

const emit = defineEmits<Emits>()

// Computed properties
const formattedValue = computed(() => {
  if (typeof props.value === 'string') {
    return props.value
  }
  
  switch (props.format) {
    case 'currency':
      return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'EUR',
        minimumFractionDigits: props.precision,
        maximumFractionDigits: props.precision
      }).format(props.value)
      
    case 'percentage':
      return `${props.value.toFixed(props.precision)}%`
      
    case 'number':
      return new Intl.NumberFormat('en-US', {
        minimumFractionDigits: props.precision,
        maximumFractionDigits: props.precision
      }).format(props.value)
      
    default:
      return props.value.toString()
  }
})

const iconColor = computed(() => {
  const colorMap: Record<string, string> = {
    'primary': 'white',
    'secondary': 'white',
    'positive': 'white',
    'negative': 'white',
    'warning': 'white',
    'info': 'white'
  }
  return colorMap[props.color] || 'white'
})

const iconBgColor = computed(() => {
  const colorMap: Record<string, string> = {
    'primary': 'var(--carrefour-blue)',
    'secondary': 'var(--carrefour-red)',
    'positive': 'var(--carrefour-green)',
    'negative': 'var(--carrefour-red)',
    'warning': 'var(--carrefour-orange)',
    'info': 'var(--carrefour-blue)'
  }
  return colorMap[props.color] || 'var(--carrefour-blue)'
})

const trendClass = computed(() => {
  if (!props.trend) return ''
  
  if (props.trend.startsWith('+')) {
    return 'metric-trend--positive'
  } else if (props.trend.startsWith('-')) {
    return 'metric-trend--negative'
  }
  return 'metric-trend--neutral'
})

const trendIcon = computed(() => {
  if (!props.trend) return ''
  
  if (props.trend.startsWith('+')) {
    return 'trending_up'
  } else if (props.trend.startsWith('-')) {
    return 'trending_down'
  }
  return 'trending_flat'
})

// Methods
const handleClick = () => {
  if (props.clickable) {
    emit('click')
  }
}
</script>

<style lang="scss" scoped>
.metric-card {
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 1px solid var(--carrefour-gray-200);
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  
  &--clickable {
    cursor: pointer;
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
      border-color: var(--carrefour-blue);
    }
    
    &:active {
      transform: translateY(0);
    }
  }
}

.metric-card-content {
  padding: 20px;
  position: relative;
}

.metric-header {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 16px;
}

.metric-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.metric-value-section {
  flex: 1;
  min-width: 0;
}

.metric-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--carrefour-dark-blue);
  line-height: 1.2;
  margin-bottom: 4px;
}

.metric-trend {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 500;
  
  &--positive {
    color: var(--carrefour-green);
  }
  
  &--negative {
    color: var(--carrefour-red);
  }
  
  &--neutral {
    color: var(--carrefour-gray-600);
  }
}

.metric-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--carrefour-gray-800);
  margin-bottom: 4px;
  line-height: 1.3;
}

.metric-description {
  font-size: 14px;
  color: var(--carrefour-gray-600);
  line-height: 1.4;
}

// Loading state
:deep(.q-inner-loading) {
  background-color: rgba(255, 255, 255, 0.9);
  border-radius: 12px;
}

// Responsive design
@media (max-width: 768px) {
  .metric-card-content {
    padding: 16px;
  }
  
  .metric-header {
    gap: 12px;
    margin-bottom: 12px;
  }
  
  .metric-icon {
    width: 40px;
    height: 40px;
    
    .q-icon {
      font-size: 20px !important;
    }
  }
  
  .metric-value {
    font-size: 24px;
  }
  
  .metric-title {
    font-size: 14px;
  }
  
  .metric-description {
    font-size: 13px;
  }
}

@media (max-width: 480px) {
  .metric-header {
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: 8px;
  }
  
  .metric-value-section {
    text-align: center;
  }
  
  .metric-value {
    font-size: 20px;
  }
}

// Animation for value changes
.metric-value {
  transition: all 0.3s ease;
}

// Hover effects for clickable cards
.metric-card--clickable:hover {
  .metric-icon {
    transform: scale(1.05);
  }
  
  .metric-value {
    color: var(--carrefour-blue);
  }
}

// Dark mode support
@media (prefers-color-scheme: dark) {
  .metric-card {
    background-color: var(--carrefour-gray-800);
    border-color: var(--carrefour-gray-700);
  }
  
  .metric-value {
    color: white;
  }
  
  .metric-title {
    color: var(--carrefour-gray-200);
  }
  
  .metric-description {
    color: var(--carrefour-gray-400);
  }
}
</style>
