import { describe, it, expect, beforeEach } from 'vitest'
import { createTestWrapper } from '../../test/utils/test-utils'
import MetricCard from '../common/MetricCard.vue'

describe('MetricCard', () => {
  const defaultProps = {
    title: 'Test Metric',
    value: '100',
    icon: 'trending_up',
    color: 'positive',
    trend: '+5%'
  }

  it('should render with default props', () => {
    const wrapper = createTestWrapper(MetricCard, {
      props: defaultProps
    })

    expect(wrapper.text()).toContain('Test Metric')
    expect(wrapper.text()).toContain('100')
    expect(wrapper.text()).toContain('+5%')
  })

  it('should render without trend when not provided', () => {
    const propsWithoutTrend = {
      ...defaultProps,
      trend: undefined
    }

    const wrapper = createTestWrapper(MetricCard, {
      props: propsWithoutTrend
    })

    expect(wrapper.text()).toContain('Test Metric')
    expect(wrapper.text()).toContain('100')
    expect(wrapper.text()).not.toContain('+5%')
  })

  it('should apply correct color classes', () => {
    const wrapper = createTestWrapper(MetricCard, {
      props: {
        ...defaultProps,
        color: 'negative'
      }
    })

    const card = wrapper.find('.metric-card')
    expect(card.classes()).toContain('metric-card--negative')
  })

  it('should display icon when provided', () => {
    const wrapper = createTestWrapper(MetricCard, {
      props: {
        ...defaultProps,
        icon: 'warning'
      }
    })

    const icon = wrapper.findComponent({ name: 'QIcon' })
    expect(icon.exists()).toBe(true)
    expect(icon.props('name')).toBe('warning')
  })

  it('should handle numeric values', () => {
    const wrapper = createTestWrapper(MetricCard, {
      props: {
        ...defaultProps,
        value: 42
      }
    })

    expect(wrapper.text()).toContain('42')
  })

  it('should handle string values', () => {
    const wrapper = createTestWrapper(MetricCard, {
      props: {
        ...defaultProps,
        value: 'Active'
      }
    })

    expect(wrapper.text()).toContain('Active')
  })

  it('should display trend with correct color', () => {
    const wrapper = createTestWrapper(MetricCard, {
      props: {
        ...defaultProps,
        trend: '-10%',
        color: 'negative'
      }
    })

    const trend = wrapper.find('.metric-trend')
    expect(trend.exists()).toBe(true)
    expect(trend.text()).toContain('-10%')
  })

  it('should handle empty title', () => {
    const wrapper = createTestWrapper(MetricCard, {
      props: {
        ...defaultProps,
        title: ''
      }
    })

    expect(wrapper.text()).toContain('100')
    expect(wrapper.text()).toContain('+5%')
  })

  it('should handle zero value', () => {
    const wrapper = createTestWrapper(MetricCard, {
      props: {
        ...defaultProps,
        value: 0
      }
    })

    expect(wrapper.text()).toContain('0')
  })

  it('should handle large numbers', () => {
    const wrapper = createTestWrapper(MetricCard, {
      props: {
        ...defaultProps,
        value: 1234567
      }
    })

    expect(wrapper.text()).toContain('1234567')
  })
})
