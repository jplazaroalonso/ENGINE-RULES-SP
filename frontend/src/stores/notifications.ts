import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import type { Notification, NotificationType } from '@/types'

export const useNotificationStore = defineStore('notifications', () => {
  // State
  const notifications = ref<Notification[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Getters
  const unreadNotifications = computed(() => 
    notifications.value.filter(n => !n.read)
  )
  
  const unreadCount = computed(() => unreadNotifications.value.length)
  
  const notificationsByType = computed(() => {
    return (type: NotificationType) => 
      notifications.value.filter(n => n.type === type)
  })

  // Actions
  const fetchNotifications = async () => {
    loading.value = true
    error.value = null

    try {
      // Mock API call - replace with actual implementation
      await new Promise(resolve => setTimeout(resolve, 500))
      
      // Mock notifications
      notifications.value = [
        {
          id: '1',
          title: 'New Rule Created',
          message: 'A new promotion rule has been created and is pending approval.',
          type: 'INFO',
          read: false,
          createdAt: new Date(Date.now() - 1000 * 60 * 5).toISOString(), // 5 minutes ago
          actionUrl: '/rules'
        },
        {
          id: '2',
          title: 'Rule Approved',
          message: 'Your "Senior Citizen Discount" rule has been approved and is now active.',
          type: 'SUCCESS',
          read: false,
          createdAt: new Date(Date.now() - 1000 * 60 * 30).toISOString(), // 30 minutes ago
          actionUrl: '/rules'
        },
        {
          id: '3',
          title: 'Campaign Performance Alert',
          message: 'Your "Black Friday" campaign has exceeded 80% of its budget.',
          type: 'WARNING',
          read: true,
          createdAt: new Date(Date.now() - 1000 * 60 * 60 * 2).toISOString(), // 2 hours ago
          actionUrl: '/campaigns'
        },
        {
          id: '4',
          title: 'System Maintenance',
          message: 'Scheduled maintenance will occur tonight from 2:00 AM to 4:00 AM.',
          type: 'INFO',
          read: true,
          createdAt: new Date(Date.now() - 1000 * 60 * 60 * 6).toISOString(), // 6 hours ago
        }
      ]
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch notifications'
    } finally {
      loading.value = false
    }
  }

  const markAsRead = async (notificationId: string) => {
    try {
      const notification = notifications.value.find(n => n.id === notificationId)
      if (notification) {
        notification.read = true
        // In a real app, you would make an API call here
      }
    } catch (err) {
      console.error('Failed to mark notification as read:', err)
    }
  }

  const markAllAsRead = async () => {
    try {
      notifications.value.forEach(n => n.read = true)
      // In a real app, you would make an API call here
    } catch (err) {
      console.error('Failed to mark all notifications as read:', err)
    }
  }

  const deleteNotification = async (notificationId: string) => {
    try {
      const index = notifications.value.findIndex(n => n.id === notificationId)
      if (index !== -1) {
        notifications.value.splice(index, 1)
        // In a real app, you would make an API call here
      }
    } catch (err) {
      console.error('Failed to delete notification:', err)
    }
  }

  const addNotification = (notification: Omit<Notification, 'id' | 'createdAt'>) => {
    const newNotification: Notification = {
      ...notification,
      id: Date.now().toString(),
      createdAt: new Date().toISOString()
    }
    
    notifications.value.unshift(newNotification)
    
    // Auto-remove after 10 seconds for success/error notifications
    if (['SUCCESS', 'ERROR'].includes(notification.type)) {
      setTimeout(() => {
        deleteNotification(newNotification.id)
      }, 10000)
    }
  }

  const showSuccess = (title: string, message?: string) => {
    addNotification({
      title,
      message: message || title,
      type: 'SUCCESS',
      read: false
    })
  }

  const showError = (title: string, message?: string) => {
    addNotification({
      title,
      message: message || title,
      type: 'ERROR',
      read: false
    })
  }

  const showWarning = (title: string, message?: string) => {
    addNotification({
      title,
      message: message || title,
      type: 'WARNING',
      read: false
    })
  }

  const showInfo = (title: string, message?: string) => {
    addNotification({
      title,
      message: message || title,
      type: 'INFO',
      read: false
    })
  }

  const clearError = () => {
    error.value = null
  }

  return {
    // State
    notifications: readonly(notifications),
    loading: readonly(loading),
    error: readonly(error),
    
    // Getters
    unreadNotifications,
    unreadCount,
    notificationsByType,
    
    // Actions
    fetchNotifications,
    markAsRead,
    markAllAsRead,
    deleteNotification,
    addNotification,
    showSuccess,
    showError,
    showWarning,
    showInfo,
    clearError
  }
})
