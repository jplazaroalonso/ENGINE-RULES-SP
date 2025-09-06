# Campaigns List Screen - Campaign Management Overview

## Overview
The Campaigns List screen provides a unified interface for viewing and managing all types of marketing campaigns including promotions, loyalty programs, and coupon campaigns. It serves as the central hub for campaign operations.

## Screen Layout

### Header Section with Campaign Type Toggle
```vue
<template>
  <div class="campaigns-header">
    <div class="page-title-section">
      <h1>Campaign Management</h1>
      <p class="subtitle">Manage promotions, loyalty programs, and coupon campaigns</p>
    </div>
    
    <div class="campaign-type-toggle">
      <q-btn-toggle
        v-model="activeCampaignType"
        :options="campaignTypeOptions"
        color="primary"
        outline
        @update:model-value="onCampaignTypeChange"
      />
    </div>
    
    <div class="header-actions">
      <q-btn-dropdown
        color="primary"
        icon="add"
        label="Create Campaign"
        auto-close
      >
        <q-list>
          <q-item clickable @click="createCampaign('promotion')">
            <q-item-section avatar>
              <q-icon name="local_offer" color="orange" />
            </q-item-section>
            <q-item-section>
              <q-item-label>Promotion Campaign</q-item-label>
              <q-item-label caption>Discounts and special offers</q-item-label>
            </q-item-section>
          </q-item>
          
          <q-item clickable @click="createCampaign('loyalty')">
            <q-item-section avatar>
              <q-icon name="stars" color="purple" />
            </q-item-section>
            <q-item-section>
              <q-item-label>Loyalty Program</q-item-label>
              <q-item-label caption>Customer loyalty and rewards</q-item-label>
            </q-item-section>
          </q-item>
          
          <q-item clickable @click="createCampaign('coupon')">
            <q-item-section avatar>
              <q-icon name="confirmation_number" color="green" />
            </q-item-section>
            <q-item-section>
              <q-item-label>Coupon Campaign</q-item-label>
              <q-item-label caption>Coupon codes and vouchers</q-item-label>
            </q-item-section>
          </q-item>
        </q-list>
      </q-btn-dropdown>
      
      <q-btn
        flat
        icon="download"
        label="Export"
        @click="showExportDialog = true"
      />
      
      <q-btn
        flat
        icon="analytics"
        label="Analytics"
        @click="navigateToAnalytics"
      />
    </div>
  </div>
</template>
```

