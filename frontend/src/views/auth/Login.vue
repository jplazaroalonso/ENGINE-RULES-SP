<template>
  <div class="login-page">
    <div class="login-container">
      <!-- Logo and branding -->
      <div class="login-header">
        <img 
          src="/favicon.svg" 
          alt="Carrefour" 
          class="login-logo"
        />
        <h1 class="login-title">Rules Engine</h1>
        <p class="login-subtitle">Carrefour Business Rules Management</p>
      </div>

      <!-- Login form -->
      <q-card class="login-card">
        <q-card-section class="login-card-header">
          <h2 class="login-form-title">Sign In</h2>
          <p class="login-form-subtitle">Enter your credentials to access the system</p>
        </q-card-section>

        <q-card-section class="login-card-body">
          <q-form @submit="handleLogin" class="login-form">
            <!-- Email field -->
            <div class="form-field">
              <q-input
                v-model="form.email"
                type="email"
                label="Email Address"
                placeholder="Enter your email"
                outlined
                dense
                :rules="emailRules"
                :error="!!errors.email"
                :error-message="errors.email"
                @blur="clearError('email')"
              >
                <template v-slot:prepend>
                  <q-icon name="email" color="grey-6" />
                </template>
              </q-input>
            </div>

            <!-- Password field -->
            <div class="form-field">
              <q-input
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                label="Password"
                placeholder="Enter your password"
                outlined
                dense
                :rules="passwordRules"
                :error="!!errors.password"
                :error-message="errors.password"
                @blur="clearError('password')"
              >
                <template v-slot:prepend>
                  <q-icon name="lock" color="grey-6" />
                </template>
                <template v-slot:append>
                  <q-icon
                    :name="showPassword ? 'visibility_off' : 'visibility'"
                    class="cursor-pointer"
                    @click="showPassword = !showPassword"
                  />
                </template>
              </q-input>
            </div>

            <!-- Remember me and forgot password -->
            <div class="form-options">
              <q-checkbox
                v-model="form.rememberMe"
                label="Remember me"
                color="primary"
              />
              <q-btn
                flat
                dense
                label="Forgot password?"
                color="primary"
                size="sm"
                @click="handleForgotPassword"
              />
            </div>

            <!-- Submit button -->
            <div class="form-actions">
              <q-btn
                type="submit"
                color="primary"
                size="lg"
                class="login-button"
                :loading="loading"
                :disable="!isFormValid"
              >
                <span v-if="!loading">Sign In</span>
                <span v-else>Signing In...</span>
              </q-btn>
            </div>

            <!-- Demo credentials -->
            <div class="demo-credentials">
              <q-separator />
              <div class="demo-title">Demo Credentials</div>
              <div class="demo-accounts">
                <div class="demo-account" @click="fillDemoCredentials('admin')">
                  <q-icon name="admin_panel_settings" color="primary" />
                  <div>
                    <div class="demo-role">Administrator</div>
                    <div class="demo-email">admin@carrefour.com</div>
                  </div>
                </div>
                <div class="demo-account" @click="fillDemoCredentials('manager')">
                  <q-icon name="manage_accounts" color="secondary" />
                  <div>
                    <div class="demo-role">Manager</div>
                    <div class="demo-email">manager@carrefour.com</div>
                  </div>
                </div>
              </div>
            </div>
          </q-form>
        </q-card-section>
      </q-card>

      <!-- Footer -->
      <div class="login-footer">
        <p class="footer-text">
          © {{ currentYear }} Carrefour. All rights reserved.
        </p>
        <div class="footer-links">
          <a href="#" class="footer-link">Privacy Policy</a>
          <a href="#" class="footer-link">Terms of Service</a>
          <a href="#" class="footer-link">Support</a>
        </div>
      </div>
    </div>

    <!-- Background decoration -->
    <div class="login-background">
      <div class="background-shape background-shape--1"></div>
      <div class="background-shape background-shape--2"></div>
      <div class="background-shape background-shape--3"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notifications'

const router = useRouter()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()

// Reactive state
const loading = ref(false)
const showPassword = ref(false)
const currentYear = ref(new Date().getFullYear())

const form = ref({
  email: '',
  password: '',
  rememberMe: false
})

const errors = ref<Record<string, string>>({})

// Validation rules
const emailRules = [
  (val: string) => !!val || 'Email is required',
  (val: string) => /.+@.+\..+/.test(val) || 'Please enter a valid email'
]

const passwordRules = [
  (val: string) => !!val || 'Password is required',
  (val: string) => val.length >= 6 || 'Password must be at least 6 characters'
]

// Computed properties
const isFormValid = computed(() => {
  return form.value.email && 
         form.value.password && 
         /.+@.+\..+/.test(form.value.email) &&
         form.value.password.length >= 6
})

