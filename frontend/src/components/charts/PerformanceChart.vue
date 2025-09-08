<template>
  <div class="performance-chart">
    <div ref="chartContainer" class="chart-container"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { Chart, registerables } from 'chart.js'

// Register Chart.js components
Chart.register(...registerables)

interface Props {
  data: {
    labels: string[]
    datasets: Array<{
      label: string
      data: number[]
      borderColor: string
      backgroundColor: string
      tension: number
    }>
  }
  period: string
}

const props = defineProps<Props>()

// Reactive state
const chartContainer = ref<HTMLCanvasElement>()
let chart: Chart | null = null

// Chart configuration
const chartConfig = {
  type: 'line' as const,
  data: props.data,
  options: {
    responsive: true,
    maintainAspectRatio: false,
    interaction: {
      intersect: false,
      mode: 'index' as const
    },
    plugins: {
      legend: {
        display: true,
        position: 'top' as const,
        labels: {
          usePointStyle: true,
          padding: 20,
          font: {
            family: 'Inter, sans-serif',
            size: 12
          }
        }
      },
      tooltip: {
        backgroundColor: 'rgba(0, 0, 0, 0.8)',
        titleColor: 'white',
        bodyColor: 'white',
        borderColor: 'var(--carrefour-blue)',
        borderWidth: 1,
        cornerRadius: 8,
        displayColors: true,
        callbacks: {
          title: (context: any) => {
            return `Day: ${context[0].label}`
          },
          label: (context: any) => {
            return `${context.dataset.label}: ${context.parsed.y.toLocaleString()}`
          }
        }
      }
    },
    scales: {
      x: {
        display: true,
        grid: {
          display: false
        },
        ticks: {
          color: 'var(--carrefour-gray-600)',
          font: {
            family: 'Inter, sans-serif',
            size: 11
          }
        }
      },
      y: {
        display: true,
        grid: {
          color: 'var(--carrefour-gray-200)',
          drawBorder: false
        },
        ticks: {
          color: 'var(--carrefour-gray-600)',
          font: {
            family: 'Inter, sans-serif',
            size: 11
          },
          callback: (value: any) => {
            return value.toLocaleString()
          }
        }
      }
    },
    elements: {
      point: {
        radius: 4,
        hoverRadius: 6,
        backgroundColor: 'white',
        borderWidth: 2
      },
      line: {
        borderWidth: 3
      }
    },
    animation: {
      duration: 1000,
      easing: 'easeInOutQuart' as const
    }
  }
}

// Methods
const createChart = () => {
  if (!chartContainer.value) return
  
  // Destroy existing chart
  if (chart) {
    chart.destroy()
  }
  
  // Create new chart
  chart = new Chart(chartContainer.value, chartConfig)
}

const updateChart = () => {
  if (!chart) return
  
  chart.data = props.data
  chart.update('active')
}

// Lifecycle
onMounted(() => {
  nextTick(() => {
    createChart()
  })
})

onUnmounted(() => {
  if (chart) {
    chart.destroy()
  }
})

// Watch for data changes
watch(() => props.data, () => {
  updateChart()
}, { deep: true })

// Watch for period changes
watch(() => props.period, () => {
  // In a real app, you would fetch new data based on the period
  updateChart()
})
</script>

<style lang="scss" scoped>
.performance-chart {
  width: 100%;
  height: 100%;
  position: relative;
}

.chart-container {
  width: 100%;
  height: 300px;
  position: relative;
}

// Responsive design
@media (max-width: 768px) {
  .chart-container {
    height: 250px;
  }
}

@media (max-width: 480px) {
  .chart-container {
    height: 200px;
  }
}
</style>