### Campaign Filters and Search
```vue
<template>
  <div class="campaigns-filters">
    <q-expansion-item 
      icon="filter_list" 
      label="Filters" 
      v-model="filtersExpanded"
      class="filters-expansion"
    >
      <q-card class="filters-card">
        <q-card-section>
          <div class="row q-gutter-md">
            <div class="col-md-2 col-sm-6 col-12">
              <q-select
                v-model="filters.status"
                label="Status"
                :options="statusOptions"
                clearable
                multiple
                use-chips
                outlined
                dense
              />
            </div>
            
            <div class="col-md-2 col-sm-6 col-12">
              <q-select
                v-model="filters.targetAudience"
                label="Target Audience"
                :options="audienceOptions"
                clearable
                multiple
                use-chips
                outlined
                dense
              />
            </div>
            
            <div class="col-md-2 col-sm-6 col-12">
              <q-select
                v-model="filters.discountType"
                label="Discount Type"
                :options="discountTypeOptions"
                clearable
                multiple
                use-chips
                outlined
                dense
                v-if="activeCampaignType === 'promotion' || activeCampaignType === 'coupon'"
              />
            </div>
            
            <div class="col-md-2 col-sm-6 col-12">
              <q-input
                v-model="filters.dateRange"
                label="Campaign Period"
                mask="##/##/#### - ##/##/####"
                outlined
                dense
              >
                <template v-slot:append>
                  <q-icon name="event" class="cursor-pointer">
                    <q-popup-proxy>
                      <q-date 
                        v-model="filters.dateRange" 
                        range 
                        @update:model-value="applyFilters"
                      />
                    </q-popup-proxy>
                  </q-icon>
                </template>
              </q-input>
            </div>
            
            <div class="col-md-2 col-sm-6 col-12">
              <q-input
                v-model.number="filters.budgetRange.min"
                label="Min Budget"
                type="number"
                outlined
                dense
                prefix="$"
              />
            </div>
            
            <div class="col-md-2 col-sm-6 col-12">
              <div class="filter-actions">
                <q-btn 
                  color="primary" 
                  label="Apply" 
                  @click="applyFilters"
                  dense
                />
                <q-btn 
                  flat 
                  label="Clear" 
                  @click="clearFilters"
                  dense
                />
              </div>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </q-expansion-item>
  </div>

  <!-- Search and Quick Stats -->
  <div class="search-section">
    <div class="row q-gutter-md items-center">
      <div class="col-md-4 col-sm-6 col-12">
        <q-input
          v-model="searchQuery"
          label="Search campaigns..."
          outlined
          dense
          clearable
          debounce="300"
          @update:model-value="onSearch"
        >
          <template v-slot:prepend>
            <q-icon name="search" />
          </template>
        </q-input>
      </div>
      
      <div class="col-auto">
        <div class="campaign-stats">
          <q-chip 
            color="positive" 
            text-color="white" 
            icon="trending_up"
            class="q-mr-sm"
          >
            Active: {{ campaignStats.active }}
          </q-chip>
          
          <q-chip 
            color="info" 
            text-color="white" 
            icon="schedule"
            class="q-mr-sm"
          >
            Scheduled: {{ campaignStats.scheduled }}
          </q-chip>
          
          <q-chip 
            color="warning" 
            text-color="white" 
            icon="pause"
            class="q-mr-sm"
          >
            Paused: {{ campaignStats.paused }}
          </q-chip>
          
          <q-chip 
            color="grey" 
            text-color="white" 
            icon="history"
          >
            Completed: {{ campaignStats.completed }}
          </q-chip>
        </div>
      </div>
      
      <div class="col-auto">
        <q-btn-dropdown
          color="grey-7"
          icon="sort"
          label="Sort"
          outline
          dense
        >
          <q-list>
            <q-item
              v-for="sortOption in sortOptions"
              :key="sortOption.value"
              clickable
              @click="applySorting(sortOption.value)"
              :class="{ 'bg-grey-2': currentSort === sortOption.value }"
            >
              <q-item-section>
                <q-item-label>{{ sortOption.label }}</q-item-label>
              </q-item-section>
              <q-item-section side v-if="currentSort === sortOption.value">
                <q-icon name="check" color="primary" />
              </q-item-section>
            </q-item>
          </q-list>
        </q-btn-dropdown>
      </div>
    </div>
  </div>
</template>
```

### Campaign Cards Grid
```vue
<template>
  <div class="campaigns-grid">
    <div class="row q-gutter-md">
      <div
        v-for="campaign in filteredCampaigns"
        :key="campaign.id"
        class="col-xl-4 col-lg-6 col-md-6 col-sm-12 col-12"
      >
        <CampaignCard
          :campaign="campaign"
          :type="activeCampaignType"
          :show-metrics="true"
          @view="viewCampaign"
          @edit="editCampaign"
          @duplicate="duplicateCampaign"
          @toggle-status="toggleCampaignStatus"
          @delete="deleteCampaign"
          @view-performance="viewCampaignPerformance"
        />
      </div>
    </div>
    
    <!-- Empty State -->
    <div v-if="filteredCampaigns.length === 0" class="empty-state">
      <q-icon 
        :name="getCampaignTypeIcon(activeCampaignType)" 
        size="4rem" 
        color="grey-4" 
      />
      <div class="empty-title">No {{ activeCampaignType }} campaigns found</div>
      <div class="empty-subtitle">
        {{ searchQuery ? 'Try adjusting your search criteria' : `Create your first ${activeCampaignType} campaign to get started` }}
      </div>
      <q-btn
        color="primary"
        :icon="getCampaignTypeIcon(activeCampaignType)"
        :label="`Create ${activeCampaignType.charAt(0).toUpperCase() + activeCampaignType.slice(1)} Campaign`"
        @click="createCampaign(activeCampaignType)"
        class="q-mt-md"
        v-if="!searchQuery"
      />
    </div>
    
    <!-- Load More / Pagination -->
    <div class="campaigns-pagination" v-if="pagination.hasMore">
      <q-btn
        color="primary"
        icon="keyboard_arrow_down"
        label="Load More Campaigns"
        @click="loadMore"
        :loading="loadingMore"
        flat
        class="full-width"
      />
    </div>
  </div>
</template>
```

## Campaign Card Component