// Methods
const handleLogin = async () => {
  if (!isFormValid.value) return
  
  loading.value = true
  errors.value = {}
  
  try {
    await authStore.login(form.value.email, form.value.password)
    
    notificationStore.showSuccess('Login successful', 'Welcome back!')
    
    // Redirect to dashboard
    await router.push('/dashboard')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Login failed'
    notificationStore.showError('Login failed', message)
    
    // Set form errors
    if (message.includes('email')) {
      errors.value.email = message
    } else if (message.includes('password')) {
      errors.value.password = message
    }
  } finally {
    loading.value = false
  }
}

const handleForgotPassword = () => {
  notificationStore.showInfo('Forgot Password', 'Please contact your system administrator to reset your password.')
}

const fillDemoCredentials = (type: 'admin' | 'manager') => {
  if (type === 'admin') {
    form.value.email = 'admin@carrefour.com'
    form.value.password = 'admin123'
  } else {
    form.value.email = 'manager@carrefour.com'
    form.value.password = 'manager123'
  }
  
  // Clear any existing errors
  errors.value = {}
}

const clearError = (field: string) => {
  if (errors.value[field]) {
    delete errors.value[field]
  }
}

// Lifecycle
onMounted(() => {
  // Check if user is already authenticated
  if (authStore.isAuthenticated) {
    router.push('/dashboard')
  }
})
</script>

<style lang="scss" scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--carrefour-blue) 0%, var(--carrefour-dark-blue) 100%);
  position: relative;
  overflow: hidden;
  padding: 20px;
}

.login-container {
  width: 100%;
  max-width: 400px;
  z-index: 2;
  position: relative;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-logo {
  width: 64px;
  height: 64px;
  filter: brightness(0) invert(1); // Make logo white
  margin-bottom: 16px;
}

.login-title {
  font-size: 32px;
  font-weight: 700;
  color: white;
  margin: 0 0 8px 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.login-subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.9);
  margin: 0;
  font-weight: 400;
}

.login-card {
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.95);
}

.login-card-header {
  text-align: center;
  padding: 32px 32px 16px;
}

.login-form-title {
  font-size: 24px;
  font-weight: 600;
  color: var(--carrefour-dark-blue);
  margin: 0 0 8px 0;
}

.login-form-subtitle {
  font-size: 14px;
  color: var(--carrefour-gray-600);
  margin: 0;
}

.login-card-body {
  padding: 16px 32px 32px;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-field {
  width: 100%;
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 8px 0;
}

.form-actions {
  margin-top: 8px;
}

.login-button {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 8px;
  text-transform: none;
}

.demo-credentials {
  margin-top: 24px;
}

.demo-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--carrefour-gray-700);
  text-align: center;
  margin: 16px 0 12px;
}

.demo-accounts {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.demo-account {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border: 1px solid var(--carrefour-gray-200);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    border-color: var(--carrefour-blue);
    background-color: var(--carrefour-light-blue);
  }
}

.demo-role {
  font-size: 14px;
  font-weight: 500;
  color: var(--carrefour-dark-blue);
}

.demo-email {
  font-size: 12px;
  color: var(--carrefour-gray-600);
}

.login-footer {
  text-align: center;
  margin-top: 32px;
}

.footer-text {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.8);
  margin: 0 0 12px 0;
}

.footer-links {
  display: flex;
  justify-content: center;
  gap: 16px;
}

.footer-link {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.8);
  text-decoration: none;
  transition: color 0.2s ease;
  
  &:hover {
    color: white;
    text-decoration: underline;
  }
}

// Background decoration
.login-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1;
  overflow: hidden;
}

.background-shape {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  animation: float 6s ease-in-out infinite;
  
  &--1 {
    width: 200px;
    height: 200px;
    top: 10%;
    left: 10%;
    animation-delay: 0s;
  }
  
  &--2 {
    width: 150px;
    height: 150px;
    top: 60%;
    right: 15%;
    animation-delay: 2s;
  }
  
  &--3 {
    width: 100px;
    height: 100px;
    bottom: 20%;
    left: 20%;
    animation-delay: 4s;
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  50% {
    transform: translateY(-20px) rotate(180deg);
  }
}

// Responsive design
@media (max-width: 480px) {
  .login-page {
    padding: 16px;
  }
  
  .login-title {
    font-size: 28px;
  }
  
  .login-card-header,
  .login-card-body {
    padding: 24px 20px;
  }
  
  .form-options {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
  }
  
  .footer-links {
    flex-direction: column;
    gap: 8px;
  }
}

// Form field focus styles
:deep(.q-field--outlined .q-field__control:before) {
  border-color: var(--carrefour-gray-300);
}

:deep(.q-field--outlined.q-field--focused .q-field__control:before) {
  border-color: var(--carrefour-blue);
  border-width: 2px;
}

:deep(.q-field__label) {
  color: var(--carrefour-gray-600);
}

:deep(.q-field--focused .q-field__label) {
  color: var(--carrefour-blue);
}

// Button styles
:deep(.q-btn--unelevated) {
  box-shadow: 0 2px 8px rgba(0, 75, 135, 0.3);
}

:deep(.q-btn--unelevated:hover) {
  box-shadow: 0 4px 12px rgba(0, 75, 135, 0.4);
}
</style>
