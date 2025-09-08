// Core entity types
export interface Rule {
  id: string
  name: string
  description: string
  dsl_content: string
  status: RuleStatus
  priority: Priority
  category: string
  version: number
  created_by: string
  created_at: string
  updated_at: string
  tags: string[]
}

export type RuleStatus = 
  | 'DRAFT' 
  | 'UNDER_REVIEW' 
  | 'APPROVED' 
  | 'ACTIVE' 
  | 'INACTIVE' 
  | 'DEPRECATED'

export type Priority = 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL'

// API types
export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: string
  message?: string
  pagination?: PaginationInfo
}

export interface PaginationInfo {
  page: number
  limit: number
  total: number
  totalPages: number
}

export interface ListRequest {
  page?: number
  limit?: number
  sort?: string
  order?: 'asc' | 'desc'
  filters?: Record<string, any>
  search?: string
}

export interface CreateRuleRequest {
  name: string
  description: string
  dsl_content: string
  priority: Priority
  category: string
  tags?: string[]
}

export interface UpdateRuleRequest {
  name?: string
  description?: string
  dsl_content?: string
  priority?: Priority
  category?: string
  tags?: string[]
}

// Evaluation types
export interface EvaluationRequest {
  rule_category: string
  dsl_content: string
  context: Record<string, any>
}

export interface EvaluationResponse {
  result: Record<string, any>
}

// Calculation types
export interface CalculationRequest {
  rule_ids: string[]
  context: Record<string, any>
}

export interface CalculationResponse {
  calculation_id: string
  value: number
  breakdown: Record<string, number>
}

// User types
export interface User {
  id: string
  name: string
  email: string
  role: UserRole
  avatar?: string
  lastLogin?: string
}

export type UserRole = 'ADMIN' | 'MANAGER' | 'USER' | 'VIEWER'

// Auth types
export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
  expiresAt: string
}

export interface UpdateUserRequest {
  name?: string
  email?: string
  role?: UserRole
  department?: string
  phone?: string
}

export interface ChangePasswordRequest {
  currentPassword: string
  newPassword: string
}

// Campaign types
export interface Campaign {
  id: string
  name: string
  description: string
  type: CampaignType
  status: CampaignStatus
  startDate: string
  endDate: string
  budget: number
  spent: number
  rules: string[]
  metrics: CampaignMetrics
}

export type CampaignType = 'PROMOTION' | 'LOYALTY' | 'COUPON' | 'TAX' | 'PAYMENT'
export type CampaignStatus = 'DRAFT' | 'SCHEDULED' | 'ACTIVE' | 'PAUSED' | 'COMPLETED' | 'CANCELLED'

export interface CampaignMetrics {
  impressions: number
  clicks: number
  conversions: number
  revenue: number
  roi: number
}

// Customer types
export interface Customer {
  id: string
  name: string
  email: string
  phone?: string
  tier: CustomerTier
  segment: string
  loyaltyPoints: number
  totalSpent: number
  lastPurchase: string
  registeredAt: string
  attributes: Record<string, any>
}

export type CustomerTier = 'BRONZE' | 'SILVER' | 'GOLD' | 'PLATINUM'

// Notification types
export interface Notification {
  id: string
  title: string
  message: string
  type: NotificationType
  read: boolean
  createdAt: string
  actionUrl?: string
}

export type NotificationType = 'INFO' | 'SUCCESS' | 'WARNING' | 'ERROR'

// Form types
export interface FormField {
  name: string
  label: string
  type: FieldType
  required?: boolean
  readonly?: boolean
  hidden?: boolean
  placeholder?: string
  helpText?: string
  options?: SelectOption[]
  validation?: FieldValidation
  dependsOn?: string[]
}

export type FieldType = 
  | 'text' 
  | 'email' 
  | 'password' 
  | 'number' 
  | 'textarea' 
  | 'select' 
  | 'multiselect' 
  | 'checkbox' 
  | 'radio' 
  | 'date' 
  | 'datetime' 
  | 'file' 
  | 'dsl-editor'

export interface SelectOption {
  label: string
  value: any
  disabled?: boolean
  group?: string
}

export interface FieldValidation {
  required?: boolean
  minLength?: number
  maxLength?: number
  min?: number
  max?: number
  pattern?: string
  custom?: (value: any) => string | null
}

// UI types
export interface TableColumn {
  name: string
  label: string
  field: string
  align?: 'left' | 'center' | 'right'
  sortable?: boolean
  format?: (value: any) => string
}

export interface FilterConfig {
  name: string
  label: string
  type: 'select' | 'multiselect' | 'date' | 'daterange' | 'text'
  options?: SelectOption[]
  placeholder?: string
}

export interface EntityAction {
  name: string
  label: string
  icon: string
  color?: string
  condition?: (item: any) => boolean
  handler: (item: any) => void
}

// Store types
export interface RulesState {
  rules: Rule[]
  currentRule: Rule | null
  loading: boolean
  error: string | null
  pagination: PaginationInfo
  filters: RuleFilters
  searchQuery: string
}

export interface RuleFilters {
  status: RuleStatus[]
  priority: Priority[]
  category: string
  createdBy: string
  dateRange: DateRange | null
  tags: string[]
}

export interface DateRange {
  start: Date
  end: Date
}

// Error types
export interface ApiError {
  message: string
  code?: string
  details?: Record<string, any>
}

// Theme types
export interface ThemeConfig {
  primary: string
  secondary: string
  accent: string
  positive: string
  negative: string
  info: string
  warning: string
}