### Promotion Campaign Card
```vue
<template>
  <q-card class="campaign-card promotion-card" :class="cardClasses">
    <q-card-section>
      <!-- Campaign Header -->
      <div class="campaign-header">
        <div class="campaign-title">
          <div class="campaign-type-icon">
            <q-icon name="local_offer" color="orange" size="sm" />
          </div>
          <div class="campaign-info">
            <h4>{{ campaign.name }}</h4>
            <div class="campaign-meta">
              <q-chip
                :color="getStatusColor(campaign.status)"
                text-color="white"
                size="sm"
                :icon="getStatusIcon(campaign.status)"
              >
                {{ campaign.status }}
              </q-chip>
              
              <q-chip
                color="orange"
                text-color="white"
                size="sm"
                icon="percent"
              >
                {{ formatDiscount(campaign.discountType, campaign.discountValue) }}
              </q-chip>
            </div>
          </div>
        </div>
        
        <div class="campaign-actions">
          <q-btn
            flat
            dense
            icon="visibility"
            @click="$emit('view', campaign)"
            size="sm"
          >
            <q-tooltip>View Details</q-tooltip>
          </q-btn>
          
          <q-btn
            flat
            dense
            icon="edit"
            @click="$emit('edit', campaign)"
            size="sm"
            v-if="canEdit"
          >
            <q-tooltip>Edit Campaign</q-tooltip>
          </q-btn>
          
          <q-btn
            flat
            dense
            icon="more_vert"
            size="sm"
          >
            <q-menu>
              <q-list>
                <q-item clickable @click="$emit('duplicate', campaign)">
                  <q-item-section avatar>
                    <q-icon name="content_copy" />
                  </q-item-section>
                  <q-item-section>Duplicate</q-item-section>
                </q-item>
                
                <q-item clickable @click="$emit('view-performance', campaign)">
                  <q-item-section avatar>
                    <q-icon name="analytics" />
                  </q-item-section>
                  <q-item-section>View Performance</q-item-section>
                </q-item>
                
                <q-separator />
                
                <q-item 
                  clickable 
                  @click="$emit('toggle-status', campaign)"
                  :class="campaign.status === 'ACTIVE' ? 'text-warning' : 'text-positive'"
                >
                  <q-item-section avatar>
                    <q-icon :name="campaign.status === 'ACTIVE' ? 'pause' : 'play_arrow'" />
                  </q-item-section>
                  <q-item-section>
                    {{ campaign.status === 'ACTIVE' ? 'Pause' : 'Activate' }}
                  </q-item-section>
                </q-item>
                
                <q-separator />
                
                <q-item 
                  clickable 
                  @click="$emit('delete', campaign)"
                  v-if="canDelete"
                  class="text-negative"
                >
                  <q-item-section avatar>
                    <q-icon name="delete" />
                  </q-item-section>
                  <q-item-section>Delete</q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </div>
      </div>
      
      <!-- Campaign Description -->
      <div class="campaign-description" v-if="campaign.description">
        <p>{{ truncateText(campaign.description, 100) }}</p>
      </div>
      
      <!-- Campaign Period -->
      <div class="campaign-period">
        <q-icon name="event" color="grey-6" size="sm" class="q-mr-xs" />
        <span class="period-text">
          {{ formatDateRange(campaign.startDate, campaign.endDate) }}
        </span>
        <q-chip
          :color="getPeriodStatusColor(campaign)"
          text-color="white"
          size="sm"
          class="q-ml-sm"
        >
          {{ getPeriodStatus(campaign) }}
        </q-chip>
      </div>
      
      <!-- Target Audience -->
      <div class="campaign-audience" v-if="campaign.targetSegments?.length">
        <q-icon name="group" color="grey-6" size="sm" class="q-mr-xs" />
        <span class="audience-text">
          {{ campaign.targetSegments.slice(0, 2).join(', ') }}
          <span v-if="campaign.targetSegments.length > 2">
            +{{ campaign.targetSegments.length - 2 }} more
          </span>
        </span>
      </div>
      
      <!-- Campaign Metrics -->
      <div class="campaign-metrics" v-if="showMetrics && metrics">
        <div class="metrics-grid">
          <div class="metric-item">
            <div class="metric-value">{{ formatNumber(metrics.impressions) }}</div>
            <div class="metric-label">Impressions</div>
          </div>
          
          <div class="metric-item">
            <div class="metric-value">{{ formatNumber(metrics.conversions) }}</div>
            <div class="metric-label">Conversions</div>
          </div>
          
          <div class="metric-item">
            <div class="metric-value">{{ formatPercentage(metrics.conversionRate) }}</div>
            <div class="metric-label">Conv. Rate</div>
          </div>
          
          <div class="metric-item">
            <div class="metric-value">{{ formatCurrency(metrics.revenue) }}</div>
            <div class="metric-label">Revenue</div>
          </div>
        </div>
        
        <!-- Performance Indicator -->
        <div class="performance-indicator">
          <q-linear-progress
            :value="metrics.budgetUsed / campaign.budgetLimit"
            size="4px"
            :color="getBudgetUsageColor(metrics.budgetUsed / campaign.budgetLimit)"
            class="q-mt-sm"
          />
          <div class="budget-info">
            <span class="budget-used">{{ formatCurrency(metrics.budgetUsed) }}</span>
            <span class="budget-separator">/</span>
            <span class="budget-total">{{ formatCurrency(campaign.budgetLimit) }}</span>
          </div>
        </div>
      </div>
    </q-card-section>
    
    <!-- Campaign Status Footer -->
    <q-card-section class="campaign-footer" v-if="campaign.status === 'ACTIVE'">
      <div class="status-indicators">
        <div class="status-item">
          <q-icon name="schedule" color="positive" size="sm" />
          <span>Live Campaign</span>
        </div>
        
        <div class="status-item" v-if="getDaysRemaining(campaign.endDate) > 0">
          <q-icon name="access_time" color="warning" size="sm" />
          <span>{{ getDaysRemaining(campaign.endDate) }} days left</span>
        </div>
        
        <div class="status-item" v-if="isNearBudgetLimit(metrics?.budgetUsed, campaign.budgetLimit)">
          <q-icon name="warning" color="negative" size="sm" />
          <span>Near budget limit</span>
        </div>
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { PromotionCampaign, CampaignMetrics } from '@/types'

interface Props {
  campaign: PromotionCampaign
  metrics?: CampaignMetrics
  showMetrics?: boolean
  canEdit?: boolean
  canDelete?: boolean
}

interface Emits {
  (e: 'view', campaign: PromotionCampaign): void
  (e: 'edit', campaign: PromotionCampaign): void
  (e: 'duplicate', campaign: PromotionCampaign): void
  (e: 'toggle-status', campaign: PromotionCampaign): void
  (e: 'delete', campaign: PromotionCampaign): void
  (e: 'view-performance', campaign: PromotionCampaign): void
}

const props = withDefaults(defineProps<Props>(), {
  showMetrics: true,
  canEdit: true,
  canDelete: true
})

const emit = defineEmits<Emits>()

const cardClasses = computed(() => ({
  'campaign-card--active': props.campaign.status === 'ACTIVE',
  'campaign-card--paused': props.campaign.status === 'PAUSED',
  'campaign-card--scheduled': props.campaign.status === 'SCHEDULED',
  'campaign-card--expired': props.campaign.status === 'EXPIRED',
  'campaign-card--near-budget': isNearBudgetLimit(props.metrics?.budgetUsed, props.campaign.budgetLimit)
}))

// Utility functions
const formatDiscount = (type: string, value: number) => {
  return type === 'PERCENTAGE' ? `${value}% OFF` : `$${value} OFF`
}

const formatDateRange = (start: string, end: string) => {
  const startDate = new Date(start).toLocaleDateString()
  const endDate = new Date(end).toLocaleDateString()
  return `${startDate} - ${endDate}`
}

const getPeriodStatus = (campaign: PromotionCampaign) => {
  const now = new Date()
  const start = new Date(campaign.startDate)
  const end = new Date(campaign.endDate)
  
  if (now < start) return 'Upcoming'
  if (now > end) return 'Expired'
  return 'Active'
}

const getPeriodStatusColor = (campaign: PromotionCampaign) => {
  const status = getPeriodStatus(campaign)
  const colors: Record<string, string> = {
    'Upcoming': 'info',
    'Active': 'positive',
    'Expired': 'negative'
  }
  return colors[status] || 'grey'
}

const getDaysRemaining = (endDate: string) => {
  const now = new Date()
  const end = new Date(endDate)
  const diffTime = end.getTime() - now.getTime()
  return Math.ceil(diffTime / (1000 * 60 * 60 * 24))
}

const isNearBudgetLimit = (used?: number, total?: number) => {
  if (!used || !total) return false
  return (used / total) >= 0.8
}

const getBudgetUsageColor = (ratio: number) => {
  if (ratio >= 0.8) return 'negative'
  if (ratio >= 0.6) return 'warning'
  return 'positive'
}

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    ACTIVE: 'positive',
    PAUSED: 'warning',
    SCHEDULED: 'info',
    EXPIRED: 'negative',
    DRAFT: 'grey'
  }
  return colors[status] || 'grey'
}

const getStatusIcon = (status: string) => {
  const icons: Record<string, string> = {
    ACTIVE: 'play_arrow',
    PAUSED: 'pause',
    SCHEDULED: 'schedule',
    EXPIRED: 'history',
    DRAFT: 'edit'
  }
  return icons[status] || 'help'
}

const truncateText = (text: string, maxLength: number) => {
  if (text.length <= maxLength) return text
  return text.substring(0, maxLength) + '...'
}

const formatNumber = (value: number) => {
  return new Intl.NumberFormat().format(value)
}

const formatPercentage = (value: number) => {
  return `${value.toFixed(1)}%`
}

const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD'
  }).format(value)
}
</script>

<style scoped lang="scss">
.campaign-card {
  transition: all 0.3s ease;
  border-radius: 12px;
  
  &:hover {
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
    transform: translateY(-4px);
  }
  
  &--active {
    border-top: 4px solid $positive;
  }
  
  &--paused {
    border-top: 4px solid $warning;
  }
  
  &--scheduled {
    border-top: 4px solid $info;
  }
  
  &--expired {
    border-top: 4px solid $negative;
    opacity: 0.8;
  }
  
  &--near-budget {
    background: linear-gradient(135deg, #fff 0%, #fff9f0 100%);
  }
  
  .campaign-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 16px;
    
    .campaign-title {
      display: flex;
      gap: 12px;
      flex: 1;
      
      .campaign-type-icon {
        background: rgba(255, 152, 0, 0.1);
        border-radius: 8px;
        padding: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
      }
      
      .campaign-info {
        flex: 1;
        
        h4 {
          margin: 0 0 8px 0;
          font-size: 1.1rem;
          font-weight: 600;
          line-height: 1.3;
        }
        
        .campaign-meta {
          display: flex;
          gap: 8px;
          flex-wrap: wrap;
        }
      }
    }
    
    .campaign-actions {
      display: flex;
      gap: 4px;
    }
  }
  
  .campaign-description {
    margin-bottom: 16px;
    
    p {
      margin: 0;
      color: $grey-7;
      line-height: 1.4;
    }
  }
  
  .campaign-period,
  .campaign-audience {
    display: flex;
    align-items: center;
    margin-bottom: 12px;
    font-size: 0.9rem;
    color: $grey-7;
  }
  
  .campaign-metrics {
    margin-top: 16px;
    padding-top: 16px;
    border-top: 1px solid $grey-3;
    
    .metrics-grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 16px;
      margin-bottom: 12px;
      
      .metric-item {
        text-align: center;
        
        .metric-value {
          font-size: 1.1rem;
          font-weight: 600;
          color: $grey-9;
        }
        
        .metric-label {
          font-size: 0.8rem;
          color: $grey-6;
          margin-top: 2px;
        }
      }
    }
    
    .performance-indicator {
      .budget-info {
        display: flex;
        justify-content: space-between;
        align-items: center;
        font-size: 0.85rem;
        margin-top: 4px;
        
        .budget-used {
          font-weight: 600;
          color: $grey-8;
        }
        
        .budget-separator {
          color: $grey-5;
        }
        
        .budget-total {
          color: $grey-6;
        }
      }
    }
  }
  
  .campaign-footer {
    background: rgba(0, 0, 0, 0.02);
    border-top: 1px solid $grey-3;
    
    .status-indicators {
      display: flex;
      gap: 16px;
      flex-wrap: wrap;
      
      .status-item {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 0.85rem;
        color: $grey-7;
      }
    }
  }
}
</style>
```

