<template>
  <q-toolbar class="app-footer-toolbar">
    <div class="app-footer-content">
      <!-- Left section -->
      <div class="app-footer-left">
        <div class="app-footer-brand">
          <img 
            src="/favicon.svg" 
            alt="Carrefour" 
            class="app-footer-logo"
          />
          <span class="app-footer-text">
            Rules Engine - Carrefour © {{ currentYear }}
          </span>
        </div>
      </div>

      <!-- Center section -->
      <div class="app-footer-center">
        <div class="app-footer-links">
          <a href="#" class="footer-link">Privacy Policy</a>
          <a href="#" class="footer-link">Terms of Service</a>
          <a href="#" class="footer-link">Support</a>
          <a href="#" class="footer-link">Documentation</a>
        </div>
      </div>

      <!-- Right section -->
      <div class="app-footer-right">
        <div class="app-footer-status">
          <q-icon 
            name="circle" 
            :color="systemStatus.color" 
            size="8px"
          />
          <span class="status-text">{{ systemStatus.text }}</span>
        </div>
        
        <div class="app-footer-version">
          v{{ appVersion }}
        </div>
      </div>
    </div>
  </q-toolbar>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

// Reactive state
const currentYear = ref(new Date().getFullYear())
const appVersion = ref('1.0.0')
const systemStatus = ref({
  color: 'positive',
  text: 'All systems operational'
})

// Computed properties
const statusInterval = ref<NodeJS.Timeout | null>(null)

// Methods
const checkSystemStatus = () => {
  // Mock system status check - replace with actual API call
  const statuses = [
    { color: 'positive', text: 'All systems operational' },
    { color: 'warning', text: 'Minor issues detected' },
    { color: 'negative', text: 'System maintenance in progress' }
  ]
  
  // Simulate status changes (in real app, this would be from API)
  const randomStatus = statuses[Math.floor(Math.random() * statuses.length)]
  systemStatus.value = randomStatus
}

// Lifecycle
onMounted(() => {
  // Check system status every 30 seconds
  statusInterval.value = setInterval(checkSystemStatus, 30000)
})

onUnmounted(() => {
  if (statusInterval.value) {
    clearInterval(statusInterval.value)
  }
})
</script>

<style lang="scss" scoped>
.app-footer-toolbar {
  background-color: var(--carrefour-gray-800);
  color: white;
  min-height: 48px;
  padding: 0 24px;
}

.app-footer-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  gap: 24px;
}

.app-footer-left {
  flex: 1;
  min-width: 0;
}

.app-footer-brand {
  display: flex;
  align-items: center;
  gap: 8px;
}

.app-footer-logo {
  width: 20px;
  height: 20px;
  filter: brightness(0) invert(1); // Make logo white
}

.app-footer-text {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.8);
  white-space: nowrap;
}

.app-footer-center {
  flex: 1;
  display: flex;
  justify-content: center;
}

.app-footer-links {
  display: flex;
  gap: 24px;
  align-items: center;
}

.footer-link {
  color: rgba(255, 255, 255, 0.7);
  text-decoration: none;
  font-size: 12px;
  transition: color 0.2s ease;
  
  &:hover {
    color: white;
    text-decoration: underline;
  }
}

.app-footer-right {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 16px;
}

.app-footer-status {
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-text {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.8);
}

.app-footer-version {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.6);
  font-family: 'Courier New', monospace;
  background-color: rgba(255, 255, 255, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
}

// Responsive design
@media (max-width: 1024px) {
  .app-footer-content {
    gap: 16px;
  }
  
  .app-footer-links {
    gap: 16px;
  }
}

@media (max-width: 768px) {
  .app-footer-toolbar {
    padding: 0 16px;
  }
  
  .app-footer-content {
    flex-direction: column;
    gap: 8px;
    text-align: center;
  }
  
  .app-footer-center,
  .app-footer-right {
    flex: none;
  }
  
  .app-footer-links {
    flex-wrap: wrap;
    justify-content: center;
    gap: 12px;
  }
  
  .app-footer-right {
    flex-direction: column;
    gap: 8px;
  }
}

@media (max-width: 480px) {
  .app-footer-links {
    display: none;
  }
  
  .app-footer-brand {
    flex-direction: column;
    gap: 4px;
  }
  
  .app-footer-text {
    font-size: 11px;
  }
}
</style>