## Component Data Structure

### Campaigns List Data Interface
```typescript
interface CampaignsListData {
  campaigns: Campaign[]
  filteredCampaigns: Campaign[]
  activeCampaignType: 'promotion' | 'loyalty' | 'coupon'
  loading: boolean
  loadingMore: boolean
  searchQuery: string
  filtersExpanded: boolean
  currentSort: string
  
  campaignStats: {
    active: number
    scheduled: number
    paused: number
    completed: number
  }
  
  filters: {
    status: CampaignStatus[]
    targetAudience: string[]
    discountType: DiscountType[]
    dateRange: string
    budgetRange: {
      min: number
      max: number
    }
  }
  
  pagination: {
    page: number
    limit: number
    total: number
    hasMore: boolean
  }
}

interface Campaign {
  id: string
  name: string
  description: string
  type: 'promotion' | 'loyalty' | 'coupon'
  status: CampaignStatus
  startDate: string
  endDate: string
  budgetLimit?: number
  targetSegments: string[]
  createdAt: string
  createdBy: string
  
  // Type-specific properties
  discountType?: DiscountType
  discountValue?: number
  tierStructure?: LoyaltyTier[]
  couponCodes?: CouponCode[]
}

interface CampaignMetrics {
  impressions: number
  conversions: number
  conversionRate: number
  revenue: number
  budgetUsed: number
  roi: number
  engagement: number
}

type CampaignStatus = 'DRAFT' | 'SCHEDULED' | 'ACTIVE' | 'PAUSED' | 'EXPIRED' | 'COMPLETED'
type DiscountType = 'PERCENTAGE' | 'FIXED_AMOUNT' | 'BUY_X_GET_Y' | 'FREE_SHIPPING'
```

## Screen Interactions

### Campaign Management Actions
- **Create Campaign**: Navigate to creation form based on type
- **View Campaign**: Navigate to campaign detail page
- **Edit Campaign**: Navigate to edit form
- **Duplicate Campaign**: Create copy with modified name
- **Toggle Status**: Activate/pause campaigns
- **Delete Campaign**: Soft delete with confirmation

### Filtering and Search
- **Campaign Type Toggle**: Switch between promotion/loyalty/coupon views
- **Status Filtering**: Filter by campaign status
- **Date Range**: Filter by campaign period
- **Audience Targeting**: Filter by target segments
- **Budget Range**: Filter by budget limits

### Performance Actions
- **View Analytics**: Navigate to campaign performance dashboard
- **Export Data**: Export campaign data in various formats
- **Performance Comparison**: Compare multiple campaign metrics

## State Management
```typescript
// Campaigns Store
export const useCampaignsStore = defineStore('campaigns', () => {
  const campaigns = ref<Campaign[]>([])
  const activeCampaignType = ref<CampaignType>('promotion')
  const loading = ref(false)
  
  const filteredCampaigns = computed(() => {
    return campaigns.value.filter(campaign => 
      campaign.type === activeCampaignType.value
    )
  })
  
  const campaignStats = computed(() => {
    const typeCampaigns = filteredCampaigns.value
    return {
      active: typeCampaigns.filter(c => c.status === 'ACTIVE').length,
      scheduled: typeCampaigns.filter(c => c.status === 'SCHEDULED').length,
      paused: typeCampaigns.filter(c => c.status === 'PAUSED').length,
      completed: typeCampaigns.filter(c => c.status === 'COMPLETED').length
    }
  })
  
  const fetchCampaigns = async (type: CampaignType) => {
    loading.value = true
    try {
      const response = await campaignsApi.list({ type })
      campaigns.value = response.data
    } finally {
      loading.value = false
    }
  }
  
  const toggleCampaignStatus = async (campaignId: string) => {
    const campaign = campaigns.value.find(c => c.id === campaignId)
    if (!campaign) return
    
    const newStatus = campaign.status === 'ACTIVE' ? 'PAUSED' : 'ACTIVE'
    await campaignsApi.updateStatus(campaignId, newStatus)
    
    campaign.status = newStatus
  }
  
  return {
    campaigns: readonly(campaigns),
    filteredCampaigns,
    activeCampaignType,
    campaignStats,
    loading: readonly(loading),
    fetchCampaigns,
    toggleCampaignStatus
  }
})
```

## Performance Optimizations
- Virtual scrolling for large campaign lists
- Lazy loading of campaign metrics
- Debounced search and filtering
- Optimistic updates for status changes
- Image lazy loading for campaign assets

## Accessibility Features
- Keyboard navigation for campaign cards
- Screen reader support for metrics
- High contrast mode for status indicators
- Focus management for modals
- ARIA labels for interactive elements

## Testing Strategy
- Unit tests for campaign filtering logic
- Integration tests for campaign operations
- E2E tests for complete campaign workflows
- Performance tests for large campaign lists
- Accessibility tests with screen readers
